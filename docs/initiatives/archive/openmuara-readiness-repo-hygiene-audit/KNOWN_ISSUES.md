> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Repository Hygiene Known Issues

> **Created:** 2026-07-08  
> **Last Updated:** 2026-07-10  
> **Status:** ✅ Delivered on dev

---

## Findings from initial audit

| ID | Finding | Area | Severity | Status | Recommendation ID |
|---|---|---|---|---|---|
| K01 | Legacy `toyol`/`Toyol` references in `docs/contributing-providers.md`, `DECISIONS.md`, `AGENTS.md`, `prompts/18-migration-guide.md`, `.gitignore`, and local `.toyol/` directory | Naming | Medium | Fixed | G1.1 |
| K02 | Duplicate `SECURITY.md`: root version and untracked `.github/SECURITY.md` | Security | Medium | Fixed | G2.1 |
| K03 | `.gitignore` duplicates `bin/` and `coverage.out`; missing `.gstack/` and `.playwright-mcp/` | Ignore rules | Low | Fixed | G2.3 |
| K04 | `.dockerignore` missing `node_modules/`, `coverage.html`, generated UI assets, and local AI directories | Ignore rules | Low | Fixed | G2.4 |
| K05 | No `.editorconfig` | Developer experience | Low | Fixed | G6.1 |
| K06 | `.gitattributes` only has `* text=auto`; no `linguist-generated` or explicit line-ending rules | Developer experience | Low | Fixed | G6.2 |
| K07 | Pre-commit hooks missing `shellcheck`, `actionlint`, and markdown lint | Developer experience | Low | Fixed | G6.3 |
| K08 | No `.github/SUPPORT.md` | Governance | Low | Fixed | G3.1 |
| K09 | No `MAINTAINERS.md` | Governance | Low | Fixed | G3.2 |
| K10 | No `.github/FUNDING.yml` | Governance | Low | Fixed | G3.3 |
| K11 | `CONTRIBUTING.md` lacks commit conventions and bug-register guidance | Governance | Low | Fixed | G3.4 |
| K12 | Issue templates lack provider-request and docs-issue forms; `config.yml` is minimal | GitHub metadata | Low | Fixed | G4.1 |
| K13 | No documented GitHub label taxonomy | GitHub metadata | Low | Fixed | G4.2 |
| K14 | No `.github/settings.yml` or documented repository-settings checklist | GitHub metadata | Low | Fixed | G4.3 |
| K15 | `main` branch last updated 2026-06-12; `dev` is significantly ahead | Branch hygiene | Medium | Fixed | G5.1 |
| K16 | Worktrees exist for merged/in-progress initiatives; `feat/mkp-fawry` is suspended | Branch hygiene | Low | Fixed | G5.2 |
| K17 | No documented commit-message convention | Branch hygiene | Low | Fixed | G5.3 |
| K18 | No `.github/release.yml` for auto-categorized release notes | Branch hygiene | Low | Fixed | G5.4 |
| K19 | `VERSION` says `1.0.0` while `CHANGELOG.md` has large `[Unreleased]` content | Release hygiene | Medium | Fixed | G7.3 |
| K20 | No AI-generated content disclosure | Publication readiness | Low | Fixed | G7.2 |

---

## Positive findings

- `gitleaks detect --source .` reports **zero leaks** across 397 commits (scanned 2026-07-10).
- No tracked binaries, coverage artifacts, or local config files.
- No merge-conflict markers or `TODO`/`FIXME` markers in Go source.
- `CODE_OF_CONDUCT.md`, `CONTRIBUTING.md`, `LICENSE`, issue templates, and PR template already exist and are in good shape.
- `CHANGELOG.md` follows Keep a Changelog format.

---

## How to update this file

As each issue is resolved, change its **Status** to `Fixed` and add the commit hash in the **Notes** column. Do not remove rows; the history is useful for audits.
