---
subcategory: "Cloud Bigtable"
layout: "google"
page_title: "Google: google_bigtable_instance"
sidebar_current: "docs-google-bigtable-instance"
description: |-
  Creates a Google Bigtable instance.
---

# google_bigtable_instance

Creates a Google Bigtable instance. For more information see
[the official documentation](https://cloud.google.com/bigtable/) and
[API](https://cloud.google.com/bigtable/docs/go/reference).

-> **Note**: It is strongly recommended to set `lifecycle { prevent_destroy = true }`
on instances in order to prevent accidental data loss. See
[Terraform docs](https://www.terraform.io/docs/configuration/resources.html#prevent_destroy)
for more information on lifecycle parameters.


## Example Usage - Production Instance

```hcl
resource "google_bigtable_instance" "production-instance" {
  name = "tf-instance"

  cluster {
    cluster_id   = "tf-instance-cluster"
    zone         = "us-central1-b"
    num_nodes    = 1
    storage_type = "HDD"
  }

  lifecycle {
    prevent_destroy = true
  }
}
```

## Example Usage - Development Instance

```hcl
resource "google_bigtable_instance" "development-instance" {
  name          = "tf-instance"
  instance_type = "DEVELOPMENT"

  cluster {
    cluster_id   = "tf-instance-cluster"
    zone         = "us-central1-b"
    storage_type = "HDD"
  }

  lifecycle {
    prevent_destroy = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name (also called Instance Id in the Cloud Console) of the Cloud Bigtable instance.

* `cluster` - (Required) A block of cluster configuration options. This can be specified 1 or 2 times. See structure below.

-----

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

* `instance_type` - (Optional) The instance type to create. One of `"DEVELOPMENT"` or `"PRODUCTION"`. Defaults to `"PRODUCTION"`.

* `display_name` - (Optional) The human-readable display name of the Bigtable instance. Defaults to the instance `name`.


-----

The `cluster` block supports the following arguments:

* `cluster_id` - (Required) The ID of the Cloud Bigtable cluster.

* `zone` - (Required) The zone to create the Cloud Bigtable cluster in. Each
cluster must have a different zone in the same region. Zones that support
Bigtable instances are noted on the [Cloud Bigtable locations page](https://cloud.google.com/bigtable/docs/locations).

* `num_nodes` - (Optional) The number of nodes in your Cloud Bigtable cluster.
Required, with a minimum of `1` for a `PRODUCTION` instance. Must be left unset
for a `DEVELOPMENT` instance.

* `storage_type` - (Optional) The storage type to use. One of `"SSD"` or
`"HDD"`. Defaults to `"SSD"`.

!> **Warning:** Modifying the `storage_type` or `zone` of an existing cluster (by
`cluster_id`) will cause Terraform to delete/recreate the entire
`google_bigtable_instance` resource. If these values are changing, use a new
`cluster_id`.

## Attributes Reference

Only the arguments listed above are exposed as attributes.

## Import

Bigtable Instances can be imported using any of these accepted formats:

```
$ terraform import google_bigtable_instance.default projects/{{project}}/instances/{{name}}
$ terraform import google_bigtable_instance.default {{project}}/{{name}}
$ terraform import google_bigtable_instance.default {{name}}
```
