import { Icon, type IconName } from './Icon';

export type NavItem = 'ledger' | 'webhooks' | 'settings';

interface SidebarNavProps {
  active: NavItem;
  onNavigate: (item: NavItem) => void;
  role?: 'admin' | 'viewer' | '';
}

const ITEMS: { key: NavItem; label: string; icon: IconName; shortcut: string; adminOnly?: boolean }[] = [
  { key: 'ledger', label: 'Ledger', icon: 'list', shortcut: '1' },
  { key: 'webhooks', label: 'Webhooks', icon: 'webhook', shortcut: '2' },
  { key: 'settings', label: 'Settings', icon: 'settings', shortcut: '3', adminOnly: true },
];

export function SidebarNav({ active, onNavigate, role }: SidebarNavProps) {
  const isViewer = role === 'viewer';
  return (
    <nav class="sidebar-nav" aria-label="Main">
      <ul class="sidebar-list" role="menubar">
        {ITEMS.filter((item) => !item.adminOnly || !isViewer).map((item) => {
          const isActive = active === item.key;
          return (
            <li key={item.key} role="none">
              <button
                class={`sidebar-link ${isActive ? 'active' : ''}`}
                role="menuitem"
                aria-current={isActive ? 'page' : undefined}
                onClick={() => onNavigate(item.key)}
                data-testid={`nav-${item.key}`}
                title={`${item.label} (${item.shortcut})`}
              >
                <Icon name={item.icon} size={18} />
                <span class="sidebar-label">{item.label}</span>
                <kbd class="sidebar-shortcut">{item.shortcut}</kbd>
              </button>
            </li>
          );
        })}
      </ul>
    </nav>
  );
}
