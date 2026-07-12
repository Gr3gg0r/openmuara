> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara UX Excellence — Decision Log

> **Created:** 2026-07-01
> **Last Updated:** 2026-07-01
> **Status:** ⬜ Draft

---

## Decisions

### D001 — Additive UX only

All UX improvements must preserve existing CLI commands, config keys, and routes. New behavior is opt-in or layered on top.

- **Status:** ✅ Decided
- **Rationale:** Prevents breaking current users and keeps the initiative focused.
- **Consequences:** Legacy config keys and non-interactive usage must continue to work.

### D002 — Wizard interactivity

`muara init` is interactive by default when stdin is a TTY. Non-TTY and explicit `--defaults` / `--non-interactive` flags skip the questions and write the generic config.

- **Status:** ✅ Decided
- **Rationale:** Gives new users guidance without breaking CI, scripts, or AI agents.
- **Consequences:** `init.go` must detect TTY and accept a skip flag.

### D003 — Onboarding state is derived, not persisted

The dashboard onboarding checklist is computed from existing runtime data (providers enabled, transactions, webhooks) rather than a separate progress file.

- **Status:** ✅ Decided
- **Rationale:** Avoids extra persistence complexity and keeps state consistent with reality.
- **Consequences:** Checklist cannot be manually dismissed permanently; it can only be hidden/collapsed in the UI.

### D004 — Config line numbers are best-effort

Validation errors report `field`, `message`, `hint`, `file`, and `line` when possible. Line numbers come from a raw YAML parse; if unavailable, the error falls back to field path + file path.

- **Status:** ✅ Decided
- **Rationale:** Viper unmarshaling loses source location; a separate YAML parse is needed for line numbers.
- **Consequences:** Tests must accept both full and fallback error formats.

### D005 — Dashboard primary view is the ledger

The default dashboard landing view is a ledger of transactions and webhooks, inspired by Mailpit's inbox but named for the payments domain. Legacy tables remain as secondary tabs.

- **Status:** ✅ Decided
- **Rationale:** "Ledger" is the correct financial term; the existing engine is already a transaction ledger.
- **Consequences:** Endpoint is `/_admin/ledger`, not `/_admin/inbox`.

### D006 — Default first-time provider recommendation

The first-run wizard recommends **Fawry** as the default provider because it has the most complete local escape flow, with **Stripe** as the secondary option.

- **Status:** ✅ Decided
- **Rationale:** Fawry's `/_admin/fawry-escape` page lets a user complete an end-to-end payment without an external SDK.
- **Consequences:** Other providers can still be selected; this is only the default suggestion.

### D007 — Dashboard admin API versioning

Admin endpoints introduced by this initiative (`/_admin/ledger`, `/_admin/onboarding`, etc.) are versioned by path prefix `/_admin/v1/` where feasible, or by treating `/_admin` as a single evolving surface with additive-only changes.

- **Status:** ✅ Decided
- **Rationale:** Prevents future dashboard changes from breaking external scripts or AI agents that parse admin JSON.
- **Consequences:** New endpoints should be additive. If a breaking change is ever needed, a new `/_admin/v2/` prefix must be introduced and the old path kept for at least one minor version.
