> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Repository Hygiene Decisions

> **Created:** 2026-07-10  
> **Status:** Draft — decisions recorded for implementation

---

## D1. Root directory name stays `toyol` for now

- **Context:** The repository directory on local disk is still named `toyol` from the previous project name.
- **Decision:** Do **not** rename the root directory in this initiative. All tracked content uses `OpenMuara` / `openmuara` / `muara` branding. The local directory name will be updated when the repository is transferred to the new GitHub organization.
- **Rationale:** Avoids breaking local tooling, worktrees, and CI paths mid-initiative.

## D2. Canonical security policy lives in `.github/SECURITY.md`

- **Context:** There are two `SECURITY.md` files: root and `.github/`.
- **Decision:** Keep `.github/SECURITY.md` as the canonical file (it is more complete). Replace root `SECURITY.md` with a minimal redirect.
- **Rationale:** GitHub surfaces `.github/SECURITY.md` in the security tab and issue picker; the root file is only for local readers.

## D3. Next release version is `1.1.0`

- **Context:** `VERSION` currently says `1.0.0`, but `CHANGELOG.md` has extensive `[Unreleased]` changes including provider manifest discovery and UI overhauls.
- **Decision:** Bump `VERSION` to `1.1.0` and retitle the `[Unreleased]` section to `## [1.1.0] - YYYY-MM-DD` at release time.
- **Rationale:** The manifest-first provider discovery is a significant change but the public API surface for end users remains compatible; reserving `2.0.0` for future breaking CLI or config changes.

## D4. AI-assisted development will be disclosed

- **Context:** The project has been developed with significant AI assistance.
- **Decision:** Add a concise disclosure in `README.md` and `CONTRIBUTING.md` stating that code and documentation are AI-assisted and human-reviewed.
- **Rationale:** Transparency builds trust with contributors and users; aligns with emerging OSS norms.

## D5. No `FUNDING.yml` until a funding model exists

- **Context:** GitHub supports `FUNDING.yml` for sponsorships.
- **Decision:** Add a placeholder `.github/FUNDING.yml` with explanatory comments stating it is intentionally empty and will be populated when a funding model is established.
- **Rationale:** Avoids an empty or misleading sponsorship page while keeping the file present for future use.

## D6. `main` will be fast-forwarded to `dev` before transfer

- **Context:** `main` is several weeks behind `dev`.
- **Decision:** Fast-forward `main` to the current `dev` head immediately before the GitHub organization transfer.
- **Rationale:** A public `main` branch should reflect the current state of the project. After transfer, releases will continue to drive `main` updates.

## D7. Stale worktrees will be removed after dependent initiatives merge

- **Context:** Worktrees exist for `feat/checkout-store-e2e-fixes`, `feat/readiness-ci-release-audit`, and `feat/readiness-docs-completeness`.
- **Decision:** Remove each worktree once its initiative is merged to `dev`. Delete or archive the suspended `feat/mkp-fawry` branch.
- **Rationale:** Keeps the local workspace clean and avoids confusion.

## D8. Pre-commit hooks stay fast; heavy checks remain in CI

- **Context:** Pre-commit hooks could become slow if they run the full test suite or heavy linters.
- **Decision:** Pre-commit runs fast checks (gofmt, go vet, shellcheck, actionlint, gitleaks). The full `go test ./...`, `golangci-lint run`, and `task quality` remain CI gates.
- **Rationale:** Developers will bypass slow hooks; CI is the right place for comprehensive checks.

## D9. Repository settings will be documented, not automated

- **Context:** `.github/settings.yml` via Probot Settings could automate repo settings but adds an external dependency.
- **Decision:** Document required settings in `APPENDIX.md` and apply them manually before/after transfer.
- **Rationale:** Manual application avoids Probot permissions and keeps the repo self-contained.
