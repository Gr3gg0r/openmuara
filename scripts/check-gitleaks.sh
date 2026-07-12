#!/usr/bin/env bash
set -euo pipefail

# check-gitleaks.sh runs gitleaks if available.
# If gitleaks is not installed, prints install instructions and exits 0.

if ! command -v gitleaks >/dev/null 2>&1; then
  echo "gitleaks is not installed."
  echo "Install it with:"
  echo "  brew install gitleaks"
  echo "  or download a release from https://github.com/gitleaks/gitleaks/releases"
  echo "Skipping secret scan."
  exit 0
fi

echo "Running gitleaks..."
gitleaks detect --source .
