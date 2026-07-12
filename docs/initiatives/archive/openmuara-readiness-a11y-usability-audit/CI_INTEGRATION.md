> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Accessibility & Usability Audit CI Integration

> **Created:** 2026-07-09
> **Last Updated:** 2026-07-09
> **Status:** ✅ Complete — CI enforcement wired

---

This document describes the CI changes that enforce the accessibility and usability audit outcomes.

## 1. A11y npm scripts in `web/dashboard/package.json`

Implemented scripts:

```json
{
  "scripts": {
    "test:a11y": "playwright test e2e/dashboard-a11y.spec.ts",
    "test:a11y:contrast": "node scripts/a11y-contrast-check.js",
    "lint:a11y": "eslint src/ --max-warnings 0",
    "typecheck": "tsc --noEmit",
    "test:ci": "vitest run"
  }
}
```

## 2. Dependencies added

```bash
cd web/dashboard
npm install --save-dev vitest-axe @axe-core/playwright eslint-plugin-jsx-a11y eslint @typescript-eslint/parser
```

Runtime dependencies were pinned in `package-lock.json` and committed.

## 3. CI gate in `.github/workflows/ci.yml`

A dedicated JSX-a11y lint step was added to the existing `ui-test` job:

```yaml
      - name: Lint JSX accessibility
        run: cd web/dashboard && npm run lint:a11y
```

The `ui-e2e` job already runs the full Playwright suite (`npm run test:e2e`), which includes `e2e/dashboard-a11y.spec.ts`, and the contrast regression check (`npm run test:a11y:contrast`). The `ui-build` job runs `npm run build`, which includes TypeScript type checking.

Resulting a11y enforcement in CI:

| Check | CI job | Command |
|---|---|---|
| Type check | `ui-build` | `npm run build` |
| Unit + component tests | `ui-test` | `npm run test:ci` |
| JSX a11y lint | `ui-test` | `npm run lint:a11y` |
| End-to-end a11y scan | `ui-e2e` | `npm run test:e2e` (includes `dashboard-a11y.spec.ts`) |
| Contrast regression | `ui-e2e` | `npm run test:a11y:contrast` |

## 4. Baseline protection

- New critical or serious axe-core violations fail the `ui-e2e` job.
- New JSX-a11y lint warnings fail the `ui-test` job (`--max-warnings 0`).
- If a violation is accepted, add it to `KNOWN_ISSUES.md` with rationale and mark the test exception with a documented comment.
- Do not suppress violations without documentation.

## 5. Local acceptance commands

```bash
cd web/dashboard
npm run typecheck
npm run lint:a11y
npm run test:ci
npm run test:a11y
npm run test:a11y:contrast
```

## 6. Rollback / exception handling

- If `vitest-axe` produces false positives, pin the version and file an issue.
- If a component cannot be made accessible without a redesign, document the limitation in `KNOWN_ISSUES.md` and defer the redesign.
- If CI browser differences cause flakes, use Playwright's project matrix or pin `ubuntu-latest`.
