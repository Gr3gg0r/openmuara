import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/preact';
import { FailedWebhookAlert } from '../src/components/FailedWebhookAlert';

describe('FailedWebhookAlert', () => {
  beforeEach(() => {
    sessionStorage.clear();
  });

  afterEach(() => {
    vi.restoreAllMocks();
    sessionStorage.clear();
  });

  it('renders alert when failed webhooks exist', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      json: vi.fn().mockResolvedValue({ results: [{ id: 'w1' }] }),
    });
    globalThis.fetch = fetchMock as unknown as typeof fetch;

    render(<FailedWebhookAlert onShowWebhooks={vi.fn()} enabled />);
    await waitFor(() => {
      expect(screen.getByRole('alert')).toHaveTextContent('Failed webhook detected');
    });
  });

  it('does not fetch when disabled', async () => {
    const fetchMock = vi.fn();
    globalThis.fetch = fetchMock as unknown as typeof fetch;

    render(<FailedWebhookAlert onShowWebhooks={vi.fn()} enabled={false} />);
    await new Promise((r) => setTimeout(r, 50));
    expect(fetchMock).not.toHaveBeenCalled();
  });

  it('does not render when already dismissed', async () => {
    sessionStorage.setItem('muara_failed_webhook_dismissed', 'true');
    const fetchMock = vi.fn();
    globalThis.fetch = fetchMock as unknown as typeof fetch;

    const { container } = render(<FailedWebhookAlert onShowWebhooks={vi.fn()} enabled />);
    await new Promise((r) => setTimeout(r, 50));
    expect(container.firstChild).toBeNull();
  });

  it('calls onShowWebhooks when link is clicked', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      json: vi.fn().mockResolvedValue({ results: [{ id: 'w1' }] }),
    });
    globalThis.fetch = fetchMock as unknown as typeof fetch;

    const onShowWebhooks = vi.fn();
    render(<FailedWebhookAlert onShowWebhooks={onShowWebhooks} enabled />);
    await waitFor(() => screen.getByRole('alert'));
    fireEvent.click(screen.getByRole('button', { name: 'Webhooks' }));
    expect(onShowWebhooks).toHaveBeenCalled();
  });

  it('dismisses alert and stores flag', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      json: vi.fn().mockResolvedValue({ results: [{ id: 'w1' }] }),
    });
    globalThis.fetch = fetchMock as unknown as typeof fetch;

    const { container } = render(<FailedWebhookAlert onShowWebhooks={vi.fn()} enabled />);
    await waitFor(() => screen.getByRole('alert'));
    fireEvent.click(screen.getByRole('button', { name: 'Dismiss' }));
    expect(sessionStorage.getItem('muara_failed_webhook_dismissed')).toBe('true');
    expect(container.firstChild).toBeNull();
  });

  it('handles fetch error gracefully', async () => {
    const fetchMock = vi.fn().mockRejectedValue(new Error('network'));
    globalThis.fetch = fetchMock as unknown as typeof fetch;

    const { container } = render(<FailedWebhookAlert onShowWebhooks={vi.fn()} enabled />);
    await new Promise((r) => setTimeout(r, 50));
    expect(container.firstChild).toBeNull();
  });
});
