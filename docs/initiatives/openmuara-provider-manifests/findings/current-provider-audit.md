> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# Current Provider Audit

Snapshot taken at `dev` HEAD `db8912d`.

| Provider | Manifest Exists | `runtime.type` | Go Package Auto-Registers | Notes |
|---|---|---|---|---|
| `default` | No | — | No | Internal fallback provider |
| `fawry` | Yes | `simple` | Yes | Needs `init()` removal (P03) |
| `senangpay` | Yes | `simple` | Yes | Needs `init()` removal (P03) |
| `ipay88` | Yes | `go` | Yes | Needs factory registration (P02) and `init()` removal (P03) |
| `billplz` | Yes | `go` | Yes | Needs factory registration (P02) and `init()` removal (P03) |
| `toyyibpay` | Yes | `go` | Yes | Needs factory registration (P02) and `init()` removal (P03) |
| `stripe` | No | — | Yes | Needs manifest (P04), factory registration (P02), and `init()` removal (P03) |

## Gaps

1. Loader prefers built-ins over manifests (P01).
2. `runtime.type: go` is declared in YAML but not honored by loader (P01/P02).
3. All Go providers auto-register in `init()` (P03).
4. No factory registry exists (P02).
5. `stripe` lacks a `gateway.yml` manifest (P04).
6. Contributor docs lack simple vs go guidance (P04).

## Target State

- Every non-default provider has a manifest.
- `runtime.type: simple` providers have no Go registration.
- `runtime.type: go` providers register a factory and are activated by the manifest.
- Removing a provider's Go package (or its manifest) cleanly removes it from discovery.

## References

- `appendices/a-provider-contract-checklist.md`
- `appendices/b-simple-vs-go-decision-tree.md`
- `prompts/01-make-loader-manifest-first.md`
- `prompts/02-add-go-factory-registry.md`
- `prompts/03-remove-builtin-auto-registration.md`
- `prompts/04-migrate-remaining-providers.md`
