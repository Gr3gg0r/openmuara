> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# 06 — Security CLI Helpers

## Goal

Add `muara security` commands that make hardening self-service for developers and CI.

## Context

Manually generating bcrypt hashes and self-signed certs is error-prone. CLI helpers reduce friction and ensure consistent output.

## Required Output

1. `muara security hash-password`:
   - Reads password from stdin or `--password` flag.
   - Outputs bcrypt hash.
   - Uses sensible default cost.
2. `muara security gen-cert`:
   - Generates self-signed RSA or ECDSA cert/key pair.
   - Writes to `--cert-out` and `--key-out` paths.
   - Supports `--host` (default `localhost`).
3. `muara security audit`:
   - Loads config and prints security posture.
   - Reports bind address, auth enabled, TLS enabled, rate limit state, hardened mode.
   - Warns on insecure settings (e.g., `0.0.0.0` without auth/TLS).
   - Redacts secrets.
4. Add tests for each command.
5. Update `TRACKING.md`, `DECISIONS.md`, `RISKS.md`, and `HANDOFF.md`.

## Decision Criteria

- Commands work offline.
- `audit` never prints secret values.
- `gen-cert` is suitable for local testing only, not production.

## Quality Gate

- `go test ./...`
- `muara security audit` runs in smoke test.
