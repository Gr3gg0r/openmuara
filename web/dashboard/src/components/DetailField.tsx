import { Badge } from './Badge';
import type { BadgeProps } from './Badge';

interface DetailFieldProps {
  label: string;
  value?: preact.ComponentChildren;
  badge?: { text: string; variant: BadgeProps['variant'] };
  fullWidth?: boolean;
}

export function DetailField({ label, value, badge, fullWidth = false }: DetailFieldProps) {
  return (
    <div class={`detail-field ${fullWidth ? 'detail-field-full' : ''}`}>
      <dt class="detail-field-label">{label}</dt>
      <dd class="detail-field-value">
        {badge ? <Badge variant={badge.variant}>{badge.text}</Badge> : value ?? '-'}
      </dd>
    </div>
  );
}
