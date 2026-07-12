> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — CI & Release Audit CI Integration Notes

> **Status:** ⬜ Draft | **Created:** 2026-07-08 | **Last Updated:** 2026-07-09

This document describes the concrete workflow changes implied by the initiative. It is planning-only; no YAML files are modified until execution is approved.

---

## 1. New / modified CI jobs

### `ci.yml` additions

| Job | Trigger | Purpose |
|-----|---------|---------|
| `docker-build` | `push`/`pull_request` to `main`/`dev` | Validates `docker build` and `docker compose up` on every PR. |
| `install-dry-run` | `push`/`pull_request` to `main`/`dev` | Runs `scripts/install.sh --dry-run` for latest and a pinned version across OS/arch matrix. |
| `changelog-check` | `pull_request` | Ensures `CHANGELOG.md` contains a section matching `VERSION` and follows Keep a Changelog format. |

### `release.yml` additions

| Step / Job | Purpose |
|------------|---------|
| `verify-version` job | Fail if `refs/tags/v${VERSION}` does not match pushed tag. |
| `verify-changelog` job | Fail if `CHANGELOG.md` lacks a `## [X.Y.Z]` section for the tag. |
| `extract-release-notes` step | Write the matching changelog section to a temp file for the release body. |
| `build-push-action` digest output | Capture image digest for signing. |
| `cosign-sign` step | Keyless-sign `checksums.txt` and the container digest. |
| `sbom-action` step | Generate image SBOM and attach to release. |
| `slsa-generator` job | Build provenance attestation for release artifacts. |
| `release-smoke` job | Download published tarball and run `scripts/smoke-test.sh`. |
| `release-container-smoke` job | Pull published image and run smoke tests inside it. |

---

## 2. Workflow snippets (illustrative)

### 2.1 Version alignment gate

```yaml
verify-version:
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@<sha> # v4
    - name: Verify VERSION matches tag
      run: |
        VERSION="$(cat VERSION)"
        TAG="${GITHUB_REF_NAME#v}"
        if [[ "$VERSION" != "$TAG" ]]; then
          echo "VERSION file ($VERSION) does not match tag ($TAG)" >&2
          exit 1
        fi
```

### 2.2 Changelog extraction

```yaml
    - name: Extract release notes
      id: notes
      run: |
        VERSION="$(cat VERSION)"
        awk "/^## \\[$VERSION\\]/{flag=1;next}/^## \\[/{flag=0}flag" CHANGELOG.md > release-notes.md
        if [[ ! -s release-notes.md ]]; then
          echo "No changelog section found for $VERSION" >&2
          exit 1
        fi
```

### 2.3 cosign signing

```yaml
    - name: Install cosign
      uses: sigstore/cosign-installer@<sha>

    - name: Sign checksums
      run: cosign sign-blob --yes dist/checksums.txt --output-signature dist/checksums.txt.sig

    - name: Sign image digest
      run: cosign sign --yes "${IMAGE}@${DIGEST}"
```

### 2.4 SLSA provenance

```yaml
  provenance:
    needs: release
    permissions:
      actions: read
      id-token: write
      contents: write
    uses: slsa-framework/slsa-github-generator/.github/workflows/generator_generic_slsa3.yml@v2.0.0
    with:
      base64-subjects: "${{ needs.release.outputs.hashes }}"
      upload-assets: true
```

The `release` job must produce a `hashes` output containing base64-encoded subject hashes of release artifacts.

### 2.5 Docker build CI job

```yaml
  docker-build:
    runs-on: ubuntu-latest
    needs: ui-build
    steps:
      - uses: actions/checkout@<sha> # v4
      - uses: actions/download-artifact@<sha> # v4
        with:
          name: dashboard-dist
          path: internal/ui/dashboard-dist/
      - name: Build image
        run: docker build -t openmuara:ci .
      - name: Smoke test container
        run: |
          docker run -d --name muara-ci -p 127.0.0.1:9000:9000 openmuara:ci
          sleep 3
          curl -fsS http://127.0.0.1:9000/healthz | grep -q '"status":"ok"'
          docker rm -f muara-ci
```

### 2.6 Trivy severity gate

```yaml
    - name: Scan image with Trivy
      uses: aquasecurity/trivy-action@<sha>
      with:
        image-ref: openmuara:scan
        format: sarif
        output: trivy-results.sarif
        severity: 'CRITICAL,HIGH'
        exit-code: '1'
```

Keep SARIF upload with `if: always()` so results are visible even when the gate fails.

### 2.7 GitHub artifact attestation

```yaml
    - name: Attest release artifacts
      uses: actions/attest-build-provenance@<sha>
      with:
        subject-path: dist/*.tar.gz

    - name: Attest container image
      uses: actions/attest-build-provenance@<sha>
      with:
        subject-name: ${{ env.IMAGE }}
        subject-digest: ${{ steps.build-push.outputs.digest }}
        push-to-registry: true
```

### 2.8 Minimal token permissions

```yaml
permissions: {}

jobs:
  release:
    permissions:
      contents: write
      packages: write
      id-token: write
      attestations: write
      security-events: write
    steps: ...
```

### 2.9 Scorecard workflow

```yaml
name: Scorecard supply-chain security
on:
  branch_protection_rule:
  schedule:
    - cron: '25 5 * * 1'
  push:
    branches: [main]

permissions: read-all

jobs:
  analysis:
    name: Scorecard analysis
    runs-on: ubuntu-latest
    permissions:
      security-events: write
      id-token: write
    steps:
      - uses: actions/checkout@<sha> # v4
        with:
          persist-credentials: false
      - uses: ossf/scorecard-action@<sha>
        with:
          results_file: results.sarif
          results_format: sarif
          publish_results: true
      - uses: github/codeql-action/upload-sarif@<sha> # v3
        with:
          sarif_file: results.sarif
```

### 2.10 workflow_dispatch release trigger

```yaml
on:
  workflow_dispatch:
    inputs:
      tag:
        description: 'Tag to release (e.g., v1.1.0)'
        required: true
        type: string
  push:
    tags:
      - 'v*'
```

### 2.11 Release failure notification

```yaml
    - name: Notify on failure
      if: failure()
      uses: actions/github-script@<sha>
      with:
        script: |
          github.rest.issues.create({
            owner: context.repo.owner,
            repo: context.repo.repo,
            title: `Release workflow failed for ${context.ref}`,
            body: `See ${context.payload.repository.html_url}/actions/runs/${context.runId}`,
            labels: ['release', 'incident']
          })
```

---

## 3. Dockerfile changes

### 3.1 Copy prebuilt dashboard

```dockerfile
# Build stage remains unchanged.

# Runtime stage
FROM alpine:3.21
RUN apk add --no-cache ca-certificates
# ... user setup ...
COPY --from=builder /bin/muara /usr/local/bin/muara
COPY --from=builder /src/scripts/docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
COPY --from=builder /src/internal/ui/dashboard-dist/ /app/internal/ui/dashboard-dist/
# Fallback: if dashboard-dist is missing, ensure a placeholder is embedded at build time.
HEALTHCHECK --interval=10s --timeout=5s --start-period=5s --retries=3 \
  CMD muara health
```

### 3.2 Optional fallback build

If the prebuilt dist is not present (local `docker build`), the Dockerfile can optionally build a placeholder:

```dockerfile
ARG DASHBOARD_DIST=internal/ui/dashboard-dist
COPY ${DASHBOARD_DIST} /app/internal/ui/dashboard-dist/ 2>/dev/null || \
  mkdir -p /app/internal/ui/dashboard-dist/ && \
  echo '<html><body>Dashboard not built. Run task ui:build first.</body></html>' \
    > /app/internal/ui/dashboard-dist/index.html
```

> Note: `2>/dev/null` in a Dockerfile `COPY` is not valid syntax; use a build stage or build-arg check instead. The final implementation should use a conditional build stage.

---

## 4. docker-compose.yml changes

```yaml
services:
  muara:
    build:
      context: .
      dockerfile: Dockerfile
    image: openmuara:latest
    container_name: openmuara
    read_only: true
    cap_drop:
      - ALL
    ports:
      - "127.0.0.1:9000:9000"
    volumes:
      - ./.muara:/app/.muara
    environment:
      - MUARA_SERVER_HOST=0.0.0.0
    command: ["start"]
    healthcheck:
      test: ["CMD", "muara", "health"]
      interval: 10s
      timeout: 5s
      retries: 3
      start_period: 5s
```

## 4.1. Distroless image build

```dockerfile
# syntax=docker/dockerfile:1
FROM golang:1.26-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -buildvcs=false -ldflags="-w -s" -o /bin/muara ./cmd/muara

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /bin/muara /usr/local/bin/muara
COPY --from=builder /src/internal/ui/dashboard-dist/ /app/internal/ui/dashboard-dist/
COPY --from=builder /src/scripts/docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
USER nonroot:nonroot
WORKDIR /app
EXPOSE 9000
ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
CMD ["start"]
```

---

## 5. Install script changes

### 5.1 Verification flow

```bash
# After downloading archive and checksums.txt
cd "$TMP_DIR"
sha256sum -c checksums.txt --strict --ignore-missing || {
  echo "Checksum verification failed" >&2
  exit 1
}

if command -v cosign >/dev/null 2>&1 && [[ -f checksums.txt.sig ]]; then
  cosign verify-blob \
    --signature checksums.txt.sig \
    --certificate-identity-regexp 'https://github.com/openmuara/openmuara/.github/workflows/release.yml@refs/tags/.*' \
    --certificate-oidc-issuer https://token.actions.githubusercontent.com \
    checksums.txt
fi
```

### 5.2 Escape hatch

```bash
if [[ "${SKIP_VERIFY:-0}" == "1" ]]; then
  echo "WARNING: skipping checksum/signature verification." >&2
else
  verify_checksums
  verify_signature
fi
```

---

## 6. Taskfile changes

Add tasks for local release validation:

```yaml
  release:verify:
    desc: Verify a downloaded release artifact locally
    cmds:
      - |
        echo "Verify checksums: sha256sum -c checksums.txt --strict --ignore-missing"
        echo "Verify cosign: cosign verify-blob --signature checksums.txt.sig checksums.txt"

  release:smoke-local:
    desc: Run smoke test against the locally built release binary
    cmds:
      - ./scripts/smoke-test.sh
```

---

## 7. Secrets and permissions

No new long-lived secrets are required. cosign uses OIDC. The release workflow already has `contents: write`, `packages: write`, and `security-events: write`. The provenance job additionally needs `id-token: write` and `actions: read`.

---

## 8. Validation commands

```bash
# Local workflow lint
actionlint .github/workflows/release.yml .github/workflows/ci.yml

# Local task quality
task quality

# Local release build
task release:build

# Local container build
task release:docker

# Verify image labels
docker inspect openmuara:latest --format='{{json .Config.Labels}}'
```
