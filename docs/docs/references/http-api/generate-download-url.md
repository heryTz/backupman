---
sidebar_position: 5
title: Generate Download URL
---

:::warning
Currently, this endpoint supports only the `local` storage provider. Therefore, **you must have** a `local` provider configured in your storage settings. We are working on adding support for other providers in the future. 
:::

# Generate Download URL

Generates a temporary, pre-signed URL to download a backup file. The URL generated will allow you to download the file without needing to authenticate.

`GET /api/backups/:id/generate-download-url`

**Path Parameters:**

| Parameter | Description |
| :--- | :--- |
| `id` | The ID of the backup to generate a URL for. |


**Example Response (200 OK):**

```json
{
  "Url": "https://storage.googleapis.com/..."
}
```

**Error Response (500 Internal Server Error):**

```json
{
  "Error": "Some error message"
}
```
