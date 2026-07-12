#!/usr/bin/env bash
# Run mutation testing on the curated package list and enforce a minimum score.
set -euo pipefail

THRESHOLD="${1:-70}"
PACKAGES=(
  ./internal/webhook
  ./internal/engine
  # ./internal/fawry is excluded for now: gremlins times out on the HTTP
  # handler/signature tests because mutations cause the test servers to hang.
  # Re-evaluate after adding faster unit tests for the pure functions.
)

failures=0

for pkg in "${PACKAGES[@]}"; do
  echo "==> Running mutation tests for $pkg"
  output=$(gremlins unleash --tags=test "$pkg" 2>&1 || true)
  echo "$output"
  efficacy=$(echo "$output" | grep 'Test efficacy:' | awk '{print $3}' | sed 's/%//')
  if [[ -z "$efficacy" ]]; then
    echo "ERROR: could not determine mutation efficacy for $pkg" >&2
    failures=$((failures + 1))
    continue
  fi
  if awk "BEGIN { exit ($efficacy >= $THRESHOLD) ? 0 : 1 }"; then
    echo "PASS: $pkg mutation efficacy $efficacy% >= $THRESHOLD%"
  else
    echo "FAIL: $pkg mutation efficacy $efficacy% < $THRESHOLD%" >&2
    failures=$((failures + 1))
  fi
done

if [[ "$failures" -gt 0 ]]; then
  exit 1
fi
