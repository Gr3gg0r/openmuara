package fawry

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
)

// ChargeItem represents a single line item in a Fawry charge request.
type ChargeItem struct {
	ItemID   string  `json:"itemId"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

// ChargeRequest matches the Fawry Express Checkout charge request shape.
type ChargeRequest struct {
	MerchantCode      string       `json:"merchantCode"`
	MerchantRefNum    string       `json:"merchantRefNum"`
	CustomerEmail     string       `json:"customerEmail"`
	CustomerName      string       `json:"customerName"`
	CustomerProfileID string       `json:"customerProfileId"`
	PaymentExpiry     int64        `json:"paymentExpiry"`
	Language          string       `json:"language"`
	ChargeItems       []ChargeItem `json:"chargeItems"`
	ReturnURL         string       `json:"returnUrl"`
	Signature         string       `json:"signature"`
}

// ChargeResponse is the OpenMuara charge response.
type ChargeResponse struct {
	Status    string `json:"status"`
	Reference string `json:"reference"`
}

// Sign computes the Fawry-style SHA256 signature.
// Concatenation order (official Fawry Express Checkout):
// merchantCode + merchantRefNum + customerProfileId (or "") + returnUrl + itemId + quantity + price(2 decimals) + secureKey
// For multiple items, sort by itemId then concatenate itemId + quantity + price for each.
func Sign(req ChargeRequest, merchantSecurityKey string) string {
	items := make([]ChargeItem, len(req.ChargeItems))
	copy(items, req.ChargeItems)
	sort.Slice(items, func(i, j int) bool { return items[i].ItemID < items[j].ItemID })

	var itemPart string
	for _, it := range items {
		itemPart += fmt.Sprintf("%s%d%s", it.ItemID, it.Quantity, fmt.Sprintf("%.2f", it.Price))
	}

	text := fmt.Sprintf("%s%s%s%s%s%s",
		req.MerchantCode,
		req.MerchantRefNum,
		req.CustomerProfileID,
		req.ReturnURL,
		itemPart,
		merchantSecurityKey,
	)
	hash := sha256.Sum256([]byte(text))
	return hex.EncodeToString(hash[:])
}

// Verify checks the request signature against the computed signature.
func Verify(req ChargeRequest, merchantSecurityKey string) bool {
	return req.Signature == Sign(req, merchantSecurityKey)
}

// StatusQuery represents the signed parameters for a Fawry payment-status request.
type StatusQuery struct {
	MerchantCode   string `json:"merchantCode"`
	MerchantRefNum string `json:"merchantRefNum"`
	Signature      string `json:"signature"`
}

// SignStatusQuery computes the SHA256 signature for a payment-status request.
// Concatenation order: merchantCode + merchantRefNum + merchantSecurityKey.
func SignStatusQuery(q StatusQuery, merchantSecurityKey string) string {
	text := q.MerchantCode + q.MerchantRefNum + merchantSecurityKey
	hash := sha256.Sum256([]byte(text))
	return hex.EncodeToString(hash[:])
}

// VerifyStatusQuery checks the status query signature against the computed signature.
func VerifyStatusQuery(q StatusQuery, merchantSecurityKey string) bool {
	return q.Signature == SignStatusQuery(q, merchantSecurityKey)
}
