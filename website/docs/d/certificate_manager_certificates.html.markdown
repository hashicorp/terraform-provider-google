---
subcategory: "Certificate manager"
description: |-
  List all certificates within a project and region.
---
# google_certificate_manager_certificates

List all certificates within Google Certificate Manager for a given project, region or filter.

## Example Usage

```tf
data "google_certificate_manager_certificates" "default" {
}
```

## Example Usage - with a filter

```tf
data "google_certificate_manager_certificates" "default" {
  filter = "name:projects/PROJECT_ID/locations/REGION/certificates/certificate-name-*"
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) Filter expression to restrict the certificates returned.
* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.
* `region` - (Optional) The region in which the resource belongs. If it is not provided, `GLOBAL` is used.

## Attributes Reference

See [google_certificate_manager_certificate](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/certificate_manager_certificate) resource for details of the available attributes.
