---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_projects"
sidebar_current: "docs-google-datasource-iam-testable-permissions"
description: |-
  Retrieve a list of testable permissions for a resource. Testable permissions mean the permissions that user can add or remove in a role at a given resource. The resource can be referenced either via the full resource name or via a URI.
---

# google\_iam\_testable\_permissions

Retrieve a list of testable permissions for a resource. Testable permissions mean the permissions that user can add or remove in a role at a given resource. The resource can be referenced either via the full resource name or via a URI.

## Example Usage - searching for projects about to be deleted in an org

```hcl
data "google_iam_testable_permissions" "perms" {
	full_resource_name = "//cloudresourcemanager.googleapis.com/projects/my-project"
	stages             = ["GA", "BETA"]
}
```

## Argument Reference

The following arguments are supported:

* `full_resource_name` - (Required) See [full resource name documentation](https://cloud.google.com/apis/design/resource_names#full_resource_name) for more detail.
* `stages` - (Optional) The acceptable release stages of the permission in the output. Note that `BETA` does not include permissions in `GA`, but you can specify both with `["GA", "BETA"]` for example. Can be a list of `"ALPHA"`, `"BETA"`, `"GA"`, `"DEPRECATED"`. Default is `["GA"]`.
* `custom_support_level` - (Optional) The level of support for custom roles. Can be one of `"NOT_SUPPORTED"`, `"SUPPORTED"`, `"TESTING"`. Default is `"SUPPORTED"`

## Attributes Reference

The following attributes are exported:

* `permissions` - A list of permissions matching the provided input. Structure is defined below.

The `permissions` block supports:

* `name` - Name of the permission.
* `title` - Human readable title of the permission.
* `stage` - Release stage of the permission.
* `custom_support_level` - The the support level of this permission for custom roles.
* `api_disabled` - Whether the corresponding API has been enabled for the resource.

