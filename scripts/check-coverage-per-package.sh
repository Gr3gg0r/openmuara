#!/bin/bash
set -euo pipefail

# Minimum coverage per package. Format: "import-path threshold"
# Overall total is still enforced by scripts/check-coverage.sh.
# Floors are calibrated to current architecture; see coverage-exemptions.yml.
MIN_FLOORS=(
  "github.com/Gr3gg0r/openmuara/internal/audit 80"
  "github.com/Gr3gg0r/openmuara/internal/plugin 80"
  "github.com/Gr3gg0r/openmuara/internal/provider/conform 79"
  "github.com/Gr3gg0r/openmuara/internal/version 70"
  "github.com/Gr3gg0r/openmuara/internal/ui 70"
  "github.com/Gr3gg0r/openmuara/internal/provider/simple 45"
)

fail=0
for entry in "${MIN_FLOORS[@]}"; do
  pkg="${entry% *}"
  floor="${entry#* }"
  cov=$(go test -cover "$pkg" 2>&1 |
    sed -E -n 's/^ok[[:space:]]+[^[:space:]]+[[:space:]]+[^[:space:]]+[[:space:]]+coverage: ([0-9.]+)%.*/\1/p')
  if [[ -z "$cov" ]]; then
    echo "WARN: no coverage data for $pkg"
    continue
  fi
  if awk "BEGIN { exit ($cov >= $floor) ? 0 : 1 }"; then
    echo "PASS: $pkg $cov% >= $floor%"
  else
    echo "FAIL: $pkg $cov% < $floor%" >&2
    fail=1
  fi
done

exit "$fail"
