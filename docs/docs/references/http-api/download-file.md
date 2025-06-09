---
sidebar_position: 6
title: Download File
---

:::warning
Currently, this endpoint supports only the `local` storage provider. Therefore, **you must have** a `local` provider configured in your storage settings. We are working on a feature to support direct downloads from any provider in the future.
:::

# Download File

Downloads the backup file directly.

`GET /api/backups/:id/download`

**Path Parameters:**

| Parameter | Description |
| :--- | :--- |
| `id` | The ID of the drive file to download. |


**Successful Response (200 OK):**

The response will be the raw file data with the following headers:

*   `Content-Disposition: attachment; filename="..."`
*   `Content-Type: <mime-type>`

**Error Response (500 Internal Server Error):**

```json
{
  "Error": "Some error message"
}
```
