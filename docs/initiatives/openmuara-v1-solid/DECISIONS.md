> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara v1 Solid — Decision Log

> **Purpose:** Running record of decisions made during this initiative.
> **Format:**
> - **Formal Decisions** (`D###`): Human-reviewed or significant architecture choices.
> - **Auto-Decisions** (`A###`): AI-made decisions in AFK mode.

---

## Formal Decisions

_No formal decisions yet._

---

## Auto-Decisions (AFK Mode)

> **Rule:** In AFK mode, the AI MUST auto-log every non-trivial decision here. A "non-trivial" decision is any choice where >1 valid option existed and the AI picked one without human input.

| ID | Prompt | Decision | Options Considered | Why Chosen | Reversible? | Date |
|----|--------|----------|-------------------|------------|-------------|------|
| | | | | | | |

---

## How to Add a New Decision

### Formal Decision
1. Pick the next number: `D###`.
2. Use the format below.
3. Update `TRACKING.md` to reference this decision if it affects future steps.

### Auto-Decision (AFK Mode)
1. After executing a prompt, ask: "Did I make any non-trivial decision?"
2. If yes, append a row to the Auto-Decisions table.
3. If no, append a row with "None" to maintain the audit trail.

### Formal Decision Template

```markdown
### D001 — [Short Title]

- **Date:** YYYY-MM-DD
- **Context:** What required a decision?
- **Options Considered:**
  - Option A: ...
  - Option B: ...
- **Decision:** We chose Option __ because ...
- **Consequences:** ...
- **Reversible?** Yes / No
- **Reversal Trigger:** If ..., revisit.
- **Logged By:** AI Agent / Human
```
