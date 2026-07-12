import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/preact';
import { FilterToolbar } from '../src/components/FilterToolbar';

describe('FilterToolbar', () => {
  it('renders search input', () => {
    render(
      <FilterToolbar
        searchId="test-search"
        query=""
        onQueryChange={vi.fn()}
      />,
    );
    expect(screen.getByLabelText('Search')).toBeInTheDocument();
  });

  it('calls onQueryChange when typing', () => {
    const onQueryChange = vi.fn();
    render(
      <FilterToolbar
        searchId="test-search"
        query=""
        onQueryChange={onQueryChange}
      />,
    );
    const input = screen.getByLabelText('Search');
    fireEvent.input(input, { target: { value: 'pi_' } });
    expect(onQueryChange).toHaveBeenCalledWith('pi_');
  });

  it('renders provider and status filters when provided', () => {
    const providers = new Map([['stripe', 'Stripe']]);
    render(
      <FilterToolbar
        searchId="test-search"
        query=""
        onQueryChange={vi.fn()}
        providerOptions={providers}
        onProviderChange={vi.fn()}
        statusOptions={[{ value: 'paid', label: 'paid' }]}
        onStatusChange={vi.fn()}
      />,
    );
    expect(screen.getByLabelText('Provider filter')).toBeInTheDocument();
    expect(screen.getByLabelText('Status filter')).toBeInTheDocument();
  });

  it('renders active filter chips', () => {
    render(
      <FilterToolbar
        searchId="test-search"
        query="pi_"
        onQueryChange={vi.fn()}
        activeFilters={[{ label: 'Search: pi_', onRemove: vi.fn() }]}
        onResetAll={vi.fn()}
      />,
    );
    expect(screen.getByText('Search: pi_')).toBeInTheDocument();
    expect(screen.getByText('Reset all')).toBeInTheDocument();
  });

  it('calls onResetAll when reset button clicked', () => {
    const onResetAll = vi.fn();
    render(
      <FilterToolbar
        searchId="test-search"
        query="pi_"
        onQueryChange={vi.fn()}
        activeFilters={[{ label: 'Search: pi_', onRemove: vi.fn() }]}
        onResetAll={onResetAll}
      />,
    );
    fireEvent.click(screen.getByText('Reset all'));
    expect(onResetAll).toHaveBeenCalled();
  });
});
