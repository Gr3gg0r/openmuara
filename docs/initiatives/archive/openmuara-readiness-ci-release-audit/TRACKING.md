# OpenMuara Readiness — CI & Release Audit Tracking

> **Status:** ✅ Delivered on `dev`  
> **Branch:** `feat/readiness-ci-release-audit` → merged into `dev`

---

## Milestones

| ID | Milestone | Status | Commit / Notes |
|---|---|---|---|
| M1 | Harden release workflow with SLSA, cosign, SBOMs, Trivy, VEX, GitHub attestations, Scorecard, post-release smoke | ✅ | `.github/workflows/release.yml`, `.github/workflows/scorecard.yml`, `vex.json` |
| M2 | Sign and harden Docker image: non-root user, HEALTHCHECK, embedded dashboard, plugins copy, distroless variant | ✅ | `Dockerfile`, `Dockerfile.distroless`, `docker-compose.yml`, `scripts/docker-entrypoint.sh` |
| M3 | Implement `muara health` CLI command and container runtime hardening | ✅ | `internal/cli/health.go`, `internal/cli/health_test.go`, registered in `internal/cli/root.go` |
| M4 | Harden install script with SHA256, cosign, SKIP_VERIFY, OS/arch matrix | ✅ | `scripts/install.sh` |
| M5 | Extend CI workflow: minimal permissions, SHA pinning, scheduled builds, docker-build, install-dry-run, changelog-check | ✅ | `.github/workflows/ci.yml` |
| M6 | Complete release documentation, branch protection docs, and README badges | ✅ | `docs/install.md`, `docs/security/cve-exceptions.md`, `runbooks/release.md`, `.actrc`, `docs/contributing.md`, `AGENTS.md`, `README.md` |

---

## Side Improvements

| Item | Status | Notes |
|---|---|---|
| SQLite concurrency fix | ✅ | `db.SetMaxOpenConns(1)` + `db.SetMaxIdleConns(1)` in `internal/cli/start.go` |
| SQLite DSN tuning | ✅ | `_busy_timeout=5000&_journal_mode=WAL` across all SQLite stores |
| Gitleaks allowlist | ✅ | `.gitleaks.toml` for shell-script env-var auth false positives |
| Tracker consistency | ✅ | Created migration-guide entry points; synced `DefaultYAML()` with `muara.yml.example` |
| Install script portability | ✅ | Manual hash comparison works on both GNU sha256sum and macOS shasum |
| Docs accuracy | ✅ | `SKIP_VERIFY` example corrected; README badge and `go install` restored |
| Local CI image | ✅ | `.actrc` uses `catthehacker/ubuntu:act-22.04` |
| Lint exclusions | ✅ | `.golangci.yml` excludes `node_modules` from Go lint scans |
| Coverage script | ✅ | `scripts/check-coverage.sh` excludes `node_modules` Go packages from coverage |

---

## Quality Gate Results

| Gate | Result |
|---|---|
| `go build ./...` | ✅ |
| `go fmt` | ✅ |
| `go vet ./...` | ✅ |
| `golangci-lint run` | ✅ 0 issues |
| `go test ./...` | ✅ |
| `actionlint .github/workflows/*.yml` | ✅ |
| Coverage ≥ 80% | ✅ 80.3% |
| UI build + tests + bundle size | ✅ |
| `./scripts/smoke-test.sh` | ✅ |
| `govulncheck` | ✅ 0 reachable |
| `gosec` | ✅ 0 issues |
| `gitleaks` | ✅ 0 leaks |
| Forbidden patterns | ✅ |
| Shell scripts | ✅ |
| Size checks | ✅ advisory warnings |
| Tracker audit | ✅ |
| Full `task quality` | ✅ |

---

## Sign-off

- **Planned by:** AI Agent + user review  
- **Implemented by:** AI Agent  
- **Verified by:** `task quality` end-to-end  
- **User sign-off:** ✅ Approved for delivery
