# Security Policy

Thank you for helping keep OpenMuara and its users safe.

## Supported versions

| Version | Supported |
|---------|-----------|
| 1.1.x   | ✅ Yes |
| < 1.1   | ❌ No |

Only the latest patch release of a supported minor version receives security updates.

## Reporting a vulnerability

Please report security issues privately so we can fix them before public disclosure.

- **Email:** <security@openmuara.dev> (recommended)
- **GitHub:** Use the [Security Advisory form](https://github.com/openmuara/openmuara/security/advisories/new)

Please include:

- A clear description of the issue.
- Steps to reproduce, or a minimal proof of concept.
- The affected version(s) and commit, if known.
- Your preferred disclosure timeline (default is 90 days).

We aim to acknowledge reports within **48 hours** and provide an initial assessment within **5 business days**.

## Disclosure policy

- We follow a coordinated disclosure process.
- We will work with you to agree on a disclosure timeline.
- We will credit you in the advisory unless you request anonymity.
- If a vulnerability is actively exploited, we may expedite the patch and disclosure.

## Security-related configuration

OpenMuara is local-first by default. For shared or production-like environments, enable the hardening options documented in [`docs/security.md`](../docs/security.md):

- Bind to `127.0.0.1` unless network exposure is required.
- Enable admin authentication and TLS when exposing the server.
- Use the `hardened: true` preset to enable rate limiting and strict security headers.
- Keep provider secrets out of version control; use environment variables.

## Security scanning

The following checks run on every release and pull request:

- `go vet ./...`
- `golangci-lint run`
- `govulncheck ./...`
- `gosec ./...`
- `gitleaks` secret scanning
- `npm audit --production` for the dashboard

## Past security advisories

None yet.
