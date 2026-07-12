> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# Appendix F — Architecture Diagram

ASCII view of the manifest-first provider discovery flow after this initiative.

```
┌─────────────────────────────────────────────────────────────────┐
│                         Startup                                  │
└───────────────────────────┬─────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│  1. Load config from .muara/config.yml                           │
│     (provider names, secrets, feature flags)                     │
└───────────────────────────┬─────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│  2. Walk plugins/ directory                                       │
│     For each plugins/<name>/gateway.yml:                         │
│     • Parse YAML                                                  │
│     • Validate schema                                             │
│     • Read runtime.type                                           │
└───────────────────────────┬─────────────────────────────────────┘
                            │
             ┌──────────────┼──────────────┐
             │              │              │
             ▼              ▼              ▼
    ┌─────────────┐  ┌─────────────┐  ┌─────────────┐
    │  simple     │  │     go      │  │   default   │
    │  runtime    │  │   runtime   │  │  fallback   │
    │             │  │             │  │             │
    │ internal/   │  │ factory.Get │  │  hard-coded │
    │ provider/   │  │  by name    │  │             │
    │ simple/     │  │             │  │             │
    └──────┬──────┘  └──────┬──────┘  └──────┬──────┘
           │                │                │
           └────────────────┼────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│  3. Provider instances available to router, engine, dashboard    │
└─────────────────────────────────────────────────────────────────┘

Key rule: No manifest → no provider (except default).
```

## Package Layout (recommended)

```
internal/
├── config/
│   └── provider_loader.go      # Walks plugins/ and routes by runtime.type
├── plugin/
│   ├── schema.go               # Gateway YAML types
│   └── validator.go            # Manifest validation
├── provider/
│   ├── simple/
│   │   └── provider.go         # runtime.type: simple implementation
│   ├── factory/
│   │   ├── registry.go         # Go factory registry
│   │   └── factory.go          # Factory type definition
│   └── conform/                # Protocol conformance tests
├── <provider>/                 # One per runtime.type: go provider
│   ├── provider.go             # Provider implementation
│   └── register.go             # Factory registration
└── server/
    └── router.go               # Uses loaded providers from config

plugins/
├── <name>/
│   └── gateway.yml             # Manifest: name, runtime.type, config
```

## Data Flow

1. **Config** says which providers the user cares about.
2. **Manifests** say how each provider is implemented.
3. **Loader** reads manifests and picks the runtime.
4. **Runtime** creates the provider instance.
5. **Router** serves provider endpoints using the instance.

## Future Extensions

```
plugins/<name>/gateway.yml
        │
        ├── runtime.type: simple  → internal/provider/simple/
        ├── runtime.type: go      → internal/provider/factory/ → internal/<name>/
        ├── runtime.type: bridge  → external service (future)
        └── runtime.type: wasm    → .muara/plugins/<name>.wasm (future)
```
