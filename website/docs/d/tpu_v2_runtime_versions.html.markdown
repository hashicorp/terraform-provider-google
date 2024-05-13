---
subcategory: "Cloud TPU v2"
description: |-
  Get available runtime versions.
---

# google_tpu_v2_runtime_versions

Get runtime versions available for a project. For more information see the [official documentation](https://cloud.google.com/tpu/docs/) and [API](https://cloud.google.com/tpu/docs/reference/rest/v2/projects.locations.runtimeVersions).

## Example Usage

```hcl
data "google_tpu_v2_runtime_versions" "available" {
}
```

## Example Usage: Configure Basic TPU VM with available version

```hcl
data "google_tpu_v2_runtime_versions" "available" {
}

resource "google_tpu_v2_vm" "tpu" {
  name = "test-tpu"
  zone = "us-central1-b"
  runtime_version = data.google_tpu_v2_runtime_versions.available.versions[0]
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The project to list versions for. If it
    is not provided, the provider project is used.

* `zone` - (Optional) The zone to list versions for. If it
    is not provided, the provider zone is used.

## Attributes Reference

The following attributes are exported:

* `versions` - The list of runtime versions available for the given project and zone.
