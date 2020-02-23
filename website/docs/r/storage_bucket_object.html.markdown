---
subcategory: "Cloud Storage"
layout: "google"
page_title: "Google: google_storage_bucket_object"
sidebar_current: "docs-google-storage-bucket-object"
description: |-
  Creates a new object inside a specified bucket
---

# google\_storage\_bucket\_object

Creates a new object inside an existing bucket in Google cloud storage service (GCS). 
[ACLs](https://cloud.google.com/storage/docs/access-control/lists) can be applied using the `google_storage_object_acl` resource.
 For more information see 
[the official documentation](https://cloud.google.com/storage/docs/key-terms#objects) 
and 
[API](https://cloud.google.com/storage/docs/json_api/v1/objects).


## Example Usage

Example creating a public object in an existing `image-store` bucket.

```hcl
resource "google_storage_bucket_object" "picture" {
  name   = "butterfly01"
  source = "/images/nature/garden-tiger-moth.jpg"
  bucket = "image-store"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) The name of the containing bucket.

* `name` - (Required) The name of the object. If you're interpolating the name of this object, see `output_name` instead.

* `metadata` - (Optional) User-provided metadata, in key/value pairs.

One of the following is required:

* `content` - (Optional, Sensitive) Data as `string` to be uploaded. Must be defined if `source` is not. **Note**: The `content` field is marked as sensitive. To view the raw contents of the object, please define an [output](/docs/configuration/outputs.html).

* `source` - (Optional) A path to the data you want to upload. Must be defined
    if `content` is not.

- - -

* `cache_control` - (Optional) [Cache-Control](https://tools.ietf.org/html/rfc7234#section-5.2)
    directive to specify caching behavior of object data. If omitted and object is accessible to all anonymous users, the default will be public, max-age=3600

* `content_disposition` - (Optional) [Content-Disposition](https://tools.ietf.org/html/rfc6266) of the object data.

* `content_encoding` - (Optional) [Content-Encoding](https://tools.ietf.org/html/rfc7231#section-3.1.2.2) of the object data.

* `content_language` - (Optional) [Content-Language](https://tools.ietf.org/html/rfc7231#section-3.1.3.2) of the object data.

* `content_type` - (Optional) [Content-Type](https://tools.ietf.org/html/rfc7231#section-3.1.1.5) of the object data. Defaults to "application/octet-stream" or "text/plain; charset=utf-8".

* `storage_class` - (Optional) The [StorageClass](https://cloud.google.com/storage/docs/storage-classes) of the new bucket object.
    Supported values include: `MULTI_REGIONAL`, `REGIONAL`, `NEARLINE`, `COLDLINE`. If not provided, this defaults to the bucket's default
    storage class or to a [standard](https://cloud.google.com/storage/docs/storage-classes#standard) class.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `crc32c` - (Computed) Base 64 CRC32 hash of the uploaded data.

* `md5hash` - (Computed) Base 64 MD5 hash of the uploaded data.

* `self_link` - (Computed) A url reference to this object.

* `output_name` - (Computed) The name of the object. Use this field in interpolations with `google_storage_object_acl` to recreate
`google_storage_object_acl` resources when your `google_storage_bucket_object` is recreated.
