---
id: on-call
title: On-Call Runbook — OpenMuara
---

# On-Call Runbook — OpenMuara

Quick reference for responding to OpenMuara alerts and operational issues.

---

## On-call checklist

When paged or investigating an issue:

1. Confirm the OpenMuara process is running: `GET /healthz`.
2. Confirm it is ready: `GET /readyz`.
3. Check recent logs for errors or panics.
4. Check `GET /metrics` for request/error/webhook rates.
5. Check the audit log: `muara audit list --limit 100`.
6. If a provider flow is failing, inspect the relevant transaction in `/_admin`.
7. If webhooks are failing, inspect attempts and replay after fixing the consumer.
8. Escalate if data integrity is in question (SQLite corruption, unexpected refunds, etc.).

---

## Health endpoints

```bash
curl -s http://127.0.0.1:9000/healthz
curl -s http://127.0.0.1:9000/readyz
```

`/readyz` returns `503` if a provider failed to initialize. Common causes:

- Missing or invalid provider config (`merchant_code`, `secret_key`, etc.).
- SQLite database is locked or missing write permissions.
- A provider plugin returned an error from `Init`.

---

## Key metrics

| Metric | What to watch |
|--------|---------------|
| `openmuara_requests_total` | Traffic volume and 5xx rate |
| `openmuara_request_duration_seconds` | Latency distribution |
| `openmuara_webhook_attempts_total` | Delivery success/failure by provider |
| `openmuara_transactions_total` | Transaction volume by provider and status |

---

## Common alerts

### OpenMuaraDown / process not responding

**Symptoms:** `GET /healthz` times out or returns no response.

**Steps:**

1. Check the process status (`docker ps`, `systemctl status openmuara`, or `ps`).
2. Look at logs for the last panic or fatal error.
3. Restart the process.
4. Verify `/healthz` and `/readyz` return `200`.

### High 5xx rate

**Symptoms:** `openmuara_requests_total{status=~"5.."}` is elevated.

**Steps:**

1. Filter by `path` to identify the failing endpoint.
2. Read logs for the matching `trace_id`.
3. Check `/readyz` for provider initialization failures.
4. If SQLite returns `database is locked`, see [SQLite contention](#sqlite-database-is-locked) below.

### Webhook delivery failures

**Symptoms:** `openmuara_webhook_attempts_total{status="failed"}` is increasing.

**Steps:**

1. Inspect attempts: `muara webhook list --provider <name>`.
2. Verify the consumer URL is reachable from OpenMuara.
3. Check consumer logs for 4xx/5xx responses.
4. Fix the consumer, then replay failed attempts from `/_admin` or with `muara webhook replay <id>`.

### SQLite database is locked

**Symptoms:** `500` responses, logs mention `database is locked`.

**Steps:**

1. Reduce concurrent write load. OpenMuara uses a single shared SQLite connection.
2. Ensure the database file is on a local filesystem, not a network share.
3. Restart the process to clear stale locks.
4. If the issue persists, see [`runbooks/debugging.md`](debugging.md#database-locked-or-slow).

### Audit log growing rapidly

**Symptoms:** Disk usage in `.muara/data/ledger.db` grows quickly.

**Steps:**

1. Check event volume: `muara audit list --limit 1000`.
2. Increase `log.level` only when needed; audit entries are independent of logs.
3. Archive or truncate old audit records manually if disk space is critical.
4. Long-term: schedule periodic pruning or move audit data to a dedicated store.

### Unauthorized admin access attempts

**Symptoms:** Repeated `401` responses on `/_admin/*`, or `security.auth.failure` audit events.

**Steps:**

1. Confirm `admin.enabled` is set and credentials are strong.
2. Check the source IPs; rate limiting will block abusive IPs with `429`.
3. Verify the server is not bound to `0.0.0.0` without TLS in an untrusted network.
4. Rotate `admin.token` and regenerate `admin.password_hash` if compromise is suspected.

### Rate limiting triggered

**Symptoms:** Clients receive `429 Too Many Requests` on admin or provider endpoints.

**Steps:**

1. Identify whether the traffic is legitimate load testing or abuse.
2. Adjust `rate_limit.requests_per_minute` if the limit is too low for your test profile.
3. Disable rate limiting only in local trusted environments; keep it enabled in CI/shared networks.

### TLS configuration errors

**Symptoms:** Server fails to start with TLS-related errors, or HTTPS clients reject the certificate.

**Steps:**

1. Verify `server.tls_cert` and `server.tls_key` are both set and readable.
2. Confirm the certificate is valid for the host name clients use.
3. For local testing, generate a self-signed cert with `muara security gen-cert`.

---

## Escalation

Escalate to the team lead or project owner when:

- Data loss or corruption is suspected.
- A fix requires changing provider emulation contracts or signature algorithms.
- A release rollback is needed.
- The issue spans multiple OpenMuara instances or shared infrastructure.

Before escalating, include:

- Time the issue started.
- `/healthz` and `/readyz` output.
- Relevant log excerpts (with `trace_id`).
- Recent metric snapshots from `/metrics`.
- Recent audit events from `muara audit list`.
