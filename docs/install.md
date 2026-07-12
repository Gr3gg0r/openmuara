# Installing OpenMuara

OpenMuara is distributed as pre-built binaries, a container image, and source.

## Recommended: release binary

Use the install script (macOS / Linux):

```bash
curl -sSL https://raw.githubusercontent.com/Gr3gg0r/openmuara/main/scripts/install.sh | bash
```

The script detects your OS and architecture, downloads the correct archive, and
verifies the SHA256 checksum. If [cosign](https://docs.sigstore.dev/cosign/overview/)
is installed, it also verifies the signature of the checksum file.

### Verify the installation

```bash
muara version
muara health
```

## Verifying a downloaded release manually

Every GitHub Release includes:

- `muara-<os>-<arch>.tar.gz` — the compressed binary
- `checksums.txt` — SHA256 hashes
- `checksums.txt.sig` — cosign signature
- `checksums.txt.crt` — cosign certificate
- `sbom*.spdx.json` — SBOMs for Go, dashboard, website, and container
- `openmuara-<version>.intoto.jsonl` — SLSA provenance attestation

### Verify the checksum

```bash
curl -LO https://github.com/Gr3gg0r/openmuara/releases/download/v0.1.1/checksums.txt
curl -LO https://github.com/Gr3gg0r/openmuara/releases/download/v0.1.1/muara-linux-amd64.tar.gz
sha256sum -c checksums.txt --strict --ignore-missing
```

### Verify the signature with cosign

```bash
curl -LO https://github.com/Gr3gg0r/openmuara/releases/download/v0.1.1/checksums.txt.sig
cosign verify-blob \
  --signature checksums.txt.sig \
  --certificate-identity-regexp 'https://github.com/Gr3gg0r/openmuara/.github/workflows/release.yml@refs/tags/.*' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com \
  checksums.txt
```

### Verify SLSA provenance

```bash
slsa-verifier verify-artifact \
  --provenance-path openmuara-v0.1.1.intoto.jsonl \
  --source-uri github.com/Gr3gg0r/openmuara \
  --source-tag v0.1.1 \
  muara-linux-amd64.tar.gz
```

## Container image

A signed multi-arch image is published to GitHub Container Registry on every
release:

```bash
docker run --rm -p 127.0.0.1:9000:9000 \
  -e MUARA_SERVER_HOST=0.0.0.0 \
  -v "$(pwd)/.muara:/app/.muara" \
  ghcr.io/gr3gg0r/openmuara:latest
```

Or use Docker Compose:

```bash
docker compose up
```

### Verify the container image signature

```bash
cosign verify \
  --certificate-identity-regexp 'https://github.com/Gr3gg0r/openmuara/.github/workflows/release.yml@refs/tags/.*' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com \
  ghcr.io/gr3gg0r/openmuara:0.1.1
```

## From source

Clone the repository and build locally:

```bash
git clone https://github.com/Gr3gg0r/openmuara.git
cd openmuara
go build -o bin/muara ./cmd/muara
./bin/muara init
./bin/muara start
```

For a full development build including the dashboard:

```bash
task build
```

## Skip verification (not recommended)

If you are in an air-gapped environment or otherwise cannot verify artifacts, set
`SKIP_VERIFY=1`:

```bash
curl -sSL https://raw.githubusercontent.com/Gr3gg0r/openmuara/main/scripts/install.sh | SKIP_VERIFY=1 bash
```

This bypasses checksum and signature verification and prints a warning.
