import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/preact';
import { Shell } from '../src/components/Shell';

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

describe('Shell', () => {
  beforeEach(() => {
    localStorage.clear();
    document.documentElement.removeAttribute('data-theme');
  });

  it('renders all tabs and children', () => {
    render(
      <Shell
        tab="overview"
        onTabChange={vi.fn()}
        showHelp={false}
        onToggleHelp={vi.fn()}
      >
        <div data-testid="child">content</div>
      </Shell>,
    );

    expect(screen.getByText('Overview')).toBeInTheDocument();
    expect(screen.getByText('Ledger')).toBeInTheDocument();
    expect(screen.getByText('Transactions')).toBeInTheDocument();
    expect(screen.getByText('Webhooks')).toBeInTheDocument();
    expect(screen.getByTestId('child')).toBeInTheDocument();
  });

  it('calls onTabChange when a tab is clicked', () => {
    const onTabChange = vi.fn();
    render(
      <Shell
        tab="overview"
        onTabChange={onTabChange}
        showHelp={false}
        onToggleHelp={vi.fn()}
      >
        <div />
      </Shell>,
    );

    fireEvent.click(screen.getByText('Ledger'));
    expect(onTabChange).toHaveBeenCalledWith('ledger');
  });

  it('renders a theme toggle button', () => {
    render(
      <Shell
        tab="overview"
        onTabChange={vi.fn()}
        showHelp={false}
        onToggleHelp={vi.fn()}
      >
        <div />
      </Shell>,
    );

    const toggle = screen.getByLabelText(/Switch to (light|dark) mode/);
    expect(toggle).toBeInTheDocument();
  });

  it('toggles the data-theme attribute when clicked', () => {
    localStorage.setItem('muara-theme', 'light');
    document.documentElement.setAttribute('data-theme', 'light');
    render(
      <Shell
        tab="overview"
        onTabChange={vi.fn()}
        showHelp={false}
        onToggleHelp={vi.fn()}
      >
        <div />
      </Shell>,
    );

    const toggle = screen.getByLabelText(/Switch to dark mode/);
    fireEvent.click(toggle);
    expect(document.documentElement.getAttribute('data-theme')).toBe('dark');
    expect(localStorage.getItem('muara-theme')).toBe('dark');
  });

  it('shows connection status indicator', () => {
    render(
      <Shell
        tab="overview"
        onTabChange={vi.fn()}
        showHelp={false}
        onToggleHelp={vi.fn()}
        connectionStatus="online"
      >
        <div />
      </Shell>,
    );

    expect(screen.getByLabelText('Connected')).toBeInTheDocument();
  });
});
