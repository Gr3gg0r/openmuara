> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This project is subordinate to it.**

# OpenMuara v1

> **Status:** 🟡 In Progress | **Started:** 2026-06-27 | **Target End:** 2026-07-15
> **Scope:** Rebrand `muara` to `OpenMuara` (CLI binary: `muara`) and ship v1: SQLite persistence, universal payment API, Stripe adapter, webhook relay, Docker, and basic web UI.
> **Owner:** AI Agent (Kimi Code) | **Human Reviewer:** ___________
> **Target Repo:** `<repo-root>`
> **Product Branch:** `dev` (consolidated; no feature branches)

---

## Project Structure

```
docs/projects/openmuara-v1/
├── README.md              # This file
├── HOWTO.md               # Step-by-step guide for AI to decompose tasks
├── PREREQUISITES.md       # Human pre-flight checklist
├── TRACKING.md            # Central execution tracker
├── HANDOFF.md             # Session continuity
├── DECISIONS.md           # Decision log
├── RISKS.md               # Risk register
├── KNOWN_ISSUES.md        # Pre-existing bugs / out-of-scope areas
├── REFERENCES.md          # Links to architecture docs and vendor APIs
├── .gitignore             # Ignore ephemeral artifacts
│
├── prompts/               # Isolated execution prompts
├── tasks/                 # Detailed specs (dual-layer pattern)
├── findings/              # Research output
├── runbooks/              # Operational docs
├── screenshots/           # Visual evidence (gitignored)
├── qa/                    # Validation artifacts (gitignored)
└── state/                 # Agent state snapshots (gitignored)
```

Planning docs live in `docs/projects/openmuara-v1/` in the same repo as the product code. All product-code commits happen on the `dev` branch.

---

## Conventions (MANDATORY)

### 1. Read `AGENTS.md` first
Root `AGENTS.md` is the source of truth for branch rules, quality gates, autonomy boundaries, and code style.

### 2. Branch Boundary
- **Default working branch:** `dev`.
- **Protected branches:** `main`, `master` — never commit directly.
- **Consolidation rule:** All changes for this project happen on `dev`.
- Planning docs and product code live in the same repo, but **do not mix them in one commit**.

### 3. Prompt Isolation
- Each prompt MUST be self-contained.
- Use dual-layer pattern: `tasks/##-*.md` = spec, `prompts/##-*.md` = execution.
- Prompts MUST NOT reference other prompts.

### 4. Status Discipline
- Update `TRACKING.md` after EVERY step.
- Update `HANDOFF.md` BEFORE exiting the session.
- Log decisions in `DECISIONS.md` as they happen.
- In AFK mode, auto-log every non-trivial decision.

### 5. Quality Gates

| Gate | Command | Purpose |
|------|---------|---------|
| Format | `task fmt` | Format all Go files |
| Vet | `task vet` | Run `go vet ./...` |
| Lint | `task lint` | Run `golangci-lint run` |
| Test | `task test` | Run `go test ./...` |
| Coverage | `task coverage` | Coverage report (≥ 50% threshold) |
| Build | `task build` | Build all binaries |
| Smoke | `task smoke` | Run smoke tests |

Run `task check` where available for the full suite.

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
- [ ] Project moved to `docs/projects/archive/done/openmuara-v1/`.
