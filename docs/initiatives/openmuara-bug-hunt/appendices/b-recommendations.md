> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This document is subordinate to it.**

# Appendix B — Recommendations & Future Enhancements

> **Updated:** 2026-07-06

All recommendations below were approved and implemented during this bug-hunt pass. They are kept here as a record of what was added and as a starting point for the next quality cycle.

| # | Recommendation | Priority | Rationale | Target Initiative / Owner |
|---|----------------|----------|-----------|---------------------------|
| E1 | Store the P01 Playwright MCP screenshot baseline in `findings/visual-baseline/` and add a diff step to P06. | High | Catches unintended UI changes automatically; gives reviewers a visual changelog. | bug-hunt follow-up |
| E2 | Add GitHub issue/PR templates that match the bug register columns (severity, area, reproduction, root cause, regression test). | Medium | Makes external contributors and future agents report bugs consistently. | `.github/` |
| E3 | Integrate `govulncheck`, `npm audit --production`, and `golangci-lint` into CI as required checks. | High | Prevents dependency vulnerabilities and lint regressions from reaching `dev`. | CI / DevOps |
| E4 | Add fuzz/property tests for signature verification (HMAC/SHA256, MD5), idempotency keys, and the transaction state machine. | Medium | Finds edge cases that unit tests miss; aligns with testing-gold-standard. | testing-gold-standard v2 |
| E5 | Run mutation testing (`gremlins` or `go-mutesting`) on the packages most changed by fixes. | Low | Validates that regression tests actually kill mutants and are not just coverage padding. | bug-hunt follow-up |
| E6 | Schedule recurring bug-hunt sprints before each release, with a fixed time-box and a pre-defined reconnaissance checklist. | Medium | Prevents quality debt from accumulating; institutionalizes the process. | release process |
| E7 | Add provider contract conformance tests with golden files for every supported provider/version. | High | Catches provider drift early and documents expected request/response shapes. | testing-gold-standard |
| E8 | Add chaos tests for webhook dispatch: slow targets, connection resets, DNS failures, retry exhaustion. | Medium | Hardens the dispatcher and surfaces timeout/retry bugs. | webhook hardening |
| E9 | Maintain a public `KNOWN_ISSUES.md` at the repo root, synced from deferred bugs in `RISKS.md`. | Low | Sets honest expectations for users and contributors. | docs |
| E10 | Add a release-notes automation step that scrapes fixed bug IDs from `TRACKING.md`. | Low | Reduces manual changelog work and ensures no fix is forgotten. | release process |
| E11 | Add a coverage-regression comment bot on PRs that flags dropped module coverage. | Low | Keeps coverage discipline visible without blocking merges. | CI / DevOps |
| E12 | Establish an error-code taxonomy across providers so bugs can be classified by symptom (e.g., `E1001` signature mismatch). | Low | Makes debugging and bug reports more precise. | ux-excellence |
