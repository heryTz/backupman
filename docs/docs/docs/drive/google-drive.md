---
sidebar_position: 2
description: "Backupman can backup your data to Google Drive."
---

# Google Drive

Backupman can backup your data to Google Drive.

> **Prerequisites**
> - A Google account with access to Google Drive. [This](https://cloud.google.com/iam/docs/service-accounts-create).
> - Share the folder inside Google Drive with the service account email.

## Configuration

After creating a service account, download the JSON key file and use it in your configuration.

```yaml title="config.yml"
drives:
  - provider: google_drive
    label: Google Drive
    folder: demo
    service_account: /path/to/your/service-account.json
```
