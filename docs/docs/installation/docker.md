---
sidebar_position: 2
description: "You can run Backupman using Docker."
---

# Docker

You can run Backupman using Docker.

## Image

The official Docker image for Backupman is available on [Docker Hub](https://hub.docker.com/r/herytz/backupman).

```bash
docker run -it --rm herytz/backupman version
```

<details>
<summary>Output</summary>
```bash
Version: x.x.x
Commit SHA: xxx
Build Date: xxx
```
</details>

:::warning
You should import your configuration file into the container. By default, Backupman will look the configuration file at `/app/config.yml`.
:::

## Full Compose Example

You can use the following `compose.yml` file to run Backupman:

```yaml title="compose.yml"
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
    volumes:
      - ./config.yml:/app/config.yml
      - ./storage:/app/storage
    ports:
      - 8080:8080
    command: serve

volumes:
  backupman_db_data:
```
