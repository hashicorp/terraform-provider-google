---
subcategory: "Compute Engine"
layout: "google"
page_title: "Google: google_compute_backend_bucket"
sidebar_current: "docs-google-datasource-compute-backend-bucket"
description: |-
  Get information about a BackendBucket.
---

# google\_compute\_backend\_bucket

Get information about a BackendBucket.

## Example Usage

```tf
data "google_compute_backend_bucket" "my-backend-bucket" {
  name = "my-backend"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the resource.

- - -

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `bucket_name` - Cloud Storage bucket name.

* `cdn_policy` - Cloud CDN configuration for this Backend Bucket. Structure is documented below.

* `description` - An optional textual description of the resource; provided by the client when the resource is created.

* `enable_cdn` - Whether Cloud CDN is enabled for this BackendBucket.

* `id` - an identifier for the resource with format `projects/{{project}}/global/backendBuckets/{{name}}`

* `creation_timestamp` - Creation timestamp in RFC3339 text format.

* `self_link` - The URI of the created resource.

The `cdn_policy` block supports:

* `signed_url_cache_max_age_sec` - Maximum number of seconds the response to a signed URL request will be considered fresh. After this time period, the response will be revalidated before being served. When serving responses to signed URL requests, Cloud CDN will internally behave as though all responses from this backend had a "Cache-Control: public, max-age=[TTL]" header, regardless of any existing Cache-Control header. The actual headers served in responses will not be altered.
