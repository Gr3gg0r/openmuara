package server

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/openmuara/openmuara/internal/audit"
	"github.com/openmuara/openmuara/internal/config"
	"github.com/openmuara/openmuara/internal/provider"
	"github.com/openmuara/openmuara/internal/provider/factory"
	"gopkg.in/yaml.v3"
)

// ConfigAdminHandlers registers dashboard endpoints for reading and writing
// muara configuration. Changes are persisted to disk but only take effect after
// a server restart (provider enablement) or are immediately reflected by
// dispatchers that read the in-memory config on each request.
func ConfigAdminHandlers(mux *http.ServeMux, cfg RouterConfig) {
	if cfg.ConfigPath == "" {
		return
	}

	mux.HandleFunc("GET /_admin/config", getConfigHandler(cfg.ConfigPath))
	mux.HandleFunc("PATCH /_admin/config/providers", patchProvidersHandler(cfg.ConfigPath))
	mux.HandleFunc("GET /_admin/config/webhooks", getWebhooksHandler(cfg.ConfigPath))
	mux.HandleFunc("PATCH /_admin/config/webhooks", patchWebhooksHandler(cfg.ConfigPath))
	mux.HandleFunc("POST /_admin/config/webhooks/test", testWebhookHandler(cfg.Hardened))
	mux.HandleFunc("POST /_admin/config/reload", reloadConfigHandler())
}

// safeConfig is the subset of configuration exposed to the dashboard.
type safeConfig struct {
	Server         safeServerConfig              `json:"server"`
	Providers      map[string]safeProviderConfig `json:"providers"`
	Webhook        safeWebhookConfig             `json:"webhook"`
	ConfigChecksum string                        `json:"config_checksum"`
}

type safeServerConfig struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	AdminPort int    `json:"admin_port"`
}

type safeProviderConfig struct {
	Enabled bool           `json:"enabled"`
	Config  map[string]any `json:"config"`
}

type safeWebhookConfig struct {
	URL        string              `json:"url"`
	MaxRetries int                 `json:"max_retries"`
	Targets    map[string]string   `json:"targets"`
	Events     map[string][]string `json:"events"`
}

func getConfigHandler(configPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		// #nosec G304 -- path is the configured muara config file, not user-supplied.
		data, err := os.ReadFile(configPath)
		if err != nil && !os.IsNotExist(err) {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		cfg, err := config.Load(configPath)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		respondJSON(w, http.StatusOK, safeConfig{
			Server: safeServerConfig{
				Host:      cfg.Server.Host,
				Port:      cfg.Server.Port,
				AdminPort: cfg.Server.AdminPort,
			},
			Providers:      redactedProviders(cfg.Providers),
			Webhook:        safeWebhookConfigFrom(cfg.Webhook),
			ConfigChecksum: sha256Hex(data),
		})
	}
}

func safeWebhookConfigFrom(cfg config.WebhookConfig) safeWebhookConfig {
	return safeWebhookConfig{
		URL:        cfg.URL,
		MaxRetries: cfg.MaxRetries,
		Targets:    cfg.Targets,
		Events:     cfg.Events,
	}
}

func redactedProviders(providers map[string]config.ProviderConfig) map[string]safeProviderConfig {
	out := make(map[string]safeProviderConfig, len(providers))
	for name, pc := range providers {
		out[name] = safeProviderConfig{
			Enabled: pc.Enabled,
			Config:  redactedConfig(pc.Config),
		}
	}
	return out
}

func redactedConfig(cfg map[string]any) map[string]any {
	out := make(map[string]any, len(cfg))
	for k, v := range cfg {
		lower := strings.ToLower(k)
		if strings.Contains(lower, "secret") ||
			strings.Contains(lower, "key") ||
			strings.Contains(lower, "token") ||
			strings.Contains(lower, "password") {
			out[k] = "***"
			continue
		}
		out[k] = v
	}
	return out
}

type patchProvidersRequest struct {
	Providers map[string]struct {
		Enabled *bool `json:"enabled"`
	} `json:"providers"`
	Checksum *string `json:"checksum"`
}

func patchProvidersHandler(configPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !RequireAdmin(r.Context(), w, r) {
			return
		}
		var req patchProvidersRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		expectedChecksum := expectedChecksumFrom(r, req.Checksum)
		if err := patchConfig(configPath, expectedChecksum, func(raw map[string]any) error {
			providersRaw, ok := raw["providers"].(map[string]any)
			if !ok {
				providersRaw = make(map[string]any)
				raw["providers"] = providersRaw
			}
			for name, update := range req.Providers {
				if update.Enabled == nil {
					continue
				}
				if _, registered := config.WizardChoiceByKey(name); !registered && name != "default" {
					if _, hasFactory := factory.Get(name); !hasFactory {
						if _, err := provider.Get(name); err != nil {
							return fmt.Errorf("provider %q is not registered", name)
						}
					}
				}
				pcRaw, ok := providersRaw[name].(map[string]any)
				if !ok {
					pcRaw = make(map[string]any)
					providersRaw[name] = pcRaw
				}
				pcRaw["enabled"] = *update.Enabled
			}
			return nil
		}); err != nil {
			if errors.Is(err, errConfigConflict) {
				respondJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
				return
			}
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		audit.FromContext(r.Context()).Log(r.Context(), "admin.config_providers_update", "config", configPath, "", "updated")
		respondJSON(w, http.StatusAccepted, map[string]string{"status": "saved", "note": "restart required for provider changes to take effect"})
	}
}

func getWebhooksHandler(configPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		cfg, err := config.Load(configPath)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		respondJSON(w, http.StatusOK, safeWebhookConfigFrom(cfg.Webhook))
	}
}

type patchWebhooksRequest struct {
	URL        *string             `json:"url"`
	MaxRetries *int                `json:"max_retries"`
	Targets    map[string]string   `json:"targets"`
	Events     map[string][]string `json:"events"`
	Checksum   *string             `json:"checksum"`
}

func patchWebhooksHandler(configPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !RequireAdmin(r.Context(), w, r) {
			return
		}
		var req patchWebhooksRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		expectedChecksum := expectedChecksumFrom(r, req.Checksum)
		if err := patchConfig(configPath, expectedChecksum, func(raw map[string]any) error {
			webhookRaw, ok := raw["webhook"].(map[string]any)
			if !ok {
				webhookRaw = make(map[string]any)
				raw["webhook"] = webhookRaw
			}
			if req.URL != nil {
				if *req.URL != "" {
					if u, err := url.Parse(*req.URL); err != nil || (u.Scheme != "http" && u.Scheme != "https") {
						return fmt.Errorf("invalid webhook url")
					}
				}
				webhookRaw["url"] = *req.URL
			}
			if req.MaxRetries != nil {
				if *req.MaxRetries < 0 {
					return fmt.Errorf("max_retries must be >= 0")
				}
				webhookRaw["max_retries"] = *req.MaxRetries
			}
			if req.Targets != nil {
				webhookRaw["targets"] = req.Targets
			}
			if req.Events != nil {
				webhookRaw["events"] = req.Events
			}
			return nil
		}); err != nil {
			if errors.Is(err, errConfigConflict) {
				respondJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
				return
			}
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		audit.FromContext(r.Context()).Log(r.Context(), "admin.config_webhooks_update", "config", configPath, "", "updated")
		respondJSON(w, http.StatusAccepted, map[string]string{"status": "saved", "note": "restart required for webhook target changes to take effect"})
	}
}

func expectedChecksumFrom(r *http.Request, bodyChecksum *string) string {
	if etag := r.Header.Get("If-Match"); etag != "" {
		return strings.Trim(etag, "\"")
	}
	if bodyChecksum != nil {
		return *bodyChecksum
	}
	return ""
}

type testWebhookRequest struct {
	URL      string `json:"url"`
	Provider string `json:"provider"`
	Secret   string `json:"secret"`
}

type testWebhookResponse struct {
	Success           bool   `json:"success"`
	Status            int    `json:"status"`
	Error             string `json:"error,omitempty"`
	LatencyMs         int64  `json:"latency_ms"`
	SignatureVerified bool   `json:"signature_verified"`
}

func testWebhookHandler(hardened bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !RequireAdmin(r.Context(), w, r) {
			return
		}
		var req testWebhookRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if req.URL == "" {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "url is required"})
			return
		}
		u, err := url.Parse(req.URL)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid url"})
			return
		}

		if hardened {
			if isPrivateOrReserved(u.Hostname()) {
				respondJSON(w, http.StatusBadRequest, testWebhookResponse{
					Success: false,
					Error:   "url resolves to a private or reserved address",
				})
				return
			}
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		timestamp := time.Now().UTC().Format(time.RFC3339)
		payload, err := json.Marshal(map[string]string{
			"event":     "muara.test_event",
			"provider":  req.Provider,
			"timestamp": timestamp,
		})
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, testWebhookResponse{Success: false, Error: err.Error()})
			return
		}

		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, req.URL, strings.NewReader(string(payload)))
		if err != nil {
			respondJSON(w, http.StatusOK, testWebhookResponse{Success: false, Error: err.Error()})
			return
		}
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("User-Agent", "openmuara-webhook-test/1.0")

		if req.Secret != "" {
			mac := hmac.New(sha256.New, []byte(req.Secret))
			mac.Write(payload)
			httpReq.Header.Set("X-Muara-Signature", hex.EncodeToString(mac.Sum(nil)))
		}

		start := time.Now()
		resp, err := http.DefaultClient.Do(httpReq)
		latencyMs := time.Since(start).Milliseconds()
		if err != nil {
			respondJSON(w, http.StatusOK, testWebhookResponse{Success: false, Error: err.Error(), LatencyMs: latencyMs})
			return
		}
		defer func() { _ = resp.Body.Close() }()
		_, _ = io.Copy(io.Discard, resp.Body)

		result := testWebhookResponse{
			Success:           resp.StatusCode >= 200 && resp.StatusCode < 300,
			Status:            resp.StatusCode,
			LatencyMs:         latencyMs,
			SignatureVerified: req.Secret != "",
		}
		if !result.Success {
			result.Error = fmt.Sprintf("received status %d", resp.StatusCode)
		}
		respondJSON(w, http.StatusOK, result)
	}
}

func reloadConfigHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !RequireAdmin(r.Context(), w, r) {
			return
		}
		audit.FromContext(r.Context()).Log(r.Context(), "admin.config_reload_signal", "config", "*", "", "ok")
		respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

var errConfigConflict = errors.New("config conflict: file changed on disk")

// patchConfig reads the YAML config at path, applies fn, validates the result,
// and writes it back atomically.
func patchConfig(path string, expectedChecksum string, fn func(map[string]any) error) error {
	// #nosec G304 -- path is the configured muara config file, not user-supplied.
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("read config: %w", err)
	}

	raw := make(map[string]any)
	if len(data) > 0 {
		if err := yaml.Unmarshal(data, &raw); err != nil {
			return fmt.Errorf("parse config: %w", err)
		}
	}

	if err := fn(raw); err != nil {
		return err
	}

	out, err := yaml.Marshal(raw)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := validateCandidateConfig(out); err != nil {
		return fmt.Errorf("validation: %w", err)
	}

	currentChecksum := sha256Hex(data)
	if expectedChecksum != "" && expectedChecksum != currentChecksum {
		return fmt.Errorf("expected checksum %s, got %s: %w", expectedChecksum, currentChecksum, errConfigConflict)
	}

	if len(data) > 0 {
		// #nosec G703 -- path is the configured muara config file, not user-supplied.
		if err := os.WriteFile(path+".bak", data, 0o600); err != nil {
			return fmt.Errorf("backup config: %w", err)
		}
	}

	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, "config-*.yml")
	if err != nil {
		return fmt.Errorf("create temp config: %w", err)
	}
	tmpPath := tmp.Name()
	if _, err := tmp.Write(out); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpPath)
		return fmt.Errorf("write temp config: %w", err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("close temp config: %w", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("rename config: %w", err)
	}
	return nil
}

func validateCandidateConfig(data []byte) error {
	cfg, err := config.LoadFromBytes(data)
	if err != nil {
		return err
	}
	return cfg.Validate()
}

func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// isPrivateOrReserved reports whether host resolves to or is a loopback,
// link-local, or private address.
func isPrivateOrReserved(host string) bool {
	if host == "" || host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	if ip != nil {
		return ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast()
	}
	ips, err := net.LookupIP(host)
	if err != nil {
		return false
	}
	for _, ip := range ips {
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
			return true
		}
	}
	return false
}
