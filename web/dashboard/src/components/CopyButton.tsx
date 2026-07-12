import { useState } from 'preact/hooks';
import { copyToClipboard } from '../api';
import { Button } from './Button';

interface CopyButtonProps {
  text: string;
  label?: string;
  size?: 'sm' | 'md';
}

export function CopyButton({ text, label = '', size = 'sm' }: CopyButtonProps) {
  const [copied, setCopied] = useState(false);

  const handleClick = async () => {
    const ok = await copyToClipboard(text);
    if (!ok) return;
    setCopied(true);
    setTimeout(() => setCopied(false), 1500);
  };

  const actionLabel = label.trim();
  const ariaLabel = copied
    ? 'Copied'
    : actionLabel
      ? `Copy ${actionLabel}`
      : 'Copy';

  return (
    <Button
      variant="ghost"
      size={size}
      icon={copied ? 'check' : 'copy'}
      onClick={handleClick}
      aria-label={ariaLabel}
      title={ariaLabel}
    >
      {copied ? 'Copied' : actionLabel || undefined}
    </Button>
  );
}
