> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Provider Manifests — Glossary

| Term | Definition |
|---|---|
| **Provider** | A module that emulates a payment provider's API surface (endpoints, signatures, webhooks). |
| **Manifest** | A `gateway.yml` file in `plugins/<name>/` that declares a provider's metadata, runtime type, and configuration schema. |
| **Runtime type** | The value of `runtime.type` in a manifest: `simple`, `go`, `bridge` (future), or `wasm` (future). |
| **Simple runtime** | A provider implemented entirely in YAML using the generic runtime in `internal/provider/simple/`. |
| **Go runtime** | A provider implemented partly or fully in Go, activated by a manifest and instantiated through a factory registry. |
| **Factory** | A Go function that constructs a provider instance when given configuration and optional dependencies. |
| **Factory registry** | A map keyed by provider name that holds factory functions. |
| **Auto-registration** | The anti-pattern of creating provider instances in package `init()` functions, causing side effects on import. |
| **Built-in provider** | A provider whose Go code lives in `internal/<name>/` and is compiled into the binary. |
| **Bridge provider** | A future runtime type for proprietary/private providers configured in `.muara/config.yml`. |
| **WASM plugin** | A future runtime type for sandboxed third-party providers loaded from `.muara/plugins/<name>.wasm`. |
| **Default provider** | The fallback provider used when no specific provider is selected; currently hard-coded. |
| **Conformance test** | A test that verifies a provider emulates its real-world protocol faithfully. |
| **Phantom provider** | A provider that appears in the runtime even though its manifest is absent or disabled. |
