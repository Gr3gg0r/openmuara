package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/pprof"

	"github.com/Gr3gg0r/openmuara/internal/api"
	"github.com/Gr3gg0r/openmuara/internal/audit"
	"github.com/Gr3gg0r/openmuara/internal/config"
	"github.com/Gr3gg0r/openmuara/internal/engine"
	"github.com/Gr3gg0r/openmuara/internal/httputil"
	"github.com/Gr3gg0r/openmuara/internal/provider"
	"github.com/Gr3gg0r/openmuara/internal/ui"
	"github.com/Gr3gg0r/openmuara/internal/webhook"
)

// RouterConfig holds the dependencies for building the HTTP router.
type RouterConfig struct {
	ActiveProvider     string
	EnabledProviders   []string
	AvailableProviders []string
	Providers          map[string]provider.Provider
	Host               string
	Port               int
	PublicBaseURL      string
	Dispatcher         *webhook.Dispatcher
	Dispatchers        map[string]*webhook.Dispatcher
	TransactionStore   engine.TransactionStore
	AuditStore         audit.Store
	CORS               config.CORSConfig
	CSRF               config.CSRFConfig
	Pprof              bool
	Auth               AuthConfig
	Hardened           bool
	RateLimit          RateLimiterConfig
	RateLimiter        *RateLimiter
	SecurityHeaders    SecurityHeadersConfig
	ConfigPath         string
	AdminBaseURL       string
}

// NewRouter builds the application router with provider-driven route registration.
// This is the single-port convenience that registers all routes on one mux.
func NewRouter(cfg RouterConfig) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthzHandler)
	mux.HandleFunc("GET /readyz", readyzHandler(cfg.EnabledProviders, cfg.AvailableProviders))
	mux.HandleFunc("GET /openapi.yaml", OpenAPIHandler())
	mux.Handle("GET /metrics", MetricsHandler())
	mux.HandleFunc("GET /__muara/payment-pages.css", paymentPagesCSSHandler)
	mux.HandleFunc("GET /{$}", dashboardHandler(cfg.ActiveProvider, cfg.AdminBaseURL))
	mux.HandleFunc("GET /_admin/{$}", dashboardHandler(cfg.ActiveProvider, cfg.AdminBaseURL))
	mux.Handle("GET /dashboard-assets/", http.StripPrefix("/dashboard-assets/", http.FileServer(ui.DashboardAssetsFS())))

	// Universal payment API.
	if cfg.TransactionStore != nil {
		mux.HandleFunc("POST /v1/pay", api.NewPayHandler(cfg.TransactionStore, cfg.EnabledProviders))
		mux.HandleFunc("GET /v1/pay/", api.NewGetPaymentHandler(cfg.TransactionStore))
		mux.HandleFunc("POST /v1/refund/", api.NewRefundHandler(cfg.TransactionStore))
	}

	// Register routes for all enabled providers from the loaded provider set.
	// Mutating admin routes are wrapped with requireAdmin so viewer credentials
	// cannot simulate payments or change provider state on either port.
	for _, name := range cfg.EnabledProviders {
		p, ok := cfg.Providers[name]
		if !ok {
			if gp, err := provider.Get(name); err == nil {
				p = gp
			} else {
				slog.Error("enabled provider not loaded", "name", name)
				continue
			}
		}
		for _, route := range p.Routes() {
			pattern := route.Method + " " + route.Path
			h := route.Handler
			if isMutatingAdminRoute(route.Method, route.Path) {
				h = requireAdmin(h)
			}
			mux.Handle(pattern, h)
			slog.Info("registered route", "method", route.Method, "path", route.Path, "provider", name)
		}
	}

	// Register Stripe webhook configuration UI if the provider exposes it.
	if p, ok := cfg.Providers["stripe"]; ok {
		if hp, ok := p.(interface {
			WebhookConfigPageHandler() http.Handler
			WebhookConfigSaveHandler() http.Handler
		}); ok {
			mux.Handle("GET /_admin/stripe/webhooks", hp.WebhookConfigPageHandler())
			mux.Handle("POST /_admin/stripe/webhooks", requireAdmin(hp.WebhookConfigSaveHandler()))
		}
	}

	mux.HandleFunc("POST /_admin/webhook-receiver", webhook.NewTestReceiverHandler())
	WebhookAdminHandlers(mux, cfg.Dispatcher, cfg.Dispatchers)
	ScenarioAdminHandlers(mux, cfg.TransactionStore)
	AdminAPIHandlers(mux, cfg)
	ConfigAdminHandlers(mux, cfg)
	AuditAdminHandlers(mux, cfg.AuditStore)

	if cfg.Pprof {
		mux.HandleFunc("GET /_admin/debug/pprof/", pprof.Index)
		mux.HandleFunc("GET /_admin/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("GET /_admin/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("GET /_admin/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("GET /_admin/debug/pprof/trace", pprof.Trace)
	}

	return wrapWithMiddleware(mux, cfg)
}

// NewProviderRouter builds the router for provider emulation and public health/metrics endpoints.
// It does not register admin JSON API routes except for the dashboard HTML fallback and
// the webhook test receiver needed by provider flows.
func NewProviderRouter(cfg RouterConfig) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthzHandler)
	mux.HandleFunc("GET /readyz", readyzHandler(cfg.EnabledProviders, cfg.AvailableProviders))
	mux.HandleFunc("GET /openapi.yaml", OpenAPIHandler())
	mux.Handle("GET /metrics", MetricsHandler())
	mux.HandleFunc("GET /__muara/payment-pages.css", paymentPagesCSSHandler)
	mux.HandleFunc("GET /{$}", dashboardHandler(cfg.ActiveProvider, cfg.AdminBaseURL))
	mux.HandleFunc("GET /_admin/{$}", dashboardHandler(cfg.ActiveProvider, cfg.AdminBaseURL))
	mux.Handle("GET /dashboard-assets/", http.StripPrefix("/dashboard-assets/", http.FileServer(ui.DashboardAssetsFS())))

	// Universal payment API.
	if cfg.TransactionStore != nil {
		mux.HandleFunc("POST /v1/pay", api.NewPayHandler(cfg.TransactionStore, cfg.EnabledProviders))
		mux.HandleFunc("GET /v1/pay/", api.NewGetPaymentHandler(cfg.TransactionStore))
		mux.HandleFunc("POST /v1/refund/", api.NewRefundHandler(cfg.TransactionStore))
	}

	// Register routes for all enabled providers from the loaded provider set.
	// Mutating admin routes are wrapped with requireAdmin so viewer credentials
	// cannot simulate payments or change provider state on either port.
	for _, name := range cfg.EnabledProviders {
		p, ok := cfg.Providers[name]
		if !ok {
			if gp, err := provider.Get(name); err == nil {
				p = gp
			} else {
				slog.Error("enabled provider not loaded", "name", name)
				continue
			}
		}
		for _, route := range p.Routes() {
			pattern := route.Method + " " + route.Path
			h := route.Handler
			if isMutatingAdminRoute(route.Method, route.Path) {
				h = requireAdmin(h)
			}
			mux.Handle(pattern, h)
			slog.Info("registered route", "method", route.Method, "path", route.Path, "provider", name)
		}
	}

	// Register Stripe webhook configuration UI if the provider exposes it.
	if p, ok := cfg.Providers["stripe"]; ok {
		if hp, ok := p.(interface {
			WebhookConfigPageHandler() http.Handler
			WebhookConfigSaveHandler() http.Handler
		}); ok {
			mux.Handle("GET /_admin/stripe/webhooks", hp.WebhookConfigPageHandler())
			mux.Handle("POST /_admin/stripe/webhooks", requireAdmin(hp.WebhookConfigSaveHandler()))
		}
	}

	// Webhook test receiver used by provider flows and the admin dashboard.
	mux.HandleFunc("POST /_admin/webhook-receiver", webhook.NewTestReceiverHandler())

	return wrapWithMiddleware(mux, cfg)
}

// NewAdminRouter builds the router for the admin web UI and JSON API endpoints.
func NewAdminRouter(cfg RouterConfig) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", dashboardHandler(cfg.ActiveProvider, cfg.AdminBaseURL))
	mux.HandleFunc("GET /_admin/{$}", dashboardHandler(cfg.ActiveProvider, cfg.AdminBaseURL))
	mux.Handle("GET /dashboard-assets/", http.StripPrefix("/dashboard-assets/", http.FileServer(ui.DashboardAssetsFS())))

	mux.HandleFunc("POST /_admin/webhook-receiver", webhook.NewTestReceiverHandler())
	WebhookAdminHandlers(mux, cfg.Dispatcher, cfg.Dispatchers)
	ScenarioAdminHandlers(mux, cfg.TransactionStore)
	AdminAPIHandlers(mux, cfg)
	ConfigAdminHandlers(mux, cfg)
	AuditAdminHandlers(mux, cfg.AuditStore)

	if cfg.Pprof {
		mux.HandleFunc("GET /_admin/debug/pprof/", pprof.Index)
		mux.HandleFunc("GET /_admin/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("GET /_admin/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("GET /_admin/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("GET /_admin/debug/pprof/trace", pprof.Trace)
	}

	return wrapWithMiddleware(mux, cfg)
}

// wrapWithMiddleware applies the standard server middleware stack to a handler.
func wrapWithMiddleware(handler http.Handler, cfg RouterConfig) http.Handler {
	var auditLogger audit.Logger = &audit.StoreLogger{Store: cfg.AuditStore, Actor: "http"}

	handler = MaxBodySizeMiddleware(handler)
	handler = AuthMiddleware(cfg.Auth)(handler)

	limiter := cfg.RateLimiter
	if limiter == nil {
		limiter = NewRateLimiter(cfg.RateLimit.RequestsPerMinute)
	}
	handler = RateLimitMiddleware(cfg.RateLimit, limiter)(handler)

	handler = SecurityHeadersMiddleware(cfg.SecurityHeaders)(handler)
	handler = MetricsMiddleware(handler)
	handler = CSRFGuardMiddleware(CSRFGuardConfig{
		Enabled:        cfg.CSRF.Enabled,
		SecureCookie:   cfg.SecurityHeaders.TLS,
		SameSiteStrict: cfg.Auth.IsAuthEnabled(),
	})(handler)
	handler = CORSMiddleware(cfg.CORS)(handler)
	handler = AuditMiddleware(auditLogger)(handler)
	handler = LoggingMiddleware(handler)
	handler = RequestIDMiddleware(handler)

	return handler
}

// requireAdmin wraps a handler so it returns 403 unless the request context
// carries RoleAdmin. It is used at route-registration time for mutating admin
// endpoints that are registered outside the server package.
func requireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !RequireAdmin(r.Context(), w, r) {
			return
		}
		next.ServeHTTP(w, r)
	})
}

// isMutatingAdminRoute reports whether a provider route should be wrapped with
// requireAdmin. It always returns false because provider pay/escape pages are
// part of the customer payment flow and must stay reachable without admin
// elevation. Dashboard mutations are wrapped explicitly in router registration.
func isMutatingAdminRoute(_, _ string) bool {
	return false
}

func healthzHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func readyzHandler(enabledProviders, availableProviders []string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":    "ok",
			"providers": enabledProviders,
			"available": availableProviders,
		})
	}
}

func paymentPagesCSSHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	_, _ = w.Write(ui.PaymentPagesCSS())
}

func dashboardHandler(activeProvider, adminBaseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := ui.DashboardData{
			ActiveProvider:  activeProvider,
			AdminAPIBaseURL: adminBaseURL,
			Role:            string(RoleFromContext(r.Context())),
		}
		if tok, ok := httputil.CSRFTokenFromContext(r.Context()); ok {
			data.CSRFToken = tok
		}
		_ = ui.ServeDashboard(w, data)
	}
}
