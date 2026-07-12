# OpenMuara Billplz — Decision Log

| ID | Decision | Status | Date | Rationale |
|----|----------|--------|------|-----------|
| D001 | Add Billplz as a first-class OpenMuara provider named `billplz`. | ✅ | 2026-07-01 | Malaysian developers commonly use Billplz for FPX/card payments; OpenMuara needs to emulate it for local testing. |
| D002 | Billplz `x_signature` is a form/query parameter, not an HTTP header. | ✅ | 2026-07-01 | Matches real Billplz v3 behavior: `x_signature` is included in the callback form body and redirect query string. |
| D003 | Callback is a server-side POST to `callback_url`; redirect is a browser GET to `redirect_url`. | ✅ | 2026-07-01 | Aligns with real Billplz terminology and flow. |
| D004 | Webhook payload is a flat form-urlencoded Bill object, not a JSON event envelope. | ✅ | 2026-07-01 | Real Billplz v3 callbacks are form-encoded Bill objects, not JSON events like `bill.paid`. |
| D005 | Bill response omits `currency`; amount is in sen (integer). | ✅ | 2026-07-01 | Real Billplz v3 bills are implicitly MYR and return amount as integer sen. |
