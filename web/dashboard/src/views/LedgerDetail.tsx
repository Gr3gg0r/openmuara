import { useEffect, useRef } from 'preact/hooks';
import { apiGet, apiPost, formatDate } from '../api';
import { Badge } from '../components/Badge';
import type { BadgeProps } from '../components/Badge';
import { Button } from '../components/Button';
import { CodeBlock } from '../components/CodeBlock';
import { CopyButton } from '../components/CopyButton';
import { DetailField } from '../components/DetailField';
import { Icon } from '../components/Icon';
import { Skeleton } from '../components/Skeleton';
import { Timeline, type TimelineItem } from '../components/Timeline';
import { useAsync } from '../hooks/useAsync';
import type { LedgerEvent, Transaction, TransactionDetailResponse, WebhookAttempt, WebhookDetailResponse } from '../types';

interface LedgerDetailProps {
  event: LedgerEvent;
  onBack: () => void;
}

function statusBadgeVariant(status?: string): BadgeProps['variant'] {
  switch (status) {
    case 'paid':
    case 'delivered':
      return 'success';
    case 'failed':
    case 'unpaid':
      return 'danger';
    case 'pending':
    case 'new':
      return 'warning';
    case 'refunded':
      return 'info';
    default:
      return 'neutral';
  }
}

function formatAmount(amount?: number, currency?: string): string {
  if (amount == null) return '-';
  const symbol = currency ? ` ${currency}` : '';
  return `${amount.toFixed(2)}${symbol}`;
}

function transactionTimeline(tx?: Transaction): TimelineItem[] {
  if (!tx) return [];
  const status = tx.status ?? 'unknown';
  const created: TimelineItem = {
    status: status === 'new' ? 'current' : 'completed',
    label: 'Created',
    time: tx.createdAt,
    detail: `Reference ${tx.reference}`,
  };

  switch (status) {
    case 'new':
      return [created];
    case 'paid':
      return [
        created,
        { status: 'current', label: 'Paid', time: tx.updatedAt },
      ];
    case 'unpaid':
      return [
        created,
        { status: 'error', label: 'Failed / Unpaid', time: tx.updatedAt },
      ];
    case 'refunded':
      return [
        created,
        { status: 'completed', label: 'Paid', time: tx.updatedAt },
        { status: 'current', label: 'Refunded', time: tx.updatedAt },
      ];
    case 'pending':
      return [
        created,
        { status: 'current', label: 'Pending settlement', time: tx.updatedAt },
      ];
    default:
      return [created];
  }
}

function webhookTimeline(wh?: WebhookAttempt): TimelineItem[] {
  if (!wh) return [];
  const events = wh.attempt_events ?? [];
  const items: TimelineItem[] = events.map((ev, index) => {
    const httpStatus = ev.status ? Number(ev.status) : 0;
    const isLast = index === events.length - 1;
    let state: TimelineItem['status'] = 'completed';
    if (isLast) {
      state = wh.status === 'failed' ? 'error' : 'current';
    } else if (httpStatus < 200 || httpStatus >= 300) {
      state = 'error';
    }
    return {
      status: state,
      label: `Delivery attempt ${index + 1}`,
      time: ev.time,
      detail: ev.error ? `Error: ${ev.error}` : httpStatus ? `HTTP ${httpStatus}` : undefined,
    };
  });
  if (items.length === 0) {
    items.push({
      status: wh.status === 'failed' ? 'error' : 'current',
      label: wh.status === 'failed' ? 'Failed' : 'Queued',
      time: wh.createdAt,
    });
  }
  return items;
}

function DetailMeta({ label, value, mono, children }: { label: string; value?: string; mono?: boolean; children?: preact.ComponentChildren }) {
  return (
    <div class="detail-meta-item">
      <span class="detail-meta-label">{label}</span>
      {children ?? (mono ? <code class="detail-meta-value">{value}</code> : <span class="detail-meta-value">{value}</span>)}
    </div>
  );
}

function DetailSkeleton() {
  return (
    <div class="detail-skeleton">
      <div class="detail-section" style={{ marginTop: 0 }}>
        <Skeleton variant="title" className="skeleton-sm" />
        <div class="detail-grid">
          <Skeleton variant="line" />
          <Skeleton variant="line" />
          <Skeleton variant="line" />
          <Skeleton variant="line" />
          <Skeleton variant="line" />
          <Skeleton variant="line" />
        </div>
      </div>
      <div class="detail-section">
        <Skeleton variant="title" className="skeleton-sm" />
        <Skeleton variant="line" />
        <Skeleton variant="line" />
      </div>
    </div>
  );
}

function TransactionDetail({ tx }: { tx?: Transaction }) {
  if (!tx) return <div class="muted">No transaction details available.</div>;

  return (
    <>
      <div class="detail-section" style={{ marginTop: 0 }}>
        <h3>Overview</h3>
        <dl class="detail-grid detail-grid-3">
          <DetailField label="Provider" value={tx.provider ?? '-'} />
          <DetailField label="Amount" value={formatAmount(tx.amount, tx.currency)} />
          <DetailField label="Status" badge={{ text: tx.status ?? 'unknown', variant: statusBadgeVariant(tx.status) }} />
          <DetailField label="Customer" value={tx.customerRef ?? '-'} />
          <DetailField label="Created" value={formatDate(tx.createdAt)} />
          <DetailField label="Updated" value={formatDate(tx.updatedAt)} />
          {tx.trace_id && (
            <DetailField
              label="Trace ID"
              fullWidth
              value={
                <div class="detail-copy-row">
                  <code>{tx.trace_id}</code>
                  <CopyButton text={tx.trace_id} label="trace ID" size="sm" />
                </div>
              }
            />
          )}
        </dl>
      </div>

      <div class="detail-section">
        <h3>Status timeline</h3>
        <Timeline items={transactionTimeline(tx)} />
      </div>

      {tx.items && tx.items.length > 0 && (
        <div class="detail-section">
          <h3>Line items</h3>
          <CodeBlock value={tx.items} title={`${tx.items.length} item${tx.items.length === 1 ? '' : 's'}`} />
        </div>
      )}
    </>
  );
}

function WebhookDetailInline({ wh }: { wh?: WebhookAttempt }) {
  if (!wh) return <div class="muted">No webhook details available.</div>;

  return (
    <>
      <div class="detail-section" style={{ marginTop: 0 }}>
        <h3>Overview</h3>
        <dl class="detail-grid detail-grid-3">
          <DetailField label="Provider" value={wh.provider_name ?? wh.provider ?? '-'} />
          <DetailField label="Status" badge={{ text: wh.status ?? 'unknown', variant: statusBadgeVariant(wh.status) }} />
          {wh.attempts != null && <DetailField label="Attempts" value={String(wh.attempts)} />}
          {wh.signature_valid != null && (
            <DetailField
              label="Signature"
              badge={{ text: wh.signature_valid ? 'Valid' : 'Invalid', variant: wh.signature_valid ? 'success' : 'danger' }}
            />
          )}
          {wh.trace_id && (
            <DetailField
              label="Trace ID"
              fullWidth
              value={
                <div class="detail-copy-row">
                  <code>{wh.trace_id}</code>
                  <CopyButton text={wh.trace_id} label="trace ID" size="sm" />
                </div>
              }
            />
          )}
          {wh.url && (
            <DetailField
              label="URL"
              fullWidth
              value={<span class="break-all">{wh.url}</span>}
            />
          )}
        </dl>
      </div>

      <div class="detail-section">
        <h3>Delivery timeline</h3>
        <Timeline items={webhookTimeline(wh)} />
      </div>

      {wh.headers && Object.keys(wh.headers).length > 0 && (
        <div class="detail-section">
          <h3>Headers</h3>
          <CodeBlock value={wh.headers} title="Request headers" />
        </div>
      )}
      {wh.payload && (
        <div class="detail-section">
          <h3>Payload</h3>
          <CodeBlock value={wh.payload} title="Webhook payload" />
        </div>
      )}
    </>
  );
}

export function LedgerDetail({ event, onBack }: LedgerDetailProps) {
  const headerRef = useRef<HTMLButtonElement>(null);

  useEffect(() => {
    headerRef.current?.focus();
  }, []);

  const { data, loading, error } = useAsync<Transaction | WebhookAttempt | null>(async () => {
    if (event.type === 'transaction' && event.reference) {
      const d = await apiGet<TransactionDetailResponse>(`/_admin/transactions/${encodeURIComponent(event.reference)}`);
      return d.transaction ?? null;
    }
    if (event.type === 'webhook' && event.reference) {
      const d = await apiGet<WebhookDetailResponse>(`/_admin/webhooks/${encodeURIComponent(event.reference)}`);
      return d.webhook ?? null;
    }
    return null;
  }, [event.reference, event.type]);

  const replay = async () => {
    if (!event.reference) return;
    await apiPost(`/_admin/transactions/${encodeURIComponent(event.reference)}/replay-webhook`);
  };

  const title = event.type === 'transaction' ? 'Transaction detail' : 'Webhook detail';
  const providerName =
    (data as WebhookAttempt | null)?.provider_name ??
    (data as Transaction | WebhookAttempt | null)?.provider ??
    event.provider ??
    'unknown provider';

  return (
    <section aria-label="Ledger detail" class="detail-page">
      <div class="detail-header">
        <Button variant="secondary" size="sm" icon="chevronLeft" onClick={onBack} ref={headerRef}>
          Back to ledger
        </Button>
        <div class="detail-title-row">
          <h2 class="detail-title">{title}</h2>
          {event.status && <Badge variant={statusBadgeVariant(event.status)}>{event.status}</Badge>}
        </div>
        <div class="detail-meta">
          <DetailMeta label="Reference" mono value={event.reference}>
            <span class="detail-meta-value detail-meta-inline">
              <code>{event.reference}</code>
              <CopyButton text={event.reference} label="reference" size="sm" />
            </span>
          </DetailMeta>
          <DetailMeta label="Provider" value={providerName} />
          <DetailMeta label="Time" value={formatDate(event.time)} />
        </div>
      </div>

      {loading ? (
        <div class="card card-padded">
          <DetailSkeleton />
        </div>
      ) : error ? (
        <div class="error-banner" role="alert">
          <Icon name="alert" size={18} />
          <span>Failed to load detail: {error}</span>
        </div>
      ) : (
        <>
          <div class="card card-padded">
            {event.type === 'transaction' && <TransactionDetail tx={data as Transaction | undefined} />}
            {event.type === 'webhook' && <WebhookDetailInline wh={data as WebhookAttempt | undefined} />}
          </div>
          {event.type === 'transaction' && event.reference && (
            <div class="detail-actions">
              <Button icon="replay" onClick={replay}>
                Replay webhook
              </Button>
            </div>
          )}
        </>
      )}
    </section>
  );
}
