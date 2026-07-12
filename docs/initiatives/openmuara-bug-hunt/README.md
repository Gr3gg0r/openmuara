> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Bug Hunt

> **Status:** 🟢 Completed | **Started:** 2026-07-06 | **Completed:** 2026-07-06
>
> **Completion note:** All approved enhancement recommendations (E1–E12) from `appendices/b-recommendations.md` were implemented, quality gates passed, and visual sign-off completed. The branch is ready for ongoing use as the baseline quality process.
> **Scope:** Systematically discover, reproduce, triage, and fix bugs across the OpenMuara v1 codebase — with explicit guardrails for UI/UX, security, and the project philosophy — before the next release.
> **Owner:** AI Agent (Kimi Code) | **Human Reviewer:** ___________
> **Target Repo:** `<repo-root>/`
> **Product Branch:** `feat/bug-hunt`
>
> **Why:** The Mailpit-style dashboard redesign (`feat/dashboard-mailpit-redesign`) and recent v1 work need a focused quality pass before promotion to `dev`. This initiative is the safety net: find regressions, latent defects, and edge-case failures in provider emulation, config handling, webhook dispatch, the admin UI, and the dual-port runtime; fix them with minimal, tested changes; and prove the result visually and through quality gates.

---

## Initiative Structure

```
docs/initiatives/openmuara-bug-hunt/
├── README.md              # This file
├── TRACKING.md            # Central execution tracker
├── HANDOFF.md             # Session continuity
├── DECISIONS.md           # Decision log
├── RISKS.md               # Risk register
├── PREREQUISITES.md       # Tools, assumptions, branch base, time-box
├── KNOWN_ISSUES.md        # Deferred bugs
├── REVIEW_CHECKLIST.md    # Pre-PR human review checklist
│
├── prompts/               # Numbered, self-contained execution prompts
│   ├── _template.md
│   ├── 01-reconnaissance.md
│   ├── 02-triage-and-prioritization.md
│   ├── 03-fix-batch-1.md
│   ├── 04-fix-batch-2.md
│   ├── 05-regression-tests-and-quality-gates.md
│   └── 06-visual-sign-off-and-philosophy-check.md
│
├── findings/              # Individual bug reports and visual baselines
│   ├── TEMPLATE.md
│   └── .gitkeep
│
└── appendices/            # Deep-dive reference material
    ├── a-security-checklist.md
    ├── b-recommendations.md
    ├── c-post-bug-hunt-checklist.md
    ├── d-gold-standard-alignment.md
    └── e-bug-register-format.md
```

Planning docs live in `docs/initiatives/openmuara-bug-hunt/` in the root repo.
Product code commits to the `feat/bug-hunt` branch, which was branched from
`feat/dashboard-mailpit-redesign`. Do not commit directly to `main`.

> **Entry point:** Read `PREREQUISITES.md` before starting P01.

---

## Goals

1. **Discover bugs** — use static analysis, test failures, code review, runtime exploration, and visual inspection to find defects.
2. **Reproduce every finding** — each bug must have a minimal reproduction (test, curl command, UI step sequence, or Playwright snapshot).
3. **Triage by severity and philosophy** — classify bugs as P0 (crash/security/data loss), P1 (broken feature or UX regression), P2 (polish/edge case); UI/UX regressions that violate the Mailpit-style contract are elevated.
4. **Fix in batches** — group related fixes; keep each commit focused on one logical change.
5. **Add regression tests** — every code fix must include a test that fails before the fix and passes after.
6. **Protect dashboard invariants** — the Mailpit-style layout (left nav, ledger default, filters, detail pages, provider settings, dual-port runtime) must not regress.
7. **Pass all quality gates** — build, test, lint, vet, frontend tests/build, bundle size, and a11y must remain green.
8. **Document known issues** — anything intentionally not fixed is recorded in `RISKS.md` and `HANDOFF.md` with rationale and target release.
9. **Visual sign-off** — before the branch is declared ready, use Playwright MCP to capture and inspect the dashboard against the project philosophy and the dashboard redesign acceptance criteria.

---

## Philosophy & Priority Alignment

Root `AGENTS.md` and the dashboard redesign initiative establish this priority stack:

**UI > UX > Performance > Usability > Philosophy > Efficiency > Memory size**

For this bug hunt, translate that stack into triage weight:

- Bugs that break **visual clarity, layout, navigation, or the Mailpit-style shell** are treated as high-impact UX issues (P1 or P0 if they block a primary flow).
- Bugs that break **correctness, security, or data integrity** still come first.
- Performance and efficiency fixes are in scope only when they fix a measurable bug or regression; optimization for its own sake is out of scope.
- Memory/bundle size is a constraint, not a driver. We will not shrink the bundle by removing accessibility, filters, or detail views that the user explicitly asked for.

The project philosophy is **local-first, simple, and explicit**:

- No external services, telemetry, or cloud dependencies.
- Provider protocol emulation must remain faithful to documented behavior.
- Every config change is explicit; every secret stays server-side and out of git.
- Fixes are minimal. No speculative refactors. No feature creep.

---

## Conventions (MANDATORY)

### 1. Read `AGENTS.md` first
Root `AGENTS.md` governs branch rules, quality gates, autonomy boundaries, and code style.

### 2. Priority stack
When trade-offs arise, decide in this order:

1. **Correctness** — wrong provider behavior, data loss, or crashes.
2. **Security** — auth bypass, secret leakage, SSRF, injection, unsafe defaults.
3. **Reliability** — flakes, races, startup failures.
4. **UX** — confusing error messages, broken navigation, broken filters/detail pages, accessibility regressions.
5. **Performance** — only after the above are satisfied.
6. **Polish** — cosmetics last.

### 3. P0 integration changes need explicit user sign-off
Per `AGENTS.md` autonomy boundaries, fixes that touch provider emulation logic, webhook signature verification, config persistence schemas, auth/billing/PII flows, or the provider plugin schema contract require user sign-off **before** implementation. Document the sign-off in `DECISIONS.md`.

### 4. One logical change per commit
Each fix gets its own commit with a clear message: `fix(scope): short description (#issue)`.

### 5. Regression test required
Every code fix must include a regression test. Docs-only fixes are exempt. If a fix is not unit-testable, add an integration test; if neither is practical, record the justification in `DECISIONS.md`.

### 6. No speculative refactors
Fix the bug with minimal changes. Do not opportunistically refactor unrelated modules, rename symbols, or change interfaces.

### 7. Quality gates
Every prompt must pass:

- `go build ./...`
- `go test ./...`
- `go vet ./...`
- `golangci-lint run`
- `cd web/dashboard && npm run test:ci`
- `cd web/dashboard && npm run build`
- `node web/dashboard/scripts/check-bundle-size.js`
- `cd web/dashboard && node scripts/a11y-contrast-check.js`

### 8. Definition of done
A bug is fixed only when:
- It is reproduced and understood.
- A minimal fix is implemented.
- A regression test is added.
- Quality gates pass.
- `TRACKING.md` and `HANDOFF.md` are updated.
- If the bug touched UI, the visual sign-off prompt (P06) is re-run before declaring the branch ready.

---

## Out of Scope

- New features or providers.
- Large architectural refactors.
- Performance optimization without a measurable bug or regression.
- Changes to the provider plugin schema contract.
- UI redesign work (covered by the dashboard redesign initiative).
- v2 features (RevenueCat, mobile receipt validation).

---

## Metrics

| Metric | Current | Target | How measured |
|--------|---------|--------|--------------|
| Known bugs filed | 0 | ≥ 5 discovered | `TRACKING.md` bug register |
| Bugs fixed | 0 | ≥ 5 fixed | `TRACKING.md` bug register |
| Regression tests added | 0 | ≥ 5 new tests | Test file diff |
| Quality gates | Passing | Passing | CI / local gate commands |
| Test coverage | Baseline | No drop on changed modules | `coverage.out` / `go test -cover` |
| Visual sign-off | Not run | Pass (P06) | Playwright MCP screenshots + checklist |
| A11y serious violations | 0 | 0 | `npm run test:a11y` / axe-core |

## Success Criteria

- [x] Approved enhancement recommendations E1–E12 implemented with tests.
- [x] All dashboard redesign invariants remain intact.
- [x] All quality gates pass.
- [x] No new warnings from `golangci-lint`, `go vet`, or TypeScript.
- [x] Test coverage does not drop on any changed module.
- [x] `RISKS.md` and `KNOWN_ISSUES.md` list any intentionally deferred items with rationale.
- [x] P06 visual sign-off confirms the dashboard still matches the Mailpit-style philosophy.
- [x] Branch is ready for PR to `dev`.

---

## Proposed Workflow

| Step | Prompt | Goal | Key Output |
|------|--------|------|------------|
| 1 | P01 Reconnaissance | Discover bugs via tests, static analysis, runtime, and visual baseline. | Bug register with ≥5 findings; `findings/` reports; visual baseline. |
| 2 | P02 Triage & Prioritization | Confirm reproducibility, assign severity, group fixes, identify sign-off needs. | Updated bug register with batches; `RISKS.md`; `DECISIONS.md`. |
| 3 | P03 Fix Batch 1 | Fix 2–3 high-impact/low-risk bugs with regression tests. | Fixed bugs + focused commits. |
| 4 | P04 Fix Batch 2 | Fix remaining bugs or defer with rationale. | Fixed/deferred bugs + focused commits. |
| 5 | P05 Regression Tests & Quality Gates | Add integration tests, verify coverage, run full gate suite. | Green gates; `CHANGELOG.md` snippet. |
| 6 | P06 Visual Sign-off & Philosophy Check | Verify dashboard invariants with Playwright MCP. | Screenshots; branch ready for PR. |

Detailed tasks for each prompt are in `prompts/`.

---

## Bug Register Format

See `appendices/e-bug-register-format.md` for the full column definitions, severity rubric, root-cause categories, and finding-file naming convention.

---

## Post-Bug-Hunt Actions

After P06 visual sign-off is green, complete the steps in `appendices/c-post-bug-hunt-checklist.md`. The high-level summary is:

1. Final documentation sweep (`TRACKING.md`, `HANDOFF.md`, `DECISIONS.md`, `RISKS.md`, `KNOWN_ISSUES.md`).
2. Add a `CHANGELOG.md` release-notes snippet for fixed bugs.
3. Update root `TRACKING.md` and the v1 master backlog.
4. Rebase/merge the latest dashboard redesign commits if needed.
5. Open a PR from `feat/bug-hunt` to `dev` (or to `feat/dashboard-mailpit-redesign` if it has not merged yet).
6. Hand off to the human reviewer using `REVIEW_CHECKLIST.md`.

---

## References

- `docs/initiatives/openmuara-dashboard-mailpit-redesign/README.md`
- `docs/initiatives/openmuara-v1-master-backlog/TRACKING.md`
- `docs/initiatives/openmuara-v1-solid-gold/README.md`
- `docs/initiatives/openmuara-testing-gold-standard/README.md`
- `docs/initiatives/openmuara-ux-excellence/README.md`
- `docs/initiatives/openmuara-a11y-usability-polish/README.md`
- `PREREQUISITES.md`
- `KNOWN_ISSUES.md`
- `REVIEW_CHECKLIST.md`
- `appendices/a-security-checklist.md`
- `appendices/b-recommendations.md`
- `appendices/c-post-bug-hunt-checklist.md`
- `appendices/d-gold-standard-alignment.md`
- `appendices/e-bug-register-format.md`
- Root `AGENTS.md`
- `runbooks/debugging.md`
- `runbooks/testing.md`
- `runbooks/quality-gates.md`

---

## Self-assessment

**Solidity: 9.7 / 10**

The initiative now links to the project philosophy and user priority stack, protects dashboard redesign invariants, requires regression tests and visual sign-off, includes a detailed security checklist, documents a `findings/` workflow, provides prerequisites/assumptions, a communication/escalation plan, a post-bug-hunt PR checklist, a human review checklist, and a roadmap of follow-up recommendations. Files are split to stay close to the 250-line guideline.

The remaining 0.3 points are execution risk: the value depends on the next agent actually following the sign-off gates, keeping fixes minimal, and completing P06 visual sign-off. No additional documentation will fix that; it requires disciplined execution.
