import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/preact';
import { LedgerDetail } from '../src/views/LedgerDetail';
import type { LedgerEvent } from '../src/types';

describe('LedgerDetail', () => {
  beforeEach(() => {
    document.body.innerHTML = '';
  });

  it('renders back button and calls onBack', () => {
    const onBack = vi.fn();
    const event: LedgerEvent = {
      id: 'tx1',
      reference: 'tx1',
      type: 'transaction',
      provider: 'stripe',
      status: 'paid',
    };

    globalThis.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: vi.fn().mockResolvedValue({ transaction: { provider: 'stripe', amount: 10, currency: 'USD', status: 'paid' } }),
    }) as unknown as typeof fetch;

    render(<LedgerDetail event={event} onBack={onBack} />);
    const back = screen.getByText('Back to ledger');
    expect(back).toBeInTheDocument();
    fireEvent.click(back);
    expect(onBack).toHaveBeenCalled();
  });

  it('shows transaction details', async () => {
    const event: LedgerEvent = {
      id: 'tx1',
      reference: 'tx1',
      type: 'transaction',
      provider: 'stripe',
      status: 'paid',
    };

    globalThis.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: vi.fn().mockResolvedValue({
        transaction: {
          provider: 'stripe',
          amount: 10,
          currency: 'USD',
          status: 'paid',
          reference: 'tx1',
        },
      }),
    }) as unknown as typeof fetch;

    render(<LedgerDetail event={event} onBack={vi.fn()} />);
    await waitFor(() => {
      expect(screen.getByText('10.00 USD')).toBeInTheDocument();
    });
  });

  it('shows error banner on fetch failure', async () => {
    const event: LedgerEvent = {
      id: 'tx1',
      reference: 'tx1',
      type: 'transaction',
      provider: 'stripe',
      status: 'paid',
    };

    globalThis.fetch = vi.fn().mockResolvedValue({
      ok: false,
      status: 500,
      text: vi.fn().mockResolvedValue('server error'),
    }) as unknown as typeof fetch;

    render(<LedgerDetail event={event} onBack={vi.fn()} />);
    await waitFor(() => {
      expect(screen.getByRole('alert')).toBeInTheDocument();
    });
  });
});
