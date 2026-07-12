> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This document is subordinate to it.**

# Appendix C — Gold-Standard Alignment

> **Updated:** 2026-07-06

This follow-up initiative builds on earlier quality initiatives and pushes the codebase toward OSS-grade solidity.

| Source initiative | What this follow-up adopts |
|---|---|
| `openmuara-v1-solid-gold` | ≥80% package coverage target, additive-only changes, local reproducibility, `task quality` discipline. |
| `openmuara-testing-gold-standard` | Race/shuffle runs, provider contract golden files, fuzz/property tests, mutation testing. |
| `openmuara-a11y-usability-polish` | Visual regression guard to protect layout, contrast, and keyboard navigation. |
| `openmuara-ux-excellence` | Stable error codes and plain-language error messages across CLI/API/dashboard. |
| `openmuara-dashboard-mailpit-redesign` | Dashboard invariants that visual baseline and errcode adoption must not regress. |
| `openmuara-bug-hunt` | E1–E12 recommendations as the baseline to automate and enforce. |

Every new gate should be treated as a signal amplifier: it makes existing quality holes visible without introducing new behavior. If a gate surfaces a problem, the fix should close the hole and add a regression test so it stays closed.
