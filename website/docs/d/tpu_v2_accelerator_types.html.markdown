---
subcategory: "Cloud TPU v2"
description: |-
  Get available accelerator types.
---

# google_tpu_v2_accelerator_types

Get accelerator types available for a project. For more information see the [official documentation](https://cloud.google.com/tpu/docs/) and [API](https://cloud.google.com/tpu/docs/reference/rest/v2/projects.locations.acceleratorTypes).

## Example Usage

```hcl
data "google_tpu_v2_accelerator_types" "available" {
}
```

## Example Usage: Configure Basic TPU VM with available type

```hcl
data "google_tpu_v2_accelerator_types" "available" {
}

data "google_tpu_v2_runtime_versions" "available" {
}

resource "google_tpu_v2_vm" "tpu" {
  name = "test-tpu"
  zone = "us-central1-b"
  runtime_version = data.google_tpu_v2_runtime_versions.available.versions[0]
  accelerator_type = data.google_tpu_v2_accelerator_types.available.types[0]
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The project to list types for. If it
    is not provided, the provider project is used.

* `zone` - (Optional) The zone to list types for. If it
    is not provided, the provider zone is used.

## Attributes Reference

The following attributes are exported:

* `types` - The list of accelerator types available for the given project and zone.
