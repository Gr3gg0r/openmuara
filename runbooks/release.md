# Release Runbook

This runbook describes how to cut a new OpenMuara release.

## Prerequisites

- Maintainer access to the GitHub repository.
- `git`, `gh` (optional), and `cosign` (for verification) installed locally.
- A clean `dev` branch with all intended changes merged.

## Pre-release checks

1. Ensure CI is green on `dev`.

   ```bash
   git checkout dev
   git pull origin dev
   task quality
   ```

2. Decide the next version using [Semantic Versioning](https://semver.org/).
   Examples: `1.1.0`, `1.0.1`, `1.1.0-rc.1`.

3. Update `VERSION`:

   ```bash
   echo "1.1.0" > VERSION
   ```

4. Update `CHANGELOG.md`. Move items from `[Unreleased]` to a new section:

   ```markdown
   ## [1.1.0] - 2026-07-10
   ```

5. Open a PR from a release-prep branch to `dev` with these two changes.
   The PR must pass the `changelog-check` and all other CI jobs.

6. Merge the PR into `dev`.

## Create the release

1. Tag the merge commit on `dev`:

   ```bash
   git checkout dev
   git pull origin dev
   git tag -a v1.1.0 -m "Release v1.1.0"
   git push origin v1.1.0
   ```

2. The `Release` workflow starts automatically. Monitor it on the
   **Actions** tab.

3. Wait for all jobs to complete:

   - `verify-version` and `verify-changelog`
   - `release` (build, sign, scan, push, create GitHub Release)
   - `provenance` (SLSA attestation)
   - `release-smoke` (binary smoke test)
   - `release-container-smoke` (container smoke test)

## Post-release validation

1. Open the GitHub Release page and verify the assets:

   - All platform tarballs
   - `checksums.txt`, `checksums.txt.sig`, `checksums.txt.crt`
   - `sbom*.spdx.json`
   - `openmuara-v1.1.0.intoto.jsonl`

2. Verify the checksum and signature:

   ```bash
   VERSION=1.1.0
   curl -LO "https://github.com/Gr3gg0r/openmuara/releases/download/v${VERSION}/checksums.txt"
   curl -LO "https://github.com/Gr3gg0r/openmuara/releases/download/v${VERSION}/checksums.txt.sig"
   curl -LO "https://github.com/Gr3gg0r/openmuara/releases/download/v${VERSION}/muara-linux-amd64.tar.gz"
   sha256sum -c checksums.txt --strict --ignore-missing
   cosign verify-blob \
     --signature checksums.txt.sig \
     --certificate-identity-regexp 'https://github.com/Gr3gg0r/openmuara/.github/workflows/release.yml@refs/tags/.*' \
     --certificate-oidc-issuer https://token.actions.githubusercontent.com \
     checksums.txt
   ```

3. Verify the container image signature:

   ```bash
   cosign verify \
     --certificate-identity-regexp 'https://github.com/Gr3gg0r/openmuara/.github/workflows/release.yml@refs/tags/.*' \
     --certificate-oidc-issuer https://token.actions.githubusercontent.com \
     "ghcr.io/gr3gg0r/openmuara:v${VERSION}"
   ```

4. Verify SLSA provenance:

   ```bash
   slsa-verifier verify-artifact \
     --provenance-path "openmuara-v${VERSION}.intoto.jsonl" \
     --source-uri github.com/Gr3gg0r/openmuara \
     --source-tag "v${VERSION}" \
     muara-linux-amd64.tar.gz
   ```

5. Smoke-test the release binary:

   ```bash
   tar -xzf muara-linux-amd64.tar.gz
   ./muara version
   ./muara init
   ./muara start &
   sleep 2
   ./muara health
   ```

## Prereleases

For release candidates or betas, use a semver prerelease tag:

```bash
git tag -a v1.1.0-rc.1 -m "Release v1.1.0-rc.1"
git push origin v1.1.0-rc.1
```

The workflow will:

- Create a GitHub **pre-release**.
- Push `ghcr.io/gr3gg0r/openmuara:1.1.0-rc.1`.
- **Not** move the `latest` container tag.

## Rollback

If a release is broken:

1. Delete the GitHub Release and the git tag (requires admin).

   ```bash
   gh release delete v1.1.0 --yes
   git push --delete origin v1.1.0
   git tag --delete v1.1.0
   ```

2. If the container `latest` tag was moved, retag it to the previous digest:

   ```bash
   docker pull ghcr.io/gr3gg0r/openmuara:<previous-digest>
   docker tag ghcr.io/gr3gg0r/openmuara:<previous-digest> ghcr.io/gr3gg0r/openmuara:latest
   docker push ghcr.io/gr3gg0r/openmuara:latest
   ```

3. Document the regression in `CHANGELOG.md` under `[Unreleased]`.

4. If needed, re-tag a known-good commit and push it to trigger a fresh release.

## Branch and worktree cleanup

After a release or when a feature initiative is fully merged, clean up stale
branches and worktrees to keep the repository tidy:

```bash
# List worktrees
git worktree list

# Remove a merged worktree (safe; does not delete the branch)
git worktree remove /path/to/worktree

# Delete a stale branch after its worktree is removed
git branch -d feat/old-branch

# Delete a suspended branch that was never merged (use with care)
git branch -D feat/suspended-branch
```

## `main` synchronization

`main` is updated from `dev` at release time. For a public transfer or major
milestone, fast-forward `main` to the current `dev` head:

```bash
git checkout main
git merge --ff-only dev
git push origin main
```

## Troubleshooting

| Symptom | Likely cause | Fix |
|---------|--------------|-----|
| `verify-version` fails | `VERSION` does not match tag | Update `VERSION` and retag |
| `verify-changelog` fails | Missing changelog section | Add `## [X.Y.Z]` section |
| cosign signing fails | Missing `id-token: write` or wrong identity regex | Check workflow permissions and regex |
| Trivy fails the build | CRITICAL/HIGH CVE in base image | Pin base image digest or patch dependencies |
| `release-smoke` fails | Binary incompatible with smoke test | Check `MUARA_BINARY` handling in `scripts/smoke-test.sh` |
