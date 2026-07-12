> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This document is subordinate to it.**

# OpenMuara Bug Hunt — Review Checklist

> **Updated:** 2026-07-06
>
> Use this checklist before declaring `feat/bug-hunt` ready for PR to `dev`.

## Prompt Completion

- [ ] P01 reconnaissance complete; bug register has ≥5 findings.
- [ ] P02 triage complete; every finding is confirmed or marked false positive.
- [ ] P03 fix batch 1 complete; each fix has a regression test and passing gates.
- [ ] P04 fix batch 2 complete; remaining bugs fixed or deferred with rationale.
- [ ] P05 regression tests and quality gates complete; coverage did not drop.
- [ ] P06 visual sign-off complete; Playwright MCP screenshots attached.

## Quality Gates

- [ ] `go build ./...` passes.
- [ ] `go test ./...` passes.
- [ ] `go test -race ./...` passes.
- [ ] `go vet ./...` clean.
- [ ] `golangci-lint run` zero issues.
- [ ] `cd web/dashboard && npm run test:ci` passes.
- [ ] `cd web/dashboard && npm run build` passes.
- [ ] Bundle size within budget.
- [ ] A11y contrast check zero violations.

## Dashboard Invariants

- [ ] `/_admin` defaults to Ledger.
- [ ] Left nav has Ledger, Webhooks, Settings.
- [ ] Every table view has a filter toolbar.
- [ ] Ledger and webhook rows navigate to detail pages.
- [ ] Webhooks view is delivery-log only.
- [ ] Provider settings include enable toggle, base URL, version tabs (when applicable), webhook URL, env vars.
- [ ] Dual-port runtime works.
- [ ] Keyboard shortcuts and a11y checks pass.

## Documentation & Process

- [ ] `TRACKING.md` reflects the final state of every prompt.
- [ ] `HANDOFF.md` is up to date and includes final visual sign-off summary.
- [ ] `DECISIONS.md` records all sign-offs for P0/P1 integration fixes.
- [ ] `RISKS.md` lists any deferred bugs and mitigations.
- [ ] `KNOWN_ISSUES.md` lists deferred bugs with rationale and target release.
- [ ] Each fixed bug has a `findings/BXXX-*.md` file.
- [ ] `CHANGELOG.md` has a release-notes snippet for fixed bugs.
- [ ] Root `TRACKING.md` active initiatives table is updated if needed.

## Scope Guardrails

- [ ] No new features or providers were added.
- [ ] No provider plugin schema contract changes.
- [ ] No speculative refactors unrelated to a documented bug.
- [ ] No secrets or real credentials committed.
