> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Accessibility & Usability Audit Glossary

> **Created:** 2026-07-09
> **Last Updated:** 2026-07-09
> **Status:** ✅ Complete — glossary published

---

This glossary defines terms used throughout the accessibility and usability audit initiative.

## A

**Accessible name** — The text a screen reader announces for a control. Can come from visible text, `aria-label`, `aria-labelledby`, or native HTML semantics (e.g., `<button>text</button>`).

**ARIA (Accessible Rich Internet Applications)** — A set of attributes that extend HTML semantics to make dynamic content and custom controls accessible.

**Assistive technology (AT)** — Software or hardware used by people with disabilities, such as screen readers, screen magnifiers, switch devices, and voice control.

## C

**Color-only cue** — Information conveyed only through color (e.g., a red dot without a text label). WCAG requires an additional non-color indicator.

**Contrast ratio** — The luminance difference between foreground and background colors. WCAG 2.1 AA requires 4.5:1 for normal text and 3:1 for large text/UI components.

## F

**Focus indicator** — The visible outline or highlight showing which element currently has keyboard focus.

**Focus management** — The practice of controlling where keyboard focus moves when content changes, such as opening a dialog or navigating routes.

**Focus trap** — Keeping keyboard focus inside a modal or dialog until it is closed, preventing users from tabbing to background content.

## K

**Keyboard accessible** — A control that can be reached and operated using only the keyboard.

## L

**Landmark** — A semantic region of a page (e.g., `main`, `nav`, `aside`, `header`) that helps screen-reader users navigate.

**Live region** — An ARIA mechanism (`aria-live`) to announce dynamic content changes without moving focus.

## R

**Reduced motion** — A user preference (`prefers-reduced-motion`) that asks interfaces to minimize animation.

## S

**Screen reader** — Assistive technology that reads screen content aloud and lets users interact via keyboard.

**Semantic HTML** — HTML elements that carry built-in meaning, such as `<button>`, `<nav>`, `<main>`, `<table>`, rather than generic `<div>` or `<span>`.

**Skip link** — A keyboard-only link at the top of a page that lets users jump to the main content.

## T

**Touch target** — The screen area that responds to a tap or click. WCAG 2.5.5 Target Size (AAA) recommends 44×44 CSS pixels; this initiative uses 36×36 as a minimum.

## V

**VoiceOver** — The built-in screen reader on macOS and iOS.

**NVDA** — A free screen reader for Windows.

## W

**WCAG (Web Content Accessibility Guidelines)** — International guidelines for web accessibility published by the W3C.

**WCAG Level A** — Minimum conformance level; addresses the most basic accessibility features.

**WCAG Level AA** — Target conformance level for most organizations; addresses major accessibility barriers.

**WCAG Level AAA** — Highest conformance level; not always practical for entire sites.

## Related documents

- [`README.md`](README.md) — initiative overview
- [`TRACKING.md`](TRACKING.md) — phases and metrics
- [`RECOMMENDATIONS.md`](RECOMMENDATIONS.md) — standards and tools
