package ipay88

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var nonDigitRe = regexp.MustCompile(`[^0-9]`)

// SignRequest computes the iPay88 payment request signature.
// Canonical string: MerchantKey + MerchantCode + RefNo + Amount(stripped) + Currency.
func SignRequest(merchantKey, merchantCode, refNo, amount, currency string) string {
	msg := merchantKey + merchantCode + refNo + stripAmount(amount) + currency
	return sha256Hex(msg)
}

// SignResponse computes the iPay88 response/backend signature.
// Canonical string: MerchantKey + MerchantCode + PaymentId + RefNo + Amount(stripped) + Currency + Status.
func SignResponse(merchantKey, merchantCode, paymentID, refNo, amount, currency, status string) string {
	msg := merchantKey + merchantCode + paymentID + refNo + stripAmount(amount) + currency + status
	return sha256Hex(msg)
}

// VerifyRequest validates the request signature.
func VerifyRequest(req PaymentRequest, merchantKey string) bool {
	expected := SignRequest(merchantKey, req.MerchantCode, req.RefNo, req.Amount, req.Currency)
	return strings.EqualFold(req.Signature, expected)
}

// VerifyResponse validates the response/backend signature.
func VerifyResponse(merchantKey, merchantCode, paymentID, refNo, amount, currency, status, signature string) bool {
	expected := SignResponse(merchantKey, merchantCode, paymentID, refNo, amount, currency, status)
	return strings.EqualFold(signature, expected)
}

// stripAmount removes thousand separators and decimal points after formatting
// to two decimal places so the result is the integer amount in the smallest
// displayed unit.
func stripAmount(amount string) string {
	value, err := parseAmount(amount)
	if err != nil {
		return nonDigitRe.ReplaceAllString(amount, "")
	}
	formatted := fmt.Sprintf("%.2f", value)
	return nonDigitRe.ReplaceAllString(formatted, "")
}

func sha256Hex(s string) string {
	sum := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", sum[:])
}

// parseAmount parses an amount string that may contain commas as thousand
// separators and a dot as the decimal separator.
func parseAmount(amount string) (float64, error) {
	clean := strings.ReplaceAll(amount, ",", "")
	return strconv.ParseFloat(clean, 64)
}

// amountInCents returns the integer amount in the smallest currency unit.
func amountInCents(amount float64) int64 {
	return int64(amount*100 + 0.5)
}
