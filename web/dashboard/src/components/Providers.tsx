import { useState } from 'preact/hooks';
import { announce } from './Announce';
import { apiGet, apiPatch, curlExample, escapeHtml } from '../api';
import { useAsync } from '../hooks/useAsync';
import { Badge } from './Badge';
import { Button } from './Button';
import { Card } from './Card';
import { Icon } from './Icon';
import type { ProvidersResponse } from '../types';

export function Providers() {
  const { data, refetch } = useAsync<ProvidersResponse>(() => apiGet('/_admin/providers'), []);
  const [copied, setCopied] = useState<string | null>(null);
  const [saving, setSaving] = useState<string | null>(null);

  const toggleProvider = async (name: string, enabled: boolean) => {
    setSaving(name);
    try {
      await apiPatch('/_admin/config/providers', { providers: { [name]: { enabled } } });
      announce(`${name} ${enabled ? 'enabled' : 'disabled'}; restart server to activate`);
      await refetch();
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      announce(`Failed to update ${name}: ${message}`);
    } finally {
      setSaving(null);
    }
  };

  const copy = async (text: string, key: string, name: string) => {
    try {
      await navigator.clipboard.writeText(text);
    } catch {
      const ta = document.createElement('textarea');
      ta.value = text;
      document.body.appendChild(ta);
      ta.select();
      document.execCommand('copy');
      document.body.removeChild(ta);
    }
    setCopied(key);
    announce(`Copied ${name} curl example to clipboard`);
    setTimeout(() => setCopied(null), 1500);
  };

  if (!data) return <div class="loading">Loading providers...</div>;

  const activeProvider = data.active || 'none';

  return (
    <section aria-label="Providers">
      <h2>Providers</h2>

      <Card className="mb-4">
        <div class="flex items-center gap-2">
          <span class="muted">Active provider:</span>
          <Badge variant={data.active ? 'active' : 'neutral'}>{activeProvider}</Badge>
        </div>
        <div class="mt-3">
          <a href="/_admin/stripe/webhooks" class="flex items-center gap-1 text-sm" target="_blank" rel="noopener noreferrer">
            Stripe webhook config
            <Icon name="external" size={14} />
          </a>
        </div>
      </Card>

      <div class="provider-grid">
        {(data.available ?? []).map((name) => {
          const info = data.providers?.[name];
          const enabled = (data.enabled ?? []).includes(name);
          const isActive = data.active === name;
          const display = info?.display_name ? `${info.display_name} (${name})` : name;

          const example = info?.sample_method && info?.sample_route
            ? curlExample(info.sample_method, info.sample_route)
            : null;

          const otherVersions = info?.versions?.filter((v) => v !== info.version) ?? [];

          return (
            <Card key={name} className="provider-card">
              <div class="flex justify-between items-start gap-3">
                <div class="flex items-center gap-2">
                  <Icon name="provider" size={18} />
                  <strong dangerouslySetInnerHTML={{ __html: escapeHtml(display) }} />
                  {info?.is_recommended_for_first_time && (
                    <Badge variant="info">recommended</Badge>
                  )}
                </div>
                <label class="toggle flex items-center gap-2 flex-shrink-0">
                  <span class="sr-only">Enable {display}</span>
                  <input
                    type="checkbox"
                    checked={enabled}
                    disabled={saving === name}
                    onChange={(e) => {
                      const target = e.currentTarget as HTMLInputElement;
                      void toggleProvider(name, target.checked);
                    }}
                    aria-label={`Enable ${display}`}
                  />
                  <span aria-hidden="true">{enabled ? 'On' : 'Off'}</span>
                </label>
              </div>

              <div class="flex flex-wrap gap-2 mt-2">
                {isActive && <Badge variant="active">active</Badge>}
                {enabled ? <Badge variant="success">enabled</Badge> : <Badge variant="neutral">disabled</Badge>}
              </div>

              <div class="card-body mt-2">
                {info?.description && (
                  <div class="muted text-sm">{info.description}</div>
                )}
                {info?.real_providers && info.real_providers.length > 0 && (
                  <div class="muted text-sm mt-1">
                    Emulates: {info.real_providers.join(', ')}
                  </div>
                )}
                {example && (
                  <>
                    <div class="muted text-sm mt-2">
                      Try:{' '}
                      <code>
                        {escapeHtml(info?.sample_method ?? 'POST')} {escapeHtml(info?.sample_route ?? '')}
                      </code>
                    </div>
                    <Button
                      variant="secondary"
                      size="sm"
                      icon={copied === name ? 'check' : 'copy'}
                      className="mt-2"
                      onClick={() => copy(example, name, display)}
                      aria-label={`Copy ${display} curl example`}
                    >
                      {copied === name ? 'Copied!' : 'Copy curl'}
                    </Button>
                  </>
                )}
                {info?.docs_path && (
                  <div class="text-sm mt-2">
                    <a href={info.docs_path} class="flex items-center gap-1" target="_blank" rel="noopener noreferrer">
                      Provider docs
                      <Icon name="external" size={14} />
                    </a>
                  </div>
                )}
                {info?.version && (
                  <div class="muted text-sm mt-2">
                    Version: <code>{info.version}</code>
                    {info.versions && info.versions.length > 1 && (
                      <> · supported: {info.versions.map((v) => <code key={v}>{v}</code>)}</>
                    )}
                    {otherVersions.length > 0 && (
                      <div class="mt-1">
                        To pilot another version, send traffic to{' '}
                        <code>/{name}/{otherVersions[0]}/...</code> or set{' '}
                        <code>providers.{name}.config.version</code>.
                      </div>
                    )}
                  </div>
                )}
              </div>
            </Card>
          );
        })}
      </div>
    </section>
  );
}
