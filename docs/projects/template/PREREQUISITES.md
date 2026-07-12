> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This template is subordinate to it.**

# <Project Name> — Prerequisites & Pre-Flight Checklist

> **Purpose:** Human-readable checklist of everything needed BEFORE the AI starts executing. Review this before copying the template and initiating work.
> **Owner:** Human
> **Status:** ⬜ Not reviewed / 🟡 In progress / ✅ Ready

---

## 1. Skills Inventory

List the Kimi skills / gstack skills expected to be used. Install any that are missing before starting.

| Skill | Path | Required? | Installed? | Notes |
|-------|------|-----------|------------|-------|
| context-save | `~/.kimi-code/skills/context-save` | Recommended | ⬜ | Save working context before long sessions |
| context-restore | `~/.kimi-code/skills/context-restore` | Recommended | ⬜ | Resume saved context |
| investigate | `~/.kimi-code/skills/investigate` | Optional | ⬜ | Root-cause debugging |
| review | `~/.kimi-code/skills/review` | Optional | ⬜ | Pre-landing diff review |
| | | | | |

**Install a skill:**
```bash
/skill install <skill-name>
```

---

## 2. MCP Inventory

List the MCP servers the AI will need. Install and verify any that are missing.

| MCP | Purpose | Required? | Installed? | Verified? |
|-----|---------|-----------|------------|-----------|
| playwright | Browser automation, screenshots, QA | Optional | ⬜ | ⬜ |
| github | PR creation, issue tracking, code search | Optional | ⬜ | ⬜ |
| | | | | |

**Verify an MCP:**
```bash
/mcp list
```

---

## 3. Environment & Access Checklist

Check each item before the AI begins. A missing item here will block execution.

| # | Item | Status | Notes |
|---|------|--------|-------|
| 1 | Repo cloned at `<repo-root>` and `git status` clean | ⬜ | |
| 2 | On `dev` branch: `git branch --show-current` returns `dev` | ⬜ | |
| 3 | Go 1.22+ installed: `go version` | ⬜ | |
| 4 | `task` installed: `task --version` | ⬜ | |
| 5 | `golangci-lint` installed: `golangci-lint --version` | ⬜ | |
| 6 | `task check` passes on baseline | ⬜ | |
| 7 | `task smoke` passes on baseline | ⬜ | |
| 8 | Required env vars set (`.env` or `.env.example` at repo root) | ⬜ | Copy from `.env.example` if exists |

---

## 4. Human Approval Gates

These are points where the AI MUST pause and ask for human approval before proceeding. Define them upfront so the AI knows when to stop.

| Gate | Trigger | Approver | Status |
|------|---------|----------|--------|
| A | Before modifying provider plugin schema or interface | Tech Lead | ⬜ |
| B | Before adding/changing persistence schema (SQLite/migrations) | Backend Lead | ⬜ |
| C | Before changing billing, pricing, or webhook signature logic | Product + Engineering | ⬜ |
| D | Before renaming module/binary/config paths (rebrand) | Project Owner | ⬜ |
| E | Before merging/pushing to `main` or `master` | Code Reviewer | ⬜ |
| F | | | ⬜ |

**How the AI uses gates:**
- When the AI reaches a gate, it stops execution and uses `AskUserQuestion` to request approval.
- The AI MUST NOT proceed past a gate without explicit human confirmation.
- If running in AFK mode, the AI skips the gate and logs it in `HANDOFF.md` for human review.

---

## 5. External Dependencies & Blockers

| Dependency | Needed By | Status | Owner | ETA |
|------------|-----------|--------|-------|-----|
| Example: Vendor API documentation for SenangPay | Prompt 04 | ⬜ | | YYYY-MM-DD |
| | | ⬜ | | |

---

## 6. Scope Guardrails

Define what is IN scope and what is OUT of scope so the AI does not wander.

**In Scope:**
- 
- 

**Out of Scope (explicitly):**
- 
- 

**Budget / Time Guardrail:**
- Max prompts: ___
- Max estimated dev time: ___ hours
- If exceeded, human review required before continuing.

---

## 7. Sign-Off

This project is cleared for AI execution when:
- [ ] All **Required** skills are installed.
- [ ] All **Required** MCPs are installed and verified.
- [ ] Environment checklist is 100% ✅.
- [ ] Human approval gates are defined and approvers notified.
- [ ] External dependencies are resolved or have committed ETAs.
- [ ] Scope guardrails are agreed upon.

**Signed off by:** ___________ **Date:** ___________
