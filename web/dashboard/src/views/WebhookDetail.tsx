import { useEffect, useMemo, useRef } from 'preact/hooks';
import { apiGet, apiPost } from '../api';
import { Badge } from '../components/Badge';
import type { BadgeProps } from '../components/Badge';
import { Button } from '../components/Button';
import { CodeBlock } from '../components/CodeBlock';
import { CopyButton } from '../components/CopyButton';
import { DetailField } from '../components/DetailField';
import { Skeleton } from '../components/Skeleton';
import { Timeline, type TimelineItem } from '../components/Timeline';
import { useAsync } from '../hooks/useAsync';
import type { WebhookAttempt, WebhookDetailResponse } from '../types';

interface WebhookDetailProps {
  webhook: WebhookAttempt;
  onBack: () => void;
}

function statusBadgeVariant(status?: string): BadgeProps['variant'] {
  switch (status) {
    case 'delivered':
      return 'success';
    case 'failed':
      return 'danger';
    case 'pending':
      return 'warning';
    default:
      return 'neutral';
  }
}

function buildTimeline(wh: WebhookAttempt): TimelineItem[] {
  const events = wh.attempt_events ?? [];
  if (events.length > 0) {
    return events.map((ev, index) => {
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
        label: `Attempt ${index + 1}`,
        time: ev.time,
        detail: ev.error ? `Error: ${ev.error}` : httpStatus ? `HTTP ${httpStatus}` : undefined,
      };
    });
  }
  return [
    {
      status: wh.status === 'failed' ? 'error' : 'current',
      label: wh.status === 'failed' ? 'Failed' : 'Queued',
      time: wh.createdAt,
    },
  ];
}

export function WebhookDetail({ webhook, onBack }: WebhookDetailProps) {
  const headerRef = useRef<HTMLButtonElement>(null);

  useEffect(() => {
    headerRef.current?.focus();
  }, []);

  const { data: detail, loading, error } = useAsync<WebhookAttempt>(async () => {
    const d = await apiGet<WebhookDetailResponse>(`/_admin/webhooks/${encodeURIComponent(webhook.ref)}`);
    return d.webhook ?? webhook;
  }, [webhook.ref]);

  const wh = detail ?? webhook;
  const timeline = useMemo(() => buildTimeline(wh), [wh]);

  const replay = async () => {
    await apiPost(`/_admin/webhooks/${encodeURIComponent(webhook.ref)}/replay`);
  };

  return (
    <section aria-label="Webhook detail" class="detail-page">
      <div class="detail-header">
        <Button variant="secondary" size="sm" icon="chevronLeft" onClick={onBack} ref={headerRef}>
          Back to webhooks
        </Button>
        <div class="detail-title-row">
          <h2 class="detail-title">Webhook detail</h2>
          {wh.status && <Badge variant={statusBadgeVariant(wh.status)}>{wh.status}</Badge>}
        </div>
        <p class="detail-subtitle">
          <code>{wh.ref}</code>
          <CopyButton text={wh.ref} label="reference" size="sm" />
          <span class="detail-dot" aria-hidden="true" />
          {wh.provider_name ?? wh.provider ?? 'unknown provider'}
        </p>
      </div>

      {loading ? (
        <div class="card card-padded">
          <Skeleton variant="title" />
          <Skeleton variant="line" />
          <Skeleton variant="line" />
          <Skeleton variant="line" />
        </div>
      ) : error ? (
        <div class="error-banner" role="alert">Failed to load detail: {error}</div>
      ) : (
        <>
          <div class="card card-padded">
            <div class="detail-section" style={{ marginTop: 0 }}>
              <h3>Overview</h3>
              <dl class="detail-grid">
                <DetailField label="Provider" value={wh.provider_name ?? wh.provider ?? '-'} />
                <DetailField label="URL" value={<span class="break-all">{wh.url ?? '-'}</span>} />
                <DetailField label="Status" badge={{ text: wh.status ?? 'unknown', variant: statusBadgeVariant(wh.status) }} />
                {wh.signature_valid != null && (
                  <DetailField
                    label="Signature"
                    badge={{
                      text: wh.signature_valid ? 'Valid' : 'Invalid',
                      variant: wh.signature_valid ? 'success' : 'danger',
                    }}
                  />
                )}
                {wh.trace_id && (
                  <DetailField
                    label="Trace ID"
                    value={
                      <div class="detail-copy-row">
                        <code>{wh.trace_id}</code>
                        <CopyButton text={wh.trace_id} label="trace ID" size="sm" />
                      </div>
                    }
                  />
                )}
                {wh.attempts != null && <DetailField label="Attempts" value={String(wh.attempts)} />}
              </dl>
            </div>

            <div class="detail-section">
              <h3>Delivery timeline</h3>
              <Timeline items={timeline} />
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
          </div>
          <div class="detail-actions">
            <Button icon="replay" onClick={replay}>
              Replay webhook
            </Button>
          </div>
        </>
      )}
    </section>
  );
}
