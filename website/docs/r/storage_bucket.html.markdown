---
subcategory: "Cloud Storage"
layout: "google"
page_title: "Google: google_storage_bucket"
sidebar_current: "docs-google-storage-bucket-x"
description: |-
  Creates a new bucket in Google Cloud Storage.
---

# google\_storage\_bucket

Creates a new bucket in Google cloud storage service (GCS).
Once a bucket has been created, its location can't be changed.
[ACLs](https://cloud.google.com/storage/docs/access-control/lists) can be applied
using the [`google_storage_bucket_acl` resource](/docs/providers/google/r/storage_bucket_acl.html).

For more information see
[the official documentation](https://cloud.google.com/storage/docs/overview)
and
[API](https://cloud.google.com/storage/docs/json_api/v1/buckets).

**Note**: If the project id is not set on the resource or in the provider block it will be dynamically
determined which will require enabling the compute api.


## Example Usage

Example creating a private bucket in standard storage, in the EU region.

```hcl
resource "google_storage_bucket" "image-store" {
  name     = "image-store-bucket"
  location = "EU"

  website {
    main_page_suffix = "index.html"
    not_found_page   = "404.html"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the bucket.

- - -

* `force_destroy` - (Optional, Default: false) When deleting a bucket, this
    boolean option will delete all contained objects. If you try to delete a
    bucket that contains objects, Terraform will fail that run.

* `location` - (Optional, Default: 'US') The [GCS location](https://cloud.google.com/storage/docs/bucket-locations)

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

* `storage_class` - (Optional, Default: 'STANDARD') The [Storage Class](https://cloud.google.com/storage/docs/storage-classes) of the new bucket. Supported values include: `STANDARD`, `MULTI_REGIONAL`, `REGIONAL`, `NEARLINE`, `COLDLINE`.

* `lifecycle_rule` - (Optional) The bucket's [Lifecycle Rules](https://cloud.google.com/storage/docs/lifecycle#configuration) configuration. Multiple blocks of this type are permitted. Structure is documented below.

* `versioning` - (Optional) The bucket's [Versioning](https://cloud.google.com/storage/docs/object-versioning) configuration.

* `website` - (Optional) Configuration if the bucket acts as a website. Structure is documented below.

* `cors` - (Optional) The bucket's [Cross-Origin Resource Sharing (CORS)](https://www.w3.org/TR/cors/) configuration. Multiple blocks of this type are permitted. Structure is documented below.

* `retention_policy` - (Optional) Configuration of the bucket's data retention policy for how long objects in the bucket should be retained. Structure is documented below.

* `labels` - (Optional) A set of key/value label pairs to assign to the bucket.

* `logging` - (Optional) The bucket's [Access & Storage Logs](https://cloud.google.com/storage/docs/access-logs) configuration.

* `encryption` - (Optional) The bucket's encryption configuration.

* `requester_pays` - (Optional, Default: false) Enables [Requester Pays](https://cloud.google.com/storage/docs/requester-pays) on a storage bucket.

* `bucket_policy_only` - (Optional, Default: false) Enables [Bucket Policy Only](https://cloud.google.com/storage/docs/bucket-policy-only) access to a bucket.

The `lifecycle_rule` block supports:

* `action` - (Required) The Lifecycle Rule's action configuration. A single block of this type is supported. Structure is documented below.

* `condition` - (Required) The Lifecycle Rule's condition configuration. A single block of this type is supported. Structure is documented below.

The `action` block supports:

* `type` - The type of the action of this Lifecycle Rule. Supported values include: `Delete` and `SetStorageClass`.

* `storage_class` - (Required if action type is `SetStorageClass`) The target [Storage Class](https://cloud.google.com/storage/docs/storage-classes) of objects affected by this Lifecycle Rule. Supported values include: `MULTI_REGIONAL`, `REGIONAL`, `NEARLINE`, `COLDLINE`.

The `condition` block supports the following elements, and requires at least one to be defined:

* `age` - (Optional) Minimum age of an object in days to satisfy this condition.

* `created_before` - (Optional) Creation date of an object in RFC 3339 (e.g. `2017-06-13`) to satisfy this condition.

* `with_state` - (Optional) Match to live and/or archived objects. Unversioned buckets have only live objects. Supported values include: `"LIVE"`, `"ARCHIVED"`, `"ANY"`.

* `matches_storage_class` - (Optional) [Storage Class](https://cloud.google.com/storage/docs/storage-classes) of objects to satisfy this condition. Supported values include: `MULTI_REGIONAL`, `REGIONAL`, `NEARLINE`, `COLDLINE`, `STANDARD`, `DURABLE_REDUCED_AVAILABILITY`.

* `num_newer_versions` - (Optional) Relevant only for versioned objects. The number of newer versions of an object to satisfy this condition.

The `versioning` block supports:

* `enabled` - (Required) While set to `true`, versioning is fully enabled for this bucket.

The `website` block supports:

* `main_page_suffix` - (Optional) Behaves as the bucket's directory index where
    missing objects are treated as potential directories.

* `not_found_page` - (Optional) The custom object to return when a requested
    resource is not found.

The `cors` block supports:

* `origin` - (Optional) The list of [Origins](https://tools.ietf.org/html/rfc6454) eligible to receive CORS response headers. Note: "*" is permitted in the list of origins, and means "any Origin".

* `method` - (Optional) The list of HTTP methods on which to include CORS response headers, (GET, OPTIONS, POST, etc) Note: "*" is permitted in the list of methods, and means "any method".

* `response_header` - (Optional) The list of HTTP headers other than the [simple response headers](https://www.w3.org/TR/cors/#simple-response-header) to give permission for the user-agent to share across domains.

* `max_age_seconds` - (Optional) The value, in seconds, to return in the [Access-Control-Max-Age header](https://www.w3.org/TR/cors/#access-control-max-age-response-header) used in preflight responses.

The `retention_policy` block supports:

* `is_locked` - (Optional) If set to `true`, the bucket will be [locked](https://cloud.google.com/storage/docs/using-bucket-lock#lock-bucket) and permanently restrict edits to the bucket's retention policy.  Caution: Locking a bucket is an irreversible action.

* `retention_period` - (Optional) The period of time, in seconds, that objects in the bucket must be retained and cannot be deleted, overwritten, or archived. The value must be less than 3,155,760,000 seconds.

The `logging` block supports:

* `log_bucket` - (Required) The bucket that will receive log objects.

* `log_object_prefix` - (Optional, Computed) The object prefix for log objects. If it's not provided,
    by default GCS sets this to this bucket's name.

The `encryption` block supports:

* `default_kms_key_name`: A Cloud KMS key that will be used to encrypt objects inserted into this bucket, if no encryption method is specified.
  You must pay attention to whether the crypto key is available in the location that this bucket is created in.
  See [the docs](https://cloud.google.com/storage/docs/encryption/using-customer-managed-keys) for more details.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `self_link` - The URI of the created resource.

* `url` - The base URL of the bucket, in the format `gs://<bucket-name>`.

## Import

Storage buckets can be imported using the `name` or  `project/name`. If the project is not
passed to the import command it will be inferred from the provider block or environment variables.
If it cannot be inferred it will be queried from the Compute API (this will fail if the API is
not enabled).

e.g.

```
$ terraform import google_storage_bucket.image-store image-store-bucket
$ terraform import google_storage_bucket.image-store tf-test-project/image-store-bucket
```

~> **Note:** Terraform will import this resource with `force_destroy` set to
`false` in state. If you've set it to `true` in config, run `terraform apply` to
update the value set in state. If you delete this resource before updating the
value, objects in the bucket will not be destroyed.

