---
subcategory: "Cloud (Stackdriver) Logging"
layout: "google"
page_title: "Google: google_logging_view"
sidebar_current: "docs-google-logging-view"
description: |-
  Manages a log bucket view.
---

# google\_logging\_view

Manages a log bucket view. For more information see
[the official logging documentation](https://cloud.google.com/logging/docs/) and
[Managing log views on your log buckets](https://cloud.google.com/logging/docs/logs-views).

~> **Note:** Log views are automatically created for a log bucket. Creating a resource of this type will fail in Terraform. The log views that are currently automatically created are `_AllLogs` and `_Default`.

## Example Usage

```hcl
resource "google_logging_project_bucket_config" "default" {
	bucket_id      = "_Default"
	project        = google_project.default.name
	location       = "global"
	retention_days = 30
}

resource "google_logging_view" "only_compute_instances" {
  view_id     = "OnlyComputeInstances"
  bucket      = google_logging_project_bucket_config.default.id
  description = "Compute instance logs"
  filter      = "resource.type = gce_instance"
}
```

## Argument Reference

The following arguments are supported:

* `view_id` - (Required) Log view identifier.

* `bucket` - (Required) Resource name of the log bucket with the following format: `{{parent}}/{{parent_id}}/locations/{{location}}/buckets/{{bucket_id}}`.

* `description` - (Optional) Description of this log view.

* `filter` - (Optional) Log view filter. Read more about filter constraints in the [official documentation](https://cloud.google.com/logging/docs/logs-views#before_you_begin).

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `id` - An identifier for the resource with format `{{parent}}/{{parent_id}}/locations/{{location}}/buckets/{{bucket_id}}/views/{{view_id}}`.

* `name` - The resource name of the log view. For example: `{{parent}}/{{parent_id}}/locations/{{location}}/buckets/{{bucket_id}}/views/{{view_id}}`.

## Import

This resource can be imported using the following format:

```
$ terraform import google_logging_view.basic {{parent}}/{{parent_id}}/locations/{{location}}/buckets/{{bucket_id}}/views/{{view_id}}
```
