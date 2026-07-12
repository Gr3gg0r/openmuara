# AGENTS.md — OpenMuara Workspace

> Rules for AI coding agents working on **OpenMuara**.
> OpenMuara is a local-first billing & payment virtualization layer. It emulates payment providers (Stripe, RevenueCat, App Store, Google Play, Fawry, SenangPay, iPay88, Billplz) so developers can test financial infrastructure offline, fast, and headlessly.
> Workspace language is **English** unless a feature explicitly requires other languages.

---

## 1. Workspace layout

- **Root repo directory** on local disk may still be named `toyol` for historical reasons. User-facing names, the module path, and all tracked content use `openmuara`.
- **Product code** lives in this repo.
- **Planning docs** live in `prompts/`, `tasks/`, `DECISIONS.md`, and `TRACKING.md`.
- `.muara/` is the **user-local workspace**. It holds `config.yml`, `data/`, `plugins/`, and runtime state. It is ignored by git.
- `.agents/` is the **AI workflow workspace**. It is local-only, uncommitted. (Create if needed for large initiatives.)

### Project stack

| Alias | Path | Stack | Status |
|---|---|---|---|
| `openmuara` | `/` | Go | ✅ Active |
| `muara-cli` | `/cmd/muara/` | Go | ✅ Implemented |
| `openmuara-mcp` | `/cmd/mcp-server/` | Go | ⬜ Planned |
| `openmuara-web` | `/internal/ui/` (embedded) | HTML/CSS/JS | ✅ Implemented |

---

## 2. Commit boundaries

- **Governance/docs** → edit and commit in root repo.
- **Product code** → edit and commit in the repository that owns the file.
- **One logical change per commit.** Run `git status` before editing.

---

## 3. Branch rules

**Protected branches:** `main` and `dev`.

- **Never** commit or merge directly to `main`.
- **Default working branch:** `dev`.
- Work happens on `dev` unless you explicitly create a feature branch: `feat/<description>` or `fix/<description>`.

### Required status checks

The following checks must pass before merging to `dev` or `main`.

- `docs` — markdown lint, link check, OpenAPI validation, Docusaurus build.
- `ui-build` — dashboard build and bundle-size check.
- `ui-test` — dashboard unit tests.
- `lint` — `gofmt`, `go vet`, `golangci-lint`.
- `unit` — Go unit tests with race detector and coverage floors.
- `smoke` — end-to-end smoke test.
- `vuln` — `govulncheck`.
- `gosec` — security scan (SARIF upload).
- `secrets` — `gitleaks` secret scan.
- `quality` — full `task quality` matrix.
- `dependency-license` — Go module verification, license check, npm audit.
- `docker-build` — container builds and reports healthy.
- `install-dry-run` — install script detection works on linux and macOS.
- `changelog-check` — `CHANGELOG.md` has a section matching `VERSION`.

### Branch protection settings

- Require a pull request before merging.
- Require status checks to pass before merging (list above).
- Require signed commits on `main`.
- Do not allow bypassing the above settings.
- Restrict pushes that create files larger than 100 MB.

---

## 4. API change contract

- Provider protocol emulation must be faithful to documented provider behavior unless explicitly marked as a limitation.
- Provider plugin schema changes (`plugins/*/gateway.yml`) are breaking. Version them carefully.
- CLI changes should be backward-compatible within a major version.

### Data flow

- **Incoming:** test app → `internal/server/router.go` → provider → `internal/engine/` (SQLite ledger).
- **Idempotency:** router checks `Idempotency-Key` against ledger before processing.
- **Mobile receipts:** App Store / Play Store receipts are treated as lookup keys in `.muara/data/unified_matrix.json` — no real crypto decoding.
- **Outgoing:** scheduler → webhook dispatcher → test app webhook handler.
- **Dashboard:** engine events → SSE stream → `/_admin`.

---

## 5. Task runner

- `go run ./cmd/muara`
- `go test ./...`
- `go build ./...`
- `go vet ./...`
- `golangci-lint run`

---

## 6. Code quality enforcement

All work should pass with **zero warnings** before commit. New features and bug fixes should include tests.

| Gate | Command |
|---|---|
| Build | `go build ./...` |
| Test | `go test ./...` |
| Lint | `golangci-lint run` |
| Vet | `go vet ./...` |

### File-size gates

- Max **250 lines** per file (recommended).
- Max **80 lines** per function (recommended).
- Max **120 chars** per line (recommended).

---

## 7. Code style

- Prefer explicit types and strong contracts.
- Avoid `any` / `interface{}` where a concrete type fits.
- No debug `fmt.Println` in committed code.
- Remove unused imports.

---

## 8. Agent autonomy boundaries

### Must ask the user before:

1. Changing **P0 integration** logic (Stripe, RevenueCat, App Store, Play Store, Fawry, SenangPay, regional gateways).
2. Creating or deleting **database migrations** or persistence schemas.
3. Modifying **auth, billing, or PII-handling** flows.
4. **Refactoring across more than two modules** in a single session.
5. Changing the **provider plugin schema contract**.

### May decide independently:

- Internal module structure and naming.
- Utility functions and helpers.
- Test coverage approach.
- Minor refactoring within a single module.

---

## 9. Agent Workflow

Before starting any multi-phase initiative, read `TRACKING.md` and `DECISIONS.md`.

For work aligned with a prompt:

1. Read the prompt in `prompts/`.
2. Check `TRACKING.md` for status.
3. Implement, test, and commit.
4. Update `TRACKING.md`.
5. Log non-trivial decisions in `DECISIONS.md`.

---

## 10. Security

- **Never commit** real `.env` files or `.muara/config.yml`.
- API secrets are server-side only. OpenMuara is local-only by design.
- Webhook signature verification must be faithfully emulated. Do not bypass HMAC/SHA256 checks in dev mode.
- OpenMuara must never proxy traffic to real provider endpoints in default mode.

---

## 11. Integrations

### P0 (core mission)

- Stripe Checkout & Billing
- RevenueCat
- App Store / Google Play
- Fawry
- SenangPay

### Other important

- iPay88, Billplz, Razer Merchant Services
- Docker / container runtime
- CI/CD pipelines

---

## 12. Local development runbook

1. Build: `go build ./...`
2. Test: `go test ./...`
3. Run: `go run ./cmd/muara start`
4. Configure test app to point at `http://localhost:9000`.
5. Use `/_admin` dashboard or CLI to inspect state and replay webhooks.

---

## 13. Legacy reference

- Previous project name: `muara`
- Migration guide: `docs/migration/openmuara-to-openmuara.md` (Prompt 18)

---

## 14. When this document is incomplete

If a command, port, or env var here disagrees with a nested **README**, **prompt**, or **task spec**, follow the more specific document and update this file to match.
