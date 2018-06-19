---
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


## Example Usage

```hcl
resource "google_bigtable_instance" "instance" {
  name         = "tf-instance"
  cluster_id   = "tf-instance-cluster"
  zone         = "us-central1-b"
  num_nodes    = 3
  storage_type = "HDD"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Cloud Bigtable instance.

* `cluster_id` - (Required) The ID of the Cloud Bigtable cluster.

* `zone` - (Required) The zone to create the Cloud Bigtable cluster in. Zones that support Bigtable instances are noted on the [Cloud Bigtable locations page](https://cloud.google.com/bigtable/docs/locations).

* `num_nodes` - (Optional) The number of nodes in your Cloud Bigtable cluster. Minimum of `3` for a `PRODUCTION` instance. Cannot be set for a `DEVELOPMENT` instance.

* `instance_type` - (Optional) The instance type to create. One of `"DEVELOPMENT"` or `"PRODUCTION"`. Defaults to `"PRODUCTION"`.

* `storage_type` - (Optional) The storage type to use. One of `"SSD"` or `"HDD"`. Defaults to `"SSD"`.

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

* `display_name` - (Optional) The human-readable display name of the Bigtable instance. Defaults to the instance `name`.

## Attributes Reference

Only the arguments listed above are exposed as attributes.
