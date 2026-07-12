> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This tracker is subordinate to it.**

# OpenMuara Accessibility & Usability Polish — Tracker

## Status

| Item | State |
|------|-------|
| Initiative README | ✅ Complete |
| Audit findings | ✅ Complete |
| Prompt 01 — Dashboard keyboard navigation | ✅ Complete |
| Prompt 02 — Dashboard labels and live regions | ✅ Complete |
| Prompt 03 — Provider pages focus and landmarks | ✅ Complete |
| Prompt 04 — Example apps accessibility | ✅ Complete |
| Prompt 05 — Shortcuts and theme polish | ✅ Complete |
| Enhancement E1 — Skip link | ✅ Complete |
| Enhancement E2 — `prefers-contrast: more` | ✅ Complete |
| Enhancement E3 — Playwright a11y smoke test | ✅ Complete |
| Enhancement E4 — CI contrast regression check | ✅ Complete |
| Quality gates | ✅ Passed |
| Merge to `dev` | ✅ Complete |
| Teardown branch | ✅ Complete |

## Branch

`feat/a11y-enhancements` (to be merged into `dev` and deleted)

## Recent commits

- `acd8a9e` — feat(a11y): dashboard keyboard navigation, labels, live regions, theme sync (P01, P02, P05)
- `4c392a5` — feat(a11y): add focus indicators and main landmarks to provider pages (P03)
- `8bb9ee3` — feat(a11y): improve example mini-app accessibility (P04)
- `f066246` — feat(a11y): skip link, prefers-contrast, E2E a11y tests, CI contrast check (E1–E4)
- `f39a0f3` — docs(a11y): update tracker, handoff, and changelog for E1–E4 enhancements
- `c12554d` — fix(a11y): wait for theme application in contrast check

## Quality gate results

- `go build ./...` ✅
- `go test ./...` ✅
- `go vet ./...` ✅
- `golangci-lint run` ✅ (0 issues)
- `cd web/dashboard && npm run test:ci` ✅ (16/16)
- `cd web/dashboard && npm run test:e2e` ✅ (4/4)
- `cd web/dashboard && npm run a11y:contrast` ✅ (light + dark)
- `cd web/dashboard && npm run bundle-size` ✅ (JS 13.24 KiB / total dist 189.33 KiB)
- `./scripts/smoke-test.sh` ✅
- `./scripts/audit-trackers.sh` ✅
- `./scripts/check-sizes.sh` ✅ (advisory warnings only)
- `./scripts/check-forbidden.sh` ✅

## Blockers

- None.

## Notes

- Bundle size increased slightly from the original a11y initiative (JS 13.21 KiB / total dist 187.73 KiB) to the enhancements pass (JS 13.24 KiB / total dist 189.33 KiB), both well within limits.
- Dark mode tokens were adjusted for WCAG AA contrast: `--color-text-muted` lightened, `--color-primary` lightened, and dedicated dark-mode link colors added.
