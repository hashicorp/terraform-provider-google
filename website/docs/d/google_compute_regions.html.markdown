---
subcategory: "Compute Engine"
layout: "google"
page_title: "Google: google_compute_regions"
sidebar_current: "docs-google-datasource-compute-regions"
description: |-
  Provides a list of available Google Compute regions
---

# google\_compute\_regions

Provides access to available Google Compute regions for a given project.
See more about [regions and regions](https://cloud.google.com/compute/docs/regions-zones/) in the upstream docs.

```hcl
data "google_compute_regions" "available" {
}

resource "google_compute_subnetwork" "cluster" {
  count         = length(data.google_compute_regions.available.names)
  name          = "my-network"
  ip_cidr_range = "10.36.${count.index}.0/24"
  network       = "my-network"
  region        = data.google_compute_regions.available.names[count.index]
}
```

## Argument Reference

The following arguments are supported:

* `project` (Optional) - Project from which to list available regions. Defaults to project declared in the provider.
* `status` (Optional) - Allows to filter list of regions based on their current status. Status can be either `UP` or `DOWN`.
  Defaults to no filtering (all available regions - both `UP` and `DOWN`).

## Attributes Reference

The following attribute is exported:

* `names` - A list of regions available in the given project
