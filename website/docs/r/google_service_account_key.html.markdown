---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_service_account_key"
sidebar_current: "docs-google-service-account-key"
description: |-
  Allows management of a Google Cloud Platform service account Key Pair
---

# google_service_account_key

Creates and manages service account key-pairs, which allow the user to establish identity of a service account outside of GCP. For more information, see [the official documentation](https://cloud.google.com/iam/docs/creating-managing-service-account-keys) and [API](https://cloud.google.com/iam/reference/rest/v1/projects.serviceAccounts.keys).

## Example Usage, creating a new Key Pair

```hcl
resource "google_service_account" "myaccount" {
  account_id   = "myaccount"
  display_name = "My Service Account"
}

resource "google_service_account_key" "mykey" {
  service_account_id = google_service_account.myaccount.name
  public_key_type    = "TYPE_X509_PEM_FILE"
}
```

## Example Usage, creating and regularly rotating a key pair

```hcl
resource "google_service_account" "myaccount" {
  account_id   = "myaccount"
  display_name = "My Service Account"
}

# note this requires the terraform to be run regularly
resource "time_rotating" "mykey_rotation" {
  rotate_days = 30
}

resource "google_service_account_key" "mykey" {
  service_account_id = google_service_account.myaccount.name

  keepers = {
    rotation_time = time_rotating.mykey_rotation.rotation_rfc3339
  }
}
```

## Example Usage, save key in Kubernetes secret - DEPRECATED

```hcl
# Workload Identity is the recommended way of accessing Google Cloud APIs from pods.
# https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity

resource "google_service_account" "myaccount" {
  account_id   = "myaccount"
  display_name = "My Service Account"
}

resource "google_service_account_key" "mykey" {
  service_account_id = google_service_account.myaccount.name
}

resource "kubernetes_secret" "google-application-credentials" {
  metadata {
    name = "google-application-credentials"
  }
  data = {
    "credentials.json" = base64decode(google_service_account_key.mykey.private_key)
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_account_id` - (Required) The Service account id of the Key Pair. This can be a string in the format
`{ACCOUNT}` or `projects/{PROJECT_ID}/serviceAccounts/{ACCOUNT}`, where `{ACCOUNT}` is the email address or
unique id of the service account. If the `{ACCOUNT}` syntax is used, the project will be inferred from the account.

* `key_algorithm` - (Optional) The algorithm used to generate the key. KEY_ALG_RSA_2048 is the default algorithm.
Valid values are listed at
[ServiceAccountPrivateKeyType](https://cloud.google.com/iam/reference/rest/v1/projects.serviceAccounts.keys#ServiceAccountKeyAlgorithm)
(only used on create)

* `public_key_type` (Optional) The output format of the public key requested. TYPE_X509_PEM_FILE is the default output format.

* `private_key_type` (Optional) The output format of the private key. TYPE_GOOGLE_CREDENTIALS_FILE is the default output format.

* `public_key_data` (Optional) Public key data to create a service account key for given service account. The expected format for this field is a base64 encoded X509_PEM and it conflicts with `public_key_type` and `private_key_type`.

* `keepers` (Optional) Arbitrary map of values that, when changed, will trigger a new key to be generated.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - an identifier for the resource with format `projects/{{project}}/serviceAccounts/{{account}}/keys/{{key}}`

* `name` - The name used for this key pair

* `public_key` - The public key, base64 encoded

* `private_key` - The private key in JSON format, base64 encoded. This is what you normally get as a file when creating
service account keys through the CLI or web console. This is only populated when creating a new key.

* `valid_after` - The key can be used after this timestamp. A timestamp in RFC3339 UTC "Zulu" format, accurate to nanoseconds. Example: "2014-10-02T15:01:23.045123456Z".

* `valid_before` - The key can be used before this timestamp.
A timestamp in RFC3339 UTC "Zulu" format, accurate to nanoseconds. Example: "2014-10-02T15:01:23.045123456Z".

## Import

This resource does not support import.
