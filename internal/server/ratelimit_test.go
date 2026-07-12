package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiterAllow(t *testing.T) {
	rl := NewRateLimiter(2)
	ip := "127.0.0.1"

	if !rl.Allow(ip) {
		t.Error("first request should be allowed")
	}
	if !rl.Allow(ip) {
		t.Error("second request should be allowed")
	}
	if rl.Allow(ip) {
		t.Error("third request should be blocked")
	}
}

func TestRateLimiterRefill(t *testing.T) {
	rl := NewRateLimiter(60) // bucket capacity is 60 tokens
	ip := "127.0.0.1"

	if !rl.Allow(ip) {
		t.Fatal("first request should be allowed")
	}
	if !rl.Allow(ip) {
		t.Error("second immediate request should be allowed (bucket has 60 tokens)")
	}

	// Exhaust the bucket.
	for i := 0; i < 58; i++ {
		if !rl.Allow(ip) {
			t.Fatalf("request %d should be allowed", i+3)
		}
	}
	if rl.Allow(ip) {
		t.Error("61st request should be blocked")
	}

	time.Sleep(1100 * time.Millisecond)
	if !rl.Allow(ip) {
		t.Error("request after ~1 second refill should be allowed")
	}
}

func TestRateLimiterMaxSize(t *testing.T) {
	rl := NewRateLimiter(10)
	rl.maxSize = 2

	if !rl.Allow("1.1.1.1") {
		t.Error("first IP should be allowed")
	}
	if !rl.Allow("2.2.2.2") {
		t.Error("second IP should be allowed")
	}
	if rl.Allow("3.3.3.3") {
		t.Error("third IP should be rejected when map is full")
	}
}

func TestRateLimitMiddlewareDisabled(t *testing.T) {
	cfg := RateLimiterConfig{Enabled: false}
	handler := RateLimitMiddleware(cfg, nil)(okHandler())

	req := httptest.NewRequest(http.MethodGet, "/_admin", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestRateLimitMiddlewareAdminOnly(t *testing.T) {
	rl := NewRateLimiter(1)
	cfg := RateLimiterConfig{Enabled: true, RequestsPerMinute: 1, AdminOnly: true}
	handler := RateLimitMiddleware(cfg, rl)(okHandler())

	// Provider route should not be rate limited.
	req1 := httptest.NewRequest(http.MethodGet, "/fawry/charge", nil)
	rr1 := httptest.NewRecorder()
	handler.ServeHTTP(rr1, req1)
	if rr1.Code != http.StatusOK {
		t.Errorf("provider route: want %d, got %d", http.StatusOK, rr1.Code)
	}

	// Admin route consumes the bucket.
	req2 := httptest.NewRequest(http.MethodGet, "/_admin", nil)
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusOK {
		t.Errorf("admin first request: want %d, got %d", http.StatusOK, rr2.Code)
	}

	// Second admin request blocked.
	req3 := httptest.NewRequest(http.MethodGet, "/_admin", nil)
	rr3 := httptest.NewRecorder()
	handler.ServeHTTP(rr3, req3)
	if rr3.Code != http.StatusTooManyRequests {
		t.Errorf("admin second request: want %d, got %d", http.StatusTooManyRequests, rr3.Code)
	}
}

func TestClientIP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	if got := clientIP(req); got != "192.168.1.1" {
		t.Errorf("want 192.168.1.1, got %q", got)
	}
}
