---
sidebar_position: 11
---

# HTTP API

This section provides a detailed reference for the Backupman HTTP API.

## Authentication

All API endpoints under the `/api` group require authentication. You must provide a valid bearer token in the `X-Api-Key` header.

```
X-Api-Key: <your-token>
```

## Usage Example

Here is an example of how to use the API to trigger a backup:

```bash
curl -X POST \
  -H "X-Api-Key: <your-token>" \
  http://localhost:8080/api/backups
```

:::info
Click here to view the full [HTTP API Reference](/docs/category/http-api)
:::

