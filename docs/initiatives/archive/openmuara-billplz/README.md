> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Billplz

> **Status:** ⬜ Not Started | **Started:** —
> **Scope:** Implement a faithful Billplz v3 API emulation so Malaysian developers can test Billplz payments locally before signing up for a real Billplz account.
> **Owner:** AI Agent (Kimi Code) | **Human Reviewer:** ___________
> **Target Repo:** `<repo-root>/`
> **Product Branch:** `dev`

---

## Initiative Structure

```
docs/initiatives/openmuara-billplz/
├── README.md              # This file
├── HOWTO.md               # Decomposition guide for AI
├── PREREQUISITES.md       # Human pre-flight checklist
├── TRACKING.md            # Central execution tracker
├── HANDOFF.md             # Session continuity
├── DECISIONS.md           # Decision log
├── RISKS.md               # Risk register
├── KNOWN_ISSUES.md        # Pre-existing bugs / out-of-scope
├── REFERENCES.md          # Links to specs, runbooks, vendor docs
├── .gitignore             # Ignore screenshots, logs, temp files
│
├── prompts/               # Numbered, self-contained execution prompts
│   ├── _template.md
│   └── 01-billplz-provider.md
│
├── tasks/                 # (Optional) Detailed specs — dual-layer
├── findings/              # Research, audit output, analysis
├── runbooks/              # Operational docs
├── screenshots/           # QA evidence (gitignored)
├── qa/                    # Validation artifacts (gitignored)
└── state/                 # Agent state snapshots (gitignored)
```

Planning docs live in `docs/initiatives/openmuara-billplz/` in the root repo. Product code commits to the `dev` branch. Do not commit directly to `main`.

---

## Why Billplz?

Billplz is one of the most popular payment gateways in Malaysia. It is especially common among SMEs because of lower fees than international providers and strong local payment method coverage (FPX, cards, e-wallets, BNPL).

OpenMuara's mission is to let developers test financial infrastructure locally before they have real provider accounts. Emulating Billplz allows Malaysian developers — including the `atur`/`potongq` project — to build paywalls and billing flows locally using the real Billplz SDK or direct API calls, then switch to production by changing only the base URL and API key.

---

## Goals

1. Implement Billplz provider registration as `billplz`.
2. Implement Billplz v3 API subset:
   - `POST /api/v3/collections` — create a collection
   - `GET /api/v3/collections/{id}` — retrieve a collection
   - `POST /api/v3/bills` — create a bill
   - `GET /api/v3/bills/{id}` — retrieve a bill
   - `DELETE /api/v3/bills/{id}` — delete a bill
   - `GET /api/v3/collections/{id}/payment_methods` — list available payment methods
3. Support payment methods: FPX (`fpx`), card (`mpgs`), e-wallets (`boost`, `touchngo`, etc.), and BNPL (`twoctwopipp`).
4. Render a local Billplz payment page at the bill URL so the customer can complete payment.
5. Implement Billplz redirect and callback flow faithfully:
   - Browser **redirect**: `GET` to the bill's optional `redirect_url` with `billplz[id]`, `billplz[paid]`, `billplz[state]`, and `x_signature` query params.
   - Server-side **callback**: `POST` to the bill's required `callback_url` with a form-urlencoded Bill object and `x_signature`.
6. Implement `x_signature` as documented: HMAC-SHA256 of sorted `key+value` pairs joined by `|`.
7. Add tests and smoke-test coverage.

---

## Conventions (MANDATORY)

### 1. Read `AGENTS.md` first
Root `AGENTS.md` governs branch rules, quality gates, autonomy boundaries, and code style. This initiative does not repeat every rule.

### 2. Provider contract fidelity
Billplz emulation must match Billplz's documented v3 behavior for the implemented subset, including:
- Request/response JSON shapes
- HTTP status codes
- Bill states: `due`, `paid`, `deleted`
- `x_signature` as a form/query parameter, **not** an HTTP header
- Callback redirect behavior

### 3. Quality gates
Every prompt must pass:
- `go build ./...`
- `go test ./...`
- `go vet ./...`
- `golangci-lint run`
- `./scripts/smoke-test.sh`

---

## Out of Scope

- Billplz disbursements / payouts
- Billplz API keys / user management endpoints
- Multi-currency beyond MYR in this initiative
- Real FPX/crypto cryptography
- Exact Billplz-hosted page styling parity
