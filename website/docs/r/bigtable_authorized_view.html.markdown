---
subcategory: "Cloud Bigtable"
description: |-
  Creates a Google Cloud Bigtable authorized view inside a table.
---

# google_bigtable_authorized_view

Creates a Google Cloud Bigtable authorized view inside a table. For more information see
[the official documentation](https://cloud.google.com/bigtable/) and
[API](https://cloud.google.com/bigtable/docs/go/reference).

-> **Note:** It is strongly recommended to set `lifecycle { prevent_destroy = true }`
on authorized views in order to prevent accidental data loss. See
[Terraform docs](https://www.terraform.io/docs/configuration/resources.html#prevent_destroy)
for more information on lifecycle parameters.


## Example Usage

```hcl
resource "google_bigtable_instance" "instance" {
  name = "tf-instance"

  cluster {
    cluster_id   = "tf-instance-cluster"
    zone         = "us-central1-b"
    num_nodes    = 3
    storage_type = "HDD"
  }

  lifecycle {
    prevent_destroy = true
  }
}

resource "google_bigtable_table" "table" {
  name          = "tf-table"
  instance_name = google_bigtable_instance.instance.name
  split_keys    = ["a", "b", "c"]

  lifecycle {
    prevent_destroy = true
  }

  column_family {
    family = "family-first"
  }

  column_family {
    family = "family-second"
  }

  change_stream_retention = "24h0m0s"
}

resource "google_bigtable_authorized_view" "authorized_view" {
  name         = "tf-authorized-view"
  instance_name = google_bigtable_instance.instance.name
  table_name = google_bigtable_table.table.name

  lifecycle {
    prevent_destroy = true
  }

  subset_view {
    row_prefixes = [base64encode("prefix#)]

    family_subsets {
      family_name = "family-first"
      qualifiers = [base64encode("qualifier"), base64encode("qualifier-second")]
    }

    family_subsets {
      family_name = "family-second"
      qualifier_prefixes = [""]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the authorized view. Must be 1-50 characters and must only contain hyphens, underscores, periods, letters and numbers.

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

* `instance_name` - (Required) The name of the Bigtable instance in which the authorized view belongs.

* `table_name` - (Required) The name of the Bigtable table in which the authorized view belongs.

* `column_family` - (Optional) A group of columns within a table which share a common configuration. This can be specified multiple times. Structure is documented below.

* `deletion_protection` - (Optional) A field to make the table protected against data loss i.e. when set to PROTECTED, deleting the table, the column families in the table, and the instance containing the table would be prohibited.
If not provided, currently deletion protection will be set to UNPROTECTED as it is the API default value. Note this field configs the deletion protection provided by the API in the backend, and should not be confused with Terraform-side deletion protection.

* `subset_view` - (Optional) An AuthorizedView permitting access to an explicit subset of a Table. Structure is documented below.

-----

`subset_view` supports the following arguments:

* `row_prefixes` - (Optional) A list of Base64-encoded row prefixes to be included in the authorized view. To provide access to all rows, include the empty string as a prefix ("").

* `family_subsets` - (Optional) A group of column family subsets to be included in the authorized view. This can be specified multiple times. Structure is documented below.

-----

`family_subsets` supports the following arguments:

* `family_name` - (Required) Name of the column family to be included in the authorized view. The specified column family must exist in the parent table of this authorized view.

* `qualifiers` - (Optional) A list of Base64-encoded individual exact column qualifiers of the column family to be included in the authorized view.

* `qualifier_prefixes` - (Optional) A list of Base64-encoded prefixes for qualifiers of the column family to be included in the authorized view.
Every qualifier starting with one of these prefixes is included in the authorized view. To provide access to all qualifiers, include the empty string as a prefix ("").

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `id` - an identifier for the resource with format `projects/{{project}}/instances/{{instance_name}}/tables/{{table_name}}/authorizedViews/{{name}}`

## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.

## Import

Bigtable Authorized Views can be imported using any of these accepted formats:

* `projects/{{project}}/instances/{{instance_name}}/tables/{{table_name}}/authorizedViews/{{name}}`
* `{{project}}/{{instance_name}}/{{table_name}}/{{name}}`
* `{{instance_name}}/{{table_name}}/{{name}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Bigtable Authorized Views using one of the formats above. For example:

```tf
import {
  id = "projects/{{project}}/instances/{{instance_name}}/tables/{{table_name}}/authorizedViews/{{name}}"
  to = google_bigtable_authorized_view.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), Bigtable Authorized Views can be imported using one of the formats above. For example:

```
$ terraform import google_bigtable_authorized_view.default projects/{{project}}/instances/{{instance_name}}/tables/{{table_name}}/authorizedViews/{{name}}
$ terraform import google_bigtable_authorized_view.default {{project}}/{{instance_name}}/{{table_name}}/{{name}}
$ terraform import google_bigtable_authorized_view.default {{instance_name}}/{{table_name}}/{{name}}
```


