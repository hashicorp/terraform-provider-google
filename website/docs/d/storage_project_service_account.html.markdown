---
subcategory: "Cloud Storage"
page_title: "Google: google_storage_project_service_account"
description: |-
  Get the email address of the project's Google Cloud Storage service account
---

# google\_storage\_project\_service\_account

Get the email address of a project's unique [automatic Google Cloud Storage service account](https://cloud.google.com/storage/docs/projects#service-accounts).

For each Google Cloud project, Google maintains a unique service account which
is used as the identity for various Google Cloud Storage operations, including
operations involving
[customer-managed encryption keys](https://cloud.google.com/storage/docs/encryption/customer-managed-keys)
and those involving
[storage notifications to pub/sub](https://cloud.google.com/storage/docs/gsutil/commands/notification).
This automatic Google service account requires access to the relevant Cloud KMS keys or pub/sub topics, respectively, in order for Cloud Storage to use
these customer-managed resources.

The service account has a well-known, documented naming format which is parameterised on the numeric Google project ID.
However, as noted in [the docs](https://cloud.google.com/storage/docs/projects#service-accounts), it is only created when certain relevant actions occur which
presuppose its existence.
These actions include calling a [Cloud Storage API endpoint](https://cloud.google.com/storage/docs/json_api/v1/projects/serviceAccount/get) to yield the
service account's identity, or performing some operations in the UI which must use the service account's identity, such as attempting to list Cloud KMS keys
on the bucket creation page.

Use of this data source calls the relevant API endpoint to obtain the service account's identity and thus ensures it exists prior to any API operations
which demand its existence, such as specifying it in Cloud IAM policy.
Always prefer to use this data source over interpolating the project ID into the well-known format for this service account, as the latter approach may cause
Terraform apply errors in cases where the service account does not yet exist.

>  When you write Terraform code which uses features depending on this service account *and* your Terraform code adds the service account in IAM policy on other resources,
   you must take care for race conditions between the establishment of the IAM policy and creation of the relevant Cloud Storage resource.
   Cloud Storage APIs will require permissions on resources such as pub/sub topics or Cloud KMS keys to exist *before* the attempt to utilise them in a
   bucket configuration, otherwise the API calls will fail.
   You may need to use `depends_on` to create an explicit dependency between the IAM policy resource and the Cloud Storage resource which depends on it.
   See the examples here and in the [`google_storage_notification`](/docs/providers/google/r/storage_notification.html) resource.

For more information see
[the API reference](https://cloud.google.com/storage/docs/json_api/v1/projects/serviceAccount).

## Example Usage – pub/sub notifications

```hcl
data "google_storage_project_service_account" "gcs_account" {
}

resource "google_pubsub_topic_iam_binding" "binding" {
  topic = google_pubsub_topic.topic.name
  role  = "roles/pubsub.publisher"

  members = ["serviceAccount:${data.google_storage_project_service_account.gcs_account.email_address}"]
}
```

## Example Usage – Cloud KMS keys

```hcl
data "google_storage_project_service_account" "gcs_account" {
}

resource "google_kms_crypto_key_iam_binding" "binding" {
  crypto_key_id = "your-crypto-key-id"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"

  members = ["serviceAccount:${data.google_storage_project_service_account.gcs_account.email_address}"]
}

resource "google_storage_bucket" "bucket" {
  name     = "kms-protected-bucket"
  location = "US"

  encryption {
    default_kms_key_name = "your-crypto-key-id"
  }

  # Ensure the KMS crypto-key IAM binding for the service account exists prior to the
  # bucket attempting to utilise the crypto-key.
  depends_on = [google_kms_crypto_key_iam_binding.binding]
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

* `member` - The Identity of the service account in the form `serviceAccount:{email_address}`. This value is often used to refer to the service account in order to grant IAM permissions.
