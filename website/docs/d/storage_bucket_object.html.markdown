---
subcategory: "Cloud Storage"
description: |-
  Get information about a Google Cloud Storage bucket object.
---


# google_storage_bucket_object

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

* `generation` - (Computed) The content generation of this object. Used for object [versioning](https://cloud.google.com/storage/docs/object-versioning) and [soft delete](https://cloud.google.com/storage/docs/soft-delete).

* `crc32c` - (Computed) Base 64 CRC32 hash of the uploaded data.

* `md5hash` - (Computed) Base 64 MD5 hash of the uploaded data.

* `self_link` - (Computed) A url reference to this object.

* `storage_class` - (Computed) The [StorageClass](https://cloud.google.com/storage/docs/storage-classes) of the new bucket object.
    Supported values include: `MULTI_REGIONAL`, `REGIONAL`, `NEARLINE`, `COLDLINE`, `ARCHIVE`. If not provided, this defaults to the bucket's default
    storage class or to a [standard](https://cloud.google.com/storage/docs/storage-classes#standard) class.

* `media_link` - (Computed) A url reference to download this object.

* `event_based_hold` - (Computed) Whether an object is under [event-based hold](https://cloud.google.com/storage/docs/object-holds#hold-types). Event-based hold is a way to retain objects until an event occurs, which is signified by the hold's release (i.e. this value is set to false). After being released (set to false), such objects will be subject to bucket-level retention (if any).

* `temporary_hold` - (Computed) Whether an object is under [temporary hold](https://cloud.google.com/storage/docs/object-holds#hold-types). While this flag is set to true, the object is protected against deletion and overwrites.

* `detect_md5hash` - (Computed) Detect changes to local file or changes made outside of Terraform to the file stored on the server. MD5 hash of the data, encoded using [base64](https://datatracker.ietf.org/doc/html/rfc4648#section-4). This field is not present for [composite objects](https://cloud.google.com/storage/docs/composite-objects). For more information about using the MD5 hash, see [Hashes and ETags: Best Practices](https://cloud.google.com/storage/docs/hashes-etags#json-api).