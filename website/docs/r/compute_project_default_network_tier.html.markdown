---
subcategory: "Compute Engine"
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

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `{{project}}`

## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options: configuration options:

- `create` - Default is 4 minutes (also used for update).

## Import

Compute Engine Default Network Tier can be imported using any of these accepted formats:

* `{{project_id}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Compute Engine Default Network Tier using one of the formats above. For example:

```tf
import {
  id = "{{project_id}}"
  to = google_compute_project_default_network_tier.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), Compute Engine Default Network Tier can be imported using one of the formats above. For example:

```
$ terraform import google_compute_project_default_network_tier.default {{project_id}}
```
