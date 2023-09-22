---
subcategory: "Certificate manager"
description: |-
  Contains the data that describes a Certificate Map
---
# google_certificate_manager_certificate_map

Get info about a Google Certificate Manager Certificate Map resource.

## Example Usage

```tf
data "google_certificate_manager_certificate_map" "default" {
 name = "cert-map"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the certificate map.

- - -
* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

See [google_certificate_manager_certificate_map](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/certificate_manager_certificate_map) resource for details of the available attributes.
