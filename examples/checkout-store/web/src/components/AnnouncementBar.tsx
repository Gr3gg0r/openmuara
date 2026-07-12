import { useState } from 'react';
import { AlertTriangle, X } from 'lucide-react';

interface AnnouncementBarProps {
  demoMode: boolean;
}

export default function AnnouncementBar({ demoMode }: AnnouncementBarProps) {
  const [dismissed, setDismissed] = useState(false);

  if (!demoMode || dismissed) {
    return null;
  }

  return (
    <div className="bg-warning text-warning-content">
      <div className="container mx-auto flex items-start gap-3 px-4 py-2 text-sm">
        <AlertTriangle className="mt-0.5 h-4 w-4 shrink-0" />
        <p className="flex-1">
          <span className="font-semibold">Demo mode</span> — running with default
          placeholder credentials. Copy{' '}
          <code className="rounded bg-warning-content/20 px-1">.env.example</code>{' '}
          to <code className="rounded bg-warning-content/20 px-1">.env</code> and set
          your provider keys (for example{' '}
          <code className="rounded bg-warning-content/20 px-1">TOYYIBPAY_USER_SECRET_KEY</code>)
          to test real payments.
        </p>
        <button
          type="button"
          className="btn btn-ghost btn-xs"
          aria-label="Dismiss announcement"
          onClick={() => setDismissed(true)}
        >
          <X className="h-4 w-4" />
        </button>
      </div>
    </div>
  );
}
