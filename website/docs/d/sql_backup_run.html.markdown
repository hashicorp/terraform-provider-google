---
subcategory: "Cloud SQL"
description: |-
  Get a  SQL backup run in Google Cloud SQL.
---

# google_sql_backup_run

Use this data source to get information about a Cloud SQL instance backup run.

## Example Usage 

```hcl
data "google_sql_backup_run" "backup" {
	instance = google_sql_database_instance.main.name
	most_recent = true
}
```

## Argument Reference

The following arguments are supported:

* `instance` - (required) The name of the instance the backup is taken from.

* `backup_id` - (optional) The identifier for this backup run. Unique only for a specific Cloud SQL instance.
    If left empty and multiple backups exist for the instance, `most_recent` must be set to `true`.

* `most_recent` - (optional) Toggles use of the most recent backup run if multiple backups exist for a 
    Cloud SQL instance.

* `project` - (Optional) The project to list instances for. If it
    is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:
    
* `location` -  Location of the backups.

* `start_time` - The time the backup operation actually started in UTC timezone in RFC 3339 format, for 
    example 2012-11-15T16:19:00.094Z.

* `status` - The status of this run. Refer to [API reference](https://cloud.google.com/sql/docs/mysql/admin-api/rest/v1beta4/backupRuns#SqlBackupRunStatus) for possible status values.