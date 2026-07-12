import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/preact';
import { CodeBlock } from '../src/components/CodeBlock';

describe('CodeBlock', () => {
  it('renders string value directly', () => {
    render(<CodeBlock value="hello world" />);
    expect(screen.getByText('hello world')).toBeInTheDocument();
  });

  it('renders formatted JSON for objects', () => {
    render(<CodeBlock value={{ status: 'ok' }} title="Response" />);
    expect(screen.getByText((text) => text.includes('status') && text.includes('ok'))).toBeInTheDocument();
    expect(screen.getByText('Response')).toBeInTheDocument();
  });

  it('omits title when not provided', () => {
    const { container } = render(<CodeBlock value={42} />);
    expect(container.querySelector('.code-block-title')).toBeNull();
    expect(screen.getByText('42')).toBeInTheDocument();
  });
});
