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
