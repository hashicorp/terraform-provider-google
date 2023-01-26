---
# ----------------------------------------------------------------------------
#
#     ***     AUTO GENERATED CODE    ***    Type: DCL     ***
#
# ----------------------------------------------------------------------------
#
#     This file is managed by Magic Modules (https:#github.com/GoogleCloudPlatform/magic-modules)
#     and is based on the DCL (https:#github.com/GoogleCloudPlatform/declarative-resource-client-library).
#     Changes will need to be made to the DCL or Magic Modules instead of here.
#
#     We are not currently able to accept contributions to this file. If changes
#     are required, please file an issue at https:#github.com/hashicorp/terraform-provider-google/issues/new/choose
#
# ----------------------------------------------------------------------------
subcategory: "Logging"
description: |-
  The Logging LogView resource
---

# google_logging_log_view

The Logging LogView resource

## Example Usage - basic
```hcl
resource "google_logging_log_view" "primary" {
  name        = "view"
  bucket      = google_logging_project_bucket_config.basic.id
  description = "A logging view configured with Terraform"
  filter      = "SOURCE(\"projects/myproject\") AND resource.type = \"gce_instance\" AND LOG_ID(\"stdout\")"
}

resource "google_logging_project_bucket_config" "basic" {
    project        = "my-project-name"
    location       = "global"
    retention_days = 30
    bucket_id      = "_Default"
}

```

## Argument Reference

The following arguments are supported:

* `bucket` -
  (Required)
  The bucket of the resource
  
* `name` -
  (Required)
  The resource name of the view. For example: `projects/my-project/locations/global/buckets/my-bucket/views/my-view`
  


- - -

* `description` -
  (Optional)
  Describes this view.
  
* `filter` -
  (Optional)
  Filter that restricts which log entries in a bucket are visible in this view. Filters are restricted to be a logical AND of ==/!= of any of the following: - originating project/folder/organization/billing account. - resource type - log id For example: SOURCE("projects/myproject") AND resource.type = "gce_instance" AND LOG_ID("stdout")
  
* `location` -
  (Optional)
  The location of the resource. The supported locations are: global, us-central1, us-east1, us-west1, asia-east1, europe-west1.
  
* `parent` -
  (Optional)
  The parent of the resource.
  


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `{{parent}}/locations/{{location}}/buckets/{{bucket}}/views/{{name}}`

* `create_time` -
  Output only. The creation timestamp of the view.
  
* `update_time` -
  Output only. The last update timestamp of the view.
  
## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

LogView can be imported using any of these accepted formats:

```
$ terraform import google_logging_log_view.default {{parent}}/locations/{{location}}/buckets/{{bucket}}/views/{{name}}
```



