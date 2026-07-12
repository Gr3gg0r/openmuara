> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Provider Manifests — Recommendations

This document records the recommended resolutions for open decisions in `DECISIONS.md`, plus future enhancements that are out of scope for this initiative.

---

## Recommended Resolutions for Open Decisions

### RD004 — Keep `default` provider hard-coded

**Status:** ✅ Approved / Pinned (D004)
**Reversibility:** High

**Recommendation:** Keep the `default` provider hard-coded as an internal fallback. Do not give it a manifest in this initiative.

**Rationale:**
- `default` is not a real payment provider; it is a bootstrapping/test helper.
- Adding a manifest would confuse contributors into thinking `default` is a provider pattern to copy.
- Keeps the initiative scope minimal.

**If rejected:** Give `default` a manifest at `plugins/default/gateway.yml` but load it implicitly. Never require users to configure it.

---

### RD006 — Factory registry package location

**Status:** ✅ Approved / Pinned (D006)
**Reversibility:** Medium

**Recommendation:** Place the registry in `internal/provider/factory/`.

**Rationale:**
- The name is explicit and discoverable.
- It keeps runtime code under `internal/provider/`, which matches the simple runtime location.
- Future runtimes (`bridge`, `wasm`) can live in sibling packages (`internal/provider/bridge/`, `internal/provider/wasm/`) without confusion.
- `internal/provider/hybrid/` is ambiguous — "hybrid" could mean many things.
- `internal/plugin/registry.go` mixes plugin schema concerns with runtime instantiation.

**Package contents:**

```
internal/provider/factory/
├── registry.go      # Register, Get, Names
├── factory.go       # Factory type definition
└── registry_test.go # Unit tests
```

**If rejected:** `internal/provider/hybrid/` is the second-best option. Avoid `internal/plugin/` for runtime code.

---

### RD007 — Migration warning for existing users

**Status:** ✅ Approved / Pinned (D007)
**Reversibility:** High

**Recommendation:** Soft landing, then hard fail.

1. **This release:** When a configured provider has no manifest, print a clear warning:
   ```
   provider "billplz" is configured but has no gateway.yml manifest; 
   auto-loading built-in providers is deprecated. 
   See docs/migration/provider-manifests.md
   ```
   Still load the provider for backwards compatibility.
2. **Next release:** Fail hard with a helpful error.

**Rationale:**
- OpenMuara is local-first and has few production users; a breaking change is acceptable with notice.
- A warning gives the sole user (and any early adopters) time to migrate.
- The migration path is trivial: add the existing built-in manifest to `plugins/`.

**If rejected:** Fail hard immediately and bump the minor/major version with a prominent migration note.

---

### RD008 — Provider factory signature

**Status:** ✅ Approved / Pinned (D008)
**Reversibility:** Medium

**Recommendation:** Start minimal:

```go
type Factory func(cfg map[string]any) (provider.Provider, error)
```

**Rationale:**
- Keeps the first implementation simple and reviewable.
- Avoids inventing a `Deps` abstraction before it is needed.
- If providers later need shared services (logger, dispatcher, webhook dispatcher), we can extend the signature or pass them through `cfg` as a non-breaking change if designed carefully.

**If rejected:** Use an explicit `Deps` struct:

```go
type Factory func(cfg map[string]any, deps Deps) (provider.Provider, error)
```

But only do this if at least two providers need the same dependency.

---

## Future Enhancements (Out of Scope)

These are intentionally not part of this initiative. They are recorded so they are not lost.

| ID | Enhancement | Priority | Rationale | Target |
|---|---|---|---|---|
| E001 | Implement `runtime.type: bridge` | Medium | Enables proprietary/private providers without open-sourcing code. | Future initiative |
| E002 | Implement `.muara/plugins/<name>.wasm` runtime | Medium | Sandbox third-party providers; no recompilation. | Future initiative |
| E003 | `muara provider diagnose <name>` | Low | Suggest `runtime.type` based on manifest content. | Future spike |
| E004 | `muara provider scaffold` CLI command | Low | Generate `gateway.yml` + `register.go` boilerplate. | Future enhancement |
| E005 | Provider versioning in manifests | Medium | Allow `v1`/`v2` tabs in dashboard per provider. | Future initiative |
| E006 | Manifest hot-reload | Low | Reload providers without restarting Muara. | Future spike |
| E007 | JSON Schema export for `gateway.yml` | Low | IDE autocompletion for contributors. | Future enhancement |
| E008 | Conformance test generator | Low | Generate conformance tests from manifest examples. | Future spike |

---

## How to Apply These Recommendations

1. ✅ Reviewed and approved by human reviewer on 2026-07-09.
2. ✅ `DECISIONS.md` marks D004, D006, D007, D008 as ✅ Pinned.
3. ✅ `TRACKING.md` closes Q001–Q004.
4. Proceed with implementation using the approved choices.
