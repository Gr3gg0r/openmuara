package ipay88

import (
	"context"
	"fmt"

	"github.com/openmuara/openmuara/internal/provider"
	"github.com/openmuara/openmuara/internal/webhook"
)

func (p *Provider) buildPayload(_ context.Context, tx provider.Transaction) ([]byte, error) {
	req, ok := p.getRequest(tx.Reference)
	if !ok {
		return nil, fmt.Errorf("payment request not found for ref %q", tx.Reference)
	}

	status := mapWebhookStatusToIPay88(tx.Status)
	paymentID := req.SelectedPaymentID
	if paymentID == "" {
		paymentID = req.PaymentID
	}

	values := responseValues(req, p.merchantCode, paymentID, req.Amount, req.Currency, status, p.merchantKey)
	return []byte(values.Encode()), nil
}

func mapWebhookStatusToIPay88(status string) string {
	if status == string(webhook.PaymentStatusPaid) {
		return string(PaymentStatusSuccess)
	}
	return string(PaymentStatusFailure)
}
