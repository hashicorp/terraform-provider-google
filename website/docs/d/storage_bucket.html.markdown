---
subcategory: "Cloud Storage"
description: |-
  Get information about a Google Cloud Storage bucket.
---

# google\_storage\_bucket

Gets an existing bucket in Google Cloud Storage service (GCS).
See [the official documentation](https://cloud.google.com/storage/docs/key-terms#buckets)
and
[API](https://cloud.google.com/storage/docs/json_api/v1/buckets).


## Example Usage

```hcl
data "google_storage_bucket" "my-bucket" {
  name = "my-bucket"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the bucket.

* `project` - (Optional) The ID of the project in which the resource belongs. If it is not provided, the provider project is used. If no value is supplied in the configuration or through provider defaults then the data source will use the Compute API to find the project id that corresponds to the project number returned from the Storage API. Supplying a value for `project` doesn't influence retrieving data about the bucket but it can be used to prevent use of the Compute API. If you do provide a `project` value ensure that it is the correct value for that bucket; the data source will not check that the project id and project number match.

## Attributes Reference

See [google_storage_bucket](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/storage_bucket#argument-reference) resource for details of the available attributes.
