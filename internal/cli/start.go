package cli

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/openmuara/openmuara/internal/audit"
	_ "github.com/openmuara/openmuara/internal/billplz" // register billplz factory
	"github.com/openmuara/openmuara/internal/config"
	"github.com/openmuara/openmuara/internal/engine"
	_ "github.com/openmuara/openmuara/internal/fawry"  // register fawry factory
	_ "github.com/openmuara/openmuara/internal/ipay88" // register ipay88 factory
	"github.com/openmuara/openmuara/internal/plugin"
	"github.com/openmuara/openmuara/internal/provider"
	_ "github.com/openmuara/openmuara/internal/provider/defaultplugin" // register default provider
	"github.com/openmuara/openmuara/internal/provider/factory"
	"github.com/openmuara/openmuara/internal/server"
	_ "github.com/openmuara/openmuara/internal/stripe"    // register stripe factory
	_ "github.com/openmuara/openmuara/internal/toyyibpay" // register toyyibpay factory
	"github.com/openmuara/openmuara/internal/webhook"
	_ "modernc.org/sqlite" // register sqlite driver with database/sql
)

// parseLogLevel converts a level string to slog.Level.
func parseLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// providerBaseURL returns the external base URL used for provider payment links.
// It prefers server.public_base_url and falls back to the bind address.
func providerBaseURL(cfg *config.Config) string {
	if cfg.Server.PublicBaseURL != "" {
		return strings.TrimSuffix(cfg.Server.PublicBaseURL, "/")
	}
	return fmt.Sprintf("http://%s:%d", cfg.Server.Host, cfg.Server.Port)
}

// adminBaseURLFromConfig returns the external base URL used for the admin UI.
// It prefers server.admin_public_base_url, then public_base_url with admin_port,
// then the bind address with admin_port, then the provider base URL.
func adminBaseURLFromConfig(cfg *config.Config) string {
	base := fmt.Sprintf("http://%s:%d", cfg.Server.Host, cfg.Server.Port)
	if cfg.Server.AdminPort != 0 {
		base = fmt.Sprintf("http://%s:%d", cfg.Server.Host, cfg.Server.AdminPort)
	}
	if cfg.Server.AdminPublicBaseURL != "" {
		return strings.TrimSuffix(cfg.Server.AdminPublicBaseURL, "/")
	}
	if cfg.Server.PublicBaseURL != "" {
		base = strings.TrimSuffix(cfg.Server.PublicBaseURL, "/")
		if cfg.Server.AdminPort != 0 {
			base = fmt.Sprintf("%s:%d", base, cfg.Server.AdminPort)
		}
	}
	return base
}

// activeProviderName returns the first enabled loaded provider that exposes an
// escape handler, falling back to the first enabled provider. Defaults to
// "fawry" when no providers are configured, preserving backward compatibility.
func activeProviderName(cfg *config.Config, loaded []config.LoadedProvider) string {
	loadedByName := make(map[string]provider.Provider, len(loaded))
	for _, lp := range loaded {
		loadedByName[lp.Name] = lp.Provider
	}

	for name, pc := range cfg.Providers {
		if !pc.Enabled {
			continue
		}
		p, ok := loadedByName[name]
		if !ok {
			continue
		}
		if p.EscapeHandler() != nil {
			return name
		}
	}
	for name, pc := range cfg.Providers {
		if pc.Enabled {
			return name
		}
	}
	return "fawry"
}

// buildProviderDispatcher creates the webhook dispatcher for a provider from webhook config.
func buildProviderDispatcher(cfg config.WebhookConfig, providerName string, providerConfig map[string]any, p provider.Provider, store webhook.AttemptStore, ledger engine.TransactionStore) *webhook.Dispatcher {
	url := cfg.URL
	if target, ok := cfg.Targets[providerName]; ok && target != "" {
		url = target
	}
	if providerConfig != nil {
		if u, ok := providerConfig["webhook_url"].(string); ok && u != "" {
			url = u
		}
	}
	if url == "" {
		return nil
	}

	d := webhook.NewDispatcherFromProvider(url, cfg.MaxRetries, p)
	d.Store = store
	d.Ledger = ledger
	if providerConfig != nil {
		if events, ok := providerConfig["enabled_events"].([]any); ok {
			d.EnabledEvents = toStringSlice(events)
		}
	}
	return d
}

func providerWebhookURLConfigured(providers map[string]config.ProviderConfig) bool {
	for name, pc := range providers {
		if name != "stripe" || !pc.Enabled {
			continue
		}
		if u, ok := pc.Config["webhook_url"].(string); ok && u != "" {
			return true
		}
	}
	return false
}

func toStringSlice(v []any) []string {
	out := make([]string, 0, len(v))
	for _, item := range v {
		if s, ok := item.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

// newPersistenceStores creates the shared transaction, audit, and webhook stores from config.
// For SQLite, all stores share a single database connection to avoid locking.
func newPersistenceStores(cfg config.PersistenceConfig) (engine.TransactionStore, audit.Store, webhook.AttemptStore, func() error, error) {
	switch cfg.Type {
	case "memory", "":
		return engine.NewMemoryStore(), audit.NewMemoryStore(), webhook.NewMemoryStore(), func() error { return nil }, nil
	case "sqlite":
		path := cfg.Path
		if path == "" {
			path = ".muara/data/ledger.db"
		}
		if dir := filepath.Dir(path); dir != "" && dir != "." {
			if err := os.MkdirAll(dir, 0o750); err != nil {
				return nil, nil, nil, nil, fmt.Errorf("create data dir: %w", err)
			}
		}
		db, err := sql.Open("sqlite", fmt.Sprintf("%s?_busy_timeout=5000&_journal_mode=WAL", path))
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("open sqlite: %w", err)
		}
		if err := db.Ping(); err != nil {
			_ = db.Close()
			return nil, nil, nil, nil, fmt.Errorf("ping sqlite: %w", err)
		}
		// Serialize SQLite access: a single connection avoids writer-lock
		// contention between the HTTP handlers and the async webhook dispatcher.
		db.SetMaxOpenConns(1)
		db.SetMaxIdleConns(1)
		ledger, err := engine.NewSQLiteStoreFromDB(db)
		if err != nil {
			_ = db.Close()
			return nil, nil, nil, nil, fmt.Errorf("init transaction store: %w", err)
		}
		auditStore, err := audit.NewSQLiteStoreFromDB(db)
		if err != nil {
			_ = db.Close()
			return nil, nil, nil, nil, fmt.Errorf("init audit store: %w", err)
		}
		webhookStore, err := webhook.NewSQLiteStoreFromDB(db)
		if err != nil {
			_ = db.Close()
			return nil, nil, nil, nil, fmt.Errorf("init webhook store: %w", err)
		}
		return ledger, auditStore, webhookStore, db.Close, nil
	default:
		return nil, nil, nil, nil, fmt.Errorf("unsupported persistence type %q", cfg.Type)
	}
}

// startRuntime holds the fully wired dependencies needed to run the server.
type startRuntime struct {
	cfg         *config.Config
	ledger      engine.TransactionStore
	auditStore  audit.Store
	closeStores func() error
	enabled     []string
	dispatchers map[string]*webhook.Dispatcher
	auditLogger audit.Logger
	srv         *server.Server
	adminSrv    *server.Server
}

// newStartRuntime wires all dependencies from config. The registry and
// newServer factory are injectable so tests can use isolated providers or a
// test server.
func newStartRuntime(
	cfg *config.Config,
	registry *provider.Registry,
	newServer func(server.Config) *server.Server,
) (*startRuntime, error) {
	sharedLedger, auditStore, sharedWebhookStore, closeStores, err := newPersistenceStores(cfg.Persistence)
	if err != nil {
		return nil, err
	}

	systemAuditLogger := &audit.StoreLogger{Store: auditStore, Actor: "system"}

	loaded, err := config.LoadEnabledProvidersWithFallback(cfg, registry, plugin.LoadBuiltin)
	if err != nil {
		_ = closeStores()
		return nil, err
	}
	systemAuditLogger.Log(context.Background(), "config.reloaded", "config", cfg.Server.Host, "", "ok")

	enabled := make([]string, 0, len(loaded))
	providersByName := make(map[string]provider.Provider, len(loaded))
	for _, lp := range loaded {
		enabled = append(enabled, lp.Name)
		providersByName[lp.Name] = lp.Provider
		systemAuditLogger.Log(context.Background(), "provider.initialized", "provider", lp.Name, "", "ok")
	}

	availableSet := make(map[string]struct{})
	for _, n := range registry.Names() {
		availableSet[n] = struct{}{}
	}
	for _, n := range factory.Names() {
		availableSet[n] = struct{}{}
	}
	if builtinPlugins, err := plugin.LoadBuiltin("plugins", "../plugins", "../../plugins"); err == nil {
		for _, p := range builtinPlugins {
			availableSet[p.Name] = struct{}{}
		}
	}
	available := make([]string, 0, len(availableSet))
	for n := range availableSet {
		available = append(available, n)
	}
	sort.Strings(available)

	activeName := activeProviderName(cfg, loaded)
	baseURL := providerBaseURL(cfg)

	providerDispatchers := make(map[string]*webhook.Dispatcher, len(enabled))

	for _, lp := range loaded {
		p := lp.Provider
		if sp, ok := p.(interface{ SetStore(engine.TransactionStore) }); ok {
			sp.SetStore(sharedLedger)
		}
		if sp, ok := p.(interface{ SetBaseURL(string) }); ok {
			sp.SetBaseURL(baseURL)
		}
		if sp, ok := p.(interface{ SetConfigPath(string) }); ok {
			sp.SetConfigPath(rootConfigPath)
		}
		if sp, ok := p.(interface{ SetDispatcher(*webhook.Dispatcher) }); ok {
			if d := buildProviderDispatcher(cfg.Webhook, lp.Name, lp.Config.Config, p, sharedWebhookStore, sharedLedger); d != nil {
				d.AuditLogger = systemAuditLogger
				providerDispatchers[lp.Name] = d
				sp.SetDispatcher(d)
			} else {
				sp.SetDispatcher(nil)
			}
		}
	}

	dispatcher := providerDispatchers[activeName]

	if cfg.Dev.Seed {
		if err := seedDevData(sharedLedger, sharedWebhookStore, enabled); err != nil {
			slog.Warn("dev seed failed", "error", err)
		}
	}

	adminBaseURL := adminBaseURLFromConfig(cfg)

	routerCfg := server.RouterConfig{
		ActiveProvider:     activeName,
		EnabledProviders:   enabled,
		AvailableProviders: available,
		Providers:          providersByName,
		Dispatcher:         dispatcher,
		Dispatchers:        providerDispatchers,
		TransactionStore:   sharedLedger,
		AuditStore:         auditStore,
		CORS:               cfg.Server.CORS,
		CSRF:               cfg.Server.CSRF,
		Pprof:              cfg.Server.Pprof,
		ConfigPath:         rootConfigPath,
		Host:               cfg.Server.Host,
		Port:               cfg.Server.Port,
		PublicBaseURL:      cfg.Server.PublicBaseURL,
		AdminBaseURL:       adminBaseURL,
		Auth: server.AuthConfig{
			Enabled:      cfg.Admin.Enabled,
			Username:     cfg.Admin.Username,
			PasswordHash: cfg.Admin.PasswordHash,
			Token:        cfg.Admin.Token,
			Viewer: server.ViewerAuthConfig{
				Enabled:      cfg.Viewer.Enabled,
				Username:     cfg.Viewer.Username,
				PasswordHash: cfg.Viewer.PasswordHash,
				Token:        cfg.Viewer.Token,
			},
		},
		Hardened: cfg.Hardened,
		RateLimit: server.RateLimiterConfig{
			Enabled:           cfg.RateLimit.Enabled || cfg.Hardened,
			RequestsPerMinute: cfg.RateLimit.RequestsPerMinute,
			AdminOnly:         !cfg.Hardened,
		},
		SecurityHeaders: server.SecurityHeadersConfig{
			Enabled: cfg.Hardened,
			TLS:     cfg.Server.TLSCert != "" && cfg.Server.TLSKey != "",
		},
	}

	var srv, adminSrv *server.Server
	if cfg.Server.AdminPort != 0 {
		// Dual-port mode: provider port serves emulation + public endpoints;
		// admin port serves the dashboard UI and JSON APIs.
		routerCfg.AdminBaseURL = adminBaseURL
		providerRouter := server.NewProviderRouter(routerCfg)
		adminRouter := server.NewAdminRouter(routerCfg)

		srv = newServer(server.Config{
			Host:    cfg.Server.Host,
			Port:    cfg.Server.Port,
			TLSCert: cfg.Server.TLSCert,
			TLSKey:  cfg.Server.TLSKey,
			Handler: providerRouter,
		})
		adminSrv = newServer(server.Config{
			Host:    cfg.Server.Host,
			Port:    cfg.Server.AdminPort,
			TLSCert: cfg.Server.TLSCert,
			TLSKey:  cfg.Server.TLSKey,
			Handler: adminRouter,
		})
	} else {
		// Single-port mode: preserve backward-compatible behavior.
		router := server.NewRouter(routerCfg)
		srv = newServer(server.Config{
			Host:    cfg.Server.Host,
			Port:    cfg.Server.Port,
			TLSCert: cfg.Server.TLSCert,
			TLSKey:  cfg.Server.TLSKey,
			Handler: router,
		})
	}

	return &startRuntime{
		cfg:         cfg,
		ledger:      sharedLedger,
		auditStore:  auditStore,
		closeStores: closeStores,
		enabled:     enabled,
		dispatchers: providerDispatchers,
		auditLogger: systemAuditLogger,
		srv:         srv,
		adminSrv:    adminSrv,
	}, nil
}

// muaraBanner is printed on server startup unless disabled.
const muaraBanner = `≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋
≋≋  muara — local payments, flowing smoothly   ≋≋
≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋≋
`

func shouldSkipBanner(noBannerFlag bool) bool {
	if noBannerFlag {
		return true
	}
	return os.Getenv("MUARA_NO_BANNER") == "1" || os.Getenv("MUARA_NO_BANNER") == "true"
}

// runStart starts the server(s) and blocks until the context is cancelled.
func runStart(ctx context.Context, rt *startRuntime, out io.Writer) error {
	defer func() { _ = rt.closeStores() }()
	_, _ = fmt.Fprintln(out, "starting muara...")

	if rt.adminSrv != nil {
		_, _ = fmt.Fprintf(out, "provider API: %s\n", rt.srv.BaseURL())
		_, _ = fmt.Fprintf(out, "admin API:    %s\n", rt.adminSrv.BaseURL())

		srvCtx, cancel := context.WithCancel(ctx)
		defer cancel()

		errCh := make(chan error, 2)
		go func() { errCh <- rt.srv.ListenAndServe(srvCtx) }()
		go func() { errCh <- rt.adminSrv.ListenAndServe(srvCtx) }()

		var firstErr error
		for i := 0; i < 2; i++ {
			if err := <-errCh; err != nil && firstErr == nil {
				firstErr = err
				cancel()
			}
		}
		return firstErr
	}

	_, _ = fmt.Fprintf(out, "server: %s\n", rt.srv.BaseURL())
	return rt.srv.ListenAndServe(ctx)
}

func newStartCommand() *cobra.Command {
	var dryRun bool
	var noBanner bool

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the muara server",
		Example: `  muara start
  muara start --config path/to/config.yml
  muara start --dry-run`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := config.Load(rootConfigPath)
			if err != nil {
				return fmt.Errorf("load config: %w\n\nDid you mean to run 'muara init' first?", err)
			}
			if valErrs := cfg.ValidateWithDetails(rootConfigPath); len(valErrs) > 0 {
				_, _ = fmt.Fprint(cmd.OutOrStderr(), config.FormatValidationErrors(valErrs))
				_, _ = fmt.Fprintln(cmd.OutOrStderr(), "\nRun 'muara doctor' for a full health check.")
				return fmt.Errorf("invalid config")
			}

			if dryRun {
				if jsonOutput {
					return respondJSON(cmd.OutOrStdout(), startDryRunResult{OK: true, Path: rootConfigPath})
				}
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "config loaded and validated: %s\n", rootConfigPath)
				return nil
			}

			if !shouldSkipBanner(noBanner) {
				_, _ = fmt.Fprint(cmd.OutOrStdout(), muaraBanner)
			}

			opts := &slog.HandlerOptions{Level: parseLogLevel(cfg.Log.Level)}
			logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
			slog.SetDefault(logger)

			if warning := config.ValidateWebhookURL(cfg.Webhook); warning != "" {
				logger.Warn(warning)
			}

			rt, err := newStartRuntime(cfg, provider.Default(), server.New)
			if err != nil {
				return err
			}

			if cfg.Server.TLSCert != "" && cfg.Server.TLSKey != "" {
				rt.auditLogger.Log(cmd.Context(), server.SecurityEventTLSEnabled, "tls", "enabled", "", "ok")
			}

			return runStart(cmd.Context(), rt, cmd.OutOrStdout())
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "load config and validate without starting the server")
	cmd.Flags().BoolVar(&noBanner, "no-banner", false, "skip the startup banner (also MUARA_NO_BANNER=1)")
	return cmd
}

// startDryRunResult is the structured output for muara start --dry-run.
type startDryRunResult struct {
	OK   bool   `json:"ok"`
	Path string `json:"path"`
}
