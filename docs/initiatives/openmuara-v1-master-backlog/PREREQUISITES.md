> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara v1 Master Backlog — Prerequisites

> **Purpose:** Lightweight pre-flight checklist for anyone updating or reading this backlog.
> **Owner:** Human or AI Agent
> **Status:** ⬜ Not reviewed / 🟡 In progress / ✅ Ready

---

## 1. Required Context

Before editing this backlog, you must have read:

- [ ] `<repo-root>/AGENTS.md`
- [ ] `<repo-root>/TRACKING.md`
- [ ] `<repo-root>/docs/projects/openmuara-v1/TRACKING.md`
- [ ] `<repo-root>/docs/initiatives/openmuara-v1-solid/TRACKING.md`

---

## 2. Environment Checklist

| # | Item | Status | Notes |
|---|------|--------|-------|
| 1 | Repo cloned and `git status` clean | ⬜ | |
| 2 | Current branch is `dev` | ⬜ | Never edit backlog on `main`. |
| 3 | Go toolchain installed (1.22+) | ⬜ | Only needed if backlog drives code work. |
| 4 | `golangci-lint` installed | ⬜ | Only needed if backlog drives code work. |

---

## 3. Scope Guardrails

**In Scope for this backlog:**
- Tracking priority, status, and ownership of v1 work.
- Cross-referencing other trackers.
- Logging new known issues and risks.

**Out of Scope for this backlog:**
- Direct product code changes.
- Modifying v2-frozen work without explicit human approval.
- Committing screenshots, QA artifacts, or temp files.

---

## 4. How to Update This Backlog

1. Run `git status`.
2. Edit `TRACKING.md` in this folder.
3. If a new risk or known issue is discovered, update `RISKS.md` or `KNOWN_ISSUES.md`.
4. Commit as a docs-only change on `dev`:
   ```bash
   git add docs/initiatives/openmuara-v1-master-backlog/
   git commit -m "docs(backlog): update v1 master backlog"
   ```
