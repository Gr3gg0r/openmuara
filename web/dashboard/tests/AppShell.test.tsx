import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/preact';
import { AppShell } from '../src/components/AppShell';

const storage: Record<string, string> = {};

Object.defineProperty(globalThis, 'localStorage', {
  value: {
    getItem: (key: string) => storage[key] ?? null,
    setItem: (key: string, value: string) => { storage[key] = value; },
    removeItem: (key: string) => { delete storage[key]; },
    clear: () => { Object.keys(storage).forEach((key) => delete storage[key]); },
  },
  writable: true,
});

Object.defineProperty(globalThis, 'matchMedia', {
  value: (query: string) => ({
    matches: query === '(prefers-color-scheme: dark)',
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
  }),
  writable: true,
});

describe('AppShell', () => {
  beforeEach(() => {
    localStorage.clear();
    document.documentElement.removeAttribute('data-theme');
  });

  it('renders sidebar navigation and children', () => {
    render(
      <AppShell
        active="ledger"
        onNavigate={vi.fn()}
        showHelp={false}
        onToggleHelp={vi.fn()}
      >
        <div data-testid="child">content</div>
      </AppShell>,
    );

    expect(screen.getByTestId('nav-ledger')).toBeInTheDocument();
    expect(screen.getByTestId('nav-webhooks')).toBeInTheDocument();
    expect(screen.getByTestId('nav-settings')).toBeInTheDocument();
    expect(screen.getByTestId('child')).toBeInTheDocument();
  });

  it('calls onNavigate when a sidebar item is clicked', () => {
    const onNavigate = vi.fn();
    render(
      <AppShell
        active="ledger"
        onNavigate={onNavigate}
        showHelp={false}
        onToggleHelp={vi.fn()}
      >
        <div />
      </AppShell>,
    );

    fireEvent.click(screen.getByTestId('nav-webhooks'));
    expect(onNavigate).toHaveBeenCalledWith('webhooks');
  });

  it('renders a theme toggle button', () => {
    render(
      <AppShell
        active="ledger"
        onNavigate={vi.fn()}
        showHelp={false}
        onToggleHelp={vi.fn()}
      >
        <div />
      </AppShell>,
    );

    const toggle = screen.getByLabelText(/Switch to (light|dark) mode/);
    expect(toggle).toBeInTheDocument();
  });

  it('toggles the data-theme attribute when theme button is clicked', () => {
    localStorage.setItem('muara-theme', 'light');
    document.documentElement.setAttribute('data-theme', 'light');
    render(
      <AppShell
        active="ledger"
        onNavigate={vi.fn()}
        showHelp={false}
        onToggleHelp={vi.fn()}
      >
        <div />
      </AppShell>,
    );

    const toggle = screen.getByLabelText(/Switch to dark mode/);
    fireEvent.click(toggle);
    expect(document.documentElement.getAttribute('data-theme')).toBe('dark');
    expect(localStorage.getItem('muara-theme')).toBe('dark');
  });

  it('shows connection status indicator', () => {
    render(
      <AppShell
        active="ledger"
        onNavigate={vi.fn()}
        showHelp={false}
        onToggleHelp={vi.fn()}
        connectionStatus="online"
      >
        <div />
      </AppShell>,
    );

    expect(screen.getByLabelText('Connected')).toBeInTheDocument();
  });

  it('calls onReload when reload button is clicked', () => {
    const onReload = vi.fn();
    render(
      <AppShell
        active="ledger"
        onNavigate={vi.fn()}
        showHelp={false}
        onToggleHelp={vi.fn()}
        onReload={onReload}
      >
        <div />
      </AppShell>,
    );

    fireEvent.click(screen.getByLabelText('Refresh dashboard data'));
    expect(onReload).toHaveBeenCalled();
  });
});
