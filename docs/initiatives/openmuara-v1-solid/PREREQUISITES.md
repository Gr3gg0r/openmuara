> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara v1 Solid — Prerequisites & Pre-Flight Checklist

> **Purpose:** Human-readable checklist of everything needed BEFORE the AI starts executing.
> **Owner:** Human
> **Status:** ⬜ Not reviewed / 🟡 In progress / ✅ Ready

---

## 1. Skills Inventory

| Skill | Path | Required? | Installed? | Notes |
|-------|------|-----------|------------|-------|
| browse | `~/.kimi-code/skills/browse` | Recommended | ⬜ | For dashboard QA, screenshots |
| qa | `~/.kimi-code/skills/qa` | Recommended | ⬜ | For test-fix-verify loops |
| review | `~/.kimi-code/skills/review` | Optional | ⬜ | For pre-landing diff review |
| investigate | `~/.kimi-code/skills/investigate` | Optional | ⬜ | For root-cause debugging |
| context-save | `~/.kimi-code/skills/context-save` | Recommended | ⬜ | Save context before long sessions |
| context-restore | `~/.kimi-code/skills/context-restore` | Recommended | ⬜ | Resume saved context |

**Install a skill:**
```bash
/skill install <skill-name>
```

---

## 2. MCP Inventory

| MCP | Purpose | Required? | Installed? | Verified? |
|-----|---------|-----------|------------|-----------|
| playwright | Browser automation, dashboard QA | Recommended | ⬜ | ⬜ |
| github | PR creation, issue tracking | Optional | ⬜ | ⬜ |

**Verify an MCP:**
```bash
/mcp list
```

---

## 3. Environment & Access Checklist

| # | Item | Status | Notes |
|---|------|--------|-------|
| 1 | Repo cloned at `<repo-root>/` and `git status` clean | ⬜ | |
| 2 | Go toolchain installed (1.22+) | ⬜ | |
| 3 | `golangci-lint` installed | ⬜ | |
| 4 | `go build ./...` passes | ⬜ | |
| 5 | `go test ./...` passes | ⬜ | |
| 6 | `./scripts/smoke-test.sh` passes | ⬜ | |
| 7 | `task` CLI available (optional) | ⬜ | `task check`, `task smoke` |

---

## 4. Human Approval Gates

| Gate | Trigger | Approver | Status |
|------|---------|----------|--------|
| A | Before modifying provider webhook signature verification | Tech Lead | ⬜ |
| B | Before changing transaction state machine rules | Tech Lead | ⬜ |
| C | Before changing OpenAPI response schemas | Tech Lead | ⬜ |
| D | Before merging to `main` | Code Reviewer | ⬜ |

**How the AI uses gates:**
- The AI stops execution and uses `AskUserQuestion` to request approval.
- The AI MUST NOT proceed past a gate without explicit human confirmation.
- In AFK mode, the AI logs the gate in `HANDOFF.md` and makes the most conservative choice.

---

## 5. External Dependencies & Blockers

| Dependency | Needed By | Status | Owner | ETA |
|------------|-----------|--------|-------|-----|
| None | — | ✅ | — | — |

---

## 6. Scope Guardrails

**In Scope:**
- Fix admin dashboard regressions from pagination.
- Sync OpenAPI spec with current API.
- Apply state machine consistently across all status transitions.
- Make Fawry escape action update the ledger.
- Verify incoming Fawry webhook signatures.
- Improve per-provider dispatcher wiring.
- Update runbooks and README.

**Out of Scope (explicitly):**
- App Store / Play Store / RevenueCat (hard frozen for v2).
- New payment providers.
- Production hosting / SaaS.
- Real money processing.

**Budget / Time Guardrail:**
- Max prompts: 6
- Max estimated dev time: 8 hours
- If exceeded, human review required before continuing.

---

## 7. Sign-Off

This initiative is cleared for AI execution when:
- [ ] All **Required** skills are installed.
- [ ] All **Required** MCPs are installed and verified.
- [ ] Environment checklist is 100% ✅.
- [ ] Human approval gates are defined and approvers notified.
- [ ] Scope guardrails are agreed upon.

**Signed off by:** ___________ **Date:** ___________
