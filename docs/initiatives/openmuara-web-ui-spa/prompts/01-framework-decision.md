> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# 01 — Framework Decision

## Goal

Choose the dashboard technology direction for OpenMuara and record the decision.

## Context

The dashboard currently lives in `internal/ui/` as vanilla HTML/JS embedded in the Go binary. It works, but the codebase is growing and becoming harder to maintain. We need a deliberate choice between:

1. **Vanilla JS** with better module structure.
2. **Alpine.js** or **petite-vue** for lightweight reactivity without a build step.
3. **Vite + Preact** for a React-compatible component model with minimal runtime.
4. **Vite + Svelte** for a compiled, no-virtual-DOM component model.

## Required Output

- Update `DECISIONS.md` with the chosen option and the reasons.
- Update `TRACKING.md` prompt 01 status to `✅`.
- Update `HANDOFF.md`.

## Decision Criteria

- Preserve single-binary, local-first deployment.
- No runtime internet dependency.
- Existing `/_admin` routes remain stable.
- Reasonable contributor onboarding (Node acceptable, but not ideal).
- Bundle size should not materially bloat the binary.

## Quality Gate

- Human review of `DECISIONS.md`.
