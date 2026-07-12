> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Repository Hygiene Audit

> **Status:** ⬜ Draft | **Started:** 2026-07-08  
> **Scope:** Prepare the repository as a polished, trustworthy, and maintainable open-source artifact: git history, ignore rules, templates, metadata, branch discipline, and publication hygiene.  
> **Owner:** AI Agent (Kimi Code) | **Human Reviewer:** ___________  
> **Target Repo:** `<repo-root>/`  
> **Product Branch:** `feat/readiness-repo-hygiene-audit` (to be created when work starts)

---

## Why this matters

A public repository is a long-term artifact. Before transferring to a new GitHub organization and opening it to contributors, we must ensure:

- No leaked secrets, local configs, or build artifacts in history or working tree.
- Consistent naming (OpenMuara, not legacy `toyol`) across all published content.
- Clear contribution paths via templates, docs, and governance files.
- Branch and release discipline that protects `main` and keeps history reviewable.
- Metadata and settings that signal a professional, sustainable project.

---

## Initiative structure

```
docs/initiatives/openmuara-readiness-repo-hygiene-audit/
├── README.md              # This file
├── TRACKING.md            # Central execution tracker
├── RECOMMENDATIONS.md     # Gap analysis and gold-standard recommendations
├── EXECUTION_PLAN.md      # Step-by-step implementation plan
├── APPENDIX.md            # Copy-paste templates and checklists
├── DECISIONS.md           # Recorded decisions
├── RISKS.md               # Risk register
├── KNOWN_ISSUES.md        # Catalog of hygiene findings
└── REVIEW_CHECKLIST.md    # Implementation review checklist
```

---

## Audit areas

1. **Git history & secrets** — full-history scan, ignore-rule audit, committed artifact review.
2. **Naming & branding** — eliminate legacy `toyol` references from tracked content.
3. **Ignore rules** — `.gitignore`, `.dockerignore`, and `.gitattributes` completeness.
4. **Governance files** — `LICENSE`, `CODE_OF_CONDUCT.md`, `CONTRIBUTING.md`, `SECURITY.md`, `SUPPORT.md`, `MAINTAINERS.md`.
5. **GitHub metadata** — issue/PR templates, labels, `FUNDING.yml`, repository settings checklist.
6. **Branch & release discipline** — `main`/`dev` sync, stale-branch policy, commit conventions.
7. **Developer experience** — pre-commit hooks, `.editorconfig`, release notes template.
8. **Publication readiness** — AI-generated content disclosure, repo-transfer checklist.

---

## Success criteria

- `gitleaks detect --source .` reports zero leaks across full history.
- No tracked file contains `toyol`/`Toyol` branding (directory name exempted).
- `.gitignore` and `.dockerignore` cover all generated/local artifacts.
- Single authoritative `SECURITY.md` (`.github/SECURITY.md`) with root symlink or redirect.
- `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md`, `LICENSE`, and `SUPPORT.md` are present and consistent.
- GitHub issue/PR templates, label taxonomy, and `FUNDING.yml` exist.
- `main` is up to date with `dev` or a documented exception exists.
- Stale worktrees and branches are removed after other readiness initiatives merge.
- Pre-commit hooks cover Go, shell, markdown, workflow, and secret checks.
- All quality gates in `task quality` pass with no new warnings.

---

## Relation to other readiness initiatives

| Initiative | Hand-off boundary |
|---|---|
| `openmuara-readiness-security-audit` | Provides authoritative `.github/SECURITY.md`; this initiative consolidates duplicates and links to it. |
| `openmuara-readiness-docs-completeness-audit` | Owns content quality of `docs/`; this initiative focuses on repo-level metadata and templates. |
| `openmuara-readiness-ci-release-audit` | Owns workflow files; this initiative ensures `.github/settings.yml` or equivalent settings are documented. |
| `openmuara-readiness-repo-hygiene-audit` | **This initiative** owns the final cleanup before public transfer. |
