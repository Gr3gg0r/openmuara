import { useEffect, useMemo, useState } from 'preact/hooks';
import { apiGet, getRole } from './api';
import { AnnounceRegion } from './components/Announce';
import { AppShell } from './components/AppShell';
import { CommandPalette } from './components/CommandPalette';
import { ErrorBoundary } from './components/ErrorBoundary';
import { FailedWebhookAlert } from './components/FailedWebhookAlert';
import { useConnectionStatus } from './hooks/useConnectionStatus';
import { initializeAppearance, listenToOSThemeChange, syncThemeAcrossTabs, toggleTheme } from './theme';
import type { LedgerEvent, OnboardingResponse, WebhookAttempt } from './types';
import { LedgerDetail } from './views/LedgerDetail';
import { LedgerView } from './views/Ledger';
import { ProviderDetail } from './views/ProviderDetail';
import { SettingsView } from './views/Settings';
import { WebhookDetail } from './views/WebhookDetail';
import { WebhooksView } from './views/Webhooks';

export type NavItem = 'ledger' | 'webhooks' | 'settings';
export type DetailView =
  | { view: 'ledger-detail'; event: LedgerEvent }
  | { view: 'webhook-detail'; webhook: WebhookAttempt }
  | { view: 'provider-detail'; name: string };

const NAV_TITLES: Record<NavItem, string> = {
  ledger: 'Ledger',
  webhooks: 'Webhooks',
  settings: 'Settings',
};

function isTyping(el: EventTarget | null): boolean {
  if (!(el instanceof HTMLElement)) return false;
  const tag = el.tagName;
  return tag === 'INPUT' || tag === 'TEXTAREA' || tag === 'SELECT';
}

function hasModifier(e: KeyboardEvent): boolean {
  return e.ctrlKey || e.altKey || e.metaKey;
}

function isNavItem(value: string): value is NavItem {
  return value === 'ledger' || value === 'webhooks' || value === 'settings';
}

function readNavFromUrl(): NavItem {
  const params = new URLSearchParams(window.location.search);
  const view = params.get('view');
  if (view && isNavItem(view)) return view;
  if (view === 'ledger-detail') return 'ledger';
  if (view === 'webhook-detail') return 'webhooks';
  if (view === 'provider-detail') return 'settings';
  // Fallback for legacy `tab` query parameter.
  const tab = params.get('tab');
  if (tab === 'webhooks') return 'webhooks';
  if (tab === 'settings') return 'settings';
  if (tab === 'ledger' || tab === 'overview' || tab === 'transactions') return 'ledger';
  return 'ledger';
}

function readDetailFromUrl(): DetailView | null {
  const params = new URLSearchParams(window.location.search);
  const view = params.get('view');
  if (view === 'ledger-detail') {
    const type = params.get('type') as LedgerEvent['type'] | null;
    const ref = params.get('ref');
    if (!ref || (type !== 'transaction' && type !== 'webhook')) return null;
    return {
      view: 'ledger-detail',
      event: {
        id: ref,
        reference: ref,
        type,
        time: params.get('time') ?? undefined,
        provider: params.get('provider') ?? undefined,
        status: params.get('status') ?? undefined,
        summary: params.get('summary') ?? undefined,
      },
    };
  }
  if (view === 'webhook-detail') {
    const ref = params.get('ref');
    if (!ref) return null;
    return {
      view: 'webhook-detail',
      webhook: {
        ref,
        provider: params.get('provider') ?? undefined,
        provider_name: params.get('provider_name') ?? undefined,
        url: params.get('url') ?? undefined,
        status: params.get('status') ?? undefined,
        attempts: params.get('attempts') ? Number(params.get('attempts')) : undefined,
      },
    };
  }
  if (view === 'provider-detail') {
    const name = params.get('provider');
    if (!name) return null;
    return { view: 'provider-detail', name };
  }
  return null;
}

function writeNavToUrl(nav: NavItem, detail?: DetailView | null): void {
  const url = new URL(window.location.href);
  if (detail?.view === 'ledger-detail') {
    url.searchParams.set('view', 'ledger-detail');
    url.searchParams.set('type', detail.event.type);
    url.searchParams.set('ref', detail.event.reference);
    if (detail.event.time) url.searchParams.set('time', detail.event.time);
    if (detail.event.provider) url.searchParams.set('provider', detail.event.provider);
    if (detail.event.status) url.searchParams.set('status', detail.event.status);
    if (detail.event.summary) url.searchParams.set('summary', detail.event.summary);
    url.searchParams.delete('provider_name');
    url.searchParams.delete('url');
    url.searchParams.delete('attempts');
  } else if (detail?.view === 'webhook-detail') {
    url.searchParams.set('view', 'webhook-detail');
    url.searchParams.set('ref', detail.webhook.ref);
    if (detail.webhook.provider) url.searchParams.set('provider', detail.webhook.provider);
    if (detail.webhook.provider_name) url.searchParams.set('provider_name', detail.webhook.provider_name);
    if (detail.webhook.url) url.searchParams.set('url', detail.webhook.url);
    if (detail.webhook.status) url.searchParams.set('status', detail.webhook.status);
    if (detail.webhook.attempts != null) url.searchParams.set('attempts', String(detail.webhook.attempts));
  } else if (detail?.view === 'provider-detail') {
    url.searchParams.set('view', 'provider-detail');
    url.searchParams.set('provider', detail.name);
    url.searchParams.delete('type');
    url.searchParams.delete('ref');
    url.searchParams.delete('time');
    url.searchParams.delete('provider_name');
    url.searchParams.delete('url');
    url.searchParams.delete('status');
    url.searchParams.delete('summary');
    url.searchParams.delete('attempts');
  } else if (nav === 'ledger') {
    url.searchParams.delete('view');
    url.searchParams.delete('type');
    url.searchParams.delete('ref');
    url.searchParams.delete('time');
    url.searchParams.delete('provider');
    url.searchParams.delete('provider_name');
    url.searchParams.delete('url');
    url.searchParams.delete('status');
    url.searchParams.delete('summary');
    url.searchParams.delete('attempts');
  } else {
    url.searchParams.set('view', nav);
    url.searchParams.delete('type');
    url.searchParams.delete('ref');
    url.searchParams.delete('time');
    url.searchParams.delete('provider');
    url.searchParams.delete('provider_name');
    url.searchParams.delete('url');
    url.searchParams.delete('status');
    url.searchParams.delete('summary');
    url.searchParams.delete('attempts');
  }
  url.searchParams.delete('tab');
  window.history.replaceState({}, '', url.toString());
}

export function App() {
  const [nav, setNav] = useState<NavItem>(() => {
    const initial = readNavFromUrl();
    return initial === 'settings' && getRole() === 'viewer' ? 'ledger' : initial;
  });
  const [detail, setDetail] = useState<DetailView | null>(() => {
    const initial = readDetailFromUrl();
    return initial?.view === 'provider-detail' && getRole() === 'viewer' ? null : initial;
  });
  const [showHelp, setShowHelp] = useState(false);
  const [showCommandPalette, setShowCommandPalette] = useState(false);
  const [onboarding, setOnboarding] = useState<OnboardingResponse | null>(null);
  const [reloadKey, setReloadKey] = useState(0);
  const connectionStatus = useConnectionStatus();
  const role = useMemo(() => getRole(), []);
  const isViewer = role === 'viewer';

  const reload = () => {
    setReloadKey((k) => k + 1);
  };

  const navigate = (next: NavItem) => {
    if (next === 'settings' && isViewer) {
      return;
    }
    setNav(next);
    setDetail(null);
    writeNavToUrl(next, null);
  };

  const showLedgerDetail = (event: LedgerEvent) => {
    const next: DetailView = { view: 'ledger-detail', event };
    setDetail(next);
    writeNavToUrl('ledger', next);
  };

  const backToLedger = () => {
    setDetail(null);
    writeNavToUrl('ledger', null);
  };

  const showWebhookDetail = (webhook: WebhookAttempt) => {
    const next: DetailView = { view: 'webhook-detail', webhook };
    setDetail(next);
    writeNavToUrl('webhooks', next);
  };

  const backToWebhooks = () => {
    setDetail(null);
    writeNavToUrl('webhooks', null);
  };

  const showProviderDetail = (name: string) => {
    const next: DetailView = { view: 'provider-detail', name };
    setDetail(next);
    writeNavToUrl('settings', next);
  };

  const backToSettings = () => {
    setDetail(null);
    writeNavToUrl('settings', null);
  };

  const focusSearch = () => {
    if (detail) return;
    const id = nav === 'ledger' ? 'ledger-search' : nav === 'webhooks' ? 'webhooks-search' : undefined;
    const el = id ? (document.getElementById(id) as HTMLInputElement | null) : null;
    el?.focus();
  };

  const pageTitle = useMemo(() => {
    if (detail?.view === 'ledger-detail') return 'Ledger detail';
    if (detail?.view === 'webhook-detail') return 'Webhook detail';
    if (detail?.view === 'provider-detail') return `${detail.name} settings`;
    return NAV_TITLES[nav];
  }, [nav, detail]);

  useEffect(() => {
    document.title = `${pageTitle} · OpenMuara Dashboard`;
  }, [pageTitle]);

  useEffect(() => {
    apiGet<OnboardingResponse>('/_admin/onboarding')
      .then(setOnboarding)
      .catch(() => setOnboarding({}));
  }, [reloadKey]);

  useEffect(() => {
    initializeAppearance();
    const cleanupOS = listenToOSThemeChange();
    const cleanupSync = syncThemeAcrossTabs();
    return () => {
      cleanupOS();
      cleanupSync();
    };
  }, []);

  useEffect(() => {
    const onKeyDown = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'k') {
        e.preventDefault();
        setShowCommandPalette((v) => !v);
      }
      if (e.key === '?' && !isTyping(e.target) && !hasModifier(e)) {
        e.preventDefault();
        setShowHelp((v) => !v);
      }
      if (e.key === 'Escape') {
        setShowHelp(false);
        setShowCommandPalette(false);
      }
      if (e.key === '/' && !isTyping(e.target) && !hasModifier(e)) {
        e.preventDefault();
        focusSearch();
      }
      if (e.key === 'd' && !isTyping(e.target) && !hasModifier(e)) {
        e.preventDefault();
        toggleTheme();
      }
      if (e.key === 'r' && !isTyping(e.target) && !hasModifier(e)) {
        e.preventDefault();
        reload();
      }
      if (e.key === '1' && !isTyping(e.target) && !hasModifier(e)) navigate('ledger');
      if (e.key === '2' && !isTyping(e.target) && !hasModifier(e)) navigate('webhooks');
      if (e.key === '3' && !isTyping(e.target) && !hasModifier(e) && !isViewer) navigate('settings');
    };
    document.addEventListener('keydown', onKeyDown);
    return () => document.removeEventListener('keydown', onKeyDown);
  }, [nav, detail]);

  return (
    <ErrorBoundary>
      <AnnounceRegion />
      <AppShell
        active={nav}
        onNavigate={navigate}
        showHelp={showHelp}
        onToggleHelp={() => setShowHelp((v) => !v)}
        onReload={reload}
        reloadKey={reloadKey}
        connectionStatus={connectionStatus}
        role={role}
      >
        <FailedWebhookAlert
          onShowWebhooks={() => navigate('webhooks')}
          enabled={onboarding?.webhooks_enabled}
        />
        {detail?.view === 'ledger-detail' && (
          <LedgerDetail event={detail.event} onBack={backToLedger} />
        )}
        {detail?.view === 'webhook-detail' && (
          <WebhookDetail webhook={detail.webhook} onBack={backToWebhooks} />
        )}
        {detail?.view === 'provider-detail' && !isViewer && (
          <ProviderDetail name={detail.name} onBack={backToSettings} />
        )}
        {!detail && nav === 'ledger' && <LedgerView onShowDetail={showLedgerDetail} />}
        {!detail && nav === 'webhooks' && <WebhooksView onShowDetail={showWebhookDetail} />}
        {!detail && nav === 'settings' && !isViewer && <SettingsView onShowProvider={showProviderDetail} />}
        {!detail && nav === 'settings' && isViewer && (
          <div class="error-banner" role="alert">
            Settings are only available to admin users.
          </div>
        )}
      </AppShell>
      <CommandPalette
        open={showCommandPalette}
        onClose={() => setShowCommandPalette(false)}
        onNavigate={navigate}
        onReload={reload}
        onFocusSearch={focusSearch}
      />
    </ErrorBoundary>
  );
}
