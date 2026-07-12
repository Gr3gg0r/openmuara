---
id: testing
title: Testing Guide — OpenMuara
---

> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This project is subordinate to it.**

# Testing Guide — OpenMuara

This guide explains how to write, run, and maintain tests for OpenMuara.

---

## Test pyramid

OpenMuara uses four layers, ordered by speed and scope:

| Layer | Speed | Scope | Examples |
|-------|-------|-------|----------|
| **Unit** | Fast | Single function or package | `engine.Transition`, `fawry.Sign`, CLI command parsing |
| **Integration** | Medium | Multiple packages with real (in-memory) dependencies | Router + provider + store + webhook dispatcher |
| **Contract** | Medium | Provider request/response shapes vs golden files | Fawry charge response, Stripe checkout session |
| **E2E / Smoke** | Slow | Whole binary against real HTTP | `scripts/smoke-test.sh` |

Write most tests at the unit and integration layers. Keep E2E tests small and deterministic.

---

## Naming conventions

- Test files: `*_test.go` in the same package as the code under test.
- Test functions: `Test<Name>_<Condition>` or `Test<Name>`.
  - Good: `TestChargeHandler_InvalidSignature`, `TestTransition_PaidToRefunded`.
  - Avoid: `Test1`, `TestStuff`.
- Integration tests: use a build tag or name suffix `_integration_test.go` if they need a real server/port.
- Fuzz tests: `Fuzz<Name>` in `*_fuzz_test.go`.
- Property-based tests: `Test<Name>_Property` or use a dedicated `*_property_test.go`.
- Benchmarks: `Benchmark<Name>`.

---

## Table-driven tests

Use table-driven tests for multiple inputs against the same function:

```go
func TestTransition(t *testing.T) {
    tests := []struct {
        name    string
        from    engine.TransactionStatus
        to      engine.TransactionStatus
        wantErr bool
    }{
        {"paid from new", engine.TransactionStatusNew, engine.TransactionStatusPaid, false},
        {"refunded from paid", engine.TransactionStatusPaid, engine.TransactionStatusRefunded, false},
        {"refunded from new", engine.TransactionStatusNew, engine.TransactionStatusRefunded, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tx := engine.Transaction{Status: tt.from}
            err := engine.Transition(&tx, tt.to)
            if (err != nil) != tt.wantErr {
                t.Fatalf("Transition(%q -> %q) error = %v, wantErr %v", tt.from, tt.to, err, tt.wantErr)
            }
        })
    }
}
```

---

## Testable design

- **Do not rely on package-level mutable state.** Pass registries, stores, and loggers as dependencies.
- **Use interfaces for external dependencies.** The transaction store, webhook sender, and provider registry all have interfaces.
- **Avoid hardcoded ports.** Use `httptest.Server` or port `0` for real servers.
- **Avoid writing to the real workspace.** Use `t.TempDir()` and `internal/testutil.TempWorkspace(t)`.
- **Do not call `slog.SetDefault` in reusable code.** Accept a `*slog.Logger` or use context.

---

## Fakes, mocks, and real dependencies

| Approach | When to use | Example |
|----------|-------------|---------|
| **Real in-memory implementation** | Fast and deterministic; preferred | `engine.NewMemoryStore()`, `audit.NewMemoryStore()` |
| **Fake** | Need to observe behavior without real side effects | `testutil.FakeDispatcher` records webhook attempts |
| **Mock** | Need strict assertions on calls (use sparingly) | Mock sender in `webhook` package tests |
| **httptest.Server** | Real HTTP round-trip without a port | Provider contract tests, SDK tests |

Avoid heavy mocking frameworks. Prefer fakes and real in-memory implementations.

---

## Golden files and testdata

Place fixtures in a `testdata/` directory next to the test file.

```
internal/fawry/
  contract_test.go
  testdata/
    charge_request.json
    charge_response.json
```

Use `testutil.GoldenFile(t, path)` to load golden files. If you add an `-update` flag, document it in the test:

```go
var update = flag.Bool("update", false, "update golden files")
```

Keep golden files readable and minimal. Do not include dynamic fields like timestamps or trace IDs unless you normalize them first.

---

## Coverage policy

- The coverage threshold is **80%**.
- Do not write tests just to hit the threshold. Every test should assert meaningful behavior.
- Do not test trivial getters/setters unless they contain logic.
- Run `task coverage` before committing.
- Aim for 100% coverage of error paths in critical packages (`engine`, `webhook`, `server`).

---

## Running tests

```bash
# Unit and integration tests
go test ./...

# With race detector
go test -race ./...

# Shuffle execution order (catches hidden dependencies)
go test -shuffle=on ./...

# Coverage
task coverage

# Smoke test
./scripts/smoke-test.sh

# Full local gate
task check
```

---

## Writing a new test

1. Decide the layer (unit/integration/contract/E2E).
2. Create `*_test.go` in the same package.
3. Use `internal/testutil` helpers for workspace, stores, and servers.
4. Use table-driven tests for multiple cases.
5. Assert behavior, not implementation details.
6. Run `task check` and fix any failures.

---

## Test anti-patterns

- ❌ Sleeping to wait for async behavior — use channels or `require.Eventually`.
- ❌ Using the real `.muara/` workspace — use `t.TempDir()`.
- ❌ Hardcoding ports like `9000` — use port `0` or `httptest.Server`.
- ❌ Testing private functions directly — test public behavior.
- ❌ Heavy mocking — prefer fakes and real in-memory stores.
- ❌ Ignoring errors from `t.Cleanup` or `defer` — handle resource cleanup.

---

## CI test jobs

CI splits testing into parallel jobs:

1. **lint** — `gofmt`, `go vet`, `golangci-lint`.
2. **unit** — `go test -race -coverprofile=coverage.out ./...` with 80% threshold.
3. **integration** — integration tests with real HTTP and random ports.
4. **smoke** — `./scripts/smoke-test.sh`.

See `.github/workflows/ci.yml` for the exact commands.

---

## Related resources

- [Testing Gold Standard initiative](https://github.com/openmuara/openmuara/tree/main/docs/initiatives/openmuara-testing-gold-standard)
- [Quality Gates runbook](quality-gates)
