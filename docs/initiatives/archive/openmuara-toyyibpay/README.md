> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara ToyyibPay

> **Status:** ⬜ Not Started | **Started:** —
> **Scope:** Implement a faithful ToyyibPay API emulation so Malaysian developers can test ToyyibPay payments locally before signing up for a real ToyyibPay account.
> **Owner:** AI Agent (Kimi Code) | **Human Reviewer:** ___________
> **Target Repo:** `<repo-root>/`
> **Product Branch:** `dev`

---

## Initiative Structure

```
docs/initiatives/openmuara-toyyibpay/
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
│   └── 01-toyyibpay-provider.md
│
├── tasks/                 # (Optional) Detailed specs — dual-layer
├── findings/              # Research, audit output, analysis
├── runbooks/              # Operational docs
├── screenshots/           # QA evidence (gitignored)
├── qa/                    # Validation artifacts (gitignored)
└── state/                 # Agent state snapshots (gitignored)
```

Planning docs live in `docs/initiatives/openmuara-toyyibpay/` in the root repo. Product code commits to the `dev` branch. Do not commit directly to `main`.

---

## Why ToyyibPay?

ToyyibPay is a Malaysian payment gateway similar to Billplz, popular with small merchants and solo operators because of simple setup and low fees. It uses a "category + bill" model and supports FPX, cards, FPX B2B, and DuitNow QR.

OpenMuara's mission is to let developers test financial infrastructure locally before they have real provider accounts. Emulating ToyyibPay gives Malaysian developers another local option, especially for projects like `atur`/`potongq` that may prefer ToyyibPay over Stripe or Billplz.

---

## Goals

1. Implement ToyyibPay provider registration as `toyyibpay`.
2. Implement ToyyibPay API subset:
   - `POST /index.php/api/createCategory` — create a category
   - `POST /index.php/api/getCategoryDetails` — retrieve category details
   - `POST /index.php/api/createBill` — create a bill
   - `POST /index.php/api/getBillTransactions` — retrieve bill transactions
   - `POST /index.php/api/inactiveBill` — deactivate a bill
3. Support payment methods: FPX, card, FPX B2B, DuitNow QR (via `billPaymentChannel`).
4. Render a local ToyyibPay payment page at the bill URL.
5. Implement ToyyibPay return-URL and callback flow faithfully:
   - Browser **return URL**: `GET` to `billReturnUrl` with `status_id`, `billcode`, `order_id`, etc.
   - Server-side **callback**: `POST` to `billCallbackUrl` with `refno`, `status`, `reason`, `billcode`, `order_id`, `amount`, `transaction_time`, and MD5 `hash`.
6. Implement MD5 callback hash: `MD5(userSecretKey + status + order_id + refno + "ok")`.
7. Add tests and smoke-test coverage.

---

## Conventions (MANDATORY)

### 1. Read `AGENTS.md` first
Root `AGENTS.md` governs branch rules, quality gates, autonomy boundaries, and code style. This initiative does not repeat every rule.

### 2. Provider contract fidelity
ToyyibPay emulation must match ToyyibPay's documented behavior for the implemented subset, including:
- Form-encoded request bodies (not JSON)
- Form-encoded callback bodies with MD5 `hash`
- Correct endpoint names (`getCategoryDetails`, not `getCategory`)
- Transaction status values (`1` success, `2` pending, `3` fail, `4` pending)

### 3. Quality gates
Every prompt must pass:
- `go build ./...`
- `go test ./...`
- `go vet ./...`
- `golangci-lint run`
- `./scripts/smoke-test.sh`

---

## Out of Scope

- ToyyibPay disbursements
- Merchant onboarding / API key management endpoints
- Multi-currency beyond MYR
- Real FPX/DuitNow QR cryptography
- Exact ToyyibPay-hosted page styling parity
