---
subcategory: "Compute Engine"
description: |-
  Manage asynchronous Persistent Disk replication.
---

# google_compute_disk_async_replication

Starts and stops asynchronous persistent disk replication. For more information
see [the official documentation](https://cloud.google.com/compute/docs/disks/async-pd/about)
and the [API](https://cloud.google.com/compute/docs/reference/rest/v1/disks).

## Example Usage

```tf
resource "google_compute_disk" "primary-disk" {
  name = "primary-disk"
  type = "pd-ssd"
  zone = "europe-west4-a"

  physical_block_size_bytes = 4096
}

resource "google_compute_disk" "secondary-disk" {
  name = "secondary-disk"
  type = "pd-ssd"
  zone = "europe-west3-a"

  async_primary_disk {
    disk = google_compute_disk.primary-disk.id
  }

  physical_block_size_bytes = 4096
}

resource "google_compute_disk_async_replication" "replication" {
  primary_disk = google_compute_disk.primary-disk.id
  secondary_disk {
    disk  = google_compute_disk.secondary-disk.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `primary_disk` - (Required) The primary disk (source of replication).

* `secondary_disk` - (Required) The secondary disk (target of replication). You can specify only one value. Structure is documented below.

The `secondary_disk` block includes:

* `disk` - (Required) The secondary disk.

* `state` - Output-only. Status of replication on the secondary disk.

- - -
