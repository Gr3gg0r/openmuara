#!/usr/bin/env bash
# Compare deferred bug IDs between the bug-hunt known issues file and the root
# KNOWN_ISSUES.md. Report drift so the user-facing register stays honest.
set -euo pipefail

BUG_HUNT_FILE="docs/initiatives/openmuara-bug-hunt/KNOWN_ISSUES.md"
ROOT_FILE="KNOWN_ISSUES.md"

# Extract BXXX bug IDs from a markdown file.
extract_ids() {
  local file="$1"
  if [[ ! -f "$file" ]]; then
    return
  fi
  grep -E '^###[[:space:]]+B[0-9]+' "$file" 2>/dev/null | sed -E 's/^###[[:space:]]+(B[0-9]+).*/\1/' | sort -u || true
}

bug_hunt_ids=$(extract_ids "$BUG_HUNT_FILE")
root_ids=$(extract_ids "$ROOT_FILE")

# Allow an explicit marker that means "intentionally not listed in root".
has_intentional_marker() {
  local file="$1"
  [[ -f "$file" ]] && grep -qE '^<!--[[:space:]]*check-known-issues:ignore[[:space:]]*-->$' "$file"
}

missing_in_root=$(comm -23 <(echo "${bug_hunt_ids}") <(echo "${root_ids}"))
missing_in_bug_hunt=$(comm -13 <(echo "${bug_hunt_ids}") <(echo "${root_ids}"))

fail=0

if [[ -n "$missing_in_root" ]]; then
  echo "WARN: deferred bugs missing from root KNOWN_ISSUES.md:"
  while IFS= read -r id; do
    printf '  %s\n' "$id"
  done <<< "$missing_in_root"
  fail=1
fi

if [[ -n "$missing_in_bug_hunt" ]]; then
  if has_intentional_marker "$ROOT_FILE"; then
    echo "INFO: root KNOWN_ISSUES.md has entries not in bug-hunt file, but intentional-ignore marker is present."
  else
    echo "WARN: root KNOWN_ISSUES.md entries missing from bug-hunt file:"
    while IFS= read -r id; do
      printf '  %s\n' "$id"
    done <<< "$missing_in_bug_hunt"
    fail=1
  fi
fi

if [[ "$fail" -eq 0 ]]; then
  echo "PASS: KNOWN_ISSUES.md and $BUG_HUNT_FILE are in sync."
  exit 0
fi

echo "See docs/bug-hunt-process.md for the sync process."
exit 1
