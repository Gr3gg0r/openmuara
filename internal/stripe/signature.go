package stripe

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/openmuara/openmuara/internal/errcode"
)

const signatureTolerance = 5 * time.Minute

// SignPayload returns a Stripe-Signature header value for the given payload.
func SignPayload(payload []byte, secret string, timestamp time.Time) string {
	timestampStr := strconv.FormatInt(timestamp.Unix(), 10)
	signedPayload := timestampStr + "." + string(payload)
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(signedPayload))
	sig := hex.EncodeToString(mac.Sum(nil))
	return fmt.Sprintf("t=%s,v1=%s", timestampStr, sig)
}

// VerifySignature checks a Stripe-Signature header against the payload and secret.
// It returns an error if the header is malformed, the timestamp is outside the
// tolerance window, or the signature does not match.
func VerifySignature(payload []byte, header, secret string) error {
	timestamp, sig, err := parseSignatureHeader(header)
	if err != nil {
		return err
	}

	if time.Since(timestamp).Abs() > signatureTolerance {
		return errcode.New(errcode.ESignatureMismatch, "stripe signature timestamp outside tolerance")
	}

	expected := SignPayload(payload, secret, timestamp)
	_, expectedSig, err := parseSignatureHeader(expected)
	if err != nil {
		return err
	}

	if !hmac.Equal([]byte(sig), []byte(expectedSig)) {
		return errcode.New(errcode.ESignatureMismatch, "stripe signature mismatch")
	}

	return nil
}

func parseSignatureHeader(header string) (time.Time, string, error) {
	var timestamp int64
	var sig string

	for _, part := range strings.Split(header, ",") {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		switch strings.TrimSpace(kv[0]) {
		case "t":
			ts, err := strconv.ParseInt(kv[1], 10, 64)
			if err != nil {
				return time.Time{}, "", errcode.Wrap(errcode.EInvalidRequest, "invalid stripe signature timestamp", err)
			}
			timestamp = ts
		case "v1":
			sig = kv[1]
		}
	}

	if timestamp == 0 {
		return time.Time{}, "", errcode.New(errcode.ESignatureMissing, "stripe signature timestamp missing")
	}
	if sig == "" {
		return time.Time{}, "", errcode.New(errcode.ESignatureMissing, "stripe signature v1 missing")
	}

	return time.Unix(timestamp, 0), sig, nil
}
