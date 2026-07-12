import { forwardRef } from 'preact/compat';
import { Icon } from './Icon';
import type { IconName } from './Icon';

interface ButtonProps {
  children: preact.ComponentChildren;
  type?: 'button' | 'submit';
  variant?: 'primary' | 'secondary' | 'ghost' | 'danger';
  size?: 'sm' | 'md';
  disabled?: boolean;
  loading?: boolean;
  icon?: IconName;
  title?: string;
  'aria-label'?: string;
  onClick?: () => void;
  className?: string;
}

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(function Button({
  children,
  type = 'button',
  variant = 'primary',
  size = 'md',
  disabled = false,
  loading = false,
  icon,
  title,
  'aria-label': ariaLabel,
  onClick,
  className = '',
}, ref) {
  const isDisabled = disabled || loading;
  return (
    <button
      ref={ref}
      type={type}
      className={`btn btn-${variant} btn-${size} ${className}`}
      disabled={isDisabled}
      title={title}
      aria-label={ariaLabel}
      onClick={onClick}
    >
      {loading && <span class="btn-spinner" aria-hidden="true" />}
      {!loading && icon && <Icon name={icon} size={size === 'sm' ? 14 : 16} />}
      <span>{children}</span>
    </button>
  );
});
