> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# How to Decompose Tasks for OpenMuara

> **Purpose:** Guide for AI agents to take a raw task and produce a complete, executable initiative using this template.
> **Input:** A raw task description from the user.
> **Output:** A populated `docs/initiatives/<initiative-name>/` folder ready for execution.

---

## When to Use This Template

Use this template when the task is:
- Multi-phase (>1 prompt needed)
- Touches multiple files or modules
- Has regression risk
- Needs tracking across sessions

**Do NOT use this template** for single-session fixes of <10 lines.

---

## Decomposition Workflow

```
Raw Task → Understand → Scaffold → Populate → Review → Execute
```

### Step 1: Understand the Task

Identify:

| Question | Why It Matters |
|----------|----------------|
| What is the problem? | Defines "What This Solves" sections. |
| Which files/modules are affected? | Determines target files and blast radius. |
| Is this multi-phase? | Decides prompt count. |
| Does it touch P0 integrations? | Triggers human approval gates. |
| Does it touch auth/billing/PII/DB? | Triggers autonomy boundary checks per `AGENTS.md`. |

### Step 2: Scaffold the Initiative Folder

Create `docs/initiatives/<initiative-kebab-case-name>/` with:

```bash
mkdir -p docs/initiatives/<initiative-name>/{prompts,tasks,findings,runbooks,screenshots,qa,state}
```

### Step 3: Populate Shell Files

Fill every placeholder in:

- `README.md` — name, dates, scope, repo, branch
- `PREREQUISITES.md` — skills, MCPs, env checklist, approval gates, scope guardrails
- `TRACKING.md` — prompt inventory, quality gate commands
- `RISKS.md` — at least one risk entry
- `KNOWN_ISSUES.md` — scan repo for pre-existing bugs
- `REFERENCES.md` — architecture docs, API specs, vendor docs
- `DECISIONS.md` & `HANDOFF.md` — placeholders only

### Step 4: Decompose into Prompts

Create `prompts/##-verb-domain-subject.md` for each logical step.

Rules:
- One prompt = one logical step.
- Prompts must be executable in isolation.
- Sequential unless marked `[PARALLEL SAFE]`.
- Max ~3 target files per prompt; split if larger.
- Each prompt ends with commit + update TRACKING.md.
- Quality gates are mandatory.

### Step 5: Dual-Layer Pattern (Optional)

For non-trivial prompts, create both:

| File | Purpose |
|------|---------|
| `tasks/##-*.md` | Spec: objective, constraints, schema, test expectations |
| `prompts/##-*.md` | Execution: commands, pre-flight, commit instructions |

### Step 6: Self-Review

Before presenting to the user, verify:
- [ ] Every prompt is self-contained (no "see Prompt 03").
- [ ] Every prompt specifies exact quality gate commands.
- [ ] Every prompt ends with "Commit & Update" instructions.
- [ ] Prompt numbering is contiguous: 01, 02, 03 ...
- [ ] Every non-trivial prompt has a matching `tasks/##-*.md`.
- [ ] `TRACKING.md` has a row for every prompt.
- [ ] `RISKS.md` has at least one entry.
- [ ] `KNOWN_ISSUES.md` was populated from repo scan.

### Step 7: Present to Human

```
I have decomposed your task into N prompts under docs/initiatives/<initiative-name>/.

**Project:** <name>
**Scope:** <one-line>
**Target Repo:** <repo-root>/
**Branch:** dev

**Prompts:**
| # | Title | Target | Est. Effort |
|---|-------|--------|-------------|
| 01 | ... | ... | ... |

**Human Approval Gates:**
- Gate A: ...

**Risks Logged:**
- R01: ...

Please review PREREQUISITES.md and sign off before I begin execution.
```

### Step 8: Execute

1. Read `TRACKING.md` Step 01.
2. Read `prompts/01-*.md` and matching `tasks/01-*.md`.
3. Run pre-flight: `git status`, `git branch --show-current` must be `dev`.
4. Execute the prompt.
5. Pass quality gates: `go build ./...`, `go test ./...`, `golangci-lint run`, `./scripts/smoke-test.sh`.
6. Commit on `dev`.
7. Update `TRACKING.md` → ✅, fill commit hash.
8. Update `HANDOFF.md`.
9. Log decisions in `DECISIONS.md`.
10. Repeat for next step.

---

## Project Closeout Checklist

After all prompts are complete:
- [ ] Delete `INIT.md`.
- [ ] Delete `prompts/_example.md`.
- [ ] Strip the `## AI Agent Quickstart` section from `README.md`.
- [ ] Verify contiguous prompt numbering.
- [ ] Verify task files match prompts.
- [ ] Update `TRACKING.md` status to ✅ Complete.
- [ ] Create stub runbooks for any risks referencing them.
- [ ] Update `HANDOFF.md` with final state.
- [ ] Move initiative to `docs/initiatives/archive/done/<initiative-name>/`.
