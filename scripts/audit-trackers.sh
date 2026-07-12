#!/usr/bin/env bash
set -euo pipefail

# audit-trackers.sh — lightweight consistency checks for OpenMuara planning docs.
# Run from the repo root.

errors=0

echo "=== Checking for absolute worktree / env paths in docs ==="
if grep -R -E "(/Volumes/|/home/agent/|/Users/[^/]+/Workspace/)" docs/ 2>/dev/null; then
  echo "ERROR: absolute worktree or environment paths found in docs/"
  errors=$((errors + 1))
else
  echo "OK: no absolute worktree or environment paths found"
fi

echo ""
echo "=== Checking master backlog entry points exist ==="
backlog="docs/initiatives/openmuara-v1-master-backlog/TRACKING.md"
# shellcheck disable=SC2016
entry_points=$(grep -E '^\|' "$backlog" | grep -oE '`[^`]+\.(md|yml)`' | tr -d '`' | sort -u)
missing=0
for ep in $entry_points; do
  target="${ep/<repo-root>/.}"
  if [ ! -f "$target" ]; then
    echo "MISSING: $target (from entry point: $ep)"
    missing=$((missing + 1))
  fi
done
if [ "$missing" -gt 0 ]; then
  echo "ERROR: $missing entry point(s) missing"
  errors=$((errors + 1))
else
  echo "OK: all entry points exist"
fi

echo ""
echo "=== Checking referenced commit hashes exist ==="
# Extract short commit hashes from tracker tables and verify them.
commits=$(grep -hoE '[a-f0-9]{7,40}' TRACKING.md docs/projects/openmuara-v1/TRACKING.md 2>/dev/null | sort -u)
bad_commits=0
for c in $commits; do
  if ! git rev-parse --verify --quiet "$c" >/dev/null; then
    echo "BAD COMMIT: $c"
    bad_commits=$((bad_commits + 1))
  fi
done
if [ "$bad_commits" -gt 0 ]; then
  echo "ERROR: $bad_commits referenced commit(s) not found in git"
  errors=$((errors + 1))
else
  echo "OK: all referenced commits exist"
fi

echo ""
echo "=== Checking SenangPay gateway YAML exists ==="
if [ ! -f "plugins/senangpay/gateway.yml" ]; then
  echo "ERROR: plugins/senangpay/gateway.yml is missing"
  errors=$((errors + 1))
else
  echo "OK: plugins/senangpay/gateway.yml exists"
fi

echo ""
echo "=== Checking example config exists and matches bundled defaults ==="
if [ ! -f "muara.yml.example" ]; then
  echo "ERROR: muara.yml.example is missing"
  errors=$((errors + 1))
else
  repo_root=$(git rev-parse --show-toplevel)
  check_file="$repo_root/tmp_check_default_yaml.go"
  default_file="$repo_root/tmp_default.yaml"
  trap 'rm -f "$check_file" "$default_file"' EXIT
  cat > "$check_file" <<'GOEOF'
package main

import (
	"fmt"
	"os"

	"github.com/openmuara/openmuara/internal/config"
)

func main() {
	if err := os.WriteFile("tmp_default.yaml", config.DefaultYAML(), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
GOEOF
  (cd "$repo_root" && go run tmp_check_default_yaml.go)
  if diff -q "$default_file" muara.yml.example >/dev/null; then
    echo "OK: muara.yml.example matches internal/config.DefaultYAML()"
  else
    echo "ERROR: muara.yml.example drifts from internal/config.DefaultYAML()"
    diff "$default_file" muara.yml.example || true
    errors=$((errors + 1))
  fi
fi

echo ""
if [ "$errors" -eq 0 ]; then
  echo "All tracker checks passed."
  exit 0
else
  echo "Tracker audit failed with $errors error(s)."
  exit 1
fi
