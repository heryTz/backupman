---
sidebar_position: 5
description: "Backupman supports various database to backup. Each source has its own configuration options."
---

# Data Sources

Backupman supports various database to backup. Each source has its own configuration options.

:::info
You can backup multiple database sources.
:::

## MySQL

You can use the following configuration:

```yaml title="config.yml"
data_sources:
  - provider: mysql
    label: MySQL 1
    host: 127.0.0.1
    port: 3306
    db_name: ChangeMe
    user: ChangeMe
    password: ChangeMe
    tls: false
    # Temporary folder used by Backupman (for example, to store dumps before uploading to cloud)
    tmp_folder: ./tmp/mysql
```
