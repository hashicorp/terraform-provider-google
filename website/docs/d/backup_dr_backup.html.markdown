---
subcategory: "Backup and DR Backup"
description: |-
  Get information about a Backupdr Backup.
---

# google_backup_dr_backup

A Backup and DR Backup.

## Example Usage

```hcl
data "google_backup_dr_backup" "foo" {
  location      = "us-central1"
  project = "project-test"
  data_source_id = "ds-test"
  backup_vault_id = "bv-test"
}
```

## Argument Reference

The following arguments are supported:

* `location` - (Required) The location in which the Backup belongs.
* `project` - (Required) The Google Cloud Project in which the Backup belongs.
* `data_source_id` - (Required) The ID of the Data Source in which the Backup belongs.
* `backup_vault_id` - (Required) The ID of the Backup Vault of the Data Source in which the Backup belongs.