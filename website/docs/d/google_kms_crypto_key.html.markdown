---
subcategory: "Cloud KMS"
layout: "google"
page_title: "Google: google_kms_crypto_key"
sidebar_current: "docs-google-datasource-kms-crypto-key"
description: |-
 Provides access to KMS key data with Google Cloud KMS.
---

# google\_kms\_crypto\_key

Provides access to a Google Cloud Platform KMS CryptoKey. For more information see
[the official documentation](https://cloud.google.com/kms/docs/object-hierarchy#key)
and
[API](https://cloud.google.com/kms/docs/reference/rest/v1/projects.locations.keyRings.cryptoKeys).

A CryptoKey is an interface to key material which can be used to encrypt and decrypt data. A CryptoKey belongs to a
Google Cloud KMS KeyRing.

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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The CryptoKey's name.
    A CryptoKey’s name belonging to the specified Google Cloud Platform KeyRing and match the regular expression `[a-zA-Z0-9_-]{1,63}`

* `key_ring` - (Required) The `self_link` of the Google Cloud Platform KeyRing to which the key belongs.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `rotation_period` - Every time this period passes, generate a new CryptoKeyVersion and set it as
    the primary. The first rotation will take place after the specified period. The rotation period has the format
    of a decimal number with up to 9 fractional digits, followed by the letter s (seconds).

* `purpose` - Defines the cryptographic capabilities of the key.

* `self_link` - The self link of the created CryptoKey. Its format is `projects/{projectId}/locations/{location}/keyRings/{keyRingName}/cryptoKeys/{cryptoKeyName}`.

