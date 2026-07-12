import { Icon } from './Icon';
import type { IconName } from './Icon';

interface InputProps {
  id?: string;
  type?: 'text' | 'search' | 'url' | 'password';
  value: string;
  placeholder?: string;
  icon?: IconName;
  clearable?: boolean;
  label?: string;
  onInput: (value: string) => void;
  onClear?: () => void;
}

export function Input({
  id,
  type = 'text',
  value,
  placeholder,
  icon,
  clearable,
  label,
  onInput,
  onClear,
}: InputProps) {
  return (
    <span class="input-wrapper">
      {label && (
        <label class="label sr-only" htmlFor={id}>
          {label}
        </label>
      )}
      {icon && <Icon name={icon} size={16} className="input-icon" />}
      <input
        id={id}
        type={type}
        class={`input ${icon ? 'input-with-icon' : ''}`}
        value={value}
        placeholder={placeholder}
        aria-label={label}
        onInput={(e) => onInput((e.target as HTMLInputElement).value)}
      />
      {clearable && value && (
        <button
          type="button"
          class="input-clear"
          onClick={onClear}
          aria-label="Clear"
          title="Clear"
        >
          <Icon name="close" size={14} />
        </button>
      )}
    </span>
  );
}
