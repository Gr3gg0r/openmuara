> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara UX Excellence

> **Status:** ⬜ Not Started | **Started:** —
> **Scope:** Make OpenMuara the most approachable local payment emulator for four audiences: developers integrating it into their workflow, AI agents bootstrapping from CLI `--help`, testers inspecting traffic through a ledger-style web UI, and contributors extending the provider set.
> **Owner:** AI Agent (Kimi Code) | **Human Reviewer:** ___________
> **Target Repo:** `<repo-root>/`
> **Product Branch:** `feat/ux-excellence`

---

## Initiative Structure

```
docs/initiatives/openmuara-ux-excellence/
├── README.md              # This file
├── TRACKING.md            # Central execution tracker
├── HANDOFF.md             # Session continuity
├── DECISIONS.md           # Decision log
├── RISKS.md               # Risk register
├── KNOWN_ISSUES.md        # Pre-existing UX gaps
│
├── prompts/               # Numbered, self-contained execution prompts
│   ├── _template.md
│   ├── 01-first-run-wizard.md
│   ├── 02-dashboard-onboarding.md
│   ├── 03-config-validation-errors.md
│   ├── 04-provider-selection-guide.md
│   ├── 05-webhook-debugger.md
│   ├── 06-transaction-search-and-replay.md
│   ├── 07-cli-help-and-structured-output.md
│   ├── 08-ledger-style-payment-view.md
│   └── 09-docs-quickstart.md
│
├── findings/              # UX audit output
├── GLOSSARY.md            # Shared terminology
└── .gitignore             # Ignore screenshots, logs
```

Planning docs live in `docs/initiatives/openmuara-ux-excellence/` in the root repo. Product code commits to the `feat/ux-excellence` branch. Do not commit directly to `main`.

---

## Why UX excellence?

OpenMuara already emulates multiple payment providers faithfully. The next differentiator is how easy it is to pick up, configure, debug, and extend. Right now a new user has to:

1. Build the binary.
2. Run `muara init` and get a generic YAML file.
3. Guess which providers to enable and how to point their test app.
4. Discover the dashboard at `/_admin` by reading docs.
5. Read provider docs to understand version choices, signatures, and webhook shapes.
6. Debug failed webhooks by reading logs.

Each of those steps is a drop-off point. This initiative turns them into a guided, forgiving experience.

---

## Personas

### Developer

A developer wants to emulate a payment provider locally while building or testing their app. They need:

- A fast first-run setup that matches their real provider.
- Clear, actionable errors when config is wrong.
- Copy-paste examples for their SDK and language.
- Confidence that OpenMuara behaves like the real provider.

### AI Agent

An AI agent (or automation script) starts from `muara --help` and needs:

- Self-describing CLI commands with runnable examples.
- Structured output options (`--json`, `--quiet`) for parsing.
- A deterministic path from install to running server.
- No hidden interactive prompts in headless mode.

### Tester

A tester (QA engineer or developer debugging a flow) lives in the web UI, similar to [Mailpit](https://github.com/axllent/mailpit). They need:

- A real-time, auto-refreshing ledger of transactions and webhooks.
- One-click inspection of request/response payloads.
- Signature verification status at a glance.
- Search, filter, and replay without writing code.
- Clear pass/fail indicators and a timeline of attempts.

### Contributor

A contributor adding a new provider or extending OpenMuara needs:

- A clear checklist for making a provider discoverable in the wizard, dashboard, and provider guide.
- Stable interfaces and examples in existing providers.
- Tests and docs conventions that are easy to follow.

---

## Goals

1. **First-run wizard** — After `muara init`, ask the user what they are building and generate a tailored config with the right providers enabled and sample credentials populated.
2. **Dashboard onboarding** — The web dashboard shows a visible getting-started checklist: server ready, providers enabled, first charge sent, webhook received.
3. **Configuration validation with actionable errors** — `muara start` and `muara doctor` report config problems in plain language, with line numbers and suggested fixes.
4. **Provider selection guide** — Dashboard and CLI explain which provider to pick based on the user's real provider (Stripe, Fawry, Billplz, etc.) and link to sample SDK code; includes a contributor checklist for adding new providers.
5. **Webhook debugger** — Dashboard exposes a webhook inspector: payload body, signature verification status, retry timeline, and one-click replay.
6. **Transaction search and replay** — Search transactions by reference, provider, amount, or status; replay a charge or webhook for any record.
7. **CLI help and structured output** — Every CLI command shows copy-paste examples and links to the matching runbook; supports `--json` and `--quiet` for AI agents with documented schemas.
8. **Dashboard ledger** — Auto-refreshing, visibility-aware ledger view of transactions and webhooks, with one-click payload inspection, signature status, search, filter, and replay. Inspired by Mailpit's inbox but shaped around OpenMuara's transaction ledger.
9. **Documentation quickstart** — A single "Quick Start" page gets a new user from zero to their first successful charge in under five minutes, with per-audience paths.

---

## Conventions (MANDATORY)

### 1. Read `AGENTS.md` first
Root `AGENTS.md` governs branch rules, quality gates, autonomy boundaries, and code style.

### 2. Backward compatibility
All changes must keep existing configs, routes, and CLI commands working. New UX is additive.

### 3. No external services
OpenMuara is local-first. Do not add telemetry, analytics, or cloud dependencies.

### 4. Progressive disclosure
Show the simplest path first; advanced options stay behind collapsible sections or `--advanced` flags.

### 5. Three-interface parity
Any feature that is useful in the web UI should also be discoverable from the CLI, and vice versa, so no persona is forced into the wrong interface.

Examples:
- The ledger view in the dashboard is mirrored by CLI commands such as `muara transactions list --json` and `muara webhooks list --json`.
- Dashboard filters (provider, status, search) use the same query parameters as the admin API.
- CLI `--json` output fields match the shapes returned by `/_admin/ledger`, `/_admin/transactions`, and `/_admin/webhooks`.

### 6. Quality gates
Every prompt must pass:

- `go build ./...`
- `go test ./...`
- `go vet ./...`
- `golangci-lint run`
- `./scripts/smoke-test.sh`

### 7. Definition of done
Beyond the quality gates, a prompt is done only when:

- The feature works end-to-end for at least one persona.
- Tests cover happy path, error path, and at least one edge case.
- The smoke test is updated if the feature changes routes, CLI flags, or default behavior.
- `HANDOFF.md` is updated with what was built and what changed.
- `TRACKING.md` marks the prompt `✅` with the commit hash.
- User-facing changes are noted for the next release notes.

---

## Out of Scope

- Redesigning the provider plugin schema.
- Adding authentication to the dashboard.
- Hosting OpenMuara as a managed service.
- Mobile apps or desktop installers.
- A/B testing frameworks or analytics.

---

## Success criteria

- **Developer:** can run `muara init`, answer ≤5 questions, and have a working config for their target provider.
- **AI Agent:** can run `muara --help`, read a runnable example, and start a working server without interactive prompts; `--json` output follows documented schemas.
- **Tester:** can open `/_admin`, see a real-time payment ledger, search/filter transactions and webhooks, inspect payloads, and replay with one click — without reading server logs.
- **Contributor:** can read the provider guide checklist and add a new provider that appears in the wizard, dashboard, and docs.
- The dashboard landing page shows a completed checklist after the first smoke-test run.
- `muara doctor` reports config errors with field path, file path, and best-effort line number.
- Each CLI command has a runnable `--help` example.
- Webhook failures can be diagnosed from the dashboard without reading server logs.
- All quality gates pass.

## Metrics

We will judge success with these measurable signals:

| Metric | Target | How measured |
|--------|--------|--------------|
| Time from `muara init` to first successful charge | ≤5 minutes | Manual quick-start validation |
| Dashboard tasks completed without reading logs | ≥90% of common flows | Smoke test + manual QA |
| CLI `--help` examples coverage | 100% of commands | Automated test |
| Config validation errors resolved on first hint | ≥80% | Manual test of common misconfigs |
| New contributor time to add a provider | ≤30 minutes | Walkthrough with contributor checklist |
