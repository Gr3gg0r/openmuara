# OpenMuara v1 Solid Gold — Decision Log

## D001 — Additive changes only

**Decision:** This initiative makes no breaking config, CLI, or API changes.
New fields are optional and new endpoints are additive.

**Rationale:** v1 is feature-complete; the goal is polish, not disruption.

**Date:** 2026-07-01

---

## D002 — P03 requires explicit approval

**Decision:** Prompt 03 (trace-ID propagation, CLI inspect, pprof) touches P0
webhook dispatch and provider integration logic. The user must approve the
approach before implementation.

**Rationale:** `AGENTS.md` requires sign-off for P0 integration changes.

**Date:** 2026-07-01

---

## D003 — Coverage target stays at 80%

**Decision:** The existing 80% coverage threshold is sufficient; the goal is to
eliminate packages below it, not to chase 100% globally.

**Rationale:** Diminishing returns beyond 80% for tooling/hygiene code; focus
on behavior-critical packages.

**Date:** 2026-07-01

---

## D004 — Test-only backfills and minor main.go refactor

**Decision:** P02 commits are test-only, grouped per package. `cmd/muara/main.go`
was refactored to use injectable `execute`/`exitFunc` variables so the error
exit path can be unit-tested without spawning a subprocess.

**Rationale:** Keeps diffs reviewable and reaches the 80% coverage target for
every package without changing runtime behavior.

**Date:** 2026-07-01

---

## D005 — Linter findings and dead-code handling

**Decision:** Enable `gosec`, `staticcheck`, `ineffassign`, `unparam`, `errcheck`,
`misspell`, `revive`, and `unused` in `.golangci.yml`. Fix trivial findings
across the codebase. Suppress false positives with `#nosec` comments that
include a brief justification. Remove genuinely unused code (e.g.
`ParseValidationPort`) only after confirming no callers. Advisory file/function
size warnings remain accepted debt.

**Rationale:** Stronger static analysis prevents drift, but v1 architecture
should not be refactored just to satisfy style/size checks.

**Date:** 2026-07-01
