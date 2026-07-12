package server

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/openmuara/openmuara/internal/config"
	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/httputil"
	"github.com/openmuara/openmuara/internal/provider"
	"github.com/openmuara/openmuara/internal/provider/factory"
	"github.com/openmuara/openmuara/internal/webhook"
)

// AdminAPIHandlers registers dashboard JSON endpoints used by the web UI.
func AdminAPIHandlers(mux *http.ServeMux, cfg RouterConfig) {
	if cfg.TransactionStore != nil {
		mux.HandleFunc("GET /_admin/transactions", listTransactionsHandler(cfg.TransactionStore))
		mux.HandleFunc("GET /_admin/transactions/{ref}", getTransactionHandler(cfg.TransactionStore))
		mux.HandleFunc("POST /_admin/transactions/{ref}/replay-webhook", replayTransactionWebhookHandler(cfg.TransactionStore, cfg.Dispatcher, cfg.Dispatchers))
		mux.HandleFunc("GET /_admin/ledger", ledgerHandler(cfg.TransactionStore, cfg.Dispatcher))
		mux.HandleFunc("POST /_admin/clean", cleanHandler(cfg))
	}
	mux.HandleFunc("GET /_admin/providers", listProvidersHandler(cfg))
	mux.HandleFunc("GET /_admin/providers/{name}", getProviderHandler(cfg))
	mux.HandleFunc("GET /_admin/providers/{name}/health", providerHealthHandler(cfg.ConfigPath))
	if cfg.TransactionStore != nil {
		mux.HandleFunc("GET /_admin/onboarding", onboardingHandler(cfg.ActiveProvider, cfg.EnabledProviders, cfg.TransactionStore, cfg.Dispatcher))
	}
}

func cleanHandler(cfg RouterConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !RequireAdmin(r.Context(), w, r) {
			return
		}

		var errs []string
		if err := cfg.TransactionStore.Clear(); err != nil {
			errs = append(errs, fmt.Sprintf("transactions: %v", err))
		}
		if cfg.AuditStore != nil {
			if err := cfg.AuditStore.Clear(); err != nil {
				errs = append(errs, fmt.Sprintf("audit: %v", err))
			}
		}
		if cfg.Dispatcher != nil && cfg.Dispatcher.Store != nil {
			if err := cfg.Dispatcher.Store.Clear(); err != nil {
				errs = append(errs, fmt.Sprintf("webhooks: %v", err))
			}
		}

		if len(errs) > 0 {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": strings.Join(errs, "; ")})
			return
		}

		respondJSON(w, http.StatusOK, map[string]bool{"ok": true})
	}
}

func listTransactionsHandler(store engine.TransactionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, offset := pageParams(r)
		q := strings.ToLower(r.URL.Query().Get("q"))
		providerFilter := r.URL.Query().Get("provider")
		statusFilter := r.URL.Query().Get("status")

		all, err := store.List(-1, 0)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		filtered := make([]engine.Transaction, 0, len(all))
		for _, tx := range all {
			if !transactionMatches(tx, q, providerFilter, statusFilter) {
				continue
			}
			filtered = append(filtered, tx)
		}

		total := len(filtered)
		if offset < 0 {
			offset = 0
		}
		end := offset + limit
		if limit <= 0 || end > total {
			end = total
		}
		if offset > total {
			offset = total
		}
		page := filtered[offset:end]

		respondJSON(w, http.StatusOK, map[string]any{
			"limit":   limit,
			"offset":  offset,
			"total":   total,
			"results": page,
		})
	}
}

func transactionMatches(tx engine.Transaction, q, providerFilter, statusFilter string) bool {
	if providerFilter != "" && tx.Provider != providerFilter {
		return false
	}
	if statusFilter != "" && string(tx.Status) != statusFilter {
		return false
	}
	if q == "" {
		return true
	}
	return strings.Contains(strings.ToLower(tx.Reference), q) ||
		strings.Contains(strings.ToLower(tx.Provider), q) ||
		strings.Contains(strings.ToLower(string(tx.Status)), q)
}

func getTransactionHandler(store engine.TransactionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ref := r.PathValue("ref")
		tx, ok, err := store.GetByReference(ref)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		if !ok {
			respondJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}

		respondJSON(w, http.StatusOK, map[string]any{
			"transaction": tx,
			"history":     []any{}, // Timeline not yet persisted; reserved for future enhancement.
		})
	}
}

func replayTransactionWebhookHandler(store engine.TransactionStore, active *webhook.Dispatcher, dispatchers map[string]*webhook.Dispatcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !RequireAdmin(r.Context(), w, r) {
			return
		}
		ref := r.PathValue("ref")
		tx, ok, err := store.GetByReference(ref)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		if !ok {
			respondJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}

		d := active
		if tx.Provider != "" {
			if pd, ok := dispatchers[tx.Provider]; ok && pd != nil {
				d = pd
			}
		}
		if d == nil {
			respondJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "webhook dispatcher not configured"})
			return
		}

		ctx := httputil.WithTraceID(r.Context(), tx.TraceID)
		attempt, err := d.Dispatch(ctx, ref, webhook.PaymentStatus(tx.Status))
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		respondJSON(w, http.StatusAccepted, attempt)
	}
}

func ledgerHandler(store engine.TransactionStore, dispatcher *webhook.Dispatcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, offset := pageParams(r)
		if limit <= 0 {
			limit = 50
		}
		q := strings.ToLower(r.URL.Query().Get("q"))
		typeFilter := r.URL.Query().Get("type")
		providerFilter := r.URL.Query().Get("provider")
		statusFilter := r.URL.Query().Get("status")

		txs, err := store.List(-1, 0)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		var attempts []*webhook.Attempt
		if dispatcher != nil && dispatcher.Store != nil {
			attempts, err = dispatcher.Store.List(-1, 0)
			if err != nil {
				respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}
		}

		events := buildLedgerEvents(txs, attempts)
		filtered := make([]ledgerEvent, 0, len(events))
		for _, ev := range events {
			if !ledgerEventMatches(ev, q, typeFilter, providerFilter, statusFilter) {
				continue
			}
			filtered = append(filtered, ev)
		}

		total := len(filtered)
		if offset < 0 {
			offset = 0
		}
		end := offset + limit
		if end > total {
			end = total
		}
		if offset > total {
			offset = total
		}

		respondJSON(w, http.StatusOK, map[string]any{
			"limit":   limit,
			"offset":  offset,
			"total":   total,
			"results": filtered[offset:end],
		})
	}
}

type ledgerEvent struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Time      time.Time `json:"time"`
	Provider  string    `json:"provider"`
	Reference string    `json:"reference"`
	Status    string    `json:"status"`
	Summary   string    `json:"summary"`
	TraceID   string    `json:"trace_id,omitempty"`
}

func buildLedgerEvents(txs []engine.Transaction, attempts []*webhook.Attempt) []ledgerEvent {
	const maxEvents = 1000
	events := make([]ledgerEvent, 0, len(txs)+len(attempts))
	for _, tx := range txs {
		ev := ledgerEvent{
			ID:        tx.ID,
			Type:      "transaction",
			Time:      tx.UpdatedAt,
			Provider:  tx.Provider,
			Reference: tx.Reference,
			Status:    string(tx.Status),
			Summary:   fmt.Sprintf("%s %.2f %s", tx.Provider, tx.Amount, tx.Currency),
			TraceID:   tx.TraceID,
		}
		if ev.Time.IsZero() {
			ev.Time = tx.CreatedAt
		}
		events = append(events, ev)
	}
	for _, a := range attempts {
		ev := ledgerEvent{
			ID:        a.ID,
			Type:      "webhook",
			Time:      a.UpdatedAt,
			Provider:  a.ProviderName,
			Reference: a.Ref,
			Status:    string(a.Status),
			Summary:   fmt.Sprintf("POST %s", a.URL),
			TraceID:   a.TraceID,
		}
		if ev.Time.IsZero() {
			ev.Time = a.CreatedAt
		}
		events = append(events, ev)
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Time.After(events[j].Time)
	})

	if len(events) > maxEvents {
		events = events[:maxEvents]
	}
	return events
}

func ledgerEventMatches(ev ledgerEvent, q, typeFilter, providerFilter, statusFilter string) bool {
	if typeFilter != "" && ev.Type != typeFilter {
		return false
	}
	if providerFilter != "" && ev.Provider != providerFilter {
		return false
	}
	if statusFilter != "" && ev.Status != statusFilter {
		return false
	}
	if q == "" {
		return true
	}
	return strings.Contains(strings.ToLower(ev.Reference), q) ||
		strings.Contains(strings.ToLower(ev.Provider), q) ||
		strings.Contains(strings.ToLower(ev.Status), q) ||
		strings.Contains(strings.ToLower(ev.Summary), q) ||
		strings.Contains(strings.ToLower(ev.TraceID), q)
}

func listProvidersHandler(cfg RouterConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		details := make(map[string]any, len(cfg.EnabledProviders))
		for _, name := range cfg.EnabledProviders {
			info := buildProviderInfo(name, cfg.ActiveProvider, cfg.Host, cfg.Port, cfg.PublicBaseURL, true, cfg.Providers)
			details[name] = info
		}

		available := cfg.AvailableProviders
		if len(available) == 0 {
			available = provider.Names()
		}

		respondJSON(w, http.StatusOK, map[string]any{
			"active":    cfg.ActiveProvider,
			"enabled":   cfg.EnabledProviders,
			"available": available,
			"providers": details,
		})
	}
}

func getProviderHandler(cfg RouterConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		if _, ok := config.WizardChoiceByKey(name); !ok {
			if _, loaded := cfg.Providers[name]; !loaded {
				if _, hasFactory := factory.Get(name); !hasFactory {
					respondJSON(w, http.StatusNotFound, map[string]string{"error": "provider not found"})
					return
				}
			}
		}

		enabled := false
		for _, n := range cfg.EnabledProviders {
			if n == name {
				enabled = true
				break
			}
		}

		info := buildProviderInfo(name, cfg.ActiveProvider, cfg.Host, cfg.Port, cfg.PublicBaseURL, enabled, cfg.Providers)
		info["webhook_target_url"] = providerWebhookTargetURL(cfg.ConfigPath, name)

		respondJSON(w, http.StatusOK, info)
	}
}

func buildProviderInfo(name, active, host string, port int, publicBaseURL string, enabled bool, providers map[string]provider.Provider) map[string]any {
	info := map[string]any{
		"name":    name,
		"enabled": enabled,
		"active":  name == active,
	}
	if choice, ok := config.WizardChoiceByKey(name); ok {
		info["display_name"] = choice.DisplayName
		info["description"] = choice.Description
		info["category"] = providerCategory(name)
		info["real_providers"] = realProvidersFor(name)
		info["sample_route"] = choice.SampleRoute
		info["sample_method"] = choice.SampleMethod
		info["docs_path"] = "/docs/providers/" + name + ".md"
		info["is_recommended_for_first_time"] = choice.IsRecommended
		info["env_vars"] = envVarNames(name, choice.EnvVarKeys)
	}

	version := "v1"
	var versions []string
	p, ok := providers[name]
	if !ok {
		if f, hasFactory := factory.Get(name); hasFactory {
			p, _ = f(nil)
		}
	}
	if p != nil {
		if vp, ok := p.(provider.VersionedProvider); ok {
			version = vp.CurrentVersion()
			versions = vp.Versions()
		}
	}
	if len(versions) == 0 {
		versions = []string{version}
	}
	info["version"] = version
	info["versions"] = versions
	info["version_details"] = buildVersionDetails(name, versions, choiceForProvider(name), host, port, publicBaseURL)
	info["base_url"] = providerBaseURL(host, port, publicBaseURL, name, version)

	return info
}

func choiceForProvider(name string) config.WizardChoice {
	if choice, ok := config.WizardChoiceByKey(name); ok {
		return choice
	}
	return config.WizardChoice{}
}

func envVarNames(provider string, keys []string) []string {
	out := make([]string, 0, len(keys))
	for _, key := range keys {
		out = append(out, config.EnvVarName(provider, key))
	}
	return out
}

func buildVersionDetails(name string, versions []string, choice config.WizardChoice, host string, port int, publicBaseURL string) map[string]any {
	details := make(map[string]any, len(versions))
	for _, v := range versions {
		sampleRoute := choice.SampleRoute
		if name == "fawry" {
			sampleRoute = "/fawry/" + v + "/charge"
		}
		details[v] = map[string]any{
			"base_url":     providerBaseURL(host, port, publicBaseURL, name, v),
			"sample_route": sampleRoute,
		}
	}
	return details
}

func providerBaseURL(host string, port int, publicBaseURL, name, version string) string {
	base := publicBaseURL
	if base == "" {
		if host == "" {
			host = "127.0.0.1"
		}
		if port <= 0 {
			port = 9000
		}
		base = fmt.Sprintf("http://%s:%d", host, port)
	}
	switch name {
	case "fawry":
		return base + "/fawry/" + version
	case "stripe":
		return base + "/v1"
	}
	return base
}

func providerWebhookTargetURL(configPath, name string) string {
	if configPath == "" {
		return ""
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		return ""
	}
	if target, ok := cfg.Webhook.Targets[name]; ok {
		return target
	}
	if pc, ok := cfg.Providers[name]; ok {
		if raw, ok := pc.Config["webhook_url"]; ok {
			if s, ok := raw.(string); ok {
				return s
			}
		}
	}
	return ""
}

type providerHealthResponse struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Reason string `json:"reason"`
}

func providerHealthHandler(configPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")

		cfg, err := config.Load(configPath)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		pc, ok := cfg.Providers[name]
		if !ok || !pc.Enabled {
			respondJSON(w, http.StatusOK, providerHealthResponse{
				Name:   name,
				Status: "disabled",
				Reason: "provider is not enabled",
			})
			return
		}

		var p provider.Provider
		if f, ok := factory.Get(name); ok {
			p, err = f(nil)
		} else {
			p, err = provider.Get(name)
		}
		if err != nil {
			respondJSON(w, http.StatusOK, providerHealthResponse{
				Name:   name,
				Status: "misconfigured",
				Reason: err.Error(),
			})
			return
		}

		if err := p.Init(pc.Config); err != nil {
			respondJSON(w, http.StatusOK, providerHealthResponse{
				Name:   name,
				Status: "misconfigured",
				Reason: err.Error(),
			})
			return
		}

		respondJSON(w, http.StatusOK, providerHealthResponse{
			Name:   name,
			Status: "healthy",
			Reason: "",
		})
	}
}

func providerCategory(name string) string {
	switch name {
	case "stripe":
		return "card"
	case "fawry":
		return "regional"
	case "billplz", "toyyibpay", "senangpay", "ipay88":
		return "redirect"
	default:
		return "diy"
	}
}

func realProvidersFor(name string) []string {
	switch name {
	case "stripe":
		return []string{"Stripe", "Stripe Checkout", "Stripe PaymentIntents"}
	case "fawry":
		return []string{"Fawry"}
	case "billplz":
		return []string{"Billplz"}
	case "toyyibpay":
		return []string{"ToyyibPay"}
	case "senangpay":
		return []string{"SenangPay"}
	case "ipay88":
		return []string{"iPay88"}
	default:
		return []string{"OpenMuara Default"}
	}
}

func onboardingHandler(active string, enabled []string, store engine.TransactionStore, dispatcher *webhook.Dispatcher) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		txs, _ := store.List(1, 0)
		var webhookReceived bool
		webhooksEnabled := dispatcher != nil && dispatcher.Store != nil
		if webhooksEnabled {
			attempts, _ := dispatcher.Store.List(1, 0)
			webhookReceived = len(attempts) > 0
		}

		respondJSON(w, http.StatusOK, map[string]any{
			"server_ready":           true,
			"providers_enabled":      len(enabled) > 0,
			"first_transaction":      len(txs) > 0,
			"first_webhook_received": webhookReceived,
			"webhooks_enabled":       webhooksEnabled,
			"active_provider":        active,
			"next_step":              nextStepForProvider(active),
		})
	}
}

func nextStepForProvider(name string) map[string]any {
	if name == "" {
		return map[string]any{"hint": "Enable a provider in config.yml and restart muara."}
	}
	if choice, ok := config.WizardChoiceByKey(name); ok {
		step := choice.NextStep()
		return map[string]any{
			"method": step.Method,
			"route":  step.Route,
			"hint":   step.Hint,
		}
	}
	return map[string]any{"hint": fmt.Sprintf("Send your first request to the %s provider.", name)}
}
