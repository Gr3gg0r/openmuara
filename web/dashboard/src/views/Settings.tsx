import { useState } from 'preact/hooks';
import { apiPost, isAdmin } from '../api';
import { Badge } from '../components/Badge';
import { Button } from '../components/Button';
import { Card } from '../components/Card';
import { ConfirmDialog } from '../components/ConfirmDialog';
import { EmptyState } from '../components/EmptyState';
import { Icon } from '../components/Icon';
import { useAsync } from '../hooks/useAsync';
import { apiGet } from '../api';
import type { ProvidersResponse } from '../types';

interface SettingsViewProps {
  onShowProvider?: (name: string) => void;
}

export function SettingsView({ onShowProvider }: SettingsViewProps) {
  const { data, loading, error } = useAsync<ProvidersResponse>(() => apiGet('/_admin/providers'), []);
  const [showClearConfirm, setShowClearConfirm] = useState(false);
  const [clearing, setClearing] = useState(false);
  const [clearError, setClearError] = useState<string | null>(null);

  const providers = data?.available ?? [];
  const enabled = new Set(data?.enabled ?? []);
  const active = data?.active ?? '';
  const admin = isAdmin();

  const handleClear = async () => {
    setClearing(true);
    setClearError(null);
    try {
      await apiPost('/_admin/clean');
      setShowClearConfirm(false);
      window.location.reload();
    } catch (err) {
      setClearing(false);
      setClearError(err instanceof Error ? err.message : 'Clear failed');
    }
  };

  return (
    <section aria-label="Settings">
      <div class="flex justify-between items-center flex-wrap gap-3 mb-3">
        <h2 style={{ margin: 0 }}>Settings</h2>
      </div>

      <p class="muted">
        Enable or disable providers, configure per-provider webhook targets, and view environment variable references.
      </p>

      {error && (
        <div class="error-banner" role="alert">
          Failed to load providers: {error}
        </div>
      )}

      {clearError && (
        <div class="error-banner" role="alert">
          Failed to clear data: {clearError}
        </div>
      )}

      {admin && (
        <Card className="danger-zone mb-4">
          <div class="flex items-center gap-2 mb-2">
            <Icon name="alert" size={18} />
            <h3 style={{ margin: 0 }}>Local data</h3>
          </div>
          <p class="muted">
            Remove all transactions, webhook attempts, and audit events from the running server. This cannot be undone.
          </p>
          <Button
            variant="danger"
            icon="trash"
            onClick={() => setShowClearConfirm(true)}
            disabled={clearing}
            loading={clearing}
          >
            Clear local data
          </Button>
        </Card>
      )}

      {loading && providers.length === 0 ? (
        <div class="grid">
          {Array.from({ length: 4 }).map((_, i) => (
            <Card key={i}>
              <div class="skeleton skeleton-title" />
              <div class="skeleton skeleton-line" />
              <div class="skeleton skeleton-line-short" />
            </Card>
          ))}
        </div>
      ) : providers.length === 0 ? (
        <EmptyState
          title="No providers available"
          description="No providers are registered. Check the server logs."
          icon="provider"
        />
      ) : (
        <div class="grid">
          {providers.map((name) => {
            const info = data?.providers?.[name];
            const isEnabled = enabled.has(name);
            const isActive = active === name;
            return (
              <Card key={name} className="provider-card">
                <div class="flex justify-between items-start gap-2">
                  <h3 class="provider-card-title">{info?.display_name ?? name}</h3>
                  <div class="flex gap-2">
                    {isActive && <Badge variant="active">Active</Badge>}
                    {isEnabled ? <Badge variant="success">Enabled</Badge> : <Badge variant="neutral">Disabled</Badge>}
                  </div>
                </div>
                <p class="muted provider-card-desc">{info?.description ?? 'Provider configuration'}</p>
                {info?.real_providers && info.real_providers.length > 0 && (
                  <p class="text-sm muted">Emulates: {info.real_providers.join(', ')}</p>
                )}
                <div class="provider-card-footer">
                  <Button
                    variant="secondary"
                    size="sm"
                    icon="settings"
                    onClick={() => onShowProvider?.(name)}
                    aria-label={`Configure ${info?.display_name ?? name}`}
                  >
                    Configure
                  </Button>
                </div>
              </Card>
            );
          })}
        </div>
      )}

      <ConfirmDialog
        open={showClearConfirm}
        title="Clear all local data?"
        message="This will permanently delete every transaction, webhook attempt, and audit event from the running server. Your config.yml will not be changed."
        confirmLabel="Yes, clear data"
        cancelLabel="Cancel"
        confirmVariant="danger"
        onConfirm={handleClear}
        onCancel={() => setShowClearConfirm(false)}
      />
    </section>
  );
}
