// Package senangpay emulates the SenangPay payment gateway.
package senangpay

import (
	// #nosec G501 -- SenangPay gateway uses MD5 for signature emulation
	"crypto/md5"
	"fmt"
	"strings"
)

// Sign computes the SenangPay MD5 signature for a charge request.
// The canonical string is: secret_key + detail + amount + order_id
func Sign(secret, detail string, amount float64, orderID string) string {
	msg := fmt.Sprintf("%s%s%.2f%s", secret, detail, amount, orderID)
	// #nosec G401 -- SenangPay signature uses MD5 by provider spec
	return fmt.Sprintf("%x", md5.Sum([]byte(msg)))
}

// Verify checks the request signature against the computed value.
func Verify(req ChargeRequest, secret string) bool {
	expected := Sign(secret, req.Detail, req.Amount, req.OrderID)
	return strings.EqualFold(req.Hash, expected)
}

// SignStatusQuery computes the SenangPay MD5 signature for a status query.
// The canonical string is: secret_key + order_id
func SignStatusQuery(secret, orderID string) string {
	msg := secret + orderID
	// #nosec G401 -- SenangPay signature uses MD5 by provider spec
	return fmt.Sprintf("%x", md5.Sum([]byte(msg)))
}

// VerifyStatusQuery checks the status query hash against the computed value.
func VerifyStatusQuery(orderID, hash, secret string) bool {
	expected := SignStatusQuery(secret, orderID)
	return strings.EqualFold(hash, expected)
}
