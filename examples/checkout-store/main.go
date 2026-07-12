// Command checkout-store is a minimal product landing page and checkout SPA
// that accepts one-time payments through OpenMuara's Fawry and Stripe emulators.
// It also receives webhooks and sends confirmation emails via Mailpit.
package main

import (
	"bytes"
	// #nosec G501 -- ToyyibPay gateway uses MD5 for signature emulation
	"crypto/md5"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
)

//go:embed all:web/dist
var distFS embed.FS

// Product is the single item sold by this example.
type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Currency    string  `json:"currency"`
	ImageURL    string  `json:"imageUrl"`
}

// Payment records a checkout attempt in the in-memory store.
type Payment struct {
	Ref         string  `json:"ref"`
	ProviderRef string  `json:"providerRef"`
	Method      string  `json:"method"`
	Status      string  `json:"status"`
	Email       string  `json:"email"`
	Name        string  `json:"name"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	ProductID   string  `json:"productId"`
	CreatedAt   int64   `json:"createdAt"`
}

type paymentStore struct {
	mu   sync.RWMutex
	data map[string]*Payment
}

func newPaymentStore() *paymentStore {
	return &paymentStore{data: make(map[string]*Payment)}
}

func (s *paymentStore) Get(ref string) (*Payment, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.data[ref]
	return p, ok
}

func (s *paymentStore) Put(p *Payment) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[p.Ref] = p
}

func (s *paymentStore) UpdateStatus(ref, status string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.data[ref]
	if !ok {
		return false
	}
	p.Status = status
	return true
}

// Default placeholder credentials. When a provider still uses one of these
// values the store reports it as "unconfigured" and the UI shows a demo-mode
// banner so it is obvious no real .env was set.
//
//nolint:gosec // G101: intentional demo placeholders, not real credentials.
const (
	defaultFawrySecurityKey   = "muara-fawry-secret"
	defaultStripeSecretKey    = "sk_test_muara"
	defaultToyyibPaySecretKey = "muara-toyyibpay-secret"
)

// Config is populated from environment variables.
type Config struct {
	Addr                   string
	OpenMuaraURL           string
	AppURL                 string
	PaymentMethods         []string
	FawryMerchantCode      string
	FawrySecurityKey       string
	FawryWebhookSecret     string
	StripeSecretKey        string
	StripeWebhookSecret    string
	ToyyibPayBaseURL       string
	ToyyibPayUserSecretKey string
	ToyyibPayCategoryCode  string
	MailpitHost            string
	MailpitPort            string
	MailFrom               string
}

func loadConfig() Config {
	return Config{
		Addr:                   envDefault("ADDR", ":8080"),
		OpenMuaraURL:           envDefault("OPENMUARA_URL", "http://127.0.0.1:9000"),
		AppURL:                 envDefault("APP_URL", "http://127.0.0.1:8080"),
		PaymentMethods:         envList("PAYMENT_METHODS", []string{"toyyibpay"}),
		FawryMerchantCode:      envDefault("FAWRY_MERCHANT_CODE", "muara-merchant-code"),
		FawrySecurityKey:       envDefault("FAWRY_SECURITY_KEY", defaultFawrySecurityKey),
		FawryWebhookSecret:     envDefault("FAWRY_WEBHOOK_SECRET", "muara-webhook-secret"),
		StripeSecretKey:        envDefault("STRIPE_SECRET_KEY", defaultStripeSecretKey),
		StripeWebhookSecret:    envDefault("STRIPE_WEBHOOK_SECRET", "whsec_muara"),
		ToyyibPayBaseURL:       normalizeBaseURL(envDefault("TOYYIBPAY_BASE_URL", envDefault("OPENMUARA_URL", "http://127.0.0.1:9000"))),
		ToyyibPayUserSecretKey: envDefault("TOYYIBPAY_USER_SECRET_KEY", defaultToyyibPaySecretKey),
		ToyyibPayCategoryCode:  envDefault("TOYYIBPAY_CATEGORY_CODE", "cat_openmuara"),
		MailpitHost:            envDefault("MAILPIT_HOST", "127.0.0.1"),
		MailpitPort:            envDefault("MAILPIT_PORT", "1025"),
		MailFrom:               envDefault("MAIL_FROM", "store@example.com"),
	}
}

func envDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// normalizeBaseURL ensures a scheme is present so values like
// "dev.toyyibpay.com" become "https://dev.toyyibpay.com".
func normalizeBaseURL(u string) string {
	u = strings.TrimRight(strings.TrimSpace(u), "/")
	if u == "" {
		return u
	}
	if strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://") {
		return u
	}
	return "https://" + u
}

func envList(key string, fallback []string) []string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	var out []string
	for _, part := range strings.Split(v, ",") {
		p := strings.TrimSpace(strings.ToLower(part))
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return fallback
	}
	return out
}

// methodEnabled reports whether a payment method is selectable.
func (c Config) methodEnabled(method string) bool {
	for _, m := range c.PaymentMethods {
		if m == method {
			return true
		}
	}
	return false
}

// providerConfigured reports whether a provider uses real (non-default)
// credentials. A false value means the .env was not customized for it.
func (c Config) providerConfigured(method string) bool {
	switch method {
	case "fawry":
		return c.FawrySecurityKey != defaultFawrySecurityKey
	case "stripe":
		return c.StripeSecretKey != defaultStripeSecretKey
	case "toyyibpay":
		return c.ToyyibPayUserSecretKey != defaultToyyibPaySecretKey
	}
	return false
}

// demoMode is true when at least one enabled provider still uses placeholder
// credentials. The UI uses it to show the announcement banner.
func (c Config) demoMode() bool {
	for _, m := range c.PaymentMethods {
		if !c.providerConfigured(m) {
			return true
		}
	}
	return false
}

func product() Product {
	return Product{
		ID:          "openmuara-course",
		Name:        "OpenMuara Course",
		Description: "A self-paced course on emulating billing and payments locally.",
		Price:       49.99,
		Currency:    "MYR",
		ImageURL:    "https://placehold.co/600x400/2563eb/ffffff?text=OpenMuara+Course",
	}
}

func main() {
	cfg := loadConfig()
	store := newPaymentStore()
	httpClient := &http.Client{Timeout: 10 * time.Second}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/product", productHandler())
	mux.HandleFunc("GET /api/config", configHandler(cfg))
	mux.HandleFunc("POST /api/checkout", checkoutHandler(cfg, store, httpClient))
	mux.HandleFunc("GET /api/payment/{ref}", paymentHandler(cfg, store, httpClient))
	mux.HandleFunc("GET /callback", callbackHandler(store))
	mux.HandleFunc("POST /webhook", webhookHandler(cfg, store))
	mux.Handle("GET /{path...}", spaHandler())

	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	slog.Info("checkout-store listening",
		"addr", cfg.Addr,
		"openmuara", cfg.OpenMuaraURL,
		"mailpit", cfg.MailpitHost+":"+cfg.MailpitPort,
	)
	if err := server.ListenAndServe(); err != nil {
		slog.Error("server exited", "error", err)
		os.Exit(1)
	}
}

func spaHandler() http.HandlerFunc {
	dist, err := fs.Sub(distFS, "web/dist")
	if err != nil {
		slog.Error("failed to open embedded dist", "error", err)
		return func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "frontend not built", http.StatusInternalServerError)
		}
	}
	fileServer := http.FileServer(http.FS(dist))

	return func(w http.ResponseWriter, r *http.Request) {
		urlPath := strings.TrimPrefix(r.URL.Path, "/")
		if urlPath == "" {
			urlPath = "."
		}

		// If the requested path matches a real file, serve it.
		cleanPath := path.Clean(urlPath)
		if cleanPath == "/" || cleanPath == "." {
			serveIndex(w, r)
			return
		}
		if _, err := fs.Stat(dist, cleanPath); err == nil {
			fileServer.ServeHTTP(w, r)
			return
		}

		// Otherwise fall back to index.html for SPA routes.
		serveIndex(w, r)
	}
}

func serveIndex(w http.ResponseWriter, _ *http.Request) {
	data, err := distFS.ReadFile("web/dist/index.html")
	if err != nil {
		http.Error(w, "index.html not found", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(data)
}

func productHandler() http.HandlerFunc {
	p := product()
	return func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, r, p)
	}
}

type providerConfigView struct {
	Enabled    bool `json:"enabled"`
	Configured bool `json:"configured"`
}

type configView struct {
	DemoMode  bool                          `json:"demoMode"`
	Providers map[string]providerConfigView `json:"providers"`
}

func configHandler(cfg Config) http.HandlerFunc {
	all := []string{"fawry", "stripe", "toyyibpay"}
	return func(w http.ResponseWriter, r *http.Request) {
		view := configView{
			DemoMode:  cfg.demoMode(),
			Providers: make(map[string]providerConfigView, len(all)),
		}
		for _, m := range all {
			view.Providers[m] = providerConfigView{
				Enabled:    cfg.methodEnabled(m),
				Configured: cfg.providerConfigured(m),
			}
		}
		respondJSON(w, r, view)
	}
}

func checkoutHandler(cfg Config, store *paymentStore, client *http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Method string `json:"method"`
			Email  string `json:"email"`
			Name   string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		if req.Method != "fawry" && req.Method != "stripe" && req.Method != "toyyibpay" {
			http.Error(w, "method must be fawry, stripe, or toyyibpay", http.StatusBadRequest)
			return
		}
		if !cfg.methodEnabled(req.Method) {
			http.Error(w, "payment method is disabled", http.StatusBadRequest)
			return
		}
		if req.Email == "" || req.Name == "" {
			http.Error(w, "email and name are required", http.StatusBadRequest)
			return
		}

		p := product()
		ref := fmt.Sprintf("cs-%d", time.Now().UnixNano())

		payment := &Payment{
			Ref:       ref,
			Method:    req.Method,
			Status:    "pending",
			Email:     req.Email,
			Name:      req.Name,
			Amount:    p.Price,
			Currency:  p.Currency,
			ProductID: p.ID,
			CreatedAt: time.Now().Unix(),
		}
		store.Put(payment)

		var redirectURL string
		var providerRef string
		var err error
		switch req.Method {
		case "fawry":
			redirectURL, err = createFawryCharge(cfg, client, payment, p)
		case "stripe":
			redirectURL, providerRef, err = createStripeSession(cfg, client, payment, p)
		case "toyyibpay":
			redirectURL, providerRef, err = createToyyibPayBill(cfg, client, payment, p)
		}
		if err != nil {
			slog.Error("checkout failed", "method", req.Method, "error", err)
			http.Error(w, "failed to create checkout", http.StatusBadGateway)
			return
		}
		payment.ProviderRef = providerRef
		store.Put(payment)

		respondJSON(w, r, map[string]any{
			"ok":          true,
			"ref":         ref,
			"redirectUrl": redirectURL,
		})
	}
}

func createFawryCharge(cfg Config, client *http.Client, payment *Payment, p Product) (string, error) {
	priceStr := fmt.Sprintf("%.2f", p.Price)
	quantity := 1
	returnURL := cfg.AppURL + "/callback?ref=" + payment.Ref

	itemPart := p.ID + strconv.Itoa(quantity) + priceStr
	message := cfg.FawryMerchantCode + payment.Ref + payment.Ref + returnURL + itemPart + cfg.FawrySecurityKey
	signature := sha256Hex(message)

	payload := map[string]any{
		"merchantCode":      cfg.FawryMerchantCode,
		"merchantRefNum":    payment.Ref,
		"customerEmail":     payment.Email,
		"customerName":      payment.Name,
		"customerProfileId": payment.Ref,
		"paymentExpiry":     time.Now().Add(24 * time.Hour).UnixMilli(),
		"language":          "en-gb",
		"chargeItems": []map[string]any{
			{"itemId": p.ID, "price": p.Price, "quantity": quantity},
		},
		"returnUrl": returnURL,
		"signature": signature,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, cfg.OpenMuaraURL+"/fawry/charge", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("fawry charge %d: %s", resp.StatusCode, string(respBody))
	}

	escapeURL := cfg.OpenMuaraURL + "/_admin/fawry-escape?ref=" + payment.Ref +
		"&returnUrl=" + url.QueryEscape(returnURL) +
		"&amount=" + priceStr
	return escapeURL, nil
}

func createStripeSession(cfg Config, client *http.Client, payment *Payment, p Product) (string, string, error) {
	payload := map[string]any{
		"payment_method_types": []string{"card"},
		"mode":                 "payment",
		"success_url":          cfg.AppURL + "/success?ref=" + payment.Ref,
		"cancel_url":           cfg.AppURL + "/cancel?ref=" + payment.Ref,
		"customer_email":       payment.Email,
		"client_reference_id":  payment.Ref,
		"line_items": []map[string]any{
			{
				"quantity": 1,
				"price_data": map[string]any{
					"currency":    strings.ToLower(p.Currency),
					"unit_amount": int64(p.Price * 100),
					"product_data": map[string]any{
						"name": p.Name,
					},
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", "", err
	}

	req, err := http.NewRequest(http.MethodPost, cfg.OpenMuaraURL+"/v1/checkout/sessions", bytes.NewReader(body))
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(cfg.StripeSecretKey, "")

	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("stripe session %d: %s", resp.StatusCode, string(respBody))
	}

	var session struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	}
	if err := json.Unmarshal(respBody, &session); err != nil {
		return "", "", err
	}
	if session.URL == "" {
		return "", "", fmt.Errorf("stripe session returned empty url")
	}
	return session.URL, session.ID, nil
}

func createToyyibPayBill(cfg Config, client *http.Client, payment *Payment, p Product) (string, string, error) {
	amountCents := int64(p.Price * 100)
	returnURL := cfg.AppURL + "/callback?ref=" + payment.Ref
	callbackURL := cfg.AppURL + "/webhook"

	form := url.Values{}
	form.Set("userSecretKey", cfg.ToyyibPayUserSecretKey)
	form.Set("categoryCode", cfg.ToyyibPayCategoryCode)
	form.Set("billName", p.Name)
	form.Set("billDescription", p.Description)
	form.Set("billPriceSetting", "1") // 1 = fixed price (required by real ToyyibPay)
	form.Set("billPayorInfo", "1")    // 1 = collect payor info (required by real ToyyibPay)
	form.Set("billTo", payment.Name)
	form.Set("billEmail", payment.Email)
	form.Set("billPhone", "0123456789")
	form.Set("billAmount", strconv.FormatInt(amountCents, 10))
	form.Set("billReturnUrl", returnURL)
	form.Set("billCallbackUrl", callbackURL)
	form.Set("billPaymentChannel", "2")
	form.Set("billExternalReferenceNo", payment.Ref)

	req, err := http.NewRequest(
		http.MethodPost,
		cfg.ToyyibPayBaseURL+"/index.php/api/createBill",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("toyyibpay createBill %d: %s", resp.StatusCode, string(respBody))
	}

	billCode, paymentLink, err := parseToyyibPayCreateBill(cfg.ToyyibPayBaseURL, respBody)
	if err != nil {
		return "", "", err
	}
	return paymentLink, billCode, nil
}

// parseToyyibPayCreateBill accepts both the real ToyyibPay array response
// (`[{"BillCode":"..."}]`) and the OpenMuara object response
// (`{"bill":{"billCode":"...","billPaymentLink":"..."}}`) so the same example
// code works against either endpoint by only changing TOYYIBPAY_BASE_URL.
func parseToyyibPayCreateBill(baseURL string, body []byte) (billCode, paymentLink string, err error) {
	// Real ToyyibPay: a JSON array of objects with BillCode.
	var arr []map[string]any
	if json.Unmarshal(body, &arr) == nil && len(arr) > 0 {
		if code, ok := arr[0]["BillCode"].(string); ok && code != "" {
			return code, strings.TrimRight(baseURL, "/") + "/" + code, nil
		}
	}

	// OpenMuara: an object wrapping a bill with a ready payment link.
	var obj struct {
		Bill struct {
			BillCode        string `json:"billCode"`
			BillPaymentLink string `json:"billPaymentLink"`
		} `json:"bill"`
	}
	if json.Unmarshal(body, &obj) == nil && obj.Bill.BillCode != "" {
		link := obj.Bill.BillPaymentLink
		if link == "" {
			link = strings.TrimRight(baseURL, "/") + "/" + obj.Bill.BillCode
		}
		return obj.Bill.BillCode, link, nil
	}

	return "", "", fmt.Errorf("unrecognized toyyibpay createBill response: %s", string(body))
}

func paymentHandler(cfg Config, store *paymentStore, client *http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ref := r.PathValue("ref")
		if !validRef(ref) {
			http.Error(w, "invalid ref", http.StatusBadRequest)
			return
		}
		p, ok := store.Get(ref)
		if !ok {
			http.Error(w, "payment not found", http.StatusNotFound)
			return
		}

		// For Stripe, sync the latest status from OpenMuara so the success page
		// reflects payment completion even when webhooks are not configured.
		if p.Method == "stripe" && p.Status == "pending" && p.ProviderRef != "" {
			if status := syncStripeStatus(cfg, client, p.ProviderRef); status != "" {
				store.UpdateStatus(ref, status)
				p, _ = store.Get(ref)
			}
		}

		respondJSON(w, r, p)
	}
}

func syncStripeStatus(cfg Config, client *http.Client, sessionID string) string {
	req, err := http.NewRequest(http.MethodGet, cfg.OpenMuaraURL+"/v1/checkout/sessions/"+sessionID, nil)
	if err != nil {
		return ""
	}
	req.SetBasicAuth(cfg.StripeSecretKey, "")

	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	if resp.StatusCode != http.StatusOK {
		return ""
	}

	var session struct {
		Status        string `json:"status"`
		PaymentStatus string `json:"payment_status"`
	}
	if err := json.Unmarshal(body, &session); err != nil {
		return ""
	}
	if session.PaymentStatus == "paid" || session.Status == "complete" {
		return "paid"
	}
	if session.Status == "expired" {
		return "canceled"
	}
	return ""
}

func callbackHandler(store *paymentStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ref := r.URL.Query().Get("ref")
		// Fawry's escape page sends orderStatus; Stripe-style callbacks use status;
		// ToyyibPay returns status_id ("1" success, "3" failed).
		status := r.URL.Query().Get("orderStatus")
		if status == "" {
			status = r.URL.Query().Get("status")
		}
		if status == "" {
			status = mapToyyibPayStatusID(r.URL.Query().Get("status_id"))
		}
		if !validRef(ref) {
			http.Error(w, "invalid ref", http.StatusBadRequest)
			return
		}
		if status == "" {
			status = "pending"
		}
		mapped := mapStatus(status)
		if !store.UpdateStatus(ref, mapped) {
			http.Error(w, "payment not found", http.StatusNotFound)
			return
		}
		redirectPath := "/success"
		if mapped == "canceled" {
			redirectPath = "/cancel"
		}
		// #nosec G710 -- ref is validated by validRef() above; redirect target is fixed.
		http.Redirect(w, r, redirectPath+"?ref="+ref, http.StatusSeeOther)
	}
}

func webhookHandler(cfg Config, store *paymentStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusBadRequest)
			return
		}
		defer func() { _ = r.Body.Close() }()

		var ref, status, method string
		if strings.Contains(r.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
			ref, status, method = extractToyyibPayWebhook(cfg, body)
		} else {
			ref, status, method = extractWebhookInfo(body)
		}

		if ref != "" {
			store.UpdateStatus(ref, status)
			if p, ok := store.Get(ref); ok {
				go sendConfirmationEmail(cfg, p)
			}
		}

		slog.Info("webhook received", "ref", ref, "status", status, "method", method)
		w.WriteHeader(http.StatusOK)
	}
}

// extractToyyibPayWebhook parses a form-encoded ToyyibPay callback. The hash is
// verified when TOYYIBPAY_USER_SECRET_KEY is set so the same code works against
// both the real gateway and OpenMuara.
func extractToyyibPayWebhook(cfg Config, body []byte) (ref, status, method string) {
	values, err := url.ParseQuery(string(body))
	if err != nil {
		return "", "", ""
	}
	ref = values.Get("order_id")
	status = mapToyyibPayStatus(values.Get("status"))
	method = "toyyibpay"

	if cfg.ToyyibPayUserSecretKey != "" {
		given := values.Get("hash")
		expected := toyyibPayHash(cfg.ToyyibPayUserSecretKey, values.Get("status"), values.Get("order_id"), values.Get("refno"))
		if given != "" && !strings.EqualFold(given, expected) {
			slog.Warn("toyyibpay callback hash mismatch", "ref", ref)
		}
	}
	return ref, status, method
}

func extractWebhookInfo(body []byte) (ref, status, method string) {
	var event map[string]any
	if err := json.Unmarshal(body, &event); err != nil {
		return "", "", ""
	}

	// Stripe-style event wrapper.
	if data, ok := event["data"].(map[string]any); ok {
		if obj, ok := data["object"].(map[string]any); ok {
			if id, ok := obj["id"].(string); ok && strings.HasPrefix(id, "cs_") {
				ref = id
				method = "stripe"
			}
			if clientRef, ok := obj["client_reference_id"].(string); ok && ref == "" {
				ref = clientRef
			}
		}
		if typ, ok := event["type"].(string); ok {
			status = mapStripeEventType(typ)
		}
		return ref, status, method
	}

	// Fawry V2 webhook shape.
	if merchantRef, ok := event["merchantRefNumber"].(string); ok {
		ref = merchantRef
		method = "fawry"
		if s, ok := event["orderStatus"].(string); ok {
			status = mapStatus(s)
		}
		return ref, status, method
	}

	return "", "", ""
}

func mapStatus(s string) string {
	switch strings.ToUpper(s) {
	case "PAID", "SUCCESS", "PAID_SUCCESS", "checkout.session.completed", "payment_intent.succeeded":
		return "paid"
	case "CANCELED", "CANCELLED", "UNPAID", "checkout.session.expired", "payment_intent.canceled":
		return "canceled"
	default:
		return "pending"
	}
}

// mapToyyibPayStatusID maps ToyyibPay's numeric status_id to the store status.
func mapToyyibPayStatusID(id string) string {
	switch id {
	case "1":
		return "paid"
	case "3":
		return "canceled"
	default:
		return "pending"
	}
}

// mapToyyibPayStatus maps ToyyibPay's callback status ("1"/"3") to the store.
func mapToyyibPayStatus(status string) string {
	return mapToyyibPayStatusID(status)
}

// toyyibPayHash computes MD5(secret + status + orderID + refno + "ok") per the
// ToyyibPay callback signature scheme.
func toyyibPayHash(secret, status, orderID, refno string) string {
	// #nosec G401 -- ToyyibPay signature uses MD5 by provider spec
	sum := md5.Sum([]byte(secret + status + orderID + refno + "ok"))
	return hex.EncodeToString(sum[:])
}

func mapStripeEventType(typ string) string {
	switch typ {
	case "checkout.session.completed", "payment_intent.succeeded":
		return "paid"
	case "checkout.session.expired", "payment_intent.canceled":
		return "canceled"
	default:
		return "pending"
	}
}

func sendConfirmationEmail(cfg Config, p *Payment) {
	if p.Status != "paid" {
		return
	}
	subject := fmt.Sprintf("Payment confirmed for %s", p.Name)
	body := fmt.Sprintf(
		"Hi %s,\r\n\r\nThank you for your purchase of %s (%.2f %s).\r\nReference: %s\r\n\r\n- OpenMuara Store",
		p.Name, product().Name, p.Amount, p.Currency, p.Ref,
	)
	msg := []byte(fmt.Sprintf("To: %s\r\nFrom: %s\r\nSubject: %s\r\n\r\n%s", p.Email, cfg.MailFrom, subject, body))

	addr := cfg.MailpitHost + ":" + cfg.MailpitPort
	if err := smtp.SendMail(addr, nil, cfg.MailFrom, []string{p.Email}, msg); err != nil {
		slog.Error("failed to send email", "error", err, "ref", p.Ref)
	} else {
		slog.Info("confirmation email sent", "ref", p.Ref, "to", p.Email)
	}
}

func respondJSON(w http.ResponseWriter, _ *http.Request, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func sha256Hex(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

func validRef(ref string) bool {
	if ref == "" || len(ref) > 128 {
		return false
	}
	for _, r := range ref {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' && r != '_' {
			return false
		}
	}
	return true
}
