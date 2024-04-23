---
subcategory: "Compute Engine"
description: |-
  Get information about a Google Compute Regional Persistent disks.
---

# google\_compute\_region\_disk

Get information about a Google Compute Regional Persistent disks.

[the official documentation](https://cloud.google.com/compute/docs/disks) and its [API](https://cloud.google.com/compute/docs/reference/rest/v1/regionDisks).

## Example Usage

```hcl
data "google_compute_region_disk" "disk" {
  name                      = "persistent-regional-disk"
  project                   = "example"
  region                    = "us-central1"
  type                      = "pd-ssd"
  physical_block_size_bytes = 4096
  
  replica_zones = ["us-central1-a", "us-central1-f"]
}

resource "google_compute_instance" "default" {
  # ...
    
  attached_disk {
    source = data.google_compute_disk.disk.self_link
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of a specific disk.

- - -

* `region` - (Optional) A reference to the region where the disk resides.

* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.

## Attributes Reference

See [google_compute_region_disk](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_region_disk) resource for details of the available attributes.