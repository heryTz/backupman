---
sidebar_position: 1
---

# CLI 

This page provides a reference for the `backupman` command-line interface.

## Main Command

The main `backupman` command is used to manage backups.

**Usage:**

```bash
backupman [command]
```

**Global Flags:**

| Flag | Description | Default |
| :--- | :--- | :--- |
| `-c`, `--config` | Path to the config file. | `./config.yml` |
| `-h`, `--help` | Help for `backupman`. | |

---

## Commands

### `auth-google`

Authenticate with Google Drive using OAuth2

**Usage:**

```bash
backupman auth-google
```

**Arguments:**

| Argument | Description | Default |
| :--- | :--- | :--- |
| `client-secret-file` | Path to the client secret file. | `google-client-secret.json` |
| `token-file` | Path that the token will be saved to. | `google-token.json` |
| `open-url` | Open the authorization URL in the default web browser. | `true` |

### `completion`

Generate the autocompletion script for the specified shell.

**Usage:**

```bash
backupman completion [shell]
```

### `health`

Perform a health check on the backup system.

**Usage:**

```bash
backupman health
```

### `retry`

Retry a failed backup.

**Usage:**

```bash
backupman retry [id]
```

**Arguments:**

| Argument | Description |
| :--- | :--- |
| `id` | The ID of the failed backup to retry. |

### `run`

Run the backup.

**Usage:**

```bash
backupman run
```

### `serve`

Serve the backup manager.

**Usage:**

```bash
backupman serve
```

**Flags:**

| Flag | Description | Default |
| :--- | :--- | :--- |
| `-p`, `--port` | Port to run the server on. | `8080` |

### `version`

Display the version information of backupman.

**Usage:**

```bash
backupman version
```
