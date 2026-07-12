---
id: openmuara-to-openmuara
title: OpenMuara-to-OpenMuara Migration Guide
---

# OpenMuara-to-OpenMuara Migration Guide

This guide covers upgrades between OpenMuara releases and how to back up and
restore your local `.muara/` workspace.

## Compatibility matrix

| From | To | Breaking changes | Action required |
|---|---|---|---|
| pre-1.0 `muara` | 1.0.0 | Module path, binary name, config path renamed | Follow the upgrade steps below |
| 1.0.0 | 1.x | Minor releases are backward-compatible | Review CHANGELOG for new provider config options |

## Before you start

Back up your workspace:

```bash
cp -r .muara .muara.backup.$(date +%Y%m%d)
```

## Upgrading from pre-1.0 `muara` to OpenMuara 1.0.0

1. **Install the new binary**

   ```bash
   go install github.com/openmuara/openmuara/cmd/muara@latest
   ```

2. **Rename the config file if needed**

   The config file remains at `.muara/config.yml`. Most keys are unchanged.

3. **Update environment variables**

   Environment variables keep the `MUARA_` prefix. No changes are required for
   existing variables.

4. **Verify the installation**

   ```bash
   muara doctor --json
   ```

5. **Run a smoke test**

   ```bash
   muara scenario success tx-smoke-test
   ```

## Backing up and restoring `.muara/`

The `.muara/` directory contains:

- `config.yml` — runtime configuration.
- `data/ledger.db` — SQLite ledger, audit log, and webhook attempts.
- `plugins/` — any custom provider plugins.

To restore from a backup:

```bash
rm -rf .muara
mv .muara.backup.20260101 .muara
```

## Breaking config changes

When a release introduces a breaking config change, the release notes will list
a migration snippet. Run `muara start --dry-run` to validate your config before
starting the server.

## Rollback

If an upgrade fails:

1. Stop the server.
2. Restore the `.muara/` backup.
3. Re-install the previous binary version.
4. Run `muara doctor` to confirm.

## Getting help

- Read the latest `CHANGELOG.md`.
- Run `muara doctor --json` and include the output in bug reports.
- Use the **Documentation** issue template for migration-guide corrections.
