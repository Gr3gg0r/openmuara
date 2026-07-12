#!/usr/bin/env bash
set -euo pipefail

fail=0

echo "=== Checking for forbidden patterns in non-test Go files ==="

# Debug prints in production code (examples/ are standalone programs).
if grep -R "fmt\.Println" --include='*.go' --exclude='*_test.go' --exclude-dir=examples .; then
  echo "FAIL: fmt.Println found in production code"
  fail=1
fi

# os.Exit outside cmd/ and examples/.
while IFS= read -r file; do
  if grep -q "os\.Exit" "$file"; then
    echo "FAIL: os.Exit in library code: $file"
    fail=1
  fi
done < <(find . -type f -name '*.go' -not -path './cmd/*' -not -path './examples/*' -not -path './vendor/*')

if [ "$fail" -eq 1 ]; then
  exit 1
fi

echo "OK: no forbidden patterns found"
