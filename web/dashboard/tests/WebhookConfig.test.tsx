import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/preact';
import { WebhookConfig } from '../src/components/WebhookConfig';

describe('WebhookConfig', () => {
  beforeEach(() => {
    document.head.innerHTML = '<meta name="csrf-token" content="tok123">';
  });

  afterEach(() => {
    vi.restoreAllMocks();
    document.head.innerHTML = '';
  });

  it('renders the webhook configuration form', async () => {
    const fetchMock = vi
      .fn()
      .mockResolvedValueOnce({
        ok: true,
        json: vi.fn().mockResolvedValue({ url: '', max_retries: 3, targets: {}, events: {} }),
      })
      .mockResolvedValueOnce({
        ok: true,
        json: vi.fn().mockResolvedValue({ enabled: ['stripe'], providers: { stripe: { display_name: 'Stripe' } } }),
      });
    globalThis.fetch = fetchMock as unknown as typeof fetch;

    render(<WebhookConfig />);

    await waitFor(() => {
      expect(screen.getByText('Webhook Configuration')).toBeInTheDocument();
    });
    expect(screen.getByLabelText('Global webhook URL')).toBeInTheDocument();
  });
});
