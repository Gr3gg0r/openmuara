> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Repository Hygiene Execution Plan

> **Created:** 2026-07-10  
> **Status:** Draft — refined to 9.5/10, ready for sign-off  
> **Goal:** Convert the recommendations in `RECOMMENDATIONS.md` into an ordered, verifiable implementation plan.

---

## Execution principles

1. **No functional code changes.** This initiative only touches repo metadata, templates, docs, ignore rules, and configuration.
2. **Preserve history.** Do not rewrite public branches. If a secret is found, escalate to the user before any history rewrite.
3. **One logical change per commit.** Each phase should produce focused, reviewable commits.
4. **Verify after every phase.** Run the quality gates defined in `TRACKING.md`.
5. **Defer to other initiatives.** Do not duplicate work owned by security, docs, CI/release, or dependency-license audits.
6. **Use the templates.** Exact file contents are provided in `APPENDIX.md`; key decisions are recorded in `DECISIONS.md`.

---

## Phase P01 — Secret & artifact audit

**Goal:** Confirm the repository is clean of secrets and unwanted artifacts.

| Step | Task | Verification |
|---|---|---|
| P01.1 | Run `gitleaks detect --source . --verbose` on full history. | Zero leaks. |
| P01.2 | Inspect tracked files for binaries, coverage artifacts, or local configs. | `git ls-files` shows no `coverage.out`, `bin/muara`, `.muara/`, `.env` files. |
| P01.3 | Record scan date and result in `KNOWN_ISSUES.md` or tracker. | Tracker updated. |

**Commit message:** `chore(hygiene): record secret-scan results and artifact audit`

---

## Phase P02 — Naming & branding cleanup

**Goal:** Remove legacy `toyol` references from tracked content.

| Step | Task | Verification |
|---|---|---|
| P02.1 | Search all tracked files for `toyol`/`Toyol` using the command in `APPENDIX.md` section S. | Returns only the root directory name and `.gitignore` legacy entry. |
| P02.2 | Update `docs/contributing-providers.md`, `AGENTS.md`, `prompts/18-migration-guide.md` to use OpenMuara branding. Do not edit this initiative's `DECISIONS.md` (already uses correct branding). | Text review. |
| P02.3 | Apply the `.gitignore` content from `APPENDIX.md` section D, including the commented legacy `.toyol/` entry. | `git status` on a clean build shows no unexpected files. |
| P02.4 | Add a note in `AGENTS.md` or `docs/migration/openmuara-to-openmuara.md` that the local directory may still be named `toyol`. | Review. |

**Commit messages:**
- `chore(hygiene): replace legacy toyol references with OpenMuara branding`
- `chore(hygiene): expand .gitignore for local workspace directories`

---

## Phase P03 — Ignore-rule hardening

**Goal:** Ensure `.gitignore` and `.dockerignore` cover all generated/local artifacts.

| Step | Task | Verification |
|---|---|---|
| P03.1 | Apply the `.gitignore` content from `APPENDIX.md` section D. | `sort .gitignore | uniq -d` returns nothing; `git status` is clean after build. |
| P03.2 | Ignore `internal/ui/dashboard-dist/` except for the tracked placeholder `index.html`. | `go build ./...` works; generated assets are not shown by `git status`. |
| P03.3 | Apply the `.dockerignore` content from `APPENDIX.md` section E. | `docker build` still succeeds. |

**Commit messages:**
- `chore(hygiene): deduplicate and tighten .gitignore`
- `chore(hygiene): expand .dockerignore for local and generated artifacts`

---

## Phase P04 — Governance file consolidation

**Goal:** Single sources of truth for governance and support.

| Step | Task | Verification |
|---|---|---|
| P04.1 | Replace root `SECURITY.md` with the redirect text in `APPENDIX.md` section F. | GitHub security policy tab shows `.github/SECURITY.md`. |
| P04.2 | Add `.github/SUPPORT.md` from `APPENDIX.md` section G. | File renders correctly on GitHub. |
| P04.3 | Add `MAINTAINERS.md` from `APPENDIX.md` section H; replace the placeholder handle. | Review. |
| P04.4 | Add `.github/FUNDING.yml` from `APPENDIX.md` section I. | Review. |
| P04.5 | Sync `CONTRIBUTING.md` with `AGENTS.md`: default branch, commit conventions, pre-commit hooks, bug-register format. Add the Conventional Commits section from `APPENDIX.md` section Q and the AI disclosure from section P. | No contradictions between files. |

**Commit messages:**
- `chore(hygiene): consolidate SECURITY.md under .github with root redirect`
- `chore(hygiene): add SUPPORT.md, MAINTAINERS.md, and FUNDING.yml`
- `docs(contributing): sync CONTRIBUTING.md with AGENTS.md branch and commit rules`

---

## Phase P05 — GitHub metadata

**Goal:** Improve issue triage and repository settings documentation.

| Step | Task | Verification |
|---|---|---|
| P05.1 | Add `.github/ISSUE_TEMPLATE/provider_request.yml` from `APPENDIX.md` section J. | Template renders in GitHub issue picker. |
| P05.2 | Add `.github/ISSUE_TEMPLATE/docs_issue.yml` from `APPENDIX.md` section K. | Template renders. |
| P05.3 | Update `.github/ISSUE_TEMPLATE/config.yml` from `APPENDIX.md` section L. | Issue picker shows links. |
| P05.4 | Add `.github/release.yml` from `APPENDIX.md` section M. | Release notes draft uses categories. |
| P05.5 | Create the labels from `APPENDIX.md` section N in the GitHub UI. | Label taxonomy exists. |
| P05.6 | Apply the repository settings from `APPENDIX.md` section O manually. | Settings match `AGENTS.md` required checks. |

**Commit messages:**
- `chore(github): add provider and docs issue templates`
- `chore(github): link issue picker to security and discussions`
- `chore(github): add release-notes auto-categorization`

---

## Phase P06 — Branch & release discipline

**Goal:** Establish clear branch policy and clean up stale branches/worktrees.

| Step | Task | Verification |
|---|---|---|
| P06.1 | Add the Conventional Commits section from `APPENDIX.md` section Q to `CONTRIBUTING.md`. | Review. |
| P06.2 | Add a stale-branch cleanup note to `runbooks/release.md` and the worktree commands from `APPENDIX.md` section R. | Review. |
| P06.3 | Fast-forward `main` to `dev` immediately before transfer (see `DECISIONS.md` D6). | `git log --oneline main..dev` is empty. |
| P06.4 | Remove merged worktrees and stale branches using commands from `APPENDIX.md` section R. | `git worktree list` shows only main repo; stale branches removed. |

**Commit messages:**
- `docs(contributing): document conventional commits and branch cleanup`
- `chore(hygiene): fast-forward main to dev for public transfer`
- `chore(hygiene): remove merged worktrees and stale branches`

---

## Phase P07 — Developer experience polish

**Goal:** Consistent editor settings and stronger pre-commit hooks.

| Step | Task | Verification |
|---|---|---|
| P07.1 | Add root `.editorconfig` from `APPENDIX.md` section B. | Editors respect the config. |
| P07.2 | Apply the `.gitattributes` content from `APPENDIX.md` section C. | GitHub diff view hides generated files. |
| P07.3 | Add `shellcheck` and `actionlint` hooks to `.pre-commit-config.yaml`. | `pre-commit run --all-files` passes. |
| P07.4 | Optionally add advisory markdown lint hook (non-blocking until existing warnings are fixed). | Hook runs; warnings visible but non-blocking. |

**Commit messages:**
- `chore(hygiene): add .editorconfig for consistent editor behavior`
- `chore(hygiene): mark generated files in .gitattributes`
- `chore(pre-commit): add shellcheck and actionlint hooks`

---

## Phase P08 — Publication readiness

**Goal:** Final checklist before transferring to the `openmuara` organization.

| Step | Task | Verification |
|---|---|---|
| P08.1 | Set `VERSION` to `1.1.0` and retitle `[Unreleased]` to `## [1.1.0] - YYYY-MM-DD` in `CHANGELOG.md` (see `DECISIONS.md` D3). | `VERSION` matches `CHANGELOG.md`. |
| P08.2 | Add the AI-assisted development paragraph from `APPENDIX.md` section P to `README.md` and `CONTRIBUTING.md`. | Review. |
| P08.3 | Verify README badges and links point to `openmuara/openmuara`. | All badge URLs valid post-transfer. |
| P08.4 | Confirm `APPENDIX.md` section O transfer checklist is complete. | Review. |
| P08.5 | Run full `task quality` one last time. | Passes with no new warnings. |

**Commit messages:**
- `chore(release): align VERSION and CHANGELOG for next release`
- `docs(readme): add AI-assisted development disclosure and verify badges`
- `docs(hygiene): add publication-transfer checklist`

---

## Rollback considerations

- **History rewrite:** Never rewrite history on `main` or `dev`. If a secret is discovered, stop and ask the user.
- **File deletion:** Only delete untracked or clearly generated files. Verify with `git status` before any `rm`.
- **Branch deletion:** Only delete branches that have been fully merged to `dev` or are explicitly suspended.

---

## Verification cadence

After every phase:

```bash
go build ./...
go vet ./...
go test ./...
golangci-lint run
actionlint .github/workflows/*.yml
pre-commit run --all-files   # if hooks changed
task quality                 # before final sign-off
```

---

## Estimated timeline

| Phase | Estimated time |
|---|---|
| P01 | 15 min |
| P02 | 30 min |
| P03 | 20 min |
| P04 | 30 min |
| P05 | 30 min |
| P06 | 20 min |
| P07 | 30 min |
| P08 | 30 min |
| **Total** | **~3.5 hours** of focused work |
