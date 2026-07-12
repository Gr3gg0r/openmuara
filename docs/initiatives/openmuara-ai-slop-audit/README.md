> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara AI Slop Audit

> **Status:** ⬜ Draft | **Started:** 2026-07-08
> **Scope:** Identify and remove AI-generated/generic/default-looking patterns from the OpenMuara product: dashboard UI/UX, provider metadata, documentation, code, test data, examples, and website.
> **Owner:** AI Agent (Kimi Code) | **Human Reviewer:** ___________
> **Target Repo:** `<repo-root>/`
> **Product Branch:** `feat/ai-slop-audit` (to be created when work starts)

---

## Why this matters

AI slop makes a product feel unfinished, interchangeable, and harder to trust. For OpenMuara — a local payment emulator that developers rely on for realistic test infrastructure — slop undermines confidence in the very thing it emulates. This initiative records findings so they can be fixed intentionally rather than chipped away ad-hoc.

## Initiative structure

```
docs/initiatives/openmuara-ai-slop-audit/
├── README.md              # This file
├── TRACKING.md            # Central execution tracker
├── HANDOFF.md             # Session continuity
├── DECISIONS.md           # Decision log
├── RISKS.md               # Risk register
└── KNOWN_ISSUES.md        # Catalog of slop findings
```

## Audit areas

1. **Dashboard UI & microcopy** — generic icons, placeholder labels, redundant badges, default-looking tables, inconsistent spacing.
2. **Provider metadata** — `gateway.yml` descriptions, repeated "Provider configuration" copy, inconsistent naming/casing.
3. **Documentation & prompts** — vague buzzwords, circular reasoning, placeholder TODOs, generic README sections.
4. **Code-level patterns** — `any`/`interface{}` overuse, copy-paste provider boilerplate, verbose/comment-only code, dead options.
5. **Test / seed data** — synthetic-looking transactions, repetitive references, unrealistic webhook payloads.
6. **Examples & website** — placeholder products, generic marketing copy, boilerplate Docusaurus content.

## Success criteria

- Every finding in `KNOWN_ISSUES.md` is either fixed or explicitly deferred with a rationale.
- Dashboard design audit score improves from current **B-** toward **B+** or higher.
- AI Slop Score remains **A** or improves.
- All quality gates pass after each fix batch.
