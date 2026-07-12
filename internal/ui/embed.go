// Package ui renders embedded HTML pages and the SPA dashboard for the OpenMuara admin interface.
package ui

import (
	"embed"
	"html/template"
	"io"
	"io/fs"
	"net/http"

	"github.com/openmuara/openmuara/internal/httputil"
)

//go:embed fawry-escape.html stripe-checkout.html stripe-webhooks.html stripe-payment-intent.html billplz-pay.html toyyibpay-pay.html ipay88-pay.html
var legacyAssets embed.FS

//go:embed dashboard-dist
var dashboardDist embed.FS

//go:embed payment-pages.css
var paymentPagesCSS []byte

// PaymentPagesCSS returns the shared stylesheet for provider payment and
// simulation pages.
func PaymentPagesCSS() []byte {
	return paymentPagesCSS
}

// EscapePageData is injected into the escape page template.
type EscapePageData struct {
	Ref       string
	ReturnURL string
	Amount    string
	CSRFToken string
}

// DashboardData is injected into the dashboard template.
type DashboardData struct {
	ActiveProvider  string
	CSRFToken       string
	AdminAPIBaseURL string
	Role            string
}

// StripeCheckoutPageData is injected into the Stripe checkout page template.
type StripeCheckoutPageData struct {
	ID                      string
	AmountTotal             int64
	Currency                string
	AmountTotalDisplay      string
	LineItems               []StripeCheckoutLineItem
	PaymentMethodTypes      []string
	ShowCard                bool
	ShowFPX                 bool
	ShowPaymentMethodToggle bool
	FPXBanks                []FPXBank
	CSRFToken               string
}

// StripeCheckoutLineItem is a single line item shown on the checkout page.
type StripeCheckoutLineItem struct {
	Name     string
	Quantity int64
	Amount   int64
	Currency string
}

// FPXBank represents a selectable Malaysian bank on the checkout page.
type FPXBank struct {
	Code string
	Name string
}

// StripeWebhookEventOption is a selectable event on the webhook config page.
type StripeWebhookEventOption struct {
	Name    string
	Checked bool
}

// StripeWebhooksPageData is injected into the Stripe webhook config template.
type StripeWebhooksPageData struct {
	URL           string
	WebhookSecret string
	Events        []StripeWebhookEventOption
	CSRFToken     string
}

// StripePaymentIntentPageData is injected into the Stripe PaymentIntent authentication template.
type StripePaymentIntentPageData struct {
	ID                 string
	Amount             int64
	Currency           string
	AmountDisplay      string
	Status             string
	PaymentMethodTypes []string
	ShowCard           bool
	ShowFPX            bool
	FPXBanks           []FPXBank
	CSRFToken          string
}

// BillplzPayPageData is injected into the Billplz payment page template.
type BillplzPayPageData struct {
	ID          string
	Amount      int64
	Description string
	Methods     []BillplzPaymentMethod
	CSRFToken   string
}

// BillplzPaymentMethod is a selectable payment method on the Billplz page.
type BillplzPaymentMethod struct {
	Code string
	Name string
}

// ToyyibPayPageData is injected into the ToyyibPay payment page template.
type ToyyibPayPageData struct {
	BillCode      string
	BillName      string
	Amount        int64
	AmountDisplay string
	Channel       int
	CSRFToken     string
}

// IPay88PayPageData is injected into the iPay88 payment page template.
type IPay88PayPageData struct {
	RefNo         string
	Amount        int64
	AmountDisplay string
	Currency      string
	Description   string
	Methods       []IPay88PaymentMethod
	CSRFToken     string
}

// IPay88PaymentMethod is a selectable payment method on the iPay88 page.
type IPay88PaymentMethod struct {
	ID   string
	Name string
}

// dashboardTemplate is the parsed SPA entry HTML with a CSRF token placeholder.
var dashboardTemplate = template.Must(template.ParseFS(dashboardDist, "dashboard-dist/index.html"))

// DashboardAssetsFS returns the embedded dashboard asset filesystem rooted at dist.
func DashboardAssetsFS() http.FileSystem {
	sub, err := fs.Sub(dashboardDist, "dashboard-dist")
	if err != nil {
		// Sub should never fail for an embedded tree; return an empty FS as fallback.
		return http.FS(emptyFS{})
	}
	return http.FS(sub)
}

type emptyFS struct{}

func (emptyFS) Open(string) (fs.File, error) { return nil, fs.ErrNotExist }

// RenderEscapePage renders the Fawry escape simulation page.
func RenderEscapePage(w io.Writer, data EscapePageData) error {
	tmpl, err := template.ParseFS(legacyAssets, "fawry-escape.html")
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}

// RenderDashboard renders the main admin dashboard SPA page.
func RenderDashboard(w io.Writer, data DashboardData) error {
	return dashboardTemplate.Execute(w, data)
}

// RenderStripeCheckoutPage renders the Stripe checkout page.
func RenderStripeCheckoutPage(w io.Writer, data StripeCheckoutPageData) error {
	tmpl, err := template.ParseFS(legacyAssets, "stripe-checkout.html")
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}

// RenderStripeWebhooksPage renders the Stripe webhook configuration page.
func RenderStripeWebhooksPage(w io.Writer, data StripeWebhooksPageData) error {
	tmpl, err := template.ParseFS(legacyAssets, "stripe-webhooks.html")
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}

// RenderStripePaymentIntentPage renders the Stripe PaymentIntent authentication page.
func RenderStripePaymentIntentPage(w io.Writer, data StripePaymentIntentPageData) error {
	tmpl, err := template.ParseFS(legacyAssets, "stripe-payment-intent.html")
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}

// RenderBillplzPayPage renders the Billplz payment page.
func RenderBillplzPayPage(w io.Writer, data BillplzPayPageData) error {
	tmpl, err := template.ParseFS(legacyAssets, "billplz-pay.html")
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}

// RenderToyyibPayPage renders the ToyyibPay payment page.
func RenderToyyibPayPage(w io.Writer, data ToyyibPayPageData) error {
	tmpl, err := template.ParseFS(legacyAssets, "toyyibpay-pay.html")
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}

// RenderIPay88PayPage renders the iPay88 payment page.
func RenderIPay88PayPage(w io.Writer, data IPay88PayPageData) error {
	tmpl, err := template.ParseFS(legacyAssets, "ipay88-pay.html")
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}

// CSRFTokenFromRequest extracts the CSRF token from the request context or cookie.
func CSRFTokenFromRequest(r *http.Request) string {
	if tok, ok := httputil.CSRFTokenFromContext(r.Context()); ok {
		return tok
	}
	c, err := r.Cookie("openmuara_csrf")
	if err == nil && c.Value != "" {
		return c.Value
	}
	return ""
}
