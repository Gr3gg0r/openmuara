> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# Runbook — Validate a Provider Manifest

## When to Use

After adding or editing a `plugins/<name>/gateway.yml` file.

## Steps

1. Lint the YAML syntax:

   ```bash
   yamllint plugins/<name>/gateway.yml
   ```

2. Validate against the OpenMuara schema:

   ```bash
   go run ./cmd/muara provider validate plugins/<name>/gateway.yml
   ```

3. Check runtime type consistency:

   - `runtime.type: simple` → no Go registration needed.
   - `runtime.type: go` → a factory must be registered in `internal/<name>/register.go`.

   See `appendices/b-simple-vs-go-decision-tree.md`.

4. Run the provider contract checklist:

   See `appendices/a-provider-contract-checklist.md`.

5. Run the provider conformance test:

   ```bash
   go test ./internal/provider/conform/... -run TestProvider/<name>
   ```

6. Run full gates:

   ```bash
   go test ./...
   go vet ./...
   golangci-lint run ./...
   ```

## Common Issues

| Issue | Fix |
|---|---|
| `runtime.type` missing | Add `runtime.type: simple` or `runtime.type: go`. |
| Factory not found | Ensure `internal/<name>/register.go` calls `hybrid.Register`. |
| Config schema mismatch | Update `internal/plugin/schema.go` or the manifest. |
| Test depends on global provider | Update test to load manifest or use factory directly. |
