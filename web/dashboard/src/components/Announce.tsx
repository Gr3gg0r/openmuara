import { useEffect, useState } from 'preact/hooks';

let announceCallback: ((message: string) => void) | null = null;

export function announce(message: string): void {
  announceCallback?.(message);
}

export function AnnounceRegion() {
  const [message, setMessage] = useState('');

  useEffect(() => {
    announceCallback = (msg: string) => {
      setMessage('');
      requestAnimationFrame(() => setMessage(msg));
    };
    return () => { announceCallback = null; };
  }, []);

  return (
    <div
      aria-live="polite"
      aria-atomic="true"
      class="sr-only"
    >
      {message}
    </div>
  );
}
