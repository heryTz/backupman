---
sidebar_position: 2
description: "Backupman can backup your data to Google Drive."
---

# Google Drive

Backupman can backup your data to Google Drive.

## Configuration

```yaml title="config.yml"
drives:
  - provider: google_drive
    label: Google Drive
    folder: demo
    client_secret_file: /path/to/your/google-client-secret.json
    token_file: /path/to/your/google-token.json
```

Backupman will create a folder named `demo` in your Google Drive root directory.

## How to get the client secret and token

### 1. Get Client Secret 

- Go to [Google Cloud Console](https://console.cloud.google.com/)
- Create a New Project
- Enable Google Drive API
- Configure OAuth Consent Screen
- Add Test Users
- Create OAuth Credentials (as a Desktop App)
- Download Credentials
- Save the file as `google-client-secret.json`
- Place this file in your root project directory

### 2. Get Token 

You can use the [prebuild binary](https://github.com/heryTz/backupman/releases) to get the token. Run the following command:

```bash
backupman auth-google
```

This will open a browser window to authenticate your Google account and generate the token file (`google-token.json` by default).

See [auth-google](/docs/references/cli) for more details on the command.
