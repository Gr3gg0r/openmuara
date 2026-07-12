import { useEffect, useMemo, useRef, useState } from 'preact/hooks';
import { apiGet, apiPost, formatDate, statusClass } from '../api';
import { Button } from '../components/Button';
import { EmptyState } from '../components/EmptyState';
import { FilterChip } from '../components/FilterChip';
import { Icon } from '../components/Icon';
import { Input } from '../components/Input';
import { Select } from '../components/Select';
import { SkeletonRows } from '../components/Skeleton';
import { useAsync } from '../hooks/useAsync';
import { usePolling } from '../hooks/usePolling';
import { useUrlState } from '../hooks/useUrlState';
import type { ProvidersResponse, Transaction, TransactionDetailResponse, TransactionsResponse } from '../types';

type SortKey = 'reference' | 'provider' | 'amount' | 'status' | 'created';
type SortDir = 'asc' | 'desc';

export function TransactionsView() {
  const [query, setQuery] = useUrlState('q');
  const [provider, setProvider] = useUrlState('provider');
  const [status, setStatus] = useUrlState('status');
  const [sortKey, setSortKey] = useUrlState<SortKey>('sort', 'created');
  const [sortDir, setSortDir] = useUrlState<SortDir>('dir', 'desc');
  const [selectedRef, setSelectedRef] = useState<string | null>(null);
  const [detail, setDetail] = useState<{ tx: Transaction; html: string } | null>(null);

  const { data: providers } = useAsync<ProvidersResponse>(() => apiGet('/_admin/providers'), []);
  const detailRef = useRef<HTMLDivElement>(null);
  const lastFocusedRow = useRef<HTMLButtonElement | null>(null);

  const fetchTransactions = useMemo(
    () => async () => {
      const params = new URLSearchParams();
      if (query) params.set('q', query);
      if (provider) params.set('provider', provider);
      if (status) params.set('status', status);
      return apiGet<TransactionsResponse>('/_admin/transactions?' + params.toString());
    },
    [query, provider, status],
  );

  const { data, loading, error, refetch } = useAsync<TransactionsResponse>(fetchTransactions, [fetchTransactions]);

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
    if (provider) filters.push({ label: `Provider: ${providerOptions.get(provider) ?? provider}`, onRemove: () => setProvider('') });
    if (status) filters.push({ label: `Status: ${status}`, onRemove: () => setStatus('') });
    if (query) filters.push({ label: `Search: ${query}`, onRemove: () => setQuery('') });
    return filters;
  }, [provider, status, query, providerOptions, setProvider, setStatus, setQuery]);

  const sortedTxs = useMemo(() => {
    const txs = [...(data?.results ?? [])];
    txs.sort((a, b) => {
      let aVal: string | number | undefined;
      let bVal: string | number | undefined;
      switch (sortKey) {
        case 'reference':
          aVal = a.reference;
          bVal = b.reference;
          break;
        case 'provider':
          aVal = a.provider;
          bVal = b.provider;
          break;
        case 'amount':
          aVal = a.amount ?? 0;
          bVal = b.amount ?? 0;
          break;
        case 'status':
          aVal = a.status;
          bVal = b.status;
          break;
        case 'created':
          aVal = a.createdAt;
          bVal = b.createdAt;
          break;
      }
      if (aVal === bVal) return 0;
      if (aVal == null) return sortDir === 'asc' ? -1 : 1;
      if (bVal == null) return sortDir === 'asc' ? 1 : -1;
      return sortDir === 'asc' ? (aVal < bVal ? -1 : 1) : aVal > bVal ? -1 : 1;
    });
    return txs;
  }, [data, sortKey, sortDir]);

  const setSort = (key: SortKey) => {
    if (sortKey === key) {
      setSortDir(sortDir === 'asc' ? 'desc' : 'asc');
    } else {
      setSortKey(key);
      setSortDir('asc');
    }
  };

  const showDetail = async (ref: string, trigger?: HTMLButtonElement | null) => {
    lastFocusedRow.current = trigger ?? null;
    setSelectedRef(ref);
    const d = await apiGet<TransactionDetailResponse>(`/_admin/transactions/${encodeURIComponent(ref)}`);
    const tx = d.transaction;
    if (!tx) return;
    setDetail({ tx, html: '' });
  };

  useEffect(() => {
    if (detail && detailRef.current) {
      requestAnimationFrame(() => {
        const firstFocusable = detailRef.current?.querySelector('button, a, input, select, textarea') as HTMLElement | null;
        firstFocusable?.focus();
      });
    }
  }, [detail]);

  const closeDetail = () => {
    setDetail(null);
    setSelectedRef(null);
    requestAnimationFrame(() => lastFocusedRow.current?.focus());
  };

  const replayWebhook = async (ref: string) => {
    await apiPost(`/_admin/transactions/${encodeURIComponent(ref)}/replay-webhook`);
  };

  const txs = sortedTxs;
  const hasFilter = query || provider || status;
  const sortIndicator = (key: SortKey) => (sortKey === key ? <Icon name={sortDir === 'asc' ? 'chevronUp' : 'chevronDown'} size={12} /> : null);

  return (
    <section aria-label="Transactions">
      <h2>Transactions</h2>
      <div class="card card-padded mb-4">
        <div class="toolbar">
          <Input
            id="transactions-search"
            type="search"
            value={query}
            placeholder="Search reference, provider, status..."
            icon="search"
            clearable
            onInput={setQuery}
            onClear={() => setQuery('')}
          />
          <Select value={provider} onChange={setProvider} label="Provider filter">
            <option value="">All providers</option>
            {Array.from(providerOptions.entries()).map(([value, label]) => (
              <option key={value} value={value}>
                {label}
              </option>
            ))}
          </Select>
          <Select value={status} onChange={setStatus} label="Status filter">
            <option value="">All statuses</option>
            <option value="new">new</option>
            <option value="paid">paid</option>
            <option value="unpaid">unpaid</option>
            <option value="refunded">refunded</option>
          </Select>
          <Button variant="secondary" size="sm" icon="refresh" onClick={refetch}>
            Refresh
          </Button>
        </div>
        {activeFilters.length > 0 && (
          <div class="active-filters">
            {activeFilters.map((f) => (
              <FilterChip key={f.label} label={f.label} onRemove={f.onRemove} />
            ))}
            <Button variant="ghost" size="sm" onClick={() => { setQuery(''); setProvider(''); setStatus(''); }}>
              Reset all
            </Button>
          </div>
        )}
      </div>

      {error && <div class="error-banner" role="alert">Failed to load transactions: {error}</div>}

      <div class="table-wrap">
        <table class="zebra">
          <thead>
            <tr>
              <th><button class="sort-header" onClick={() => setSort('reference')}>Reference {sortIndicator('reference')}</button></th>
              <th><button class="sort-header" onClick={() => setSort('provider')}>Provider {sortIndicator('provider')}</button></th>
              <th><button class="sort-header" onClick={() => setSort('amount')}>Amount {sortIndicator('amount')}</button></th>
              <th><button class="sort-header" onClick={() => setSort('status')}>Status {sortIndicator('status')}</button></th>
              <th class="hide-narrow"><button class="sort-header" onClick={() => setSort('created')}>Created {sortIndicator('created')}</button></th>
              <th>Action</th>
            </tr>
          </thead>
          <tbody>
            {loading && txs.length === 0 ? (
              <SkeletonRows rows={5} columns={6} />
            ) : txs.length === 0 ? (
              <tr>
                <td colspan={6} class="empty">
                  <EmptyState
                    title={hasFilter ? 'No transactions match' : 'No transactions yet'}
                    description={hasFilter ? 'Try clearing your filters.' : 'Send a test charge to see transactions.'}
                    icon="empty"
                    action={hasFilter ? { label: 'Clear filters', onClick: () => { setQuery(''); setProvider(''); setStatus(''); } } : undefined}
                  />
                </td>
              </tr>
            ) : (
              txs.map((tx) => (
                <tr
                  key={tx.reference}
                  class={`row-click ${selectedRef === tx.reference ? 'row-selected' : ''}`}
                >
                  <td>
                    <button
                      class="link-button"
                      onClick={(e) => {
                        const trigger = e.currentTarget as HTMLButtonElement;
                        void showDetail(tx.reference, trigger);
                      }}
                      aria-label={`View transaction details for ${tx.reference}`}
                    >
                      <pre>{tx.reference}</pre>
                    </button>
                  </td>
                  <td>{tx.provider ?? '-'}</td>
                  <td>{tx.amount != null ? tx.amount.toFixed(2) + ' ' + (tx.currency ?? '') : '-'}</td>
                  <td class={statusClass(tx.status)}>{tx.status ?? '-'}</td>
                  <td class="hide-narrow">{formatDate(tx.createdAt)}</td>
                  <td>
                    <Button variant="secondary" size="sm" icon="replay" onClick={() => replayWebhook(tx.reference)}>
                      Replay
                    </Button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {detail && (
        <div class="card card-padded mt-4" ref={detailRef}>
          <div class="flex justify-between items-center mb-3">
            <strong>Transaction: {detail.tx.reference}</strong>
            <Button variant="secondary" size="sm" icon="close" onClick={closeDetail} aria-label="Close detail">
              Close
            </Button>
          </div>
          <div>
            <div class="muted">Provider</div><div>{detail.tx.provider ?? '-'}</div>
            <div class="muted">Amount</div><div>{detail.tx.amount != null ? detail.tx.amount.toFixed(2) + ' ' + (detail.tx.currency ?? '') : '-'}</div>
            <div class="muted">Status</div><div class={statusClass(detail.tx.status)}>{detail.tx.status ?? '-'}</div>
            {detail.tx.trace_id && <><div class="muted">Trace ID</div><code>{detail.tx.trace_id}</code></>}
            <div class="muted">Customer</div><div>{detail.tx.customerRef ?? '-'}</div>
            <div class="muted">Created</div><div>{formatDate(detail.tx.createdAt)}</div>
            <div class="muted">Updated</div><div>{formatDate(detail.tx.updatedAt)}</div>
            {detail.tx.items?.length ? <><div class="muted">Items</div><pre>{JSON.stringify(detail.tx.items, null, 2)}</pre></> : null}
          </div>
          <div class="mt-3">
            <Button icon="replay" onClick={() => replayWebhook(detail.tx.reference)}>
              Replay webhook
            </Button>
          </div>
        </div>
      )}
    </section>
  );
}
