package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// Signer computes a webhook message signature.
type Signer interface {
	Sign(payload FawryV2Payload) (string, error)
}

// HMACSigner computes an HMAC-SHA256 signature over a canonical JSON representation.
// This is an OpenMuara approximation of the Fawry V2 messageSignature.
type HMACSigner struct {
	secret string
}

// NewHMACSigner creates a signer using the given secret.
func NewHMACSigner(secret string) *HMACSigner {
	return &HMACSigner{secret: secret}
}

// Sign returns a hex-encoded HMAC-SHA256 of the canonical payload JSON.
func (s *HMACSigner) Sign(payload FawryV2Payload) (string, error) {
	canonical, err := canonicalPayload(payload)
	if err != nil {
		return "", fmt.Errorf("canonicalize payload: %w", err)
	}

	mac := hmac.New(sha256.New, []byte(s.secret))
	if _, err := mac.Write(canonical); err != nil {
		return "", fmt.Errorf("compute hmac: %w", err)
	}

	return hex.EncodeToString(mac.Sum(nil)), nil
}

// canonicalPayload returns a stable JSON encoding of the payload fields.
// The messageSignature field is excluded because it is what we are computing.
func canonicalPayload(payload FawryV2Payload) ([]byte, error) {
	v := struct {
		RequestID             string      `json:"requestId"`
		FawryRefNumber        string      `json:"fawryRefNumber"`
		MerchantRefNumber     string      `json:"merchantRefNumber"`
		CustomerMobile        string      `json:"customerMobile"`
		CustomerMail          string      `json:"customerMail"`
		CustomerMerchantID    string      `json:"customerMerchantId"`
		PaymentAmount         float64     `json:"paymentAmount"`
		OrderAmount           float64     `json:"orderAmount"`
		FawryFees             float64     `json:"fawryFees"`
		OrderStatus           string      `json:"orderStatus"`
		PaymentMethod         string      `json:"paymentMethod"`
		PaymentTime           int64       `json:"paymentTime"`
		PaymentRefrenceNumber string      `json:"paymentRefrenceNumber"`
		OrderExpiryDate       int64       `json:"orderExpiryDate"`
		OrderItems            []OrderItem `json:"orderItems"`
	}{
		RequestID:             payload.RequestID,
		FawryRefNumber:        payload.FawryRefNumber,
		MerchantRefNumber:     payload.MerchantRefNumber,
		CustomerMobile:        payload.CustomerMobile,
		CustomerMail:          payload.CustomerMail,
		CustomerMerchantID:    payload.CustomerMerchantID,
		PaymentAmount:         payload.PaymentAmount,
		OrderAmount:           payload.OrderAmount,
		FawryFees:             payload.FawryFees,
		OrderStatus:           payload.OrderStatus,
		PaymentMethod:         payload.PaymentMethod,
		PaymentTime:           payload.PaymentTime,
		PaymentRefrenceNumber: payload.PaymentRefrenceNumber,
		OrderExpiryDate:       payload.OrderExpiryDate,
		OrderItems:            payload.OrderItems,
	}

	return json.Marshal(v)
}

// Verify checks whether the given signature matches the payload.
func (s *HMACSigner) Verify(payload FawryV2Payload, signature string) (bool, error) {
	expected, err := s.Sign(payload)
	if err != nil {
		return false, err
	}
	return hmac.Equal([]byte(expected), []byte(signature)), nil
}
