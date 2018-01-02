---
layout: "google"
page_title: "Google: google_service_account_key"
sidebar_current: "docs-google-service-account-key"
description: |-
  Allows management of a Google Cloud Platform service account Key Pair
---

# google\_service\_account\_key

Creates and manages service account key-pairs, which allow the user to establish identity of a service account outside of GCP. For more information, see [the official documentation](https://cloud.google.com/iam/docs/creating-managing-service-account-keys) and [API](https://cloud.google.com/iam/reference/rest/v1/projects.serviceAccounts.keys).


## Example Usage, creating a new Key Pair

```hcl
resource "google_service_account" "acceptance" {
  account_id = "%v"
  display_name = "%v"
}

resource "google_service_account_key" "acceptance" {
  service_account_id = "${google_service_account.acceptance.id}"
  public_key_type = "TYPE_X509_PEM_FILE"
}
```

## Create new Key Pair, encrypting the private key with a PGP Key

```hcl
resource "google_service_account" "acceptance" {
  account_id = "%v"
  display_name = "%v"
}

resource "google_service_account_key" "acceptance" {
  service_account_id = "${google_service_account.acceptance.id}"
  pgp_key = "keybase:keybaseusername"
  public_key_type = "TYPE_X509_PEM_FILE"
}
```

## Argument Reference

The following arguments are supported:

* `service_account_id` - (Required) The Service account id of the Key Pair.

* `key_algorithm` - (Optional) The algorithm used to generate the key. KEY_ALG_RSA_2048 is the default algorithm.
Valid values are listed at
[ServiceAccountPrivateKeyType](https://cloud.google.com/iam/reference/rest/v1/projects.serviceAccounts.keys#ServiceAccountKeyAlgorithm)
(only used on create)

* `public_key_type` (Optional) The output format of the public key requested. X509_PEM is the default output format.

* `private_key_type` (Optional) The output format of the private key. GOOGLE_CREDENTIALS_FILE is the default output format.

* `pgp_key` – (Optional) An optional PGP key to encrypt the resulting private
key material. Only used when creating or importing a new key pair. May either be
a base64-encoded public key or a `keybase:keybaseusername` string for looking up
in Vault.

~> **NOTE:** a PGP key is not required, however it is strongly encouraged.
Without a PGP key, the private key material will be stored in state unencrypted.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `name` - The name used for this key pair

* `public_key` - The public key, base64 encoded

* `private_key` - The private key, base64 encoded. This is only populated
when creating a new key, and when no `pgp_key` is provided

* `private_key_encrypted` – The private key material, base 64 encoded and
encrypted with the given `pgp_key`. This is only populated when creating a new
key and `pgp_key` is supplied

* `private_key_fingerprint` - The MD5 public key fingerprint for the encrypted
private key. This is only populated when creating a new key and `pgp_key` is supplied

* `valid_after` - The key can be used after this timestamp. A timestamp in RFC3339 UTC "Zulu" format, accurate to nanoseconds. Example: "2014-10-02T15:01:23.045123456Z".

* `valid_before` - The key can be used before this timestamp.
A timestamp in RFC3339 UTC "Zulu" format, accurate to nanoseconds. Example: "2014-10-02T15:01:23.045123456Z".

