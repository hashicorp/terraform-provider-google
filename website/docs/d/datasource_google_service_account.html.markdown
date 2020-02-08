---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_service_account"
sidebar_current: "docs-google-datasource-service-account"
description: |-
  Get the service account from a project.
---

# google\_service\_account

Get the service account from a project. For more information see
the official [API](https://cloud.google.com/compute/docs/access/service-accounts) documentation.

## Example Usage

```hcl
data "google_service_account" "object_viewer" {
  account_id = "object-viewer"
}
```

## Example Usage, save key in Kubernetes secret
```hcl
data "google_service_account" "myaccount" {
  account_id = "myaccount-id"
}

resource "google_service_account_key" "mykey" {
  service_account_id = data.google_service_account.myaccount.name
}

resource "kubernetes_secret" "google-application-credentials" {
  metadata {
    name = "google-application-credentials"
  }
  data = {
    credentials.json = base64decode(google_service_account_key.mykey.private_key)
  }
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Required) The Service account id.  (This is the part of the service account's email field that comes before the @ symbol.)

* `project` - (Optional) The ID of the project that the service account is present in.
    Defaults to the provider project configuration.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `email` - The e-mail address of the service account. This value
    should be referenced from any `google_iam_policy` data sources
    that would grant the service account privileges.

* `unique_id` - The unique id of the service account.

* `name` - The fully-qualified name of the service account.

* `display_name` - The display name for the service account.
