---
subcategory: "Parameter Manager"
description: |-
  Get information about an Parameter Manager Regional Parameter Version with Rendered Payload Data.
---

# google_parameter_manager_regional_parameter_version_render

Get the value and metadata from a Parameter Manager Regional Parameter version with rendered payload data. For this datasource to work as expected, the principal of the parameter must be provided with the [Secret Manager Secret Accessor](https://cloud.google.com/secret-manager/docs/access-control#secretmanager.secretAccessor) role. For more information see the [official documentation](https://cloud.google.com/secret-manager/parameter-manager/docs/overview)  and [API](https://cloud.google.com/secret-manager/parameter-manager/docs/reference/rest/v1/projects.locations.parameters.versions/render).

~> **Warning:** This datasource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta datasources.

~> **Warning:** To use this data source, we must grant the `Secret Manager Secret Accessor` role to the principal of the parameter. Please note that it can take up to 7 minutes for the role to take effect. Hence, we might need to wait approximately 7 minutes after granting  `Secret Manager Secret Accessor` role to the principal of the parameter. For more information see the [access change propagation documentation](https://cloud.google.com/iam/docs/access-change-propagation).

## Example Usage

```hcl
data "google_parameter_manager_regional_parameter_version_render" "basic" {
  parameter            = "test-regional-parameter"
  parameter_version_id = "test-regional-parameter-version"
  location             = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The project for retrieving the Regional Parameter Version. If it's not
    specified, the provider project will be used.

* `parameter` - (Required) The Parameter for obtaining the Regional Parameter Version.
    This can be either the reference of the parameter as in `projects/{{project}}/locations/{{location}}/parameters/{{parameter_id}}` or only the name of the parameter as in `{{parameter_id}}`.

* `location` - (Optional) Location of Parameter Manager regional Parameter resource.
    It must be provided when the `parameter` field provided consists of only the name of the regional parameter.

* `parameter_version_id` - (Required) The version of the regional parameter to get.

## Attributes Reference

The following attributes are exported:

* `parameter_data` - The Parameter data.

* `rendered_parameter_data` - The Rendered Parameter Data specifies that if you use `__REF__()` to reference a secret and the format is JSON or YAML, the placeholder `__REF__()` will be replaced with the actual secret value. However, if the format is UNFORMATTED, it will stay the same as the original `parameter_data`.

* `name` - The resource name of the RegionalParameterVersion. Format:
  `projects/{{project}}/locations/{{location}}/parameters/{{parameter_id}}/versions/{{parameter_version_id}}`

* `disabled` -  The current state of the Regional Parameter Version. 
