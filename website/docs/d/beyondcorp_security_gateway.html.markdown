---
subcategory: "BeyondCorp"
description: |-
  Get information about a Google BeyondCorp Security Gateway.
---

# google_beyondcorp_security_gateway

Get information about a Google BeyondCorp Security Gateway.

## Example Usage

```hcl
data "google_beyondcorp_security_gateway" "my-beyondcorp-security-gateway" {
  security_gateway_id = "my-beyondcorp-security-gateway"
}
```

## Argument Reference

The following arguments are supported:

* `security_gateway_id` - (Required) The name of the Security Gateway resource.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

See [google_beyondcorp_security_gateway](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/beyondcorp_security_gateway) resource for details of the available attributes.
