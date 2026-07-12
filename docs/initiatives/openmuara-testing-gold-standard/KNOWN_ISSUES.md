> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Testing Gold Standard — Known Issues & Boundaries

> **Purpose:** Prevent wasted effort on pre-existing testability problems or out-of-scope areas.

---

## Pre-existing Testability Issues

These issues are in-scope for this initiative and will be addressed in the listed prompts.

| ID | Issue | Location | Impact | Planned Fix |
|----|-------|----------|--------|-------------|
| K01 | Package-level provider registry default | `internal/provider/registry.go` | Tests cannot isolate providers | P21 — inject registry |
| K02 | Package-level plugin registry default | `internal/plugin/registry.go` | Tests cannot isolate plugins | P21 — inject registry |
| K03 | CLI start constructs dependencies inline | `internal/cli/start.go` | Cannot unit-test startup logic | P21 — `StartDeps` struct |
| K04 | Server starts real `http.Server` directly | `internal/server/server.go` | Cannot bind to port 0 easily | P21 — accept listener/port 0 |
| K05 | Global slog default side effect | `internal/cli/start.go` | Pollutes test logs | P21 — accept logger |
| K06 | Smoke test uses fixed port 9000 | `scripts/smoke-test.sh` | Can collide on busy machines | P25 — random port |
| K07 | No shared test utilities | scattered tests | Inconsistent setup | P22 — `internal/testutil` |
| K08 | No contract/golden tests | provider packages | Emulation drift risk | P24 — golden files |
| K09 | Many CLI functions 0% covered | `internal/cli/*.go` | Low confidence in CLI | P23 — backfill |

---

## Out-of-Scope Areas

| Area | Reason | Boundary |
|------|--------|----------|
| Production load testing | Local emulator scope | Do not implement |
| Real provider sandbox testing | Out of project vision | Document only |
| v2-frozen providers (App Store, Play Store, RevenueCat) | Hard frozen | Do not add contract tests for these |

---

## How to Use This File

1. Before starting a prompt, scan this file for landmines.
2. If you hit an unrelated pre-existing bug, log it here and move on.
3. Update status when a known issue is fixed.
