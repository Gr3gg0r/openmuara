> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This template is subordinate to it.**

# <Project Name> — Decision Log

> **Purpose:** Running record of decisions made during this project. Prevents re-debate and provides audit trail.
> **Format:**
> - **Formal Decisions** (`D###`): Human-reviewed or significant architecture choices. Full ADR format.
> - **Auto-Decisions** (`A###`): AI-made decisions in AFK mode. Lightweight, logged automatically per prompt.

---

## Formal Decisions

### D001 — [Short Decision Title]

- **Date:** YYYY-MM-DD
- **Context:** What was the situation that required a decision?
- **Options Considered:**
  - Option A: __________ (pros: ___, cons: ___)
  - Option B: __________ (pros: ___, cons: ___)
- **Decision:** We chose Option __ because __________.
- **Consequences:** __________
- **Reversible?** Yes / No
- **Reversal Trigger:** If __________, revisit this decision.
- **Logged By:** AI Agent / Human ___________

---

## Auto-Decisions (AFK Mode)

> **Rule:** In AFK mode, the AI MUST auto-log every non-trivial decision here. A "non-trivial" decision is any choice where >1 valid option existed and the AI picked one without human input.
> **Format:** One entry per prompt. Append to the table. If a prompt had zero decisions, write "None".

| ID | Prompt | Decision | Options Considered | Why Chosen | Reversible? | Date |
|----|--------|----------|-------------------|------------|-------------|------|
| A01 | | | | | | |
| A02 | | | | | | |

---

## How to Add a New Decision

### Formal Decision
1. Pick the next number: `D###`.
2. Copy the D001 template above.
3. Fill all fields. Be specific — vague decisions are useless later.
4. Update `TRACKING.md` to reference this decision if it affects future steps.

### Auto-Decision (AFK Mode)
1. After executing a prompt, ask: "Did I make any non-trivial decision?"
2. If yes, append a row to the Auto-Decisions table with `A###` ID.
3. If no, append a row with "None" to maintain the audit trail.
