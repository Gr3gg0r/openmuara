package ui

import "net/http"

// ServeEscapePage writes the rendered Fawry escape page to w.
func ServeEscapePage(w http.ResponseWriter, data EscapePageData) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return RenderEscapePage(w, data)
}

// ServeDashboard writes the rendered admin dashboard to w.
func ServeDashboard(w http.ResponseWriter, data DashboardData) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return RenderDashboard(w, data)
}

// ServeStripeCheckoutPage writes the rendered Stripe checkout page to w.
func ServeStripeCheckoutPage(w http.ResponseWriter, data StripeCheckoutPageData) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return RenderStripeCheckoutPage(w, data)
}

// ServeStripeWebhooksPage writes the rendered Stripe webhook config page to w.
func ServeStripeWebhooksPage(w http.ResponseWriter, data StripeWebhooksPageData) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return RenderStripeWebhooksPage(w, data)
}

// ServeStripePaymentIntentPage writes the rendered Stripe PaymentIntent authentication page to w.
func ServeStripePaymentIntentPage(w http.ResponseWriter, data StripePaymentIntentPageData) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return RenderStripePaymentIntentPage(w, data)
}

// ServeBillplzPayPage writes the rendered Billplz payment page to w.
func ServeBillplzPayPage(w http.ResponseWriter, data BillplzPayPageData) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return RenderBillplzPayPage(w, data)
}

// ServeToyyibPayPage writes the rendered ToyyibPay payment page to w.
func ServeToyyibPayPage(w http.ResponseWriter, data ToyyibPayPageData) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return RenderToyyibPayPage(w, data)
}

// ServeIPay88PayPage writes the rendered iPay88 payment page to w.
func ServeIPay88PayPage(w http.ResponseWriter, data IPay88PayPageData) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return RenderIPay88PayPage(w, data)
}
