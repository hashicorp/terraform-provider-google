---
subcategory: "Cloud TPU"
layout: "google"
page_title: "Google: google_tpu_tensorflow_versions"
sidebar_current: "docs-google-datasource-tpu-tensorflow-versions"
description: |-
  Get available TensorFlow versions.
---

# google\_tpu\_tensorflow\_versions

Get TensorFlow versions available for a project. For more information see the [official documentation](https://cloud.google.com/tpu/docs/) and [API](https://cloud.google.com/tpu/docs/reference/rest/v1/projects.locations.tensorflowVersions).

## Example Usage

```hcl
data "google_tpu_tensorflow_versions" "available" {
}
```

## Example Usage: Configure Basic TPU Node with available version

```hcl
data "google_tpu_tensorflow_versions" "available" {
}

resource "google_tpu_node" "tpu" {
  name = "test-tpu"
  zone = "us-central1-b"

  accelerator_type   = "v3-8"
  tensorflow_version = data.google_tpu_tensorflow_versions.available.versions[0]
  cidr_block         = "10.2.0.0/29"
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

* `versions` - The list of TensorFlow versions available for the given project and zone.
