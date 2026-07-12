function getCsrfToken(): string {
  const meta = document.querySelector<HTMLMetaElement>('meta[name="csrf-token"]');
  return meta?.content ?? '';
}

function getCookie(name: string): string | undefined {
  const match = document.cookie.match(new RegExp('(?:^|; )' + name.replace(/([.$?*|{}()[\]\\/+^])/g, '\\$1') + '=([^;]*)'));
  return match ? decodeURIComponent(match[1]) : undefined;
}

export function getFallbackCsrfToken(): string {
  return getCookie('openmuara_csrf') ?? '';
}

export class ApiError extends Error {
  constructor(
    message: string,
    public status: number,
    public body?: string,
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

/**
 * Resolve a path to an absolute URL for fetch, stripping any credentials that
 * may be present in window.location. The Fetch API rejects URLs containing
 * username/password, which can happen when the admin page is loaded with HTTP
 * Basic Auth credentials embedded in the URL.
 */
function getAdminApiBaseUrl(): string {
  const meta = document.querySelector<HTMLMetaElement>('meta[name="muara-admin-api"]');
  return meta?.content?.trim() ?? '';
}

export function getRole(): 'admin' | 'viewer' | '' {
  const meta = document.querySelector<HTMLMetaElement>('meta[name="muara-role"]');
  const value = meta?.content?.trim() ?? '';
  if (value === 'admin' || value === 'viewer') {
    return value;
  }
  return '';
}

export function isAdmin(): boolean {
  return getRole() === 'admin';
}

function resolveApiUrl(path: string): string {
  const baseUrl = getAdminApiBaseUrl();
  const url = new URL(path, baseUrl || window.location.href);
  url.username = '';
  url.password = '';
  return url.toString();
}

async function parseError(res: Response): Promise<string> {
  const text = await res.text();
  try {
    const json = JSON.parse(text);
    return json.error ?? json.message ?? text;
  } catch {
    return text || res.statusText;
  }
}

export async function apiGet<T>(path: string): Promise<T> {
  const res = await fetch(resolveApiUrl(path), {
    headers: { Accept: 'application/json' },
  });
  if (!res.ok) {
    throw new ApiError(await parseError(res), res.status);
  }
  return res.json() as Promise<T>;
}

export async function apiPost(path: string, body?: object): Promise<Response> {
  const token = getCsrfToken() || getFallbackCsrfToken();
  const headers: Record<string, string> = {
    Accept: 'application/json',
    'X-CSRF-Token': token,
  };
  const init: RequestInit = {
    method: 'POST',
    headers,
  };
  if (body) {
    headers['Content-Type'] = 'application/json';
    init.body = JSON.stringify(body);
  }
  const res = await fetch(resolveApiUrl(path), init);
  if (!res.ok) {
    throw new ApiError(await parseError(res), res.status);
  }
  return res;
}

export async function apiPatch(path: string, body?: object): Promise<Response> {
  const token = getCsrfToken() || getFallbackCsrfToken();
  const headers: Record<string, string> = {
    Accept: 'application/json',
    'X-CSRF-Token': token,
  };
  const init: RequestInit = {
    method: 'PATCH',
    headers,
  };
  if (body) {
    headers['Content-Type'] = 'application/json';
    init.body = JSON.stringify(body);
  }
  const res = await fetch(resolveApiUrl(path), init);
  if (!res.ok) {
    throw new ApiError(await parseError(res), res.status);
  }
  return res;
}

export function formatDate(s?: string): string {
  return s ? new Date(s).toLocaleString() : '-';
}

export function escapeHtml(s: unknown): string {
  return String(s).replace(/[&<>"']/g, (c) => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' }[c] ?? c));
}

export function statusClass(s?: string): string {
  return 'status-' + (s || 'unknown');
}

export function formatJSON(value: unknown): string {
  try {
    return JSON.stringify(value, null, 2);
  } catch {
    return String(value);
  }
}

export async function copyToClipboard(text: string): Promise<boolean> {
  try {
    await navigator.clipboard.writeText(text);
    return true;
  } catch {
    return false;
  }
}

export function curlExample(method?: string, route?: string): string {
  const m = (method || 'POST').toUpperCase();
  const body = m === 'GET' ? '' : " -H 'Content-Type: application/json' -d '{}'";
  return `curl -X ${m} http://127.0.0.1:9000${route}${body}`;
}

export async function getProvider(name: string): Promise<import('./types').ProviderDetailResponse> {
  return apiGet<import('./types').ProviderDetailResponse>(`/_admin/providers/${encodeURIComponent(name)}`);
}
