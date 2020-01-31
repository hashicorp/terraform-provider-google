---
subcategory: "Cloud KMS"
layout: "google"
page_title: "Google: google_kms_secret"
sidebar_current: "docs-google-kms-secret"
description: |-
  Provides access to secret data encrypted with Google Cloud KMS
---

# google\_kms\_secret

This data source allows you to use data encrypted with Google Cloud KMS
within your resource definitions.

For more information see
[the official documentation](https://cloud.google.com/kms/docs/encrypt-decrypt).

~> **NOTE**: Using this data provider will allow you to conceal secret data within your
resource definitions, but it does not take care of protecting that data in the
logging output, plan output, or state output.  Please take care to secure your secret
data outside of resource definitions.

## Example Usage

First, create a KMS KeyRing and CryptoKey using the resource definitions:

```hcl
resource "google_kms_key_ring" "my_key_ring" {
  project  = "my-project"
  name     = "my-key-ring"
  location = "us-central1"
}

resource "google_kms_crypto_key" "my_crypto_key" {
  name     = "my-crypto-key"
  key_ring = google_kms_key_ring.my_key_ring.self_link
}
```

Next, use the [Cloud SDK](https://cloud.google.com/sdk/gcloud/reference/kms/encrypt) to encrypt some
sensitive information:

```bash
$ echo -n my-secret-password | gcloud kms encrypt \
> --project my-project \
> --location us-central1 \
> --keyring my-key-ring \
> --key my-crypto-key \
> --plaintext-file - \
> --ciphertext-file - \
> | base64
CiQAqD+xX4SXOSziF4a8JYvq4spfAuWhhYSNul33H85HnVtNQW4SOgDu2UZ46dQCRFl5MF6ekabviN8xq+F+2035ZJ85B+xTYXqNf4mZs0RJitnWWuXlYQh6axnnJYu3kDU=
```

Finally, reference the encrypted ciphertext in your resource definitions:

```hcl
data "google_kms_secret" "sql_user_password" {
  crypto_key = google_kms_crypto_key.my_crypto_key.self_link
  ciphertext = "CiQAqD+xX4SXOSziF4a8JYvq4spfAuWhhYSNul33H85HnVtNQW4SOgDu2UZ46dQCRFl5MF6ekabviN8xq+F+2035ZJ85B+xTYXqNf4mZs0RJitnWWuXlYQh6axnnJYu3kDU="
}

resource "random_id" "db_name_suffix" {
  byte_length = 4
}

resource "google_sql_database_instance" "master" {
  name = "master-instance-${random_id.db_name_suffix.hex}"

  settings {
    tier = "db-f1-micro"
  }
}

resource "google_sql_user" "users" {
  name     = "me"
  instance = google_sql_database_instance.master.name
  host     = "me.com"
  password = data.google_kms_secret.sql_user_password.plaintext
}
```

This will result in a Cloud SQL user being created with password `my-secret-password`.

## Argument Reference

The following arguments are supported:

* `ciphertext` (Required) - The ciphertext to be decrypted, encoded in base64
* `crypto_key` (Required) - The id of the CryptoKey that will be used to
  decrypt the provided ciphertext. This is represented by the format
  `{projectId}/{location}/{keyRingName}/{cryptoKeyName}`.

## Attributes Reference

The following attribute is exported:

* `plaintext` - Contains the result of decrypting the provided ciphertext.
