import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/preact';
import { CommandPalette } from '../src/components/CommandPalette';

describe('CommandPalette', () => {
  beforeEach(() => {
    document.body.innerHTML = '';
  });

  it('does not render when closed', () => {
    render(
      <CommandPalette
        open={false}
        onClose={vi.fn()}
        onNavigate={vi.fn()}
        onReload={vi.fn()}
        onFocusSearch={vi.fn()}
      />,
    );
    expect(screen.queryByLabelText('Command palette')).not.toBeInTheDocument();
  });

  it('renders navigation commands when open', () => {
    render(
      <CommandPalette
        open
        onClose={vi.fn()}
        onNavigate={vi.fn()}
        onReload={vi.fn()}
        onFocusSearch={vi.fn()}
      />,
    );
    expect(screen.getByLabelText('Command palette')).toBeInTheDocument();
    expect(screen.getByRole('option', { name: /Go to Ledger/ })).toBeInTheDocument();
    expect(screen.getByRole('option', { name: /Go to Webhooks/ })).toBeInTheDocument();
    expect(screen.getByRole('option', { name: /Go to Settings/ })).toBeInTheDocument();
  });

  it('filters commands by query', () => {
    render(
      <CommandPalette
        open
        onClose={vi.fn()}
        onNavigate={vi.fn()}
        onReload={vi.fn()}
        onFocusSearch={vi.fn()}
      />,
    );
    const input = screen.getByLabelText('Command palette');
    fireEvent.input(input, { target: { value: 'reload' } });
    expect(screen.getByRole('option', { name: /Reload dashboard data/ })).toBeInTheDocument();
    expect(screen.queryByRole('option', { name: /Go to Ledger/ })).not.toBeInTheDocument();
  });

  it('calls onClose when Escape is pressed', () => {
    const onClose = vi.fn();
    render(
      <CommandPalette
        open
        onClose={onClose}
        onNavigate={vi.fn()}
        onReload={vi.fn()}
        onFocusSearch={vi.fn()}
      />,
    );
    fireEvent.keyDown(document, { key: 'Escape' });
    expect(onClose).toHaveBeenCalled();
  });
});
