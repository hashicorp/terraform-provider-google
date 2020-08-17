---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_project_organization_policy"
sidebar_current: "docs-google-datasource-project-organization-policy"
description: |-
  Retrieve Organization policies for a Google Project.
---

# google\_project\_organization\_policy

Allows management of Organization policies for a Google Project. For more information see
[the official
documentation](https://cloud.google.com/resource-manager/docs/organization-policy/overview)

## Example Usage

```hcl
data "google_project_organization_policy" "policy" {
  project    = "project-id"
  constraint = "constraints/serviceuser.services"
}

output "version" {
  value = data.google_project_organization_policy.policy.version
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required) The project ID.

* `constraint` - (Required) (Required) The name of the Constraint the Policy is configuring, for example, `serviceuser.services`. Check out the [complete list of available constraints](https://cloud.google.com/resource-manager/docs/organization-policy/understanding-constraints#available_constraints).


## Attributes Reference

See [google_project_organization_policy](https://www.terraform.io/docs/providers/google/r/google_project_organization_policy.html) resource for details of the available attributes.

