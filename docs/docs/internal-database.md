---
sidebar_position: 4
description: "Backupman uses an internal database to store information about backups. Many providers are supported."
---

# Internal Database

Backupman uses an internal database to store information about backups. Many providers are supported.

## MySQL

You can use the following configuration:

```yaml title="config.yml"
database:
  provider: mysql
  host: 127.0.0.1
  port: 3306
  db_name: ChangeMe
  user: ChangeMe
  password: ChangeMe
  tls: false
```

## PostgreSQL

You can use the following configuration:

```yaml title="config.yml"
database:
  provider: postgres
  host: 127.0.0.1
  port: 5432
  db_name: ChangeMe
  user: ChangeMe
  password: ChangeMe
  tls: false
```

## SQLite

You can use the following configuration:

```yaml title="config.yml"
database:
  provider: sqlite
  db_path: /var/lib/backupman/backupman.db
```

:::note
The database file and directory will be created automatically if they don't exist. Make sure the application has write permissions to the specified path.
:::

## Memory

The memory provider stores data in RAM and is useful for development and testing. **Data is lost when the application stops.**

You can use the following configuration:

```yaml title="config.yml"
database:
  provider: memory
```
