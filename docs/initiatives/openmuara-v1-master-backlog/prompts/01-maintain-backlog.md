> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

## 01 — Maintain the Master Backlog

### Context

The master backlog must stay current as product code changes. After any product-code commit, the matching row in this backlog (and its source tracker) needs its status, commit hash, and notes updated.

### Current State

- **Repo:** `<repo-root>`
- **Branch:** `dev`
- **Target Files:** `docs/initiatives/openmuara-v1-master-backlog/TRACKING.md`, source tracker

### Scope

- **In scope:**
  - Update backlog rows after a product-code change.
  - Add new known issues or risks discovered during work.
  - Keep cross-references accurate.
- **Out of scope:**
  - Product code changes.
  - Reprioritizing frozen v2 items without human approval.

### Pre-flight

```bash
cd <repo-root>
git status
git branch --show-current  # must be dev
```

### Execution

1. Identify which backlog item changed.
2. Update its status and commit hash in `TRACKING.md`.
3. Update the source tracker:
   - Root `TRACKING.md` for P## / T## items.
   - `docs/projects/openmuara-v1/TRACKING.md` for phase items.
   - `docs/initiatives/openmuara-v1-solid/TRACKING.md` for S## items.
4. If a new risk or issue was found, update `RISKS.md` or `KNOWN_ISSUES.md`.
5. Update `HANDOFF.md` in this initiative.

### Quality Gates

- Backlog files contain no hard-coded absolute filesystem paths.
- Cross-references to source trackers are still valid.

### Commit

```bash
git add docs/initiatives/openmuara-v1-master-backlog/
git commit -m "docs(backlog): update v1 master backlog"
```

### Post-completion

1. Verify `README.md` status line is current.
2. Hand off via `HANDOFF.md` if ending the session.
