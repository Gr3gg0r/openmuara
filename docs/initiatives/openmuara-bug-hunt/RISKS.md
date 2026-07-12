> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Bug Hunt — Risk Register

> **Updated:** 2026-07-06
>
> **Post-completion note:** This bug-hunt pass is complete. All risks below were either mitigated by the E1–E12 implementation or accepted as ongoing process risks and are now marked Closed.

| ID | Risk | Likelihood | Impact | Mitigation | Owner | Status |
|----|------|------------|--------|------------|-------|--------|
| R001 | Scope creep — fixing bugs turns into opportunistic refactors. | Medium | Medium | Enforce "minimal fix only" convention; require rationale for any change beyond the bug. | AI Agent | Closed |
| R002 | Test flakes mask real bugs or produce false positives. | Low | Medium | Run tests multiple times; isolate flaky tests and fix separately. | AI Agent | Closed |
| R003 | Fixes conflict with concurrent work on `dev` or other feature branches. | Medium | Low | Keep commits small and focused; rebase frequently; communicate via HANDOFF. | AI Agent | Closed |
| R004 | Deferred P2 bugs are forgotten. | Low | Low | All deferred items recorded in this register with explicit reason and target release. | AI Agent | Closed |
| R005 | Local-only bugs reproduce only in agent environment. | Low | Low | Document exact reproduction steps and environment details. | AI Agent | Closed |
| R006 | Fixing a bug reveals a deeper architectural issue. | Medium | Medium | Stop and document; do not expand scope without user sign-off and a plan update. | AI Agent | Closed |
| R007 | User sign-off for P0/P1 integration fixes delays the batch. | Medium | Medium | Identify sign-off needs during P02 triage; request sign-off early and batch dependent fixes. | AI Agent | Closed |
| R008 | Regression test coverage gaps hide reintroduced bugs. | Medium | High | Require a regression test for every fix; enforce the coverage gate. | AI Agent | Closed |
| R009 | A fix accidentally regresses a dashboard redesign invariant. | Medium | High | Maintain the invariant checklist; re-run P06 visual sign-off after any UI-impacting fix. | AI Agent | Closed |
| R010 | False positives or low-value findings consume the budget. | Low | Medium | Validate reproducibility in P02; downgrade/remove false positives with rationale. | AI Agent | Closed |
| R011 | Late regression discovered after P05 gates. | Low | High | Keep P05 focused on integration tests and full gate suite; reserve capacity for follow-up fixes. | AI Agent | Closed |

---

## Risk Template

```markdown
| ID | Risk | Likelihood | Impact | Mitigation | Owner | Status |
```

When adding a new risk:

1. Use the next sequential ID.
2. Link to the related bug ID when applicable.
3. Update status as the risk is mitigated, accepted, or closed.
