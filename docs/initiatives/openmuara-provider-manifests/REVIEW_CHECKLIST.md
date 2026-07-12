> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Provider Manifests — Pre-PR Review Checklist

For the human reviewer before merging to `dev`.

---

## Architecture

- [ ] The manifest-first loader is the only path for non-default provider discovery.
- [ ] `runtime.type: simple` providers run without any Go registration.
- [ ] `runtime.type: go` providers run through a factory registry keyed by name.
- [ ] No built-in provider auto-registers its instance in `init()`.
- [ ] Removing a provider's manifest cleanly removes it from discovery.
- [ ] `default` provider behavior is preserved or explicitly changed with sign-off.

## Code Quality

- [ ] `go build ./...` passes.
- [ ] `go test ./...` passes.
- [ ] `go vet ./...` passes.
- [ ] `golangci-lint run ./...` passes with zero warnings.
- [ ] `go test -race ./...` passes.
- [ ] Test coverage does not drop on changed modules.

## Provider Correctness

- [ ] Conformance tests exist for each migrated provider.
- [ ] `examples/checkout-store` Fawry flow works.
- [ ] `examples/checkout-store` Stripe flow works.
- [ ] `muara provider validate plugins/*/gateway.yml` passes for all manifests.

## Documentation

- [ ] `docs/provider-contract.md` is accurate post-change.
- [ ] `docs/contributing-providers.md` explains simple vs go vs bridge vs wasm.
- [ ] `docs/migration/provider-manifests.md` exists if breaking changes affect users.
- [ ] `CHANGELOG.md` has a release-notes snippet.

## Trackers

- [ ] `TRACKING.md` is up to date.
- [ ] `HANDOFF.md` is up to date.
- [ ] `DECISIONS.md` is up to date.
- [ ] `RISKS.md` is up to date.
- [ ] `KNOWN_ISSUES.md` reflects resolved and remaining issues.
- [ ] `RECOMMENDATIONS.md` is reviewed and approved or explicitly overridden.
- [ ] `appendices/e-test-scenarios.md` is reviewed.
- [ ] `appendices/f-architecture-diagram.md` matches the implementation.

## Security & Safety

- [ ] No secrets committed.
- [ ] Provider factories do not bypass config validation.
- [ ] No new global mutable state introduced.
- [ ] Error messages do not leak internal paths or secrets.

## Scope

- [ ] No `bridge` or `wasm` runtime implementation snuck in.
- [ ] No unrelated refactors.
- [ ] Each commit is one logical change.
