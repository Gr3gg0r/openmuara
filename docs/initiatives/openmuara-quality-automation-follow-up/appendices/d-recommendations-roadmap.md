> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This document is subordinate to it.**

# Appendix D — Recommendations Roadmap

> **Updated:** 2026-07-06

This appendix collects every suggestion, enhancement, good-to-have, and gold-standard practice identified while reviewing the bug-hunt recommendations and the follow-up initiative. Items are grouped by area and prioritized. Nothing in this appendix is executed until it is promoted into a prompt or a new initiative.

## How to use this roadmap

- **P01–P06 scope:** Items marked `Core` should be considered for the current initiative.
- **Future initiative:** Items marked `Future` are kept here as seeds for the next quality cycle.
- **Deprioritized:** Items marked `Low` are recorded so they are not lost, but they are not drivers for this pass.

---

## CI / CD Automation

| # | Recommendation | Priority | Rationale | Target | Note |
|---|----------------|----------|-----------|--------|------|
| R01 | Visual baseline diff as a required CI check | Core | Catches unintended UI changes before merge. | P01 | Start non-blocking; promote after stability. |
| R02 | Mutation testing job in CI with a score threshold | Core | Ensures tests actually detect bugs. | P02 | Target changed packages first. |
| R03 | Coverage-regression gate that fails PRs | Core | Prevents silent coverage drops on changed modules. | P03 | Compare against target branch baseline. |
| R04 | Scheduled bug-hunt prep issue before each release | Core | Institutionalizes recurring quality sprints. | P05 | Must be idempotent. |
| R05 | KNOWN_ISSUES sync check in CI | Core | Keeps user-facing known issues honest. | P05 | Start as warning. |
| R06 | CI job summary with gate status and links | Future | Faster developer feedback and debugging. | Post-P06 | Use GitHub Actions job summaries. |
| R07 | Required status checks documented in branch protection rules | Future | Prevents merges that bypass gates. | Post-P06 | Repo admin action. |
| R08 | Nightly full quality matrix (not just PR) | Future | Catches environmental/tooling drift early. | Future | Could use scheduled workflow. |
| R09 | Pin all CI action versions and enable Dependabot for actions | Future | Reduces supply-chain surprises. | Future | Security best practice. |
| R10 | CI artifact retention policy for screenshots and coverage | Future | Keeps storage costs bounded. | Future | 30-day default is usually fine. |

---

## Testing & Quality

| # | Recommendation | Priority | Rationale | Target | Note |
|---|----------------|----------|-----------|--------|------|
| R11 | Provider-wide `errcode` adoption | Core | Stable error codes across all providers and APIs. | P04 | Additive only. |
| R12 | Golden-file auto-update guard | Future | Prevents accidental golden-file drift. | Future | CI check that fails if golden files change without explicit flag. |
| R13 | Property-based tests for transaction state machine | Core | Already started in bug hunt; expand transitions. | P02 / P04 | Use fuzzing + invariants. |
| R14 | Shuffle tests (`-shuffle=on`) in CI | Future | Exposes order-dependent tests. | Future | May need test fixes first. |
| R15 | Table-driven test linter / consistency check | Low | Keeps provider tests readable and uniform. | Future | Tooling optional. |
| R16 | Mock/stub contract tests for external HTTP clients | Future | Prevents drift in provider HTTP shape. | Future | Aligns with conformance tests. |
| R17 | Replay-based regression tests for webhooks | Future | Catches payload-builder regressions. | Future | Store sample payloads. |
| R18 | Test-timing regression check | Low | Flags newly slow tests. | Future | Requires historical data. |

---

## Visual / UI / UX

| # | Recommendation | Priority | Rationale | Target | Note |
|---|----------------|----------|-----------|--------|------|
| R19 | Visual baseline per theme (light/dark) | Core | Catches theme-specific regressions. | P01 extension | Second set of snapshots. |
| R20 | Visual baseline per viewport size | Future | Catches responsive layout regressions. | Future | Mobile and desktop. |
| R21 | Mask dynamic data generically instead of per-test CSS | Core | Makes future visual tests easier to stabilize. | P01 | Add a shared `data-visual-mask` helper. |
| R22 | Component-level Storybook/Vitest snapshot tests | Future | Faster feedback than full-page Playwright. | Future | Optional if visual baseline is stable. |
| R23 | Axe-core serious-violation gate in CI | Future | Automates a11y regression detection. | Future | Existing `a11y:contrast` is a good start. |
| R24 | Keyboard-navigation regression test | Future | Ensures shortcuts continue to work. | Future | Playwright or unit test. |
| R25 | Focus-management audit after modal/dialog additions | Low | Accessibility polish. | Future | Per-feature check. |

---

## Error Handling & Observability

| # | Recommendation | Priority | Rationale | Target | Note |
|---|----------------|----------|-----------|--------|------|
| R26 | `errcode` exposed in API error responses | Core | Consumers can rely on stable codes. | P04 | JSON field `error_code`. |
| R27 | Structured logs include `errcode` | Future | Easier log-based alerting and debugging. | Future | Add to slog attributes. |
| R28 | Trace-ID propagation through provider handlers | Future | End-to-end request correlation. | Future | Already in some paths; complete it. |
| R29 | Request/response logging guard for PII | Future | Prevent secrets leaking in logs. | Future | Skip payload bodies for sensitive endpoints. |
| R30 | Health/ready endpoint telemetry for gate failures | Low | Faster incident triage. | Future | Not a local-first requirement. |

---

## Process & Documentation

| # | Recommendation | Priority | Rationale | Target | Note |
|---|----------------|----------|-----------|--------|------|
| R31 | Update `runbooks/quality-gates.md` with every new gate | Core | Contributors must know how to run and fix gates. | P06 | Mandatory. |
| R32 | `CONTRIBUTING.md` section on visual baseline updates | Future | Onboarding for UI contributors. | Future | Link to runbook. |
| R33 | Release checklist that includes all gates | Future | Prevents shipping with disabled checks. | Future | Tie to version bump script. |
| R34 | Decision log for every gate threshold change | Core | Thresholds are explicit and reviewable. | All prompts | Update `DECISIONS.md`. |
| R35 | Post-mortem template for gate flakes | Future | Learn from flaky failures. | Future | Keep in `.github/` or runbooks. |
| R36 | Initiative handoff template improved from this cycle | Low | Faster setup for next quality initiative. | Future | Refine `HANDOFF.md` format. |

---

## Security & Hardening

| # | Recommendation | Priority | Rationale | Target | Note |
|---|----------------|----------|-----------|--------|------|
| R37 | `gosec` findings treated as required | Future | Currently runs with `-no-fail`. | Future | Fix or suppress findings first. |
| R38 | Secret-scanning pre-commit hook | Future | Catches secrets before commit. | Future | `gitleaks` pre-commit. |
| R39 | Dependency vulnerability gate for UI devDependencies | Future | `npm audit --production` already covers prod deps. | Future | May be noisy. |
| R40 | SBOM generation on release | Low | Supply-chain transparency. | Future | Use `go version -m` or `syft`. |

---

## Performance & Efficiency

| # | Recommendation | Priority | Rationale | Target | Note |
|---|----------------|----------|-----------|--------|------|
| R41 | CI parallelization strategy documented | Future | Keep feedback under 10 minutes. | Future | Job dependencies and caching. |
| R42 | Cache Go modules and `golangci-lint` across jobs | Future | Faster CI. | Future | Already partially done; review. |
| R43 | Bundle-size gate tightened over time | Future | Protects dashboard performance. | Future | Lower thresholds as bundle shrinks. |
| R44 | Visual baseline job runs only on dashboard changes | Core | Avoids unnecessary CI minutes. | P01 | Path filter. |
| R45 | Mutation testing runs only on Go package changes | Core | Avoids unnecessary CI minutes. | P02 | Path filter. |

---

## Promotion Rules

1. A `Core` item is promoted into a prompt when it is selected for this initiative.
2. A `Future` item is promoted into a new initiative when the current one is complete and the item has a clear owner.
3. A `Low` item stays in this roadmap until it becomes a blocker or is explicitly rejected with rationale in `DECISIONS.md`.
