> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Provider Conformance Audit Review Checklist

> **Created:** 2026-07-09
> **Last Updated:** 2026-07-09
> **Status:** ✅ Complete — all checklist items verified; external review request deferred as follow-up.

---

Use this checklist to sign off the provider conformance audit initiative.

## Framework

- [x] L0–L6 maturity model is documented and understandable.
- [x] `internal/provider/conform` supports behavior snapshots (request/response) in addition to route snapshots.
- [x] Golden-file update workflow (`UPDATE_GOLDEN=1`) is documented.

## Provider coverage

- [x] Fawry: L1–L4 tests exist and pass.
- [x] Stripe: L1–L4 tests exist and pass.
- [x] Billplz: L1–L4 tests exist and pass.
- [x] ToyyibPay: L1–L4 tests exist and pass.
- [x] SenangPay: L1–L4 tests exist and pass.
- [x] iPay88: L1–L4 tests exist and pass.
- [x] Default provider: L0–L2 covered as reference implementation.

## Documentation

- [x] Every P0 provider doc (`docs/providers/<provider>.md`) states the emulated version.
- [x] Every P0 provider doc has a "Known limitations" or "Deviations" section.
- [x] `KNOWN_ISSUES.md` contains all confirmed gaps and accepted deviations.
- [x] No undocumented deviation from real provider behavior remains.

## Signature & webhook

- [x] Every P0 provider has valid-signature tests.
- [x] Every P0 provider has invalid/tampered-signature tests.
- [x] Every P0 provider has webhook payload tests.
- [x] Every P0 provider has webhook signature-header tests.

## External validation

- [ ] Review request sent to Fawry team.
- [x] Review request template prepared; sending deferred and tracked.

## CI & quality gates

- [x] Conformance regression test runs in CI.
- [x] `go build ./...` passes.
- [x] `go test ./...` passes.
- [x] `go vet ./...` passes.
- [x] `golangci-lint run` passes.
- [x] `scripts/check-coverage.sh 81` passes.
- [x] `scripts/check-coverage-per-package.sh` passes.
- [x] `npm run typecheck` passes.
- [x] `npm run test:ci` passes.

## Sign-off

| Role | Name | Date | Signature |
|---|---|---|---|
| AI Agent | Kimi Code | | |
| Human Reviewer | | | |
| Maintainer | | | |
