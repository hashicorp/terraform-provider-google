---
subcategory: "Privileged Access Manager"
description: |-
  Get information about a Google Cloud Privileged Access Manager Entitlement.
---

# google_privileged_access_manager_entitlement

Use this data source to get information about a Google Cloud Privileged Access Manager Entitlement.

To get more information about Privileged Access Manager, see:

* [API Documentation](https://cloud.google.com/iam/docs/reference/pam/rest)
* How-to guides
  * [Official documentation](https://cloud.google.com/iam/docs/pam-overview)

## Example Usage

```hcl
data "google_privileged_access_manager_entitlement" "my-entitlement" {
  parent  = "projects/my-project"
  location = "global"
  entitlement_id = "my-entitlement"
}
```

## Argument Reference

The following arguments are supported:

* `parent` - (Required) The project or folder or organization that contains the resource. Format: projects/{project-id|project-number} or folders/{folder-number}  or organizations/{organization-number}
* `location` - (Required) The region of the Entitlement resource.
* `entitlement_id` - (Required) ID of the Entitlement resource. This is the last part of the Entitlement's full name which is of the format `{parent}/locations/{location}/entitlements/{entitlement_id}`.

## Attribute Reference

See [google_privileged_access_manager_entitlement](https://registry.terraform.io/providers/hashicorp/google-beta/latest/docs/resources/privileged_access_manager_entitlement#argument-reference) for details of the available attributes.
