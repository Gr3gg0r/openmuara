import { Icon } from './Icon';
import { formatDate } from '../api';

export interface TimelineItem {
  status: 'completed' | 'current' | 'pending' | 'error';
  label: string;
  time?: string;
  detail?: string;
}

interface TimelineProps {
  items: TimelineItem[];
}

function statusIcon(status: TimelineItem['status']) {
  switch (status) {
    case 'completed':
      return 'checkCircle';
    case 'current':
      return 'circle';
    case 'error':
      return 'alert';
    default:
      return 'circle';
  }
}

function statusClass(status: TimelineItem['status']) {
  switch (status) {
    case 'completed':
      return 'timeline-item-completed';
    case 'current':
      return 'timeline-item-current';
    case 'error':
      return 'timeline-item-error';
    default:
      return 'timeline-item-pending';
  }
}

export function Timeline({ items }: TimelineProps) {
  if (items.length === 0) return null;

  return (
    <ol class="timeline" role="list" aria-label="Status timeline">
      {items.map((item, index) => (
        <li key={index} class={`timeline-item ${statusClass(item.status)}`}>
          <div class="timeline-marker" aria-hidden="true">
            <Icon name={statusIcon(item.status)} size={16} />
          </div>
          <div class="timeline-content">
            <div class="timeline-title">{item.label}</div>
            {item.time && <div class="timeline-time">{formatDate(item.time)}</div>}
            {item.detail && <div class="timeline-detail">{item.detail}</div>}
          </div>
        </li>
      ))}
    </ol>
  );
}
