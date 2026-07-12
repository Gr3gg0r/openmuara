import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { apiGet, apiPost, apiPatch, escapeHtml, formatDate, statusClass, curlExample, getFallbackCsrfToken, getRole, isAdmin } from '../src/api';

describe('escapeHtml', () => {
  it('escapes special HTML characters', () => {
    expect(escapeHtml('<script>alert("x")\'s</script>')).toBe(
      '&lt;script&gt;alert(&quot;x&quot;)&#39;s&lt;/script&gt;',
    );
  });
});

describe('formatDate', () => {
  it('formats an ISO date', () => {
    const s = '2026-07-03T04:00:00.000Z';
    expect(formatDate(s)).toContain('2026');
  });

  it('returns dash for missing date', () => {
    expect(formatDate(undefined)).toBe('-');
  });
});

describe('statusClass', () => {
  it('prefixes status with status-', () => {
    expect(statusClass('paid')).toBe('status-paid');
    expect(statusClass(undefined)).toBe('status-unknown');
  });
});

describe('curlExample', () => {
  it('builds a POST curl example', () => {
    expect(curlExample('POST', '/stripe/v1/charges')).toBe(
      "curl -X POST http://127.0.0.1:9000/stripe/v1/charges -H 'Content-Type: application/json' -d '{}'",
    );
  });

  it('builds a GET curl example without body', () => {
    expect(curlExample('GET', '/healthz')).toBe('curl -X GET http://127.0.0.1:9000/healthz');
  });
});

describe('apiGet', () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('returns parsed JSON on success', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      json: vi.fn().mockResolvedValue({ status: 'ok' }),
    });
    globalThis.fetch = fetchMock as unknown as typeof fetch;

    const result = await apiGet<{ status: string }>('/healthz');
    expect(result).toEqual({ status: 'ok' });
    expect(fetchMock).toHaveBeenCalledWith('http://localhost:3000/healthz', {
      headers: { Accept: 'application/json' },
    });
  });

  it('uses muara-admin-api meta tag as base URL when present', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      json: vi.fn().mockResolvedValue({ status: 'ok' }),
    });
    globalThis.fetch = fetchMock as unknown as typeof fetch;

    document.head.innerHTML = '<meta name="muara-admin-api" content="http://127.0.0.1:9001">';

    const result = await apiGet<{ status: string }>('/_admin/transactions');
    expect(result).toEqual({ status: 'ok' });
    expect(fetchMock).toHaveBeenCalledWith('http://127.0.0.1:9001/_admin/transactions', {
      headers: { Accept: 'application/json' },
    });

    document.head.innerHTML = '';
  });

  it('strips embedded credentials from the resolved URL', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      json: vi.fn().mockResolvedValue({ status: 'ok' }),
    });
    globalThis.fetch = fetchMock as unknown as typeof fetch;

    const originalHref = window.location.href;
    Object.defineProperty(window, 'location', {
      value: { href: 'http://admin:secret@127.0.0.1:9000/_admin/' },
      configurable: true,
    });

    await apiGet('/_admin/ledger');
    expect(fetchMock).toHaveBeenCalledWith('http://127.0.0.1:9000/_admin/ledger', {
      headers: { Accept: 'application/json' },
    });

    Object.defineProperty(window, 'location', {
      value: { href: originalHref },
      configurable: true,
    });
  });

  it('throws ApiError on failure', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: false,
      status: 500,
      text: vi.fn().mockResolvedValue('boom'),
    });
    globalThis.fetch = fetchMock as unknown as typeof fetch;

    await expect(apiGet('/x')).rejects.toThrow('boom');
  });
});

describe('apiPost', () => {
  beforeEach(() => {
    document.head.innerHTML = '<meta name="csrf-token" content="tok123">';
  });

  afterEach(() => {
    vi.restoreAllMocks();
    document.head.innerHTML = '';
  });

  it('sends CSRF token header', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      status: 200,
    });
    globalThis.fetch = fetchMock as unknown as typeof fetch;

    await apiPost('/_admin/replay');
    const call = fetchMock.mock.calls[0];
    expect(call[1].headers['X-CSRF-Token']).toBe('tok123');
  });

  it('strips embedded credentials from the resolved URL', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      status: 200,
    });
    globalThis.fetch = fetchMock as unknown as typeof fetch;

    Object.defineProperty(window, 'location', {
      value: { href: 'http://admin:secret@127.0.0.1:9000/_admin/' },
      configurable: true,
    });

    await apiPost('/_admin/replay');
    expect(fetchMock).toHaveBeenCalledWith(
      'http://127.0.0.1:9000/_admin/replay',
      expect.objectContaining({ headers: expect.objectContaining({ 'X-CSRF-Token': 'tok123' }) }),
    );
  });
});

describe('apiPatch', () => {
  beforeEach(() => {
    document.head.innerHTML = '<meta name="csrf-token" content="tok123">';
  });

  afterEach(() => {
    vi.restoreAllMocks();
    document.head.innerHTML = '';
  });

  it('sends PATCH method with JSON body', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      status: 202,
    });
    globalThis.fetch = fetchMock as unknown as typeof fetch;

    await apiPatch('/_admin/config/providers', { providers: { fawry: { enabled: false } } });
    const call = fetchMock.mock.calls[0];
    expect(call[1].method).toBe('PATCH');
    expect(call[1].headers['Content-Type']).toBe('application/json');
    expect(call[1].headers['X-CSRF-Token']).toBe('tok123');
    expect(call[1].body).toBe('{"providers":{"fawry":{"enabled":false}}}');
  });
});

describe('getFallbackCsrfToken', () => {
  it('reads cookie when meta tag is absent', () => {
    document.cookie = 'openmuara_csrf=fallback; path=/';
    expect(getFallbackCsrfToken()).toBe('fallback');
    document.cookie = 'openmuara_csrf=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/';
  });
});

describe('getRole', () => {
  afterEach(() => {
    document.head.innerHTML = '';
  });

  it('returns admin when meta tag is admin', () => {
    document.head.innerHTML = '<meta name="muara-role" content="admin">';
    expect(getRole()).toBe('admin');
    expect(isAdmin()).toBe(true);
  });

  it('returns viewer when meta tag is viewer', () => {
    document.head.innerHTML = '<meta name="muara-role" content="viewer">';
    expect(getRole()).toBe('viewer');
    expect(isAdmin()).toBe(false);
  });

  it('returns empty string when meta tag is missing', () => {
    document.head.innerHTML = '';
    expect(getRole()).toBe('');
    expect(isAdmin()).toBe(false);
  });

  it('returns empty string for unknown role', () => {
    document.head.innerHTML = '<meta name="muara-role" content="superuser">';
    expect(getRole()).toBe('');
    expect(isAdmin()).toBe(false);
  });
});
