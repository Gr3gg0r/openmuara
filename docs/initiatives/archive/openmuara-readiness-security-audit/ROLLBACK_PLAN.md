> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Security Rollback & Incident Response Plan

> **Created:** 2026-07-08
> **Last Updated:** 2026-07-08
> **Status:** ⬜ Draft

---

This document describes how to respond if the security audit discovers a leaked secret, a critical vulnerability, or another security incident before or after publication.

## Scope

- Leaked credentials in source code, git history, CI logs, or artifacts.
- High/critical vulnerabilities discovered in dependencies or OpenMuara code.
- Compromised release artifacts or signing keys.
- Reports from external security researchers.

## Response team

| Role | Responsibility |
|---|---|
| Human Reviewer | Decides severity, approves response, communicates externally |
| Maintainer | Rotates secrets, cuts fixed release, updates `SECURITY.md` |
| AI Agent | Implements code fixes, adds regression tests, updates docs |

## Leaked secret response

1. **Contain** — revoke or rotate the exposed credential immediately.
2. **Assess** — determine what the secret could access and whether it was exploited.
3. **Remove from future code** — ensure no remaining references in working tree.
4. **History decision** —
   - If the repository is still **private**: rewrite history with `git filter-repo` or BFG and force-push.
   - If the repository is **public**: do **not** rewrite public history; rely on rotation and disclosure.
5. **Verify** — re-run `gitleaks` and `trufflehog` to confirm no other leaks.
6. **Document** — add a private incident note; update `SECURITY.md` if user impact exists.

## Vulnerability response

1. **Triage** — reproduce and score severity (CVSS-like: Low / Medium / High / Critical).
2. **Fix** — implement the smallest safe fix on a private branch if the repo is public and the issue is exploitable.
3. **Test** — add a regression test and run the full quality matrix.
4. **Release** — cut a patch release with clear release notes.
5. **Disclose** — publish a security advisory via GitHub Security Advisories and email the disclosure contact.

## Compromised release artifact

1. **Yank or deprecate** the affected release asset.
2. **Investigate** the CI/build pipeline for compromise.
3. **Re-build** from a clean, pinned CI run.
4. **Re-sign** with a trusted key and publish new checksums.
5. **Notify** users via `SECURITY.md` contact, release notes, and GitHub advisory.

## External researcher reports

1. Acknowledge receipt within 48 hours.
2. Validate the finding within 5 business days.
3. Coordinate disclosure timeline with the researcher (default 90 days).
4. Credit the researcher in the advisory unless they request anonymity.

## Communication templates

### Security advisory (GitHub)

```
Title: OpenMuara <version> — <short summary>
Severity: <Low|Medium|High|Critical>
Affected versions: <x.y.z>
Patched versions: <x.y.z>
Description: <what the issue is and how it could be exploited>
Mitigation: <how users can protect themselves before upgrading>
Credits: <researcher name> (optional)
```

### User notification email

```
Subject: Security update for OpenMuara <version>

We have released OpenMuara <patched version> to address a <severity> security issue.
Affected versions: <range>
Recommended action: Upgrade to <patched version>.
Details: <link to GitHub Security Advisory>
```

## Post-incident review

After any incident:

- Update `RISKS.md` and `DECISIONS.md` if new risks emerged.
- Add a regression test for the vulnerability class.
- Review whether `REVIEW_CHECKLIST.md` needs new items.
