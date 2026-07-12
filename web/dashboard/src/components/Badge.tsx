export interface BadgeProps {
  children: preact.ComponentChildren;
  variant?: 'success' | 'danger' | 'warning' | 'info' | 'neutral' | 'active';
  className?: string;
}

export function Badge({ children, variant = 'neutral', className = '' }: BadgeProps) {
  return <span className={`badge badge-${variant} ${className}`}>{children}</span>;
}
