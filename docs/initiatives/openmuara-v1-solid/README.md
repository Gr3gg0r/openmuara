> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara v1 Solid

> **Status:** COMPLETE | **Started:** 2026-06-28 | **Completed:** 2026-06-29
> **Scope:** Close regressions and gaps introduced by recent v1 improvements; make the v1 runtime solid for daily use.
> **Owner:** AI Agent (Kimi Code) | **Human Reviewer:** ___________
> **Target Repo:** `<repo-root>/`
> **Product Branch:** `dev`

---

## Initiative Structure

```
docs/initiatives/openmuara-v1-solid/
├── README.md              # This file
├── HOWTO.md               # Decomposition guide for AI
├── PREREQUISITES.md       # Human pre-flight checklist
├── TRACKING.md            # Central execution tracker
├── HANDOFF.md             # Session continuity
├── DECISIONS.md           # Decision log
├── RISKS.md               # Risk register
├── KNOWN_ISSUES.md        # Pre-existing bugs / out-of-scope
├── REFERENCES.md          # Links to specs, runbooks, vendor docs
├── .gitignore             # Ignore screenshots, logs, temp files
│
├── prompts/               # Numbered, self-contained execution prompts
│   ├── _template.md
│   └── ##-verb-domain.md
│
├── tasks/                 # (Optional) Detailed specs — dual-layer
│   ├── _template.md
│   └── ##-verb-domain.md
│
├── findings/              # Research, audit output, analysis
├── runbooks/              # Operational docs
├── screenshots/           # QA evidence (gitignored)
├── qa/                    # Validation artifacts (gitignored)
└── state/                 # Agent state snapshots (gitignored)
```

Planning docs live in `docs/initiatives/openmuara-v1-solid/` in the root repo. Product code commits to the `dev` branch. Do not commit directly to `main`.

---

## Conventions (MANDATORY)

### 1. Read `AGENTS.md` first
Root `AGENTS.md` governs branch rules, quality gates, autonomy boundaries, and code style. This initiative does not repeat every rule.

### 2. Repo & Branch Boundaries
- **Planning docs** (this folder) → commit in root `muara` repo on `dev`.
- **Product code** → commit on `dev` in `<repo-root>/`.
- **Never** commit directly to `main`.
- **Never** mix a planning-doc commit with a product-code commit.

### 3. Prompt Isolation
- Each prompt MUST be self-contained. An AI reading only `prompts/05-*.md` must be able to execute it.
- If a prompt needs spec context, use the dual-layer pattern: `tasks/05-*.md` = spec, `prompts/05-*.md` = execution.
- Prompts MUST NOT reference other prompts ("see Prompt 03" is forbidden). Reference task specs or this README instead.

### 4. File Naming
- Prompts: `##-verb-domain-subject.md`
- Tasks: `##-verb-domain-subject.md` (mirror prompt when dual-layer)
- Findings: `##-short-descriptive-name.md`
- Runbooks: `kebab-case-topic.md`

### 5. Status Discipline
- Update `TRACKING.md` after EVERY step: status, commit hash, decisions, notes.
- Update `HANDOFF.md` BEFORE exiting the session.
- Log decisions in `DECISIONS.md` as they happen.
- **AFK mode:** Every non-trivial decision MUST be auto-logged in `DECISIONS.md` Auto-Decisions table before finishing the prompt. "None" is a valid entry.

### 6. Screenshots
- Save to `screenshots/` or `##/screenshots/`.
- Name descriptively: `dashboard-webhooks-empty.png`, `openapi-readyz.png`.
- Screenshots and QA artifacts are gitignored. Do not commit them.

### 7. Branch & Commit Rules
- Work happens on `dev`.
- One logical change per commit.
- Commit format: `type(scope): imperative description`
  - Types: `feat`, `fix`, `chore`, `docs`, `refactor`, `test`, `ci`, `perf`
  - Examples:
    - `fix(ui): handle paginated admin responses`
    - `docs(openapi): sync spec with readyz and pagination`

### 8. Quality Gates (Go)
| Gate | Command | Target |
|------|---------|--------|
| Build | `go build ./...` | Compiles |
| Test | `go test ./...` | All pass |
| Race | `go test -race ./...` | All pass |
| Lint | `golangci-lint run` | Zero issues |
| Vet | `go vet ./...` | Clean |
| Smoke | `./scripts/smoke-test.sh` | Passes |

### 9. Regression Protection
These flows must NEVER break. Verify after any change that touches:
- **Universal payment API:** `POST /v1/pay`, `GET /v1/pay/{ref}`, `POST /v1/refund/{ref}`
- **Providers:** Fawry charge/escape/webhook, Stripe checkout/simulation, SenangPay charge/callback/webhook
- **Outgoing webhooks:** dispatch, replay, admin list
- **Admin dashboard:** `/_admin` loads and reflects state

---

## Completion Criteria

This initiative is DONE when:
- [ ] All steps in `TRACKING.md` are ✅ DONE.
- [ ] All quality gates in `TRACKING.md` show ✅ Pass.
- [ ] `HANDOFF.md` is updated with final state.
- [ ] `DECISIONS.md` captures all major decisions.
- [ ] `RISKS.md` shows all risks as mitigated or accepted.
- [ ] `README.md` header updated to `Status: COMPLETE | Date: YYYY-MM-DD`.
- [ ] `KNOWN_ISSUES.md` and `REFERENCES.md` are up to date.
- [ ] Initiative moved to `docs/initiatives/archive/done/openmuara-v1-solid/`.
