---
subcategory: "Cloud Key Management Service"
description: |-
 Provides access to the latest KMS key version data with Google Cloud KMS.
---

# google_kms_crypto_key_latest_version

Provides access to the latest Google Cloud Platform KMS CryptoKeyVersion in a CryptoKey. For more information see
[the official documentation](https://cloud.google.com/kms/docs/object-hierarchy#key_version)
and
[API](https://cloud.google.com/kms/docs/reference/rest/v1/projects.locations.keyRings.cryptoKeys.cryptoKeyVersions).

## Example Usage

```hcl
data "google_kms_key_ring" "my_key_ring" {
  name     = "my-key-ring"
  location = "us-central1"
}

data "google_kms_crypto_key" "my_crypto_key" {
  name     = "my-crypto-key"
  key_ring = data.google_kms_key_ring.my_key_ring.id
}

data "google_kms_crypto_key_latest_version" "my_crypto_key_latest_version" {
  crypto_key = data.google_kms_crypto_key.my_key.id
}
```

## Argument Reference

The following arguments are supported:

* `crypto_key` - (Required) The `id` of the Google Cloud Platform CryptoKey to which the key version belongs. This is also the `id` field of the 
`google_kms_crypto_key` resource/datasource.

* `filter` - (Optional) The filter argument is used to add a filter query parameter that limits which type of cryptoKeyVersion is retrieved as the latest by the data source: ?filter={{filter}}. When no value is provided there is no filtering.

Example filter values if filtering on state.

* `"state:ENABLED"` will retrieve the latest cryptoKeyVersion that has the state "ENABLED".

[See the documentation about using filters](https://cloud.google.com/kms/docs/sorting-and-filtering)

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `state` - The current state of the latest CryptoKeyVersion. See the [state reference](https://cloud.google.com/kms/docs/reference/rest/v1/projects.locations.keyRings.cryptoKeys.cryptoKeyVersions#CryptoKeyVersion.CryptoKeyVersionState) for possible outputs.

* `protection_level` - The ProtectionLevel describing how crypto operations are performed with this CryptoKeyVersion. See the [protection_level reference](https://cloud.google.com/kms/docs/reference/rest/v1/ProtectionLevel) for possible outputs.

* `algorithm` - The CryptoKeyVersionAlgorithm that this CryptoKeyVersion supports. See the [algorithm reference](https://cloud.google.com/kms/docs/reference/rest/v1/CryptoKeyVersionAlgorithm) for possible outputs.

* `public_key` -  If the enclosing CryptoKey has purpose `ASYMMETRIC_SIGN` or `ASYMMETRIC_DECRYPT`, this block contains details about the public key associated to this CryptoKeyVersion. Structure is [documented below](#nested_public_key).

<a name="nested_public_key"></a>The `public_key` block, if present, contains:

* `pem` - The public key, encoded in PEM format. For more information, see the RFC 7468 sections for General Considerations and Textual Encoding of Subject Public Key Info.

* `algorithm` - The CryptoKeyVersionAlgorithm that this CryptoKeyVersion supports.