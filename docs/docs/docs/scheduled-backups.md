---
sidebar_position: 13 
title: Scheduled Backups
---

:::warning
The cron automation is currently part of the HTTP server. In the future, this will be a standalone feature.
:::

# Scheduled Backups

The backup job can be automated using a cron schedule.

To enable the backup job, you need to configure the `http.backup_job` section in your configuration file.

```yaml title="config.yml"
http:
  backup_job:
    enabled: true
    cron: "* * * * * *"
```

The `cron` field uses the standard cron format. You can find more information about the cron format [here](https://crontab.guru/).
