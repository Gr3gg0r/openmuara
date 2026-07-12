# OpenMuara iPay88 — Decision Log

| ID | Decision | Status | Date | Rationale |
|----|----------|--------|------|-----------|
| D001 | Add iPay88 as a first-class OpenMuara provider named `ipay88`. | ✅ | 2026-07-01 | iPay88 is a long-standing Southeast Asian gateway; OpenMuara needs to emulate it for local testing. |
| D002 | Target the Malaysia classic form-based ePayment API (`/ePayment/entry.asp`), not the Indonesia JSON API. | ✅ | 2026-07-01 | The classic form API is the most widely referenced iPay88 integration in Malaysia. |
| D003 | Request/response/backend bodies are form-encoded, not JSON. | ✅ | 2026-07-01 | Matches real iPay88 classic ePayment behavior. |
| D004 | Callback status codes are `1` success, `0` fail, `6` pending; requery uses `00` success. | ✅ | 2026-07-01 | Matches real iPay88 status code conventions. |
| D005 | Backend post must be acknowledged with plain text `RECEIVEOK`. | ✅ | 2026-07-01 | Real iPay88 expects this acknowledgement and retries otherwise. |
| D006 | Amount is normalized (separators stripped) before SHA256 signature. | ✅ | 2026-07-01 | Matches real iPay88 signature algorithm. |
| D007 | Common PaymentId mapping for the local payment page: `1` credit/debit card, `2` FPX, `33` Touch 'n Go eWallet, `34` Boost, `35` GrabPay. | ✅ | 2026-07-01 | Provides a simplified but representative selector for local testing; real iPay88 merchant accounts may use provider-specific codes. |
