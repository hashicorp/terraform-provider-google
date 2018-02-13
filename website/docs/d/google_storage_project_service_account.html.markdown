---
layout: "google"
page_title: "Google: google_storage_project_service_account"
sidebar_current: "docs-google-datasource-storage-project-service-account"
description: |-
  Get the email address of the project's Google Cloud Storage service account
---

# google\_storage\_project\_service\_account

Use this data source to get the email address of the project's Google Cloud Storage service account.
 For more information see 
[API](https://cloud.google.com/storage/docs/json_api/v1/projects/serviceAccount).

## Example Usage

```hcl
data "google_storage_project_service_account" "gcs_account" {}

resource "google_pubsub_topic_iam_binding" "binding" {
	topic       = "${google_pubsub_topic.topic.name}"
	role        = "roles/pubsub.publisher"
		  
	members     = ["${data.google_storage_project_service_account.gcs_account.id}"]
}
```

## Argument Reference

There are no arguments available for this data source.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service account, which is its email_address

* `email_address` - The email_address for this account
