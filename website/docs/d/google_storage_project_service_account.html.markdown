---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_storage_project_service_account"
sidebar_current: "docs-google-datasource-storage-project-service-account"
description: |-
  Get the email address of the project's Google Cloud Storage service account
---

# google\_storage\_project\_service\_account

Get the email address of a project's unique Google Cloud Storage service account.

Each Google Cloud project has a unique service account for use with Google Cloud Storage. Only this
special service account can be used to set up `google_storage_notification` resources.

For more information see
[the API reference](https://cloud.google.com/storage/docs/json_api/v1/projects/serviceAccount).

## Example Usage

```hcl
data "google_storage_project_service_account" "gcs_account" {
}

resource "google_pubsub_topic_iam_binding" "binding" {
  topic = google_pubsub_topic.topic.name
  role  = "roles/pubsub.publisher"

  members = ["serviceAccount:${data.google_storage_project_service_account.gcs_account.email_address}"]
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The project the unique service account was created for. If it is not provided, the provider project is used.

* `user_project` - (Optional) The project the lookup originates from. This field is used if you are making the request
from a different account than the one you are finding the service account for.

## Attributes Reference

The following attributes are exported:

* `email_address` - The email address of the service account. This value is often used to refer to the service account
in order to grant IAM permissions.
