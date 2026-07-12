package server

import "context"

// Role represents the authorization level of an authenticated request.
type Role string

const (
	// RoleAnonymous is the default role for unauthenticated requests.
	RoleAnonymous Role = ""
	// RoleViewer grants read-only access to the admin dashboard and APIs.
	RoleViewer Role = "viewer"
	// RoleAdmin grants full access to the admin dashboard and APIs.
	RoleAdmin Role = "admin"
)

type roleKey struct{}

// WithRole returns a new context with the given role.
func WithRole(ctx context.Context, role Role) context.Context {
	return context.WithValue(ctx, roleKey{}, role)
}

// RoleFromContext returns the role stored in the context, or RoleAnonymous.
func RoleFromContext(ctx context.Context) Role {
	if r, ok := ctx.Value(roleKey{}).(Role); ok {
		return r
	}
	return RoleAnonymous
}

// IsAdmin reports whether the context has the admin role.
func IsAdmin(ctx context.Context) bool {
	return RoleFromContext(ctx) == RoleAdmin
}
