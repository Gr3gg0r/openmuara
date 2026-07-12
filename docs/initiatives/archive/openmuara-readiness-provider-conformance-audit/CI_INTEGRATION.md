> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Provider Conformance Audit CI Integration

> **Created:** 2026-07-09
> **Last Updated:** 2026-07-09
> **Status:** ✅ Complete — conformance regression gate added to `.github/workflows/ci.yml`.

---

This document describes the CI changes implemented by the provider conformance audit initiative.

## 1. Conformance test command

Existing command:

```bash
go test ./internal/provider/conform/...
```

After extending `conform` to cover behavior snapshots, this command will also validate request/response golden files.

## 2. Golden-file update workflow

Add a documented environment variable for regenerating golden files:

```bash
UPDATE_GOLDEN=1 go test ./internal/provider/conform/...
```

PRs that update golden files must explain the provider contract change in the PR description.

## 3. CI gate in `.github/workflows/ci.yml`

Applied in the `unit` job:

```yaml
      - name: Provider conformance regression
        run: go test -race ./internal/provider/conform/...
```

This ensures any route or behavior change that drifts from golden files fails CI. The full `go test -race ./...` step already covers provider-specific conformance tests; this step makes the conformance gate explicit.

## 4. Optional: dedicated conformance job

For larger provider matrices, add a separate job:

```yaml
  provider-conformance:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@34e114876b0b11c390a56381ad16ebd13914f8d5 # v4

      - name: Set up Go
        uses: actions/setup-go@40f1582b2485089dde7abd97c1529aa768e1baff # v5
        with:
          go-version: '1.26'

      - name: Download modules
        run: go mod download

      - name: Provider conformance tests
        run: go test ./internal/provider/conform/...
```

## 5. Golden-file protection

Treat `internal/provider/conform/testdata/golden/*.json` as generated files:

- Review them in PR diffs.
- Require PR description to explain any `-update` regeneration.
- Do not accept golden-file changes without corresponding code or doc updates.

## 6. Local acceptance commands

```bash
# Run conformance tests
go test ./internal/provider/conform/...

# Update golden files after intentional contract changes
UPDATE_GOLDEN=1 go test ./internal/provider/conform/...

# Run provider-specific tests
go test ./internal/fawry/... ./internal/stripe/...
```

## 7. Rollback / exception handling

- If a provider contract cannot be fully emulated, add the deviation to `KNOWN_ISSUES.md` and update the golden file.
- If CI becomes flaky due to non-deterministic snapshots, isolate the unstable provider into its own job or mark the test as helper-only.
