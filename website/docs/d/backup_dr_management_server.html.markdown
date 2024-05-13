---
subcategory: "BackupDR Management Server"
description: |-
  Get information about a Backupdr Management server.
---

# google_backup_dr_management_server

Get information about a Google Backup DR Management server.

~> **Warning:** This resource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta resources.

## Example Usage

```hcl
data google_backup_dr_management_server my-backup-dr-management-server {
   location =  "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required) The region in which the management server resource belongs.

- - -

## Attributes Reference

See [google_backupdr_management_server](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/backup_dr_management_server) resource for details of the available attributes.
