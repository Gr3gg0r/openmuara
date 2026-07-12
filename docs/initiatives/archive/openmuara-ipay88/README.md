> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara iPay88

> **Status:** ⬜ Not Started | **Started:** —
> **Scope:** Implement a faithful iPay88 Malaysia classic ePayment API emulation so Southeast Asian developers can test iPay88 payments locally before signing up for a real iPay88 account.
> **Owner:** AI Agent (Kimi Code) | **Human Reviewer:** ___________
> **Target Repo:** `<repo-root>/`
> **Product Branch:** `dev`

---

## Initiative Structure

```
docs/initiatives/openmuara-ipay88/
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
│   └── 01-ipay88-provider.md
│
├── tasks/                 # (Optional) Detailed specs — dual-layer
├── findings/              # Research, audit output, analysis
├── runbooks/              # Operational docs
├── screenshots/           # QA evidence (gitignored)
├── qa/                    # Validation artifacts (gitignored)
└── state/                 # Agent state snapshots (gitignored)
```

Planning docs live in `docs/initiatives/openmuara-ipay88/` in the root repo. Product code commits to the `dev` branch. Do not commit directly to `main`.

---

## Why iPay88?

iPay88 is a long-standing payment gateway serving Malaysia and Southeast Asia. It supports FPX, credit/debit cards, and e-wallets through a single form-based integration. Many established Malaysian businesses use iPay88 because of its long history and broad payment-method support.

OpenMuara's mission is to let developers test financial infrastructure locally before they have real provider accounts. Emulating iPay88 gives developers another regional option for local testing.

This initiative targets the **Malaysia classic form-based ePayment API** (`/ePayment/entry.asp`), not the newer Indonesia JSON API.

---

## Goals

1. Implement iPay88 provider registration as `ipay88`.
2. Implement iPay88 Malaysia classic ePayment flow:
   - `POST /ePayment/entry.asp` — submit payment request (redirect form)
   - `POST /ePayment/enquiry.asp` — requery payment status
3. Support payment methods: FPX, credit/debit card, e-wallet (via numeric `PaymentId`).
4. Render a local iPay88 payment page at the redirect URL.
5. Implement iPay88 response/backend callback flow faithfully:
   - Browser **response**: `POST` to the merchant's `ResponseURL` with form fields and signature.
   - Server-side **backend**: `POST` to the merchant's `BackendURL` with form fields and signature; the merchant must respond with plain text `RECEIVEOK`.
6. Implement iPay88 SHA256 signature algorithm with correct field concatenation and amount normalization.
7. Add tests and smoke-test coverage.

---

## Conventions (MANDATORY)

### 1. Read `AGENTS.md` first
Root `AGENTS.md` governs branch rules, quality gates, autonomy boundaries, and code style. This initiative does not repeat every rule.

### 2. Provider contract fidelity
iPay88 emulation must match iPay88's documented behavior for the implemented subset, including:
- Form-encoded request and response bodies
- SHA256 signature algorithm with `SignatureType=SHA256`
- Payment status codes (`1` success, `0` fail, `6` pending) for callbacks
- Requery status codes (`00` success, others failure)
- `RECEIVEOK` acknowledgement for backend posts

### 3. Quality gates
Every prompt must pass:
- `go build ./...`
- `go test ./...`
- `go vet ./...`
- `golangci-lint run`
- `./scripts/smoke-test.sh`

---

## Out of Scope

- iPay88 merchant onboarding / admin portal APIs
- iPay88 Indonesia JSON API 2.0+
- Recurring payments / tokenization
- Multi-currency beyond MYR
- Real FPX/crypto cryptography
- Exact iPay88-hosted page styling parity
