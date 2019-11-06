---
subcategory: "Compute Engine"
layout: "google"
page_title: "Google: google_compute_project_default_network_tier"
sidebar_current: "docs-google-compute-project-default-network-tier"
description: |-
 Configures the default network tier for a project.
---

# google\_compute\_project\_default\_network\_tier

Configures the Google Compute Engine
[Default Network Tier](https://cloud.google.com/network-tiers/docs/using-network-service-tiers#setting_the_tier_for_all_resources_in_a_project)
for a project.

For more information, see,
[the Project API documentation](https://cloud.google.com/compute/docs/reference/rest/v1/projects/setDefaultNetworkTier).

## Example Usage

```hcl
resource "google_compute_project_default_network_tier" "default" {
  network_tier = "PREMIUM"
}
```

## Argument Reference

The following arguments are supported:

* `network_tier` - (Required) The default network tier to be configured for the project.
   This field can take the following values: `PREMIUM` or `STANDARD`.

- - -

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

Only the arguments listed above are exposed as attributes.

## Import

This resource can be imported using the project ID:

`terraform import google_compute_project_default_network_tier.default project-id`
