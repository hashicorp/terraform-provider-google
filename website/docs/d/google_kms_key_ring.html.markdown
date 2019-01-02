---
layout: "google"
page_title: "Google: google_kms_key_ring"
sidebar_current: "docs-google-kms-key-ring"
description: |-
  Provides read access to key rings in Google Cloud KMS
---

# google\_kms\_key\_ring

This data source allows you to query key ring by name in Google Cloud KMS
within your resource definitions.

## Example Usage

```hcl
resource "google_kms_key_ring" "main" {
  project  = "my-project"
  name     = "my-key-ring"
  location = "us-central1"
}

data "google_kms_secret" "info" {
  name = "${google_kms_key_ring.main.name}"
}

output "key_self_link" {
  value = "${data.google_kms_secret.info.self_link}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The KeyRing's name. A KeyRing’s name must be unique within a location and match the regular expression `[a-zA-Z0-9_-]{1,63}`
* `location` - (Required) The Google Cloud Platform location for the KeyRing. A full list of valid locations can be found by running `gcloud kms locations list`.
* `project` - (Optional) The project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

The following attribute is exported:

* `self_link` - The self link of the created KeyRing. Its format is `projects/{projectId}/locations/{location}/keyRings/{keyRingName}`.
