package server

import (
	"net/http"
	"net/http/httptest"
)

// newAdminRequest returns an httptest request whose context carries RoleAdmin.
func newAdminRequest(method, target string) *http.Request {
	req := httptest.NewRequest(method, target, nil)
	req = req.WithContext(WithRole(req.Context(), RoleAdmin))
	return req
}

// newViewerRequest returns an httptest request whose context carries RoleViewer.
func newViewerRequest(method, target string) *http.Request {
	req := httptest.NewRequest(method, target, nil)
	req = req.WithContext(WithRole(req.Context(), RoleViewer))
	return req
}
