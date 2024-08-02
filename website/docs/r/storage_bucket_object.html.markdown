---
subcategory: "Cloud Storage"
description: |-
  Creates a new object inside a specified bucket
---

# google_storage_bucket_object

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

Example creating an empty folder in an existing `image-store` bucket.

```hcl
resource "google_storage_bucket_object" "empty_folder" {
  name   = "empty_folder/" # folder name should end with '/'
  content = " "            # content is ignored but should be non-empty
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

* `customer_encryption` - (Optional) Enables object encryption with Customer-Supplied Encryption Key (CSEK). [Google [documentation about](#nested_customer_encryption) CSEK.](https://cloud.google.com/storage/docs/encryption/customer-supplied-keys)
    Structure is [documented below](#nested_customer_encryption).

* `retention` - (Optional) The [object retention](http://cloud.google.com/storage/docs/object-lock) settings for the object. The retention settings allow an object to be retained until a provided date. Structure is [documented below](#nested_retention).

* `event_based_hold` - (Optional) Whether an object is under [event-based hold](https://cloud.google.com/storage/docs/object-holds#hold-types). Event-based hold is a way to retain objects until an event occurs, which is signified by the hold's release (i.e. this value is set to false). After being released (set to false), such objects will be subject to bucket-level retention (if any).

* `temporary_hold` - (Optional) Whether an object is under [temporary hold](https://cloud.google.com/storage/docs/object-holds#hold-types). While this flag is set to true, the object is protected against deletion and overwrites.

* `detect_md5hash` - (Optional) Detect changes to local file or changes made outside of Terraform to the file stored on the server. MD5 hash of the data, encoded using [base64](https://datatracker.ietf.org/doc/html/rfc4648#section-4). This field is not present for [composite objects](https://cloud.google.com/storage/docs/composite-objects). For more information about using the MD5 hash, see [Hashes and ETags: Best Practices](https://cloud.google.com/storage/docs/hashes-etags#json-api).

* `storage_class` - (Optional) The [StorageClass](https://cloud.google.com/storage/docs/storage-classes) of the new bucket object.
    Supported values include: `MULTI_REGIONAL`, `REGIONAL`, `NEARLINE`, `COLDLINE`, `ARCHIVE`. If not provided, this defaults to the bucket's default
    storage class or to a [standard](https://cloud.google.com/storage/docs/storage-classes#standard) class.

* `kms_key_name` - (Optional) The resource name of the Cloud KMS key that will be used to [encrypt](https://cloud.google.com/storage/docs/encryption/using-customer-managed-keys) the object.

---

<a name="nested_customer_encryption"></a>The `customer_encryption` block supports:

* `encryption_algorithm` - (Optional) Encryption algorithm. Default: AES256

* `encryption_key` - (Required) Base64 encoded Customer-Supplied Encryption Key.

<a name="nested_retention"></a>The `retention` block supports:

* `mode` - (Required) The retention policy mode. Either `Locked` or `Unlocked`.

* `retain_until_time` - (Required) The time to retain the object until in RFC 3339 format, for example 2012-11-15T16:19:00.094Z.

<a name>

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `generation` - (Computed) The content generation of this object. Used for object [versioning](https://cloud.google.com/storage/docs/object-versioning) and [soft delete](https://cloud.google.com/storage/docs/soft-delete).

* `crc32c` - (Computed) Base 64 CRC32 hash of the uploaded data.

* `md5hash` - (Computed) Base 64 MD5 hash of the uploaded data.

* `self_link` - (Computed) A url reference to this object.

* `output_name` - (Computed) The name of the object. Use this field in interpolations with `google_storage_object_acl` to recreate
`google_storage_object_acl` resources when your `google_storage_bucket_object` is recreated.

* `media_link` - (Computed) A url reference to download this object.

## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options: configuration options:

- `create` - Default is 4 minutes.
- `update` - Default is 4 minutes.
- `delete` - Default is 4 minutes.

## Import

This resource does not support import.
