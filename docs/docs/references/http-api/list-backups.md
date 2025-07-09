---
sidebar_position: 3
---

# List Backups

Retrieves a list of all backups.

`GET /api/backups`

**Example Response (200 OK):**

```json
{
  "Results": [
    {
      "Id": "b4c3f1a0-1b1e-4b0e-8b4a-0e1b0e1b0e1b",
      "Status": "finished",
      "Label": "my-database-backup",
      "DumpPath": "/tmp/backupman/dump/backup-20250702100000.sql.gz",
      "CreatedAt": "2025-07-02T10:00:00Z",
      "UpdatedAt": "2025-07-02T10:05:00Z",
      "DriveFiles": [
        {
          "Id": "d1f1e1a0-1b1e-4b0e-8b4a-0e1b0e1b0e1c",
          "BackupId": "b4c3f1a0-1b1e-4b0e-8b4a-0e1b0e1b0e1b",
          "Provider": "google_drive",
          "Label": "backup-20250702100000.sql.gz",
          "Path": "backups/backup-20250702100000.sql.gz",
          "Status": "finished",
          "CreatedAt": "2025-07-02T10:01:00Z",
          "UpdatedAt": "2025-07-02T10:05:00Z"
        }
      ]
    }
  ]
}
```

**Response Schema:**

The `Results` array contains a list of backup objects with the following properties:

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
