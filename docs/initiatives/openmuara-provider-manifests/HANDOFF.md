> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Provider Manifests — Session Handoff

## Context

This initiative makes every non-default provider discoverable and configurable through `plugins/<name>/gateway.yml`. The architecture is pinned in `DECISIONS.md`.

## Current Branch

- `dev`

## Known Working Tree State

- All product-code changes for milestones 1–5 are committed on `dev`.
- Planning docs under `docs/initiatives/openmuara-provider-manifests/` are updated and closed.
- Base commit for implementation: `8e23a15` (`docs(initiatives): add OpenMuara provider manifests initiative`).
- Gold-standard planning docs were added/enhanced in the most recent session.

## How to Resume

1. Run `git status` and note any unstaged product-code changes.
2. Run gates (`go test ./...`, `go vet ./...`, `golangci-lint run ./...`).
3. Read `TRACKING.md` for the current milestone and action log.
4. Pick up the next uncompleted prompt in `prompts/`.
5. Update `TRACKING.md` and this file after each prompt.

## Recent Decisions

- D001–D003: Architecture pinned (simple/go/bridge/wasm, manifest-first, no auto-registration).
- D004–D008: Approved and pinned on 2026-07-09.
- D005: Planning docs committed first.

## Approved Choices

| Decision | Approved Choice |
|---|---|
| D004 — `default` provider | Keep hard-coded as internal fallback. |
| D006 — Factory registry location | `internal/provider/factory/`. |
| D007 — Migration path | Soft landing: warning this release, fail hard next release. |
| D008 — Factory signature | `type Factory func(cfg map[string]any) (provider.Provider, error)`. |

## Reviewer Sign-off

| Reviewer | Signed Off | Date |
|---|---|---|
| Human reviewer (user) | ✅ | 2026-07-09 |

P03 auto-registration removal completed with human sign-off captured in `DECISIONS.md` / tracker.

## Active Risks

- R005: Loader ordering changes may affect existing `.muara/config.yml` users; mitigated by D007 soft-landing warning.

## Blockers

| ID | Blocker | Owner | Next Step |
|---|---|---|---|
| — | None — initiative closed | — | — |

## Common Pitfalls

- Do not mix planning-doc commits with product-code commits.
- `internal/provider/conform/conform_test.go` and `internal/server/providers_test.go` depend on globally registered providers; they will need updates when auto-registration is removed.
- `runtime.type: go` must be implemented before removing auto-registration, or existing Go providers will fail to load.
- Do not expand scope into `bridge` or `wasm` runtimes.

## Contacts

- Owner: AI Agent (Kimi Code)
- Human Reviewer: ___________
- Domain Experts: Fawry team (invited for provider validation)
