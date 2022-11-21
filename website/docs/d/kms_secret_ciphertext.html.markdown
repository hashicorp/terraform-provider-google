---
subcategory: "Cloud Key Management Service"
page_title: "Google: google_kms_secret_ciphertext"
description: |-
  Encrypts secret data with Google Cloud KMS and provides access to the ciphertext
---

# google\_kms\_secret\_ciphertext

!> **Warning:** This data source is deprecated. Use the [`google_kms_secret_ciphertext`](../r/kms_secret_ciphertext.html) **resource** instead.

This data source allows you to encrypt data with Google Cloud KMS and use the
ciphertext within your resource definitions.

For more information see
[the official documentation](https://cloud.google.com/kms/docs/encrypt-decrypt).

~> **NOTE:** Using this data source will allow you to conceal secret data within your
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
  key_ring = google_kms_key_ring.my_key_ring.id
}
```

Next, encrypt some sensitive information and use the encrypted data in your resource definitions:

```hcl
data "google_kms_secret_ciphertext" "my_password" {
  crypto_key = google_kms_crypto_key.my_crypto_key.id
  plaintext  = "my-secret-password"
}

resource "google_compute_instance" "instance" {
  name         = "test"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }

  network_interface {
    network = "default"

    access_config {
    }
  }

  metadata = {
    password = data.google_kms_secret_ciphertext.my_password.ciphertext
  }
}
```

The resulting instance can then access the encrypted password from its metadata
and decrypt it, e.g. using the [Cloud SDK](https://cloud.google.com/sdk/gcloud/reference/kms/decrypt)):

```bash
$ curl -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/instance/attributes/password \
> | base64 -d | gcloud kms decrypt \
> --project my-project \
> --location us-central1 \
> --keyring my-key-ring \
> --key my-crypto-key \
> --plaintext-file - \
> --ciphertext-file - \
my-secret-password
```

## Argument Reference

The following arguments are supported:

* `plaintext` (Required) - The plaintext to be encrypted
* `crypto_key` (Required) - The id of the CryptoKey that will be used to
  encrypt the provided plaintext. This is represented by the format
  `{projectId}/{location}/{keyRingName}/{cryptoKeyName}`.

## Attributes Reference

The following attribute is exported:

* `ciphertext` - Contains the result of encrypting the provided plaintext, encoded in base64.

## User Project Overrides

This data source supports [User Project Overrides](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/provider_reference#user_project_override).
