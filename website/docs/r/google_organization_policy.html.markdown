---
layout: "google"
page_title: "Google: google_organization_policy"
sidebar_current: "docs-google-organization-policy"
description: |-
 Allows management of Organization policies for a Google Organization.
---

# google\_organization\_policy

Allows management of Organization policies for a Google Organization. For more information see
[the official
documentation](https://cloud.google.com/resource-manager/docs/organization-policy/overview) and
[API](https://cloud.google.com/resource-manager/reference/rest/v1/organizations/setOrgPolicy).

## Example Usage

To set policy with a [boolean constraint](https://cloud.google.com/resource-manager/docs/organization-policy/quickstart-boolean-constraints):

```hcl
resource "google_folder_organization_policy" "serial_port_policy" {
  org_id     = "123456789"
  constraint = "compute.disableSerialPortAccess"

  boolean_policy {
    enforced = true
  }
}
```


To set a policy with a [list contraint](https://cloud.google.com/resource-manager/docs/organization-policy/quickstart-list-constraints):

```hcl
resource "google_folder_organization_policy" "services_policy" {
  org_id     = "123456789"
  constraint = "serviceuser.services"

  list_policy {
    deny {
      all = true
    }
  }
}
```


Or to allow just some services, use the following instead:

```hcl
resource "google_folder_organization_policy" "services_policy" {
  org_id     = "123456789"
  constraint = "serviceuser.services"

  list_policy {
    allow {
      values = ["cloudresourcemanager.googleapis.com", "compute.googleapis.com"]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `org_id` - (Required) The numeric ID of the organization to set the policy for.

* `constraint` - (Required) The name of the Constraint the Policy is configuring, for example, `serviceuser.services`. Check out the [complete list of available constraints](https://cloud.google.com/resource-manager/docs/organization-policy/understanding-constraints#available_constraints).

* `version` - (Optional) Version of the Policy. Default version is 0.

TODO(rosbo): Add boolean and list policy block once I get agreement on the schema format.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `etag` - (Computed) The etag of the organization policy. `etag` is used for optimistic concurrency control as a way to help prevent simultaneous updates of a policy from overwriting each other. 

* `update_time` - (Computed) The timestamp in RFC3339 UTC "Zulu" format, accurate to nanoseconds, representing when the variable was last updated. Example: "2016-10-09T12:33:37.578138407Z".