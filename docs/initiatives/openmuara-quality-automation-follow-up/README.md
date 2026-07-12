> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Quality Automation Follow-Up

> **Status:** 🟡 Planned | **Started:** 2026-07-06
> **Scope:** Mature the E1–E12 bug-hunt recommendations from "implemented" to "automated and enforced," and close the remaining gaps in visual regression, mutation testing, recurring process, coverage gating, and error-code taxonomy.
> **Owner:** AI Agent (Kimi Code) | **Human Reviewer:** ___________
> **Target Repo:** `<repo-root>/`
> **Product Branch:** `feat/quality-automation-follow-up`
>
> **Why:** The bug hunt landed E1–E12 as tools and documentation. This initiative makes those tools self-sustaining: visual diffs run in CI, mutation tests guard the test suite, the bug-hunt process runs on a schedule, coverage regressions block merges, and the error-code taxonomy is adopted consistently across every provider.

---

## Initiative Structure

```
docs/initiatives/openmuara-quality-automation-follow-up/
├── README.md                    # This file
├── TRACKING.md                  # Central execution tracker
├── PREREQUISITES.md             # Tools, assumptions, branch base, time-box
├── RISKS.md                     # Risk register
├── DECISIONS.md                 # Decision log
├── REVIEW_CHECKLIST.md          # Pre-PR human review checklist
│
├── prompts/                     # Numbered, self-contained execution prompts
│   ├── _template.md
│   ├── 01-ci-visual-baseline.md
│   ├── 02-mutation-testing-gate.md
│   ├── 03-coverage-regression-gate.md
│   ├── 04-provider-errcode-adoption.md
│   ├── 05-recurring-process-and-known-issues-sync.md
│   └── 06-final-gates-and-documentation.md
│
└── appendices/                  # Deep-dive reference material
    ├── a-gaps-and-enhancements.md
    ├── b-best-practices.md
    ├── c-gold-standard-alignment.md
    └── d-recommendations-roadmap.md
```

Planning docs live in `docs/initiatives/openmuara-quality-automation-follow-up/` in the root repo.
Product code commits to the `feat/quality-automation-follow-up` branch, branched from `dev`.
Do not commit directly to `main` or `dev`.

> **Entry point:** Read `PREREQUISITES.md` before starting P01.

---

## Goals

1. **Visual baseline in CI** — run `npm run test:visual-baseline` on PRs that touch `web/dashboard/`, with deterministic baselines stored in `docs/initiatives/openmuara-bug-hunt/findings/visual-baseline/`.
2. **Mutation testing gate** — integrate `gremlins` into CI for the packages most changed by recent work; fail if the mutation score drops below a chosen threshold.
3. **Recurring bug-hunt automation** — add a scheduled GitHub issue/label workflow that opens a bug-hunt prep issue before each release.
4. **Coverage-regression gate** — evolve the PR comment bot into a required check that fails when any changed module drops coverage.
5. **Provider-wide error codes** — adopt `internal/errcode` in all provider packages and API error responses so every failure has a stable code.
6. **Known-issues automation** — add a script or CI check that warns when deferred items in `docs/initiatives/openmuara-bug-hunt/RISKS.md` are not reflected in root `KNOWN_ISSUES.md`.

---

## Philosophy & Priority Alignment

Root `AGENTS.md` and the dashboard redesign initiative establish this priority stack:

**UI > UX > Performance > Usability > Philosophy > Efficiency > Memory size**

For this follow-up, translate that stack into design choices:

- A CI gate that catches a **visual regression** is treated as high-value UX protection.
- A gate that produces **flaky failures** or slows the loop below usability is a bug to fix before it blocks anyone.
- **Performance and CI efficiency** matter, but not at the cost of dropping accessibility, visual diff, or coverage gates the user already approved.
- **Memory/bundle size** is a constraint; do not add heavy third-party services or telemetry to implement these gates.

The project philosophy is **local-first, simple, and explicit**:

- No external services, telemetry, or cloud dependencies.
- Every gate must be runnable locally with the same commands used in CI.
- Every config change is explicit; every secret stays server-side and out of git.
- Changes are minimal and additive. No speculative refactors.

---

## Conventions (MANDATORY)

### 1. Read `AGENTS.md` first
Root `AGENTS.md` governs branch rules, quality gates, autonomy boundaries, and code style.

### 2. Priority stack
When trade-offs arise, decide in this order:

1. **Correctness** — gates must not hide real failures.
2. **Reliability** — flaky gates are worse than no gates; fix or remove them.
3. **UX** — error messages, dashboard stability, and developer feedback come first.
4. **Performance** — only after the above are satisfied.
5. **Polish** — cosmetics last.

### 3. P0 integration changes need explicit user sign-off
Provider emulation logic, webhook signature verification, config persistence schemas, auth/billing/PII flows, and the provider plugin schema contract remain protected. The errcode adoption in P04 must not change provider behavior; it only wraps error values. Document any exception in `DECISIONS.md`.

### 4. One logical change per commit
Each gate or adoption step gets its own commit: `feat(scope): short description`.

### 5. Every gate must be testable locally
If a contributor cannot reproduce a CI failure with a single command, the gate is not ready.

### 6. No speculative refactors
Fix only what is required to automate or enforce the E1–E12 recommendations. Do not opportunistically refactor unrelated modules.

---

## Approach Options

| Option | Description | Pros | Cons | Recommended |
|---|---|---|---|---|
| A — Phased rollout | Add one gate per prompt; keep each optional until proven stable, then make required. | Low blast radius; easy to tune thresholds; builds trust. | Slightly longer timeline. | ✅ Yes |
| B — Big-bang | Add all gates at once and make them required immediately. | Fastest protection. | High flake risk; hard to debug which gate is noisy. | No |
| C — Commentary-only | Keep all gates as non-blocking comments/bots. | Zero CI disruption. | Does not enforce quality; regressions can still merge. | No |

Selected approach: **Option A — Phased rollout**.

---

## Metrics

| Metric | Current | Target | How measured |
|---|---|---|---|
| Visual baseline diff failures caught before merge | 0 | ≥1 per cycle where UI changes | CI logs / `npm run test:visual-baseline` |
| Mutation score (changed packages) | Not measured | ≥70% | `gremlins unleash` output |
| Coverage regressions blocked before merge | 0 | 100% of module drops caught | Coverage-comment workflow + gate |
| Provider packages using `errcode` | 1 (webhook) | All provider packages + API errors | `grep -R "errcode\." internal/*/...` |
| Recurring bug-hunt issues opened automatically | 0 | 1 per release | GitHub scheduled workflow |
| KNOWN_ISSUES sync drift | Possible | 0 | CI script diff |
| CI pipeline duration | ~8–10 min | <10 min | GitHub Actions wall time |



## Dependencies & Constraints

- Depends on the bug-hunt initiative being complete and merged to `dev`.
- Depends on GitHub Actions for CI automation.
- Depends on `gremlins` and Playwright being installable in CI within the existing 10-minute wall-time budget.
- All gates must remain runnable locally with the same commands used in CI.
- No SaaS or cloud-only tooling.
- Total CI wall time must stay under 10 minutes.

## Out of Scope

- New providers or features.
- Dashboard redesign changes.
- Large architectural refactors not driven by the items above.
- Replacing existing local-first tooling with SaaS or cloud-dependent services.
- Adding telemetry, analytics, or external dashboards.

---

## Success Criteria

- [ ] Visual baseline diff runs in CI and fails on unintended UI changes.
- [ ] Mutation testing runs in CI with a documented threshold.
- [ ] A recurring bug-hunt issue is created automatically before each release.
- [ ] Coverage regressions on changed modules block PR merges.
- [ ] All provider packages import and use `internal/errcode` for signature, config, and transaction errors.
- [ ] `KNOWN_ISSUES.md` sync check passes in CI.
- [ ] All quality gates still pass and total build time remains under 10 minutes.
- [ ] Every new gate is documented in `runbooks/quality-gates.md` or equivalent.

## Definition of Done

This initiative is done when:

1. P01–P06 are marked `✅` in `TRACKING.md` with commit hashes.
2. All success criteria above are checked and verified.
3. Every new required gate has been stable for at least three unrelated PRs, or any unstable gate is documented as non-blocking in `DECISIONS.md`.
4. `RISKS.md` shows all resolved risks as `Closed` or `Accepted`.
5. `HANDOFF.md` is updated with the final state and no open blockers.
6. `REVIEW_CHECKLIST.md` is completed by the human reviewer.
7. `CHANGELOG.md` has a release-notes snippet summarizing the new automation.
8. Root `TRACKING.md` is updated from `🟡 Planned` to `🟢 Completed`.

## Completion Verification

Run the commands in `TRACKING.md` → **Quality Gate Results** locally and confirm CI matches.
Then verify in GitHub Actions:

- Visual baseline job runs and passes on a dashboard-only change.
- Mutation testing job reports a score ≥ the documented threshold.
- Coverage-regression check passes or correctly fails on a coverage drop.
- Scheduled bug-hunt workflow creates an issue (can be triggered manually for verification).
- KNOWN_ISSUES sync check passes.

---

## Recommendations Roadmap

See `appendices/d-recommendations-roadmap.md` for the full list of suggestions, enhancements, good-to-haves, and gold-standard practices identified during the review. The roadmap covers CI/CD, testing, visual/UI, error handling, observability, process, documentation, security, and performance. Items are tagged `Core`, `Future`, or `Low` and are promoted into prompts only when selected.

## Gaps Covered

See `appendices/a-gaps-and-enhancements.md` for the full gap analysis.
High-level gaps this initiative closes:

- Visual baseline is currently a local script; no CI enforcement.
- Mutation testing is documented but not run or gated.
- Bug-hunt process is manual; no release-cadence reminder.
- Coverage bot only comments; it does not block regressions.
- `errcode` taxonomy exists but is only used in the webhook dispatcher.
- `KNOWN_ISSUES.md` can drift from the bug-hunt risk register.

---

## Gold-Standard Alignment

See `appendices/c-gold-standard-alignment.md`.

This initiative builds on:

- `openmuara-testing-gold-standard` — coverage, fuzz, contract, and mutation discipline.
- `openmuara-a11y-usability-polish` — visual regression protects a11y markup and layout.
- `openmuara-v1-solid-gold` — additive changes, local reproducibility, and OSS-grade gates.

---

## Post-PR Actions

After the final PR from `feat/quality-automation-follow-up` to `dev` is merged:

1. Update root `TRACKING.md` to mark this initiative `🟢 Completed`.
2. Add a `CHANGELOG.md` snippet summarizing the new gates and automation.
3. Update `runbooks/quality-gates.md` if not already done in P06.
4. Delete the feature branch locally after merge.
5. Announce the new gates in any contributor channel so the team knows how to update snapshots, override coverage gates, and interpret mutation scores.

## References

- `docs/initiatives/openmuara-bug-hunt/README.md`
- `docs/initiatives/openmuara-bug-hunt/appendices/b-recommendations.md`
- `docs/initiatives/openmuara-quality-automation-follow-up/appendices/a-gaps-and-enhancements.md`
- `docs/initiatives/openmuara-quality-automation-follow-up/appendices/b-best-practices.md`
- `docs/initiatives/openmuara-quality-automation-follow-up/appendices/c-gold-standard-alignment.md`
- `docs/initiatives/openmuara-quality-automation-follow-up/appendices/d-recommendations-roadmap.md`
- `docs/initiatives/openmuara-testing-gold-standard/README.md`
- `docs/initiatives/openmuara-a11y-usability-polish/README.md`
- `docs/initiatives/openmuara-v1-solid-gold/README.md`
- Root `AGENTS.md`
- `runbooks/quality-gates.md`
