---
subcategory: "BeyondCorp"
description: |-
  Get information about a Google BeyondCorp App Connection.
---

# google_beyondcorp_app_connection

Get information about a Google BeyondCorp App Connection.

## Example Usage

```hcl
data "google_beyondcorp_app_connection" "my-beyondcorp-app-connection" {
  name = "my-beyondcorp-app-connection"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the App Connection.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

* `region` - (Optional) The region in which the resource belongs. If it
    is not provided, the provider region is used.

## Attributes Reference

See [google_beyondcorp_app_connection](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/beyondcorp_app_connection) resource for details of the available attributes.
