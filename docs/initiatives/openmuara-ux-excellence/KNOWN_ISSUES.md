> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara UX Excellence — Known Issues

> **Created:** 2026-07-01

---

## Pre-existing UX gaps

1. **`muara init` is silent.** It writes a generic config and exits. A new user has to open the file and guess what to change.
2. **Dashboard has no onboarding.** A user who opens `/_admin` for the first time sees empty tables with no indication of what to do next.
3. **Config errors are reported by the underlying library.** Viper errors can be cryptic and do not include line numbers.
4. **Provider choice is undocumented in-product.** Users must read `docs/providers.md` to understand which provider maps to their real gateway.
5. **Webhook debugging requires log reading.** Failed webhooks are visible in the dashboard table, but the reason and payload are not.
6. **Transactions are not searchable.** The table grows linearly and there is no filter or replay action.
7. **CLI help lacks examples.** Commands list flags but do not show a common invocation.
8. **No single quick-start page.** A new user must cross-reference `runbooks/local-development.md`, `docs/providers.md`, and `docs/operations.md`.

---

## Out of scope for this initiative

- Real provider account linking.
- Cloud-hosted dashboards.
- Authentication or multi-user support.
