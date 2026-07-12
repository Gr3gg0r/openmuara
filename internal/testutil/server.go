package testutil

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// NewTestServer starts an httptest.Server with the given handler and registers
// a cleanup function to close it when the test ends. It returns the server so
// callers can use srv.URL and srv.Client().
func NewTestServer(t *testing.T, handler http.Handler) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return srv
}
