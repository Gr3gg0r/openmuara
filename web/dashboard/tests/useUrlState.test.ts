import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { renderHook, act } from '@testing-library/preact';
import { useUrlState } from '../src/hooks/useUrlState';

describe('useUrlState', () => {
  beforeEach(() => {
    window.history.pushState({}, '', '/');
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('reads initial value from URL', () => {
    window.history.pushState({}, '', '/?q=hello');
    const { result } = renderHook(() => useUrlState('q'));
    expect(result.current[0]).toBe('hello');
  });

  it('falls back to default value when URL param is absent', () => {
    const { result } = renderHook(() => useUrlState('sort', 'time'));
    expect(result.current[0]).toBe('time');
  });

  it('updates URL when value changes', () => {
    const { result } = renderHook(() => useUrlState('provider'));
    act(() => {
      result.current[1]('stripe');
    });
    expect(result.current[0]).toBe('stripe');
    expect(window.location.search).toContain('provider=stripe');
  });

  it('removes URL param when value matches default', () => {
    window.history.pushState({}, '', '/?sort=time');
    const { result } = renderHook(() => useUrlState('sort', 'time'));
    act(() => {
      result.current[1]('time');
    });
    expect(window.location.search).not.toContain('sort=time');
  });
});
