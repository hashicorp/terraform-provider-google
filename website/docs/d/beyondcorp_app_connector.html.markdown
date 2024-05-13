---
subcategory: "BeyondCorp"
description: |-
  Get information about a Google BeyondCorp App Connector.
---

# google_beyondcorp_app_connector

Get information about a Google BeyondCorp App Connector.

## Example Usage

```hcl
data "google_beyondcorp_app_connector" "my-beyondcorp-app-connector" {
  name = "my-beyondcorp-app-connector"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the App Connector.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

* `region` - (Optional) The region in which the resource belongs. If it
    is not provided, the provider region is used.

## Attributes Reference

See [google_beyondcorp_app_connector](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/beyondcorp_app_connector) resource for details of the available attributes.
