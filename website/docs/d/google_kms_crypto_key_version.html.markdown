---
subcategory: "Cloud Key Management Service"
layout: "google"
page_title: "Google: google_kms_crypto_key_version"
sidebar_current: "docs-google-datasource-kms-crypto-key-version"
description: |-
 Provides access to KMS key version data with Google Cloud KMS.
---

# google\_kms\_crypto\_key\_version

Provides access to a Google Cloud Platform KMS CryptoKeyVersion. For more information see
[the official documentation](https://cloud.google.com/kms/docs/object-hierarchy#key_version)
and
[API](https://cloud.google.com/kms/docs/reference/rest/v1/projects.locations.keyRings.cryptoKeys.cryptoKeyVersions).

A CryptoKeyVersion represents an individual cryptographic key, and the associated key material.

## Example Usage

```hcl
data "google_kms_key_ring" "my_key_ring" {
  name     = "my-key-ring"
  location = "us-central1"
}

data "google_kms_crypto_key" "my_crypto_key" {
  name     = "my-crypto-key"
  key_ring = data.google_kms_key_ring.my_key_ring.self_link
}

data "google_kms_crypto_key_version" "my_crypto_key_version" {
  crypto_key = data.google_kms_key.my_key.self_link
}
```

## Argument Reference

The following arguments are supported:

* `crypto_key` - (Required) The `self_link` of the Google Cloud Platform CryptoKey to which the key version belongs.

* `version` - (Optional) The version number for this CryptoKeyVersion. Defaults to `1`.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `state` - The current state of the CryptoKeyVersion. See the [state reference](https://cloud.google.com/kms/docs/reference/rest/v1/projects.locations.keyRings.cryptoKeys.cryptoKeyVersions#CryptoKeyVersion.CryptoKeyVersionState) for possible outputs.

* `protection_level` - The ProtectionLevel describing how crypto operations are performed with this CryptoKeyVersion. See the [protection_level reference](https://cloud.google.com/kms/docs/reference/rest/v1/ProtectionLevel) for possible outputs.

* `algorithm` - The CryptoKeyVersionAlgorithm that this CryptoKeyVersion supports. See the [algorithm reference](https://cloud.google.com/kms/docs/reference/rest/v1/CryptoKeyVersionAlgorithm) for possible outputs.

* `public_key` -  If the enclosing CryptoKey has purpose `ASYMMETRIC_SIGN` or `ASYMMETRIC_DECRYPT`, this block contains details about the public key associated to this CryptoKeyVersion. Structure is documented below.

The `public_key` block, if present, contains:

* `pem` - The public key, encoded in PEM format. For more information, see the RFC 7468 sections for General Considerations and Textual Encoding of Subject Public Key Info.

* `algorithm` - The CryptoKeyVersionAlgorithm that this CryptoKeyVersion supports.


