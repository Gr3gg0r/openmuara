import { describe, it, expect } from 'vitest';
import { render, waitFor } from '@testing-library/preact';
import { AnnounceRegion, announce } from '../src/components/Announce';

describe('AnnounceRegion', () => {
  it('renders live region and reflects announced messages', async () => {
    const { container } = render(<AnnounceRegion />);
    const region = container.querySelector('[aria-live="polite"]');
    expect(region).toHaveClass('sr-only');

    announce('Payment received');
    await waitFor(() => {
      expect(container.querySelector('[aria-live="polite"]')).toHaveTextContent('Payment received');
    });
  });

  it('clears callback on unmount', () => {
    const { container, unmount } = render(<AnnounceRegion />);
    unmount();
    announce('ignored');
    expect(document.body.textContent).not.toContain('ignored');
    expect(container.querySelector('[aria-live="polite"]')).toBeNull();
  });
});
