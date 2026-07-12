> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Coverage Audit CI Integration

> **Created:** 2026-07-09
> **Last Updated:** 2026-07-09
> **Status:** ✅ Complete — changes applied

---

This document contains the exact files and CI changes applied by the coverage audit initiative.

## 1. New files created

### `scripts/check-coverage-per-package.sh`

```bash
#!/bin/bash
set -euo pipefail

# Minimum coverage per package. Format: "import-path threshold"
# Overall total is still enforced by scripts/check-coverage.sh.
# Floors are calibrated to current architecture; see coverage-exemptions.yml.
MIN_FLOORS=(
  "github.com/openmuara/openmuara/internal/audit 80"
  "github.com/openmuara/openmuara/internal/plugin 80"
  "github.com/openmuara/openmuara/internal/provider/conform 79"
  "github.com/openmuara/openmuara/internal/version 70"
  "github.com/openmuara/openmuara/internal/ui 70"
  "github.com/openmuara/openmuara/internal/provider/simple 45"
)

fail=0
for entry in "${MIN_FLOORS[@]}"; do
  pkg="${entry% *}"
  floor="${entry#* }"
  cov=$(go test -cover "$pkg" 2>&1 |
    sed -E -n 's/^ok[[:space:]]+[^[:space:]]+[[:space:]]+[^[:space:]]+[[:space:]]+coverage: ([0-9.]+)%.*/\1/p')
  if [[ -z "$cov" ]]; then
    echo "WARN: no coverage data for $pkg"
    continue
  fi
  if awk "BEGIN { exit ($cov >= $floor) ? 0 : 1 }"; then
    echo "PASS: $pkg $cov% >= $floor%"
  else
    echo "FAIL: $pkg $cov% < $floor%" >&2
    fail=1
  fi
done

exit "$fail"
```

### `coverage-exemptions.yml`

See the committed file at repo root. It records all Go and dashboard coverage exemptions with rationale, owner, and review date.

## 2. Dashboard changes

### Coverage provider installed

`@vitest/coverage-v8@^2.1.9` added to `web/dashboard/package.json`.

### `web/dashboard/vitest.config.ts`

Configured with v8 coverage provider, reporters `text`/`json`/`json-summary`/`html`, thresholds 60/55/55/55, and exclusions for entry points and type files.

### `web/dashboard/package.json` scripts

Added:

```json
"test:coverage": "vitest run --coverage"
```

### `.gitignore`

Added `web/dashboard/coverage/`, `coverage.out`, and `coverage.html`.

## 3. Go CI updates

In `.github/workflows/ci.yml`, the `unit` job now:

- Enforces 81% overall coverage.
- Enforces per-package floors via `scripts/check-coverage-per-package.sh`.
- Generates `coverage.html`.
- Uploads `coverage.out` and `coverage.html` as the `go-coverage` artifact.

## 4. Dashboard coverage CI job

Added `dashboard-coverage` job to `.github/workflows/ci.yml`:

- Checks out the repo.
- Sets up Node 20.
- Installs UI dependencies.
- Runs `npm run test:coverage`.
- Uploads `web/dashboard/coverage/` as the `dashboard-coverage` artifact.

## 5. Coverage comment workflow update

Updated `.github/workflows/coverage-comment.yml` to:

- Run Go tests with coverage.
- Run dashboard coverage and parse `coverage-summary.json`.
- Post a single PR comment with Go total, lowest package, dashboard metrics table, and regression gate status.

## 6. Local acceptance commands

```bash
# Go baseline + per-package floors
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
scripts/check-coverage.sh 81
scripts/check-coverage-per-package.sh

# Dashboard coverage
cd web/dashboard
npm run test:coverage

# Regression (from a PR branch)
./scripts/check-coverage-regression.sh origin/main 10 1.0
```

## 7. Rollback / exception handling

- If a package cannot reach its floor, add it to `KNOWN_ISSUES.md` and `coverage-exemptions.yml` with rationale and review date.
- If CI becomes flaky due to strict dashboard thresholds, temporarily lower them in `vitest.config.ts` and record the decision in `DECISIONS.md`.
