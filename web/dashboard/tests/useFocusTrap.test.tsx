import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render } from '@testing-library/preact';
import { useFocusTrap } from '../src/hooks/useFocusTrap';

function Trap({ active }: { active: boolean }) {
  const ref = useFocusTrap<HTMLDivElement>(active);
  return (
    <div ref={ref} data-testid="trap">
      <button>first</button>
      <input type="text" />
      <button>last</button>
    </div>
  );
}

describe('useFocusTrap', () => {
  beforeEach(() => {
    document.body.innerHTML = '';
  });

  afterEach(() => {
    document.body.innerHTML = '';
  });

  it('focuses first focusable element when active', () => {
    const { getByText } = render(<Trap active />);
    expect(document.activeElement).toBe(getByText('first'));
  });

  it('cycles focus forward on Tab', () => {
    const { getByTestId, getByText } = render(<Trap active />);
    const container = getByTestId('trap');
    const first = getByText('first');
    const last = getByText('last');
    last.focus();

    const event = new KeyboardEvent('keydown', { key: 'Tab', bubbles: true });
    const preventDefault = vi.spyOn(event, 'preventDefault');
    container.dispatchEvent(event);

    expect(preventDefault).toHaveBeenCalled();
    expect(document.activeElement).toBe(first);
  });

  it('cycles focus backward on Shift+Tab', () => {
    const { getByTestId, getByText } = render(<Trap active />);
    const container = getByTestId('trap');
    const first = getByText('first');
    const last = getByText('last');
    first.focus();

    const event = new KeyboardEvent('keydown', { key: 'Tab', shiftKey: true, bubbles: true });
    const preventDefault = vi.spyOn(event, 'preventDefault');
    container.dispatchEvent(event);

    expect(preventDefault).toHaveBeenCalled();
    expect(document.activeElement).toBe(last);
  });

  it('ignores non-Tab keys', () => {
    const { getByTestId } = render(<Trap active />);
    const container = getByTestId('trap');
    const event = new KeyboardEvent('keydown', { key: 'Escape', bubbles: true });
    const preventDefault = vi.spyOn(event, 'preventDefault');
    container.dispatchEvent(event);

    expect(preventDefault).not.toHaveBeenCalled();
  });
});
