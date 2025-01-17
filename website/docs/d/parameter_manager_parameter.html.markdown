---
subcategory: "Parameter Manager"
description: |-
  Get information about a Parameter Manager Parameter.
---

# google_parameter_manager_parameter

Use this data source to get information about a Parameter Manager Parameter.

~> **Warning:** This datasource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta datasources.

## Example Usage 

```hcl
data "google_parameter_manager_parameter" "parameter_datasource" {
  parameter_id = "foobar"
}
```

## Argument Reference

The following arguments are supported:

* `parameter_id` - (required) The name of the parameter.

* `project` - (optional) The ID of the project in which the resource belongs.

## Attributes Reference
See [google_parameter_manager_parameter](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/parameter_manager_parameter) resource for details of all the available attributes.
