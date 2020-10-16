---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_project_default_service_accounts"
sidebar_current: "docs-google-project-default-service-accounts-x"
description: |-
  Allows management of Google Cloud Platform project default service accounts.
---

# google_project_default_service_accounts

Allows management of a Google Cloud Platform project default service accounts.

When certain Services API are enabled, Google Cloud Platform automatically creates service accounts to help to get started, but
this is not reocmended for production environment.
See the [Organization documentation](https://cloud.google.com/resource-manager/docs/quickstarts) for more details.

## Example Usage

```hcl
resource "google_project_default_service_accounts" "my_project" {
  project = "my-project-id"
  action = "delete"
}
```

To try to reactivate the default service account on the resource destroy

```hcl
resource "google_project_default_service_accounts" "my_project" {
  project = "my-project-id"
  action = "disable"
  restore_policy = "REACTIVATE"
}

```

## Argument Reference

The following arguments are supported:

- `project` - (Required) The project ID. Changing this forces a new project to be created.

- `action` - (Optional) The action to be performed in the default service accounts. Valid values are: deprivilege, delete, disable.

- `restore_policy` - (Optional) The action to be performed in the default service accounts on the resource destroy.
  Valid values are NONE and REACTIVATE. If set to REACTIVATE it will attempt to restore all default SAs.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

- `id` - an identifier for the resource with format `projects/{{project}}`

## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 10 minutes.
- `update` - Default is 10 minutes.
- `delete` - Default is 10 minutes.

## Import

This resource does not support import
