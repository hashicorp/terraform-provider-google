---
subcategory: "Cloud SQL"
description: |-
  Get Latest Recovery Time for a given instance.
---

# google_sql_database_instance_latest_recovery_time

Get Latest Recovery Time for a given instance. For more information see the
[official documentation](https://cloud.google.com/sql/)
and
[API](https://cloud.google.com/sql/docs/postgres/backup-recovery/pitr#get-the-latest-recovery-time).


## Example Usage

```hcl
data "google_sql_database_instance_latest_recovery_time" "default" {
  instance = "sample-instance"
}

output "latest_recovery_time" {
  value = data.google_sql_database_instance_latest_recovery_time.default
}
```

## Argument Reference

The following arguments are supported:

* `instance` - (Required) The name of the instance.

## Attributes Reference

The following attributes are exported:

* `instance` - The name of the instance.
* `project` - The ID of the project in which the resource belongs.
* `latest_recovery_time` - Timestamp, identifies the latest recovery time of the source instance.
