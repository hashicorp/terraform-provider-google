---
subcategory: "Access Context Manager (VPC Service Controls)"
description: |-
  Fetches an AccessPolicy from Access Context Manager.
---

# google_access_context_manager_access_policy

Get information about an Access Context Manager AccessPolicy.

## Example Usage

```tf
data "google_access_context_manager_access_policy" "policy-org" {
  parent       = "organizations/1234567"
}

data "google_access_context_manager_access_policy" "policy-scoped" {
  parent       = "organizations/1234567"
  scopes       = ["projects/1234567"]
}

```

## Argument Reference

The following arguments are supported:

* `parent` - (Required) The parent of this AccessPolicy in the Cloud Resource Hierarchy. Format: `organizations/{{organization_id}}`

* `scopes` - (Optional) Folder or project on which this policy is applicable. Format: `folders/{{folder_id}}` or `projects/{{project_number}}`


## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `name` - Resource name of the AccessPolicy.

* `title` - Human readable title. Does not affect behavior.
