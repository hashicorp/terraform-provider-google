---
subcategory: "Parameter Manager"
description: |-
  Get information about a Parameter Manager Regional Parameter.
---

# google_parameter_manager_regional_parameter

Use this data source to get information about a Parameter Manager Regional Parameter.

~> **Warning:** This datasource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta datasources.

## Example Usage 

```hcl
data "google_parameter_manager_regional_parameter" "reg_parameter_datasource" {
  parameter_id = "foobar"
  location = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `parameter_id` - (required) The name of the regional parameter.

* `location` - (required) The location of the regional parameter. eg us-central1

* `project` - (optional) The ID of the project in which the resource belongs.

## Attributes Reference
See [google_parameter_manager_regional_parameter](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/parameter_manager_regional_parameter) resource for details of all the available attributes.
