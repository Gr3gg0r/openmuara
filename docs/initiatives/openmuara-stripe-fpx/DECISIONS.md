# OpenMuara Stripe FPX — Decision Log

| ID | Decision | Status | Date | Rationale |
|----|----------|--------|------|-----------|
| D001 | Scope FPX as a dedicated charge + escape flow within the Stripe provider, modeled after Fawry, not as a Stripe Checkout Session payment method. | ✅ | 2026-06-30 | The user explicitly requested the Fawry charge + escape pattern. It is simpler to implement, test, and document than extending Checkout Sessions with FPX-specific redirect semantics. |
| D002 | Supersede the custom `/v1/stripe/fpx/*` and `/v1/stripe/card/*` routes with Stripe's real Checkout Sessions and PaymentIntents APIs. | ❄️ | 2026-07-01 | Custom routes break Stripe SDK parity. Real APIs let developers point the official Stripe SDK at OpenMuara and switch to production by changing only the base URL and API key. |
| D003 | Archive this initiative rather than delete it. | ❄️ | 2026-07-03 | Preserve the decision record, implementation history, and lessons learned for future provider emulation work. |
