---
subcategory: "Cloud Key Management Service"
description: |-
 Provides access to data about all KMS keys within a key ring with Google Cloud KMS.
---

# google_kms_crypto_keys

Provides access to all Google Cloud Platform KMS CryptoKeys in a given KeyRing. For more information see
[the official documentation](https://cloud.google.com/kms/docs/object-hierarchy#key)
and
[API](https://cloud.google.com/kms/docs/reference/rest/v1/projects.locations.keyRings.cryptoKeys).

A CryptoKey is an interface to key material which can be used to encrypt and decrypt data. A CryptoKey belongs to a
Google Cloud KMS KeyRing.

## Example Usage

```hcl
// Get all keys in the key ring
data "google_kms_crypto_keys" "all_crypto_keys" {
  key_ring = data.google_kms_key_ring.my_key_ring.id
}

// Get keys in the key ring that have "foobar" in their name
data "google_kms_crypto_keys" "all_crypto_keys" {
  key_ring = data.google_kms_key_ring.my_key_ring.id
  filter   = "name:foobar"
}
```

## Argument Reference

The following arguments are supported:

* `key_ring` - (Required) The key ring that the keys belongs to. Format: 'projects/{{project}}/locations/{{location}}/keyRings/{{keyRing}}'.,

* `filter` - (Optional) The filter argument is used to add a filter query parameter that limits which keys are retrieved by the data source: ?filter={{filter}}. When no value is provided there is no filtering.

Example filter values if filtering on name. Note: names take the form projects/{{project}}/locations/{{location}}/keyRings/{{keyRing}}/cryptoKeys/{{cryptoKey}}.

* `"name:my-key-"` will retrieve keys that contain "my-key-" anywhere in their name.
* `"name=projects/my-project/locations/global/keyRings/my-key-ring/cryptoKeys/my-key-1"` will only retrieve a key with that exact name.

[See the documentation about using filters](https://cloud.google.com/kms/docs/sorting-and-filtering)



## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `keys` - A list of all the retrieved keys from the provided key ring. This list is influenced by the provided filter argument.

See [google_kms_crypto_key](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/kms_crypto_key) resource for details of the available attributes on each key.

