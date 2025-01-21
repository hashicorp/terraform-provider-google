---
subcategory: "Cloud Key Management Service"
description: |-
 Provides access to KMS key handle data with Google Cloud KMS.
---

# google_kms_key_handles

Provides access to Google Cloud Platform KMS KeyHandle. A key handle is a Cloud KMS resource that helps you safely span the separation of duties to create new Cloud KMS keys for CMEK using Autokey.

~> **Warning:** This resource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta resources.

For more information see
[the official documentation](https://cloud.google.com/kms/docs/resource-hierarchy#key_handles)
and
[API](https://cloud.google.com/kms/docs/reference/rest/v1/projects.locations.keyHandles/list).


## Example Usage

```hcl
data "google_kms_key_handles" "my_key_handles" {
  project = "resource-project-id"
  location = "us-central1"
  resource_type_selector = "storage.googleapis.com/Bucket" 
}
```

## Argument Reference

The following arguments are supported:

* `location` - (Required) The Google Cloud Platform location for the KeyHandle.
    A full list of valid locations can be found by running `gcloud kms locations list`.

* `resource_type_selector` - (Required) The resource type by which to filter KeyHandle e.g. {SERVICE}.googleapis.com/{TYPE}. See documentation for supported resource types. 

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `name` - The name of the KeyHandle. Its format is `projects/{projectId}/locations/{location}/keyHandles/{keyHandleName}`.

* `kms_key` - The identifier of the KMS Key created for the KeyHandle. Its format is `projects/{projectId}/locations/{location}/keyRings/{keyRingName}/cryptoKeys/{cryptoKeyName}`.

* `location` - The location of the KMS Key and KeyHandle.

* `project`  - The identifier of the project where KMS KeyHandle is created.

* `resource_type_selector` - Indicates the resource type that the resulting CryptoKey is meant to protect, e.g. {SERVICE}.googleapis.com/{TYPE}. See documentation for supported resource types.


