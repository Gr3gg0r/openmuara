import { useEffect, useState } from 'preact/hooks';
import { announce } from './Announce';
import { apiGet, apiPatch, apiPost } from '../api';
import { Badge } from './Badge';
import { Button } from './Button';
import { useAsync } from '../hooks/useAsync';
import type { WebhookConfigResponse, ProvidersResponse } from '../types';

const DEFAULT_EVENTS: Record<string, string[]> = {
  stripe: ['checkout.session.completed', 'payment_intent.succeeded', 'payment_intent.payment_failed'],
  fawry: ['charge_paid'],
};

export function WebhookConfig() {
  const { data: config, refetch } = useAsync<WebhookConfigResponse>(() => apiGet('/_admin/config/webhooks'), []);
  const { data: providers } = useAsync<ProvidersResponse>(() => apiGet('/_admin/providers'), []);
  const [url, setUrl] = useState('');
  const [targets, setTargets] = useState<Record<string, string>>({});
  const [events, setEvents] = useState<Record<string, string[]>>({});
  const [testSecrets, setTestSecrets] = useState<Record<string, string>>({});
  const [showSecrets, setShowSecrets] = useState<Record<string, boolean>>({});
  const [saving, setSaving] = useState(false);
  const [testing, setTesting] = useState<string | null>(null);
  const [testResult, setTestResult] = useState<{ ok: boolean; message: string } | null>(null);

  useEffect(() => {
    if (config) {
      setUrl(config.url ?? '');
      setTargets(config.targets ?? {});
      setEvents(config.events ?? {});
    }
  }, [config]);

  const enabledProviders = providers?.enabled ?? [];

  const save = async (e: Event) => {
    e.preventDefault();
    setSaving(true);
    try {
      await apiPatch('/_admin/config/webhooks', {
        url: url || undefined,
        targets: Object.keys(targets).length ? targets : undefined,
        events: Object.keys(events).length ? events : undefined,
      });
      announce('Webhook config saved; restart server to activate target changes');
      await refetch();
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      announce(`Failed to save webhook config: ${message}`);
    } finally {
      setSaving(false);
    }
  };

  const testURL = async (targetUrl: string, provider: string) => {
    setTesting(provider);
    setTestResult(null);
    try {
      const secret = testSecrets[provider];
      const res = await apiPost('/_admin/config/webhooks/test', { url: targetUrl, provider, secret });
      const body = (await res.json()) as { success: boolean; latency_ms?: number; error?: string };
      setTestResult({
        ok: body.success,
        message: body.success
          ? `Delivered in ${body.latency_ms ?? '?'}ms`
          : body.error || 'Failed',
      });
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      setTestResult({ ok: false, message });
    } finally {
      setTesting(null);
    }
  };

  const toggleEvent = (provider: string, event: string) => {
    setEvents((prev) => {
      const current = new Set(prev[provider] ?? []);
      if (current.has(event)) {
        current.delete(event);
      } else {
        current.add(event);
      }
      return { ...prev, [provider]: Array.from(current) };
    });
  };

  const setSecret = (provider: string, value: string) => {
    setTestSecrets((prev) => ({ ...prev, [provider]: value }));
  };

  const toggleSecretVisibility = (provider: string) => {
    setShowSecrets((prev) => ({ ...prev, [provider]: !prev[provider] }));
  };

  return (
    <section aria-label="Webhook configuration" class="card card-padded mb-4">
      <h2 class="mt-0">Webhook Configuration</h2>
      <form onSubmit={save}>
        <div class="mb-3">
          <label class="label" for="webhook-global-url">Global webhook URL</label>
          <input
            id="webhook-global-url"
            type="url"
            class="input input-full"
            value={url}
            onInput={(e) => setUrl((e.currentTarget as HTMLInputElement).value)}
            placeholder="https://example.com/webhook"
          />
          <div class="muted text-sm mt-1">
            Used for all providers unless a per-provider target is set.
          </div>
        </div>

        {enabledProviders.map((name) => {
          const display = providers?.providers?.[name]?.display_name ?? name;
          const providerEvents = DEFAULT_EVENTS[name] ?? [];
          return (
            <fieldset key={name} class="fieldset">
              <legend>{display}</legend>
              <div class="mb-2">
                <label class="label" for={`webhook-target-${name}`}>Target URL</label>
                <input
                  id={`webhook-target-${name}`}
                  type="url"
                  class="input input-full"
                  value={targets[name] ?? ''}
                  onInput={(e) =>
                    setTargets((prev) => ({ ...prev, [name]: (e.currentTarget as HTMLInputElement).value }))
                  }
                  placeholder="Leave empty to use global URL"
                />
              </div>
              <div class="mb-2">
                <label class="label" for={`webhook-secret-${name}`}>Test signing secret</label>
                <div class="flex gap-2">
                  <input
                    id={`webhook-secret-${name}`}
                    type={showSecrets[name] ? 'text' : 'password'}
                    class="input input-full"
                    value={testSecrets[name] ?? ''}
                    onInput={(e) => setSecret(name, (e.currentTarget as HTMLInputElement).value)}
                    placeholder="Optional secret for test signature"
                  />
                  <Button
                    variant="secondary"
                    size="sm"
                    icon={showSecrets[name] ? 'eyeOff' : 'eye'}
                    onClick={() => toggleSecretVisibility(name)}
                    aria-label={`${showSecrets[name] ? 'Hide' : 'Show'} ${display} secret`}
                  >
                    {showSecrets[name] ? 'Hide' : 'Show'}
                  </Button>
                </div>
                <Button
                  variant="secondary"
                  size="sm"
                  className="mt-2"
                  loading={testing === name}
                  icon="provider"
                  onClick={() => testURL(targets[name] || url, name)}
                >
                  Test delivery
                </Button>
                {testResult && testing === null && (
                  <Badge variant={testResult.ok ? 'success' : 'danger'} className="ml-2">
                    {testResult.message}
                  </Badge>
                )}
              </div>
              {providerEvents.length > 0 && (
                <div>
                  <div class="mb-1">Enabled events</div>
                  {providerEvents.map((event) => (
                    <label key={event} class="label checkbox-label">
                      <input
                        type="checkbox"
                        checked={(events[name] ?? []).includes(event)}
                        onChange={() => toggleEvent(name, event)}
                      />{' '}
                      {event}
                    </label>
                  ))}
                </div>
              )}
            </fieldset>
          );
        })}

        <div>
          <Button type="submit" loading={saving}>
            Save webhook config
          </Button>
        </div>
      </form>
    </section>
  );
}
