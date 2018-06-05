---
layout: "google"
page_title: "Google: google_project_organization_policy"
sidebar_current: "docs-google-project-organization-policy"
description: |-
 Allows management of Organization policies for a Google Project.
---

# google\_project\_organization\_policy

Allows management of Organization policies for a Google Project. For more information see
[the official
documentation](https://cloud.google.com/resource-manager/docs/organization-policy/overview) and
[API](https://cloud.google.com/resource-manager/reference/rest/v1/projects/setOrgPolicy).

## Example Usage

To set policy with a [boolean constraint](https://cloud.google.com/resource-manager/docs/organization-policy/quickstart-boolean-constraints):

```hcl
resource "google_project_organization_policy" "serial_port_policy" {
  project    = "your-project-id"
  constraint = "compute.disableSerialPortAccess"

  boolean_policy {
    enforced = true
  }
}
```


To set a policy with a [list contraint](https://cloud.google.com/resource-manager/docs/organization-policy/quickstart-list-constraints):

```hcl
resource "google_project_organization_policy" "services_policy" {
  project    = "your-project-id"
  constraint = "serviceuser.services"

  list_policy {
    allow {
      all = true
    }
  }
}
```


Or to deny some services, use the following instead:

```hcl
resource "google_project_organization_policy" "services_policy" {
  project    = "your-project-id"
  constraint = "serviceuser.services"

  list_policy {
    suggested_values = "compute.googleapis.com"

    deny {
      values = ["cloudresourcemanager.googleapis.com"]
    }
  }
}
```

To restore the default project organization policy, use the following instead:

```hcl
resource "google_project_organization_policy" "services_policy" {
  project    = "your-project-id"
  constraint = "serviceuser.services"

  restore_policy {
    default = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required) The project id of the project to set the policy for.

* `constraint` - (Required) The name of the Constraint the Policy is configuring, for example, `serviceuser.services`. Check out the [complete list of available constraints](https://cloud.google.com/resource-manager/docs/organization-policy/understanding-constraints#available_constraints).

- - -

* `version` - (Optional) Version of the Policy. Default version is 0.

* `boolean_policy` - (Optional) A boolean policy is a constraint that is either enforced or not. Structure is documented below.

* `list_policy` - (Optional) A policy that can define specific values that are allowed or denied for the given constraint. It can also be used to allow or deny all values. Structure is documented below.

* `restore_policy` - (Optional) A restore policy is a constraint to restore the default policy. Structure is documented below. 

- - -

The `boolean_policy` block supports:

* `enforced` - (Required) If true, then the Policy is enforced. If false, then any configuration is acceptable.

The `list_policy` block supports:

* `allow` or `deny` - (Optional) One or the other must be set.

* `suggested_values` - (Optional) The Google Cloud Console will try to default to a configuration that matches the value specified in this field.

The `allow` or `deny` blocks support:

* `all` - (Optional) The policy allows or denies all values.

* `values` - (Optional) The policy can define specific values that are allowed or denied.

The `restore_policy` block supports:

* `default` - (Required) May only be set to true. If set, then the default Policy is restored.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `etag` - (Computed) The etag of the organization policy. `etag` is used for optimistic concurrency control as a way to help prevent simultaneous updates of a policy from overwriting each other.

* `update_time` - (Computed) The timestamp in RFC3339 UTC "Zulu" format, accurate to nanoseconds, representing when the variable was last updated. Example: "2016-10-09T12:33:37.578138407Z".
