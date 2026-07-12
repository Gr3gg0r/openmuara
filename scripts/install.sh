#!/usr/bin/env bash
# Install the latest (or a specific) OpenMuara release binary.
# Usage:
#   curl -sSL https://raw.githubusercontent.com/openmuara/openmuara/main/scripts/install.sh | bash
#   curl -sSL ... | bash -s -- -p ~/.local/bin
#   VERSION=v1.0.0 curl -sSL ... | bash
#
# Verification:
#   By default the installer verifies the SHA256 checksum of the downloaded
#   archive. If cosign is installed, it also verifies the signature of the
#   checksum file. Set SKIP_VERIFY=1 to skip verification (not recommended).

set -euo pipefail

REPO="openmuara/openmuara"
INSTALL_PREFIX="${INSTALL_PREFIX:-/usr/local/bin}"
DRY_RUN=0
SKIP_VERIFY="${SKIP_VERIFY:-0}"

print_usage() {
  cat <<EOF
Usage: install.sh [options]

Options:
  -p, --prefix <dir>   Installation directory (default: /usr/local/bin)
  -v, --version <tag>  Release tag to install (default: latest)
  -d, --dry-run        Print the download URL and target path, but do not install
  -h, --help           Show this help message

Environment:
  INSTALL_PREFIX       Installation directory
  VERSION              Release tag to install
  SKIP_VERIFY=1        Skip checksum/signature verification (not recommended)
EOF
}

VERSION="${VERSION:-latest}"

while [[ $# -gt 0 ]]; do
  case "$1" in
    -p|--prefix)
      INSTALL_PREFIX="$2"
      shift 2
      ;;
    -v|--version)
      VERSION="$2"
      shift 2
      ;;
    -d|--dry-run)
      DRY_RUN=1
      shift
      ;;
    -h|--help)
      print_usage
      exit 0
      ;;
    *)
      echo "Unknown option: $1" >&2
      print_usage >&2
      exit 1
      ;;
  esac
done

# Normalize repo to lowercase to avoid case-sensitivity issues.
REPO="$(echo "$REPO" | tr '[:upper:]' '[:lower:]')"

# Normalize version into a download path.
if [[ "$VERSION" == "latest" ]]; then
  DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download"
else
  DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}"
fi

# Detect OS.
case "$(uname -s)" in
  Linux*)     OS=linux;;
  Darwin*)    OS=darwin;;
  CYGWIN*|MINGW*|MSYS*) OS=windows;;
  *)
    echo "Unsupported operating system: $(uname -s)" >&2
    exit 1
    ;;
esac

# Detect architecture.
case "$(uname -m)" in
  x86_64|amd64) ARCH=amd64;;
  arm64|aarch64) ARCH=arm64;;
  *)
    echo "Unsupported architecture: $(uname -m)" >&2
    exit 1
    ;;
esac

if [[ "$OS" == "windows" ]]; then
  BIN_NAME="muara.exe"
  ARCHIVE="muara-${OS}-${ARCH}.tar.gz"
else
  BIN_NAME="muara"
  ARCHIVE="muara-${OS}-${ARCH}.tar.gz"
fi

URL="${DOWNLOAD_URL}/${ARCHIVE}"
CHECKSUMS_URL="${DOWNLOAD_URL}/checksums.txt"
SIGNATURE_URL="${DOWNLOAD_URL}/checksums.txt.sig"
TARGET="${INSTALL_PREFIX}/${BIN_NAME}"

echo "OpenMuara installer"
echo "  OS:        ${OS}"
echo "  Arch:      ${ARCH}"
echo "  Version:   ${VERSION}"
echo "  Archive:   ${URL}"
echo "  Target:    ${TARGET}"

if [[ "$DRY_RUN" -eq 1 ]]; then
  exit 0
fi

if [[ ! -d "$INSTALL_PREFIX" ]]; then
  echo "Installation directory does not exist: ${INSTALL_PREFIX}" >&2
  echo "Create it first or set INSTALL_PREFIX to an existing directory." >&2
  exit 1
fi

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

# Prefer sha256sum; fall back to shasum -a 256 on macOS.
sha256_cmd() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$@"
  else
    shasum -a 256 "$@"
  fi
}

verify_checksums() {
  echo "Verifying checksum..."
  cd "$TMP_DIR"
  local expected actual
  expected=$(awk -v file="${ARCHIVE}" '$2 == file {print $1}' checksums.txt)
  if [ -z "$expected" ]; then
    echo "No checksum found for ${ARCHIVE}" >&2
    exit 1
  fi
  actual=$(sha256_cmd "${ARCHIVE}" | awk '{print $1}')
  if [ "$expected" != "$actual" ]; then
    echo "Checksum verification failed for ${ARCHIVE}" >&2
    exit 1
  fi
  cd - >/dev/null
  echo "Checksum verified."
}

verify_signature() {
  if ! command -v cosign >/dev/null 2>&1; then
    echo "cosign not found; skipping signature verification."
    return 0
  fi
  if [[ ! -f "${TMP_DIR}/checksums.txt.sig" ]]; then
    echo "No signature file found; skipping signature verification."
    return 0
  fi
  echo "Verifying signature with cosign..."
  if cosign verify-blob \
       --signature "${TMP_DIR}/checksums.txt.sig" \
       --certificate-identity-regexp "https://github.com/${REPO}/.github/workflows/release.yml@refs/tags/.*" \
       --certificate-oidc-issuer https://token.actions.githubusercontent.com \
       "${TMP_DIR}/checksums.txt" >/dev/null 2>&1; then
    echo "Signature verified."
  else
    echo "Signature verification failed." >&2
    exit 1
  fi
}

echo "Downloading ${ARCHIVE}..."
if ! curl -fsSL -o "${TMP_DIR}/${ARCHIVE}" "$URL"; then
  echo "Download failed: ${URL}" >&2
  echo "Check that the release exists and supports ${OS}/${ARCH}." >&2
  exit 1
fi

echo "Downloading checksums..."
if ! curl -fsSL -o "${TMP_DIR}/checksums.txt" "$CHECKSUMS_URL"; then
  echo "Download failed: ${CHECKSUMS_URL}" >&2
  echo "Check that the release includes checksums.txt." >&2
  exit 1
fi

# Signature download is optional; failure is not fatal.
if ! curl -fsSL -o "${TMP_DIR}/checksums.txt.sig" "$SIGNATURE_URL" 2>/dev/null; then
  echo "Signature file not available; continuing without signature verification."
fi

if [[ "$SKIP_VERIFY" == "1" ]]; then
  echo "WARNING: SKIP_VERIFY is set; skipping checksum and signature verification." >&2
else
  verify_checksums
  verify_signature
fi

echo "Extracting..."
tar -xzf "${TMP_DIR}/${ARCHIVE}" -C "$TMP_DIR"

if [[ ! -f "${TMP_DIR}/${BIN_NAME}" ]]; then
  echo "Archive did not contain expected binary: ${BIN_NAME}" >&2
  exit 1
fi

if [[ ! -w "$INSTALL_PREFIX" ]]; then
  echo "Cannot write to ${INSTALL_PREFIX}. Re-run with a writable prefix or use sudo." >&2
  exit 1
fi

chmod +x "${TMP_DIR}/${BIN_NAME}"
mv "${TMP_DIR}/${BIN_NAME}" "$TARGET"

echo "Installed ${TARGET}"
echo "Run 'muara init && muara start' to get started."
