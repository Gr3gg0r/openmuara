import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/preact';
import { WebhooksView } from '../src/views/Webhooks';
import type { WebhookAttempt } from '../src/types';

const providersResponse = {
  available: ['stripe', 'fawry'],
  providers: {
    stripe: { display_name: 'Stripe' },
    fawry: { display_name: 'Fawry' },
  },
};

const webhooksResponse = {
  results: [
    { ref: 'wh_stripe_1', provider_name: 'stripe', url: 'https://example.com/stripe', status: 'delivered', attempts: 1, last_error: '' },
    { ref: 'wh_fawry_1', provider_name: 'fawry', url: 'https://example.com/fawry', status: 'failed', attempts: 3, last_error: 'timeout' },
  ] as WebhookAttempt[],
};

function createFetchMock() {
  return vi.fn().mockImplementation((url: string) => {
    if (url.includes('/_admin/config/webhooks')) {
      return Promise.resolve({
        ok: true,
        json: vi.fn().mockResolvedValue({ url: '', max_retries: 3, targets: {}, events: {} }),
      });
    }
    if (url.includes('/_admin/providers')) {
      return Promise.resolve({
        ok: true,
        json: vi.fn().mockResolvedValue(providersResponse),
      });
    }
    if (url.includes('/_admin/webhooks')) {
      const params = new URL(url, window.location.href).searchParams;
      const provider = params.get('provider');
      const status = params.get('status');
      let results = webhooksResponse.results;
      if (provider) {
        results = results.filter((wh) => wh.provider_name === provider);
      }
      if (status) {
        results = results.filter((wh) => wh.status === status);
      }
      return Promise.resolve({
        ok: true,
        json: vi.fn().mockResolvedValue({ results }),
      });
    }
    return Promise.resolve({ ok: false, status: 404, text: vi.fn().mockResolvedValue('not found') });
  }) as unknown as typeof fetch;
}

async function waitForWebhooks() {
  await waitFor(() => {
    expect(screen.getByText('wh_stripe_1')).toBeInTheDocument();
    expect(screen.getByText('wh_fawry_1')).toBeInTheDocument();
  }, { timeout: 3000 });
}

function changeSelect(label: string, value: string) {
  const select = screen.getByLabelText(label) as HTMLSelectElement;
  select.value = value;
  select.dispatchEvent(new Event('change', { bubbles: true }));
}

describe('WebhooksView', () => {
  beforeEach(() => {
    document.head.innerHTML = '<meta name="csrf-token" content="tok123">';
    document.body.innerHTML = '';
    window.history.replaceState({}, '', '/');
  });

  afterEach(() => {
    vi.restoreAllMocks();
    document.head.innerHTML = '';
  });

  it('renders webhook attempts table', async () => {
    globalThis.fetch = createFetchMock();

    render(<WebhooksView />);

    await waitForWebhooks();
    expect(screen.getByText('wh_stripe_1')).toBeInTheDocument();
    expect(screen.getByText('wh_fawry_1')).toBeInTheDocument();
  });

  it('calls onShowDetail when clicking a reference', async () => {
    globalThis.fetch = createFetchMock();
    const onShowDetail = vi.fn();

    render(<WebhooksView onShowDetail={onShowDetail} />);

    await waitForWebhooks();

    fireEvent.click(screen.getByText('wh_stripe_1'));
    expect(onShowDetail).toHaveBeenCalledWith(expect.objectContaining({ ref: 'wh_stripe_1' }));
  });

  it('filters attempts by search query', async () => {
    globalThis.fetch = createFetchMock();

    render(<WebhooksView />);

    await waitForWebhooks();

    const search = screen.getByLabelText('Search');
    fireEvent.input(search, { target: { value: 'fawry' } });

    await waitFor(() => {
      expect(screen.queryAllByText('wh_stripe_1')).toHaveLength(0);
      expect(screen.getByText('wh_fawry_1')).toBeInTheDocument();
    }, { timeout: 3000 });
  });

  it('filters attempts by status', async () => {
    globalThis.fetch = createFetchMock();

    render(<WebhooksView />);

    await waitForWebhooks();

    changeSelect('Status filter', 'failed');

    await waitFor(() => {
      expect(screen.queryAllByText('wh_stripe_1')).toHaveLength(0);
      expect(screen.getByText('wh_fawry_1')).toBeInTheDocument();
    }, { timeout: 3000 });
  });

  it('filters attempts by provider', async () => {
    globalThis.fetch = createFetchMock();

    render(<WebhooksView />);

    await waitForWebhooks();

    changeSelect('Provider filter', 'fawry');

    await waitFor(() => {
      expect(screen.queryAllByText('wh_stripe_1')).toHaveLength(0);
      expect(screen.getByText('wh_fawry_1')).toBeInTheDocument();
    }, { timeout: 3000 });
  });

  it('resets filters when reset all is clicked', async () => {
    globalThis.fetch = createFetchMock();

    render(<WebhooksView />);

    await waitForWebhooks();

    fireEvent.input(screen.getByLabelText('Search'), { target: { value: 'fawry' } });
    await waitFor(() => {
      expect(screen.queryAllByText('wh_stripe_1')).toHaveLength(0);
    }, { timeout: 3000 });

    fireEvent.click(screen.getByText('Reset all'));
    await waitFor(() => {
      expect(screen.getByText('wh_stripe_1')).toBeInTheDocument();
      expect(screen.getByText('wh_fawry_1')).toBeInTheDocument();
    }, { timeout: 3000 });
  });
});
