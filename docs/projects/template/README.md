> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This template is subordinate to it.**

# <PROJECT_NAME>

> **Status:** INIT | **Started:** YYYY-MM-DD | **Target End:** YYYY-MM-DD
> **Scope:** One-line description of what this project covers.
> **Owner:** AI Agent (Kimi Code) | **Human Reviewer:** ___________
> **Target Repo:** `<repo-root>`
> **Product Branch:** `dev` (consolidated; no feature branches)

---

## AI Agent Quickstart

> ⚠️ **AI AGENT: This entire section is SCAFFOLDING. After you finish initializing this project, you MUST delete this `## AI Agent Quickstart` section from `README.md`. Do not leave bootstrap instructions in the final document.**

### If you are starting this project NOW

1. Read `PREREQUISITES.md` — ensure all skills, MCPs, and environment items are ready.
2. Read `INIT.md` and run the initialization checklist.
3. Replace `<PROJECT_NAME>` and all `YYYY-MM-DD` / `___________` placeholders in this file.
4. **Delete this `## AI Agent Quickstart` section once initialized.**
5. Create real prompts in `prompts/` and tasks in `tasks/`.
6. Begin execution from `TRACKING.md` Step 01.

### If you are resuming this project

1. Read `HANDOFF.md` first — it contains the last known state.
2. Read `TRACKING.md` to see what is done vs. remaining.
3. Read `DECISIONS.md` and `RISKS.md` for constraints.
4. Pick up from the first non-DONE step in `TRACKING.md`.

### If you are creating prompts for this project

1. Read `HOWTO.md` — it explains how to decompose a raw task into prompts.
2. Read `prompts/_template.md` — it is the authoring standard.
3. Use `tasks/_template.md` if using the dual-layer pattern.
4. Every prompt MUST reference its target repo path explicitly.

---

## Project Structure

```
docs/projects/<project-name>/
├── README.md              # This file — project context and conventions
├── HOWTO.md               # Step-by-step guide for AI to decompose tasks
├── PREREQUISITES.md       # Human pre-flight: skills, MCPs, env, approval gates
├── INIT.md                # One-time setup checklist (delete after init)
├── TRACKING.md            # Central execution tracker — UPDATE AFTER EVERY STEP
├── HANDOFF.md             # Session continuity — UPDATE BEFORE EVERY EXIT
├── DECISIONS.md           # Decision log (formal + AFK auto-decisions)
├── RISKS.md               # Risk register and rollback plans
├── KNOWN_ISSUES.md        # Pre-existing bugs AI must NOT try to fix
├── REFERENCES.md          # Links to architecture docs, API specs, designs
├── .gitignore             # Ignore screenshots, logs, temp files
│
├── prompts/               # Isolated execution prompts for AI sessions
│   ├── _template.md       # Prompt authoring standard
│   ├── _example.md        # Example of a good prompt
│   └── ##-verb-domain.md  # Real prompts (numbered, self-contained)
│
├── tasks/                 # (Optional) Detailed specs — dual-layer pattern
│   ├── _template.md       # Task spec authoring standard
│   └── ##-verb-domain.md  # Real task specs
│
├── findings/              # Research output, audit results, analysis
├── runbooks/              # Operational docs for post-execution reference
├── screenshots/           # Visual evidence (gitignored)
├── qa/                    # Validation artifacts, test reports, diffs (gitignored)
└── state/                 # Agent state snapshots (session-only, gitignored)
```

Planning docs live in `docs/projects/<project-name>/` in the **same repo** (`<repo-root>`). Product code commits also happen in this repo on the `dev` branch. Planning-doc commits should still be separate from product-code commits.

---

## Conventions (MANDATORY)

### 1. Read `AGENTS.md` first
Root `AGENTS.md` is the source of truth for branch rules, quality gates, autonomy boundaries, and code style. This template does not repeat every rule — it references them.

### 2. Branch Boundary
- **Default working branch:** `dev`.
- **Protected branches:** `main`, `master` — never commit directly.
- **Consolidation rule:** All product-code changes for this project happen on `dev` unless the human explicitly requests a feature branch.
- Planning docs (this folder) and product code live in the same repo, but **do not mix them in one commit**.

### 3. Prompt Isolation
- Each prompt MUST be self-contained. An AI reading only `prompts/05-*.md` must be able to execute it.
- If a prompt needs spec context, use the dual-layer pattern: `tasks/05-*.md` = spec, `prompts/05-*.md` = execution.
- Prompts MUST NOT reference other prompts ("see Prompt 03" is forbidden). Reference task specs or this README instead.

### 4. File Naming
- Prompts: `##-verb-domain-subject.md` (e.g., `01-rename-go-module.md`)
- Tasks: `##-verb-domain-subject.md` (mirror prompt when dual-layer)
- Findings: `##-short-descriptive-name.md` (e.g., `01-provider-adapter-inventory.md`)
- Runbooks: `kebab-case-topic.md` (e.g., `investigating-webhook-retries.md`)

### 5. Status Discipline
- Update `TRACKING.md` after EVERY step: status, commit hash, decisions, notes.
- Update `HANDOFF.md` BEFORE exiting the session: last action, next action, blockers, auto-decisions.
- Log decisions in `DECISIONS.md` as they happen. Do not rely on memory.
- **AFK mode:** Every non-trivial decision MUST be auto-logged in `DECISIONS.md` (Auto-Decisions table) before finishing the prompt. "None" is a valid entry.

### 6. Screenshots
- Save to `screenshots/` or `##/screenshots/`.
- Name descriptively: `webhook-replay-success.png`, `fawry-checkout-page.png`.
- Screenshots, QA artifacts, and state snapshots are gitignored. Do not commit them.
- Planning docs (README, TRACKING, HANDOFF, DECISIONS, RISKS, prompts, tasks, findings, runbooks) commit in root `muara`.

### 7. Commit Rules
- Commit format: `type(scope): imperative description`
- Types: `feat`, `fix`, `chore`, `docs`, `refactor`, `test`, `ci`, `perf`
- Scope examples: `repo`, `server`, `webhook`, `fawry`, `stripe`, `cli`, `provider`, `config`
- Examples:
  - `feat(provider): add stripe adapter`
  - `fix(webhook): resolve race condition in delivery worker`
  - `chore(repo): update golangci-lint config`

### 8. Subagent Delegation Rules
Use the `Agent` tool when:
- Exploration requires >3 file reads or searches.
- A task is independent and can run in parallel.
- The task is read-only (explore agent) vs. write (coder agent).

Do NOT use subagents when:
- The change is <10 lines in a known file.
- You are in the middle of a step and need immediate feedback.

### 9. Regression Protection
These P0 flows must NEVER break. Verify after any change that touches:
- **Fawry provider:** charge endpoint, webhook signature verification, escape page, callback redirect.
- **Webhook delivery:** dispatch, retries, replay, store persistence.
- **CLI:** `start`, `init`, `webhook`, `version`, `doctor`.
- **Provider interface:** existing provider implementations and registry behavior.
- **Config:** environment loading, defaults, `.muara/` / `.muara/` paths.

### 10. Codebase Boundaries
- All new features belong in the existing Go module structure (`cmd/`, `internal/`, `pkg/`, `plugins/`).
- Do not introduce new top-level directories without explicit agreement.
- Keep files under 250 lines and functions under 80 lines where possible.

### 11. Quality Gates

| Stack | Command | Purpose |
|-------|---------|---------|
| Go | `task fmt` | Format all Go files |
| Go | `task vet` | Run `go vet ./...` |
| Go | `task lint` | Run `golangci-lint run` |
| Go | `task test` | Run `go test ./...` |
| Go | `task coverage` | Coverage report (≥ 50% threshold) |
| Go | `task build` | Build all binaries |
| Go | `task smoke` | Run smoke tests |

Run `task check` to run the full gate suite where available.

---

## Completion Criteria

This project is DONE when:
- [ ] All steps in `TRACKING.md` are ✅ DONE.
- [ ] All quality gates in `TRACKING.md` show ✅ Pass.
- [ ] `HANDOFF.md` is updated with final state.
- [ ] `DECISIONS.md` captures all major decisions.
- [ ] `RISKS.md` shows all risks as mitigated or accepted.
- [ ] `README.md` header updated to `Status: COMPLETE | Date: YYYY-MM-DD`.
- [ ] `KNOWN_ISSUES.md` and `REFERENCES.md` are up to date.
- [ ] `dev` branch is green and pushed.
- [ ] Project moved to `docs/projects/archive/done/<project-name>/`.
