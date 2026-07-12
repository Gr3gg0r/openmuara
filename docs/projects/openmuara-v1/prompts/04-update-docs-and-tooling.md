> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This project is subordinate to it.**

## Update docs, scripts, CI, and tooling for OpenMuara / `muara`

### Context
After renaming the module, binary, and config paths, all remaining references to `muara` in documentation and tooling must be updated. The project is called OpenMuara, the module is `github.com/openmuara/openmuara`, and the CLI is `muara`.

### Current State
- **Repo:** `<repo-root>`
- **Branch:** `dev`
- **Last step:** Prompt 03 completed (config paths rename).

### Scope
- **In scope:**
  - Update `README.md` with new name, install instructions, and quick start.
  - Update `AGENTS.md` references to `muara` / CLI name.
  - Update `.github/workflows/*.yml` to build `./cmd/muara`.
  - Update `scripts/*.sh` to reference `muara`.
  - Update `.pre-commit-config.yaml` if needed.
  - Update other docs under `docs/`.
- **Out of scope:**
  - Git history.
  - Renaming the GitHub org/repo itself (can be done later).

### Pre-flight
```bash
cd <repo-root>
git status
git branch --show-current  # must be dev
grep -R "openmuara" --include="*.md" --include="*.yml" --include="*.yaml" --include="*.sh" . | grep -v ".git/" | grep -v "docs/projects/openmuara-v1/" | grep -v "session-6786b9f5" | head -50
```

### Execution
1. Update `README.md`:
   - Project name and tagline.
   - Install command (`go install github.com/openmuara/openmuara/cmd/muara@latest`).
   - Quick start (`muara init`, `muara start`).
2. Update `AGENTS.md` to reference `muara` / OpenMuara where appropriate.
3. Update `.github/workflows/*.yml` build/test paths to `./cmd/muara` and `muara` binary.
4. Update `scripts/*.sh` binary names.
5. Update `.pre-commit-config.yaml` if it references `muara`.
6. Update other markdown docs, excluding this planning workspace.

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
```bash
git add -A
git commit -m "docs(repo): update README, AGENTS, CI, and scripts for OpenMuara / muara rebrand"
```

### Post-completion
1. Update `TRACKING.md` Step 04 → ✅, fill commit hash.
2. Log any auto-decisions in `DECISIONS.md`.
3. Update `HANDOFF.md`.
