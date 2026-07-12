> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This project is subordinate to it.**

## Run full quality gates to complete the rebrand

### Context
The rebrand is nearly complete. This step performs a final sweep for remaining `muara` references and runs all quality gates to ensure `dev` is green.

### Current State
- **Repo:** `<repo-root>`
- **Branch:** `dev`
- **Last step:** Prompt 04 completed (docs and tooling update).

### Scope
- **In scope:**
  - Search for remaining `muara` references across the repo.
  - Run all quality gates.
  - Fix any rebrand-related failures.
- **Out of scope:**
  - New features.
  - Pre-existing bugs unrelated to rebrand.

### Pre-flight
```bash
cd <repo-root>
git status
git branch --show-current  # must be dev
grep -R "openmuara" --include="*.go" --include="*.md" --include="*.yml" --include="*.yaml" --include="*.sh" --include="*.json" . | grep -v ".git/" | grep -v "docs/projects/openmuara-v1/" | grep -v "session-6786b9f5"
```

### Execution
1. Review every remaining `muara` reference. Rename unless it is intentionally historical (e.g., old commit messages in changelog).
2. Run quality gates.
3. Fix any failures.

### Quality Gates
```bash
task fmt
task vet
task lint
task test
task coverage
task build
task smoke
```

### Commit
If any fixes were needed:
```bash
git add -A
git commit -m "chore(repo): final rebrand sweep and quality gate fixes"
```

If no fixes were needed, this step is verification-only; still update `TRACKING.md`.

### Post-completion
1. Update `TRACKING.md` Step 05 → ✅, fill commit hash or `N/A`.
2. Log any auto-decisions in `DECISIONS.md`.
3. Update `HANDOFF.md`.
4. Push `dev` to origin if not already pushed.
