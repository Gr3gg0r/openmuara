package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"
)

// Config contains the settings needed to start an HTTP server.
type Config struct {
	Host    string
	Port    int
	TLSCert string
	TLSKey  string
	Handler http.Handler
}

// Address returns the host:port listening address.
func (c Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// Server wraps an http.Server with graceful shutdown support.
type Server struct {
	httpServer *http.Server
	mu         sync.RWMutex
	listener   net.Listener
	tlsCert    string
	tlsKey     string
}

// New creates a new Server from config.
func New(cfg Config) *Server {
	srv := &http.Server{
		Addr:              cfg.Address(),
		Handler:           cfg.Handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	if cfg.TLSCert != "" && cfg.TLSKey != "" {
		srv.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	return &Server{
		httpServer: srv,
		tlsCert:    cfg.TLSCert,
		tlsKey:     cfg.TLSKey,
	}
}

// Addr returns the actual listening address. If the server has not started
// listening yet, it returns the configured address. When the configured port
// is 0 (auto-allocate) and no listener is ready yet, it returns an empty
// string so callers can wait for the bound address.
func (s *Server) Addr() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.listener != nil {
		return s.listener.Addr().String()
	}
	if s.httpServer.Addr == "" {
		return ""
	}
	_, port, err := net.SplitHostPort(s.httpServer.Addr)
	if err == nil && port == "0" {
		return ""
	}
	return s.httpServer.Addr
}

// BaseURL returns a URL for the listening address.
func (s *Server) BaseURL() string {
	scheme := "http"
	if s.tlsCert != "" && s.tlsKey != "" {
		scheme = "https"
	}
	return scheme + "://" + s.Addr()
}

// ListenAndServe starts the server and blocks until ctx is cancelled.
func (s *Server) ListenAndServe(ctx context.Context) error {
	ln, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", s.httpServer.Addr, err)
	}
	s.mu.Lock()
	s.listener = ln
	s.mu.Unlock()

	errCh := make(chan error, 1)
	go func() {
		scheme := "http"
		if s.tlsCert != "" && s.tlsKey != "" {
			scheme = "https"
		}
		slog.Info("starting http server", "addr", s.Addr(), "scheme", scheme)
		var serveErr error
		if s.tlsCert != "" && s.tlsKey != "" {
			serveErr = s.httpServer.ServeTLS(ln, s.tlsCert, s.tlsKey)
		} else {
			serveErr = s.httpServer.Serve(ln)
		}
		if serveErr != nil && serveErr != http.ErrServerClosed {
			errCh <- serveErr
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.httpServer.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}
