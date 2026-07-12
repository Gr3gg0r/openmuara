package toyyibpay

import (
	// #nosec G501 -- ToyyibPay gateway uses MD5 for signature emulation
	"crypto/md5"
	"fmt"
	"net/url"
	"strings"
)

// ComputeHash calculates the ToyyibPay callback MD5 hash.
// hash = MD5(userSecretKey + status + order_id + refno + "ok")
func ComputeHash(secret, status, orderID, refno string) string {
	msg := secret + status + orderID + refno + "ok"
	// #nosec G401 -- ToyyibPay signature uses MD5 by provider spec
	return fmt.Sprintf("%x", md5.Sum([]byte(msg)))
}

// VerifyCallback checks the MD5 hash on an incoming callback payload.
func VerifyCallback(secret string, values url.Values) bool {
	given := values.Get("hash")
	if given == "" {
		return false
	}
	expected := ComputeHash(secret, values.Get("status"), values.Get("order_id"), values.Get("refno"))
	return strings.EqualFold(given, expected)
}
