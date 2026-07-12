# Prompt 01 — Project Bootstrap & Rebrand

## Goal
Complete the rebrand from `muara` to `OpenMuara` while preserving all existing functionality.

## Acceptance Criteria
- [x] Module path updated from `github.com/openmuara/openmuara` to `github.com/openmuara/openmuara`
- [x] Binary name updated to `muara`
- [x] Default workspace directory changed from `.muara/` to `.muara/`
- [x] Config file path changed from `.muara/config.yml` to `.muara/config.yml`
- [x] CLI command renamed from `muara` to `muara`
- [x] README updated with new name and migration note
- [x] `AGENTS.md` updated to reflect OpenMuara conventions
- [x] Build passes: `go build ./...`
- [x] Tests pass: `go test ./...`

## Files Changed
- `go.mod` — module path
- `cmd/openmuara/` → `cmd/muara/`
- `internal/config/config.go` — default paths
- `README.md`
- `AGENTS.md`
- `Taskfile.yml` — command names
- `scripts/smoke-test.sh` — command names

## Response Shape
Return:
1. List of renamed/moved files
2. Any config key mappings that changed
3. Verification commands and their output

## Test Notes
- Run `go build ./...`
- Run `go test ./...`
- Run `./bin/muara version`
