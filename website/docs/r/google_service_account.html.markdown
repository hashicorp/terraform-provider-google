---
layout: "google"
page_title: "Google: google_service_account"
sidebar_current: "docs-google-service-account-x"
description: |-
 Allows management of a Google Cloud Platform service account.
---

# google\_service\_account

Allows management of a [Google Cloud Platform service account](https://cloud.google.com/compute/docs/access/service-accounts)

## Example Usage

This snippet creates a service account, then gives it objectViewer
permission in a project.

```hcl
resource "google_service_account" "object_viewer" {
  account_id   = "object-viewer"
  display_name = "Object viewer"
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Required) The account id that is used to generate the service
    account email address and a stable unique id. It is unique within a project,
    must be 6-30 characters long, and match the regular expression `[a-z]([-a-z0-9]*[a-z0-9])`
    to comply with RFC1035. Changing this forces a new service account to be created.

* `display_name` - (Optional) The display name for the service account.
    Can be updated without creating a new resource.

* `project` - (Optional) The ID of the project that the service account will be created in.
    Defaults to the provider project configuration.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `email` - The e-mail address of the service account. This value
    should be referenced from any `google_iam_policy` data sources
    that would grant the service account privileges.

* `name` - The fully-qualified name of the service account.

* `unique_id` - The unique id of the service account.

## Import

Service accounts can be imported using their URI, e.g.

```
$ terraform import google_service_account.my_sa projects/my-project/serviceAccounts/my-sa@my-project.iam.gserviceaccount.com
```
