import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { renderHook } from '@testing-library/preact';
import { useUrlStateSynced } from '../src/hooks/useUrlState';

describe('useUrlStateSynced', () => {
  beforeEach(() => {
    window.history.pushState({}, '', '/');
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('reads initial value from URL and calls setter', () => {
    window.history.pushState({}, '', '/?tab=settings');
    const setValue = vi.fn();
    renderHook(() => useUrlStateSynced('tab', 'overview', setValue));
    expect(setValue).toHaveBeenCalledWith('settings');
  });

  it('does not call setter when URL matches current value', () => {
    window.history.pushState({}, '', '/?tab=overview');
    const setValue = vi.fn();
    renderHook(() => useUrlStateSynced('tab', 'overview', setValue));
    expect(setValue).not.toHaveBeenCalled();
  });

  it('updates URL when value changes', () => {
    const setValue = vi.fn();
    const { rerender } = renderHook(
      ({ value }) => useUrlStateSynced('tab', value, setValue, 'overview'),
      { initialProps: { value: 'overview' } },
    );
    rerender({ value: 'settings' });
    expect(window.location.search).toContain('tab=settings');
  });

  it('removes param when value matches default', () => {
    window.history.pushState({}, '', '/?tab=settings');
    const setValue = vi.fn();
    const { rerender } = renderHook(
      ({ value }) => useUrlStateSynced('tab', value, setValue, 'overview'),
      { initialProps: { value: 'settings' } },
    );
    rerender({ value: 'overview' });
    expect(window.location.search).not.toContain('tab=settings');
    expect(window.location.search).not.toContain('tab=overview');
  });
});
