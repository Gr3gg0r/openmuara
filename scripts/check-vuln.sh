#!/bin/bash
set -euo pipefail

# check-vuln.sh runs govulncheck if available.
# If govulncheck is not installed, prints install instructions and exits 0.

if ! command -v govulncheck >/dev/null 2>&1; then
    echo "govulncheck is not installed."
    echo "Install it with:"
    echo "  go install golang.org/x/vuln/cmd/govulncheck@latest"
    echo "Skipping vulnerability scan."
    exit 0
fi

echo "Running govulncheck..."
govulncheck ./...
