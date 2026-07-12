package server

import (
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	// maxRateLimiterEntries caps the number of tracked IPs to prevent
	// unbounded memory growth.
	maxRateLimiterEntries = 10_000
	// defaultRequestsPerMinute is the default request threshold.
	defaultRequestsPerMinute = 100
)

// RateLimiterConfig holds in-memory rate limiting settings.
type RateLimiterConfig struct {
	Enabled           bool
	RequestsPerMinute int
	AdminOnly         bool
}

// bucket tracks a token bucket for a single IP.
type bucket struct {
	tokens    float64
	lastSeen  time.Time
	updatedAt time.Time
}

// RateLimiter is an in-memory token-bucket rate limiter with bounded size
// and TTL eviction. It is safe for concurrent use.
type RateLimiter struct {
	mu        sync.RWMutex
	buckets   map[string]*bucket
	reqPerMin float64
	ttl       time.Duration
	maxSize   int
}

// NewRateLimiter creates a rate limiter. If requestsPerMinute is <= 0,
// defaultRequestsPerMinute is used.
func NewRateLimiter(requestsPerMinute int) *RateLimiter {
	rpm := float64(requestsPerMinute)
	if rpm <= 0 {
		rpm = defaultRequestsPerMinute
	}
	return &RateLimiter{
		buckets:   make(map[string]*bucket),
		reqPerMin: rpm,
		ttl:       2 * time.Minute,
		maxSize:   maxRateLimiterEntries,
	}
}

// Allow reports whether a request from the given IP is allowed.
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.evictExpiredLocked()

	now := time.Now()
	b, ok := rl.buckets[ip]
	if !ok {
		if len(rl.buckets) >= rl.maxSize {
			// Map is full; reject new IPs to prevent growth.
			return false
		}
		b = &bucket{tokens: rl.reqPerMin - 1, lastSeen: now, updatedAt: now}
		rl.buckets[ip] = b
		return true
	}

	elapsed := now.Sub(b.updatedAt).Minutes()
	b.tokens = minFloat64(b.tokens+elapsed*rl.reqPerMin, rl.reqPerMin)
	b.updatedAt = now
	b.lastSeen = now

	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

// evictExpiredLocked removes stale buckets when the map is at capacity.
// Caller must hold rl.mu.
func (rl *RateLimiter) evictExpiredLocked() {
	if len(rl.buckets) < rl.maxSize {
		return
	}
	cutoff := time.Now().Add(-rl.ttl)
	for ip, b := range rl.buckets {
		if b.lastSeen.Before(cutoff) {
			delete(rl.buckets, ip)
		}
	}
}

// RateLimitMiddleware returns a middleware that rate-limits requests.
func RateLimitMiddleware(cfg RateLimiterConfig, limiter *RateLimiter) Middleware {
	if !cfg.Enabled || limiter == nil {
		return func(next http.Handler) http.Handler { return next }
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cfg.AdminOnly && !isAdminRoute(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			ip := clientIP(r)
			if !limiter.Allow(ip) {
				logSecurityEventFromRequest(r, SecurityEventRateLimit, ip)
				headers := w.Header()
				headers.Set("Content-Type", "application/json")
				headers.Set("Retry-After", "60")
				w.WriteHeader(http.StatusTooManyRequests)
				_, _ = w.Write([]byte(`{"error":"rate limit exceeded"}`))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// clientIP returns the remote IP address, preferring X-Forwarded-For only
// when the direct peer is a trusted local address. For OpenMuara's local-first
// model, the direct remote address is used by default to prevent spoofing.
func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func minFloat64(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
