# Prompt 03 — Configuration & Environment

## Goal
Implement a robust configuration system for OpenMuara supporting YAML, environment variables, and validation.

## Acceptance Criteria
- [x] Config loaded from `.muara/config.yml` (or path override)
- [x] Environment variable prefix `MUARA_`
- [x] All provider configs nested under `providers.*`
- [x] Validation errors return clear messages
- [x] Bundled default config embedded in `internal/config.DefaultYAML()`
- [x] No secrets committed to defaults
- [x] Backward-compatible fallback for legacy `.muara/config.yml` (warn, migrate)

## Files Changed
- `internal/config/config.go` — struct, defaults, validation
- `internal/config/config_test.go` — env var tests
- `internal/config/migrate.go` — legacy path warning

## Response Shape
Return:
1. Config struct (YAML/JSON shape)
2. Env var mapping table
3. Validation rules
4. Default config path

## Test Notes
- `go test ./internal/config/...`
- Test env override for nested keys
- Test legacy path warning
