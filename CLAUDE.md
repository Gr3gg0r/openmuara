# CLAUDE.md — OpenMuara

> Project-specific notes for AI coding agents working on OpenMuara.
> See `AGENTS.md` for workspace rules, branch policy, and quality gates.

## Health Stack

- build: `go build ./...`
- vet: `go vet ./...`
- test: `go test ./...`
- lint: `golangci-lint run ./...`
- shell: `shellcheck scripts/*.sh`
