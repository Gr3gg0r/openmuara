> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara AI Slop Audit — Risk Register

> **Created:** 2026-07-08
> **Last Updated:** 2026-07-08

---

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Refactoring shared provider code introduces emulation bugs | Medium | High | Keep provider-specific quirks explicit; add conformance tests before abstracting. |
| Changing seed data breaks existing tests or screenshots | Medium | Medium | Update baselines and smoke tests together; gate seeding behind a flag. |
| New font or icon assets increase bundle size / offline friction | Low | Low | Use system fallback; prefer lightweight SVG icons and self-hosted or no external font. |
| Docs rewrite accidentally removes important technical detail | Medium | Medium | Review diffs for accuracy; keep decision rationale in `DECISIONS.md`. |
| `any`/`interface{}` removal causes API breakage | Low | High | Target internal helpers first; avoid changing exported interfaces without versioning. |
