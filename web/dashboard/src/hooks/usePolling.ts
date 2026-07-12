import { useEffect, useRef } from 'preact/hooks';

export function usePolling(callback: () => void | Promise<void>, intervalMs: number, enabled = true) {
  const callbackRef = useRef(callback);
  callbackRef.current = callback;

  useEffect(() => {
    if (!enabled) return;

    const tick = () => {
      void callbackRef.current();
    };

    tick();
    const id = setInterval(() => {
      if (!document.hidden) tick();
    }, intervalMs);

    const onVisibility = () => {
      if (document.hidden) {
        clearInterval(id);
      } else {
        tick();
      }
    };
    document.addEventListener('visibilitychange', onVisibility);

    return () => {
      clearInterval(id);
      document.removeEventListener('visibilitychange', onVisibility);
    };
  }, [intervalMs, enabled]);
}
