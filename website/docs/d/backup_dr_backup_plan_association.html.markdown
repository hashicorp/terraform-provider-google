---
subcategory: "Backup and DR BackupPlanAssociation"
description: |-
  Get information about a Backupdr BackupPlanAssociation.
---

# google_backup_dr_backup_plan_association

A Backup and DR BackupPlanAssociation.

~> **Warning:** This resource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta resources.

## Example Usage

```hcl
data "google_backup_dr_backup_plan_association" "my-backupplan-association" {
  location =  "us-central1"
  backup_plan_association_id="bpa-id"
}
```

## Argument Reference

The following arguments are supported:

* `location` - (Required) The location in which the Backupplan association resource belongs.
* `backup_plan_association_id` - (Required) The id of Backupplan association resource.

- - -

## Attributes Reference

See [google_backup_dr_backup_plan_association](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/backup_dr_backup_plan_association) resource for details of the available attributes.
