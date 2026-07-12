import { Icon } from './Icon';
import type { IconName } from './Icon';

interface EmptyStateProps {
  title: string;
  description?: string;
  icon?: IconName;
  action?: {
    label: string;
    onClick: () => void;
  };
}

export function EmptyState({ title, description, icon = 'empty', action }: EmptyStateProps) {
  return (
    <div class="empty-state">
      <div class="empty-state-icon">
        <Icon name={icon} size={48} />
      </div>
      <div class="empty-state-title">{title}</div>
      {description && <div class="empty-state-desc">{description}</div>}
      {action && (
        <button class="btn btn-primary mt-3" onClick={action.onClick}>
          {action.label}
        </button>
      )}
    </div>
  );
}
