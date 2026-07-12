> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Provider Manifests — Known Issues

Pre-existing gaps discovered during planning. These are not bugs to fix in this initiative unless they block a prompt; they are recorded so future work does not lose them.

| ID | Issue | Area | Severity | Blocking? | Workaround | Target Resolution |
|---|---|---|---|---|---|---|
| K001 | Loader prefers built-ins over manifests | `internal/config/provider_loader.go` | High | Yes — blocks P01 | None; fix in P01 | This initiative |
| K002 | `runtime.type: go` not implemented in loader | `internal/config/provider_loader.go` | High | Yes — blocks P02/P03 | None; fix in P01/P02 | This initiative |
| K003 | Built-in providers auto-register in `init()` | `internal/<provider>/provider.go` | High | Yes — blocks P03 | None; fix in P03 | This initiative |
| K004 | No Go factory registry | `internal/provider/` | High | Yes — blocks P02 | None; fix in P02 | This initiative |
| K005 | Tests depend on global `provider.Get("fawry")` | `internal/provider/conform/conform_test.go`, `internal/server/providers_test.go` | Medium | Yes — blocks P03 | Update tests to load manifests | This initiative |
| K006 | Stripe has no `gateway.yml` manifest | `plugins/stripe/` | Medium | Yes — blocks P04 | Create manifest in P04 | This initiative |
| K007 | Contributor docs lack simple vs go decision tree | `docs/contributing-providers.md` | Low | No | Use `appendices/b-simple-vs-go-decision-tree.md` as temporary guide | This initiative (P04) |
| K008 | No migration guide for users relying on auto-registration | `docs/migration/` | Low | No | Follow `RECOMMENDATIONS.md#rd007`; add warning in P03; write guide in P04 | This initiative (P04) |
| K009 | `muara doctor` does not list registered factories | `internal/cli/doctor.go` | Low | No | Manual code inspection | Future enhancement |
| K010 | `runtime.type: bridge` and `wasm` are architecture-only | `internal/plugin/schema.go` | Low | No | N/A | Future initiatives |

---

## Issues Resolved by This Initiative

When the initiative closes, the following should be removable from this list:

- K001, K002, K003, K004, K005, K006, K007, K008.
