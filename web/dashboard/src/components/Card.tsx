interface CardProps {
  children: preact.ComponentChildren;
  className?: string;
  padded?: boolean;
}

export function Card({ children, className = '', padded = true }: CardProps) {
  return <div className={`card ${padded ? 'card-padded' : ''} ${className}`}>{children}</div>;
}
