> **Initiative:** OpenMuara Accessibility & Usability Polish

# Decision Log

## D01 — No component library

**Decision:** Fix accessibility with the existing Preact + vanilla CSS stack. Do not add Daisy UI, shadcn/ui, Tailwind, or similar libraries.

**Rationale:**
- The dashboard JS bundle is ~12 KiB; adding a library would multiply that.
- Provider pages are intentionally self-contained HTML files with no build step.
- The issues are small ARIA/focus gaps, not a lack of components.

**Consequences:** More manual code, but the bundle stays tiny and pages stay dependency-free.

## D02 — Clickable rows via inline button

**Decision:** Make table rows keyboard-accessible by adding a visible or visually hidden action button inside the row, rather than making the entire `<tr>` focusable.

**Rationale:**
- Rows already contain a "Replay webhook" button; adding a second "View details" button keeps the DOM semantic.
- Avoids the complexity and ARIA pitfalls of `role="button"` on a `<tr>`.

## D03 — Global announcements via single aria-live region

**Decision:** Use one app-level `aria-live="polite"` region for copy-to-clipboard and other transient announcements.

**Rationale:** Simple, reusable, and avoids scattering live regions across components.
