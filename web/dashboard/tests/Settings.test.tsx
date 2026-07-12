import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/preact';
import { SettingsView } from '../src/views/Settings';

const providersResponse = {
  active: 'fawry',
  enabled: ['fawry', 'default'],
  available: ['fawry', 'stripe', 'default'],
  providers: {
    fawry: {
      name: 'fawry',
      display_name: 'Fawry',
      description: 'Egyptian payment gateway',
      enabled: true,
      active: true,
      category: 'regional',
      real_providers: ['Fawry'],
      version: 'v1',
      versions: ['v1', 'v2'],
      env_vars: ['MUARA_FAWRY_MERCHANT_CODE'],
    },
    stripe: {
      name: 'stripe',
      display_name: 'Stripe',
      description: 'Stripe Checkout emulation',
      enabled: false,
      active: false,
      category: 'card',
      real_providers: ['Stripe'],
      version: 'v1',
      versions: [],
      env_vars: ['MUARA_STRIPE_SECRET_KEY'],
    },
    default: {
      name: 'default',
      display_name: 'Default / DIY',
      description: 'Minimal provider',
      enabled: true,
      active: false,
      category: 'diy',
      real_providers: [],
      version: 'v1',
      versions: [],
      env_vars: [],
    },
  },
};

function createFetchMock(cleanOk = true, cleanError?: string) {
  return vi.fn().mockImplementation((url: string) => {
    if (url.includes('/_admin/providers')) {
      return Promise.resolve({
        ok: true,
        json: vi.fn().mockResolvedValue(providersResponse),
      });
    }
    if (url.includes('/_admin/clean')) {
      if (cleanError) {
        return Promise.resolve({
          ok: false,
          status: 500,
          text: vi.fn().mockResolvedValue(cleanError),
        });
      }
      return Promise.resolve({ ok: cleanOk, status: 200, text: vi.fn().mockResolvedValue('') });
    }
    return Promise.resolve({ ok: false, status: 404, text: vi.fn().mockResolvedValue('not found') });
  }) as unknown as typeof fetch;
}

describe('SettingsView', () => {
  beforeEach(() => {
    document.head.innerHTML = '<meta name="csrf-token" content="tok123"><meta name="muara-role" content="admin">';
    document.body.innerHTML = '';
    window.history.replaceState({}, '', '/');
  });

  afterEach(() => {
    vi.restoreAllMocks();
    document.head.innerHTML = '';
  });

  it('renders provider cards', async () => {
    globalThis.fetch = createFetchMock();
    render(<SettingsView />);

    await waitFor(() => {
      expect(screen.getByText('Fawry')).toBeInTheDocument();
      expect(screen.getByText('Stripe')).toBeInTheDocument();
      expect(screen.getByText('Default / DIY')).toBeInTheDocument();
    }, { timeout: 3000 });
  });

  it('shows enable and active badges', async () => {
    globalThis.fetch = createFetchMock();
    render(<SettingsView />);

    await waitFor(() => {
      expect(screen.getByText('Active')).toBeInTheDocument();
      expect(screen.getAllByText('Enabled').length).toBeGreaterThanOrEqual(1);
      expect(screen.getAllByText('Disabled').length).toBeGreaterThanOrEqual(1);
    }, { timeout: 3000 });
  });

  it('calls onShowProvider when configure button clicked', async () => {
    globalThis.fetch = createFetchMock();
    const onShowProvider = vi.fn();
    render(<SettingsView onShowProvider={onShowProvider} />);

    await waitFor(() => {
      expect(screen.getByText('Fawry')).toBeInTheDocument();
    }, { timeout: 3000 });

    fireEvent.click(screen.getByLabelText('Configure Fawry'));
    expect(onShowProvider).toHaveBeenCalledWith('fawry');
  });

  it('shows real providers list', async () => {
    globalThis.fetch = createFetchMock();
    render(<SettingsView />);

    await waitFor(() => {
      expect(screen.getByText(/Emulates: Fawry/)).toBeInTheDocument();
      expect(screen.getByText(/Emulates: Stripe/)).toBeInTheDocument();
    }, { timeout: 3000 });
  });

  it('shows clear data button for admin', async () => {
    globalThis.fetch = createFetchMock();
    render(<SettingsView />);

    await waitFor(() => {
      expect(screen.getByRole('button', { name: /Clear local data/i })).toBeInTheDocument();
    }, { timeout: 3000 });
  });

  it('hides clear data button for viewer', async () => {
    document.head.innerHTML = '<meta name="csrf-token" content="tok123"><meta name="muara-role" content="viewer">';
    globalThis.fetch = createFetchMock();
    render(<SettingsView />);

    await waitFor(() => {
      expect(screen.queryByRole('button', { name: /Clear local data/i })).not.toBeInTheDocument();
    }, { timeout: 3000 });
  });

  it('opens confirmation dialog and clears data', async () => {
    globalThis.fetch = createFetchMock();
    const reload = vi.fn();
    Object.defineProperty(window, 'location', {
      value: { href: 'http://localhost/', reload },
      writable: true,
    });

    render(<SettingsView />);

    await waitFor(() => {
      expect(screen.getByRole('button', { name: /Clear local data/i })).toBeInTheDocument();
    }, { timeout: 3000 });

    fireEvent.click(screen.getByRole('button', { name: /Clear local data/i }));
    expect(screen.getByRole('alertdialog')).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: /Yes, clear data/i }));

    await waitFor(() => {
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/_admin/clean'),
        expect.objectContaining({ method: 'POST' }),
      );
    });
    expect(reload).toHaveBeenCalled();
  });

  it('shows error when clear fails', async () => {
    globalThis.fetch = createFetchMock(true, '{"error":"database locked"}');
    render(<SettingsView />);

    await waitFor(() => {
      expect(screen.getByRole('button', { name: /Clear local data/i })).toBeInTheDocument();
    }, { timeout: 3000 });

    fireEvent.click(screen.getByRole('button', { name: /Clear local data/i }));
    fireEvent.click(screen.getByRole('button', { name: /Yes, clear data/i }));

    await waitFor(() => {
      expect(screen.getByText(/database locked/)).toBeInTheDocument();
    }, { timeout: 3000 });
  });
});
