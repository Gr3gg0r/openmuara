import { Button } from './Button';
import { FilterChip } from './FilterChip';
import { Input } from './Input';
import { Select } from './Select';

interface SelectOption {
  value: string;
  label: string;
}

interface Filter {
  label: string;
  onRemove: () => void;
}

interface FilterToolbarProps {
  searchId: string;
  query: string;
  onQueryChange: (value: string) => void;
  providerOptions?: Map<string, string>;
  provider?: string;
  onProviderChange?: (value: string) => void;
  statusOptions?: SelectOption[];
  status?: string;
  onStatusChange?: (value: string) => void;
  sortValue?: string;
  sortOptions?: SelectOption[];
  onSortChange?: (value: string) => void;
  onRefresh?: () => void;
  lastRefresh?: string;
  activeFilters?: Filter[];
  onResetAll?: () => void;
}

export function FilterToolbar({
  searchId,
  query,
  onQueryChange,
  providerOptions,
  provider = '',
  onProviderChange,
  statusOptions,
  status = '',
  onStatusChange,
  sortValue,
  sortOptions,
  onSortChange,
  onRefresh,
  lastRefresh,
  activeFilters = [],
  onResetAll,
}: FilterToolbarProps) {
  return (
    <div class="card card-padded mb-4">
      <div class="toolbar">
        <Input
          id={searchId}
          type="search"
          value={query}
          placeholder="Search reference, provider, status..."
          icon="search"
          clearable
          label="Search"
          onInput={onQueryChange}
          onClear={() => onQueryChange('')}
        />
        {providerOptions && onProviderChange && (
          <Select value={provider} onChange={onProviderChange} label="Provider filter">
            <option value="">All providers</option>
            {Array.from(providerOptions.entries()).map(([value, label]) => (
              <option key={value} value={value}>
                {label}
              </option>
            ))}
          </Select>
        )}
        {statusOptions && onStatusChange && (
          <Select value={status} onChange={onStatusChange} label="Status filter">
            <option value="">All statuses</option>
            {statusOptions.map((opt) => (
              <option key={opt.value} value={opt.value}>
                {opt.label}
              </option>
            ))}
          </Select>
        )}
        {sortOptions && onSortChange && (
          <Select value={sortValue ?? ''} onChange={onSortChange} label="Sort by">
            <option value="" disabled>Sort by...</option>
            {sortOptions.map((opt) => (
              <option key={opt.value} value={opt.value}>
                {opt.label}
              </option>
            ))}
          </Select>
        )}
        {onRefresh && (
          <Button variant="secondary" size="sm" icon="refresh" onClick={onRefresh}>
            Refresh
          </Button>
        )}
        {lastRefresh && <span class="refresh-status" data-visual-mask>{lastRefresh}</span>}
      </div>
      {activeFilters.length > 0 && onResetAll && (
        <div class="active-filters">
          {activeFilters.map((f) => (
            <FilterChip key={f.label} label={f.label} onRemove={f.onRemove} />
          ))}
          <Button variant="ghost" size="sm" onClick={onResetAll}>
            Reset all
          </Button>
        </div>
      )}
    </div>
  );
}
