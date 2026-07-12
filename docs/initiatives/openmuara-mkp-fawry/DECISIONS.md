# OpenMuara MKP Fawry Integration — Decision Log

## D001 — Backward compatibility

**Decision:** All new Fawry fields and statuses are optional. Existing configs,
routes, and payloads keep working unchanged.

**Rationale:** MKP is one consumer; other users rely on the current Fawry
contract.

**Date:** 2026-07-01

---

## D002 — Response delay scope

**Decision:** `fawry.response_delay_ms` delays the outgoing webhook dispatch
only. The escape-page redirect returns immediately.

**Rationale:** Testers should not wait for the redirect; the delay simulates a
slow provider webhook.

**Date:** 2026-07-01

---

## D003 — Status values

**Decision:** Add `canceled` and `expired` as core `engine.TransactionStatus`
values so other providers can reuse them.

**Rationale:** These are generic payment lifecycle states, not Fawry-specific.

**Date:** 2026-07-01
