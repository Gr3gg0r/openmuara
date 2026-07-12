> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Accessibility & Usability Audit Known Issues

> **Created:** 2026-07-08
> **Last Updated:** 2026-07-09
> **Status:** ✅ Complete — all findings triaged and remediated

---

## How to record an a11y/usability finding

When the audit identifies a gap, record it using this template:

```markdown
### AXXX — <short summary>

- **Area:** Automated scan / Keyboard / Screen reader / Color & contrast / Motion / Forms / Mobile / Dynamic content / Usability heuristic
- **View/Component:** <name>
- **Severity:** Critical / Serious / Moderate / Minor
- **WCAG criterion:** e.g., 1.3.1 Info and Relationships (A), 2.1.1 Keyboard (A), 1.4.3 Contrast (Minimum) (AA)
- **Real impact:** Concrete consequence for users
- **Current behavior:** What the dashboard does now
- **Expected behavior:** What WCAG or best practice requires
- **Recommended fix:** Specific code/test/doc change
- **Decision:** Fix / Document deviation / Accept limitation
- **Owner:** AI Agent / Maintainer / External reviewer
- **Review date:** YYYY-MM-DD
```

## Severity definitions

| Severity | Definition | Example |
|---|---|---|
| Critical | Blocks a user group from completing a primary task | Button not keyboard-operable |
| Serious | Major difficulty for a user group; workaround possible | Missing label on a form input |
| Moderate | Noticeable friction; does not block core workflow | Focus indicator barely visible |
| Minor | Cosmetic or low-impact issue | Redundant announcement |

## Findings resolved during this initiative

### A001 — Command palette dialog lacked role and accessible name

- **Area:** Screen reader / Dynamic content
- **View/Component:** `CommandPalette.tsx`
- **Severity:** Serious
- **WCAG criterion:** 4.1.2 Name, Role, Value (A), 1.3.1 Info and Relationships (A)
- **Real impact:** Screen-reader users could not identify the overlay as a dialog or distinguish it from page content.
- **Current behavior:** The palette container was a generic `<div>` with no role or label.
- **Expected behavior:** Modal overlay exposes `role="dialog"`, `aria-modal="true"`, and an accessible name.
- **Fixed in:** `web/dashboard/src/components/CommandPalette.tsx` — added `role="dialog"`, `aria-modal="true"`, `aria-label="Command palette"`.
- **Decision:** ✅ Fix
- **Owner:** AI Agent
- **Review date:** 2026-07-09

### A002 — Confirm dialog backdrop was not keyboard-accessible and closed on inner clicks

- **Area:** Keyboard / Dynamic content
- **View/Component:** `ConfirmDialog.tsx`
- **Severity:** Serious
- **WCAG criterion:** 2.1.1 Keyboard (A), 2.1.2 No Keyboard Trap (A), 4.1.2 Name, Role, Value (A)
- **Real impact:** Clicking inside the modal message accidentally closed the dialog; the backdrop could not be operated by keyboard.
- **Current behavior:** Backdrop used `role="presentation"` and `onClick={onCancel}` on the parent without target checks.
- **Expected behavior:** Backdrop click-to-close only fires when the backdrop itself is clicked and is keyboard dismissible.
- **Fixed in:** `web/dashboard/src/components/ConfirmDialog.tsx` — converted backdrop to a focusable close control with `role="button"`, `tabIndex={-1}`, `aria-label="Close dialog"`, target/currentTarget guard, and `Escape`/`Enter`/`Space` handling.
- **Decision:** ✅ Fix
- **Owner:** AI Agent
- **Review date:** 2026-07-09

### A003 — Form labels were not programmatically associated with inputs

- **Area:** Forms / Screen reader
- **View/Component:** `WebhookConfig.tsx`, `ProviderDetail.tsx`
- **Severity:** Serious
- **WCAG criterion:** 1.3.1 Info and Relationships (A), 3.3.2 Labels or Instructions (A), 4.1.2 Name, Role, Value (A)
- **Real impact:** Screen-reader users could not determine which label described which input.
- **Current behavior:** Labels used JSX `for=` attribute, which does not render the HTML `for` attribute in Preact/React.
- **Expected behavior:** Labels use `htmlFor=` and are associated with matching `id` attributes on inputs.
- **Fixed in:** `web/dashboard/src/components/WebhookConfig.tsx`, `web/dashboard/src/views/ProviderDetail.tsx`.
- **Decision:** ✅ Fix
- **Owner:** AI Agent
- **Review date:** 2026-07-09

### A004 — Timeline list had redundant ARIA role

- **Area:** Screen reader
- **View/Component:** `Timeline.tsx`
- **Severity:** Minor
- **WCAG criterion:** 4.1.2 Name, Role, Value (A)
- **Real impact:** Redundant `role="list"` on an `<ol>` added noise without harm.
- **Fixed in:** `web/dashboard/src/components/Timeline.tsx` — removed redundant `role="list"`.
- **Decision:** ✅ Fix
- **Owner:** AI Agent
- **Review date:** 2026-07-09

### A005 — Unused `Shell` component and test added maintenance noise

- **Area:** Usability heuristic / Code health
- **View/Component:** `Shell.tsx`, `tests/Shell.test.tsx`
- **Severity:** Minor
- **WCAG criterion:** N/A
- **Real impact:** Dead code increased cognitive load and risk of stale patterns.
- **Fixed in:** Deleted `web/dashboard/src/components/Shell.tsx` and `web/dashboard/tests/Shell.test.tsx`.
- **Decision:** ✅ Fix
- **Owner:** AI Agent
- **Review date:** 2026-07-09

### A006 — Provider display text used `dangerouslySetInnerHTML`

- **Area:** Usability / Security
- **View/Component:** `Providers.tsx`
- **Severity:** Moderate
- **WCAG criterion:** N/A (best practice)
- **Real impact:** Escaped HTML could render unpredictably and complicate screen-reader parsing.
- **Fixed in:** `web/dashboard/src/components/Providers.tsx` — rendered values as plain text and removed the `enabled`/`disabled` color-only badges in favor of a single explicit "active" badge.
- **Decision:** ✅ Fix
- **Owner:** AI Agent
- **Review date:** 2026-07-09

## Accepted deviations

None. All triaged findings were fixed. Any future deviations must be recorded here with rationale and a review date.

| ID | Area | Deviation | Rationale | Owner | Review date |
|---|---|---|---|---|---|
| — | — | — | — | — | — |

## Active findings

None.
