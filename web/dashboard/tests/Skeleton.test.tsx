import { describe, it, expect } from 'vitest';
import { render } from '@testing-library/preact';
import { Skeleton, SkeletonRows, SkeletonCards } from '../src/components/Skeleton';

describe('Skeleton', () => {
  it('renders card variant by default', () => {
    const { container } = render(<Skeleton />);
    expect(container.querySelector('.skeleton-card')).toBeInTheDocument();
  });

  it('renders line and title variants', () => {
    const { container: line } = render(<Skeleton variant="line" />);
    expect(line.querySelector('.skeleton-line')).toBeInTheDocument();

    const { container: title } = render(<Skeleton variant="title" className="extra" />);
    expect(title.querySelector('.skeleton-title.extra')).toBeInTheDocument();
  });

  it('renders configured table rows', () => {
    const { container } = render(<SkeletonRows rows={2} columns={3} />);
    expect(container.querySelectorAll('tr')).toHaveLength(2);
    expect(container.querySelectorAll('td')).toHaveLength(6);
  });

  it('renders configured card skeletons', () => {
    const { container } = render(<SkeletonCards count={2} />);
    expect(container.querySelectorAll('.skeleton-title')).toHaveLength(2);
  });
});
