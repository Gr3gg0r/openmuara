const THEME_KEY = 'muara-theme';
const DENSITY_KEY = 'muara-density';

export type Theme = 'light' | 'dark';
export type Density = 'comfortable' | 'compact';

function getMetaThemeColor(): HTMLMetaElement | null {
  return document.getElementById('theme-color') as HTMLMetaElement | null;
}

export function getTheme(): Theme {
  const stored = localStorage.getItem(THEME_KEY) as Theme | null;
  if (stored === 'light' || stored === 'dark') {
    return stored;
  }
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
}

export function applyTheme(theme: Theme): void {
  document.documentElement.setAttribute('data-theme', theme);
  const meta = getMetaThemeColor();
  if (meta) {
    meta.content = theme === 'dark' ? '#0f172a' : '#f8fafc';
  }
}

export function setTheme(theme: Theme): void {
  applyTheme(theme);
  localStorage.setItem(THEME_KEY, theme);
  window.dispatchEvent(new StorageEvent('storage', { key: THEME_KEY, newValue: theme }));
}

export function toggleTheme(): Theme {
  const next = getTheme() === 'dark' ? 'light' : 'dark';
  setTheme(next);
  return next;
}

export function getDensity(): Density {
  const stored = localStorage.getItem(DENSITY_KEY) as Density | null;
  return stored === 'compact' ? 'compact' : 'comfortable';
}

export function applyDensity(density: Density): void {
  document.documentElement.setAttribute('data-density', density);
}

export function setDensity(density: Density): void {
  applyDensity(density);
  localStorage.setItem(DENSITY_KEY, density);
  window.dispatchEvent(new StorageEvent('storage', { key: DENSITY_KEY, newValue: density }));
}

export function toggleDensity(): Density {
  const next = getDensity() === 'compact' ? 'comfortable' : 'compact';
  setDensity(next);
  return next;
}

export function listenToOSThemeChange(): () => void {
  const mq = window.matchMedia('(prefers-color-scheme: dark)');
  const handler = () => {
    if (!localStorage.getItem(THEME_KEY)) {
      applyTheme(mq.matches ? 'dark' : 'light');
    }
  };
  mq.addEventListener('change', handler);
  return () => mq.removeEventListener('change', handler);
}

export function syncThemeAcrossTabs(): () => void {
  const handler = (e: StorageEvent) => {
    if (e.key === THEME_KEY && (e.newValue === 'light' || e.newValue === 'dark')) {
      applyTheme(e.newValue);
    }
    if (e.key === DENSITY_KEY && (e.newValue === 'comfortable' || e.newValue === 'compact')) {
      applyDensity(e.newValue);
    }
  };
  window.addEventListener('storage', handler);
  return () => window.removeEventListener('storage', handler);
}

export function initializeAppearance(): void {
  applyTheme(getTheme());
  applyDensity(getDensity());
}
