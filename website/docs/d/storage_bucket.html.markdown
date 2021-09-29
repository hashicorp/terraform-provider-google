---
subcategory: "Cloud Storage"
layout: "google"
page_title: "Google: google_storage_bucket"
sidebar_current: "docs-google-datasource-storage-bucket"
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

## Attributes Reference

See [google_storage_bucket](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/storage_bucket#argument-reference) resource for details of the available attributes.
