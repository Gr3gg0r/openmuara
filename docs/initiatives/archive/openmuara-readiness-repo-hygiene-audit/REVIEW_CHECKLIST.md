> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Repository Hygiene Review Checklist

> **Created:** 2026-07-10  
> **Status:** Draft

Use this checklist when reviewing the implementation of the repo-hygiene initiative.

---

## Pre-review

- [ ] Branch is `feat/readiness-repo-hygiene-audit` or a PR into `dev`.
- [ ] `git status` shows only expected files.
- [ ] No secrets, local configs, or build artifacts are present.

---

## P01 — Secret & artifact audit

- [ ] `gitleaks detect --source . --verbose` reports zero leaks.
- [ ] `git ls-files` contains no `coverage.out`, `bin/muara`, `.env`, `.muara/config.yml`, or `.toyol/` files.
- [ ] `KNOWN_ISSUES.md` records the scan date and result.

---

## P02 — Naming & branding cleanup

- [ ] `grep -Ril --exclude-dir=.git --exclude-dir=node_modules toyol .` returns only the root directory name.
- [ ] `.gitignore` has a commented legacy entry for `.toyol/`.
- [ ] `AGENTS.md` and migration docs use OpenMuara branding consistently.

---

## P03 — Ignore-rule hardening

- [ ] `.gitignore` has no duplicate entries.
- [ ] `.gitignore` includes `.gstack/` and `.playwright-mcp/`.
- [ ] Generated files in `internal/ui/dashboard-dist/` are ignored except the placeholder `index.html`.
- [ ] `.dockerignore` includes `node_modules/`, `coverage.html`, generated UI assets, and local AI directories.
- [ ] `docker build -t openmuara:hygiene-test -f Dockerfile .` succeeds.

---

## P04 — Governance consolidation

- [ ] Root `SECURITY.md` is a redirect to `.github/SECURITY.md`.
- [ ] `.github/SUPPORT.md` exists and renders on GitHub.
- [ ] `MAINTAINERS.md` exists with named maintainers and areas.
- [ ] `.github/FUNDING.yml` exists (placeholder acceptable).
- [ ] `CONTRIBUTING.md` mentions `dev` default branch, Conventional Commits, pre-commit hooks, and bug-register format.

---

## P05 — GitHub metadata

- [ ] `.github/ISSUE_TEMPLATE/provider_request.yml` renders.
- [ ] `.github/ISSUE_TEMPLATE/docs_issue.yml` renders.
- [ ] `.github/ISSUE_TEMPLATE/config.yml` links to Security Advisories and Discussions.
- [ ] `.github/release.yml` exists and uses recognized labels.
- [ ] `APPENDIX.md` contains the label taxonomy and repo-settings checklist.

---

## P06 — Branch & release discipline

- [ ] `CONTRIBUTING.md` documents Conventional Commits.
- [ ] `runbooks/release.md` or `APPENDIX.md` documents stale-branch cleanup.
- [ ] `main` is fast-forwarded to `dev` (or documented exception).
- [ ] Merged worktrees are removed.
- [ ] Suspended/stale branches are deleted or archived.

---

## P07 — Developer experience

- [ ] `.editorconfig` exists and covers Go, YAML, Markdown, JSON, TSX, CSS, shell.
- [ ] `.gitattributes` marks generated files as `linguist-generated`.
- [ ] `.pre-commit-config.yaml` includes `shellcheck` and `actionlint`.
- [ ] `pre-commit run --all-files` passes (markdown hook advisory only).

---

## P08 — Publication readiness

- [ ] `VERSION` matches the upcoming release section in `CHANGELOG.md` (target `1.1.0`).
- [ ] `README.md` has a Contributing section linking to `CONTRIBUTING.md` and `CODE_OF_CONDUCT.md`.
- [ ] `README.md` includes the AI-assisted development disclosure.
- [ ] Badge URLs point to `openmuara/openmuara` (post-transfer verification).
- [ ] `APPENDIX.md` contains the repo-transfer checklist.

---

## Final quality gates

- [ ] `go build ./...`
- [ ] `go vet ./...`
- [ ] `go test ./...`
- [ ] `golangci-lint run`
- [ ] `actionlint .github/workflows/*.yml`
- [ ] `gitleaks detect --source .`
- [ ] `pre-commit run --all-files`
- [ ] `task quality`

---

## Sign-off

- [ ] Reviewer: ___________
- [ ] Date: ___________
- [ ] Approved for merge to `dev`: ⬜
