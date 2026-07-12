> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# Appendix A — Provider Contract Checklist

Use this checklist when adding or reviewing a provider manifest.

## Manifest Structure

- [ ] File is at `plugins/<name>/gateway.yml`.
- [ ] `name` matches the directory name `<name>`.
- [ ] `runtime.type` is one of: `simple`, `go`.
- [ ] `version` is present and follows semver.
- [ ] `display_name` and `description` are human-readable.

## Simple Runtime (`runtime.type: simple`)

- [ ] `simple` block is present.
- [ ] Required endpoints are declared.
- [ ] Signature scheme is declared (if applicable).
- [ ] Webhook config is declared (if applicable).
- [ ] No Go registration exists for this provider.

## Go Runtime (`runtime.type: go`)

- [ ] A Go package exists at `internal/<name>/`.
- [ ] A factory registration exists at `internal/<name>/register.go`.
- [ ] The factory signature matches the registry's `Factory` type.
- [ ] The manifest activates the factory by declaring `runtime.type: go`.
- [ ] Provider protocol emulation is covered by conformance tests.

## Config Schema

- [ ] Required config keys are documented.
- [ ] Defaults are explicit.
- [ ] Secrets are loaded from env vars or `.muara/config.yml`, never hard-coded.
- [ ] Validation rules are implemented in `internal/plugin/validator.go`.

## Tests

- [ ] Unit tests for manifest parsing.
- [ ] Conformance tests for protocol emulation.
- [ ] Loader integration test proving the provider is discovered.
- [ ] Regression test if fixing a bug.

## Documentation

- [ ] Provider is listed in `docs/providers.md` (or equivalent).
- [ ] `docs/contributing-providers.md` references the provider as an example if it demonstrates a pattern.
