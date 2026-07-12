package stripe

import (
	"testing"
	"time"
)

func FuzzSignVerify(f *testing.F) {
	f.Add([]byte("payload"), "whsec_secret", time.Now().Unix())

	f.Fuzz(func(t *testing.T, payload []byte, secret string, ts int64) {
		if len(payload) == 0 || secret == "" {
			t.Skip()
		}

		timestamp := time.Unix(ts, 0)
		header := SignPayload(payload, secret, timestamp)
		if err := VerifySignature(payload, header, secret); err != nil {
			t.Errorf("fresh signature should verify: %v", err)
		}

		if err := VerifySignature(append(payload, 'x'), header, secret); err == nil {
			t.Error("mutated payload should not verify")
		}
	})
}
