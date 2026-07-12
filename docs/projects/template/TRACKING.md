> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This template is subordinate to it.**

# <Project Name> — Execution Tracker

> **Updated:** YYYY-MM-DD | **Total Prompts:** N | **Status:** ⬜ Init / 🟡 In Progress / ✅ Complete
>
> **AI Agent:** When prompts are ready for execution, update the Status above from ⬜ Init to 🟡 In Progress. Do not leave it as ⬜ Init after plan creation.
> **Scope:** One-line description of what this project covers.
> **Repo:** `<repo-root>`
> **Product Branch:** `dev`
> **Last Agent Action:** ________________________________
> **Next Agent Action:** ________________________________

---

## ⚠️ Repo & Branch Boundary Rules (READ BEFORE EVERY PROMPT)

- `docs/projects/<project-name>/` is for **planning docs only**.
- All implementation code commits to the **`dev`** branch in `<repo-root>`.
- **Never commit directly to `main` or `master`.**
- **Do not mix planning-doc commits with product-code commits.**
- Pre-flight check: `git branch --show-current` must return `dev` before any product-code commit.

---

## Legend

| Icon | Meaning |
|------|---------|
| ⬜ | To Do |
| 🟡 | In Progress |
| ✅ | Completed |
| ❌ | Blocked |
| 🔀 | Parallel Safe (can run simultaneously with others) |

---

## Execution Rules

1. Execute steps in order unless marked **[PARALLEL SAFE]**.
2. Every step MUST end with: tests passing → git commit → update this file to `✅`.
3. If a step fails a quality gate, STOP. Do not proceed. Log the blocker in `RISKS.md`.
4. After EVERY step, update `HANDOFF.md` with what was done and what comes next.
5. Product-code commits happen on `dev`.

---

## Prompt / Task Inventory

| Step | Title | Target Files | Depends On | Status | Commit | Decisions | Notes |
|------|-------|--------------|------------|--------|--------|-----------|-------|
| 01 | | | — | ⬜ | — | — | |
| 02 | | | 01 | ⬜ | — | — | |
| 03 | | | 02 | ⬜ | — | — | |

---

## Quality Gate Results

### Stack: Go (`<repo-root>`)

| Gate | Command | Target | Status |
|------|---------|--------|--------|
| Format | `task fmt` | Clean | ⬜ |
| Vet | `task vet` | Clean | ⬜ |
| Lint | `task lint` | Zero issues | ⬜ |
| Test | `task test` | All pass | ⬜ |
| Race | `task test-race` (if available) | All pass | ⬜ |
| Coverage | `task coverage` | ≥ 50% | ⬜ |
| Build | `task build` | Compiles | ⬜ |
| Smoke | `task smoke` | Passes | ⬜ |

---

## Findings Inventory

| Step | Finding File | Status | Summary |
|------|-------------|--------|---------|
| | | ⬜ | |

---

## Runbooks Inventory

| Runbook | File | Status | Purpose |
|---------|------|--------|---------|
| | | ⬜ | |

---

## Merge / Completion Checklist

Before marking this project complete:

- [ ] All prompts in this tracker are ✅.
- [ ] All quality gates show ✅ Pass.
- [ ] `DECISIONS.md` is up to date.
- [ ] `RISKS.md` shows all risks as mitigated or accepted.
- [ ] `KNOWN_ISSUES.md` is accurate.
- [ ] `HANDOFF.md` reflects final state.
- [ ] Latest `dev` rebased/merged into the working tree and gates verified one last time.
- [ ] `dev` pushed to origin.

---

## Context Links

| Resource | Path |
|----------|------|
| This Tracker | `docs/projects/<project-name>/TRACKING.md` |
| Handoff | `docs/projects/<project-name>/HANDOFF.md` |
| Decisions | `docs/projects/<project-name>/DECISIONS.md` |
| Risks | `docs/projects/<project-name>/RISKS.md` |
| Prompts | `docs/projects/<project-name>/prompts/` |
| Tasks | `docs/projects/<project-name>/tasks/` |
| Findings | `docs/projects/<project-name>/findings/` |
| Runbooks | `docs/projects/<project-name>/runbooks/` |
| QA | `docs/projects/<project-name>/qa/` |
