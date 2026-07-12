# Prompt 19 — Post-Launch Monitoring

## Goal
Prepare runbooks and monitoring guidance for operating OpenMuara.

## Acceptance Criteria
- [ ] `runbooks/on-call.md` with common alerts
- [ ] `runbooks/debugging.md` — how to inspect state, replay webhooks, check audit log
- [ ] `KNOWN_ISSUES.md` documented
- [ ] `RISKS.md` updated with mitigations
- [ ] Alerting rules example for Prometheus
- [ ] Log aggregation guidance

## Files to Create/Change
- `runbooks/on-call.md`
- `runbooks/debugging.md`
- `KNOWN_ISSUES.md`
- `RISKS.md`
- `docs/operations.md`

## Response Shape
Return:
1. Runbook index
2. Key metrics and alerts
3. Known issues summary

## Test Notes
- Review docs for completeness
- Verify alerting rule syntax
