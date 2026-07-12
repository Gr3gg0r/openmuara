import { useMemo, useState } from 'preact/hooks';
import { apiGet, apiPost, ApiError, statusClass } from '../api';
import { Button } from '../components/Button';
import { EmptyState } from '../components/EmptyState';
import { FilterToolbar } from '../components/FilterToolbar';
import { Icon } from '../components/Icon';
import { SkeletonRows } from '../components/Skeleton';
import { useAsync } from '../hooks/useAsync';
import { usePolling } from '../hooks/usePolling';
import { useUrlState } from '../hooks/useUrlState';
import type { ProvidersResponse, WebhookAttempt, WebhooksResponse } from '../types';

type SortKey = 'reference' | 'provider' | 'url' | 'status' | 'attempts';
type SortDir = 'asc' | 'desc';

interface WebhooksViewProps {
  onShowDetail?: (attempt: WebhookAttempt) => void;
}

const STATUS_OPTIONS = [
  { value: 'delivered', label: 'delivered' },
  { value: 'failed', label: 'failed' },
  { value: 'pending', label: 'pending' },
];

const SORT_OPTIONS = [
  { value: 'reference:asc', label: 'Reference (A–Z)' },
  { value: 'reference:desc', label: 'Reference (Z–A)' },
  { value: 'provider:asc', label: 'Provider (A–Z)' },
  { value: 'provider:desc', label: 'Provider (Z–A)' },
  { value: 'url:asc', label: 'URL (A–Z)' },
  { value: 'url:desc', label: 'URL (Z–A)' },
  { value: 'status:asc', label: 'Status (A–Z)' },
  { value: 'status:desc', label: 'Status (Z–A)' },
  { value: 'attempts:asc', label: 'Attempts (low first)' },
  { value: 'attempts:desc', label: 'Attempts (high first)' },
];

export function WebhooksView({ onShowDetail }: WebhooksViewProps) {
  const [query, setQuery] = useUrlState('q');
  const [provider, setProvider] = useUrlState('provider');
  const [status, setStatus] = useUrlState('status');
  const [sortKey, setSortKey] = useUrlState<SortKey>('sort', 'reference');
  const [sortDir, setSortDir] = useUrlState<SortDir>('dir', 'asc');
  const [lastRefresh, setLastRefresh] = useState<string | undefined>();

  const { data: providers } = useAsync<ProvidersResponse>(() => apiGet('/_admin/providers'), []);

  const fetchWebhooks = async () => {
    const params = new URLSearchParams();
    if (status) params.set('status', status);
    if (provider) params.set('provider', provider);
    try {
      const data = await apiGet<WebhooksResponse>('/_admin/webhooks?' + params.toString());
      setLastRefresh('Updated ' + new Date().toLocaleTimeString());
      return data;
    } catch (err) {
      if (err instanceof ApiError && err.status === 404) {
        return { results: [] };
      }
      throw err;
    }
  };

  const { data, loading, error, refetch } = useAsync<WebhooksResponse>(fetchWebhooks, [status, provider]);

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

  const filteredAttempts = useMemo(() => {
    let attempts = [...(data?.results ?? [])];
    if (query) {
      const q = query.toLowerCase();
      attempts = attempts.filter(
        (a) =>
          (a.ref ?? '').toLowerCase().includes(q) ||
          (a.provider_name ?? '').toLowerCase().includes(q) ||
          (a.url ?? '').toLowerCase().includes(q) ||
          (a.status ?? '').toLowerCase().includes(q),
      );
    }
    attempts.sort((a, b) => {
      let aVal: string | number | undefined;
      let bVal: string | number | undefined;
      switch (sortKey) {
        case 'reference':
          aVal = a.ref;
          bVal = b.ref;
          break;
        case 'provider':
          aVal = a.provider_name;
          bVal = b.provider_name;
          break;
        case 'url':
          aVal = a.url;
          bVal = b.url;
          break;
        case 'status':
          aVal = a.status;
          bVal = b.status;
          break;
        case 'attempts':
          aVal = a.attempts ?? 0;
          bVal = b.attempts ?? 0;
          break;
      }
      if (aVal === bVal) return 0;
      if (aVal == null) return sortDir === 'asc' ? -1 : 1;
      if (bVal == null) return sortDir === 'asc' ? 1 : -1;
      return sortDir === 'asc' ? (aVal < bVal ? -1 : 1) : aVal > bVal ? -1 : 1;
    });
    return attempts;
  }, [data, query, sortKey, sortDir]);

  const handleSortChange = (value: string) => {
    const [key, dir] = value.split(':') as [SortKey, SortDir];
    setSortKey(key);
    setSortDir(dir);
  };

  const handleRowClick = (attempt: WebhookAttempt, trigger?: HTMLButtonElement | null) => {
    if (onShowDetail) {
      onShowDetail(attempt);
      return;
    }
    trigger?.focus();
  };

  const replay = async (ref: string) => {
    await apiPost(`/_admin/webhooks/${encodeURIComponent(ref)}/replay`);
    setTimeout(refetch, 500);
  };

  const attempts = filteredAttempts;
  const hasFilter = status || provider || query;
  const sortIndicator = (key: SortKey) => (sortKey === key ? <Icon name={sortDir === 'asc' ? 'chevronUp' : 'chevronDown'} size={12} /> : null);

  return (
    <section aria-label="Webhook attempts">
      <h2>Webhook Attempts</h2>

      <FilterToolbar
        searchId="webhooks-search"
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

      {error && <div class="error-banner" role="alert">Failed to load webhooks: {error}</div>}

      <div class="table-wrap">
        <table class="zebra">
          <thead>
            <tr>
              <th><button class="sort-header" onClick={() => handleSortChange('reference:asc')}>Reference {sortIndicator('reference')}</button></th>
              <th><button class="sort-header" onClick={() => handleSortChange('provider:asc')}>Provider {sortIndicator('provider')}</button></th>
              <th class="hide-narrow hide-medium"><button class="sort-header" onClick={() => handleSortChange('url:asc')}>URL {sortIndicator('url')}</button></th>
              <th><button class="sort-header" onClick={() => handleSortChange('status:asc')}>Status {sortIndicator('status')}</button></th>
              <th><button class="sort-header" onClick={() => handleSortChange('attempts:asc')}>Attempts {sortIndicator('attempts')}</button></th>
              <th class="hide-narrow hide-medium">Last Error</th>
              <th>Action</th>
            </tr>
          </thead>
          <tbody>
            {loading && attempts.length === 0 ? (
              <SkeletonRows rows={5} columns={7} />
            ) : attempts.length === 0 ? (
              <tr>
                <td colspan={7} class="empty">
                  <EmptyState
                    title={hasFilter ? 'No webhook attempts match' : 'No webhook attempts yet'}
                    description={hasFilter ? 'Try clearing your filters.' : 'Send a test charge and wait for the webhook to be dispatched.'}
                    icon="empty"
                    action={hasFilter ? { label: 'Clear filters', onClick: () => { setQuery(''); setProvider(''); setStatus(''); } } : undefined}
                  />
                </td>
              </tr>
            ) : (
              attempts.map((a) => (
                <tr
                  key={a.ref}
                  class="row-click"
                >
                  <td>
                    <button
                      class="link-button"
                      onClick={(e) => {
                        const trigger = e.currentTarget as HTMLButtonElement;
                        void handleRowClick(a, trigger);
                      }}
                      aria-label={`View webhook details for ${a.ref}`}
                    >
                      <pre>{a.ref}</pre>
                    </button>
                  </td>
                  <td>{a.provider_name ?? '-'}</td>
                  <td class="hide-narrow hide-medium">{a.url}</td>
                  <td class={statusClass(a.status)}>{a.status ?? '-'}</td>
                  <td>{a.attempts ?? 0}</td>
                  <td class="muted hide-narrow hide-medium">{a.last_error ?? '-'}</td>
                  <td>
                    <Button variant="secondary" size="sm" icon="replay" onClick={() => replay(a.ref)}>
                      Replay
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
