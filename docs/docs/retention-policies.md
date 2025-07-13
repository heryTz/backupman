---
sidebar_position: 7
description: "Retention policies refer to the rules that determine how long data is retained in a system before it is deleted."
---

# Retention Policies

Retention policies refer to the rules that determine how long data is retained in a system before it is deleted.

## By Age

Retention policies by age specify that data should be retained for a certain period of time.

```yaml title="config.yml"
retention:
  enabled: true
  by: age
  value: 30 #Only days periods are supported
```
