import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/preact';
import { EmptyState } from '../src/components/EmptyState';

describe('EmptyState', () => {
  it('renders title and description', () => {
    render(<EmptyState title="No data" description="Nothing to show" />);
    expect(screen.getByText('No data')).toBeInTheDocument();
    expect(screen.getByText('Nothing to show')).toBeInTheDocument();
  });

  it('renders default empty icon', () => {
    render(<EmptyState title="No data" />);
    expect(document.querySelector('svg')).toBeInTheDocument();
  });

  it('calls action on button click', () => {
    const onClick = vi.fn();
    render(<EmptyState title="No data" action={{ label: 'Create', onClick }} />);
    fireEvent.click(screen.getByRole('button', { name: 'Create' }));
    expect(onClick).toHaveBeenCalled();
  });
});
