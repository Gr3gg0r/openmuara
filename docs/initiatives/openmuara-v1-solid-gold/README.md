> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara v1 Solid Gold

> **Status:** 🟡 In Progress | **Started:** 2026-07-01
> **Scope:** Close remaining v1 hygiene, testing, debuggability, and usability gaps so OpenMuara reaches OSS-grade "solid gold" quality.
> **Owner:** AI Agent (Kimi Code) | **Human Reviewer:** ___________
> **Target Repo:** `<repo-root>/`
> **Product Branch:** `feat/v1-solid-gold`
>
> **Why:** v1 feature work is complete. This initiative is purely about making the project more reliable, easier to debug, and nicer to use.

---

## Initiative Structure

```
docs/initiatives/openmuara-v1-solid-gold/
├── README.md              # This file
├── TRACKING.md            # Central execution tracker
├── HANDOFF.md             # Session continuity
├── DECISIONS.md           # Decision log
├── RISKS.md               # Risk register
│
└── prompts/               # Numbered, self-contained execution prompts
    ├── _template.md
    ├── 01-tooling-hygiene.md
    ├── 02-coverage-backfill.md
    ├── 03-debuggability.md
    ├── 04-dashboard-usability.md
    └── 05-best-practices-and-tooling.md
```

Planning docs live in `docs/initiatives/openmuara-v1-solid-gold/` in the root repo.
Product code commits to the `feat/v1-solid-gold` branch. Do not commit directly to
`main`.

---

## Goals

1. **Tooling hygiene** — Fix all current local quality-gate failures (config drift,
   shellcheck warnings) and make CI run the same matrix as `task quality`.
2. **Coverage backfill** — Bring every package to ≥80% coverage, especially the
   weakest ones (`internal/ui`, `internal/fawry/v2`, `internal/cli`, etc.).
3. **Debuggability** — Add trace-ID propagation, CLI inspect commands, and optional
   pprof endpoints so failed webhooks and provider calls are easy to trace.
4. **Dashboard usability** — Surface failed webhooks, add copy-paste curl buttons,
   and improve mobile layout.
5. **Best practices / tooling** — Strengthen linters, add pre-commit hooks, dependency
   automation, and reproducible build flags.

---

## Conventions (MANDATORY)

### 1. Read `AGENTS.md` first
Root `AGENTS.md` governs branch rules, quality gates, autonomy boundaries, and
code style.

### 2. Backward compatibility
All changes are additive or hygiene fixes. No breaking API or config changes.

### 3. P0 integration changes need explicit approval
Prompts that touch webhook dispatch, provider signature verification, or core
billing flows (P03) require user sign-off per `AGENTS.md`.

### 4. Quality gates
Every prompt must pass:

- `go build ./...`
- `go test ./...`
- `go vet ./...`
- `golangci-lint run`
- `./scripts/smoke-test.sh`
- `task quality` (for P01)

### 5. Definition of done
Beyond the quality gates, a prompt is done only when:

- The change is tested or justified as untestable.
- `HANDOFF.md` is updated with what was built.
- `TRACKING.md` marks the prompt `✅` with the commit hash.
- User-facing changes are noted for the next release notes.

---

## Out of Scope

- New providers or payment methods.
- Changes to the provider plugin schema contract.
- Authentication or authorization in the dashboard.
- v2 features (App Store / Play Store receipts, RevenueCat).

---

## Metrics

| Metric | Current | Target | How measured |
|--------|---------|--------|--------------|
| `task quality` pass rate | Fails on audit/scripts | 100% pass | `task quality` |
| Lowest package coverage | 21.4% (`internal/ui`) | ≥80% every package | `go test -cover ./...` |
| Shell script hygiene | 2 shellcheck warnings | 0 warnings | `scripts/check-scripts.sh` |
| CI/local gate parity | CI skips vuln/forbidden/scripts/sizes/audit | Identical | `.github/workflows/ci.yml` vs `task quality` |
| Dashboard first-failure visibility | Manual scan | Alert within 2s of load | Manual / Playwright test |

## Success Criteria

- `task quality` passes locally and in CI with zero warnings.
- `go test -cover ./...` shows every package ≥80%.
- A failed webhook can be traced end-to-end with a single ID.
- The dashboard surfaces failures without opening server logs.
- New contributors can run one command (`task quality`) and get the same gate
  results as CI.
