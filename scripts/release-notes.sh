#!/usr/bin/env bash
# Generate a release-notes snippet from fixed bugs in the bug hunt TRACKING.md.
# Usage: scripts/release-notes.sh [path-to-TRACKING.md]

set -euo pipefail

TRACKER="${1:-docs/initiatives/openmuara-bug-hunt/TRACKING.md}"

if [[ ! -f "$TRACKER" ]]; then
  echo "Tracker not found: $TRACKER" >&2
  exit 1
fi

echo "## Bug Fixes"
echo ""

# Extract fixed bug rows from the bug register table.
# The table columns are:
# ID | Severity | Area | Summary | Reproduction | Finding File | Root Cause Category | Regression Test | Status | Commit | Introduced By | Fixed By
awk -F'|' '
  /^\|[[:space:]]*B[0-9]+/ {
    gsub(/^[[:space:]]+|[[:space:]]+$/, "", $1) # ID
    gsub(/^[[:space:]]+|[[:space:]]+$/, "", $2) # Severity
    gsub(/^[[:space:]]+|[[:space:]]+$/, "", $3) # Area
    gsub(/^[[:space:]]+|[[:space:]]+$/, "", $4) # Summary
    gsub(/^[[:space:]]+|[[:space:]]+$/, "", $10) # Status
    gsub(/^[[:space:]]+|[[:space:]]+$/, "", $11) # Commit
    id = $1; severity = $2; area = $3; summary = $4; status = $10; commit = $11
    if (status ~ /fixed|Fixed|✅/) {
      printf "- **%s** (%s / %s): %s", id, severity, area, summary
      if (commit && commit != "—" && commit != "-") {
        printf " (%s)", commit
      }
      printf "\n"
    }
  }
' "$TRACKER"
