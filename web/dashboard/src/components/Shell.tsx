import { useEffect, useRef, useState } from 'preact/hooks';
import { useFocusTrap } from '../hooks/useFocusTrap';
import { Button } from './Button';
import { Icon } from './Icon';
import { getDensity, getTheme, setDensity, setTheme, type Density, type Theme } from '../theme';
import type { ConnectionStatus } from '../hooks/useConnectionStatus';

type Tab = 'overview' | 'ledger' | 'transactions' | 'webhooks';

interface ShellProps {
  tab: Tab;
  onTabChange: (tab: Tab) => void;
  showHelp: boolean;
  onToggleHelp: () => void;
  onReload?: () => void;
  reloadKey?: number;
  connectionStatus?: ConnectionStatus;
  children: preact.ComponentChildren;
}

const TABS: { key: Tab; label: string }[] = [
  { key: 'overview', label: 'Overview' },
  { key: 'ledger', label: 'Ledger' },
  { key: 'transactions', label: 'Transactions' },
  { key: 'webhooks', label: 'Webhooks' },
];

export function Shell({
  tab,
  onTabChange,
  showHelp,
  onToggleHelp,
  onReload,
  reloadKey = 0,
  connectionStatus,
  children,
}: ShellProps) {
  const [theme, setThemeState] = useState<Theme>(getTheme());
  const [density, setDensityState] = useState<Density>(getDensity());
  const helpRef = useFocusTrap<HTMLDivElement>(showHelp);
  const closeHelpButtonRef = useRef<HTMLButtonElement>(null);
  const tabRefs = useRef<Map<Tab, HTMLButtonElement>>(new Map());

  useEffect(() => {
    const onStorage = (e: StorageEvent) => {
      if (e.key === 'muara-theme' && (e.newValue === 'light' || e.newValue === 'dark')) {
        setThemeState(e.newValue);
      }
      if (e.key === 'muara-density' && (e.newValue === 'comfortable' || e.newValue === 'compact')) {
        setDensityState(e.newValue);
      }
    };
    window.addEventListener('storage', onStorage);
    return () => window.removeEventListener('storage', onStorage);
  }, []);

  useEffect(() => {
    if (showHelp) {
      closeHelpButtonRef.current?.focus();
    }
  }, [showHelp]);

  const handleToggleTheme = () => {
    const next = getTheme() === 'dark' ? 'light' : 'dark';
    setTheme(next);
    setThemeState(next);
  };

  const handleToggleDensity = () => {
    const next = getDensity() === 'compact' ? 'comfortable' : 'compact';
    setDensity(next);
    setDensityState(next);
  };

  const focusTab = (key: Tab) => {
    tabRefs.current.get(key)?.focus();
  };

  const handleTabKeyDown = (e: KeyboardEvent, index: number) => {
    let nextKey: Tab | null = null;
    if (e.key === 'ArrowRight') {
      e.preventDefault();
      nextKey = TABS[(index + 1) % TABS.length].key;
    } else if (e.key === 'ArrowLeft') {
      e.preventDefault();
      nextKey = TABS[(index - 1 + TABS.length) % TABS.length].key;
    } else if (e.key === 'Home') {
      e.preventDefault();
      nextKey = TABS[0].key;
    } else if (e.key === 'End') {
      e.preventDefault();
      nextKey = TABS[TABS.length - 1].key;
    }
    if (nextKey) {
      onTabChange(nextKey);
      requestAnimationFrame(() => focusTab(nextKey!));
    }
  };

  const statusLabel = connectionStatus === 'online' ? 'Connected' : connectionStatus === 'offline' ? 'Disconnected' : 'Checking connection';

  return (
    <>
      <a href="#main-content" class="skip-link">Skip to main content</a>

      <header class="dashboard-header">
        <div class="header-toolbar">
          <h1>OpenMuara Dashboard</h1>
          <div class="header-actions">
            {connectionStatus && (
              <span
                class={`connection-status connection-${connectionStatus} ${connectionStatus === 'checking' ? 'connection-pulse' : ''}`}
                aria-label={statusLabel}
                title={statusLabel}
              >
                <span class="connection-dot" aria-hidden="true" />
                {statusLabel}
              </span>
            )}
            <button
              class="command-trigger"
              onClick={onToggleHelp}
              aria-label="Open command palette"
              title="Command palette (Cmd+K)"
            >
              <Icon name="command" size={14} />
              <span>Command</span>
              <span class="command-shortcut">⌘K</span>
            </button>
            {onReload && (
              <Button
                variant="secondary"
                size="sm"
                icon="refresh"
                title="Refresh data (r)"
                onClick={onReload}
                aria-label="Refresh dashboard data"
              >
                Reload
              </Button>
            )}
            <Button
              variant="secondary"
              size="sm"
              icon={theme === 'dark' ? 'sun' : 'moon'}
              title={`Switch to ${theme === 'dark' ? 'light' : 'dark'} mode`}
              onClick={handleToggleTheme}
              aria-label={`Switch to ${theme === 'dark' ? 'light' : 'dark'} mode`}
            >
              {theme === 'dark' ? 'Light' : 'Dark'}
            </Button>
            <Button
              variant="secondary"
              size="sm"
              icon="settings"
              title={`Density: ${density}`}
              onClick={handleToggleDensity}
              aria-label={`Switch to ${density === 'compact' ? 'comfortable' : 'compact'} density`}
            >
              {density === 'compact' ? 'Compact' : 'Comfort'}
            </Button>
            <Button
              variant="secondary"
              size="sm"
              icon="help"
              onClick={onToggleHelp}
              aria-label="Show keyboard shortcuts"
            >
              Help
            </Button>
          </div>
        </div>
      </header>

      <nav class="tabs" role="tablist" aria-label="Dashboard sections">
        {TABS.map((t, index) => (
          <button
            key={t.key}
            ref={(el) => { if (el) tabRefs.current.set(t.key, el); }}
            class={`tab ${tab === t.key ? 'active' : ''}`}
            role="tab"
            aria-selected={tab === t.key}
            tabIndex={tab === t.key ? 0 : -1}
            onClick={() => onTabChange(t.key)}
            onKeyDown={(e) => handleTabKeyDown(e, index)}
          >
            {t.label}
          </button>
        ))}
      </nav>

      <main id="main-content" tabIndex={-1} key={reloadKey}>{children}</main>

      {showHelp && (
        <div
          class="help-modal active"
          onClick={(e) => {
            if (e.target === e.currentTarget) onToggleHelp();
          }}
          role="dialog"
          aria-modal="true"
          aria-label="Keyboard shortcuts"
          ref={helpRef}
        >
          <div class="help-box">
            <div class="flex justify-between items-center mb-3">
              <h2 style={{ margin: 0 }}>Keyboard shortcuts</h2>
              <button class="btn btn-ghost btn-sm" onClick={onToggleHelp} ref={closeHelpButtonRef} aria-label="Close help">
                <Icon name="close" size={16} />
              </button>
            </div>
            <ul>
              <li><kbd>?</kbd> — show or hide this help</li>
              <li><kbd>Cmd</kbd> / <kbd>Ctrl</kbd> + <kbd>K</kbd> — open command palette</li>
              <li><kbd>/</kbd> — focus the current view's search box</li>
              <li><kbd>1</kbd> / <kbd>2</kbd> / <kbd>3</kbd> / <kbd>4</kbd> — switch to Overview / Ledger / Transactions / Webhooks</li>
              <li><kbd>r</kbd> — refresh dashboard data</li>
              <li><kbd>d</kbd> — toggle dark/light mode</li>
              <li><kbd>esc</kbd> — close panels and this help</li>
            </ul>
            <h2>How the dashboard works</h2>
            <p class="muted">
              The dashboard is a control plane for OpenMuara providers and webhooks.
              Enable providers on the Overview tab, configure webhooks on the Webhooks tab,
              and inspect the ledger and transactions in real time. Data refreshes automatically
              every 2 seconds while the page is visible, and pauses when you switch browser tabs.
            </p>
          </div>
        </div>
      )}
    </>
  );
}
