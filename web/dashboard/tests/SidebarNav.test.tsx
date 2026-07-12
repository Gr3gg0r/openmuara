import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/preact';
import { SidebarNav, type NavItem } from '../src/components/SidebarNav';

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

describe('SidebarNav', () => {
  beforeEach(() => {
    localStorage.clear();
    document.documentElement.removeAttribute('data-theme');
  });

  it('renders three navigation items', () => {
    render(<SidebarNav active="ledger" onNavigate={vi.fn()} />);

    expect(screen.getByTestId('nav-ledger')).toBeInTheDocument();
    expect(screen.getByTestId('nav-webhooks')).toBeInTheDocument();
    expect(screen.getByTestId('nav-settings')).toBeInTheDocument();
  });

  it('marks the active item with aria-current', () => {
    render(<SidebarNav active="webhooks" onNavigate={vi.fn()} />);

    expect(screen.getByTestId('nav-webhooks')).toHaveAttribute('aria-current', 'page');
    expect(screen.getByTestId('nav-ledger')).not.toHaveAttribute('aria-current');
    expect(screen.getByTestId('nav-settings')).not.toHaveAttribute('aria-current');
  });

  it('calls onNavigate with the selected item', () => {
    const onNavigate = vi.fn();
    render(<SidebarNav active="ledger" onNavigate={onNavigate} />);

    fireEvent.click(screen.getByTestId('nav-settings'));
    expect(onNavigate).toHaveBeenCalledWith('settings');
  });

  it.each(['ledger', 'webhooks', 'settings'] as NavItem[])('displays label for %s', (item) => {
    render(<SidebarNav active={item} onNavigate={vi.fn()} />);

    const label = item.charAt(0).toUpperCase() + item.slice(1);
    expect(screen.getByText(label)).toBeInTheDocument();
  });

  it('hides settings for viewer role', () => {
    render(<SidebarNav active="ledger" onNavigate={vi.fn()} role="viewer" />);

    expect(screen.getByTestId('nav-ledger')).toBeInTheDocument();
    expect(screen.getByTestId('nav-webhooks')).toBeInTheDocument();
    expect(screen.queryByTestId('nav-settings')).not.toBeInTheDocument();
  });
});
