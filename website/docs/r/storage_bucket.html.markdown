---
subcategory: "Cloud Storage"
page_title: "Google: google_storage_bucket"
description: |-
  Creates a new bucket in Google Cloud Storage.
---

# google\_storage\_bucket

Creates a new bucket in Google cloud storage service (GCS).
Once a bucket has been created, its location can't be changed.

For more information see
[the official documentation](https://cloud.google.com/storage/docs/overview)
and
[API](https://cloud.google.com/storage/docs/json_api/v1/buckets).

**Note**: If the project id is not set on the resource or in the provider block it will be dynamically
determined which will require enabling the compute api.


## Example Usage - creating a private bucket in standard storage, in the EU region. Bucket configured as static website and CORS configurations

```hcl
resource "google_storage_bucket" "static-site" {
  name          = "image-store.com"
  location      = "EU"
  force_destroy = true

  uniform_bucket_level_access = true

  website {
    main_page_suffix = "index.html"
    not_found_page   = "404.html"
  }
  cors {
    origin          = ["http://image-store.com"]
    method          = ["GET", "HEAD", "PUT", "POST", "DELETE"]
    response_header = ["*"]
    max_age_seconds = 3600
  }
}
```

## Example Usage - Life cycle settings for storage bucket objects

```hcl
resource "google_storage_bucket" "auto-expire" {
  name          = "auto-expiring-bucket"
  location      = "US"
  force_destroy = true

  lifecycle_rule {
    condition {
      age = 3
    }
    action {
      type = "Delete"
    }
  }

  lifecycle_rule {
    condition {
      age = 1
    }
    action {
      type = "AbortIncompleteMultipartUpload"
    }
  }
}
```

## Example Usage - Enabling public access prevention

```hcl
resource "google_storage_bucket" "auto-expire" {
  name          = "no-public-access-bucket"
  location      = "US"
  force_destroy = true

  public_access_prevention = "enforced"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the bucket.

* `location` - (Required) The [GCS location](https://cloud.google.com/storage/docs/bucket-locations).

- - -

* `force_destroy` - (Optional, Default: false) When deleting a bucket, this
    boolean option will delete all contained objects. If you try to delete a
    bucket that contains objects, Terraform will fail that run.

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

* `storage_class` - (Optional, Default: 'STANDARD') The [Storage Class](https://cloud.google.com/storage/docs/storage-classes) of the new bucket. Supported values include: `STANDARD`, `MULTI_REGIONAL`, `REGIONAL`, `NEARLINE`, `COLDLINE`, `ARCHIVE`.

* `lifecycle_rule` - (Optional) The bucket's [Lifecycle Rules](https://cloud.google.com/storage/docs/lifecycle#configuration) configuration. Multiple blocks of this type are permitted. Structure is [documented below](#nested_lifecycle_rule).

* `versioning` - (Optional) The bucket's [Versioning](https://cloud.google.com/storage/docs/object-versioning) configuration.  Structure is [documented below](#nested_versioning).

* `website` - (Optional) Configuration if the bucket acts as a website. Structure is [documented below](#nested_website).

* `cors` - (Optional) The bucket's [Cross-Origin Resource Sharing (CORS)](https://www.w3.org/TR/cors/) configuration. Multiple blocks of this type are permitted. Structure is [documented below](#nested_cors).

* `default_event_based_hold` - (Optional) Whether or not to automatically apply an eventBasedHold to new objects added to the bucket.

* `retention_policy` - (Optional) Configuration of the bucket's data retention policy for how long objects in the bucket should be retained. Structure is [documented below](#nested_retention_policy).

* `labels` - (Optional) A map of key/value label pairs to assign to the bucket.

* `logging` - (Optional) The bucket's [Access & Storage Logs](https://cloud.google.com/storage/docs/access-logs) configuration. Structure is [documented below](#nested_logging).

* `encryption` - (Optional) The bucket's encryption configuration. Structure is [documented below](#nested_encryption).

* `requester_pays` - (Optional, Default: false) Enables [Requester Pays](https://cloud.google.com/storage/docs/requester-pays) on a storage bucket.

* `uniform_bucket_level_access` - (Optional, Default: false) Enables [Uniform bucket-level access](https://cloud.google.com/storage/docs/uniform-bucket-level-access) access to a bucket.

* `public_access_prevention` - (Optional) Prevents public access to a bucket. Acceptable values are "inherited" or "enforced". If "inherited", the bucket uses [public access prevention](https://cloud.google.com/storage/docs/public-access-prevention). only if the bucket is subject to the public access prevention organization policy constraint. Defaults to "inherited".

* `custom_placement_config` - (Optional) The bucket's custom location configuration, which specifies the individual regions that comprise a dual-region bucket. If the bucket is designated a single or multi-region, the parameters are empty. Structure is [documented below](#nested_custom_placement_config).

<a name="nested_lifecycle_rule"></a>The `lifecycle_rule` block supports:

* `action` - (Required) The Lifecycle Rule's action configuration. A single block of this type is supported. Structure is [documented below](#nested_action).

* `condition` - (Required) The Lifecycle Rule's condition configuration. A single block of this type is supported. Structure is [documented below](#nested_condition).

<a name="nested_action"></a>The `action` block supports:

* `type` - The type of the action of this Lifecycle Rule. Supported values include: `Delete`, `SetStorageClass` and `AbortIncompleteMultipartUpload`.

* `storage_class` - (Required if action type is `SetStorageClass`) The target [Storage Class](https://cloud.google.com/storage/docs/storage-classes) of objects affected by this Lifecycle Rule. Supported values include: `STANDARD`, `MULTI_REGIONAL`, `REGIONAL`, `NEARLINE`, `COLDLINE`, `ARCHIVE`.

<a name="nested_condition"></a>The `condition` block supports the following elements, and requires at least one to be defined. If you specify multiple conditions in a rule, an object has to match all of the conditions for the action to be taken:

* `age` - (Optional) Minimum age of an object in days to satisfy this condition.

* `created_before` - (Optional) A date in the RFC 3339 format YYYY-MM-DD. This condition is satisfied when an object is created before midnight of the specified date in UTC.

* `with_state` - (Optional) Match to live and/or archived objects. Unversioned buckets have only live objects. Supported values include: `"LIVE"`, `"ARCHIVED"`, `"ANY"`.

* `matches_storage_class` - (Optional) [Storage Class](https://cloud.google.com/storage/docs/storage-classes) of objects to satisfy this condition. Supported values include: `STANDARD`, `MULTI_REGIONAL`, `REGIONAL`, `NEARLINE`, `COLDLINE`, `ARCHIVE`, `DURABLE_REDUCED_AVAILABILITY`.

* `matches_prefix` - (Optional) One or more matching name prefixes to satisfy this condition.

* `matches_suffix` - (Optional) One or more matching name suffixes to satisfy this condition.

* `num_newer_versions` - (Optional) Relevant only for versioned objects. The number of newer versions of an object to satisfy this condition.

* `custom_time_before` - (Optional) A date in the RFC 3339 format YYYY-MM-DD. This condition is satisfied when the customTime metadata for the object is set to an earlier date than the date used in this lifecycle condition.

* `days_since_custom_time` - (Optional)	Days since the date set in the `customTime` metadata for the object. This condition is satisfied when the current date and time is at least the specified number of days after the `customTime`.

* `days_since_noncurrent_time` - (Optional) Relevant only for versioned objects. Number of days elapsed since the noncurrent timestamp of an object.

* `noncurrent_time_before` - (Optional) Relevant only for versioned objects. The date in RFC 3339 (e.g. `2017-06-13`) when the object became nonconcurrent.

<a name="nested_versioning"></a>The `versioning` block supports:

* `enabled` - (Required) While set to `true`, versioning is fully enabled for this bucket.

<a name="nested_website"></a>The `website` block supports the following elements, and requires at least one to be defined:

* `main_page_suffix` - (Optional) Behaves as the bucket's directory index where
    missing objects are treated as potential directories.

* `not_found_page` - (Optional) The custom object to return when a requested
    resource is not found.

<a name="nested_cors"></a>The `cors` block supports:

* `origin` - (Optional) The list of [Origins](https://tools.ietf.org/html/rfc6454) eligible to receive CORS response headers. Note: "*" is permitted in the list of origins, and means "any Origin".

* `method` - (Optional) The list of HTTP methods on which to include CORS response headers, (GET, OPTIONS, POST, etc) Note: "*" is permitted in the list of methods, and means "any method".

* `response_header` - (Optional) The list of HTTP headers other than the [simple response headers](https://www.w3.org/TR/cors/#simple-response-header) to give permission for the user-agent to share across domains.

* `max_age_seconds` - (Optional) The value, in seconds, to return in the [Access-Control-Max-Age header](https://www.w3.org/TR/cors/#access-control-max-age-response-header) used in preflight responses.

<a name="nested_retention_policy"></a>The `retention_policy` block supports:

* `is_locked` - (Optional) If set to `true`, the bucket will be [locked](https://cloud.google.com/storage/docs/using-bucket-lock#lock-bucket) and permanently restrict edits to the bucket's retention policy.  Caution: Locking a bucket is an irreversible action.

* `retention_period` - (Required) The period of time, in seconds, that objects in the bucket must be retained and cannot be deleted, overwritten, or archived. The value must be less than 2,147,483,647 seconds.

<a name="nested_logging"></a>The `logging` block supports:

* `log_bucket` - (Required) The bucket that will receive log objects.

* `log_object_prefix` - (Optional, Computed) The object prefix for log objects. If it's not provided,
    by default GCS sets this to this bucket's name.

<a name="nested_encryption"></a>The `encryption` block supports:

* `default_kms_key_name`: The `id` of a Cloud KMS key that will be used to encrypt objects inserted into this bucket, if no encryption method is specified.
  You must pay attention to whether the crypto key is available in the location that this bucket is created in.
  See [the docs](https://cloud.google.com/storage/docs/encryption/using-customer-managed-keys) for more details.

-> As per [the docs](https://cloud.google.com/storage/docs/encryption/using-customer-managed-keys) for customer-managed encryption keys, the IAM policy for the
  specified key must permit the [automatic Google Cloud Storage service account](https://cloud.google.com/storage/docs/projects#service-accounts) for the bucket's
  project to use the specified key for encryption and decryption operations.
  Although the service account email address follows a well-known format, the service account is created on-demand and may not necessarily exist for your project
  until a relevant action has occurred which triggers its creation.
  You should use the [`google_storage_project_service_account`](/docs/providers/google/d/storage_project_service_account.html) data source to obtain the email
  address for the service account when configuring IAM policy on the Cloud KMS key.
  This data source calls an API which creates the account if required, ensuring your Terraform applies cleanly and repeatedly irrespective of the
  state of the project.
  You should take care for race conditions when the same Terraform manages IAM policy on the Cloud KMS crypto key. See the data source page for more details.

<a name="nested_custom_placement_config"></a>The `custom_placement_config` block supports:

* `data_locations` - (Required) The list of individual regions that comprise a dual-region bucket. See [Cloud Storage bucket locations](https://cloud.google.com/storage/docs/dual-regions#availability) for a list of acceptable regions. **Note**: If any of the data_locations changes, it will [recreate the bucket](https://cloud.google.com/storage/docs/locations#key-concepts).

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `self_link` - The URI of the created resource.

* `url` - The base URL of the bucket, in the format `gs://<bucket-name>`.

## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 4 minutes.
- `update` - Default is 4 minutes.
- `read` - Default is 4 minutes.

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
