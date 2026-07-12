import { useCallback, useEffect, useState } from 'preact/hooks';

export interface UrlStateOptions {
  replace?: boolean;
}

function readParams(): URLSearchParams {
  return new URLSearchParams(window.location.search);
}

function getParam(key: string): string | undefined {
  const value = readParams().get(key);
  return value ?? undefined;
}

function setParams(params: URLSearchParams, options?: UrlStateOptions): void {
  const url = new URL(window.location.href);
  url.search = params.toString();
  window.history[options?.replace ? 'replaceState' : 'pushState']({}, '', url.toString());
}

export function useUrlState<T extends string>(key: string, defaultValue?: T): [T, (value: T) => void] {
  const [value, setValue] = useState<T>(() => (getParam(key) as T | undefined) ?? defaultValue ?? ('' as T));

  useEffect(() => {
    const onPopState = () => {
      setValue((getParam(key) as T | undefined) ?? defaultValue ?? ('' as T));
    };
    window.addEventListener('popstate', onPopState);
    return () => window.removeEventListener('popstate', onPopState);
  }, [key, defaultValue]);

  const update = useCallback(
    (next: T) => {
      const params = readParams();
      if (next && next !== defaultValue) {
        params.set(key, next);
      } else {
        params.delete(key);
      }
      setParams(params);
      setValue(next);
    },
    [key, defaultValue],
  );

  return [value, update];
}

export function useUrlStateSynced<T extends string>(key: string, value: T, setValue: (value: T) => void, defaultValue?: T): void {
  useEffect(() => {
    const initial = getParam(key);
    if (initial !== undefined && initial !== value) {
      setValue(initial as T);
    }
    // Only run once on mount.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    const params = readParams();
    if (value && value !== defaultValue) {
      params.set(key, value);
    } else {
      params.delete(key);
    }
    setParams(params, { replace: true });
  }, [key, value, defaultValue]);
}
