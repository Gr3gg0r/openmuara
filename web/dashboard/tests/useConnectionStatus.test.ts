import { describe, it, expect, vi, afterEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/preact';
import { useConnectionStatus } from '../src/hooks/useConnectionStatus';

describe('useConnectionStatus', () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('reports online when ping succeeds', async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({ ok: true }) as unknown as typeof fetch;

    const { result } = renderHook(() => useConnectionStatus('/_admin/config'));
    await waitFor(() => expect(result.current).toBe('online'), { timeout: 2000 });
  });

  it('reports offline when ping fails', async () => {
    globalThis.fetch = vi.fn().mockRejectedValue(new Error('network error')) as unknown as typeof fetch;

    const { result } = renderHook(() => useConnectionStatus('/_admin/config'));
    await waitFor(() => expect(result.current).toBe('offline'), { timeout: 2000 });
  });

  it('reacts to browser offline event', async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({ ok: true }) as unknown as typeof fetch;

    const { result } = renderHook(() => useConnectionStatus('/_admin/config'));
    await waitFor(() => expect(result.current).toBe('online'), { timeout: 2000 });

    act(() => {
      window.dispatchEvent(new Event('offline'));
    });
    await waitFor(() => expect(result.current).toBe('offline'), { timeout: 1000 });
  });
});
