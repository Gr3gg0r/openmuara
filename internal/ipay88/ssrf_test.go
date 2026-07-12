package ipay88

import (
	"testing"
)

func TestIsPublicURL(t *testing.T) {
	cases := []struct {
		name    string
		url     string
		allowed bool
	}{
		{"public http", "http://example.com/callback", true},
		{"public https", "https://example.com/callback", true},
		{"loopback", "http://127.0.0.1:9000/callback", false},
		{"localhost", "http://localhost/callback", false},
		{"private range", "http://192.168.1.1/callback", false},
		{"link local", "http://169.254.1.1/callback", false},
		{"non http scheme", "ftp://example.com/callback", false},
		{"missing scheme", "example.com/callback", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := IsPublicURL(tc.url)
			if tc.allowed && err != nil {
				t.Errorf("expected %q to be allowed, got error: %v", tc.url, err)
			}
			if !tc.allowed && err == nil {
				t.Errorf("expected %q to be rejected", tc.url)
			}
		})
	}
}
