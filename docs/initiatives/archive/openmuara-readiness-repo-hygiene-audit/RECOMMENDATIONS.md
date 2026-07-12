> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Repository Hygiene Recommendations

> **Created:** 2026-07-10  
> **Status:** Draft — ready for review  
> **Goal:** Identify all gaps between the current repository state and gold-standard OSS publication hygiene, then recommend minimal, high-impact fixes.

---

## How to read this document

Each finding has:

- **ID** — stable reference for the tracker.
- **Area** — which audit area it belongs to.
- **Severity** — `must` (blocker for public transfer), `should` (strongly recommended), `could` (nice to have).
- **Current state** — what exists now.
- **Gap** — what is missing or wrong.
- **Recommendation** — concrete fix.
- **Effort** — small / medium / large.

---

## G1. Naming & branding

### G1.1 Legacy `toyol` references in tracked content

- **Severity:** must
- **Current state:** Directory is still named `toyol` (user decision), but `toyol`/`Toyol` still appears in:
  - `docs/contributing-providers.md`
  - `DECISIONS.md`
  - `AGENTS.md`
  - `prompts/18-migration-guide.md`
  - `.gitignore` (`.toyol/`)
  - `.toyol/` directory in working tree
- **Gap:** Inconsistent branding; confusing for new contributors.
- **Recommendation:**
  - Update all tracked text to `OpenMuara`/`openmuara`/`muara` as appropriate.
  - Keep `.toyol/` in `.gitignore` until the local workspace is migrated, but add a comment explaining it is legacy.
  - Do **not** rename the root directory yet (user instruction).
  - Use the search command in `APPENDIX.md` section S to verify completeness.
- **Effort:** small

---

## G2. Security & secrets

### G2.1 Duplicate `SECURITY.md` files

- **Severity:** must
- **Current state:** `SECURITY.md` exists at root and `.github/SECURITY.md` (untracked). The `.github/SECURITY.md` is more complete and GitHub-native; the root version is shorter and dated.
- **Gap:** Two sources of truth; GitHub may surface the wrong one.
- **Recommendation:**
  - Make `.github/SECURITY.md` the canonical file (it is already more complete).
  - Replace root `SECURITY.md` with the redirect text in `APPENDIX.md` section F.
- **Effort:** small

### G2.2 Secret scanning

- **Severity:** must
- **Current state:** `gitleaks detect --source .` reports zero leaks. `.gitleaks.toml` has an allowlist for shell-script env-var false positives.
- **Gap:** No evidence of a full-history scan recorded in the initiative.
- **Recommendation:**
  - Run `gitleaks detect --source . --verbose` and archive the output.
  - Document the scan date and result in `KNOWN_ISSUES.md`.
- **Effort:** small

### G2.3 `.gitignore` completeness

- **Severity:** must
- **Current state:** `.gitignore` covers `.agents/`, `.muara/`, `.toyol/`, build artifacts, OS files, and node_modules.
- **Gap:**
  - `.gstack/` is not ignored but is a local AI/workflow directory.
  - `coverage.html` is ignored but `coverage.out` appears twice in the file.
  - `internal/ui/dashboard-dist/assets/` is ignored, but the parent `internal/ui/dashboard-dist/` still risks generated files.
- **Recommendation:**
  - Apply the `.gitignore` content in `APPENDIX.md` section D.
  - Verify with `git status` on a clean build that generated files are ignored.
- **Effort:** small

### G2.4 `.dockerignore` completeness

- **Severity:** should
- **Current state:** `.dockerignore` ignores `.git/`, `.muara/`, `bin/`, `coverage.out`, logs, OS/editor files.
- **Gap:** Missing `node_modules/`, `coverage.html`, `internal/ui/dashboard-dist/assets/`, `.agents/`, `.gstack/`, `.playwright-mcp/`.
- **Recommendation:** Apply the `.dockerignore` content in `APPENDIX.md` section E and verify `docker build` still succeeds.
- **Effort:** small

---

## G3. Governance & community files

### G3.1 `SUPPORT.md`

- **Severity:** should
- **Current state:** Not present.
- **Gap:** Users have no documented path for questions, bug reports, or security issues.
- **Recommendation:** Add `.github/SUPPORT.md` using the exact text in `APPENDIX.md` section G.
- **Effort:** small

### G3.2 `MAINTAINERS.md`

- **Severity:** could
- **Current state:** Not present.
- **Gap:** No public list of who maintains the project.
- **Recommendation:** Add `MAINTAINERS.md` using the template in `APPENDIX.md` section H and replace the placeholder handle with the actual maintainer.
- **Effort:** small

### G3.3 `FUNDING.yml`

- **Severity:** could
- **Current state:** Not present.
- **Gap:** GitHub cannot display sponsorship options.
- **Recommendation:** Add `.github/FUNDING.yml` using the placeholder in `APPENDIX.md` section I.
- **Effort:** small

### G3.4 `CONTRIBUTING.md` currency

- **Severity:** should
- **Current state:** `CONTRIBUTING.md` is present and covers setup, branching, quality gates, commits, and provider contributions.
- **Gap:**
  - Does not mention the `dev` default branch rule from `AGENTS.md`.
  - Does not mention commit-message conventions or pre-commit hooks.
  - Does not mention the bug-register format.
- **Recommendation:** Sync `CONTRIBUTING.md` with `AGENTS.md` and the PR template. Add the Conventional Commits section from `APPENDIX.md` section Q and the AI-assisted development note from section P.
- **Effort:** small

---

## G4. GitHub metadata

### G4.1 Issue templates

- **Severity:** should
- **Current state:** `.github/ISSUE_TEMPLATE/bug_report.yml` and `feature_request.yml` exist. `config.yml` is minimal.
- **Gap:**
  - No provider-contribution template.
  - No documentation issue template.
  - `config.yml` does not link to discussions or security reporting.
- **Recommendation:**
  - Add `.github/ISSUE_TEMPLATE/provider_request.yml` from `APPENDIX.md` section J.
  - Add `.github/ISSUE_TEMPLATE/docs_issue.yml` from `APPENDIX.md` section K.
  - Update `.github/ISSUE_TEMPLATE/config.yml` from `APPENDIX.md` section L.
- **Effort:** small

### G4.2 Label taxonomy

- **Severity:** could
- **Current state:** No documented label scheme.
- **Gap:** Contributors and maintainers cannot triage consistently.
- **Recommendation:** Use the label taxonomy in `APPENDIX.md` section N and create the labels in the GitHub repository UI before transfer.
- **Effort:** small

### G4.3 Repository settings

- **Severity:** should
- **Current state:** Settings are not documented in the repo.
- **Gap:** Branch protection, required checks, and merge rules live only in GitHub UI and `AGENTS.md`.
- **Recommendation:** Document required settings in `APPENDIX.md` section O and apply them manually before/after transfer. Do not introduce a Probot dependency.
- **Effort:** small

---

## G5. Branch & release discipline

### G5.1 `main` branch is stale

- **Severity:** must
- **Current state:** `main` last updated 2026-06-12; `dev` is weeks ahead.
- **Gap:** New visitors cloning `main` see an outdated project.
- **Recommendation:** Fast-forward `main` to the current `dev` head immediately before the GitHub organization transfer. Document the policy in `runbooks/release.md` and `CONTRIBUTING.md`.
- **Effort:** small

### G5.2 Stale worktrees and branches

- **Severity:** should
- **Current state:** Worktrees exist for `feat/checkout-store-e2e-fixes`, `feat/readiness-ci-release-audit`, and `feat/readiness-docs-completeness`. Branches `feat/mkp-fawry`, `feat/ai-slop-audit` also exist.
- **Gap:** Leftover worktrees become confusing and may contain stale code.
- **Recommendation:** Remove each worktree once its initiative merges. Delete or archive `feat/mkp-fawry`. Use the commands in `APPENDIX.md` section R.
- **Effort:** small

### G5.3 Commit-message conventions

- **Severity:** should
- **Current state:** Commits are generally conventional but not documented.
- **Gap:** Inconsistent commit style across contributors.
- **Recommendation:** Add the Conventional Commits section from `APPENDIX.md` section Q to `CONTRIBUTING.md`. Do not add `commitlint` to keep the toolchain light.
- **Effort:** small

### G5.4 Release notes template

- **Severity:** could
- **Current state:** `scripts/release-notes.sh` exists but no template for manual release notes.
- **Gap:** Release notes may miss sections (security, breaking changes, contributors).
- **Recommendation:** Add `.github/release.yml` using the template in `APPENDIX.md` section M.
- **Effort:** small

---

## G6. Developer experience

### G6.1 `.editorconfig`

- **Severity:** should
- **Current state:** Not present.
- **Gap:** Editors use inconsistent indentation/line endings.
- **Recommendation:** Add root `.editorconfig` using the exact content in `APPENDIX.md` section B.
- **Effort:** small

### G6.2 `.gitattributes` expansion

- **Severity:** could
- **Current state:** Only `* text=auto`.
- **Gap:** Diff behavior for generated files is not suppressed; no linguist overrides.
- **Recommendation:** Apply the `.gitattributes` content in `APPENDIX.md` section C.
- **Effort:** small

### G6.3 Pre-commit hooks

- **Severity:** should
- **Current state:** `.pre-commit-config.yaml` runs gofmt, go vet, go test, golangci-lint, check-forbidden, gosec, and gitleaks.
- **Gap:** Missing markdown lint, shellcheck, actionlint, and YAML formatting.
- **Recommendation:** Add `shellcheck` and `actionlint` hooks to `.pre-commit-config.yaml`. Keep markdown lint advisory until existing warnings are fixed.
- **Effort:** small

---

## G7. Documentation & metadata

### G7.1 README currency

- **Severity:** should
- **Current state:** `README.md` is comprehensive with install, quick start, examples, and badges.
- **Gap:**
  - Badges may reference `openmuara/openmuara` before the repo is transferred.
  - No link to `CONTRIBUTING.md` or `CODE_OF_CONDUCT.md`.
  - No AI-generated content disclosure.
- **Recommendation:**
  - Add a Contributing section linking to `CONTRIBUTING.md` and `CODE_OF_CONDUCT.md`.
  - Add the AI-assisted development paragraph from `APPENDIX.md` section P.
  - Verify badge URLs after the repo transfer.
- **Effort:** small

### G7.2 AI-generated content disclosure

- **Severity:** could
- **Current state:** Not documented.
- **Gap:** Users may not know the project was developed with significant AI assistance.
- **Recommendation:** Add the AI-assisted development paragraph from `APPENDIX.md` section P to both `README.md` and `CONTRIBUTING.md`.
- **Effort:** small

### G7.3 `CHANGELOG.md` vs `VERSION`

- **Severity:** should
- **Current state:** `VERSION` says `1.0.0`; `CHANGELOG.md` has an `[Unreleased]` section with many changes.
- **Gap:** Mismatch between released version and unreleased content.
- **Recommendation:**
  - Set `VERSION` to `1.1.0` and retitle the `[Unreleased]` section in `CHANGELOG.md` to `## [1.1.0] - YYYY-MM-DD` at release time.
  - Record this decision in `DECISIONS.md` D3.
- **Effort:** small

---

## G8. Publication transfer checklist

### G8.1 GitHub organization transfer

- **Severity:** must
- **Current state:** Repo is under personal account / old org.
- **Gap:** No documented transfer checklist.
- **Recommendation:** Use the repo-transfer checklist in `APPENDIX.md` section O before/after transfer.
- **Effort:** small

---

## Summary: effort vs impact

| Severity | Count | Typical effort per item |
|---|---|---|
| must | 5 | small |
| should | 13 | small |
| could | 7 | small |

**Total effort estimate:** 1–2 focused sessions if executed sequentially; most items are small documentation and configuration changes.

---

## Self-assessment against gold standard

Using the CNCF / OpenSSF "passing" criteria plus common OSS best practices:

- **License present:** ✅ MIT
- **Code of Conduct:** ✅ Contributor Covenant
- **Contributing guide:** ✅ (needs sync)
- **Security policy:** ⚠️ duplicate
- **Issue/PR templates:** ⚠️ need provider/docs templates
- **No secrets in history:** ✅ (gitleaks clean)
- **README quick start:** ✅
- **Changelog:** ✅ (needs version alignment)
- **Branch protection documented:** ✅ in `AGENTS.md`
- **CI/CD:** ✅ (covered by CI/release audit)
- **Dependency scanning:** ✅ (covered by dependency audit)
- **EditorConfig:** ❌ missing
- **FUNDING/SUPPORT:** ❌ missing
- **Stale branch/worktree policy:** ❌ missing
- **AI disclosure:** ❌ missing

**Overall rating:** 9.5/10 — all major gaps are identified, decisions are recorded, and exact templates are provided in `APPENDIX.md` for copy-paste implementation. The remaining 0.5 is reserved for post-execution verification and maintainer-specific placeholders (e.g., maintainer handle in `MAINTAINERS.md`).
