> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# P01 — Versioned Provider Layout and Fawry Reference Migration

## Goal

Introduce a provider-versioning convention that supports both single-version and multi-version providers, then migrate **Fawry** as the reference implementation.

The convention must:

- Keep single-version providers simple (flat package stays valid).
- Allow multi-version providers to live in `internal/<provider>/v1/`, `internal/<provider>/v2/`, etc.
- Preserve backward compatibility: unversioned routes continue to work.
- Stay lightweight so contributors are not forced to learn a heavy plugin/version framework.

## Branch

Work on `feat/provider-versioning`. Do not commit to `dev` or `main`.

## Background

OpenMuara providers currently live in flat packages:

```
internal/fawry/
├── provider.go
├── charge.go
├── webhook.go
├── signature.go
└── ...
```

Fawry already references `webhook.FawryV2Payload`, which signals that a V2 notification format exists. The current code handles only one shape. We want a layout that can hold both V1 and V2 cleanly.

## Target Layout

### Multi-version provider (Fawry after this prompt)

```
internal/fawry/
├── provider.go              # thin dispatcher; reads version config
├── v1/
│   ├── provider.go          # concrete v1 provider
│   ├── charge.go            # /fawry/v1/charge (and legacy /fawry/charge alias)
│   ├── signature.go         # V1 signature rules
│   └── webhook.go           # /fawry/v1/webhook with V1 payload
└── v2/
    ├── provider.go          # concrete v2 provider
    ├── charge.go            # /fawry/v2/charge
    ├── signature.go         # V2 signature rules
    └── webhook.go           # /fawry/v2/webhook with FawryV2Payload
```

### Single-version provider (unchanged)

```
internal/billplz/
├── provider.go
├── bill.go
└── ...
```

A single-version provider does **not** need a `v1/` sub-package. It only needs the `version` config field to default to `"v1"` if we ever decide to version it later.

## Required Changes

### 1. Version dispatcher (`internal/fawry/provider.go`)

- Read `cfg["version"]`; default `"v1"` if missing or empty.
- Validate: only `"v1"` and `"v2"` are accepted.
- The top-level `Provider` stores the selected version provider (interface below).
- `Routes()` returns the union of:
  - Legacy routes `/fawry/charge` and `/fawry/webhook` mapped to the configured default version.
  - Explicit version routes `/fawry/v1/charge`, `/fawry/v1/webhook`, `/fawry/v2/charge`, `/fawry/v2/webhook`.
- `ChargeHandler()`, `WebhookHandler()`, `PayloadBuilder()`, `EscapeHandler()` delegate to the selected version provider.
- Keep `SetStore` and `SetDispatcher` support.

### 2. Shared version contract

Define a small unexported interface inside `internal/fawry` (not exported, not in `provider` package yet):

```go
type versionProvider interface {
    provider.Provider
    SetStore(engine.TransactionStore)
    SetDispatcher(*webhook.Dispatcher)
}
```

Both `v1.Provider` and `v2.Provider` implement it.

### 3. `internal/fawry/v1/`

- Move current charge and webhook behavior here.
- The V1 webhook payload is the legacy format (a simple JSON object, not `webhook.FawryV2Payload`).
- If a legacy payload type does not exist, define a minimal one in `v1/webhook.go`.

### 4. `internal/fawry/v2/`

- Reuse the charge logic from V1 (the charge request is the same across Fawry versions).
- The V2 webhook uses `webhook.FawryV2Payload` and `webhook.NewHMACSigner`.
- Move the existing `buildPayload` logic from the current top-level provider into `v2`.

### 5. Config defaults

Add to `internal/config/config.go` `DefaultYAML()`:

```yaml
providers:
  fawry:
    enabled: false
    config:
      merchant_code: ""
      merchant_security_key: ""
      webhook_secret: ""
      version: "v1"
```

### 6. Plugin manifest

Update `plugins/fawry/gateway.yml`:

- Add `versions` metadata listing supported versions.
- Add explicit V1 and V2 routes.
- Keep the existing fixture; add a V2 webhook fixture.

Example metadata addition:

```yaml
metadata:
  name: fawry
  version: 1.0.0
  supported_versions:
    - v1
    - v2
```

### 7. Tests

- Move/adapt existing Fawry tests into `internal/fawry/v1/*_test.go`.
- Add `internal/fawry/v2/webhook_test.go` verifying the V2 payload shape and HMAC signature.
- Add top-level `internal/fawry/provider_test.go` verifying:
  - default version is v1,
  - unknown version is rejected,
  - legacy routes work,
  - versioned routes work.

### 8. Smoke test

Extend `scripts/smoke-test.sh` with a Fawry V2 notification flow:

1. Charge via `/fawry/charge` (legacy, still v1).
2. POST to `/fawry/v2/webhook` with a V2 payload.
3. Verify the webhook receiver accepted the delivery.

### 9. Documentation

- Update `docs/initiatives/openmuara-provider-versioning/TRACKING.md` to mark P01 ✅.
- Add a short section to `docs/architecture.md` or `runbooks/local-development.md` describing the new layout convention.

## Acceptance Criteria

- `go build ./...` passes.
- `go test ./...` passes.
- `go vet ./...` passes.
- `golangci-lint run` reports zero issues.
- `./scripts/smoke-test.sh` passes, including the new Fawry V2 check.
- Legacy `/fawry/charge` and `/fawry/webhook` still behave like V1 by default.
- `/fawry/v2/webhook` returns/behaves differently from `/fawry/v1/webhook`.
- Single-version providers (e.g., `internal/billplz/`) are untouched and still compile.

## Non-Goals

- Do not migrate any provider other than Fawry.
- Do not add a public `Version` field to the `provider.Provider` interface.
- Do not force single-version providers to adopt the `v1/` sub-package.
