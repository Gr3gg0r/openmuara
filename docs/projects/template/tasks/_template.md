> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This template is subordinate to it.**

# Step ## — [Task Title]

> **Purpose:** Detailed specification of WHAT and WHY. The "brain" behind the prompt.
> **Audience:** AI Agent + Human reviewer
> **Related Prompt:** `prompts/##-*.md` (the execution command)

---

## Objective

[One-paragraph description of what this step achieves and why it matters.]

---

## Target Files

| # | File | Action | Repo Path |
|---|------|--------|-----------|
| 1 | `internal/path/to/file.go` | New / Modify / Delete | `<repo-root>/internal/path/to/file.go` |

---

## Constraints & Security

- Constraint 1: __________
- Constraint 2: __________
- Security consideration: __________
- Performance consideration: __________

---

## Schema / Interface (if applicable)

```go
type Example struct {
    Field string
}
```

---

## Error Handling Requirements

- What happens on the unhappy path? Specify expected error responses, fallback behavior, and logging requirements.
- Do NOT leak sensitive data (PII, tokens, internal paths) in errors.

---

## Performance & Observability

- Any concurrency risks? Specify mutex/atomic strategy.
- Any new metrics or logs needed?

---

## BDD / TDD Quality Gates

- [ ] Unit test: scenario A → expected outcome.
- [ ] Unit test: scenario B → expected outcome.
- [ ] Unit test: edge case (empty input, max length, invalid format) → expected outcome.
- [ ] Integration test: component A + component B → expected outcome.
- [ ] Build / test / lint passes.
- [ ] File-size gates respected (250 lines/file, 80 lines/function).

---

## Rollback Trigger

If any of the following happens, STOP and consult `RISKS.md`:
- __________
- __________
