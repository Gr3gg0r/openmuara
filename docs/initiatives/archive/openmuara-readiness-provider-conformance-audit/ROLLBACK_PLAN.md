> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Provider Conformance Audit Rollback Plan

> **Created:** 2026-07-09
> **Last Updated:** 2026-07-09
> **Status:** ✅ Complete — rollback plan in place and aligned with CI changes.

---

This plan is active. Use it if a conformance change causes regression, CI instability, or incorrect emulation.

## 1. Conformance test regression

**Scenario:** A new test fails in CI after a provider code change.

**Response:**
1. Determine if the failure is due to:
   - A real bug in the provider code.
   - An incorrect or outdated golden file.
   - A flaky test fixture.
2. If it is a bug, fix the provider code and add a regression test.
3. If the golden file is outdated and the new behavior is correct, run `UPDATE_GOLDEN=1 go test ./internal/provider/conform/...` and commit the updated golden file with a clear explanation.
4. If the test is flaky, move it to a separate job or mark it as helper-only while investigating.

## 2. Wrong provider emulation merged

**Scenario:** A PR changes provider behavior to match an incorrect interpretation of the real contract.

**Response:**
1. Revert the PR if it is safe to do so.
2. Add the incorrect behavior to `KNOWN_ISSUES.md` if the revert cannot happen immediately.
3. Re-map the provider contract and write the correct test.
4. Update `docs/providers/<provider>.md` to reflect the correct behavior.

## 3. Golden-file drift without explanation

**Scenario:** A golden file changes in a PR but the PR description does not explain the contract change.

**Response:**
1. Request an explanation in the PR review.
2. Do not merge until the change is documented in `KNOWN_ISSUES.md` or provider docs.
3. If the change is accidental, regenerate golden files from `main`.

## 4. External review feedback contradicts current behavior

**Scenario:** The Fawry team reports that OpenMuara does not match the real Fawry contract.

**Response:**
1. Thank the reviewer and log the feedback in `KNOWN_ISSUES.md`.
2. Create a fix branch with a new conformance test that captures the corrected behavior.
3. Update golden files and docs.
4. Send the fix back to the reviewer for confirmation if possible.

## 5. CI performance regression

**Scenario:** Conformance tests slow CI significantly.

**Response:**
1. Profile the slow tests.
2. Move heavy scenario tests to a dedicated `provider-conformance` job that runs in parallel.
3. Keep the lightweight `internal/provider/conform` tests in the `unit` job.

## Communication template

For significant regressions, open a GitHub issue with:

```markdown
## Conformance regression: <provider> <area>

- **Introduced in:** commit/PR
- **Provider:** 
- **Level:** L1/L2/L3/L4/L5
- **Expected behavior:** 
- **Actual behavior:** 
- **Impact:** 
- **Proposed fix:** 
```
