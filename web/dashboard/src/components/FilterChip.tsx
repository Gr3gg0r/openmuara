import { Icon } from './Icon';

interface FilterChipProps {
  label: string;
  onRemove: () => void;
}

export function FilterChip({ label, onRemove }: FilterChipProps) {
  return (
    <span class="filter-chip">
      <span>{label}</span>
      <button
        type="button"
        class="filter-chip-remove"
        onClick={onRemove}
        aria-label={`Remove ${label} filter`}
        title={`Remove ${label} filter`}
      >
        <Icon name="close" size={12} />
      </button>
    </span>
  );
}
