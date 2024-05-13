---
subcategory: "Cloud (Stackdriver) Logging"
description: |-
  Manages a folder-level logging exclusion.
---

# google_logging_folder_exclusion

Manages a folder-level logging exclusion. For more information see:

* [API documentation](https://cloud.google.com/logging/docs/reference/v2/rest/v2/folders.exclusions)
* How-to Guides
    * [Excluding Logs](https://cloud.google.com/logging/docs/exclusions)

~> You can specify exclusions for log sinks created by terraform by using the exclusions field of `google_logging_folder_sink`

## Example Usage

```hcl
resource "google_logging_folder_exclusion" "my-exclusion" {
  name   = "my-instance-debug-exclusion"
  folder = google_folder.my-folder.name

  description = "Exclude GCE instance debug logs"

  # Exclude all DEBUG or lower severity messages relating to instances
  filter = "resource.type = gce_instance AND severity <= DEBUG"
}

resource "google_folder" "my-folder" {
  display_name = "My folder"
  parent       = "organizations/123456"
}
```

## Argument Reference

The following arguments are supported:

* `folder` - (Required) The folder to be exported to the sink. Note that either [FOLDER_ID] or "folders/[FOLDER_ID]" is
    accepted.

* `name` - (Required) The name of the logging exclusion.

* `description` - (Optional) A human-readable description.

* `disabled` - (Optional) Whether this exclusion rule should be disabled or not. This defaults to
    false.

* `filter` - (Required) The filter to apply when excluding logs. Only log entries that match the filter are excluded.
    See [Advanced Log Filters](https://cloud.google.com/logging/docs/view/advanced-filters) for information on how to
    write a filter.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `folders/{{folder}}/exclusions/{{name}}`

## Import

Folder-level logging exclusions can be imported using their URI, e.g.

* `folders/{{folder}}/exclusions/{{name}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import folder-level logging exclusions using one of the formats above. For example:

```tf
import {
  id = "folders/{{folder}}/exclusions/{{name}}"
  to = google_logging_folder_exclusion.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), folder-level logging exclusions can be imported using one of the formats above. For example:

```
$ terraform import google_logging_folder_exclusion.default folders/{{folder}}/exclusions/{{name}}
```
