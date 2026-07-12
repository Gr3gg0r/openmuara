package httputil

import (
	"strings"
	"testing"
)

func TestValidateWebhookURL(t *testing.T) {
	cases := []struct {
		name     string
		url      string
		hardened bool
		wantErr  string
	}{
		{"empty", "", false, "webhook url is required"},
		{"whitespace only", "   ", false, "webhook url is required"},
		{"valid http", "http://example.com/webhook", false, ""},
		{"valid https", "https://example.com/webhook", false, ""},
		{"invalid scheme", "ftp://example.com/webhook", false, "not allowed"},
		{"missing host", "http:///webhook", false, "host is required"},
		{"localhost non-hardened", "http://127.0.0.1/webhook", false, ""},
		{"localhost hardened", "http://127.0.0.1/webhook", true, "private or internal"},
		{"private ip hardened", "https://10.0.0.1/webhook", true, "private or internal"},
		{"public ip hardened", "https://1.1.1.1/webhook", true, ""},
		{"localhost hostname hardened", "http://localhost/webhook", true, "private or internal"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateWebhookURL(tc.url, tc.hardened)
			if tc.wantErr == "" {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Errorf("expected error containing %q, got %v", tc.wantErr, err)
			}
		})
	}
}
