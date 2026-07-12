> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# Appendix B — Simple vs Go vs Bridge vs WASM Decision Tree

For contributors deciding which runtime type to use.

---

## Decision Tree

```
Is the provider publicly documented and emulatable?
├── No → Is it proprietary or private?
│   ├── Yes → Use bridge (future: configure in .muara/config.yml)
│   └── No → Is it a third-party experiment you don't want to compile in?
│       ├── Yes → Use wasm (future: drop .muara/plugins/<name>.wasm)
│       └── No → Open an issue to discuss architecture.
└── Yes → Does it fit the simple runtime model?
    ├── Yes → Use runtime.type: simple
    │           (common REST endpoints, standard signatures, no state machine)
    └── No → Use runtime.type: go
                (custom crypto, complex state, protocol quirks, webhooks)
```

---

## Simple Runtime

**Use when:**
- Provider is mostly request/response shaped.
- Signature verification is a standard HMAC/SHA256 or similar.
- No complex state machine between request and callback.
- Errors map cleanly to HTTP status codes.

**Examples in OpenMuara:**
- Fawry
- SenangPay

**Pros:**
- No Go code to write.
- Fast to iterate.
- Easy for non-Go contributors.

**Cons:**
- Limited to what the simple runtime supports.

---

## Go Runtime

**Use when:**
- Provider has custom signature schemes.
- Provider requires a state machine (e.g., checkout session → payment intent).
- Provider has non-standard webhook semantics.
- You need custom CLI commands or doctor checks.

**Examples in OpenMuara:**
- iPay88
- Billplz
- ToyyibPay
- Stripe

**Pros:**
- Full control over behavior.
- Can implement any protocol faithfully.

**Cons:**
- Requires Go knowledge.
- Must register a factory.
- More code to maintain.

---

## Bridge Runtime (Future)

**Use when:**
- Provider is proprietary and cannot be open-sourced.
- Provider logic lives in a private repository or service.
- You need to integrate OpenMuara with an existing internal emulator.

**Path:** `providers.<name>.type: bridge` in `.muara/config.yml`.

---

## WASM Runtime (Future)

**Use when:**
- You want to distribute a provider as a sandboxed plugin.
- You don't want to recompile OpenMuara to add a provider.
- You want to isolate a provider for security.

**Path:** `.muara/plugins/<name>.wasm`.
