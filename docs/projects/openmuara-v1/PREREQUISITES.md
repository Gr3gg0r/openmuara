> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This project is subordinate to it.**

# OpenMuara v1 — Prerequisites & Pre-Flight Checklist

> **Purpose:** Human-readable checklist of everything needed BEFORE the AI starts executing.
> **Owner:** Human
> **Status:** 🟡 In progress

---

## 1. Skills Inventory

| Skill | Path | Required? | Installed? | Notes |
|-------|------|-----------|------------|-------|
| context-save | `~/.kimi-code/skills/context-save` | Recommended | ⬜ | Save working context before long sessions |
| context-restore | `~/.kimi-code/skills/context-restore` | Recommended | ⬜ | Resume saved context |
| investigate | `~/.kimi-code/skills/investigate` | Optional | ⬜ | Root-cause debugging |
| review | `~/.kimi-code/skills/review` | Optional | ⬜ | Pre-landing diff review |

---

## 2. MCP Inventory

| MCP | Purpose | Required? | Installed? | Verified? |
|-----|---------|-----------|------------|-----------|
| playwright | Browser automation, screenshots, QA | Optional | ⬜ | ⬜ |
| github | PR creation, issue tracking | Optional | ⬜ | ⬜ |

---

## 3. Environment & Access Checklist

| # | Item | Status | Notes |
|---|------|--------|-------|
| 1 | Repo cloned at `<repo-root>` and `git status` clean | ✅ | |
| 2 | On `dev` branch: `git branch --show-current` returns `dev` | ✅ | |
| 3 | Go 1.22+ installed: `go version` | ✅ | `/usr/local/go/bin/go` |
| 4 | `task` installed: `task --version` | ✅ | `<go-bin>/task` |
| 5 | `golangci-lint` installed: `golangci-lint --version` | ✅ | `<go-bin>/golangci-lint` |
| 6 | PATH includes `/usr/local/go/bin` and `<go-bin>` | ✅ | |
| 7 | `TMPDIR` set to executable dir (e.g., `<tmp-dir>`) | ✅ | `/tmp` is mounted `noexec` |
| 8 | `task check` passes on baseline | ✅ | 69.4% coverage |
| 9 | `task smoke` passes on baseline | ✅ | |

---

## 4. Human Approval Gates

| Gate | Trigger | Approver | Status |
|------|---------|----------|--------|
| A | Before renaming module/binary/config paths (rebrand) | Project Owner | ⬜ |
| B | Before adding/changing persistence schema (SQLite) | Backend Lead | ⬜ |
| C | Before changing webhook signature or provider plugin schema | Tech Lead | ⬜ |
| D | Before merging/pushing to `main` or `master` | Code Reviewer | ⬜ |

**How the AI uses gates:**
- When the AI reaches a gate, it stops execution and uses `AskUserQuestion` to request approval.
- The AI MUST NOT proceed past a gate without explicit human confirmation.
- If running in AFK mode, the AI skips the gate and logs it in `HANDOFF.md` for human review.

---

## 5. External Dependencies & Blockers

| Dependency | Needed By | Status | Owner | ETA |
|------------|-----------|--------|-------|-----|
| SenangPay API docs / signature spec | Prompt 10 | ⬜ | | 2026-07-05 |
| Stripe API reference for Checkout session shape | Prompt 09 | ⬜ | | Available online |
| | | ⬜ | | |

---

## 6. Scope Guardrails

**In Scope:**
- Rebrand `muara` → `OpenMuara` (module, binary, CLI, config, docs).
- SQLite persistence for transactions and webhook attempts.
- Universal payment API (`/v1/pay`, `/v1/pay/{ref}`, `/v1/refund/{ref}`).
- Stripe and SenangPay provider adapters.
- Scenario commands for deterministic testing.
- Webhook relay with multi-destination forwarding.
- Basic web UI (webhook inspector + provider status).
- Docker image and Docker Compose setup.
- OpenAPI spec and Go test SDK.

**Out of Scope (explicitly):**
- RevenueCat, App Store, Play Store adapters (v2 — hard frozen for v1).
- Multi-port runtime where each provider has its own port (v1.2+).
- MCP server (v1.2+).
- Production hosted service / SaaS.
- Real money processing.

**Budget / Time Guardrail:**
- Max prompts: 18
- Max estimated dev time: 40 hours
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
