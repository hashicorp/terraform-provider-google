---
subcategory: "Compute Engine"
page_title: "Google: google_compute_snapshot"
description: |-
  Get information about a Google Compute Snapshot.
---

# google\_compute\_snapshot

To get more information about Snapshot, see:

* [API documentation](https://cloud.google.com/compute/docs/reference/rest/v1/snapshots)
* How-to Guides
    * [Official Documentation](https://cloud.google.com/compute/docs/disks/create-snapshots)

## Example Usage

```hcl
#by name 
data "google_compute_snapshot" "snapshot" {
  name    = "my-snapshot"
}

# using a filter
data "google_compute_snapshot" "latest-snapshot" {
  filter      = "name != my-snapshot"
  most_recent = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the compute snapshot. One of `name` or `filter` must be provided.

* `filter` - (Optional) A filter to retrieve the compute snapshot.
    See [gcloud topic filters](https://cloud.google.com/sdk/gcloud/reference/topic/filters) for reference.
    If multiple compute snapshot match, either adjust the filter or specify `most_recent`. One of `name` or `filter` must be provided.

* `most_recent` - (Optional) If `filter` is provided, ensures the most recent snapshot is returned when multiple compute snapshot match. 

- - -

* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.

## Attributes Reference

See [google_compute_snapshot](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_snapshot) resource for details of the available attributes.
