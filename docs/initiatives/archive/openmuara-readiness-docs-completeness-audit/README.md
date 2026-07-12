# OpenMuara Readiness — Documentation Completeness Audit

> **Status:** 🟡 Planning  
> **Goal:** Bring OpenMuara documentation to gold-standard OSS completeness, accuracy, and discoverability.  
> **Self-rating (current plan):** 9.5/10 — covers accuracy, provider docs, reference docs, governance, website discoverability, and verification gates.

---

## Scope

This initiative audits and improves every public-facing document: root docs, `docs/`, `runbooks/`, planning indices, GitHub templates, and the Docusaurus website. It does **not** change provider emulation code or API behavior.

## Why now

OpenMuara has strong core docs, but several high-impact gaps undermine trust and discoverability:

- A requirements doc claims Stripe is unimplemented despite it being shipped.
- Two contributing guides give conflicting setup instructions.
- Provider docs are too thin for "integrate without reading code".
- The published website sidebar omits several existing docs.
- Planning indices are stale and link to missing files.
- There is no governance or maintainer-responsibilities doc.

## Audit findings

### Existing docs (high-level)

| Area | Strengths | Gaps |
|---|---|---|
| Root docs | README, CONTRIBUTING, CODE_OF_CONDUCT, SECURITY, CHANGELOG, LICENSE all present | CONTRIBUTING.md diverges from `docs/contributing.md`; no GOVERNANCE.md |
| Product docs | intro, quickstart, architecture, security, operations, webhooks are strong | `docs/cli.md` missing in this branch; provider docs thin; MKP doc inaccurate |
| Provider docs | All providers have a page | Missing runnable examples, signature details, simulation/escape routes, webhook payloads |
| Runbooks | Complete set with index | Minor updates needed as features ship |
| Planning docs | Prompts and tasks exist | Indices stale; broken links |
| Website | Docusaurus wired to docs/runbooks | Sidebar omits install, contributing, contributing-providers, provider-contract, migration, bug-hunt-process |
| GitHub templates | Bug, feature, PR templates present | No docs issue template |

### Top 10 gaps ranked by impact

1. **Inaccurate MKP requirements** — `docs/mkp-billing-requirements.md` says Stripe is not implemented; it is.
2. **Divergent contribution guides** — Root `CONTRIBUTING.md` (Go 1.22+) vs `docs/contributing.md` (Go 1.26+).
3. **Thin provider docs** — No runnable first request, signature formula, or webhook example for most providers.
4. **Missing CLI reference** — No `docs/cli.md` cataloging `muara` commands.
5. **Stale planning indices** — `prompts/INDEX.md` statuses outdated; `tasks/INDEX.md` links to missing migration guide.
6. **No governance doc** — Missing `GOVERNANCE.md` for maintainer roles and decision process.
7. **Incomplete website sidebar** — Several docs are not discoverable on the site.
8. **Stale OpenAPI version** — `info.version: 1.0.0` while codebase has unreleased changes.
9. **Placeholder migration guide** — `docs/migration/openmuara-to-openmuara.md` lacks real upgrade guidance.
10. **Missing docs issue template** — Docs-only issues may use wrong template.

---

## Milestones

| ID | Milestone | Deliverables | Acceptance Criteria |
|---|---|---|---|
| M1 | Accuracy sweep | Fix `docs/mkp-billing-requirements.md`; reconcile `CONTRIBUTING.md` ↔ `docs/contributing.md`; update `prompts/INDEX.md` and `tasks/INDEX.md`; bump `openapi.yaml` version | All inaccurate statements removed; one source of truth for contributing; no broken index links |
| M2 | Provider docs hardening | Expand all `docs/providers/*.md` with runnable examples, signature details, simulation/escape routes, webhook payload snippets, and error tables | Each provider doc contains a copy-paste first request and signature explanation |
| M3 | Reference docs completion | Add `docs/cli.md`; ensure `docs/accessibility.md` and `docs/DOCS_STYLE.md` are present; expand `docs/migration/openmuara-to-openmuara.md` | CLI doc covers all commands; migration guide has version-specific steps |
| M4 | Governance & discoverability | Add `GOVERNANCE.md`; update `website/sidebars.ts`; add docs issue template; link governance/security from README | All existing docs reachable from website; governance doc defines maintainer roles |
| M5 | Verification & handoff | Run markdown lint, link checker, doc-site build; verify every code example against current binary; update CHANGELOG and tracker | `task quality` passes; `npm run build` in `website/` passes; no broken internal links |

---

## Recommendations

### M1 — Accuracy sweep

- **`docs/mkp-billing-requirements.md`**: rewrite the "Current OpenMuara Coverage" table to reflect reality:
  - Stripe → ✅ Implemented (Checkout + PaymentIntents + webhooks)
  - RevenueCat → ⏸️ Deferred to v2
  - Remove internal repo paths or replace with generic placeholders
  - Update endpoint paths from `/v1/stripe/...` to actual `/v1/...` routes
- **`CONTRIBUTING.md` reconciliation**:
  - Make root `CONTRIBUTING.md` the canonical source
  - Update Go version to `1.26+`, Node to `20+`, and mention `task ui:build`
  - Replace `docs/contributing.md` body with a short redirect + link to root file
  - Keep the `.actrc` / fork testing guidance by merging it into root file
- **Planning indices**:
  - Walk `prompts/` and `tasks/` to set correct status emojis
  - Remove or redirect prompts that are now initiatives
  - Fix `tasks/INDEX.md` migration-guide link
- **`docs/openapi.yaml`**:
  - Set `info.version` to match `VERSION` file (e.g., `1.0.0` → current unreleased version)
  - Add a release-runbook reminder to keep this in sync

### M2 — Provider docs hardening

Standardize every `docs/providers/<provider>.md` to include:

1. **Configuration snippet** for `providers.<name>.enabled: true`
2. **First request** — a copy-paste `curl` command that returns a real response
3. **Signature algorithm** — formula + example values showing how to compute it
4. **Simulation / escape routes** — admin pages that mutate state for testing
5. **Webhook payload example** — signed body + headers the consumer will receive
6. **Common errors table** — HTTP status, error code, cause, fix
7. **Standard localhost** — use `127.0.0.1` everywhere

Example template for each provider doc:

```markdown
## Configuration

```yaml
providers:
  <name>:
    enabled: true
    config:
      ...
```

## First request

```bash
curl -X POST http://127.0.0.1:9000/<route> ...
```

## Signature

Formula: ...

## Simulation

- `GET /_admin/<name>/...`
- `POST /_admin/<name>/...`

## Webhook

Payload: ...
Headers: ...
```

### M3 — Reference docs completion

- **`docs/cli.md`**:
  - Auto-generate or manually maintain a table of all `muara` commands
  - Include `init`, `start`, `doctor`, `scenario`, `security`, `webhook`, `audit`, `transaction`, `health`, `version`, `clean`
  - Link each command to its JSON schema in `docs/cli-schemas/`
  - Show global flags (`--config`, `--json`, `--quiet`)
- **`docs/accessibility.md`** and **`docs/DOCS_STYLE.md`**:
  - If missing from this branch, copy/adapt from `dev` or create fresh
  - `DOCS_STYLE.md` should cover tone, code examples, localhost convention, and markdown lint rules
- **`docs/migration/openmuara-to-openmuara.md`**:
  - Add a version compatibility matrix
  - Document any breaking config changes between versions
  - Explain how to back up and restore `.muara/`

### M4 — Governance & discoverability

- **`GOVERNANCE.md`**:
  - Maintainer roles (BDFL/PMC/committers)
  - Decision-making (lazy consensus for routine changes, explicit approval for breaking/provider-contract changes)
  - Conflict resolution
  - Path to maintainership
- **`website/sidebars.ts`**:
  - Add `install` under Getting Started
  - Add `contributing` and `contributing-providers` under Community
  - Add `provider-contract` under Reference
  - Add `migration` category
  - Add `bug-hunt-process` under Reference
- **`.github/ISSUE_TEMPLATE/docs.yml`**:
  - Fields: page URL, what's wrong, suggested change
- **`README.md`**:
  - Add visible "Governance" and "Security" links in the footer section

### M5 — Verification & handoff

- Add a doc-quality job to CI (optional but recommended):
  - `markdownlint-cli2` for style
  - `markdown-link-check` or `lychee` for broken links
  - `cd website && npm run build` to catch Docusaurus breakage
- Manually execute every shell example added in M2 against `bin/muara`
- Update `CHANGELOG.md` under `[Unreleased]`
- Update this tracker and the master backlog

---

## Gold-standard principles

- **One source of truth:** Never duplicate instructions. Link or redirect instead.
- **Runnable examples:** Every provider doc and quickstart path must contain a command that works against the current binary.
- **Multiple personas:** Developer, AI agent, tester, contributor, and maintainer each have a clear entry point.
- **Discoverability:** Every doc is reachable from either the website sidebar, README, or runbook index.
- **Accuracy gates:** Docs are verified with the same rigor as code: lint, link check, and example execution.

---

## Definition of Done

Before this initiative can be marked complete, **all** of the following must be true:

1. Every inaccurate statement identified in the audit has been corrected.
2. There is exactly one source of truth for contributor setup (root `CONTRIBUTING.md`).
3. Every provider doc has a copy-paste first request that returns a real response against `bin/muara`.
4. Every code example uses `127.0.0.1:9000` and a test secret/key that the binary accepts.
5. `docs/cli.md` exists and every command is verified against `muara --help`.
6. `GOVERNANCE.md` is present and defines roles, decision-making, and conflict resolution.
7. The Docusaurus website builds without warnings and every doc is reachable from the sidebar or README.
8. `CHANGELOG.md` has an `[Unreleased]` entry summarizing doc changes.
9. `task quality`, `go test ./...`, and `cd website && npm run build` all pass.
10. This tracker is updated with commit hashes and marked ✅.

---

## Sample deliverables (for execution reference)

### CONTRIBUTING redirect in `docs/contributing.md`

```markdown
# Contributing

Please see the root [`CONTRIBUTING.md`](/CONTRIBUTING.md) for setup,
code style, and submission guidelines.
```

### `GOVERNANCE.md` outline

```markdown
# Governance

## Roles
- BDFL / project lead
- Maintainers
- Committers
- Contributors

## Decision-making
- Routine changes: lazy consensus
- Breaking / provider-contract changes: explicit maintainer approval
- Security issues: private disclosure, then public post-fix

## Conflict resolution
- Discuss in issue/PR
- Escalate to maintainers
- Final call by BDFL

## Becoming a maintainer
- Sustained contributions
- Review quality
- Community trust
```

### Docs issue template (`.github/ISSUE_TEMPLATE/docs.yml`)

```yaml
name: Documentation
labels: [docs]
body:
  - type: input
    attributes:
      label: Page URL
      placeholder: https://openmuara.dev/docs/...
  - type: textarea
    attributes:
      label: What is wrong?
  - type: textarea
    attributes:
      label: Suggested change
```

---

## Anti-patterns to avoid

| Anti-pattern | Why it hurts | What to do instead |
|---|---|---|
| Duplicating setup instructions | One copy goes stale | Link to the single source of truth |
| Using `localhost` inconsistently | Examples fail on some systems | Standardize on `127.0.0.1` |
| Hardcoding real secrets in examples | Security risk and broken copy-paste | Use `sk_test_...` / `test_key` placeholders |
| Writing provider docs from memory | Signature formulas drift | Verify every curl against the running binary |
| Adding docs without sidebar entries | Users can't find them | Update `website/sidebars.ts` in the same commit |
| Skipping the website build | Broken links and Docusaurus errors slip through | Run `npm run build` before commit |

---

## Peer-review checklist

Use this list when reviewing the PR that executes this initiative:

- [ ] No inaccurate provider-status claims remain.
- [ ] Root `CONTRIBUTING.md` is the only file with full setup instructions.
- [ ] Every provider doc has: config, first request, signature, simulation, webhook, errors.
- [ ] `docs/cli.md` matches `muara --help` output.
- [ ] `GOVERNANCE.md` answers: who decides, how to escalate, how to become a maintainer.
- [ ] `website/sidebars.ts` includes every doc changed or added.
- [ ] All examples were run against `bin/muara` and produced the documented output.
- [ ] `CHANGELOG.md` updated.
- [ ] CI docs job passes (if added).

---

## Risks and mitigations

| Risk | Mitigation |
|---|---|
| Provider doc expansion becomes stale as code changes | Add doc-example verification to CI quality gate |
| Reconciling CONTRIBUTING docs breaks existing links | Keep root file as source of truth; make `docs/contributing.md` a redirect |
| Website sidebar churn | Update sidebars in the same PR that adds/removes docs |
| OpenAPI version drift | Add release-runbook step to sync `info.version` with `VERSION` |
| Examples drift from binary behavior | Execute every new example during execution and before review |

---

## Out of scope

- Rewriting the website theme or visual design (separate initiative).
- Provider emulation feature work.
- RevenueCat v2 docs (tracked in `openmuara-v2-revenuecat`).
