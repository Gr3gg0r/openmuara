import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/preact';
import { WebhookDetail } from '../src/views/WebhookDetail';
import type { WebhookAttempt } from '../src/types';

describe('WebhookDetail', () => {
  beforeEach(() => {
    document.body.innerHTML = '';
  });

  it('renders back button and calls onBack', () => {
    const onBack = vi.fn();
    const webhook: WebhookAttempt = {
      ref: 'wh_1',
      provider: 'stripe',
      url: 'https://example.com/webhook',
      status: 'delivered',
      attempts: 1,
    };

    globalThis.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: vi.fn().mockResolvedValue({ webhook: { ref: 'wh_1', provider: 'stripe', status: 'delivered' } }),
    }) as unknown as typeof fetch;

    render(<WebhookDetail webhook={webhook} onBack={onBack} />);
    const back = screen.getByText('Back to webhooks');
    expect(back).toBeInTheDocument();
    fireEvent.click(back);
    expect(onBack).toHaveBeenCalled();
  });

  it('shows webhook details', async () => {
    const webhook: WebhookAttempt = {
      ref: 'wh_1',
      provider_name: 'Stripe',
      url: 'https://example.com/webhook',
      status: 'delivered',
      attempts: 2,
    };

    globalThis.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: vi.fn().mockResolvedValue({
        webhook: {
          ref: 'wh_1',
          provider_name: 'Stripe',
          url: 'https://example.com/webhook',
          status: 'delivered',
          attempts: 2,
        },
      }),
    }) as unknown as typeof fetch;

    render(<WebhookDetail webhook={webhook} onBack={vi.fn()} />);
    await waitFor(() => {
      expect(screen.getByText('https://example.com/webhook')).toBeInTheDocument();
    });
    expect(screen.getByText('Webhook detail')).toBeInTheDocument();
    expect(screen.getByText('Stripe', { selector: 'dd' })).toBeInTheDocument();
    expect(screen.getByText('delivered', { selector: 'dd .badge' })).toBeInTheDocument();
    expect(screen.getByText('2', { selector: 'dd' })).toBeInTheDocument();
  });

  it('shows error banner on fetch failure', async () => {
    const webhook: WebhookAttempt = {
      ref: 'wh_1',
      provider: 'stripe',
      status: 'failed',
    };

    globalThis.fetch = vi.fn().mockResolvedValue({
      ok: false,
      status: 500,
      text: vi.fn().mockResolvedValue('server error'),
    }) as unknown as typeof fetch;

    render(<WebhookDetail webhook={webhook} onBack={vi.fn()} />);
    await waitFor(() => {
      expect(screen.getByRole('alert')).toBeInTheDocument();
    });
  });
});
