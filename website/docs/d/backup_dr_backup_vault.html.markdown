---
subcategory: "Backup and DR BackupVault"
description: |-
  Get information about a Backupdr BackupVault.
---

# google_backup_dr_backup_vault

A Backup and DRBackupVault.

## Example Usage

```hcl
data "google_backup_dr_backup_vault" "my-backup-vault" {
  location =  "us-central1"
  backup_vault_id="bv-1"
}
```

## Argument Reference

The following arguments are supported:

* `location` - (Required) The location in which the Backup Vault resource belongs.
* `backup_vault_id` - (Required) The id of Backup Vault resource.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

See [google_backup_dr_backup_vault](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/backup_dr_backup_vault) resource for details of the available attributes.