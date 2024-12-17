---
subcategory: "Cloud Key Management Service"
description: |-
 Provides access to KMS key handle data with Google Cloud KMS.
---

# google_kms_key_handle

Provides access to Google Cloud Platform KMS KeyHandle. For more information see
[the official documentation](https://cloud.google.com/kms/docs/resource-hierarchy#key_handles)
and
[API](https://cloud.google.com/kms/docs/reference/rest/v1/projects.locations.keyHandles).

A key handle is a Cloud KMS resource that helps you safely span the separation of duties to create new Cloud KMS keys for CMEK using Autokey.

## Example Usage

```hcl
data "google_kms_key_handle" "my_key_handle" {
  name     = "eed58b7b-20ad-4da8-ad85-ba78a0d5ab87"
  location = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The KeyHandle's name.
    A KeyHandle name must exist within the provided location and must be valid UUID.

* `location` - (Required) The Google Cloud Platform location for the KeyHandle.
    A full list of valid locations can be found by running `gcloud kms locations list`.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `id` - The identifier of the created KeyHandle. Its format is `projects/{projectId}/locations/{location}/keyHandles/{keyHandleName}`.
