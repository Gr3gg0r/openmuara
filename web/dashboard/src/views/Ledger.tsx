import { useMemo, useState } from 'preact/hooks';
import { apiGet, formatDate, statusClass } from '../api';
import { Badge } from '../components/Badge';
import { Button } from '../components/Button';
import { EmptyState } from '../components/EmptyState';
import { FilterToolbar } from '../components/FilterToolbar';
import { Icon } from '../components/Icon';
import { SkeletonRows } from '../components/Skeleton';
import { useAsync } from '../hooks/useAsync';
import { usePolling } from '../hooks/usePolling';
import { useUrlState } from '../hooks/useUrlState';
import type { LedgerEvent, LedgerResponse, ProvidersResponse } from '../types';

type SortKey = 'time' | 'provider' | 'reference' | 'status';
type SortDir = 'asc' | 'desc';

interface LedgerViewProps {
  onShowDetail?: (event: LedgerEvent) => void;
}

const STATUS_OPTIONS = [
  { value: 'new', label: 'new' },
  { value: 'paid', label: 'paid' },
  { value: 'unpaid', label: 'unpaid' },
  { value: 'refunded', label: 'refunded' },
  { value: 'pending', label: 'pending' },
  { value: 'delivered', label: 'delivered' },
  { value: 'failed', label: 'failed' },
];

const SORT_OPTIONS = [
  { value: 'time:desc', label: 'Time (newest first)' },
  { value: 'time:asc', label: 'Time (oldest first)' },
  { value: 'provider:asc', label: 'Provider (A–Z)' },
  { value: 'provider:desc', label: 'Provider (Z–A)' },
  { value: 'reference:asc', label: 'Reference (A–Z)' },
  { value: 'reference:desc', label: 'Reference (Z–A)' },
  { value: 'status:asc', label: 'Status (A–Z)' },
  { value: 'status:desc', label: 'Status (Z–A)' },
];

export function LedgerView({ onShowDetail }: LedgerViewProps) {
  const [query, setQuery] = useUrlState('q');
  const [provider, setProvider] = useUrlState('provider');
  const [status, setStatus] = useUrlState('status');
  const [sortKey, setSortKey] = useUrlState<SortKey>('sort', 'time');
  const [sortDir, setSortDir] = useUrlState<SortDir>('dir', 'desc');
  const [lastRefresh, setLastRefresh] = useState<string | undefined>();

  const { data: providers } = useAsync<ProvidersResponse>(() => apiGet('/_admin/providers'), []);

  const fetchLedger = useMemo(
    () => async () => {
      const params = new URLSearchParams();
      if (query) params.set('q', query);
      if (provider) params.set('provider', provider);
      if (status) params.set('status', status);
      const data = await apiGet<LedgerResponse>('/_admin/ledger?' + params.toString());
      setLastRefresh('Updated ' + new Date().toLocaleTimeString());
      return data;
    },
    [query, provider, status],
  );

  const { data: ledger, loading, error, refetch } = useAsync<LedgerResponse>(fetchLedger, [fetchLedger]);

  usePolling(refetch, 2000);

  const providerOptions = useMemo(() => {
    const map = new Map<string, string>();
    (providers?.available ?? []).forEach((name) => {
      const info = providers?.providers?.[name];
      const label = info?.display_name ? `${info.display_name} (${name})` : name;
      map.set(name, label);
    });
    return map;
  }, [providers]);

  const activeFilters = useMemo(() => {
    const filters: { label: string; onRemove: () => void }[] = [];
    if (provider) {
      const label = providerOptions.get(provider) ?? provider;
      filters.push({ label: `Provider: ${label}`, onRemove: () => setProvider('') });
    }
    if (status) filters.push({ label: `Status: ${status}`, onRemove: () => setStatus('') });
    if (query) filters.push({ label: `Search: ${query}`, onRemove: () => setQuery('') });
    return filters;
  }, [provider, status, query, providerOptions, setProvider, setStatus, setQuery]);

  const sortedEvents = useMemo(() => {
    const events = [...(ledger?.results ?? [])];
    events.sort((a, b) => {
      let aVal: string | undefined;
      let bVal: string | undefined;
      switch (sortKey) {
        case 'time':
          aVal = a.time;
          bVal = b.time;
          break;
        case 'provider':
          aVal = a.provider;
          bVal = b.provider;
          break;
        case 'reference':
          aVal = a.reference;
          bVal = b.reference;
          break;
        case 'status':
          aVal = a.status;
          bVal = b.status;
          break;
      }
      if (aVal === bVal) return 0;
      if (aVal == null) return sortDir === 'asc' ? -1 : 1;
      if (bVal == null) return sortDir === 'asc' ? 1 : -1;
      return sortDir === 'asc' ? (aVal < bVal ? -1 : 1) : aVal > bVal ? -1 : 1;
    });
    return events;
  }, [ledger, sortKey, sortDir]);

  const handleSortChange = (value: string) => {
    const [key, dir] = value.split(':') as [SortKey, SortDir];
    setSortKey(key);
    setSortDir(dir);
  };

  const handleRowClick = (ev: LedgerEvent, trigger?: HTMLButtonElement | null) => {
    if (onShowDetail) {
      onShowDetail(ev);
      return;
    }
    // Fallback: keep legacy inline behavior if no detail handler provided.
    trigger?.focus();
  };

  const events = sortedEvents;
  const hasFilter = query || provider || status;
  const sortIndicator = (key: SortKey) => (sortKey === key ? <Icon name={sortDir === 'asc' ? 'chevronUp' : 'chevronDown'} size={12} /> : null);

  return (
    <section aria-label="Ledger">
      <h2>Ledger</h2>
      <FilterToolbar
        searchId="ledger-search"
        query={query}
        onQueryChange={setQuery}
        providerOptions={providerOptions}
        provider={provider}
        onProviderChange={setProvider}
        statusOptions={STATUS_OPTIONS}
        status={status}
        onStatusChange={setStatus}
        sortValue={`${sortKey}:${sortDir}`}
        sortOptions={SORT_OPTIONS}
        onSortChange={handleSortChange}
        onRefresh={refetch}
        lastRefresh={lastRefresh}
        activeFilters={activeFilters}
        onResetAll={() => { setQuery(''); setProvider(''); setStatus(''); }}
      />

      {error && <div class="error-banner" role="alert">Failed to load ledger: {error}</div>}

      <div class="table-wrap">
        <table class="zebra">
          <thead>
            <tr>
              <th><button class="sort-header" onClick={() => handleSortChange('time:desc')}>Time {sortIndicator('time')}</button></th>
              <th>Type</th>
              <th><button class="sort-header" onClick={() => handleSortChange('provider:asc')}>Provider {sortIndicator('provider')}</button></th>
              <th><button class="sort-header" onClick={() => handleSortChange('reference:asc')}>Reference {sortIndicator('reference')}</button></th>
              <th><button class="sort-header" onClick={() => handleSortChange('status:asc')}>Status {sortIndicator('status')}</button></th>
              <th class="hide-narrow hide-medium">Summary</th>
              <th class="hide-narrow hide-medium">Action</th>
            </tr>
          </thead>
          <tbody>
            {loading && events.length === 0 ? (
              <SkeletonRows rows={5} columns={7} />
            ) : events.length === 0 ? (
              <tr>
                <td colspan={7} class="empty">
                  <EmptyState
                    title={hasFilter ? 'No events match' : 'No events yet'}
                    description={hasFilter ? 'Try clearing your filters.' : 'Simulate your first charge to see activity.'}
                    icon="empty"
                    action={hasFilter ? { label: 'Clear filters', onClick: () => { setQuery(''); setProvider(''); setStatus(''); } } : { label: 'Simulate charge', onClick: () => {} }}
                  />
                </td>
              </tr>
            ) : (
              events.map((ev) => (
                <tr
                  key={ev.reference + ev.time}
                  class="row-click"
                >
                  <td>{formatDate(ev.time)}</td>
                  <td>
                    <Badge variant={ev.type === 'transaction' ? 'info' : 'neutral'}>{ev.type}</Badge>
                  </td>
                  <td>{ev.provider ?? '-'}</td>
                  <td>
                    <button
                      class="link-button"
                      onClick={(e) => {
                        const trigger = e.currentTarget as HTMLButtonElement;
                        void handleRowClick(ev, trigger);
                      }}
                      aria-label={`View ${ev.type} details for ${ev.reference ?? 'unknown'}`}
                    >
                      <pre>{ev.reference}</pre>
                    </button>
                  </td>
                  <td class={statusClass(ev.status)}>{ev.status ?? '-'}</td>
                  <td class="muted hide-narrow hide-medium">{ev.summary ?? '-'}</td>
                  <td class="hide-narrow hide-medium">
                    <Button variant="secondary" size="sm" onClick={() => handleRowClick(ev)}>
                      View
                    </Button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </section>
  );
}
