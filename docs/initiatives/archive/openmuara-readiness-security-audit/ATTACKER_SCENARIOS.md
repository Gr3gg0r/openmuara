> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Attacker Scenarios

> **Created:** 2026-07-08
> **Last Updated:** 2026-07-09
> **Status:** ✅ Reviewed

---

Use these scenarios to design regression tests and acceptance criteria. Each scenario states the attacker goal, the expected defense, and the test to verify it.

## AS-01 — Unauthorized dashboard access

- **Attacker goal:** View transactions or replay webhooks without credentials.
- **Prerequisites:** Server bound to `0.0.0.0` with admin auth enabled.
- **Expected defense:** `/_admin/*` returns `401 Unauthorized` without valid Basic Auth or bearer token.
- **Test:** `TestSecurityIntegrationAdminRequiresAuth` pattern; extend to cover viewer-vs-admin differentiation.

## AS-02 — Viewer escalates to admin

- **Attacker goal:** Use a viewer account to change config or replay webhooks.
- **Prerequisites:** Viewer auth enabled.
- **Expected defense:** Viewer requests to mutating admin endpoints return `403 Forbidden`.
- **Test:** Add table-driven tests for every mutating admin route using viewer credentials.

## AS-03 — Webhook replay without signature

- **Attacker goal:** Inject a fake provider webhook to mutate ledger state.
- **Prerequisites:** Provider webhook endpoint configured.
- **Expected defense:** Webhook endpoint rejects requests with missing or invalid signature.
- **Test:** Negative signature tests for each provider; verify ledger is unchanged.

## AS-04 — Signature bypass via algorithm confusion or empty signature

- **Attacker goal:** Bypass HMAC verification by sending empty, null, or malformed signature fields.
- **Prerequisites:** Provider charge/webhook endpoint.
- **Expected defense:** Signature is required and compared in constant time; malformed input returns an error.
- **Test:** Fuzz and table tests for empty, too-short, and tampered signatures.

## AS-05 — XSS via webhook payload rendered in dashboard

- **Attacker goal:** Execute JavaScript in an admin’s browser by injecting a payload into a webhook or transaction field.
- **Prerequisites:** Webhook payload saved and rendered in dashboard.
- **Expected defense:** CSP blocks inline scripts; UI escapes rendered JSON; no `dangerouslySetInnerHTML` on untrusted data.
- **Test:** Add a unit test that renders a payload containing `<script>alert(1)</script>` and assert no script execution.

## AS-06 — CSRF against admin mutation

- **Attacker goal:** Trick an authenticated admin into performing a state-changing action.
- **Prerequisites:** Admin authenticated, CSRF enabled.
- **Expected defense:** `X-CSRF-Token` header or `csrf_token` form field must match the double-submit cookie.
- **Test:** `internal/server/csrf_test.go` already exists; extend to cover all admin POST/PUT/DELETE routes.

## AS-07 — SSRF via webhook URL or provider return URL

- **Attacker goal:** Make the server call an internal service (e.g., `http://169.254.169.254/` or `file://`).
- **Prerequisites:** Admin-configured `webhook.url` or `webhook.targets`.
- **Expected defense:** `httputil.ValidateWebhookURL` rejects non-HTTP(S) schemes always; rejects loopback, link-local, and private IPs when `hardened: true`.
- **Test:** `internal/httputil/url_test.go` and `internal/config/validation_test.go` cover valid/invalid schemes and private-IP rejection in hardened mode.

## AS-08 — SQL injection via query parameters

- **Attacker goal:** Extract or modify ledger data through unsanitized query inputs.
- **Prerequisites:** Admin API with search/filter parameters.
- **Expected defense:** All SQL queries use parameterized statements or ORM equivalents.
- **Test:** Fuzz search parameters with SQL metacharacters; assert no error/leakage.

## AS-09 — Path traversal via file endpoints

- **Attacker goal:** Read arbitrary files from the server.
- **Prerequisites:** Any endpoint that accepts a file path.
- **Expected defense:** Paths are validated and joined with a safe base directory; `..` rejected.
- **Test:** Attempt `../../../etc/passwd` style paths on any file-serving endpoint.

## AS-10 — Rate-limit bypass or DoS

- **Attacker goal:** Exhaust the rate limiter or crash the server with oversized payloads.
- **Prerequisites:** Rate limiting enabled.
- **Expected defense:** Per-IP token bucket rejects excess requests; body size limited to 1 MiB; map bounded.
- **Test:** Burst requests beyond threshold; send >1 MiB body; assert `429` and no panic.

## AS-11 — Secret leakage in logs or error messages

- **Attacker goal:** Obtain admin tokens or provider secrets from logs or error responses.
- **Prerequisites:** Any flow that logs headers, query strings, or form data.
- **Expected defense:** Secrets redacted before logging; generic error messages to clients.
- **Test:** Trigger errors with secret-bearing headers; inspect logs and responses.

## AS-12 — Container escape via root user

- **Attacker goal:** Exploit a vulnerability in the container to gain host root.
- **Prerequisites:** Container runs as root.
- **Expected defense:** Container runs as non-root user; read-only root fs where feasible.
- **Test:** `docker run --entrypoint whoami openmuara:latest` returns a non-root user.

## Scenario-to-phase mapping

| Scenario | Primary phase | Test file suggestion |
|---|---|---|
| AS-01 | P03 | `internal/server/auth_test.go` |
| AS-02 | P03 | `internal/server/auth_test.go` |
| AS-03 | P06 | `internal/<provider>/webhook_test.go` |
| AS-04 | P04 / P06 | `internal/<provider>/signature_fuzz_test.go` |
| AS-05 | P05 | `web/dashboard/tests/WebhookDetail.test.tsx` |
| AS-06 | P03 / P08 | `internal/server/csrf_test.go` |
| AS-07 | P05 | `internal/webhook/dispatcher_test.go` |
| AS-08 | P05 | `internal/engine/sqlite_test.go` |
| AS-09 | P05 | `internal/server/*_test.go` |
| AS-10 | P03 / P08 | `internal/server/ratelimit_test.go` |
| AS-11 | P07 | `internal/audit/audit_test.go` |
| AS-12 | P09 | CI container scan + `docker run` test |
