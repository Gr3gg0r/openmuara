> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Security Audit Execution Plan

> **Created:** 2026-07-08
> **Last Updated:** 2026-07-08
> **Status:** ⬜ Draft

---

## Goal

Run a gold-standard security audit of OpenMuara and remediate all high/critical findings before public release. The plan is planning-only; no code changes until the initiative is explicitly approved for execution.

## Milestones

| Milestone | Target | Deliverables | Owner |
|-----------|--------|--------------|-------|
| M1 — Discovery & baseline | Day 1–2 | Threat model, scan reports, baseline coverage, populated `KNOWN_ISSUES.md` | AI Agent |
| M2 — Quick wins | Day 3–5 | `SECURITY.md`, Dockerfile non-root user, CI action pinning, checksums in release workflow | AI Agent |
| M3 — Core security hardening | Week 2 | Auth/authorization tests, input-validation audit, webhook negative tests, audit-log improvements | AI Agent |
| M4 — Supply chain & release | Week 2–3 | SBOM generation, image scanning, release signing/attestation | AI Agent |
| M5 — Validation & sign-off | Week 3 | All scans clean, `REVIEW_CHECKLIST.md` complete, human reviewer approval | Human Reviewer |

## Phase sequence

1. **P01 — Threat modeling & asset inventory** (M1)
2. **P02 — Static & dependency analysis** (M1)
3. **P11 — Incident response & disclosure** (M2)
4. **P09 — Container & supply-chain security** (M2)
5. **P10 — CI/CD & release security** (M2)
6. **P03 — Authentication & authorization** (M3)
7. **P04 — Cryptography review** (M3)
8. **P05 — Input validation & web attack surface** (M3)
9. **P06 — Webhook security** (M3)
10. **P07 — Audit logging & PII handling** (M3)
11. **P08 — Configuration & defaults** (M3)
12. **P12 — Remediation & regression tests** (M4–M5)

## RACI (detailed)

| Activity | AI Agent | Human Reviewer | Maintainer |
|---|---|---|---|
| Run scans & triage findings | R | A | C |
| Write threat model | R | A | C |
| Approve accepted-risk list | C | A | R |
| Create `SECURITY.md` | R | A | C |
| Harden Dockerfile | R | A | C |
| Pin CI actions / restrict permissions | R | A | C |
| Implement release signing | R | A | C |
| Add auth/webhook regression tests | R | A | C |
| Review audit-log / PII changes | C | A | R |
| Final sign-off | C | A | R |

## Resource assumptions

- 1 AI agent working in focused sessions.
- 1 human reviewer available for daily async feedback.
- No new paid services required (uses free OSS scanners and GitHub Actions).

## Risk to schedule

| Risk | Mitigation |
|---|---|
| History contains a secret requiring rewrite | Front-load P02; if found, pause and follow `ROLLBACK_PLAN.md` |
| Scanner reports many low-severity noise | Use allowlists and document accepted risks; do not chase zero findings at all costs |
| Container changes break Docker Compose | Run `docker compose up` smoke test in CI |

## Definition of done

- All success metrics in `README.md` are met.
- `REVIEW_CHECKLIST.md` is fully checked.
- Human reviewer approves the final PR.
