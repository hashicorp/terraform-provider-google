---
subcategory: "Cloud Key Management Service"
description: |-
 Provides access to data about all KMS key rings within a location with Google Cloud KMS.
---

# google_kms_crypto_key_rings

Provides access to all Google Cloud Platform KMS CryptoKeyRings in a set location. For more information see
[the official documentation](https://cloud.google.com/kms/docs/resource-hierarchy#key_rings)
and
[API](https://cloud.google.com/kms/docs/reference/rest/v1/projects.locations.keyRings).

A key ring organizes keys in a specific Google Cloud location and lets you manage access control on groups of keys. A key ring's name does not need to be unique across a Google Cloud project, but must be unique within a given location. After creation, a key ring cannot be deleted. Key rings don't incur any costs.

## Example Usage

```hcl
// Get all key rings in us-west1
data "google_kms_key_rings" "all_crypto_key_rings" {
  location = "us-west1"
}

// Get key rings from us-west1 that have "foobar" in their name
data "google_kms_key_rings" "all_crypto_key_rings" {
  location = "us-west1"
  filter   = "name:foobar"
}
```

## Argument Reference

The following arguments are supported:

* `location` - (Required) The location that the underlying key ring resides in. e.g us-west1

* `project` - (Optional) The Project ID of the project.

* `filter` - (Optional) The filter argument is used to add a filter query parameter that limits which key rings are retrieved by the data source: ?filter={{filter}}. When no value is provided there is no filtering.

Example filter values if filtering on name. Note: names take the form projects/{{project}}/locations/{{location}}/keyRings/{{keyRing}}.

* `"name:my-key-"` will retrieve key rings that contain "my-key-" anywhere in their name.
* `"name=projects/my-project/locations/global/keyRings/my-key-ring"` will only retrieve a key with that exact name.

[See the documentation about using filters](https://cloud.google.com/kms/docs/sorting-and-filtering)



## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `key_rings` - A list of all the retrieved key rings from the provided location. This list is influenced by the provided filter argument.

See [google_kms_key_ring](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/kms_key_ring) resource for details of the available attributes on each key ring.

