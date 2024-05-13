---
subcategory: "BeyondCorp"
description: |-
  Get information about a Google BeyondCorp App Gateway.
---

# google_beyondcorp_app_gateway

Get information about a Google BeyondCorp App Gateway.

## Example Usage

```hcl
data "google_beyondcorp_app_gateway" "my-beyondcorp-app-gateway" {
  name = "my-beyondcorp-app-gateway"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the App Gateway.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

* `region` - (Optional) The region in which the resource belongs. If it
    is not provided, the provider region is used.

## Attributes Reference

See [google_beyondcorp_app_gateway](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/beyondcorp_app_gateway) resource for details of the available attributes.
