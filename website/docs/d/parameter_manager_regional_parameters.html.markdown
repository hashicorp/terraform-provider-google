---
subcategory: "Parameter Manager"
description: |-
  List the Parameter Manager Regional Parameters.
---

# google_parameter_manager_regional_parameters

Use this data source to list the Parameter Manager Regional Parameters

~> **Warning:** This datasource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta datasources.

## Example Usage

```hcl
data "google_parameter_manager_regional_parameters" "regional-parameters" {
  location = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

- `project` - (optional) The ID of the project.

- `filter` - (optional) Filter string, adhering to the rules in List-operation filtering. List only parameters matching the filter. If filter is empty, all regional parameters are listed.

- `location` - (Required) The location of regional parameter.


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

- `parameters` - A list of regional parameters matching the filter. Structure is [defined below](#nested_parameters).

<a name="nested_parameters"></a>The `parameters` block supports:

- `format` - The format type of the regional parameter.

- `labels` - The labels assigned to the regional parameter.

- `create_time` - The time at which the regional parameter was created.

- `update_time` - The time at which the regional parameter was updated.

- `project` - The ID of the project in which the resource belongs.

- `parameter_id` - The unique name of the resource.

- `name` - The resource name of the regional parameter. Format: `projects/{{project}}/locations/{{location}}/parameters/{{parameter_id}}`

- `policy_member` - An object containing a unique resource identity tied to the regional parameter. Structure is [documented below](#nested_policy_member).

<a name="nested_policy_member"></a>The `policy_member` block contains:

* `iam_policy_uid_principal` - IAM policy binding member referring to a Google Cloud resource by system-assigned unique identifier.
If a resource is deleted and recreated with the same name, the binding will not be applicable to the
new resource. Format:
`principal://parametermanager.googleapis.com/projects/{{project}}/uid/locations/{{location}}/parameters/{{uid}}`

* `iam_policy_name_principal` - AM policy binding member referring to a Google Cloud resource by user-assigned name. If a resource is deleted and recreated with the same name, the binding will be applicable to the
new resource. Format:
`principal://parametermanager.googleapis.com/projects/{{project}}/name/locations/{{location}}/parameters/{{parameter_id}}`