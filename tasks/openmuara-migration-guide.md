# Task T02 — OpenMuara-to-OpenMuara Migration Guide

> **Status:** ✅ Completed  
> **Related prompt:** [`prompts/18-migration-guide.md`](../prompts/18-migration-guide.md)

## Background

The project was rebranded from `muara` to `openmuara` while keeping the local workspace at `.muara/`,
the binary/CLI named `muara`, and the module path `github.com/openmuara/openmuara`. Users who started
on the legacy layout need a concise guide to align their workspace and configuration with the current
release.

## Migration checklist

1. **Back up the workspace**

   ```bash
   cp -r .muara .muara.backup.$(date +%Y%m%d)
   ```

2. **Update the binary**

   ```bash
   go install github.com/openmuara/openmuara/cmd/muara@latest
   ```

3. **Rename legacy config keys if present**

   | Legacy key | Current key |
   |---|---|
   | `server.port` | `server.port` (unchanged) |
   | `database.path` | `database.path` (unchanged) |
   | Provider sections under `providers.*` | Provider sections under `providers.*` (unchanged) |

   Most keys are unchanged; the rebrand affected names and module path only.

4. **Verify the installation**

   ```bash
   muara doctor --json
   ```

5. **Run smoke tests**

   ```bash
   muara scenario success
   ```

## Rollback

If anything fails, restore the backup:

```bash
rm -rf .muara
mv .muara.backup.20260101 .muara
```

## Entry point

The public guide lives at `docs/migration/openmuara-to-openmuara.md`.
