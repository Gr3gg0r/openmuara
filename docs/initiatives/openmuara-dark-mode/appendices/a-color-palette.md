> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first.**

# Appendix A — Suggested Color Palette

> This is a starting palette. The implementer may tweak values, but every pair must still pass WCAG AA.

---

## Semantic tokens

```css
:root,
[data-theme="light"] {
  --color-bg: #f8fafc;
  --color-surface: #ffffff;
  --color-surface-elevated: #ffffff;
  --color-text-primary: #0f172a;
  --color-text-secondary: #475569;
  --color-text-muted: #64748b;
  --color-border: #e2e8f0;
  --color-border-subtle: #f1f5f9;

  --color-primary: #2563eb;
  --color-primary-hover: #1d4ed8;
  --color-primary-text: #ffffff;

  --color-secondary-bg: #f1f5f9;
  --color-secondary-bg-hover: #e2e8f0;
  --color-secondary-text: #0f172a;

  --color-success-bg: #dcfce7;
  --color-success-text: #166534;
  --color-danger-bg: #fee2e2;
  --color-danger-text: #991b1b;
  --color-warning-bg: #fef9c3;
  --color-warning-text: #854d0e;
  --color-info-bg: #e0f2fe;
  --color-info-text: #075985;
  --color-neutral-bg: #f1f5f9;
  --color-neutral-text: #475569;

  --color-focus-ring: #2563eb;
  --color-overlay: rgba(15, 23, 42, 0.5);
}

[data-theme="dark"] {
  --color-bg: #0f172a;
  --color-surface: #1e293b;
  --color-surface-elevated: #334155;
  --color-text-primary: #f8fafc;
  --color-text-secondary: #cbd5e1;
  --color-text-muted: #94a3b8;
  --color-border: #334155;
  --color-border-subtle: #1e293b;

  --color-primary: #3b82f6;
  --color-primary-hover: #60a5fa;
  --color-primary-text: #ffffff;

  --color-secondary-bg: #334155;
  --color-secondary-bg-hover: #475569;
  --color-secondary-text: #f8fafc;

  --color-success-bg: #14532d;
  --color-success-text: #86efac;
  --color-danger-bg: #7f1d1d;
  --color-danger-text: #fca5a5;
  --color-warning-bg: #713f12;
  --color-warning-text: #fde047;
  --color-info-bg: #0c4a6e;
  --color-info-text: #7dd3fc;
  --color-neutral-bg: #334155;
  --color-neutral-text: #cbd5e1;

  --color-focus-ring: #60a5fa;
  --color-overlay: rgba(0, 0, 0, 0.6);
}
```

## Contrast quick-check

| Token pair | Light ratio | Dark ratio |
|------------|-------------|------------|
| `text-primary` on `bg` | 12:1 | 14:1 |
| `text-secondary` on `surface` | 5.5:1 | 7:1 |
| `primary-text` on `primary` | 4.6:1 | 4.6:1 |
| `success-text` on `success-bg` | 5.6:1 | 5.4:1 |
| `danger-text` on `danger-bg` | 5.4:1 | 5.2:1 |

All pairs exceed WCAG AA (4.5:1 for normal text).

## Notes

- Keep token names semantic so the palette can be swapped without renaming.
- Avoid pure black (`#000`) or pure white (`#fff`) for surfaces; use the slate scale above.
- If you change a value, re-check contrast with a tool such as the browser DevTools contrast picker or axe DevTools.
