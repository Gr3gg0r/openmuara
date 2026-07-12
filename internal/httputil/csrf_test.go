package httputil

import (
	"context"
	"testing"
)

func TestCSRFTokenContextRoundTrip(t *testing.T) {
	tok := "test-csrf-token"
	ctx := WithCSRFToken(context.Background(), tok)
	got, ok := CSRFTokenFromContext(ctx)
	if !ok {
		t.Fatal("expected token to be present in context")
	}
	if got != tok {
		t.Fatalf("token mismatch: want %q, got %q", tok, got)
	}
}

func TestCSRFTokenFromContextMissing(t *testing.T) {
	_, ok := CSRFTokenFromContext(context.Background())
	if ok {
		t.Fatal("expected no token in empty context")
	}
}
