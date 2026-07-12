// Command checkout-store is a minimal product landing page and checkout SPA
// that accepts one-time payments through OpenMuara's Fawry and Stripe emulators.
// It also receives webhooks and sends confirmation emails via Mailpit.
package main

import (
	"bytes"
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

// Config is populated from environment variables.
type Config struct {
	Addr                string
	OpenMuaraURL        string
	AppURL              string
	FawryMerchantCode   string
	FawrySecurityKey    string
	FawryWebhookSecret  string
	StripeSecretKey     string
	StripeWebhookSecret string
	MailpitHost         string
	MailpitPort         string
	MailFrom            string
}

func loadConfig() Config {
	return Config{
		Addr:                envDefault("ADDR", ":8080"),
		OpenMuaraURL:        envDefault("OPENMUARA_URL", "http://127.0.0.1:9000"),
		AppURL:              envDefault("APP_URL", "http://127.0.0.1:8080"),
		FawryMerchantCode:   envDefault("FAWRY_MERCHANT_CODE", "muara-merchant-code"),
		FawrySecurityKey:    envDefault("FAWRY_SECURITY_KEY", "muara-fawry-secret"),
		FawryWebhookSecret:  envDefault("FAWRY_WEBHOOK_SECRET", "muara-webhook-secret"),
		StripeSecretKey:     envDefault("STRIPE_SECRET_KEY", "sk_test_muara"),
		StripeWebhookSecret: envDefault("STRIPE_WEBHOOK_SECRET", "whsec_muara"),
		MailpitHost:         envDefault("MAILPIT_HOST", "127.0.0.1"),
		MailpitPort:         envDefault("MAILPIT_PORT", "1025"),
		MailFrom:            envDefault("MAIL_FROM", "store@example.com"),
	}
}

func envDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func product() Product {
	return Product{
		ID:          "openmuara-course",
		Name:        "OpenMuara Course",
		Description: "A self-paced course on emulating billing and payments locally.",
		Price:       49.99,
		Currency:    "EGP",
		ImageURL:    "https://placehold.co/600x400/2563eb/ffffff?text=OpenMuara+Course",
	}
}

func main() {
	cfg := loadConfig()
	store := newPaymentStore()
	httpClient := &http.Client{Timeout: 10 * time.Second}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/product", productHandler())
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
		if req.Method != "fawry" && req.Method != "stripe" {
			http.Error(w, "method must be fawry or stripe", http.StatusBadRequest)
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
		// Fawry's escape page sends orderStatus; Stripe-style callbacks use status.
		status := r.URL.Query().Get("orderStatus")
		if status == "" {
			status = r.URL.Query().Get("status")
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

		ref, status, method := extractWebhookInfo(body)
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
