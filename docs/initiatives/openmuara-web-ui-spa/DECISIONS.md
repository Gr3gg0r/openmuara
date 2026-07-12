> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Web UI SPA — Decision Log

| ID | Decision | Status | Date | Notes |
|----|----------|--------|------|-------|
| D001 | Dashboard framework selection | ✅ Decided | 2026-07-02 | Vite + Preact (Option C). Embedded into Go binary at build time. |
| D002 | Build assets must be embeddable by Go | ✅ Decided | 2026-07-02 | No runtime Node dependency; `//go:embed dist`. |
| D003 | Preserve existing `/_admin` route stability | ✅ Decided | 2026-07-02 | Escape/pay URLs must not change. |
