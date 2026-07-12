import { useEffect, useState } from 'preact/hooks';

export type ConnectionStatus = 'online' | 'offline' | 'checking';

export function useConnectionStatus(pingUrl = '/_admin/config'): ConnectionStatus {
  const [status, setStatus] = useState<ConnectionStatus>('checking');

  useEffect(() => {
    let cancelled = false;
    const abort = new AbortController();

    const check = async () => {
      try {
        const res = await fetch(pingUrl, {
          method: 'HEAD',
          signal: abort.signal,
          cache: 'no-store',
        });
        if (!cancelled) {
          setStatus(res.ok ? 'online' : 'offline');
        }
      } catch {
        if (!cancelled) {
          setStatus('offline');
        }
      }
    };

    const handleOnline = () => {
      setStatus('checking');
      void check();
    };
    const handleOffline = () => setStatus('offline');

    window.addEventListener('online', handleOnline);
    window.addEventListener('offline', handleOffline);

    void check();
    const interval = setInterval(check, 30000);

    return () => {
      cancelled = true;
      abort.abort();
      window.removeEventListener('online', handleOnline);
      window.removeEventListener('offline', handleOffline);
      clearInterval(interval);
    };
  }, [pingUrl]);

  return status;
}
