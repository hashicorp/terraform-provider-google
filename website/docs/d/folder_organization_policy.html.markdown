---
subcategory: "Cloud Platform"
page_title: "Google: google_folder_organization_policy"
description: |-
  Retrieve Organization policies for a Google Folder
---

# google\_folder\_organization\_policy

Allows management of Organization policies for a Google Folder. For more information see
[the official
documentation](https://cloud.google.com/resource-manager/docs/organization-policy/overview)

## Example Usage

```hcl
data "google_folder_organization_policy" "policy" {
  folder     = "folders/folderid"
  constraint = "constraints/compute.trustedImageProjects"
}

output "version" {
  value = data.google_folder_organization_policy.policy.version
}
```

## Argument Reference

The following arguments are supported:

* `folder` - (Required) The resource name of the folder to set the policy for. Its format is folders/{folder_id}.

* `constraint` - (Required) (Required) The name of the Constraint the Policy is configuring, for example, `serviceuser.services`. Check out the [complete list of available constraints](https://cloud.google.com/resource-manager/docs/organization-policy/understanding-constraints#available_constraints).


## Attributes Reference

See [google_folder_organization_policy](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_folder_organization_policy) resource for details of the available attributes.
