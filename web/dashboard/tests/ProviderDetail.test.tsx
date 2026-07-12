import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/preact';
import { ProviderDetail } from '../src/views/ProviderDetail';

const fawryDetailResponse = {
  name: 'fawry',
  display_name: 'Fawry',
  description: 'Egyptian payment gateway',
  enabled: true,
  active: true,
  category: 'regional',
  real_providers: ['Fawry'],
  version: 'v1',
  versions: ['v1', 'v2'],
  version_details: {
    v1: { base_url: 'http://127.0.0.1:9000/fawry/v1', sample_route: '/fawry/v1/charge' },
    v2: { base_url: 'http://127.0.0.1:9000/fawry/v2', sample_route: '/fawry/v2/charge' },
  },
  env_vars: ['MUARA_FAWRY_MERCHANT_CODE', 'MUARA_FAWRY_WEBHOOK_SECRET'],
  docs_path: '/docs/providers/fawry.md',
  sample_method: 'POST',
  sample_route: '/fawry/charge',
  base_url: 'http://127.0.0.1:9000/fawry/v1',
  webhook_target_url: 'https://example.com/fawry-webhook',
};

function createFetchMock() {
  return vi.fn().mockImplementation((url: string, init?: RequestInit) => {
    if (url.includes('/_admin/providers/fawry')) {
      return Promise.resolve({
        ok: true,
        json: vi.fn().mockResolvedValue(fawryDetailResponse),
      });
    }
    if (url.includes('/_admin/config') && (!init || init.method === 'GET')) {
      return Promise.resolve({
        ok: true,
        json: vi.fn().mockResolvedValue({
          server: { host: '127.0.0.1', port: 9000 },
          webhook: { url: '', max_retries: 3, targets: { fawry: 'https://example.com/fawry-webhook' }, events: {} },
        }),
      });
    }
    if (url.includes('/_admin/config/providers') && init?.method === 'PATCH') {
      return Promise.resolve({ ok: true, json: vi.fn().mockResolvedValue({ status: 'saved' }) });
    }
    if (url.includes('/_admin/config/webhooks') && init?.method === 'PATCH') {
      return Promise.resolve({ ok: true, json: vi.fn().mockResolvedValue({ status: 'saved' }) });
    }
    return Promise.resolve({ ok: false, status: 404, text: vi.fn().mockResolvedValue('not found') });
  }) as unknown as typeof fetch;
}

describe('ProviderDetail', () => {
  beforeEach(() => {
    document.head.innerHTML = '<meta name="csrf-token" content="tok123">';
    document.body.innerHTML = '';
    window.history.replaceState({}, '', '/');
  });

  afterEach(() => {
    vi.restoreAllMocks();
    document.head.innerHTML = '';
  });

  it('renders provider details with version tabs', async () => {
    globalThis.fetch = createFetchMock();
    render(<ProviderDetail name="fawry" onBack={vi.fn()} />);

    await waitFor(() => {
      expect(screen.getByText('Fawry')).toBeInTheDocument();
      expect(screen.getByText('Egyptian payment gateway')).toBeInTheDocument();
      expect(screen.getByRole('tab', { name: 'v1' })).toBeInTheDocument();
      expect(screen.getByRole('tab', { name: 'v2' })).toBeInTheDocument();
      expect(screen.getByText('http://127.0.0.1:9000/fawry/v1')).toBeInTheDocument();
    }, { timeout: 3000 });
  });

  it('calls onBack when back button clicked', async () => {
    globalThis.fetch = createFetchMock();
    const onBack = vi.fn();
    render(<ProviderDetail name="fawry" onBack={onBack} />);

    await waitFor(() => {
      expect(screen.getByText('Back to Settings')).toBeInTheDocument();
    }, { timeout: 3000 });

    fireEvent.click(screen.getByText('Back to Settings'));
    expect(onBack).toHaveBeenCalled();
  });

  it('switches base url when version tab clicked', async () => {
    globalThis.fetch = createFetchMock();
    render(<ProviderDetail name="fawry" onBack={vi.fn()} />);

    await waitFor(() => {
      expect(screen.getByText('http://127.0.0.1:9000/fawry/v1')).toBeInTheDocument();
    }, { timeout: 3000 });

    fireEvent.click(screen.getByRole('tab', { name: 'v2' }));

    await waitFor(() => {
      expect(screen.getByText('http://127.0.0.1:9000/fawry/v2')).toBeInTheDocument();
    }, { timeout: 3000 });
  });

  it('lists environment variables', async () => {
    globalThis.fetch = createFetchMock();
    render(<ProviderDetail name="fawry" onBack={vi.fn()} />);

    await waitFor(() => {
      expect(screen.getByText('MUARA_FAWRY_MERCHANT_CODE')).toBeInTheDocument();
      expect(screen.getByText('MUARA_FAWRY_WEBHOOK_SECRET')).toBeInTheDocument();
    }, { timeout: 3000 });
  });

  it('saves webhook target url', async () => {
    globalThis.fetch = createFetchMock();
    render(<ProviderDetail name="fawry" onBack={vi.fn()} />);

    await waitFor(() => {
      expect(screen.getByLabelText('Target URL')).toBeInTheDocument();
    }, { timeout: 3000 });

    fireEvent.input(screen.getByLabelText('Target URL'), {
      target: { value: 'https://example.com/new-fawry-webhook' },
    });
    fireEvent.click(screen.getByText('Save webhook target'));

    await waitFor(() => {
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/_admin/config/webhooks'),
        expect.objectContaining({ method: 'PATCH' }),
      );
    }, { timeout: 3000 });
  });
});
