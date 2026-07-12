import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { renderHook } from '@testing-library/preact';
import { usePolling } from '../src/hooks/usePolling';

describe('usePolling', () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('calls callback immediately and on interval', () => {
    const cb = vi.fn();
    renderHook(() => usePolling(cb, 1000));
    expect(cb).toHaveBeenCalledTimes(1);
    vi.advanceTimersByTime(1000);
    expect(cb).toHaveBeenCalledTimes(2);
  });

  it('does not call callback when disabled', () => {
    const cb = vi.fn();
    renderHook(() => usePolling(cb, 1000, false));
    expect(cb).not.toHaveBeenCalled();
    vi.advanceTimersByTime(2000);
    expect(cb).not.toHaveBeenCalled();
  });

  it('pauses when document is hidden and resumes on visible', () => {
    const cb = vi.fn();
    Object.defineProperty(document, 'hidden', { value: false, writable: true, configurable: true });

    renderHook(() => usePolling(cb, 1000));
    expect(cb).toHaveBeenCalledTimes(1);

    vi.advanceTimersByTime(1000);
    expect(cb).toHaveBeenCalledTimes(2);

    (document as { hidden: boolean }).hidden = true;
    document.dispatchEvent(new Event('visibilitychange'));
    vi.advanceTimersByTime(2000);
    expect(cb).toHaveBeenCalledTimes(2);

    (document as { hidden: boolean }).hidden = false;
    document.dispatchEvent(new Event('visibilitychange'));
    expect(cb).toHaveBeenCalledTimes(3);
  });
});
