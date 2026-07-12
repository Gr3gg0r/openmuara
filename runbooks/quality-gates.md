---
id: quality-gates
title: Local Quality Gates
---

# Local Quality Gates

This runbook describes the local quality gates for the OpenMuara Go codebase.

---

## Running all gates

```bash
task check
```

`task check` runs, in order:

1. `gofmt -l .` — ensure all Go files are formatted.
2. `go vet ./...` — run static analysis.
3. `golangci-lint run` — run configured linters.
4. `go test -race -coverprofile=coverage.out ./...` — run tests with race detector.
5. `scripts/check-coverage.sh 80` — enforce 80% minimum coverage.

If any gate fails, fix the issue before committing.

Run the full quality matrix (including smoke, vulnerability scan, forbidden-pattern check, shell-script check, size advisory, and tracker audit) with:

```bash
task quality
```

---

## CI jobs

The GitHub Actions workflow (`.github/workflows/ci.yml`) runs the same gates in parallel jobs for fast feedback:

- `lint` — `gofmt`, `go vet`, and `golangci-lint`.
- `unit` — `go test -race -coverprofile=coverage.out ./...` with an 80% coverage gate.
- `smoke` — `./scripts/smoke-test.sh` end-to-end test using a random port and isolated workspace.
- `vuln` — `govulncheck ./...` vulnerability scan.
- `gosec` — Go security linter via `securego/gosec`.
- `secrets` — `gitleaks-action@v2` secret scan.
- `quality` — Full local quality matrix via `task quality`.

Additional dedicated workflows:

- `.github/workflows/coverage-comment.yml` — posts a coverage summary and runs the per-module coverage regression gate. Non-blocking during the phased rollout.
- `.github/workflows/visual-baseline.yml` — Playwright visual diff for dashboard changes, filtered to `web/dashboard/**` and `internal/ui/**`. Runs on light and dark themes.
- `.github/workflows/mutation.yml` — Gremlins mutation testing for targeted Go packages, filtered to `internal/**/*.go`, `go.mod`, `go.sum`, and `scripts/mutation-test.sh`. Non-blocking during the phased rollout.

---

## Individual gates

### Build

```bash
go build ./...
```

### Tests with coverage

```bash
task test
```

Generates `coverage.out` and prints total coverage.

### Coverage gate

```bash
task coverage
```

Exits with code 1 if total coverage is below 80%. Override the threshold:

```bash
scripts/check-coverage.sh 85
```

### Lint

```bash
task lint
```

### Security scan

```bash
task security
```

Runs `gosec ./...` if installed. If `gosec` is missing, the task prints install instructions and exits 0 so it does not block development.

Install gosec:

```bash
go install github.com/securego/gosec/v2/cmd/gosec@latest
```

### Secret scan

```bash
./scripts/check-gitleaks.sh
```

Runs `gitleaks detect --source .` if installed. CI uses `gitleaks/gitleaks-action@v2` for full-history scanning.

Install gitleaks:

```bash
brew install gitleaks
# or download a release from https://github.com/gitleaks/gitleaks/releases
```

### Vulnerability scan

```bash
task vuln
```

Runs `govulncheck ./...` if installed. If `govulncheck` is missing, the task prints install instructions and exits 0 so it does not block development.

Install govulncheck:

```bash
go install golang.org/x/vuln/cmd/govulncheck@latest
```

### Smoke test

```bash
task smoke
```

Runs the end-to-end smoke test against a locally started server.

### Forbidden-pattern check

```bash
task forbidden
```

Ensures production Go code does not contain `fmt.Println` or `os.Exit` outside `cmd/`.

### Shell-script check

```bash
task scripts
```

Runs `shellcheck` on `./scripts/*.sh` if `shellcheck` is installed. Skips cleanly otherwise.

### Size advisory

```bash
task sizes
```

Prints advisory warnings for Go files, functions, or lines that exceed the recommended limits (250 lines/file, 80 lines/function, 120 characters/line). This task never fails; it is meant for ongoing refactoring.

### Environment check

```bash
./bin/muara doctor
```

Reports whether required tools (`go`, `golangci-lint`, `task`) and the optional `govulncheck` are on PATH.

### Frontend build

```bash
cd web/dashboard && npm run build
```

Produces the dashboard bundle that is embedded by `internal/ui/`.

### Frontend tests

```bash
cd web/dashboard && npm run test:ci
```

Runs unit tests in CI mode.

### Accessibility contrast check

```bash
cd web/dashboard && node scripts/a11y-contrast-check.js
```

Fails if any dashboard color combination violates WCAG contrast requirements.

### Bundle size check

```bash
node web/dashboard/scripts/check-bundle-size.js
```

Fails if the dashboard bundle exceeds the configured budget.

### Visual baseline

```bash
cd web/dashboard && npm run test:visual-baseline
```

Compares dashboard screenshots against baselines in `web/dashboard/e2e/baselines/`. Baselines are captured for both **light** and **dark** themes (`-light.png` and `-dark.png` suffixes). Dynamic elements (e.g., live refresh timestamps) are hidden with the shared `[data-visual-mask]` attribute instead of per-test CSS.

Update baselines intentionally with:

```bash
cd web/dashboard && npm run test:visual-baseline -- --update-snapshots
```

### Mutation testing

```bash
./scripts/mutation-test.sh 70
```

Runs Gremlins on `internal/webhook` and `internal/engine`. Install Gremlins first:

```bash
go install github.com/go-gremlins/gremlins/cmd/gremlins@latest
```

The initial threshold is 70%; raise it once the baseline is stable. `internal/fawry` is currently excluded because mutations cause its HTTP-handler tests to time out; re-evaluate after adding faster pure-function tests. The CI mutation job is non-blocking during the phased rollout so scores are reported without blocking merges while stability is measured.

### Coverage regression

```bash
./scripts/check-coverage-regression.sh origin/dev 10 1.0
```

Compares per-package coverage between the current checkout and the base ref. Fails if any changed Go package drops coverage by more than the tolerance (1.0%). The gate is non-blocking in CI during the phased rollout.

### Release / reproducible builds

Release binaries are built with `-trimpath` so full build paths are stripped, making builds more reproducible across machines:

```bash
task release:build
```

To verify a binary was built with `-trimpath`:

```bash
go version -m dist/muara-linux-amd64 | grep -q '^\s*path\s' || echo 'trimpath not detected'
```

The `path` line should be absent or only the module path, not a local filesystem path.

---

## Pre-commit hooks

Install hooks so fast gates run on every commit:

```bash
pre-commit install
```

Run manually on all files:

```bash
pre-commit run --all-files
```

The configured hooks are:

- `gofmt`
- `go vet ./...`
- `go test ./...`
- `golangci-lint run`
- `scripts/check-forbidden.sh`

Run `task check` before pushing for the full gate matrix including race detection and coverage.

---

## Troubleshooting

### Coverage gate fails

1. Run `task test` to see per-package coverage.
2. Add focused unit tests for uncovered pure functions and handlers.
3. Avoid adding tests only to game the metric.

### golangci-lint reports issues

1. Run `golangci-lint run --fix` to auto-fix formatting issues.
2. Address remaining issues manually.

### govulncheck is not found

Install it with:

```bash
go install golang.org/x/vuln/cmd/govulncheck@latest
```

If you choose not to install it, `task vuln` and `task check` will still pass.
