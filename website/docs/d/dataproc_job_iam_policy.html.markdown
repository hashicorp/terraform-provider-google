---
subcategory: "Dataproc"
description: |-
  A datasource to retrieve the IAM policy state for a Dataproc job.
---


# `google_dataproc_job_iam_policy`
Retrieves the current IAM policy data for a Dataproc job.

## example

```hcl
data "google_dataproc_job_iam_policy" "policy" {
  job_id      = google_dataproc_job.pyspark.reference[0].job_id
  region      = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `job_id` - (Required) The name or relative resource id of the job to manage IAM policies for.

* `project` - (Optional) The project in which the job belongs. If it
    is not provided, Terraform will use the provider default.

* `region` - (Optional) The region in which the job belongs. If it
    is not provided, Terraform will use the provider default.

## Attributes Reference

The attributes are exported:

* `etag` - (Computed) The etag of the IAM policy.

* `policy_data` - (Computed) The policy data
