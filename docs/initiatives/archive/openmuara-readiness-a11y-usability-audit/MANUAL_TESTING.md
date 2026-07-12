> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Accessibility & Usability Audit Manual Testing Guide

> **Created:** 2026-07-09
> **Last Updated:** 2026-07-09
> **Status:** ✅ Complete — manual test guide published

---

This guide documents the manual tests that automated scanners cannot fully cover. Run these tests after the automated baseline scan and after fixing keyboard, screen-reader, and mobile issues.

## Environment

- **Browser:** Latest Chrome or Safari.
- **Screen reader (primary):** VoiceOver on macOS.
- **Screen reader (secondary):** NVDA on Windows, if available.
- **Input methods:** Keyboard only; touch emulation at 320 px and 375 px widths.
- **Themes:** Light and dark.

## Setup

1. Start OpenMuara:
   ```bash
   go run ./cmd/muara start
   ```
2. Open the dashboard at `http://localhost:9000/_admin`.
3. Generate sample state if needed:
   ```bash
   go run ./cmd/muara seed --providers fawry,stripe
   ```
4. Enable VoiceOver: `Cmd + F5` on macOS.

## Test 1 — Keyboard navigation

**Goal:** All primary workflows work without a mouse.

**Steps:**
1. Press `Tab` repeatedly. Verify every interactive control receives focus in a logical order.
2. Verify focus indicators are visible on all focused elements.
3. Press `Shift+Tab` to navigate backwards.
4. Press `Enter` or `Space` to activate focused buttons and links.
5. Press `Escape` to close the command palette and any open dialogs.
6. Use `Ctrl+K` (or `Cmd+K`) to open the command palette; navigate it with arrow keys; press `Enter` to select; press `Escape` to close.
7. Navigate to **Ledger**, open a detail row with `Enter`, and return with `Escape` or focus.
8. Navigate to **Webhooks**, select a webhook, and use the replay button via keyboard.

**Pass criteria:**
- No keyboard traps.
- Focus order matches visual order.
- All actions reachable and operable.

## Test 2 — Screen-reader labels and announcements

**Goal:** Controls and dynamic updates are announced correctly.

**Steps:**
1. With VoiceOver on, navigate to the sidebar.
2. Verify each navigation link is announced with a meaningful name (not "link" alone).
3. Tab to icon-only buttons (e.g., copy, refresh, filter). Verify each announces its purpose.
4. Open the command palette. Verify it announces as a dialog/search.
5. Trigger a webhook replay. Verify the result is announced via the live region.
6. Submit an invalid form. Verify the error message is announced.
7. Switch themes. Verify the change is announced or the toggle has an accessible name.

**Pass criteria:**
- No unlabeled controls.
- Status changes are announced without moving focus unexpectedly.

## Test 3 — Color and contrast

**Goal:** Information is not conveyed by color alone and contrast is sufficient.

**Steps:**
1. Open the dashboard in light mode.
2. Use a contrast checker (e.g., browser devtools, axe DevTools) on body text, labels, buttons, and status badges.
3. Switch to dark mode and repeat.
4. Identify any status that relies only on color (e.g., a green dot without "success" text). Verify it has an icon, label, or text alternative.

**Pass criteria:**
- Normal text ≥ 4.5:1.
- Large text and UI components ≥ 3:1.
- No color-only information cues.

## Test 4 — Mobile touch targets

**Goal:** Interactive elements are easy to tap on small screens.

**Steps:**
1. Open browser devtools and set viewport to 375 × 667.
2. Inspect every button, link, tab, and form control.
3. Verify each target is at least 36×36 CSS pixels (44×44 for primary actions).
4. Attempt to use pinch zoom. Verify the viewport does not block scaling.

**Pass criteria:**
- Zero targets below 36×36 px.
- Zoom is allowed.

## Test 5 — Reduced motion

**Goal:** Animations respect user preference.

**Steps:**
1. Enable reduced motion in macOS: System Settings → Accessibility → Display → Reduce Motion.
2. Open the dashboard and trigger animations (command palette open/close, toast, loading skeleton).
3. Verify motion is minimized or removed.
4. Disable reduced motion and verify functional animation returns.

**Pass criteria:**
- No essential information is lost when reduced motion is active.
- No vestibular-triggering motion remains.

## Test 6 — Focus management in dynamic content

**Goal:** Focus moves predictably when content changes.

**Steps:**
1. Open a dialog (e.g., command palette or confirm dialog). Verify focus moves into the dialog.
2. Press `Tab` inside the dialog. Verify focus is trapped.
3. Close the dialog. Verify focus returns to the triggering control.
4. Navigate to a different view. Verify focus moves to the main content or page title.
5. Trigger a toast notification. Verify focus does not move unexpectedly.

**Pass criteria:**
- Focus is managed for dialogs, route changes, and notifications.

## Test 7 — Form accessibility

**Goal:** Forms are usable with keyboard and screen reader.

**Steps:**
1. Navigate to **Settings** or **Webhook Config**.
2. Tab through each field. Verify each has a visible label or `aria-label`.
3. Submit the form with empty required fields. Verify errors are associated with inputs and announced.
4. Correct an error and verify the error state is cleared accessibly.

**Pass criteria:**
- 100% of inputs have accessible names.
- Errors are keyboard-focusable and announced.

## Recording findings

Record any issue in `KNOWN_ISSUES.md` using the finding template. Include:
- Test number and step.
- Environment (browser, screen reader, viewport).
- Expected vs. actual behavior.
- Severity and WCAG criterion.

## Related documents

- [`TRACKING.md`](TRACKING.md) — phases and metrics
- [`KNOWN_ISSUES.md`](KNOWN_ISSUES.md) — finding template
- [`REVIEW_CHECKLIST.md`](REVIEW_CHECKLIST.md) — sign-off checklist
