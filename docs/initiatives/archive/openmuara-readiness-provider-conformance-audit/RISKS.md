> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Provider Conformance Audit Risk Register

> **Created:** 2026-07-08
> **Last Updated:** 2026-07-09
> **Status:** ✅ Complete — risks accepted or mitigated; register updated with residual risk owners.

---

| ID | Risk | Likelihood | Impact | Mitigation | Owner |
|---|---|---|---|---|---|
| R01 | Real provider docs are ambiguous, versioned, or change without notice | High | Medium | Pin emulated version in `gateway.yml` and provider docs; document known limitations; schedule periodic re-review | AI Agent |
| R02 | Conformance tests become a maintenance burden and are skipped | Medium | Medium | Use table-driven tests, shared fixtures, golden files, and clear `-update` workflows; keep tests close to provider package | AI Agent |
| R03 | Undocumented deviation misleads users into shipping broken integrations | Medium | High | Require every deviation to be listed in `KNOWN_ISSUES.md` and provider docs; treat undocumented deviations as bugs | AI Agent / Maintainer |
| R04 | External provider team does not respond to review request | Medium | Medium | Send request early; document the attempt; use public docs and community reports as secondary validation | Maintainer |
| R05 | Emulating exact provider behavior requires secrets or sandbox accounts unavailable to the project | Medium | Medium | Rely on published docs and community fixtures; clearly label behavior validated by docs vs. behavior validated by sandbox | AI Agent |
| R06 | Provider-specific test fixtures drift from real behavior over time | Medium | Medium | Version fixtures with provider API version; add `docs/providers/<provider>.md` changelog; re-run audit yearly | Maintainer |
| R07 | Conformance tests slow CI below acceptable thresholds | Medium | Low | Run provider contract tests in parallel; keep golden-file tests fast; isolate slow scenario tests | AI Agent |
| R08 | Simple-runtime providers hide conformance gaps behind YAML abstraction | Medium | High | Add explicit contract tests for each `gateway.yml` route; validate YAML against schema | AI Agent |
| R09 | Webhook dispatch differences are hard to observe without real receiver | Medium | Medium | Provide reference webhook receiver fixtures; test payload + signature against known examples | AI Agent |
| R10 | Over-emulating provider quirks delays release | Low | Medium | Use maturity model to ship L1–L4 first; defer L5/L6 to follow-up sprints | Maintainer |

## Risk treatment summary

- **Accept:** R05, R10 (documented and bounded).
- **Mitigate:** R01, R02, R03, R04, R06, R07, R08, R09 (controls implemented via tests, golden files, CI gate, and `KNOWN_ISSUES.md`).
- **Transfer:** None.
- **Avoid:** None.

## Residual risks

| ID | Residual risk | Owner | Monitoring |
|---|---|---|---|
| R04 | Fawry team may not respond to review request | Maintainer | Track in `KNOWN_ISSUES.md`; send follow-up after 30 days |
| R06 | Provider fixtures drift over time | Maintainer | Yearly audit; watch provider changelog |
| R10 | Over-emulating quirks delays future releases | Maintainer | Use maturity model to bound scope per release |
