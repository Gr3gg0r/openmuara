> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This document is subordinate to it.**

# OpenMuara Quality Automation Follow-Up — Risk Register

> **Updated:** 2026-07-06

| ID | Risk | Likelihood | Impact | Mitigation | Owner | Status |
|----|------|------------|--------|------------|-------|--------|
| R001 | Visual baseline is flaky (font rendering, animations, timing). | Medium | High | Start non-blocking; hide dynamic elements; use deterministic DB; retry with tolerance. | AI Agent | Accepted |
| R002 | Mutation testing is slow and breaks the 10-minute CI budget. | Medium | Medium | Target only changed packages; cache tool install; run on schedule instead of every PR if needed. | AI Agent | Accepted |
| R003 | Coverage gate blocks legitimate refactors that temporarily drop coverage. | Medium | Medium | Compare only changed modules; allow override with `DECISIONS.md` rationale; start as commentary. | AI Agent | Accepted |
| R004 | Provider errcode adoption accidentally changes API error messages. | Medium | High | Add codes alongside existing messages; never remove or alter public message text without sign-off. | AI Agent | Accepted |
| R005 | Scheduled bug-hunt workflow creates noise or duplicates. | Low | Low | Use deterministic issue title; skip if an open bug-hunt issue already exists. | AI Agent | Accepted |
| R006 | KNOWN_ISSUES sync script produces false positives. | Low | Medium | Allow an explicit "intentionally empty" marker; run as warning before promotion. | AI Agent | Accepted |
| R007 | New gates slow down the local developer loop. | Medium | Medium | Keep every gate runnable via `task quality`; document fast vs. full modes. | AI Agent | Accepted |
| R008 | Tooling dependencies (gremlins, Playwright) drift or break. | Low | Medium | Pin versions in CI; document install commands in `PREREQUISITES.md`. | AI Agent | Accepted |
| R009 | Provider errcode wrapping changes `errors.Is` / equality behavior in existing tests. | Medium | High | Wrap with `fmt.Errorf("...: %w", errcode.New(...))` so existing `errors.Is` checks still work; add tests for both old and new behavior. | AI Agent | Accepted |
| R010 | Visual baseline PNGs increase repository size. | Low | Low | Keep baselines minimal; delete obsolete baselines; consider Git LFS if size grows >1 MB. | AI Agent | Accepted |
| R011 | Coverage gate false-positives on small packages (e.g., one-line change drops coverage 100% → 0%). | Medium | Medium | Set a minimum changed-line threshold (e.g., ignore packages with <10 changed lines); document override path. | AI Agent | Accepted |
| R012 | Required gate is promoted too early and blocks urgent fixes. | Low | High | Demote to commentary immediately; require three stable PRs before re-promotion; document in `DECISIONS.md`. | AI Agent | Accepted |

---

## Risk Template

```markdown
| ID | Risk | Likelihood | Impact | Mitigation | Owner | Status |
```

When adding a new risk:

1. Use the next sequential ID.
2. Link to the related prompt ID when applicable.
3. Update status as the risk is mitigated, accepted, or closed.
