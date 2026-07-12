interface SelectProps {
  value: string;
  onChange: (value: string) => void;
  label?: string;
  children: preact.ComponentChildren;
}

export function Select({ value, onChange, label, children }: SelectProps) {
  return (
    <span class="select-wrapper">
      {label && <span class="select-label sr-only">{label}</span>}
      <select
        class="select"
        value={value}
        onChange={(e) => onChange((e.target as HTMLSelectElement).value)}
        aria-label={label}
      >
        {children}
      </select>
    </span>
  );
}
