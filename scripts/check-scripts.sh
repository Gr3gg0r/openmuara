#!/usr/bin/env bash
set -euo pipefail

if ! command -v shellcheck >/dev/null 2>&1; then
  echo "SKIP: shellcheck not installed"
  exit 0
fi

echo "=== Running shellcheck on scripts ==="
find ./scripts -type f -name '*.sh' -print0 | xargs -0 shellcheck
