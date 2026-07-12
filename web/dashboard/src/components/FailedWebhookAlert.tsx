import { useEffect, useState } from 'preact/hooks';
import { apiGet } from '../api';
import { Button } from './Button';
import { Icon } from './Icon';
import type { FailedWebhooksResponse } from '../types';

interface FailedWebhookAlertProps {
  onShowWebhooks: () => void;
  enabled?: boolean;
}

export function FailedWebhookAlert({ onShowWebhooks, enabled }: FailedWebhookAlertProps) {
  const [dismissed, setDismissed] = useState(() => sessionStorage.getItem('muara_failed_webhook_dismissed') === 'true');
  const [hasFailed, setHasFailed] = useState(false);

  useEffect(() => {
    if (dismissed || !enabled) return;
    apiGet<FailedWebhooksResponse>('/_admin/webhooks?status=failed&limit=1')
      .then((data) => setHasFailed((data.results ?? []).length > 0))
      .catch(() => setHasFailed(false));
  }, [dismissed, enabled]);

  if (dismissed || !hasFailed) return null;

  return (
    <div class="alert active" role="alert">
      <Icon name="alert" size={20} />
      <span class="flex items-center gap-2 flex-wrap">
        <span class="sr-only">Warning:</span>
        <strong>Failed webhook detected.</strong>
        <span>Check the{' '}
          <button
            class="link-button"
            onClick={onShowWebhooks}
          >
            Webhooks
          </button>{' '}
          tab.
        </span>
      </span>
      <Button
        variant="ghost"
        size="sm"
        className="ml-auto"
        onClick={() => {
          sessionStorage.setItem('muara_failed_webhook_dismissed', 'true');
          setDismissed(true);
        }}
      >
        Dismiss
      </Button>
    </div>
  );
}
