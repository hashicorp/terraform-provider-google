---
subcategory: "Cloud Key Management Service"
layout: "google"
page_title: "Google: google_kms_key_ring"
sidebar_current: "docs-google-datasource-kms-key-ring"
description: |-
 Provides access to KMS key ring data with Google Cloud KMS.
---

# google\_kms\_key\_ring

Provides access to Google Cloud Platform KMS KeyRing. For more information see
[the official documentation](https://cloud.google.com/kms/docs/object-hierarchy#key_ring)
and
[API](https://cloud.google.com/kms/docs/reference/rest/v1/projects.locations.keyRings).

A KeyRing is a grouping of CryptoKeys for organizational purposes. A KeyRing belongs to a Google Cloud Platform Project
and resides in a specific location.

## Example Usage

```hcl
data "google_kms_key_ring" "my_key_ring" {
  name     = "my-key-ring"
  location = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The KeyRing's name.
    A KeyRing name must exist within the provided location and match the regular expression `[a-zA-Z0-9_-]{1,63}`

* `location` - (Required) The Google Cloud Platform location for the KeyRing.
    A full list of valid locations can be found by running `gcloud kms locations list`.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `self_link` - The self link of the created KeyRing. Its format is `projects/{projectId}/locations/{location}/keyRings/{keyRingName}`.
