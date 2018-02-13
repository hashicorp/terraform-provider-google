---
layout: "google"
page_title: "Google: google_project_usage_export_bucket"
sidebar_current: "docs-google-project-usage-export-bucket"
description: |-
  Creates a dataset resource for Google BigQuery.
---

# google_project_usage_export_bucket

Sets up a usage export bucket for a particular project.  A usage export bucket
is a pre-configured GCS bucket which is set up to receive daily and monthly
reports of the GCE resources used.

For more information see the [Docs](https://cloud.google.com/compute/docs/usage-export)
and for further details, the
[API Documentation](https://cloud.google.com/compute/docs/reference/rest/beta/projects/setUsageExportBucket).


## Example Usage

```hcl
resource "google_project_usage_export_bucket" "export" {
  project      = "foo"
  bucket_name  = "bar"
}
```

## Argument Reference
* `project`: (Required) The project to set the export bucket on.
* `bucket_name`: (Required) The bucket to store reports in.
* `prefix`: (Optional) A prefix for the reports, for instance, the project name.

## Note

You should specify only one of these per project.  If there are two or more
they will fight over which bucket the reports should be stored in.  It is
safe to have multiple resources with the same backing bucket.
