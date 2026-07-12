> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Dependency & License Audit CI Integration

> **Created:** 2026-07-09
> **Last Updated:** 2026-07-09
> **Status:** ✅ Implemented

---

This document contains the exact CI changes recommended by the dependency & license audit initiative. Treat these as implementation-ready plans; apply them during execution.

## 1. Dependabot configuration

Update `.github/dependabot.yml` to add npm ecosystems for both npm packages:

```yaml
version: 2

updates:
  - package-ecosystem: gomod
    directory: /
    schedule:
      interval: weekly
    open-pull-requests-limit: 5
    labels:
      - dependencies
      - go

  - package-ecosystem: npm
    directory: /web/dashboard
    schedule:
      interval: weekly
    open-pull-requests-limit: 5
    labels:
      - dependencies
      - npm
      - dashboard

  - package-ecosystem: npm
    directory: /website
    schedule:
      interval: weekly
    open-pull-requests-limit: 5
    labels:
      - dependencies
      - npm
      - website

  - package-ecosystem: github-actions
    directory: /
    schedule:
      interval: weekly
    open-pull-requests-limit: 5
    labels:
      - dependencies
      - github-actions
```

## 2. CI job: dependency & license audit

Add a new job to `.github/workflows/ci.yml` (or create `.github/workflows/dependency-license-audit.yml`):

```yaml
  dependency-license:
    runs-on: ubuntu-latest
    permissions:
      contents: read
    steps:
      - name: Checkout
        uses: actions/checkout@34e114876b0b11c390a56381ad16ebd13914f8d5 # v4

      - name: Set up Go
        uses: actions/setup-go@40f1582b2485089dde7abd97c1529aa768e1baff # v5
        with:
          go-version: '1.26'

      - name: Set up Node
        uses: actions/setup-node@49933ea5288caeca8642d1e84afbd3f7d6820020 # v4
        with:
          node-version: '20'
          cache: 'npm'
          cache-dependency-path: |
            web/dashboard/package-lock.json
            website/package-lock.json

      - name: Verify Go modules
        run: |
          go mod tidy
          go mod verify
          git diff --exit-code go.mod go.sum

      - name: Install go-licenses
        run: go install github.com/google/go-licenses/v2@v2.0.1

      - name: Check Go licenses
        run: go-licenses check ./...

      - name: Audit dashboard production dependencies
        run: cd web/dashboard && npm ci && npm audit --omit=dev

      - name: Audit website production dependencies (allowed to fail)
        run: cd website && npm ci && npm audit --omit=dev
        continue-on-error: true

      - name: Generate license matrix
        run: |
          go-licenses csv ./... > /tmp/licenses-go.csv
          echo "TODO: merge npm licenses into docs/initiatives/openmuara-readiness-dependency-license-audit/LICENSE_MATRIX.md"
```

## 3. Release workflow: npm SBOMs

Update `.github/workflows/release.yml` to generate and attach npm SBOMs after the Go SBOM step:

```yaml
      - name: Generate Go SBOM
        run: |
          go install github.com/anchore/syft/cmd/syft@v1.46.0
          syft dir:. -o spdx-json=dist/sbom.spdx.json \
            --source-name openmuara \
            --source-version "${GITHUB_REF_NAME}"

      - name: Generate npm SBOMs
        run: |
          cd web/dashboard && npm ci && npm sbom --package-lock-only --sbom-format=spdx > ../../dist/sbom-dashboard.spdx.json
          cd ../../website && npm ci && npm sbom --package-lock-only --sbom-format=spdx > ../dist/sbom-website.spdx.json

      - name: Create GitHub Release
        uses: softprops/action-gh-release@3bb12739c298aeb8a4eeaf626c5b8d85266b0e65 # v2
        with:
          files: |
            dist/*.tar.gz
            dist/checksums.txt
            dist/sbom.spdx.json
            dist/sbom-dashboard.spdx.json
            dist/sbom-website.spdx.json
          body_path: CHANGELOG.md
          fail_on_unmatched_paths: true
```

## 4. Container image scanning in release

Add an image-scan step to the release workflow before pushing the image:

```yaml
      - name: Build image for scanning
        run: docker build -t openmuara:scan .

      - name: Scan image with Trivy
        uses: aquasecurity/trivy-action@6c175e9c4083fbb36b69a55b791cd3bce129a53e # v0.30.0
        with:
          image-ref: openmuara:scan
          format: sarif
          output: trivy-results.sarif

      - name: Upload Trivy SARIF
        uses: github/codeql-action/upload-sarif@02c5e83432fe5497fd85b873b6c9f16a8578e1d9 # v3
        if: always()
        with:
          sarif_file: trivy-results.sarif
```

## 5. Dockerfile base-image pinning (optional)

For reproducible builds, pin base images by digest and let Dependabot update the digests:

```dockerfile
# Build stage
FROM golang:1.26-alpine@sha256:0178a641fbb4858c5f1b48e34bdaabe0350a330a1b1149aabd498d0699ff5fb2 AS builder

# Runtime stage
FROM alpine:3.21@sha256:48b0309ca019d89d40f670aa1bc06e426dc0931948452e8491e3d65087abc07d
```

*Note: digests above are examples; use the actual digests at execution time.*

## 6. Acceptance scripts per phase

### P01 — Go dependency review

```bash
#!/bin/bash
set -euo pipefail
go mod tidy
go mod verify
git diff --exit-code go.mod go.sum
govulncheck ./...
```

### P02 — npm dependency review

```bash
#!/bin/bash
set -euo pipefail
for dir in web/dashboard website; do
  pushd "$dir"
  npm ci
  npm audit --production
  npm outdated || true
  npx depcheck || true
  popd
done
```

### P03 — License matrix generation

```bash
#!/bin/bash
set -euo pipefail
go install github.com/google/go-licenses/v2@v2.0.1
go-licenses csv ./... > /tmp/licenses-go.csv
# TODO: append npm license data and write to LICENSE_MATRIX.md
```

## 7. Rollback / exception handling

- If `go-licenses` flags a false positive, add the package to an allowlist file (e.g., `.github/license-allowlist.csv`) and record the decision in `DECISIONS.md`.
- If website `npm audit` fails due to Docusaurus build-time deps, keep `continue-on-error: true` until upstream fixes are available, and document in `KNOWN_ISSUES.md`.
