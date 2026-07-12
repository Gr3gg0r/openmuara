> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first.**

# Appendix B — Implementation Snippets

> Copy-paste starting points. Adapt to the actual codebase and keep them minimal.

---

## 1. Blocking theme script for dashboard `index.html`

```html
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta name="color-scheme" content="light dark">
  <meta name="theme-color" content="#f8fafc" id="theme-color">
  <title>OpenMuara Dashboard</title>
  <script>
    (function () {
      const STORAGE_KEY = 'muara-theme';
      const stored = localStorage.getItem(STORAGE_KEY);
      const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
      const theme = stored || (prefersDark ? 'dark' : 'light');
      document.documentElement.setAttribute('data-theme', theme);
      document.getElementById('theme-color').setAttribute(
        'content',
        theme === 'dark' ? '#0f172a' : '#f8fafc'
      );
    })();
  </script>
</head>
```

## 2. Theme toggle helper

```typescript
const STORAGE_KEY = 'muara-theme';

type Theme = 'light' | 'dark';

export function getTheme(): Theme {
  const stored = localStorage.getItem(STORAGE_KEY) as Theme | null;
  if (stored) return stored;
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
}

export function setTheme(theme: Theme): void {
  document.documentElement.setAttribute('data-theme', theme);
  localStorage.setItem(STORAGE_KEY, theme);
  const meta = document.getElementById('theme-color') as HTMLMetaElement | null;
  if (meta) meta.content = theme === 'dark' ? '#0f172a' : '#f8fafc';
  window.dispatchEvent(new StorageEvent('storage', { key: STORAGE_KEY, newValue: theme }));
}

export function toggleTheme(): void {
  setTheme(getTheme() === 'dark' ? 'light' : 'dark');
}

export function listenToOSThemeChange(cb: (theme: Theme) => void): () => void {
  const mq = window.matchMedia('(prefers-color-scheme: dark)');
  const handler = () => {
    if (!localStorage.getItem(STORAGE_KEY)) {
      const theme = mq.matches ? 'dark' : 'light';
      document.documentElement.setAttribute('data-theme', theme);
      cb(theme);
    }
  };
  mq.addEventListener('change', handler);
  return () => mq.removeEventListener('change', handler);
}

export function syncThemeAcrossTabs(cb: (theme: Theme) => void): () => void {
  const handler = (e: StorageEvent) => {
    if (e.key === STORAGE_KEY && e.newValue) {
      document.documentElement.setAttribute('data-theme', e.newValue);
      cb(e.newValue as Theme);
    }
  };
  window.addEventListener('storage', handler);
  return () => window.removeEventListener('storage', handler);
}
```

## 3. Theme toggle button in Preact

```tsx
import { toggleTheme, getTheme } from '../theme';

export function ThemeToggle() {
  const [theme, setTheme] = useState(getTheme());

  useEffect(() => {
    return syncThemeAcrossTabs(setTheme);
  }, []);

  return (
    <button
      class="secondary"
      aria-label={`Switch to ${theme === 'dark' ? 'light' : 'dark'} mode`}
      onClick={() => {
        toggleTheme();
        setTheme(getTheme());
      }}
    >
      {theme === 'dark' ? '☀' : '☾'}
    </button>
  );
}
```

## 4. Reduced-motion safe transition

```css
@media (prefers-reduced-motion: no-preference) {
  body,
  .card,
  button,
  input,
  select,
  table,
  th,
  td,
  .alert,
  .help-box {
    transition: background-color 150ms ease, color 150ms ease, border-color 150ms ease;
  }
}
```

## 5. Provider pay page minimal dark mode

```html
<head>
  <meta name="color-scheme" content="light dark">
  <style>
    :root {
      --bg: #f8fafc;
      --surface: #ffffff;
      --text: #0f172a;
      --muted: #64748b;
      --border: #e2e8f0;
      --primary: #2563eb;
      --primary-text: #ffffff;
    }
    @media (prefers-color-scheme: dark) {
      :root {
        --bg: #0f172a;
        --surface: #1e293b;
        --text: #f8fafc;
        --muted: #94a3b8;
        --border: #334155;
        --primary: #3b82f6;
        --primary-text: #ffffff;
      }
    }
    body { background: var(--bg); color: var(--text); }
    .container { background: var(--surface); border: 1px solid var(--border); }
    button[type="submit"] { background: var(--primary); color: var(--primary-text); }
  </style>
</head>
```

## 6. Example mini-app toggle script

```html
<script>
  (function () {
    const key = 'muara-theme';
    const root = document.documentElement;
    const stored = localStorage.getItem(key);
    const prefersDark = matchMedia('(prefers-color-scheme: dark)').matches;
    const theme = stored || (prefersDark ? 'dark' : 'light');
    root.setAttribute('data-theme', theme);

    document.getElementById('theme-toggle').addEventListener('click', () => {
      const next = root.getAttribute('data-theme') === 'dark' ? 'light' : 'dark';
      root.setAttribute('data-theme', next);
      localStorage.setItem(key, next);
    });
  })();
</script>
```
