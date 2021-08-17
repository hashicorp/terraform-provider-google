---
subcategory: "Compute Engine"
layout: "google"
page_title: "Google: google_compute_health_check"
sidebar_current: "docs-google-datasource-compute-health-check"
description: |-
  Get information about a HealthCheck.
---

# google\_compute\_health\_check

Get information about a HealthCheck.

## Example Usage

```tf
data "google_compute_health_check" "health_check" {
  name = "my-hc"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the resource.

- - -

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

See [google_compute_health_check](https://www.terraform.io/docs/providers/google/r/compute_health_check.html) resource for details of the available attributes.
