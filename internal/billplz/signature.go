package billplz

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
)

// Sign computes the Billplz v3 x_signature for a flat key/value map.
// Keys are sorted case-insensitively ascending; each pair is key+value;
// pairs are joined with "|" and signed with the configured x_signature_key.
func Sign(values map[string]string, key string) string {
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return strings.ToLower(keys[i]) < strings.ToLower(keys[j])
	})

	pairs := make([]string, 0, len(keys))
	for _, k := range keys {
		pairs = append(pairs, k+values[k])
	}
	msg := strings.Join(pairs, "|")

	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(msg))
	return hex.EncodeToString(mac.Sum(nil))
}

// Verify checks the provided x_signature against the computed value.
func Verify(values map[string]string, key, signature string) bool {
	expected := Sign(values, key)
	return hmac.Equal([]byte(expected), []byte(signature))
}
