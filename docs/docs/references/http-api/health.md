---
sidebar_position: 2
---

# Health Check

Performs a health check on the system and its components (database, drives, etc.).

`GET /health`

**Example Response (200 OK):**

```json
{
  "Version": "1.0.0",
  "CommitSHA": "a1b2c3d",
  "BuildDate": "2025-07-02T12:00:00Z",
  "Status": "UP",
  "Details": {
    "Database": {
      "Status": "UP",
      "Components": null
    },
    "Drives": {
      "Status": "UP",
      "Components": {
        "my_google_drive_label": {
          "Status": "UP",
          "Components": null
        }
      }
    },
    "DataSources": {
      "Status": "UP",
      "Components": {
        "my_database_label": {
          "Status": "UP",
          "Components": null
        }
      }
    }
  }
}
```
