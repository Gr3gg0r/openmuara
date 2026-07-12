import { useAsync } from '../hooks/useAsync';
import { apiGet, curlExample } from '../api';
import { Badge } from './Badge';
import { Button } from './Button';
import { Card } from './Card';
import { EmptyState } from './EmptyState';
import { SkeletonCards } from './Skeleton';
import type { ConfigResponse, OnboardingResponse, ProvidersResponse } from '../types';

interface OverviewProps {
  onboarding: OnboardingResponse | null;
  onShowLedger: () => void;
  onShowProviders: () => void;
  onShowWebhooks: () => void;
}

export function Overview({ onboarding, onShowLedger, onShowProviders, onShowWebhooks }: OverviewProps) {
  const { data: providers, loading: providersLoading } = useAsync<ProvidersResponse>(() => apiGet('/_admin/providers'), []);
  const { data: config } = useAsync<ConfigResponse>(() => apiGet('/_admin/config'), []);

  const enabled = providers?.enabled ?? [];
  const available = providers?.available ?? [];
  const sampleRoute = (() => {
    for (const name of available) {
      const info = providers?.providers?.[name];
      if (info?.sample_method && info?.sample_route) {
        return `${info.sample_method} ${info.sample_route}`;
      }
    }
    return 'POST /stripe/v1/charges';
  })();

  const steps = [
    { label: 'Enable a provider', done: enabled.length > 0, action: onShowProviders },
    { label: 'Send a test charge', done: onboarding?.first_transaction ?? false, action: onShowLedger },
    { label: 'Configure a webhook', done: (config?.webhook?.url ?? '').length > 0, action: onShowWebhooks },
    { label: 'Receive or replay a webhook', done: onboarding?.first_webhook_received ?? false, action: onShowWebhooks },
  ];

  return (
    <section aria-label="Overview">
      <h2>Overview</h2>

      <Card className="mb-4">
        <h3 class="mt-0">Getting started</h3>
        <ol class="onboarding-list">
          {steps.map((step) => (
            <li key={step.label} class="flex items-center gap-2 onboarding-item">
              <Badge variant={step.done ? 'success' : 'neutral'}>{step.done ? 'Done' : 'Todo'}</Badge>
              <span>{step.label}</span>
              {!step.done && (
                <Button variant="secondary" size="sm" className="ml-auto" onClick={step.action}>
                  Start
                </Button>
              )}
            </li>
          ))}
        </ol>
      </Card>

      <h3>Providers</h3>
      {providersLoading ? (
        <SkeletonCards count={3} />
      ) : available.length === 0 ? (
        <EmptyState
          title="No providers available"
          description="Check your configuration and restart the server."
          icon="provider"
        />
      ) : (
        <div class="provider-grid">
          {available.map((name) => {
            const info = providers?.providers?.[name];
            const isEnabled = enabled.includes(name);
            const isActive = providers?.active === name;
            return (
              <Card key={name} className="provider-card">
                <div class="flex justify-between items-start">
                  <strong>{info?.display_name ?? name}</strong>
                  <Badge variant={isEnabled ? 'success' : 'neutral'}>{isEnabled ? 'enabled' : 'disabled'}</Badge>
                </div>
                {isActive && <Badge variant="info" className="mt-2">active</Badge>}
                <div class="card-body muted text-sm mt-2">{info?.description}</div>
                <div class="mt-3">
                  <Button variant="secondary" size="sm" icon="settings" onClick={onShowProviders}>
                    Configure
                  </Button>
                </div>
              </Card>
            );
          })}
        </div>
      )}

      {!onboarding?.first_transaction && (
        <Card className="mt-4">
          <h3 class="mt-0">Simulate your first charge</h3>
          <div class="muted text-sm">Copy and run this curl command:</div>
          <pre class="mt-2">
            <code>{curlExample('POST', sampleRoute.replace(/^POST\s+/, ''))}</code>
          </pre>
        </Card>
      )}
    </section>
  );
}
