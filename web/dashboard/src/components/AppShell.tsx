import { useEffect, useRef, useState } from 'preact/hooks';
import { useFocusTrap } from '../hooks/useFocusTrap';
import { Button } from './Button';
import { Icon } from './Icon';
import { SidebarNav, type NavItem } from './SidebarNav';
import { getDensity, getTheme, setDensity, setTheme, type Density, type Theme } from '../theme';
import type { ConnectionStatus } from '../hooks/useConnectionStatus';

interface AppShellProps {
  active: NavItem;
  onNavigate: (item: NavItem) => void;
  showHelp: boolean;
  onToggleHelp: () => void;
  onReload?: () => void;
  reloadKey?: number;
  connectionStatus?: ConnectionStatus;
  role?: 'admin' | 'viewer' | '';
  children: preact.ComponentChildren;
}

export function AppShell({
  active,
  onNavigate,
  showHelp,
  onToggleHelp,
  onReload,
  reloadKey = 0,
  connectionStatus,
  role,
  children,
}: AppShellProps) {
  const [theme, setThemeState] = useState<Theme>(getTheme());
  const [density, setDensityState] = useState<Density>(getDensity());
  const helpRef = useFocusTrap<HTMLDivElement>(showHelp);
  const closeHelpButtonRef = useRef<HTMLButtonElement>(null);

  useEffect(() => {
    const onStorage = (e: StorageEvent) => {
      if (e.key === 'muara-theme' && (e.newValue === 'light' || e.newValue === 'dark')) {
        setThemeState(e.newValue);
      }
      if (e.key === 'muara-density' && (e.newValue === 'compact' || e.newValue === 'comfortable')) {
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

  const statusLabel = connectionStatus === 'online'
    ? 'Connected'
    : connectionStatus === 'offline'
      ? 'Disconnected'
      : 'Checking connection';

  return (
    <>
      <a href="#main-content" class="skip-link">Skip to main content</a>

      <div class="app-shell">
        <header class="app-header">
          <div class="app-header-inner">
            <h1 class="app-title">OpenMuara</h1>
            <div class="app-header-actions">
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

        <SidebarNav active={active} onNavigate={onNavigate} role={role} />

        <main id="main-content" class="app-main" tabIndex={-1} key={reloadKey}>
          {children}
        </main>
      </div>

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
              <li><kbd>1</kbd> / <kbd>2</kbd> / <kbd>3</kbd> — switch to Ledger / Webhooks / Settings</li>
              <li><kbd>r</kbd> — refresh dashboard data</li>
              <li><kbd>d</kbd> — toggle dark/light mode</li>
              <li><kbd>esc</kbd> — close panels and this help</li>
            </ul>
            <h2>How the dashboard works</h2>
            <p class="muted">
              The dashboard is a control plane for OpenMuara providers and webhooks.
              Enable providers and configure per-provider webhook targets on the Settings page,
              inspect webhook delivery attempts on the Webhooks page, and view the ledger in real time.
              Data refreshes automatically every 2 seconds while the page is visible, and pauses when
              you switch browser tabs.
            </p>
          </div>
        </div>
      )}
    </>
  );
}
