---
subcategory: "Parameter Manager"
description: |-
  List the Parameter Manager Parameters.
---

# google_parameter_manager_parameters

Use this data source to list the Parameter Manager Parameters.

~> **Warning:** This datasource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta datasources.

## Example Usage 

```hcl
data "google_parameter_manager_parameters" "parameters" {
}
```

## Argument Reference

The following arguments are supported:

* `project` - (optional) The ID of the project.

* `filter` - (optional) Filter string, adhering to the rules in List-operation filtering. List only parameters matching the filter. If filter is empty, all parameters are listed.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `parameters` - A list of parameters matching the filter. Structure is [defined below](#nested_parameters).

<a name="nested_parameters"></a>The `parameters` block supports:

* `format` - The format type of the parameter.

* `labels` - The labels assigned to the parameter.

* `create_time` - The time at which the parameter was created.

* `update_time` - The time at which the parameter was updated.

* `project` - The ID of the project in which the resource belongs.

* `parameter_id` - The unique name of the resource.

* `name` - The resource name of the parameter. Format: `projects/{{project}}/locations/global/parameters/{{parameter_id}}`

* `policy_member` - An object containing a unique resource identity tied to the parameter. Structure is [documented below](#nested_policy_member).

<a name="nested_policy_member"></a>The `policy_member` block contains:

* `iam_policy_uid_principal` - IAM policy binding member referring to a Google Cloud resource by system-assigned unique identifier.
If a resource is deleted and recreated with the same name, the binding will not be applicable to the
new resource. Format:
`principal://parametermanager.googleapis.com/projects/{{project}}/uid/locations/global/parameters/{{uid}}`

* `iam_policy_name_principal` - AM policy binding member referring to a Google Cloud resource by user-assigned name. If a resource is deleted and recreated with the same name, the binding will be applicable to the
new resource. Format:
`principal://parametermanager.googleapis.com/projects/{{project}}/name/locations/global/parameters/{{parameter_id}}`