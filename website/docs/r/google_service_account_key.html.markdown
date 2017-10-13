---
layout: "google"
page_title: "Google: google_service_accout_key"
sidebar_current: "docs-google-service-account-key"
description: |-
  Allows management of a Google Cloud Platform service account Key Pair
---

# google\_service\_account\_key

Allows management of a key, and must be created or imported for use with
[Google Cloud Platform service account](https://cloud.google.com/compute/docs/access/service-accounts).


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
* `name` - The name used for this key pair (not used on create)

* `service_account_id` - (Required) The Serice account id of the Key Pair.

* `key_algorithm` - (Optional) The output format of the private key. GOOGLE_CREDENTIALS_FILE is the default output format. Valid values [ServiceAccountPrivateKeyType](https://cloud.google.com/iam/reference/rest/v1/projects.serviceAccounts.keys#ServiceAccountPrivateKeyType) (only used on create)

* `private_key_type` (Optional) The output format of the private key. GOOGLE_CREDENTIALS_FILE is the default output format.

* `pgp_key` – (Optional) An optional PGP key to encrypt the resulting private
key material. Only used when creating or importing a new key pair

~> **NOTE:** a PGP key is not required, however it is strongly encouraged.
Without a PGP key, the private key material will be stored in state unencrypted.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `fingerprint` - The MD5 public key fingerprint as specified in section 4 of RFC 4716.
* `public_key` - the public key, base64 encoded
* `private_key` - the private key, base64 encoded. This is only populated
when creating a new key, and when no `pgp_key` is provided
* `encrypted_private_key` – the private key material, base 64 encoded and
encrypted with the given `pgp_key`. This is only populated when creating a new
key and `pgp_key` is supplied
* `encrypted_fingerprint` - The MD5 public key fingerprint for the encrypted
private key

## Import

Lightsail Key Pairs cannot be imported, because the private and public key are
only available on initial creation.