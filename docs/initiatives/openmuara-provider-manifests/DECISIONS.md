> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Provider Manifests — Decision Log

## How to Use This Log

- Each decision gets a unique ID (`D###`).
- Pinned decisions may not be reversed without human sign-off.
- Open decisions must be resolved before the initiative closes.
- Record reversibility: can we undo this decision later without a breaking change?

---

## D001 — Provider Runtime Architecture

**Date:** 2026-07-08
**Status:** ✅ Pinned
**Reversibility:** Low — changes provider discovery contract.

We support four provider shapes:

| Use case | Path | Notes |
|---|---|---|
| Public common provider | `plugins/<name>/gateway.yml` with `runtime.type: simple` | No Go code required. Ideal for common REST-ish providers that fit the simple runtime model. |
| Public complex provider | `plugins/<name>/gateway.yml` with `runtime.type: go` + Go package in `internal/<name>/` | Manifest activates a Go factory registered at build time. |
| Proprietary/private provider | `providers.<name>.type: bridge` in `.muara/config.yml` | Future initiative. Keeps proprietary logic out of the public repo. |
| Sandboxed runtime plugin | `.muara/plugins/<name>.wasm` | Future initiative. Allows third-party providers without recompilation. |

**Rationale:**
- Lowers the barrier for contributors (start with 60% coverage in YAML).
- Lets advanced providers graduate to Go without changing the discovery path.
- Gives proprietary providers an escape hatch.
- Gives power users a sandboxed extension path.

**Consequences:**
- The provider contract becomes manifest-centric.
- Built-in Go packages must be refactored to expose factories.

---

## D002 — Manifest-First Discovery

**Date:** 2026-07-08
**Status:** ✅ Pinned
**Reversibility:** Low — changes loader behavior.

The loader must read `plugins/<name>/gateway.yml` first. Built-in Go providers register a factory, but the manifest decides whether that factory is activated.

**Rationale:**
- Removing a Go package must not break config loading if the manifest is absent.
- Tests should not rely on `provider.Get("fawry")` unless the manifest is loaded.

**Consequences:**
- Provider instantiation moves from package import side effects to explicit loader calls.
- Startup error messages must clearly identify missing/invalid manifests.

---

## D003 — No Built-in Auto-Registration

**Date:** 2026-07-08
**Status:** ✅ Pinned
**Reversibility:** Medium — can be partially undone by adding explicit registration calls.

Built-in providers must not register themselves in `init()` in a way that creates a provider instance. They must expose a factory function and let the loader or registry call it when the manifest says so.

**Rationale:**
- Prevents "phantom" providers from appearing in tests and production.
- Makes the manifest the single source of truth.

**Consequences:**
- Tests must be updated to load manifests or use factories explicitly.
- Router and other callers must get providers from the loader/config.

---

## D004 — Keep `default` Provider Special

**Date:** 2026-07-08
**Status:** ✅ Pinned
**Approved:** 2026-07-09 by human reviewer
**Reversibility:** High — default provider is an implementation detail.

The `default` provider remains hard-coded as a fallback for bootstrapping and tests. It does not get a manifest in this initiative.

**Rationale:**
- `default` is not a real payment provider; it is a bootstrapping/test helper.
- Adding a manifest would confuse contributors into thinking `default` is a provider pattern to copy.
- Keeps the initiative scope minimal.

**If revisited later:** Give `default` a manifest at `plugins/default/gateway.yml` but load it implicitly. Never require users to configure it.

---

## D005 — Planning Docs Commit First

**Date:** 2026-07-09
**Status:** ✅ Pinned
**Reversibility:** High — docs-only.

All planning docs for this initiative are committed in one docs-only commit before product code changes are committed.

**Rationale:**
- Keeps reviewable history.
- Avoids mixing architecture decisions with implementation noise.

---

## D006 — Factory Registry Package Location

**Date:** 2026-07-09
**Status:** ✅ Pinned
**Approved:** 2026-07-09 by human reviewer
**Reversibility:** Medium — package move is mechanical.

The Go factory registry lives in `internal/provider/factory/`.

**Options considered:**
1. `internal/provider/hybrid/` — emphasizes "hybrid runtime" (YAML + Go).
2. `internal/provider/factory/` — explicit and discoverable. **Selected.**
3. `internal/plugin/registry.go` — keeps plugin-related code together.

**Rationale:**
- The name is explicit and discoverable.
- It keeps runtime code under `internal/provider/`, matching the simple runtime location.
- Future runtimes (`bridge`, `wasm`) can live in sibling packages without confusion.
- `internal/provider/hybrid/` is ambiguous; `internal/plugin/` mixes schema with runtime.

**Expected package layout:**

```
internal/provider/factory/
├── registry.go      # Register, Get, Names
├── factory.go       # Factory type definition
└── registry_test.go # Unit tests
```

**Consequences:**
- Any existing code using `internal/provider/hybrid/` for the factory registry must be migrated to `internal/provider/factory/` during implementation.
- The hybrid provider wrapper (mixed simple+Go) may remain as a separate runtime helper or be folded into the loader.

---

## D007 — Migration Warning for Existing Users

**Date:** 2026-07-09
**Status:** ✅ Pinned
**Approved:** 2026-07-09 by human reviewer
**Reversibility:** High — warning text only.

Existing `.muara/config.yml` files may list providers that were previously auto-loaded. After this initiative, those providers need manifests.

**Decision:** Soft landing.

1. **This release:** When a configured provider has no manifest, print a clear warning and still load the provider for backwards compatibility.
2. **Next release:** Fail hard with a helpful error.

**Warning template:**

```
provider "<name>" is configured but has no gateway.yml manifest;
auto-loading built-in providers is deprecated.
See docs/migration/provider-manifests.md
```

**Documentation:** Create `docs/migration/provider-manifests.md` as part of P04.

**Rationale:**
- OpenMuara is local-first with few production users; a breaking change is acceptable with notice.
- The migration path is trivial: add the existing built-in manifest to `plugins/`.

**Consequences:**
- Loader must distinguish "manifest present" from "manifest absent" for configured providers.
- Tests should cover both the warning path and the future fail-hard path.

---

## D008 — Provider Factory Signature

**Date:** 2026-07-09
**Status:** ✅ Pinned
**Approved:** 2026-07-09 by human reviewer
**Reversibility:** Medium — signature changes touch all factories.

The factory function signature is minimal:

```go
type Factory func(cfg map[string]any) (provider.Provider, error)
```

**Rationale:**
- Keeps the first implementation simple and reviewable.
- Avoids inventing a `Deps` abstraction before it is needed.
- If providers later need shared services (logger, dispatcher, webhook dispatcher), extend the signature or pass them through `cfg` as a non-breaking change.

**Consequences:**
- Each Go provider package exposes a factory matching this signature.
- The registry maps provider names to factories.
- Shared dependencies are injected after instantiation via `SetStore`, `SetBaseURL`, `SetDispatcher`, etc.

---

## D009 — Conformance Test Strategy

**Date:** 2026-07-09
**Status:** ✅ Pinned
**Reversibility:** High — test-only.

Before changing provider loading, add or verify conformance tests for each provider. These tests must pass before and after the refactor to ensure protocol emulation does not regress.

---

## D010 — Documentation Updates Are Mandatory

**Date:** 2026-07-09
**Status:** ✅ Pinned
**Reversibility:** High — docs-only.

`docs/provider-contract.md` and `docs/contributing-providers.md` must be updated before the initiative closes. Docs are part of P04 acceptance criteria.
