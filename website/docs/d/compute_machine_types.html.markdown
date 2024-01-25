---
subcategory: "Compute Engine"
description: |-
  Provides a list of available Google Compute machine types
---

# google\_compute\_machine\_types

Provides access to available Google Compute machine types in a zone for a given project.
See more about [machine type availability](https://cloud.google.com/compute/docs/regions-zones#available) in the upstream docs.

To get more information about machine types, see:

* [API Documentation](https://cloud.google.com/compute/docs/reference/rest/v1/machineTypes/list)
* [Comparison Guide](https://cloud.google.com/compute/docs/machine-resource)

## Example Usage - Machine Type properties

Configure a Google Kubernetes Engine (GKE) cluster with node auto-provisioning, using memory constraints matching the memory of the provided machine type.

```hcl
data "google_compute_machine_types" "example" {
  filter = "name = 'n1-standard-1'"
  zone   = "us-central1-a"
}

resource "google_container_cluster" "example" {
  name = "my-gke-cluster"

  cluster_autoscaling {
    enabled = true
    resource_limits {
      resource_type = "memory"
      minimum       =  2 * data.google_compute_machine_types.example.machine_types[0].memory_mb
      maximum       =  4 * data.google_compute_machine_types.example.machine_types[0].memory_mb
    }
}
```

## Example Usage - Property-based availability

Create a VM instance template for each machine type with 16GB of memory and 8 CPUs available in the provided zone.

```hcl
data "google_compute_machine_types" "example" {
  filter = "memoryMb = 16384 AND guestCpus = 8"
  zone   = var.zone
}

resource "google_compute_instance_template" "example" {
  for_each     = toset(data.google_compute_machine_types.example.machine_types[*].name)
  machine_type = each.value

  disk {
    source_image = "debian-cloud/debian-11"
    auto_delete  = true
    boot         = true
  }
}
```

## Example Usage - Machine Family preference

Create an instance template, preferring `c3` machine family if available in the provided zone, otherwise falling back to `c2` and finally `n2`.

```hcl
data "google_compute_machine_types" "example" {
  filter = "memoryMb = 16384 AND guestCpus = 4"
  zone   = var.zone
}

resource "google_compute_instance_template" "example" {
  machine_type = coalescelist(
    [for mt in data.google_compute_machine_types.example.machine_types: mt.name if startswith(mt.name, "c3-")],
    [for mt in data.google_compute_machine_types.example.machine_types: mt.name if startswith(mt.name, "c2-")],
    [for mt in data.google_compute_machine_types.example.machine_types: mt.name if startswith(mt.name, "n2-")],
  )[0]

  disk {
    source_image = "debian-cloud/debian-11"
    auto_delete  = true
    boot         = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` (Required) - A filter expression that filters machine types listed in the response.

* `zone` (Required) - Zone from which to list machine types.

* `project` (Optional) - Project from which to list available zones. Defaults to project declared in the provider.

## Attributes Reference

The following attributes are exported:

* `machine_types` - The list of machine types matching the provided filter. Structure is [documented below](#nested_machine_types).

<a name="nested_machine_types"></a>The `machine_types` block supports:

* `name` - The name of the machine type.

* `description` - A textual description of the machine type.

* `bundled_local_ssds` - ([Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html)) The configuration of bundled local SSD for the machine type. Structure is [documented below](#nested_bundled_local_ssds).

* `deprecated` - The deprecation status associated with this machine type. Structure is [documented below](#nested_deprecated).

* `guest_cpus` - The number of virtual CPUs that are available to the instance.

* `memory_mb` - The amount of physical memory available to the instance, defined in MB.

* `maximum_persistent_disks` - The maximum persistent disks allowed.

* `maximum_persistent_disks_size_gb` - The maximum total persistent disks size (GB) allowed.

* `is_shared_cpus` - Whether this machine type has a shared CPU.

* `accelerators` - A list of accelerator configurations assigned to this machine type. Structure is [documented below](#nested_accelerators).

* `self_link` - The server-defined URL for the machine type.

<a name="nested_bundled_local_ssds"></a>The `bundled_local_ssds` block supports:

* `default_interface` - ([Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html)) The default disk interface if the interface is not specified.

* `partition_count` - ([Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html)) The number of partitions.

<a name="nested_deprecated"></a>The `deprecated` block supports:

* `replacement` - The URL of the suggested replacement for a deprecated machine type.

* `state` - The deprecation state of this resource. This can be `ACTIVE`, `DEPRECATED`, `OBSOLETE`, or `DELETED`.

<a name="nested_accelerators"></a>The `accelerators` block supports:

* `guest_accelerator_type` - The accelerator type resource name, not a full URL, e.g. `nvidia-tesla-t4`.

* `guest_accelerator_count` - Number of accelerator cards exposed to the guest.
