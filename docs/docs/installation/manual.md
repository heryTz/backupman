---
sidebar_position: 1
description: "You have several options to install Backupman on your system."
---

# Manual Installation

You have several options to install Backupman on your system.

## Prebuilt binaries

You can download prebuilt binaries for your platform from the [releases page](https://github.com/heryTz/backupman/releases).

```bash
tar xvf backupman-<version>.tar.gz
chmod +x ./backupman
./backupman version
```

<details>
<summary>Output</summary>
```bash
Version: x.x.x
Commit SHA: xxx
Build Date: xxx
```
</details>

## Build from source

To build Backupman from source:

```bash
git clone https://github.com/heryTz/backupman.git
cd backupman
go build -o backupman
./backupman version
```

<details>
<summary>Output</summary>
```bash
Version: x.x.x
Commit SHA: xxx
Build Date: xxx
```
</details>
