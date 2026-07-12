> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# 01 — Generator Decision

## Goal

Choose the static-site generator for the OpenMuara documentation website and record the decision.

## Context

OpenMuara has rich Markdown documentation in `docs/`, `runbooks/`, `README.md`, and `CHANGELOG.md`. A proper docs site will improve discoverability and onboarding. The main candidates are:

1. **VitePress** — Node/Vite/Vue, fast, modern.
2. **Docusaurus** — Node/React, mature, versioning.
3. **MkDocs Material** — Python, stable, great default theme.
4. **Hugo + Docsy** — Go-based, extremely fast.
5. **GitHub-rendered Markdown** — do nothing.

## Required Output

- Update `DECISIONS.md` with the chosen generator and the reasons.
- Update `TRACKING.md` prompt 01 status to `✅`.
- Update `HANDOFF.md`.

## Decision Criteria

- Markdown source stays in the root repo.
- Authors can contribute without learning a complex framework.
- Search and navigation work out of the box or with minimal config.
- CI/CD deploy is simple and reliable.
- Fits the project's local-first, low-ceremony culture.

## Quality Gate

- Human review of `DECISIONS.md`.
