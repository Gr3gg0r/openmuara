import { formatJSON } from '../api';
import { CopyButton } from './CopyButton';

interface CodeBlockProps {
  value: unknown;
  title?: string;
}

export function CodeBlock({ value, title }: CodeBlockProps) {
  const text = typeof value === 'string' ? value : formatJSON(value);

  return (
    <div class="code-block">
      <div class="code-block-header">
        {title && <span class="code-block-title">{title}</span>}
        <CopyButton text={text} />
      </div>
      <pre>{text}</pre>
    </div>
  );
}
