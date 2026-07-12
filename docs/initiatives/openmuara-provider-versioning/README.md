> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Provider API Versioning

> **Status:** ⬜ Not Started | **Started:** —
> **Scope:** Introduce a convention for emulating multiple API versions of the same payment provider.
> **Owner:** AI Agent (Kimi Code) | **Human Reviewer:** ___________
> **Target Repo:** `<repo-root>/`
> **Product Branch:** `feat/provider-versioning`

---

## Initiative Structure

```
docs/initiatives/openmuara-provider-versioning/
├── README.md              # This file
├── TRACKING.md            # Central execution tracker
└── prompts/
    └── 01-versioned-provider-layout.md
```

Planning docs live in `docs/initiatives/openmuara-provider-versioning/` in the root repo. Product code commits to the `dev` branch. Do not commit directly to `main`.

---

## Why version providers?

Real payment providers evolve their APIs over time:

- **Stripe** uses a `Stripe-Version` header and dated API versions.
- **Billplz** has v3 today and may add v4.
- **Fawry** already has both legacy and V2 notification formats (the current codebase references `webhook.FawryV2Payload`).
- **iPay88** has classic and newer variants.

Right now each OpenMuara provider is a flat package (`internal/fawry/`, `internal/billplz/`, etc.). That works for a single version but makes it hard to keep multiple provider API versions faithful and side-by-side. A clear versioning convention lets us:

1. Emulate the exact version a user’s app targets.
2. Keep new versions from breaking existing tests and smoke flows.
3. Make the provider layout predictable for contributors.

---

## Goals

1. Define and document a provider-versioning package layout.
2. Add a `version` config field to provider configuration with a safe default (`v1`).
3. Build a thin top-level provider dispatcher that selects the requested version.
4. Migrate **Fawry** as the reference implementation because it already contains V2 payload concepts.
5. Update the Fawry plugin manifest (`plugins/fawry/gateway.yml`) to declare supported versions.
6. Add smoke-test coverage for Fawry V2.
7. Pass all quality gates.

---

## Target Layout

```
internal/<provider>/
├── provider.go            # version dispatcher; implements provider.Provider
├── v1/
│   └── provider.go        # v1 concrete implementation
│   └── charge.go
│   └── signature.go
│   └── webhook.go
├── v2/
│   └── provider.go        # v2 concrete implementation
│   └── charge.go
│   └── signature.go
│   └── webhook.go
```

- The top-level `Provider.Name()` returns the provider name (e.g., `fawry`).
- `Provider.Init()` reads `cfg["version"]` (default `"v1"`), validates it, and instantiates the matching version package.
- `Provider.Routes()` exposes versioned routes. Default routing convention:
  - Unversioned legacy routes continue to work and map to the configured default version.
  - Explicit versioned routes are available under `/<provider>/v1/...` and `/<provider>/v2/...`.
- Version packages return `provider.Route` slices and use the same `engine.TransactionStore`, `webhook.Dispatcher`, and config injection interfaces (`SetStore`, `SetDispatcher`, etc.).

---

## Conventions (MANDATORY)

### 1. Read `AGENTS.md` first
Root `AGENTS.md` governs branch rules, quality gates, autonomy boundaries, and code style.

### 2. One reference provider only
Do **not** rewrite every provider in this initiative. Migrate only Fawry. Other providers can be versioned later using the same convention.

### 3. Backward compatibility
Existing `/fawry/charge` and `/fawry/webhook` routes must continue to behave like v1 unless a user explicitly opts into v2 via config or a versioned path.

### 4. Quality gates
Every prompt must pass:

- `go build ./...`
- `go test ./...`
- `go vet ./...`
- `golangci-lint run`
- `./scripts/smoke-test.sh`

---

## Out of Scope

- Versioning the provider plugin schema itself (`gateway.yml` `schema_version`).
- Migrating Stripe, Billplz, ToyyibPay, iPay88, or SenangPay to versioned layouts.
- Header-based API version negotiation beyond a simple `version` config field.
- Breaking changes to the existing `provider.Provider` interface.
