import { useEffect, useMemo, useRef, useState } from 'preact/hooks';
import { getRole } from '../api';
import { Icon } from './Icon';
import type { NavItem } from './SidebarNav';

interface Command {
  id: string;
  label: string;
  icon: Parameters<typeof Icon>[0]['name'];
  shortcut?: string;
  action: () => void;
}

interface CommandPaletteProps {
  open: boolean;
  onClose: () => void;
  onNavigate: (nav: NavItem) => void;
  onReload?: () => void;
  onFocusSearch?: () => void;
}

export function CommandPalette({ open, onClose, onNavigate, onReload, onFocusSearch }: CommandPaletteProps) {
  const [query, setQuery] = useState('');
  const [selectedIndex, setSelectedIndex] = useState(0);
  const inputRef = useRef<HTMLInputElement>(null);
  const listRef = useRef<HTMLDivElement>(null);

  const commands = useMemo<Command[]>(() => {
    const isViewer = getRole() === 'viewer';
    const list: Command[] = [
      { id: 'nav-ledger', label: 'Go to Ledger', icon: 'list', shortcut: '1', action: () => { onNavigate('ledger'); onClose(); } },
      { id: 'nav-webhooks', label: 'Go to Webhooks', icon: 'webhook', shortcut: '2', action: () => { onNavigate('webhooks'); onClose(); } },
    ];
    if (!isViewer) {
      list.push({ id: 'nav-settings', label: 'Go to Settings', icon: 'settings', shortcut: '3', action: () => { onNavigate('settings'); onClose(); } });
    }
    list.push(
      { id: 'reload', label: 'Reload dashboard data', icon: 'refresh', shortcut: 'r', action: () => { onReload?.(); onClose(); } },
      { id: 'search', label: 'Focus search', icon: 'search', shortcut: '/', action: () => { onFocusSearch?.(); onClose(); } },
    );
    return list;
  }, [onNavigate, onReload, onFocusSearch, onClose]);

  const filtered = useMemo(() => {
    const q = query.toLowerCase().trim();
    if (!q) return commands;
    return commands.filter((c) => c.label.toLowerCase().includes(q));
  }, [commands, query]);

  useEffect(() => {
    setSelectedIndex(0);
  }, [query]);

  useEffect(() => {
    if (open) {
      setQuery('');
      setSelectedIndex(0);
      inputRef.current?.focus();
    }
  }, [open]);

  useEffect(() => {
    if (!open) return;
    const onKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        e.preventDefault();
        onClose();
        return;
      }
      if (e.key === 'ArrowDown') {
        e.preventDefault();
        setSelectedIndex((i) => (i + 1) % filtered.length);
      }
      if (e.key === 'ArrowUp') {
        e.preventDefault();
        setSelectedIndex((i) => (i - 1 + filtered.length) % filtered.length);
      }
      if (e.key === 'Enter') {
        e.preventDefault();
        filtered[selectedIndex]?.action();
      }
    };
    document.addEventListener('keydown', onKeyDown);
    return () => document.removeEventListener('keydown', onKeyDown);
  }, [open, filtered, selectedIndex, onClose]);

  useEffect(() => {
    const selected = listRef.current?.querySelector(`[data-index="${selectedIndex}"]`) as HTMLElement | null;
    if (selected && typeof selected.scrollIntoView === 'function') {
      selected.scrollIntoView({ block: 'nearest' });
    }
  }, [selectedIndex]);

  if (!open) return null;

  return (
    <div class="command-palette" onClick={(e) => { if (e.target === e.currentTarget) onClose(); }}>
      <div class="command-box">
        <input
          ref={inputRef}
          class="command-input"
          type="text"
          placeholder="Type a command or search..."
          value={query}
          onInput={(e) => setQuery((e.target as HTMLInputElement).value)}
          aria-label="Command palette"
        />
        <div ref={listRef} class="command-list" role="listbox">
          {filtered.length === 0 ? (
            <div class="empty-state p-4">
              <div class="empty-state-title">No commands found</div>
            </div>
          ) : (
            filtered.map((cmd, index) => (
              <button
                key={cmd.id}
                class={`command-item ${index === selectedIndex ? 'active' : ''}`}
                data-index={index}
                role="option"
                aria-selected={index === selectedIndex}
                onClick={cmd.action}
                onMouseEnter={() => setSelectedIndex(index)}
              >
                <Icon name={cmd.icon} size={18} />
                <span>{cmd.label}</span>
                {cmd.shortcut && <span class="command-shortcut">{cmd.shortcut}</span>}
              </button>
            ))
          )}
        </div>
      </div>
    </div>
  );
}
