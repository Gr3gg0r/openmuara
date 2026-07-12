> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara UX Excellence — Risk Register

> **Created:** 2026-07-01
> **Status:** ⬜ Draft

---

## Active Risks

| ID | Risk | Likelihood | Impact | Mitigation | Owner |
|----|------|------------|--------|------------|-------|
| R001 | Interactive `muara init` breaks CI or headless usage | Medium | High | Add `--defaults` / `--non-interactive` flag and keep existing non-interactive behavior as fallback. | TBD |
| R002 | Dashboard onboarding state adds persistence complexity | Low | Medium | Derive checklist state from existing data (transactions, webhooks, providers) instead of a new state file. | TBD |
| R003 | Config validation becomes too strict and rejects legacy configs | Medium | High | Validation warnings first, errors later; support legacy top-level `fawry`/`stripe` keys during transition. | TBD |
| R004 | Webhook debugger exposes sensitive payloads on `/_admin` | Low | High | Keep debugger local-only; respect existing CORS/CSRF settings; do not log secrets. | TBD |
| R005 | Scope creep turns UX initiative into a full UI rewrite | Medium | Medium | Stick to the 9 prompts in `TRACKING.md`; defer non-listed ideas to future initiatives. | TBD |
| R006 | Config line numbers require raw YAML parsing | Medium | Medium | Implement best-effort line numbers with a YAML parser fallback; accept field path + file path when line numbers are unavailable. | TBD |
| R010 | Dashboard polling creates noise or load | Low | Low | Pause polling on hidden tabs; cap default ledger size to 50 events; make refresh interval configurable. | TBD |
| R011 | New admin endpoints break external scripts | Low | Medium | Treat `/_admin` as additive-only (D007); path-prefix with `/_admin/v1/` if a breaking change is ever required. | TBD |
| R012 | Zero-data state is unhelpful | Low | Medium | Design a clear empty ledger state with a copy-paste first-charge example per active provider. | TBD |

---

## Resolved Risks

| ID | Risk | Resolution |
|----|------|------------|
| R007 | What should the dashboard primary view be called? | Decided to use "ledger" (D005) because it matches the payments domain and the existing transaction ledger. |
| R008 | Should the wizard be interactive by default? | Decided yes in TTY, with `--defaults` flag for headless usage (D002). |
| R009 | Should onboarding progress be persisted? | Decided no; derive state from existing data (D003). |
