#!/bin/bash
# License compliance check for OpenMuara.
# Verifies that all Go dependencies use licenses compatible with MIT distribution.
# Modernc.org/mathutil is ignored because go-licenses cannot detect its BSD-3-Clause
# license text automatically; it was manually verified (see LICENSE_MATRIX.md).

set -euo pipefail

GO_LICENSES_VERSION="v2.0.1"

if ! command -v go-licenses >/dev/null 2>&1; then
  echo "Installing go-licenses ${GO_LICENSES_VERSION}..."
  go install "github.com/google/go-licenses/v2@${GO_LICENSES_VERSION}"
fi

allowed="MIT,Apache-2.0,BSD-2-Clause,BSD-3-Clause,ISC,MPL-2.0"

echo "Checking Go dependency licenses..."
go-licenses check \
  --allowed_licenses="$allowed" \
  --ignore="modernc.org/mathutil" \
  ./...

echo "All Go dependency licenses are compatible."
