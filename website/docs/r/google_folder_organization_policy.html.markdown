---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_folder_organization_policy"
sidebar_current: "docs-google-folder-organization-policy"
description: |-
 Allows management of Organization policies for a Google Folder.
---

# google\_folder\_organization\_policy

Allows management of Organization policies for a Google Folder. For more information see
[the official
documentation](https://cloud.google.com/resource-manager/docs/organization-policy/overview) and
[API](https://cloud.google.com/resource-manager/reference/rest/v1/folders/setOrgPolicy).

## Example Usage

To set policy with a [boolean constraint](https://cloud.google.com/resource-manager/docs/organization-policy/quickstart-boolean-constraints):

```hcl
resource "google_folder_organization_policy" "serial_port_policy" {
  folder     = "folders/123456789"
  constraint = "compute.disableSerialPortAccess"

  boolean_policy {
    enforced = true
  }
}
```


To set a policy with a [list constraint](https://cloud.google.com/resource-manager/docs/organization-policy/quickstart-list-constraints):

```hcl
resource "google_folder_organization_policy" "services_policy" {
  folder     = "folders/123456789"
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
resource "google_folder_organization_policy" "services_policy" {
  folder     = "folders/123456789"
  constraint = "serviceuser.services"

  list_policy {
    suggested_value = "compute.googleapis.com"

    deny {
      values = ["cloudresourcemanager.googleapis.com"]
    }
  }
}
```

To restore the default folder organization policy, use the following instead:

```hcl
resource "google_folder_organization_policy" "services_policy" {
  folder     = "folders/123456789"
  constraint = "serviceuser.services"

  restore_policy {
    default = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `folder` - (Required) The resource name of the folder to set the policy for. Its format is folders/{folder_id}.

* `constraint` - (Required) The name of the Constraint the Policy is configuring, for example, `serviceuser.services`. Check out the [complete list of available constraints](https://cloud.google.com/resource-manager/docs/organization-policy/understanding-constraints#available_constraints).

- - -

* `version` - (Optional) Version of the Policy. Default version is 0.

* `boolean_policy` - (Optional) A boolean policy is a constraint that is either enforced or not. Structure is documented below.

* `list_policy` - (Optional) A policy that can define specific values that are allowed or denied for the given constraint. It
can also be used to allow or deny all values. Structure is documented below.

* `restore_policy` - (Optional) A restore policy is a constraint to restore the default policy. Structure is documented below.

~> **Note:** If none of [`boolean_policy`, `list_policy`, `restore_policy`] are defined the policy for a given constraint will
effectively be unset. This is represented in the UI as the constraint being 'Inherited'.

- - -

The `boolean_policy` block supports:

* `enforced` - (Required) If true, then the Policy is enforced. If false, then any configuration is acceptable.

The `list_policy` block supports:

* `allow` or `deny` - (Optional) One or the other must be set.

* `suggested_value` - (Optional) The Google Cloud Console will try to default to a configuration that matches the value specified in this field.

* `inherit_from_parent` - (Optional) If set to true, the values from the effective Policy of the parent resource
are inherited, meaning the values set in this Policy are added to the values inherited up the hierarchy.

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

## Import

Folder organization policies can be imported using any of the follow formats:

```
$ terraform import google_folder_organization_policy.policy folders/folder-1234/constraints/serviceuser.services
$ terraform import google_folder_organization_policy.policy folder-1234/serviceuser.services
```
