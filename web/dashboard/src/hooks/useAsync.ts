import { useCallback, useEffect, useState } from 'preact/hooks';

interface AsyncState<T> {
  data: T | undefined;
  loading: boolean;
  error: string | undefined;
}

export function useAsync<T>(fetcher: () => Promise<T>, deps: unknown[] = []) {
  const [state, setState] = useState<AsyncState<T>>({ data: undefined, loading: true, error: undefined });

  const run = useCallback(() => {
    setState((s) => ({ ...s, loading: true, error: undefined }));
    fetcher()
      .then((data) => setState({ data, loading: false, error: undefined }))
      .catch((err: Error) => setState({ data: undefined, loading: false, error: err.message }));
  }, deps);

  useEffect(() => {
    run();
  }, [run]);

  return { ...state, refetch: run };
}
