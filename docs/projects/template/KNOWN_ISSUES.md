> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This template is subordinate to it.**

# <Project Name> — Known Issues & Out-of-Scope List

> **Purpose:** Prevent the AI from wasting time on pre-existing bugs or out-of-scope problems. If the AI encounters these, it logs them here and moves on — it does NOT fix them unless explicitly instructed.
> **Rule:** Update this file whenever the AI discovers a pre-existing issue unrelated to the current project.

---

## Pre-Existing Bugs (Do NOT Fix)

| ID | Issue | Location | Impact | Why Out of Scope |
|----|-------|----------|--------|------------------|
| K01 | | | | |
| K02 | | | | |

---

## Out-of-Scope Areas

| Area | Reason | Boundary |
|------|--------|----------|
| | | |
| | | |

---

## Pre-Existing Test Failures

If the AI runs tests and sees failures BEFORE making changes, it MUST log them here and NOT attempt to fix them unless they are directly caused by the current project's changes.

| Test Suite | Failing Test | Error | Logged Date |
|------------|-------------|-------|-------------|
| | | | |

---

## How to Use This File

1. **Before starting a step:** Scan this file to know what landmines to avoid.
2. **During execution:** If you encounter a pre-existing bug unrelated to your task, STOP trying to fix it. Log it here with a `K##` ID.
3. **In the prompt:** Reference this file if the step touches code near a known issue: "Do NOT fix K01. Only modify ..."
