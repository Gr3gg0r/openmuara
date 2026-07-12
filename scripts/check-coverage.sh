#!/bin/bash
set -euo pipefail

# check-coverage.sh runs Go tests with coverage and enforces a minimum threshold.
# Usage: scripts/check-coverage.sh [threshold]
# Default threshold: 50

threshold="${1:-50}"
profile="coverage.out"

run_tests() {
    # Exclude vendored JS dependency Go packages from coverage calculations.
    local packages=()
    while IFS= read -r pkg; do
        packages+=("${pkg}")
    done < <(go list ./... | grep -v -E 'web/dashboard/node_modules|website/node_modules')
    go test -race -coverprofile="${profile}" "${packages[@]}"
}

extract_total() {
    go tool cover -func="${profile}" | awk '/^total:/ {print $3}' | tr -d '%'
}

main() {
    if ! [[ "${threshold}" =~ ^[0-9]+(\.[0-9]+)?$ ]]; then
        echo "Invalid threshold: ${threshold}" >&2
        exit 1
    fi

    run_tests

    local total
    total=$(extract_total)

    echo "Total coverage: ${total}% (threshold: ${threshold}%)"

    if awk "BEGIN {exit !(${total} >= ${threshold})}"; then
        echo "Coverage gate passed."
        exit 0
    fi

    echo "Coverage gate failed: ${total}% < ${threshold}%" >&2
    exit 1
}

main "$@"
