> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Repository Hygiene Risk Register

> **Created:** 2026-07-08  
> **Last Updated:** 2026-07-10  
> **Status:** Draft

---

| ID | Risk | Likelihood | Impact | Mitigation | Owner |
|---|---|---|---|---|---|
| R01 | Secret in git history becomes public after transfer | Low | Critical | Run `gitleaks` full-history scan before transfer; if a leak is found, rewrite history **before** transfer or rotate the secret. | AI Agent |
| R02 | Legacy `toyol` references confuse contributors or search engines | Medium | Medium | Systematically replace tracked references; keep directory name only. | AI Agent |
| R03 | Duplicate `SECURITY.md` files cause GitHub to surface outdated policy | Medium | Low | Consolidate under `.github/SECURITY.md`; redirect root file. | AI Agent |
| R04 | Stale `main` branch gives visitors an outdated impression | Medium | Medium | Fast-forward `main` to `dev` before public transfer; document release branch policy. | AI Agent |
| R05 | Leftover worktrees or branches cause merge conflicts or stale code | Medium | Low | Remove merged worktrees and delete/archive stale branches after initiative delivery. | AI Agent |
| R06 | `.gitignore` / `.dockerignore` gaps commit local or generated files | Medium | Medium | Audit ignore rules against current workspace; test on a clean clone. | AI Agent |
| R07 | Missing governance files (SUPPORT, MAINTAINERS) slow contributor onboarding | Medium | Low | Add governance files before public transfer. | AI Agent |
| R08 | Pre-commit hooks become too slow and are bypassed | Low | Medium | Keep hooks fast; run heavy checks (full test suite) in CI, not pre-commit. | AI Agent |
| R09 | Repository settings not replicated after org transfer | Low | Medium | Document required settings in `APPENDIX.md`; verify post-transfer. | AI Agent |
| R10 | AI-generated content not disclosed, reducing trust | Low | Low | Add a concise AI-assisted development note in README or CONTRIBUTING. | AI Agent |

---

## Risk acceptance

- **R01** (secret in history) is the only risk that could block public transfer. Current scan is clean; a final scan must be run immediately before transfer.
- **R02–R10** are polish and process risks. They are mitigated by the execution plan and do not block transfer if incomplete, but they materially affect first impressions.
