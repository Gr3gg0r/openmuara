package ipay88

import (
	"strings"
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/webhook"
)

func TestValidateEntryRequestMissingFields(t *testing.T) {
	base := PaymentRequest{
		MerchantCode:  "M00001",
		RefNo:         "REF-1",
		Amount:        "12.50",
		Currency:      "MYR",
		ProdDesc:      "Test",
		UserName:      "User",
		UserEmail:     "user@example.com",
		UserContact:   "012",
		Signature:     "sig",
		SignatureType: "SHA256",
		ResponseURL:   "http://example.com/r",
		BackendURL:    "http://example.com/b",
	}
	cases := []struct {
		name    string
		mutate  func(*PaymentRequest)
		wantErr string
	}{
		{"merchant code", func(r *PaymentRequest) { r.MerchantCode = "" }, "merchant code is required"},
		{"ref no", func(r *PaymentRequest) { r.RefNo = "" }, "ref no is required"},
		{"amount", func(r *PaymentRequest) { r.Amount = "" }, "amount is required"},
		{"currency", func(r *PaymentRequest) { r.Currency = "" }, "currency is required"},
		{"product description", func(r *PaymentRequest) { r.ProdDesc = "" }, "product description is required"},
		{"user name", func(r *PaymentRequest) { r.UserName = "" }, "user name is required"},
		{"user email", func(r *PaymentRequest) { r.UserEmail = "" }, "user email is required"},
		{"user contact", func(r *PaymentRequest) { r.UserContact = "" }, "user contact is required"},
		{"signature", func(r *PaymentRequest) { r.Signature = "" }, "signature is required"},
		{"signature type", func(r *PaymentRequest) { r.SignatureType = "" }, "signature type is required"},
		{"response url", func(r *PaymentRequest) { r.ResponseURL = "" }, "response url is required"},
		{"backend url", func(r *PaymentRequest) { r.BackendURL = "" }, "backend url is required"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := base
			tc.mutate(&r)
			err := validateEntryRequest(r)
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Errorf("want error containing %q, got %v", tc.wantErr, err)
			}
		})
	}
}

func TestMapIPay88StatusToWebhook(t *testing.T) {
	if got := mapIPay88StatusToWebhook("1"); got != webhook.PaymentStatusPaid {
		t.Errorf("success: want paid, got %q", got)
	}
	if got := mapIPay88StatusToWebhook("0"); got != webhook.PaymentStatusUnpaid {
		t.Errorf("failure: want unpaid, got %q", got)
	}
	if got := mapIPay88StatusToWebhook("6"); got != webhook.PaymentStatusUnpaid {
		t.Errorf("pending: want unpaid, got %q", got)
	}
}

func TestStripAmountInvalid(t *testing.T) {
	got := stripAmount("not-a-number")
	if got != "" {
		t.Errorf("want empty, got %q", got)
	}
}
