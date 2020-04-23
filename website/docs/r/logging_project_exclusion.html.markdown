---
subcategory: "Cloud (Stackdriver) Logging"
layout: "google"
page_title: "Google: google_logging_project_exclusion"
sidebar_current: "docs-google-logging-project-exclusion"
description: |-
  Manages a project-level logging exclusion.
---

# google\_logging\_project\_exclusion

Manages a project-level logging exclusion. For more information see
[the official documentation](https://cloud.google.com/logging/docs/) and
[Excluding Logs](https://cloud.google.com/logging/docs/exclusions).

Note that you must have the "Logs Configuration Writer" IAM role (`roles/logging.configWriter`)
granted to the credentials used with Terraform.

## Example Usage

```hcl
resource "google_logging_project_exclusion" "my-exclusion" {
  name = "my-instance-debug-exclusion"

  description = "Exclude GCE instance debug logs"

  # Exclude all DEBUG or lower severity messages relating to instances
  filter = "resource.type = gce_instance AND severity <= DEBUG"
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Required) The filter to apply when excluding logs. Only log entries that match the filter are excluded.
    See [Advanced Log Filters](https://cloud.google.com/logging/docs/view/advanced-filters) for information on how to
    write a filter.

* `name` - (Required) The name of the logging exclusion.

* `description` - (Optional) A human-readable description.

* `disabled` - (Optional) Whether this exclusion rule should be disabled or not. This defaults to
    false.

* `project` - (Optional) The project to create the exclusion in. If omitted, the project associated with the provider is
    used.

## Import

Project-level logging exclusions can be imported using their URI, e.g.

```
$ terraform import google_logging_project_exclusion.my_exclusion projects/my-project/exclusions/my-exclusion
```
