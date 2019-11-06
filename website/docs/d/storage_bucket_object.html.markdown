---
subcategory: "Cloud Storage"
layout: "google"
page_title: "Google: google_storage_bucket_object"
sidebar_current: "docs-google-datasource-storage-bucket-object"
description: |-
  Get information about a Google Cloud Storage bucket object.
---


# google\_storage\_bucket\_object

Gets an existing object inside an existing bucket in Google Cloud Storage service (GCS).
See [the official documentation](https://cloud.google.com/storage/docs/key-terms#objects)
and
[API](https://cloud.google.com/storage/docs/json_api/v1/objects).


## Example Usage

Example picture stored within a folder.

```hcl
data "google_storage_bucket_object" "picture" {
  name   = "folder/butterfly01.jpg"
  bucket = "image-store"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) The name of the containing bucket.

* `name` - (Required) The name of the object.

## Attributes Reference

The following attributes are exported:

* `cache_control` - (Computed) [Cache-Control](https://tools.ietf.org/html/rfc7234#section-5.2)
    directive to specify caching behavior of object data. If omitted and object is accessible to all anonymous users, the default will be public, max-age=3600

* `content_disposition` - (Computed) [Content-Disposition](https://tools.ietf.org/html/rfc6266) of the object data.

* `content_encoding` - (Computed) [Content-Encoding](https://tools.ietf.org/html/rfc7231#section-3.1.2.2) of the object data.

* `content_language` - (Computed) [Content-Language](https://tools.ietf.org/html/rfc7231#section-3.1.3.2) of the object data.

* `content_type` - (Computed) [Content-Type](https://tools.ietf.org/html/rfc7231#section-3.1.1.5) of the object data. Defaults to "application/octet-stream" or "text/plain; charset=utf-8".

* `crc32c` - (Computed) Base 64 CRC32 hash of the uploaded data.

* `md5hash` - (Computed) Base 64 MD5 hash of the uploaded data.

* `self_link` - (Computed) A url reference to this object.

* `storage_class` - (Computed) The [StorageClass](https://cloud.google.com/storage/docs/storage-classes) of the new bucket object.
    Supported values include: `MULTI_REGIONAL`, `REGIONAL`, `NEARLINE`, `COLDLINE`. If not provided, this defaults to the bucket's default
    storage class or to a [standard](https://cloud.google.com/storage/docs/storage-classes#standard) class.
