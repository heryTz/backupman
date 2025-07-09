---
sidebar_position: 2 
description: Run Backupman in 5 Minutes.
---

# Quickstart

Run Backupman in 5 Minutes

> **Prerequisites for this quickstart**
> - Docker installed (see [Docker Installation Guide](https://docs.docker.com/get-docker/))
> - `curl` for API requests
> - `jq` for JSON processing

## 1. Project Directory

Create a directory for the demo project:

```bash
mkdir backupman-demo
cd backupman-demo
```

## 2. Docker Compose File 

Create a `compose.yml` file with the following content:

```yml title="compose.yml"
services:
  backupman_db:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: backupman
    volumes:
      - backupman_db_data:/var/lib/mysql

  backupman:
    image: herytz/backupman
    ports:
      - 8080:8080
    volumes:
      - ./config.yml:/app/config.yml
      - ./storage:/app/storage
    command: serve 

volumes:
  backupman_db_data:
```

## 3. Configuration File

Create a `config.yml` file with the following content:

```yaml title="config.yml"
database:
  provider: mysql
  host: backupman_db
  port: 3306
  db_name: backupman
  user: root
  password: root
  tls: false

# To simplify this quickstart, we use backupman database as the production database.
# You should use your own production database.
data_sources:
  - provider: mysql
    label: MySQL 1
    host: backupman_db
    port: 3306
    db_name: backupman
    user: root
    password: root
    tmp_folder: ./storage/tmp/mysql
    tls: false

drives:
  - provider: local
    label: Local Drive
    folder: ./storage/drive

http:
  app_url: http://localhost:8080
  api_keys:
    - apikey1
```

:::info

The configuration should contain:

1. [database](/docs/internal-database) : Backupman's internal database.
2. [data sources](/docs/backup-sources) : The databases you want to back up.
3. [drives](/docs/storage) : The storage drives where backups will be saved.

:::

## 4. Launch Services 

Run the following command to start the services:

```bash
docker compose up -d
```

## 5. Check Health Status

Check the health status of Backupman:

```bash
curl http://localhost:8080/health | jq .
```

<details>
<summary>Output</summary>
```json
{
  "Version": "xxx",
  "CommitSHA": "xxx",
  "BuildDate": "xxx",
  "Status": "UP",
  "Details": {
    "DataSources": {
      "Status": "UP",
      "Components": {
        "MySQL 1": {
          "Status": "UP",
          "Components": null
        }
      }
    },
    "Database": {
      "Status": "UP",
      "Components": null
    },
    "Drives": {
      "Status": "UP",
      "Components": {
        "Local Drive": {
          "Status": "UP",
          "Components": null
        }
      }
    }
  }
}
```
</details>

## 6. Trigger a Backup

We are going to trigger a backup via HTTP API but you can also use the [cli tool](/docs/cli-commands) or a [cron job](/docs/http-server).

```bash
curl -H "X-Api-Key: apikey1" -X POST http:/localhost:8080/api/backups | jq .
```

<details>
<summary>Output</summary>
```json
{
  "Message": "Backup started"
}
```
</details>

Wait a few minutes for the backup to complete.

:::note

The endpoint to monitor the backup status will be available soon.

:::

## 7. Verify Backup File 

Check the backup file in the local drive:

```bash
ls storage/drive
```

<details>
<summary>Output</summary>
```bash
1234xxxxxx.sql
```
</details>
