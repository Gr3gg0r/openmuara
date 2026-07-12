---
id: cli
title: CLI Reference
---

# CLI Reference

The `muara` command-line tool manages workspaces, runs the server, inspects
state, and simulates payment flows.

## Global flags

| Flag | Default | Description |
|---|---|---|
| `--config` | `.muara/config.yml` | Path to the configuration file |
| `--json` | false | Output results as JSON where supported |
| `--quiet` | false | Suppress non-error output |
| `--help` | — | Show help for any command |

## Commands

| Command | Purpose |
|---|---|
| `muara init` | Initialize a local `.muara/` workspace |
| `muara start` | Start the OpenMuara server |
| `muara doctor` | Check the environment and configuration |
| `muara scenario` | Simulate payment outcomes |
| `muara webhook` | Inspect and replay outgoing webhooks |
| `muara audit` | Inspect the audit log |
| `muara transaction` | Inspect transactions |
| `muara plugins` | List and validate provider plugins |
| `muara provider` | Test and scaffold YAML-driven simple providers |
| `muara security` | Security helpers |
| `muara clean` | Reset local data |
| `muara completion` | Generate shell completion scripts |
| `muara version` | Print version information |

---

## `muara init`

Initialize a local workspace.

```bash
# Interactive wizard
muara init

# Non-interactive default config
muara init --defaults

# Preview without writing files
muara init --dry-run

# Overwrite an existing config
muara init --force
```

Flags: `--defaults`, `--dry-run`, `--force`

---

## `muara start`

Start the server.

```bash
muara start
muara start --config path/to/config.yml
muara start --dry-run
```

Flags: `--dry-run`, `--no-banner`

---

## `muara doctor`

Check the environment and config.

```bash
muara doctor
muara doctor --json
muara doctor --check-webhook
```

Flags: `--check-webhook`

---

## `muara scenario`

Simulate payment outcomes for a transaction reference.

```bash
muara scenario success tx-123
muara scenario fail tx-123
muara scenario timeout tx-123
```

Subcommands: `success`, `fail`, `timeout`

---

## `muara webhook`

Inspect and replay outgoing webhooks.

```bash
muara webhook list
muara webhook inspect tx-123
muara webhook replay tx-123
```

Subcommands: `list`, `inspect`, `replay`

---

## `muara audit`

Inspect the audit log.

```bash
muara audit list
muara audit list --since 2026-01-01T00:00:00Z
```

Subcommands: `list`

---

## `muara transaction`

Inspect transactions.

```bash
muara transaction list
muara transaction inspect tx-123
```

Subcommands: `list`, `inspect`

---

## `muara plugins`

List and validate declarative provider plugins.

```bash
muara plugins list
muara plugins validate
```

Subcommands: `list`, `validate`

---

## `muara provider`

Test and scaffold YAML-driven simple providers.

```bash
muara provider test fawry
muara provider init my-gateway
```

Subcommands: `init`, `test`

---

## `muara security`

Security helpers.

```bash
muara security audit
muara security hash-password --password mypassword
muara security gen-cert
```

Subcommands: `audit`, `gen-cert`, `hash-password`

---

## `muara clean`

Remove the local SQLite database used for the ledger, audit log, and webhook
attempts.

```bash
muara clean
muara clean --force
```

Flag: `--force`

---

## `muara completion`

Generate shell completion scripts.

```bash
muara completion bash > /usr/local/etc/bash_completion.d/muara
muara completion zsh > "${fpath[1]}/_muara"
muara completion fish > ~/.config/fish/completions/muara.fish
```

---

## `muara version`

Print version information.

```bash
muara version
muara version --json
```

## JSON schemas

Command-specific JSON schemas live in `docs/cli-schemas/`.
