---
subcategory: "Cloud Key Management Service"
description: |-
 Provides access to the KMS key versions data with Google Cloud KMS.
---

# google_kms_crypto_key_versions

Provides access to Google Cloud Platform KMS CryptoKeyVersions. For more information see
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

data "google_kms_crypto_key_versions" "my_crypto_key_versions" {
  crypto_key = data.google_kms_crypto_key.my_key.id
}
```

## Argument Reference

The following arguments are supported:

* `crypto_key` - (Required) The `id` of the Google Cloud Platform CryptoKey to which the key version belongs. This is also the `id` field of the 
`google_kms_crypto_key` resource/datasource.

* `filter` - (Optional) The filter argument is used to add a filter query parameter that limits which versions are retrieved by the data source: ?filter={{filter}}. When no value is provided there is no filtering.

Example filter values if filtering on name. Note: names take the form projects/{{project}}/locations/{{location}}/keyRings/{{keyRing}}/cryptoKeys/{{cryptoKey}}/cryptoKeyVersions.

* `"name:my-key-"` will retrieve cryptoKeyVersions that contain "my-key-" anywhere in their name.
* `"name=projects/my-project/locations/global/keyRings/my-key-ring/cryptoKeys/my-key-1/cryptoKeyVersions/my-version-1"` will only retrieve a key with that exact name.

[See the documentation about using filters](https://cloud.google.com/kms/docs/sorting-and-filtering)

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `versions` - A list of all the retrieved crypto key versions from the provided crypto key. This list is influenced by the provided filter argument.

<a name="nested_public_key"></a>The `public_key` block, if present, contains:

* `pem` - The public key, encoded in PEM format. For more information, see the RFC 7468 sections for General Considerations and Textual Encoding of Subject Public Key Info.

* `algorithm` - The CryptoKeyVersionAlgorithm that this CryptoKeyVersion supports.

See [google_kms_crypto_key_version](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/kms_crypto_key_version) resource for details of the available attributes on each crypto key version.

