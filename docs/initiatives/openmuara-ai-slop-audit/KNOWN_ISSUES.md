> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara AI Slop Audit — Known Issues

> **Created:** 2026-07-08
> **Last Updated:** 2026-07-08
> **Status:** ⬜ Draft

---

## Dashboard

### F001 — Sad-face empty state is off-tone
- **Evidence:** `web/dashboard/src/components/EmptyState.tsx`, `screenshots/webhooks-list-desktop.png`
- **Detail:** A frowny face for an empty webhook list feels apologetic rather than neutral/informative.
- **Suggested fix:** Swap `frown` icon for `inbox`, `search`, or a neutral illustration.
- **Status:** ⬜ Open

### F002 — Provider cards use generic placeholder descriptions
- **Evidence:** `web/dashboard/src/views/Settings.tsx`, `/_admin/providers` response
- **Detail:** billplz, ipay88, senangpay, stripe, and toyyibpay cards all show **"Provider configuration"** while Default/DIY and Fawry have real descriptions.
- **Suggested fix:** Add short, specific descriptions to each provider in `plugins/*/gateway.yml`.
- **Status:** ⬜ Open

### F003 — Redundant status badges on Default/DIY provider card
- **Evidence:** `screenshots/settings-desktop.png`
- **Detail:** Default/DIY shows both `ACTIVE` and `ENABLED`, which communicate nearly the same thing.
- **Suggested fix:** Keep one status badge per card (prefer `ENABLED`).
- **Status:** ⬜ Open

### F004 — Generic system-ui font stack
- **Evidence:** `web/dashboard/src/styles.css` (`--font-sans`, `body` font-family)
- **Detail:** The skill flags `system-ui`/`-apple-system` as the primary display font as the "gave up on typography" signal.
- **Suggested fix:** Choose one real typeface (e.g., Inter) and load it via Google Fonts or a local file.
- **Status:** ⬜ Open | **Deferred decision:** font asset choice needed.

### F005 — Preseed transactions are visible on first load
- **Evidence:** `screenshots/ledger-list-desktop.png`
- **Detail:** The dev build ships with preseed rows. The user has requested a runtime flag so production builds start empty while dev can keep seeding.
- **Suggested fix:** Gate seeding behind an env var or config flag (e.g., `MUARA_SEED_DATA=1`). Default to empty.
- **Status:** ⬜ Open | **Deferred:** requires backend/config change.

### F006 — Header action buttons used 28px touch targets
- **Evidence:** `web/dashboard/src/components/AppShell.tsx`, `web/dashboard/src/components/Shell.tsx`
- **Detail:** Buttons used `size="sm"` → `min-height: 28px`, below the 44px accessibility minimum.
- **Suggested fix:** Use `size="md"` or enforce `min-height: 36px` for header actions. *(Already fixed during 2026-07-08 responsiveness pass; verify and close.)*
- **Status:** ✅ Partially fixed — verify.

## Provider metadata

### F007 — Inconsistent provider descriptions and categories
- **Evidence:** `plugins/*/gateway.yml`
- **Detail:** Some providers have descriptions, categories, and emulated real-provider lists; others are bare.
- **Suggested fix:** Normalize every `gateway.yml` to include a concise `description`, `category`, and `real_providers` list where applicable.
- **Status:** ⬜ Open

## Documentation / prompts

### F008 — Vague or buzzword-heavy copy to audit
- **Evidence:** `docs/`, `prompts/`, `runbooks/`, `website/`
- **Detail:** Areas to scan for slop terms: *seamless, leverage, cutting-edge, revolutionary, empower, unlock, delightful, robust, scalable, innovative, game-changing, next-generation, world-class*.
- **Suggested fix:** Replace vague claims with concrete behavior, examples, or remove them.
- **Status:** ⬜ Open

### F009 — Placeholder TODOs / unfinished prompt sections
- **Evidence:** `prompts/`, `docs/initiatives/`
- **Detail:** Initiatives and prompts may contain `TODO`, `FIXME`, ellipses, or empty rationale sections.
- **Suggested fix:** Resolve or convert to explicit backlog items.
- **Status:** ⬜ Open

## Code-level slop

### F010 — `any` / `interface{}` usage
- **Evidence:** `internal/**/*.go` (to be audited)
- **Detail:** `AGENTS.md` prefers explicit types. Sloppy `any`/`interface{}` often hides half-finished abstractions.
- **Suggested fix:** Audit with `golangci-lint` and replace with concrete types where possible.
- **Status:** ⬜ Open

### F011 — Copy-paste provider emulation boilerplate
- **Evidence:** `internal/fawry/`, `internal/senangpay/`, `internal/ipay88/`, etc.
- **Detail:** Similar handlers, response builders, and signature logic duplicated across providers with only names changed.
- **Suggested fix:** Extract shared helpers without over-abstracting; keep provider-specific quirks explicit.
- **Status:** ⬜ Open

### F012 — Verbose comments that restate code
- **Evidence:** `internal/**/*.go`, `web/dashboard/src/**/*.tsx`
- **Detail:** Comments like `// fetch the config` above `fetchConfig()` add noise.
- **Suggested fix:** Delete or replace with *why* comments.
- **Status:** ⬜ Open

## Test / seed data

### F013 — Synthetic-looking seed transactions
- **Evidence:** Ledger preseed data
- **Detail:** Rows share the same provider, similar references, and predictable amounts, making the dashboard look like a fixture rather than realistic traffic.
- **Suggested fix:** If seeding stays in dev, vary providers, references, amounts, currencies, and statuses.
- **Status:** ⬜ Open

## Examples & website

### F014 — Generic example store copy
- **Evidence:** `examples/checkout-store/`
- **Detail:** Placeholder products, generic README instructions, and hardcoded values.
- **Suggested fix:** Give the example a coherent scenario (e.g., a fictional course/merchant) and realistic product SKUs.
- **Status:** ⬜ Open

### F015 — Docusaurus marketing copy
- **Evidence:** `website/docs/`, `website/src/pages/`
- **Detail:** Landing/docs copy may drift into generic AI phrasing.
- **Suggested fix:** Audit for slop terms and replace with specific OpenMuara behaviors and examples.
- **Status:** ⬜ Open
