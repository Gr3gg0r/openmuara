> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This document is subordinate to it.**

# Appendix D — Gold-Standard Alignment

> **Updated:** 2026-07-06

This bug hunt should stand on the shoulders of earlier quality initiatives and push the codebase toward OSS-grade solidity.

| Source initiative | What this bug hunt adopts |
|---|---|
| `openmuara-v1-solid-gold` | ≥80% package coverage target, trace-ID debuggability, additive-only changes, `task quality` discipline. |
| `openmuara-testing-gold-standard` | Race/shuffle runs, random-port integration tests, provider contract golden files, fuzz/property tests where applicable. |
| `openmuara-a11y-usability-polish` | WCAG AA contrast, keyboard shortcuts, axe-core zero serious violations, focus management. |
| `openmuara-ux-excellence` | Three-interface parity (CLI / API / dashboard), progressive disclosure, copy-paste examples, plain-language errors. |
| `openmuara-dashboard-mailpit-redesign` | The full set of dashboard invariants this bug hunt must protect. |

Every confirmed bug should be treated as a signal about a missing test, a missing invariant, or a missing guardrail. The fix should close the immediate hole **and** add a regression test so the hole stays closed.
