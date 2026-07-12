package ipay88

import (
	"strings"
	"testing"
)

func TestIsPublicURLEdgeCases(t *testing.T) {
	cases := []struct {
		name    string
		rawURL  string
		wantErr string
	}{
		{"unparseable", "http://[::1:invalid/", "parse"},
		{"empty host", "http:///path", "host is required"},
		{"public ip", "http://1.1.1.1/callback", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := IsPublicURL(tc.rawURL)
			if tc.wantErr == "" {
				if err != nil {
					t.Errorf("expected allowed, got %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Errorf("expected error containing %q, got %v", tc.wantErr, err)
			}
		})
	}
}
