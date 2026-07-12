interface SkeletonProps {
  rows?: number;
  columns?: number;
}

export function SkeletonRows({ rows = 5, columns = 6 }: SkeletonProps) {
  return (
    <>
      {Array.from({ length: rows }).map((_, r) => (
        <tr key={r}>
          {Array.from({ length: columns }).map((_, c) => (
            <td key={c}>
              <span class="skeleton skeleton-inline" />
            </td>
          ))}
        </tr>
      ))}
    </>
  );
}

export function SkeletonCards({ count = 3 }: { count?: number }) {
  return (
    <div class="grid">
      {Array.from({ length: count }).map((_, i) => (
        <div key={i} class="card card-padded">
          <span class="skeleton skeleton-title" />
          <span class="skeleton skeleton-line" />
          <span class="skeleton skeleton-line skeleton-line-short" />
        </div>
      ))}
    </div>
  );
}

export function Skeleton({ variant = 'card', className = '' }: { variant?: 'card' | 'line' | 'title'; className?: string }) {
  return <span class={`skeleton skeleton-${variant} ${className}`} />;
}
