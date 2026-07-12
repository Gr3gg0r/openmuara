import { useEffect, useMemo, useState } from 'preact/hooks';
import { announce } from '../components/Announce';
import { apiGet, apiPatch, apiPost } from '../api';
import { Badge } from '../components/Badge';
import { Button } from '../components/Button';
import { Card } from '../components/Card';
import { Icon } from '../components/Icon';
import { useAsync } from '../hooks/useAsync';
import type { ConfigResponse, ProviderDetailResponse } from '../types';

interface ProviderDetailProps {
  name: string;
  onBack: () => void;
}

const DEFAULT_EVENTS: Record<string, string[]> = {
  stripe: ['checkout.session.completed', 'payment_intent.succeeded', 'payment_intent.payment_failed'],
  fawry: ['charge_paid'],
};

export function ProviderDetail({ name, onBack }: ProviderDetailProps) {
  const { data, loading, error, refetch } = useAsync<ProviderDetailResponse>(
    () => apiGet(`/_admin/providers/${encodeURIComponent(name)}`),
    [name],
  );
  const { data: config, refetch: refetchConfig } = useAsync<ConfigResponse>(() => apiGet('/_admin/config'), []);

  const [selectedVersion, setSelectedVersion] = useState<string>(data?.version ?? '');
  const [enabled, setEnabled] = useState<boolean>(data?.enabled ?? false);
  const [webhookUrl, setWebhookUrl] = useState<string>(data?.webhook_target_url ?? '');
  const [testSecret, setTestSecret] = useState<string>('');
  const [showSecret, setShowSecret] = useState<boolean>(false);
  const [saving, setSaving] = useState<boolean>(false);
  const [testing, setTesting] = useState<boolean>(false);
  const [testResult, setTestResult] = useState<{ ok: boolean; message: string } | null>(null);
  const [restartNotice, setRestartNotice] = useState<string | null>(null);

  useEffect(() => {
    if (data) {
      setEnabled(data.enabled ?? false);
      setSelectedVersion(data.version ?? '');
      setWebhookUrl(data.webhook_target_url ?? '');
    }
  }, [data]);

  const versions = useMemo(() => {
    const list = data?.versions ?? [];
    if (list.length === 0 && data?.version) return [data.version];
    return list;
  }, [data]);

  const versionDetails = useMemo(() => {
    if (!data?.version_details) return undefined;
    return selectedVersion ? data.version_details[selectedVersion] : undefined;
  }, [data, selectedVersion]);

  const displayName = data?.display_name ?? name;
  const isActive = data?.active ?? false;
  const providerEvents = DEFAULT_EVENTS[name] ?? [];

  const saveEnabled = async (next: boolean) => {
    setSaving(true);
    try {
      await apiPatch('/_admin/config/providers', {
        providers: { [name]: { enabled: next } },
      });
      setEnabled(next);
      setRestartNotice('Provider enablement changed. Restart the server for changes to take effect.');
      announce(`${displayName} ${next ? 'enabled' : 'disabled'}; restart required`);
      await refetch();
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      announce(`Failed to update ${displayName}: ${message}`);
    } finally {
      setSaving(false);
    }
  };

  const saveWebhook = async (e: Event) => {
    e.preventDefault();
    setSaving(true);
    try {
      const targets = { ...(config?.webhook?.targets ?? {}), [name]: webhookUrl };
      await apiPatch('/_admin/config/webhooks', { targets });
      setRestartNotice('Webhook target changed. Restart the server for changes to take effect.');
      announce(`${displayName} webhook target saved; restart required`);
      await refetchConfig();
      await refetch();
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      announce(`Failed to save webhook target: ${message}`);
    } finally {
      setSaving(false);
    }
  };

  const testWebhook = async () => {
    const targetUrl = webhookUrl || config?.webhook?.url;
    if (!targetUrl) {
      setTestResult({ ok: false, message: 'No webhook URL configured' });
      return;
    }
    setTesting(true);
    setTestResult(null);
    try {
      const res = await apiPost('/_admin/config/webhooks/test', { url: targetUrl, provider: name, secret: testSecret });
      const body = (await res.json()) as { success: boolean; latency_ms?: number; error?: string };
      setTestResult({
        ok: body.success,
        message: body.success ? `Delivered in ${body.latency_ms ?? '?'}ms` : body.error || 'Failed',
      });
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      setTestResult({ ok: false, message });
    } finally {
      setTesting(false);
    }
  };

  const copyToClipboard = async (text: string) => {
    try {
      await navigator.clipboard.writeText(text);
      announce('Copied to clipboard');
    } catch {
      announce('Copy failed');
    }
  };

  if (loading && !data) {
    return (
      <section aria-label={`${displayName} settings`}>
        <Button variant="ghost" size="sm" icon="chevronLeft" onClick={onBack}>
          Back to Settings
        </Button>
        <h2 class="mt-4">Loading {displayName}...</h2>
        <div class="skeleton skeleton-title" />
        <div class="skeleton skeleton-line" />
        <div class="skeleton skeleton-line" />
      </section>
    );
  }

  if (error || !data) {
    return (
      <section aria-label={`${displayName} settings`}>
        <Button variant="ghost" size="sm" icon="chevronLeft" onClick={onBack}>
          Back to Settings
        </Button>
        <div class="error-banner mt-4" role="alert">
          Failed to load provider details: {error}
        </div>
      </section>
    );
  }

  return (
    <section aria-label={`${displayName} settings`}>
      <Button variant="ghost" size="sm" icon="chevronLeft" onClick={onBack}>
        Back to Settings
      </Button>

      <Card className="provider-detail-header mt-4">
        <div class="flex justify-between items-start flex-wrap gap-3">
          <div>
            <h2 style={{ margin: '0 0 4px' }}>{displayName}</h2>
            <p class="muted" style={{ margin: 0 }}>{data.description}</p>
            <div class="flex gap-2 mt-2 flex-wrap">
              {isActive && <Badge variant="active">Active</Badge>}
              {enabled ? <Badge variant="success">Enabled</Badge> : <Badge variant="neutral">Disabled</Badge>}
              {data.category && <Badge variant="info">{data.category}</Badge>}
            </div>
          </div>
          <Button
            variant={enabled ? 'secondary' : 'primary'}
            loading={saving}
            onClick={() => saveEnabled(!enabled)}
            aria-pressed={enabled}
          >
            {enabled ? 'Disable' : 'Enable'}
          </Button>
        </div>
      </Card>

      {restartNotice && (
        <div class="alert info-banner mt-4 active" role="status">
          <Icon name="info" size={18} />
          <span>{restartNotice}</span>
          <Button variant="ghost" size="sm" onClick={() => setRestartNotice(null)} aria-label="Dismiss restart notice">
            Dismiss
          </Button>
        </div>
      )}

      {versions.length > 1 && (
        <div class="tabs mt-4" role="tablist" aria-label="API versions">
          {versions.map((v) => (
            <button
              key={v}
              class={`tab ${selectedVersion === v ? 'active' : ''}`}
              role="tab"
              aria-selected={selectedVersion === v}
              onClick={() => setSelectedVersion(v)}
            >
              {v}
            </button>
          ))}
        </div>
      )}

      <Card className="mt-4">
        <h3 class="mt-0">Base URL</h3>
        <div class="flex items-center gap-2 flex-wrap">
          <code>{versionDetails?.base_url ?? data.base_url ?? '-'}</code>
          {(versionDetails?.base_url ?? data.base_url) && (
            <Button
              variant="secondary"
              size="sm"
              icon="copy"
              onClick={() => copyToClipboard(versionDetails?.base_url ?? data.base_url ?? '')}
              aria-label="Copy base URL"
            >
              Copy
            </Button>
          )}
        </div>

        <h3>Sample endpoint</h3>
        <p class="muted">
          <code>
            {data.sample_method ?? 'POST'} {versionDetails?.sample_route ?? data.sample_route ?? '-'}
          </code>
        </p>

        {data.docs_path && (
          <p class="muted">
            <a href={data.docs_path} target="_blank" rel="noopener noreferrer">
              Provider documentation <Icon name="external" size={12} />
            </a>
          </p>
        )}
      </Card>

      <Card className="mt-4">
        <h3 class="mt-0">Webhook target URL</h3>
        <p class="muted">
          Configure a webhook URL specific to {displayName}. Leave empty to fall back to the global webhook URL.
        </p>
        <form onSubmit={saveWebhook}>
          <div class="mb-3">
            <label class="label" for={`provider-webhook-url-${name}`}>
              Target URL
            </label>
            <input
              id={`provider-webhook-url-${name}`}
              type="url"
              class="input input-full"
              value={webhookUrl}
              onInput={(e) => setWebhookUrl((e.currentTarget as HTMLInputElement).value)}
              placeholder="https://example.com/webhook"
            />
          </div>

          {providerEvents.length > 0 && (
            <div class="mb-3">
              <div class="label">Events dispatched</div>
              <ul class="muted" style={{ margin: 0, paddingLeft: '1.2rem' }}>
                {providerEvents.map((event) => (
                  <li key={event}>{event}</li>
                ))}
              </ul>
            </div>
          )}

          <div class="mb-3">
            <label class="label" for={`provider-webhook-secret-${name}`}>
              Test signing secret
            </label>
            <div class="flex gap-2 flex-wrap">
              <input
                id={`provider-webhook-secret-${name}`}
                type={showSecret ? 'text' : 'password'}
                class="input input-full"
                value={testSecret}
                onInput={(e) => setTestSecret((e.currentTarget as HTMLInputElement).value)}
                placeholder="Optional secret for test signature"
              />
              <Button
                variant="secondary"
                size="sm"
                icon={showSecret ? 'eyeOff' : 'eye'}
                onClick={() => setShowSecret((v) => !v)}
                aria-label={`${showSecret ? 'Hide' : 'Show'} test secret`}
              >
                {showSecret ? 'Hide' : 'Show'}
              </Button>
            </div>
          </div>

          <div class="flex gap-2 flex-wrap items-center">
            <Button type="submit" loading={saving}>
              Save webhook target
            </Button>
            <Button variant="secondary" loading={testing} onClick={testWebhook}>
              Test delivery
            </Button>
            {testResult && !testing && (
              <Badge variant={testResult.ok ? 'success' : 'danger'}>{testResult.message}</Badge>
            )}
          </div>
        </form>
      </Card>

      <Card className="mt-4">
        <h3 class="mt-0">Environment variables</h3>
        <p class="muted">Reference only — values are not exposed in the UI.</p>
        {data.env_vars && data.env_vars.length > 0 ? (
          <ul class="env-var-list">
            {data.env_vars.map((envVar) => (
              <li key={envVar} class="env-var-item">
                <code>{envVar}</code>
                <Button
                  variant="ghost"
                  size="sm"
                  icon="copy"
                  onClick={() => copyToClipboard(envVar)}
                  aria-label={`Copy ${envVar}`}
                >
                  Copy
                </Button>
              </li>
            ))}
          </ul>
        ) : (
          <p class="muted">No environment variables defined for this provider.</p>
        )}
      </Card>
    </section>
  );
}
