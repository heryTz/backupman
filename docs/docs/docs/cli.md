---
sidebar_position: 12
---

# CLI

Backupman provides a command-line interface to manage backups.

## Main Command

Here is the main help output for the `backupman` command:

```bash
A command line tool for managing backups.

Usage:
  backupman [flags]
  backupman [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  health      Health check
  help        Help about any command
  retry       Retry a failed backup
  run         Run the backup
  serve       Serve the backup manager
  version     Version information

Flags:
  -c, --config string   Path to the config file (default "./config.yml")
  -h, --help            help for backupman

Use "backupman [command] --help" for more information about a command.
```

:::info
For a more detailed breakdown of each command and its flags, please see the [CLI Reference](/docs/references/cli).
:::
