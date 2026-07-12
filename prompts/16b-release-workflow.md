# Prompt 16b — Release Workflow

## Goal

Define and automate the release process.

## Acceptance Criteria

- [ ] Versioning follows semantic versioning
- [ ] `CHANGELOG.md` maintained
- [ ] `docs/openapi.yaml` `info.version` is bumped to match `VERSION` before tagging
- [ ] Git tag triggers release workflow
- [ ] Release artifacts:
  - binaries for linux/darwin/windows (amd64 + arm64)
  - Docker image
- [ ] GitHub Release notes auto-generated from changelog

## Files to Create/Change

- `CHANGELOG.md`
- `.github/workflows/release.yml`
- `Taskfile.yml` — release tasks
- `README.md` — install instructions

## Response Shape

Return:

1. Version bump process
2. CI release job outline
3. Artifact list

## Test Notes

- Create a test tag and verify workflow (dry-run if possible)
