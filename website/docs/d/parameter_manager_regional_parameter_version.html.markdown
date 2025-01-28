---
subcategory: "Parameter Manager"
description: |-
  Get information about an Parameter Manager Regional Parameter Version
---

# google_parameter_manager_regional_parameter_version

Get the value and metadata from a Parameter Manager Regional Parameter version. For more information see the [official documentation](https://cloud.google.com/secret-manager/parameter-manager/docs/overview) and [API](https://cloud.google.com/secret-manager/parameter-manager/docs/reference/rest/v1/projects.locations.parameters.versions).

~> **Warning:** This datasource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta datasources.

## Example Usage

```hcl
data "google_parameter_manager_regional_parameter_version" "basic" {
  parameter            = "test-regional-parameter"
  parameter_version_id = "test-regional-parameter-version"
  location             = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The project for retrieving the Regional Parameter Version. If it's not specified, 
    the provider project will be used.

* `parameter` - (Required) The parameter for obtaining the Regional Parameter Version.
    This can be either the reference of the regional parameter as in `projects/{{project}}/locations/{{location}}/parameters/{{parameter_id}}` or only the name of the regional parameter as in `{{parameter_id}}`.

* `parameter_version_id` - (Required) The version of the regional parameter to get.

* `location` - (Optional) The location of regional parameter.


## Attributes Reference

The following attributes are exported:

* `parameter_data` - The regional parameter data.

* `name` - The resource name of the Regional Parameter Version. Format:
  `projects/{{project}}/locations/{{location}}/parameters/{{parameter_id}}/versions/{{parameter_version_id}}`

* `create_time` - The time at which the Regional Parameter Version was created.

* `update_time` - The time at which the Regional Parameter Version was last updated.

* `disabled` -  The current state of the Regional Parameter Version. 
