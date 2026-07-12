> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Coverage Audit Risk Register

> **Created:** 2026-07-08
> **Last Updated:** 2026-07-09
> **Status:** ✅ Planning complete

---

## Risk matrix

| ID | Risk | Likelihood | Impact | Score | Treatment | Owner |
|----|------|------------|--------|-------|-----------|-------|
| R01 | Coverage target is unrealistic for thin wrappers / generated packages | Medium | Low | Low | Accept with documented exemptions | AI Agent |
| R02 | Tests added only for coverage are brittle | Medium | Medium | Medium | Mitigate through review focus on behavior tests | AI Agent |
| R03 | Dashboard tests miss integration paths (API + state + UI) | Medium | High | High | Mitigate with MSW/stubs and user-flow tests | AI Agent |
| R04 | Strict per-package thresholds slow down delivery | Low | Medium | Low | Mitigate by phasing in thresholds | AI Agent |
| R05 | Coverage tooling breaks with Vitest/Vite major updates | Low | Low | Low | Mitigate by pinning provider to Vitest major | AI Agent |
| R06 | CI coverage gates become flaky due to non-deterministic tests | Low | High | Medium | Mitigate with `-count=1`, mocked time, stable test data | AI Agent |
| R07 | Overall coverage hides weak packages | High | High | High | Mitigate with per-package floors and lowest-package reporting | AI Agent |

*Score = Likelihood × Impact mapped to Low/Medium/High.*

## Detailed risks

### R01 — Unrealistic targets for thin wrappers

Some packages (`internal/ui`, `internal/provider/factory`, `examples/checkout-store`) are thin wrappers, generated code, or examples. Requiring 80% coverage adds little value.

**Mitigation:**
- Maintain `coverage-exemptions.yml` with rationale and review dates.
- Set realistic floors for embedding layers (70%) and exclude examples.

### R02 — Brittle tests added for coverage

Developers may write tests that assert implementation details just to hit a number.

**Mitigation:**
- Review tests for behavior and contract, not line coverage.
- Reject tests that break on harmless refactors.
- Include this rule in `CONTRIBUTING.md`.

### R03 — Dashboard integration paths untested

Unit tests may cover components in isolation but miss API-error → UI-state paths.

**Mitigation:**
- Use MSW or stub `fetch` in `tests/setup.ts`.
- Add user-flow tests for: load ledger → create charge → see webhook.
- Keep Playwright E2E tests as a backstop.

### R04 — Thresholds slow delivery

New features in undertested areas may drop coverage below the gate.

**Mitigation:**
- Phase in dashboard thresholds: 60/55/55/55 now, 70/65/65/65 later.
- Allow temporary threshold lowering with a documented recovery plan.

### R05 — Coverage tooling incompatibility

`@vitest/coverage-v8` major version may drift from `vitest`.

**Mitigation:**
- Pin to `^2.0.0` to match current Vitest major.
- Update both together via Dependabot grouping.

### R06 — Flaky CI gates

Time-based or race-prone tests can make coverage non-deterministic.

**Mitigation:**
- Use `go test -count=1` in CI.
- Mock `time.Now` in Go tests where needed.
- Avoid real timers in dashboard tests; use `vi.useFakeTimers()`.

### R07 — Overall metric hides weak packages

A healthy overall number can mask a package at 45%.

**Mitigation:**
- Enforce `scripts/check-coverage-per-package.sh`.
- Report the lowest package in every PR comment.

## Treatment summary

| Treatment | Risks |
|---|---|
| Accept | R01 |
| Mitigate | R02, R03, R04, R05, R06, R07 |
| Transfer | None |
| Avoid | None |

## Monitoring

- Review coverage trends in PR comments.
- Revisit exemption list quarterly or before each release.
- If three consecutive PRs fail a gate, treat it as a signal to adjust thresholds or tests, not to disable the gate.
