# OpenMuara ToyyibPay — Decision Log

| ID | Decision | Status | Date | Rationale |
|----|----------|--------|------|-----------|
| D001 | Add ToyyibPay as a first-class OpenMuara provider named `toyyibpay`. | ✅ | 2026-07-01 | ToyyibPay is a popular Malaysian gateway for small merchants; OpenMuara needs to emulate it for local testing. |
| D002 | ToyyibPay endpoints accept form-encoded request bodies, not JSON. | ✅ | 2026-07-01 | Matches real ToyyibPay API behavior. |
| D003 | ToyyibPay callback uses MD5 hash, not an HTTP header. | ✅ | 2026-07-01 | Real ToyyibPay callback includes `hash = MD5(userSecretKey + status + order_id + refno + "ok")` as a form field. |
| D004 | Browser return URL is a GET to `billReturnUrl`; server callback is a POST to `billCallbackUrl`. | ✅ | 2026-07-01 | Aligns with real ToyyibPay flow. |
| D005 | Endpoint name is `getCategoryDetails`, not `getCategory`. | ✅ | 2026-07-01 | Matches real ToyyibPay API endpoint naming. |
