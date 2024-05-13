---
subcategory: "Cloud SQL"
description: |-
  Get all available Cloud SQL tiers for the given project.
---

# google_sql_tiers

Get all available machine types (tiers) for a project, for example, db-custom-1-3840. For more information see the
[official documentation](https://cloud.google.com/sql/)
and
[API](https://cloud.google.com/sql/docs/mysql/admin-api/rest/v1beta4/tiers/list).


## Example Usage

```hcl
data "google_sql_tiers" "tiers" {
  project = "sample-project"
}

locals {
  all_available_tiers = [for v in data.google_sql_tiers.tiers.tiers : v.tier]
}

output "avaialble_tiers" {
  description = "List of all available tiers for give project."
  value       = local.all_available_tiers
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The Project ID for which to list tiers. If `project` is not provided, the project defined within the default provider configuration is used.

## Attributes Reference

The following attributes are exported:

* `tiers` - A list of all available machine types (tiers) for project. Each contains:
  * `tier` - An identifier for the machine type, for example, db-custom-1-3840.
  * `ram` - The maximum ram usage of this tier in bytes.
  * `disk_quota` - The maximum disk size of this tier in bytes.
  * `region` - The applicable regions for this tier.
