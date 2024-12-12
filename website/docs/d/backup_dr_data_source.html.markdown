---
subcategory: "Backup and DR Data Source"
description: |-
  Get information about a Backupdr Data Source.
---

# google_backup_dr_data_source

A Backup and DR Data Source.

## Example Usage

```hcl
data "google_backup_dr_data_source" "foo" {
    location      = "us-central1"
    project = "project-test"
    data_source_id = "ds-test"
    backup_vault_id = "bv-test"
}
```

## Argument Reference

The following arguments are supported:

* `location` - (Required) The location in which the Data Source belongs.
* `project` - (Required) The Google Cloud Project in which the Data Source belongs.
* `data_source_id` - (Required) The ID of the Data Source.
* `backup_vault_id` - (Required) The ID of the Backup Vault in which the Data Source belongs.