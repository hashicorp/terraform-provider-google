---
subcategory: "BigQuery"
page_title: "Google: google_bigquery_default_service_account"
description: |-
  Get the email address of the project's BigQuery service account
---

# google\_bigquery\_default\_service\_account

Get the email address of a project's unique BigQuery service account.

Each Google Cloud project has a unique service account used by BigQuery. When using
BigQuery with [customer-managed encryption keys](https://cloud.google.com/bigquery/docs/customer-managed-encryption),
this account needs to be granted the
`cloudkms.cryptoKeyEncrypterDecrypter` IAM role on the customer-managed Cloud KMS key used to protect the data.

For more information see
[the API reference](https://cloud.google.com/bigquery/docs/reference/rest/v2/projects/getServiceAccount).

## Example Usage

```hcl
data "google_bigquery_default_service_account" "bq_sa" {
}

resource "google_kms_crypto_key_iam_member" "key_sa_user" {
  crypto_key_id = google_kms_crypto_key.key.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:${data.google_bigquery_default_service_account.bq_sa.email}"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The project the unique service account was created for. If it is not provided, the provider project is used.

## Attributes Reference

The following attributes are exported:

* `email` - The email address of the service account. This value is often used to refer to the service account
in order to grant IAM permissions.

* `member` - The Identity of the service account in the form `serviceAccount:{email}`. This value is often used to refer to the service account in order to grant IAM permissions.
