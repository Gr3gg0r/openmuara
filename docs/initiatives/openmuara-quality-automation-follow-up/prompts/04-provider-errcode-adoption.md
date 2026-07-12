> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# P04 — Provider errcode Adoption

> **Initiative:** OpenMuara Quality Automation Follow-Up
> **Depends on:** P01
> **Target files:** `internal/errcode/`, `internal/fawry/`, `internal/stripe/`, `internal/senangpay/`, `internal/ipay88/`, `internal/billplz/`, `internal/toyyibpay/`, `internal/api/`, `internal/server/`
> **Status:** ⬜

## Goal

Adopt the existing `internal/errcode` taxonomy across all provider packages and API error responses without changing public error message text.

## Tasks

- [ ] Audit each provider package for error-return paths (signature, config, validation, transaction lookup).
- [ ] Add `errcode` wrapping to the highest-value paths first: signature mismatch/missing, config missing/invalid, provider disabled, transaction not found/duplicate.
- [ ] Ensure API error responses include the code where they already include a message.
- [ ] Add or update tests that assert the code is present.
- [ ] Verify no existing error message text changes unless user sign-off is recorded.

## Acceptance Criteria

- [ ] Every provider package imports and uses `internal/errcode`.
- [ ] Signature, config, and transaction errors have stable codes.
- [ ] Existing tests still pass; new tests verify codes.
- [ ] No provider behavior changes (only error metadata changes).

## Completion Checklist

- [ ] Provider packages wrapped and tested.
- [ ] `TRACKING.md` updated with status `✅`, commit hash, and packages touched.
- [ ] `HANDOFF.md` updated with what was done and P05 next steps.
- [ ] `DECISIONS.md` updated if any public error message had to change.

## Quality Gates

Run the full suite before committing:

```bash
go build ./...
go test ./...
go test -race ./...
go vet ./...
golangci-lint run
cd web/dashboard && npm run test:ci
cd web/dashboard && npm run build
node web/dashboard/scripts/check-bundle-size.js
cd web/dashboard && node scripts/a11y-contrast-check.js
cd web/dashboard && npm run test:visual-baseline
```

## Notes

- This is an additive change. If a public error message must change to include a code, request user sign-off in `DECISIONS.md`.
- Keep commits per provider package for easier review.
