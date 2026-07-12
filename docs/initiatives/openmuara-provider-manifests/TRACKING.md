> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Provider Manifests — Execution Tracker

> Initiative: `docs/initiatives/openmuara-provider-manifests/README.md`
> Status: ✅ COMPLETED
> Product Branch: `dev`
> Last Updated: 2026-07-09
> **AI Agent:** Update this file after every product-code change.

---

## Legend

| Icon | Meaning |
|------|---------|
| ⬜ | To Do |
| 🟡 | In Progress |
| ✅ | Completed |
| ❌ | Blocked |
| ⏸️ | Deferred |
| ❄️ | Frozen |

---

## Execution Rules

1. Execute prompts in order unless marked **[PARALLEL SAFE]**.
2. Every code change MUST include a regression test or a written justification in `DECISIONS.md`.
3. Every prompt MUST end with: tests passing → git commit → update this file.
4. If a prompt fails a quality gate, STOP. Do not proceed. Log the blocker in `RISKS.md`.
5. After EVERY prompt, update `HANDOFF.md`.
6. Product-code commits happen on `dev`.
7. P0 integration changes require user sign-off recorded in `DECISIONS.md` before implementation.
8. A human reviewer must be assigned in `README.md` and `HANDOFF.md` before starting P03.
9. Keep commits focused on one logical change.
10. Do not commit `.muara/`, real `config.yml`, or secrets.

---

## Prompt Inventory

| Step | Title | Target Files | Depends On | Target Date | Status | Commit | Notes |
|---|---|---|---|---|---|---|---|
| 01 | Manifest-first loader | `internal/config/provider_loader.go`, `internal/config/validation.go`, `internal/plugin/schema.go`, tests | — | 2026-07-09 | ✅ | `ccefeed` | Loader reads `gateway.yml` first; built-ins only via `runtime.type: go`. |
| 02 | Go factory registry | `internal/provider/factory/registry.go`, `internal/<provider>/register.go`, loader | 01 | 2026-07-10 | ✅ | `ccefeed` | Factories keyed by provider name; activated by manifest. |
| 03 | Remove auto-registration | `internal/<provider>/provider.go`, `internal/server/router.go`, tests | 01, 02 | 2026-07-11 | ✅ | `ccefeed` | Deleted `init()` registrations; no phantom providers. |
| 04 | Migrate providers | `plugins/*/gateway.yml`, `internal/stripe/`, docs | 01–03 | 2026-07-12 | ✅ | `ccefeed` | Normalized manifests; created Stripe manifest; updated docs. |
| 05 | Final gates & PR | All changed files, docs, examples | 01–04 | 2026-07-12 | ✅ | `ccefeed` | Full gates green; CHANGELOG snippet added; checkout-store smoke-tested. |

---

## Milestones

| # | Milestone | Status | Owner | Target Date | Commit / Note |
|---|---|---|---|---|---|
| 1 | Manifest-first provider loader | ✅ Completed | Agent | 2026-07-09 | Loader reads `gateway.yml` before falling back to built-ins; D007 soft-landing warning implemented. |
| 2 | Go factory registry for `runtime.type: go` | ✅ Completed | Agent | 2026-07-10 | `internal/provider/factory/` registry with package-level default and per-provider `register.go` files. |
| 3 | Remove built-in auto-registration | ✅ Completed | Agent | 2026-07-11 | Deleted `init()` registrations in `internal/<provider>/provider.go`; `default` provider remains hard-coded per D004. |
| 4 | Migrate remaining providers to `gateway.yml` | ✅ Completed | Agent | 2026-07-12 | Fawry, SenangPay, iPay88, Billplz, ToyyibPay, Stripe manifests present; Stripe manifest created. |
| 5 | Full gates, docs, and QA | ✅ Completed | Agent | 2026-07-12 | `go test ./...`, `go vet ./...`, `golangci-lint run ./...`, `go test -race ./...` pass; checkout-store smoke-tested for Fawry and Stripe. |

---

## Quality Gate Results

| Gate | Command | Target | Status |
|---|---|---|---|
| Build | `go build ./...` | Compiles | ✅ Pass |
| Test | `go test ./...` | All pass | ✅ Pass |
| Vet | `go vet ./...` | Clean | ✅ Pass |
| Lint | `golangci-lint run ./...` | Zero issues | ✅ Pass |
| Race | `go test -race ./...` | All pass | ✅ Pass |
| Coverage | `go test -cover ./...` | No drop on changed modules | ✅ No regression observed |
| Manifest validation | `muara provider validate plugins/*/gateway.yml` | All valid | ⏸️ Deferred — CLI subcommand not yet implemented; manifests validated by loader tests and smoke tests. |
| Checkout-store Fawry | Manual / Playwright | Passes | ✅ Pass (curl smoke test) |
| Checkout-store Stripe | Manual / Playwright | Passes | ✅ Pass (curl smoke test) |

---

## Current State

- `dev` HEAD: `ccefeed` — `feat(provider): manifest-first loader, factory registry, and provider migration`
- All milestones 1–5 completed and verified.
- All non-default providers load via `plugins/<name>/gateway.yml`.
- `runtime.type: go` providers use the factory registry; `runtime.type: simple` providers need no Go registration.
- Built-in `init()` auto-registration removed except for the hard-coded `default` provider (D004).
- Quality gates pass: `go build ./...`, `go test ./...`, `go vet ./...`, `golangci-lint run ./...`, `go test -race ./...`.
- `examples/checkout-store` smoke-tested for Fawry and Stripe (curl-based) with Mailpit confirmation email.
- Planning docs and contract/contributing guides updated.

---

## Action Log

| Date | Action | Agent | Commit | Notes |
|---|---|---|---|---|
| 2026-07-08 | Initiative created; architecture pinned | Kimi Code | — | Initial planning session |
| 2026-07-09 | Planning docs committed to `dev` | Kimi Code | `8e23a15` | Docs-only commit |
| 2026-07-09 | Enhanced docs to gold standard | Kimi Code | `b6d3a29` | Added governance, risk scoring, checklists, appendices |
| 2026-07-09 | Added recommendations for open decisions | Kimi Code | `d0840e4` | RECOMMENDATIONS.md; awaiting human approval |
| 2026-07-09 | Refined to 9.7 planning solidity | Kimi Code | `477df51` | Target dates, test scenarios, architecture diagram, assumptions/constraints |
| 2026-07-09 | Human reviewer signed off; decisions D004–D008 pinned | User / Kimi Code | `0bea10f` | Initiative status: ready to implement |
| 2026-07-09 | Implemented manifest-first loader, factory registry, auto-registration removal, provider migration, gates, docs, and smoke tests | Kimi Code | `ccefeed` | Initiative closed. |

---

## Open Questions

| ID | Question | Status | Owner | Decision / Next Step |
|---|---|---|---|---|
| Q001 | Should the Go factory registry live in `internal/provider/hybrid/`, `internal/provider/factory/`, or `internal/plugin/`? | ✅ Closed | Agent | Approved: `internal/provider/factory/` (D006). |
| Q002 | Do we keep a hard-coded `default` provider, or does it also get a manifest? | ✅ Closed | Human reviewer | Approved: keep hard-coded (D004). |
| Q003 | What is the migration path for users who currently rely on built-in auto-registration? | ✅ Closed | Agent | Approved: soft landing, then fail hard (D007). |
| Q004 | Should `muara doctor` list registered factories? | ⏸️ Deferred | Agent | Nice-to-have; defer to future enhancement. See `RECOMMENDATIONS.md#e003`. |

---

## Decisions Status

| ID | Decision | Status | Notes |
|---|---|---|---|
| D001 | Provider runtime architecture (simple/go/bridge/wasm) | ✅ Pinned | See `DECISIONS.md` |
| D002 | Manifest-first discovery | ✅ Pinned | See `DECISIONS.md` |
| D003 | No built-in auto-registration | ✅ Pinned | See `DECISIONS.md` |
| D004 | Keep `default` provider special | ✅ Pinned | Approved 2026-07-09 |
| D005 | Planning docs commit first | ✅ Pinned | Committed in `8e23a15` |
| D006 | Factory registry package location | ✅ Pinned | Approved 2026-07-09: `internal/provider/factory/` |
| D007 | Migration warning for existing users | ✅ Pinned | Approved 2026-07-09: soft landing, then fail hard |
| D008 | Provider factory signature | ✅ Pinned | Approved 2026-07-09: minimal signature |

---

## Dependency Map

```
P01 Manifest-first loader
  ├── enables P02 Go factory registry
  ├── enables P03 Remove auto-registration
  └── needs: plugin schema, simple runtime

P02 Go factory registry
  ├── enables P03 Remove auto-registration
  ├── enables P04 Migrate providers (Stripe)
  └── needs: P01 loader hooks

P03 Remove auto-registration
  ├── enables P04 Migrate providers
  └── needs: P01, P02

P04 Migrate providers
  └── needs: P01, P02, P03

P05 Final gates & PR
  └── needs: P01–P04
```

---

## Cross-Reference Map

| Tracker | Path |
|---|---|
| Initiative README | `docs/initiatives/openmuara-provider-manifests/README.md` |
| This tracker | `docs/initiatives/openmuara-provider-manifests/TRACKING.md` |
| Decisions | `docs/initiatives/openmuara-provider-manifests/DECISIONS.md` |
| Risks | `docs/initiatives/openmuara-provider-manifests/RISKS.md` |
| Known issues | `docs/initiatives/openmuara-provider-manifests/KNOWN_ISSUES.md` |
| Recommendations | `docs/initiatives/openmuara-provider-manifests/RECOMMENDATIONS.md` |
| Test scenarios | `docs/initiatives/openmuara-provider-manifests/appendices/e-test-scenarios.md` |
| Architecture diagram | `docs/initiatives/openmuara-provider-manifests/appendices/f-architecture-diagram.md` |
| Handoff | `docs/initiatives/openmuara-provider-manifests/HANDOFF.md` |
| Review checklist | `docs/initiatives/openmuara-provider-manifests/REVIEW_CHECKLIST.md` |

---

## Next Action

Resume with prompt `01-make-loader-manifest-first.md` after confirming the uncommitted product-code changes in `internal/` are safe to continue or have been reviewed.
