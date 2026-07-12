#!/usr/bin/env bash
set -euo pipefail

MAX_FILE_LINES=250
MAX_FUNC_LINES=80
MAX_LINE_LEN=120

echo "=== Checking file line counts (advisory) ==="
while IFS= read -r file; do
  lines=$(wc -l < "$file")
  if [ "$lines" -gt "$MAX_FILE_LINES" ]; then
    echo "WARN: $file has $lines lines (recommended max $MAX_FILE_LINES)"
  fi
done < <(find . -type f -name '*.go' -not -path './vendor/*')

echo "=== Checking function line counts (advisory) ==="
while IFS= read -r file; do
  awk '
    /^func / { start=NR; name=$0 }
    start && /^}/ {
      len = NR - start + 1
      if (len > '"$MAX_FUNC_LINES"') {
        printf "WARN: %s:%d function spans %d lines (recommended max %d)\n", FILENAME, start, len, '"$MAX_FUNC_LINES"'
      }
      start=0
    }
  ' "$file" || true
done < <(find . -type f -name '*.go' -not -path './vendor/*')

echo "=== Checking line lengths (advisory) ==="
while IFS= read -r file; do
  awk 'length($0) > '"$MAX_LINE_LEN"' { printf "WARN: %s:%d line length %d (recommended max %d)\n", FILENAME, NR, length($0), '"$MAX_LINE_LEN"' }' "$file" || true
done < <(find . -type f -name '*.go' -not -path './vendor/*')

echo "OK: size checks completed (warnings are advisory)"
