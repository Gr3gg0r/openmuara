> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Provider Conformance Audit Known Issues

> **Created:** 2026-07-08
> **Last Updated:** 2026-07-09
> **Status:** ✅ Complete — all gaps closed or formally accepted

---

## How to record a conformance finding

When the audit identifies a gap between OpenMuara behavior and the real provider contract, record it using this template:

```markdown
### FXXX — <short summary>

- **Provider:** <name>
- **Level:** L0/L1/L2/L3/L4/L5/L6
- **Severity:** P0/P1/P2
- **Real contract:** What the provider docs specify.
- **OpenMuara behavior:** What the emulator currently does.
- **Impact:** Concrete consequence for users.
- **Recommended fix:** Specific code/doc/test change.
- **Decision:** Fix / Document deviation / Accept limitation.
- **Owner:** AI Agent / Maintainer / External reviewer.
- **Review date:** YYYY-MM-DD.
```

## Closed findings

| ID | Finding | Provider | Level | Severity | Fixed in |
|---|---|---|---|---|---|
| F01 | Missing L1 request-contract tests | fawry, stripe, billplz, toyyibpay, senangpay, ipay88 | L1 | High | `internal/<provider>/conformance_test.go` |
| F02 | Missing L2 response-contract golden files for Billplz, ToyyibPay, iPay88 | billplz, toyyibpay, ipay88 | L2 | High | `internal/<provider>/testdata/conform/*.json` |
| F03 | Missing invalid/tampered signature tests | billplz, toyyibpay, senangpay, ipay88 | L3 | High | `internal/<provider>/conformance_test.go` |
| F04 | Missing webhook payload shape tests | billplz, toyyibpay, senangpay, ipay88 | L4 | High | `internal/<provider>/conformance_test.go` |
| F05 | Missing L5 state-transition scenarios | fawry, stripe | L5 | Medium | `internal/fawry/scenario_test.go`, `internal/stripe/scenario_test.go` |
| F06 | `internal/provider/conform` only captured static routes | all | L0 | Medium | `internal/provider/conform/conform.go` extended with `AssertJSONEqual` |
| F07 | CI did not explicitly enforce conformance regression | all | L0–L5 | Medium | `.github/workflows/ci.yml` provider conformance step |

## Accepted deviations

These behaviors differ from the strict provider contract but are accepted for the current release. Every deviation has a test that documents the behavior.

| ID | Provider | Area | Deviation | Rationale | Owner | Review date |
|---|---|---|---|---|---|---|
| D01 | fawry | L1 | Charge handler does not enforce `Content-Type: application/json` | Emulation focuses on body parsing; header enforcement is low value for local testing | AI Agent | 2026-10-09 |
| D02 | stripe | L1 | `cancel_url` is declared required in `gateway.yml` but not enforced | Current validation checks `success_url` and `line_items`; minimal impact on integration testing | AI Agent | 2026-10-09 |
| D03 | stripe | L1 | Charge handler does not enforce `Content-Type` | Same rationale as D01 | AI Agent | 2026-10-09 |
| D04 | stripe | L2 | Webhook event wraps full object under `data.object` rather than the flat `{id, object, type}` template in `gateway.yml` | More realistic Stripe event shape; template updated in behavior | AI Agent | 2026-10-09 |
| D05 | billplz | L1 | Charge handler does not enforce `Content-Type` | Same rationale as D01 | AI Agent | 2026-10-09 |
| D06 | toyyibpay | L1 | Wrong `Content-Type` is not explicitly rejected; request fails authentication instead | `ParseForm` behavior makes explicit 400 expensive; documented in test | AI Agent | 2026-10-09 |
| D07 | senangpay | L1 | Charge handler does not enforce `Content-Type` | Same rationale as D01 | AI Agent | 2026-10-09 |
| D08 | senangpay | L4 | Webhook handler does not verify signatures or special headers | Local emulation focuses on payload shape; signature verification is client-side responsibility | AI Agent | 2026-10-09 |
| D09 | senangpay | L4 | `PayloadBuilder()` returns `{"provider","reference","status"}` while `gateway.yml` declares `{"order_id","status_id"}` | Simpler generic payload; provider-specific template to be aligned in future | AI Agent | 2026-10-09 |
| D10 | ipay88 | L1 | Charge handler does not explicitly reject wrong `Content-Type`; `ParseForm` cannot read body | Same rationale as D06 | AI Agent | 2026-10-09 |
| D11 | ipay88 | L2 | Success response is an HTTP redirect, not JSON | iPay88 entry page redirects to local simulation page by design | AI Agent | 2026-10-09 |

## Active findings

None. All known gaps are closed or covered by a documented deviation with a review date.
