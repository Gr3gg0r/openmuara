# Prompt 04 — Provider Registry Refactor

## Goal
Unify the provider registry and declarative plugin system into a single runtime model.

## Acceptance Criteria
- [ ] Single registry used at runtime
- [ ] Providers can be loaded from Go code or from `plugins/<provider>/gateway.yml`
- [ ] Provider interface supports:
  - `Init(cfg map[string]any) error`
  - `Routes() []Route`
  - `PayloadBuilder() func(...)`
  - `PayloadHeaders() (optional)`
  - `EscapeHandler() http.Handler`
- [ ] `openmuara start` loads both built-in providers and discovered plugins
- [ ] No duplicate route registration
- [ ] Dormant plugin code removed or integrated

## Files to Create/Change
- `internal/provider/provider.go` — interface
- `internal/provider/registry.go` — discovery
- `internal/plugin/loader.go` — adapter to provider interface
- `internal/cli/start.go` — unified load path
- `internal/server/router.go` — single route registration path

## Response Shape
Return:
1. Unified provider interface
2. Load order (built-ins vs plugins)
3. How conflicts are resolved
4. Files removed/merged

## Test Notes
- `go test ./internal/provider/... ./internal/plugin/...`
- Start server with Fawry built-in + Stripe plugin, verify routes
