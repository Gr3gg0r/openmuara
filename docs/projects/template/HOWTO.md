> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This template is subordinate to it.**

# How to Decompose a Task into the OpenMuara Project Template

> **Purpose:** Step-by-step guide for AI agents to take a raw human task and produce a complete, executable project using the template.
> **Audience:** AI Agent (Kimi Code)
> **Input:** A raw task description from the user.
> **Output:** A fully populated `docs/projects/<project-name>/` folder ready for execution.

---

## Overview

When the user says something like:

> "Add SQLite persistence to the transaction ledger"

You do NOT start coding immediately. You follow this decomposition workflow to produce a robust, trackable, executable plan.

```
Raw Task → Decompose → Populate Template → Human Review → Execute on dev
```

---

## Step 1: Understand the Task

Read the user's request. Identify:

| Question | Why It Matters |
|----------|----------------|
| What is the problem? | Defines the "What This Solves" sections. |
| Which files/packages are affected? | Determines target paths. |
| Is this multi-phase? (>1 prompt needed) | Decides if the template is needed at all. |
| Does it touch P0 flows? | Triggers human approval gates in PREREQUISITES.md. |
| Does it touch auth/billing/PII/DB? | Triggers autonomy boundary checks. |
| Are there existing issues nearby? | Populates KNOWN_ISSUES.md scan list. |

**If the task is single-session or <10 lines, do NOT use the template.** Just do the work on `dev`.

---

## Step 2: Scaffold the Project Folder

Copy this template:

```bash
cp -R docs/projects/template docs/projects/<project-kebab-case-name>
```

Replace `<project-kebab-case-name>` with a descriptive name: `openmuara-rebrand`, `sqlite-persistence`, `stripe-provider-adapter`, etc.

---

## Step 3: Populate the Shell Files

Run through these in order. Fill every placeholder. Do not skip fields.

### 3.1 `README.md`
- `<PROJECT_NAME>` → Real project name
- `Started`, `Target End`, `Human Reviewer` → Fill or leave blank for human
- `Scope` → One-line summary
- **Delete the `## AI Agent Quickstart` section** once initialized.

### 3.2 `PREREQUISITES.md`
- Fill Skills Inventory (mark which are installed)
- Fill MCP Inventory (mark which are installed/verified)
- Fill Environment & Access Checklist
- Define Human Approval Gates (if auth/billing/PII/DB touched, create gates)
- Fill Scope Guardrails
- Leave Sign-Off for human

### 3.3 `TRACKING.md`
- Replace `<Project Name>` with the real name.
- Fill Scope line.
- In Prompt Inventory, create one row per logical step you anticipate.
- Do NOT guess every detail — leave Title blank if unsure, but create the rows.
- Fill Quality Gate Results with the correct commands (`task ...`).

### 3.4 `RISKS.md`
- Create at least one risk entry for every project with >3 steps.
- Generic risks to consider: breaking API contract, provider emulation drift, DB migration failure, coverage drop.

### 3.5 `KNOWN_ISSUES.md`
- Scan the affected packages for known bugs near the files you will touch.
- Run tests to find pre-existing failures.
- Log anything you find.

### 3.6 `REFERENCES.md`
- Add architecture docs relevant to the task.
- Add API specs, vendor docs (Stripe, Fawry, SenangPay, etc.).

### 3.7 `DECISIONS.md` & `HANDOFF.md`
- Replace `<Project Name>` placeholders.
- Leave content empty for now — they fill during execution.

---

## Step 4: Decompose into Prompts

This is the critical step. Break the raw task into numbered, self-contained prompts.

### 4.1 Rules for Decomposition

| Rule | Example |
|------|---------|
| **One prompt = one logical step** | "Create the SQLite store" is one prompt. "Create store AND migrate all handlers" is too big. |
| **Prompts must be executable in isolation** | An AI reading only `prompts/05-*.md` must be able to execute it without reading 01-04. |
| **Sequential unless parallel-safe** | If Step 2 does not depend on Step 1's output, mark `[PARALLEL SAFE]`. |
| **Max 3 target files per prompt** | If a step touches >3 files, split it. |
| **Each prompt ends with commit** | Tests pass → git commit → update TRACKING.md. No open-ended prompts. |
| **Quality gates are mandatory** | Every prompt must specify exact `task ...` commands. |

### 4.2 Prompt Structure

For each step, create `prompts/##-verb-domain-subject.md` following `prompts/_template.md`:

```
## Context
## Current State
## Scope
## Pre-flight
## Execution
## Quality Gates
## Commit
## Post-completion
## AFK Auto-Decision Log
```

### 4.3 Dual-Layer Pattern

Use the dual-layer pattern when the step is non-trivial:

| File | Purpose | Contains |
|------|---------|----------|
| `tasks/##-*.md` | Spec | Objective, constraints, schema, BDD/TDD gates, rollback triggers |
| `prompts/##-*.md` | Execution | Step-by-step commands, pre-flight, commit instructions |

**Rule:** Create `tasks/##-*.md` BEFORE creating `prompts/##-*.md`. The prompt references the task spec, but must not reference other prompts.

**Exception:** If a step is purely read-only (e.g., audit/discovery), the task file can be minimal.

---

## Step 5: Self-Review Before Presenting

Before showing the populated template to the user, verify:

- [ ] Every prompt is self-contained (no "see Prompt 03" references).
- [ ] Every prompt has a repo path reminder.
- [ ] Every prompt has exact quality gate commands using `task ...`.
- [ ] Every prompt ends with "Commit & Update" instructions.
- [ ] AFK Auto-Decision Log section is present in every prompt.
- [ ] **Prompt numbering is contiguous: 01 through N with no gaps.**
- [ ] **Every non-trivial `prompts/##-*.md` has a matching `tasks/##-*.md`.**
- [ ] `TRACKING.md` has a row for every prompt.
- [ ] `PREREQUISITES.md` has human approval gates for auth/billing/PII/DB.
- [ ] `RISKS.md` has at least one entry.
- [ ] `KNOWN_ISSUES.md` was populated from repo scan.
- [ ] Branch check in pre-flight enforces `dev`.

---

## Step 6: Present to Human for Review

Output a summary like this:

```
I have decomposed your task into N prompts under docs/projects/<project-name>/.

**Project:** <name>
**Scope:** <one-line>
**Repo:** <repo-root>
**Branch:** dev

**Prompts:**
| # | Title | Target | Est. Effort |
|---|-------|--------|-------------|
| 01 | ... | ... | ... |

**Human Approval Gates:**
- Gate A: Before touching ...

**Risks Logged:**
- R01: ...

Please review PREREQUISITES.md and sign off before I begin execution.
```

**Do NOT start execution until the human signs off PREREQUISITES.md.**

If the user is in AFK mode, present the plan and proceed with conservative choices, logging all non-trivial decisions in `DECISIONS.md`.

---

## Step 7: Execute

Once approved:

1. Ensure you are on `dev`: `git branch --show-current`.
2. Read `TRACKING.md` Step 01.
3. Read `prompts/01-*.md` and matching `tasks/01-*.md`.
4. Run pre-flight:
   ```bash
   git status
   git branch --show-current  # must be dev
   task check                 # baseline must pass
   ```
5. Execute the prompt.
6. Pass quality gates:
   ```bash
   task fmt && task vet && task lint && task test && task coverage && task build && task smoke
   ```
7. Git commit on `dev`.
8. Update `TRACKING.md` → ✅, fill commit hash, decisions, notes.
9. Update `HANDOFF.md`.
10. If decisions made, log in `DECISIONS.md`.
11. If risks materialized, update `RISKS.md`.
12. Repeat for Step 02.

---

## Step 8: Project Closeout (Mandatory Before Finishing)

After all prompts are created and the plan is presented, perform these cleanup steps. Do NOT skip them.

- [ ] **Delete `INIT.md`** — it is a one-time checklist, not a permanent document.
- [ ] **Delete `prompts/_example.md`** — it is instructional scaffolding only.
- [ ] **Strip the `## AI Agent Quickstart` section from `README.md`** — it is bootstrap instructions for the next agent, not project documentation.
- [ ] **Verify contiguous prompt numbering:** 01 through N. No gaps. Prompt 01 must exist even if it is a read-only audit/discovery step.
- [ ] **Verify task files:** Every non-trivial `prompts/##-*.md` has a matching `tasks/##-*.md`.
- [ ] **Update `TRACKING.md` status** from ⬜ Init to 🟡 In Progress.
- [ ] **Create stub runbooks** for any risks in `RISKS.md` that reference a runbook.
- [ ] **Update `HANDOFF.md`** with final state before ending session.
- [ ] **Save context** before ending session if a context-save skill is available.

---

## Example: Raw Task → Decomposed Project

### Raw Task
> "Add SQLite persistence to the transaction ledger so restarts don't lose state."

### Decomposition

| Step | Prompt | Why This Split |
|------|--------|----------------|
| 01 | Audit current transaction store and engine interfaces | Need to understand current in-memory model before designing schema. |
| 02 | Design SQLite schema and migration strategy | Schema must exist before code uses it. |
| 03 | Implement SQLite-backed `TransactionStore` | Core logic. Depends on 02. |
| 04 | Wire SQLite store into `engine` and CLI startup | Registration. Depends on 03. |
| 05 | Add unit tests for SQLite store | Testing. Depends on 03. |
| 06 | Update smoke tests to verify persistence across restart | E2E verification. Depends on 04. |

### Approval Gates
- Gate A: Before changing persistence schema or config paths

### Risks
- R01: SQLite migration breaks existing in-memory tests → mitigation: keep in-memory store as test option
- R02: Coverage drops → mitigation: add tests for new store before merging

---

## Anti-Patterns to Avoid

| Anti-Pattern | Why It's Bad | Fix |
|-------------|-------------|-----|
| One giant prompt | AI loses context, quality drops, hard to debug | Split into 3-5 file-sized steps |
| Missing quality gates | Code ships untested or unlinted | Every prompt MUST have gate checkboxes |
| Prompts reference each other | Breaks isolation, next agent can't resume | Each prompt is fully self-contained |
| No branch check in pre-flight | Commits to protected branch | Pre-flight checks `git branch --show-current` |
| Skipping KNOWN_ISSUES scan | AI wastes time fixing pre-existing bugs | Pre-flight scans known issues |
| No AFK auto-decision log | Decisions are lost, same debates repeat | AFK Auto-Decision Log is mandatory |
| Leaving template scaffolding | INIT.md, _example.md, Quickstart section clutter the workspace | Follow Step 8 closeout checklist strictly |
