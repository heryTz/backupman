---
sidebar_position: 8
description: "Backupman supports a variety of notifiers to send notifications about backup events."
---

# Notifiers

Backupman supports a variety of notifiers to send notifications about backup events.

:::info
For now, you will be notified about:
- `backup_report` : A report after a backup is completed.
:::

## Mail

It sends a notification via email. You can configure it:

```yaml title="config.yml"
notifiers:
  mail:
    enabled: true
    smtp_host: smtp.example.com
    smtp_port: 587
    smtp_user: user
    smtp_password: password
    smtp_crypto: ssl
    destinations:
      - name: John Doe
        email: john.doe@example.com
```

## Webhook

It sends a notification via a webhook. You can configure it:

```yaml title="config.yml"
notifiers:
  webhook:
    enabled: true
    endpoints:
      - name: My Webhook
        url: http://localhost:8080/webhook
        token: xxx
```

:::info
The request sent via webhook is a POST request with the event payload in the body. The content type is `application/json`.
:::

### Authentication

The token defines inside `notifiers.webhook.endpoints[].token` is used to authenticate the webhook requests. When Backupman sends a webhook event, it includes this token in the `X-Webhook-Token` header of the request.

### Events Payload

Here are the events payloads that Backupman sends to the webhook endpoints:

<details>
<summary><code>backup_report</code></summary>

This event is triggered when a backup job completes.

**Payload Example:**

```json
{
  "Event": "backup_report",
  "Payload": {
    "Id": "b4c3f1a0-1b1e-4b0e-8b4a-0e1b0e1b0e1b",
    "Status": "finished",
    "Label": "Mysql 1",
    "DumpPath": "/tmp/backupman/dump/backup-20250702100000.sql.gz",
    "CreatedAt": "2025-07-02T10:00:00Z",
    "UpdatedAt": "2025-07-02T10:05:00Z",
    "DriveFiles": [
      {
        "Id": "d1f1e1a0-1b1e-4b0e-8b4a-0e1b0e1b0e1c",
        "BackupId": "b4c3f1a0-1b1e-4b0e-8b4a-0e1b0e1b0e1b",
        "Provider": "google_drive",
        "Label": "My Google Drive",
        "Path": "backups/backup-20250702100000.sql.gz",
        "Status": "finished",
        "CreatedAt": "2025-07-02T10:01:00Z",
        "UpdatedAt": "2025-07-02T10:05:00Z"
      }
    ]
  }
}
```

**Payload Schema:**

The `Payload` object contains the following properties:

| Field | Type | Description | Nullable |
| :--- | :--- | :--- | :--- |
| `Id` | `string` | The unique identifier for the backup. | No |
| `Status` | `string` | The overall status of the backup (`pending`, `finished`, `failed`). | No |
| `Label` | `string` | The user-defined name for the backup job. | No |
| `DumpPath`| `string` | The local path where the database dump is stored. | Yes |
| `CreatedAt` | `string` | The timestamp when the backup was created (ISO 8601). | No |
| `UpdatedAt` | `string` | The timestamp when the backup was last updated (ISO 8601). | Yes |
| `DriveFiles`| `array` | A list of files associated with the backup, uploaded to different storage providers. | Yes |

**`DriveFiles` Object Schema:**

Each object in the `DriveFiles` array has the following structure:

| Field | Type | Description | Nullable |
| :--- | :--- | :--- | :--- |
| `Id` | `string` | The unique identifier for the drive file record. | No |
| `BackupId` | `string` | The ID of the parent backup. | No |
| `Provider` | `string` | The storage provider name (e.g., `google_drive`, `s3`). | No |
| `Label` | `string` | The name of the file on the storage provider. | No |
| `Path` | `string` | The full path or identifier for the file on the storage provider. | Yes |
| `Status` | `string` | The upload status for this specific file (`pending`, `finished`, `failed`). | No |
| `CreatedAt` | `string` | The timestamp when the file record was created (ISO 8601). | No |
| `UpdatedAt` | `string` | The timestamp when the file record was last updated (ISO 8601). | Yes |

</details>
