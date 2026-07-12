#!/usr/bin/env bash
# Compare per-package coverage between the current checkout (PR) and a base ref.
# Fail if any changed Go package has lower coverage than the base.
set -euo pipefail

BASE_REF="${1:-origin/main}"
MIN_LINES="${2:-10}"
TOLERANCE="${3:-1.0}"

# Parse `go test -cover` output into "<import-path> <coverage>" lines.
# Handles both cached and non-cached output, and ignores packages with no tests.
run_coverage() {
  go test -cover ./... 2>&1 |
    sed -E -n 's/^ok[[:space:]]+([^[:space:]]+)[[:space:]]+[^[:space:]]+[[:space:]]+coverage: ([0-9.]+)%.*/\1 \2/p'
}

# PR coverage.
pr_cov=$(mktemp)
trap 'rm -f "$pr_cov"; [[ -n "${worktree:-}" ]] && git worktree remove -f "$worktree" 2>/dev/null || true' EXIT
run_coverage > "$pr_cov"

# Fetch base ref if not present.
if ! git rev-parse --verify "$BASE_REF" >/dev/null 2>&1; then
  git fetch origin "$BASE_REF"
fi

# Total changed Go lines.
changed_lines=$(git diff --numstat "$BASE_REF" -- '*.go' | awk '{s+=$1+$2} END {print s+0}')
if [[ "$changed_lines" -lt "$MIN_LINES" ]]; then
  echo "Only $changed_lines changed Go lines (<$MIN_LINES); skipping coverage regression check."
  exit 0
fi

# Base coverage in a disposable worktree.
worktree=$(mktemp -d)
git worktree add -q "$worktree" "$BASE_REF"
(
  cd "$worktree"
  go mod download
  run_coverage > coverage.base.txt
)
base_cov="$worktree/coverage.base.txt"

# Map changed Go files to package import paths.
mapfile -t changed_files < <(git diff --name-only "$BASE_REF" -- '*.go')
declare -A changed_packages
for f in "${changed_files[@]}"; do
  dir=$(dirname "$f")
  if import_path=$(cd "$dir" && go list . 2>/dev/null); then
    changed_packages["$import_path"]=1
  fi
done

fail=0
for pkg in "${!changed_packages[@]}"; do
  pr=$(awk -v p="$pkg" '$1 == p {print $2}' "$pr_cov")
  base=$(awk -v p="$pkg" '$1 == p {print $2}' "$base_cov")
  if [[ -z "$pr" || -z "$base" ]]; then
    echo "WARN: no coverage data for $pkg; skipping"
    continue
  fi
  if awk "BEGIN { exit ($pr + $TOLERANCE >= $base) ? 0 : 1 }"; then
    echo "PASS: $pkg coverage $pr% >= base $base% (tolerance ${TOLERANCE}%)"
  else
    echo "FAIL: $pkg coverage dropped from $base% to $pr% (tolerance ${TOLERANCE}%)" >&2
    fail=1
  fi
done

exit "$fail"
