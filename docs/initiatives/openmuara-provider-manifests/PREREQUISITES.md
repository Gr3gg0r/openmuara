> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Provider Manifests — Prerequisites

Use this checklist before starting implementation.

---

## Definition of Ready

This initiative is ready to implement when:

- [ ] The architecture in `DECISIONS.md` is pinned.
- [ ] A human reviewer has been assigned in `README.md` and `HANDOFF.md`.
- [ ] The human reviewer has acknowledged the scope and branch (`dev`).
- [ ] Recommended resolutions in `RECOMMENDATIONS.md` are approved or explicitly overridden.
- [ ] No unresolved P0 integration sign-off items remain.
- [ ] The current `dev` branch passes all gates.
- [ ] `TRACKING.md` is up to date.
- [ ] `appendices/e-test-scenarios.md` and `appendices/f-architecture-diagram.md` have been reviewed.

---

## Repository

- [ ] You are on branch `dev`.
- [ ] `git status` shows only expected changes.
- [ ] `go build ./...` passes on current `dev`.
- [ ] `go test ./...` passes on current `dev`.
- [ ] `go vet ./...` passes on current `dev`.
- [ ] `golangci-lint run ./...` passes on current `dev`.

## Understanding

- [ ] Read `docs/initiatives/openmuara-provider-manifests/README.md`.
- [ ] Read `docs/initiatives/openmuara-provider-manifests/DECISIONS.md`.
- [ ] Read `docs/initiatives/openmuara-provider-manifests/RISKS.md`.
- [ ] Understand the four provider runtime paths:
  - `runtime.type: simple` — pure YAML manifest
  - `runtime.type: go` — YAML manifest + Go factory
  - `providers.<name>.type: bridge` — proprietary/private (future)
  - `.muara/plugins/<name>.wasm` — sandboxed plugin (future)

## Tooling

- [ ] Go 1.22+ installed.
- [ ] `golangci-lint` installed.
- [ ] SQLite development headers available (usually included on macOS).
- [ ] Node.js available for `examples/checkout-store` and dashboard tests.

## Safety

- [ ] Do not commit `.muara/` or real `config.yml` files.
- [ ] Do not modify `main` branch.
- [ ] Keep product-code commits separate from planning-doc commits.
- [ ] P0 integration changes require user sign-off before implementation.
- [ ] Review `README.md` Assumptions & Constraints before starting.

## Time-box

- Planning: complete before product code.
- Implementation: 1–2 agent sessions.
- QA and docs: half a session.

If implementation exceeds the time-box, escalate in `RISKS.md` and `HANDOFF.md`.

---

## Pre-Flight Commands

Run these before starting P01:

```bash
git checkout dev
git pull origin dev
go build ./...
go test ./...
go vet ./...
golangci-lint run ./...
```

If any fail, fix or escalate before proceeding.
