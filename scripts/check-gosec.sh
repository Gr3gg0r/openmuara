#!/usr/bin/env bash
set -euo pipefail

# check-gosec.sh runs gosec if available.
# If gosec is not installed, prints install instructions and exits 0.

if ! command -v gosec >/dev/null 2>&1; then
  echo "gosec is not installed."
  echo "Install it with:"
  echo "  go install github.com/securego/gosec/v2/cmd/gosec@latest"
  echo "Skipping gosec scan."
  exit 0
fi

echo "Running gosec..."
gosec ./...
