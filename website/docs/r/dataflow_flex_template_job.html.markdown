---
subcategory: "Dataflow"
layout: "google"
page_title: "Google: google_dataflow_flex_template_job"
sidebar_current: "docs-google-dataflow-flex-template-job"
description: |-
  Creates a job in Dataflow based on a Flex Template.
---

# google\_dataflow\_flex\_template\_job

Creates a [Flex Template](https://cloud.google.com/dataflow/docs/guides/templates/using-flex-templates)
job on Dataflow, which is an implementation of Apache Beam running on Google
Compute Engine. For more information see the official documentation for [Beam](https://beam.apache.org)
and [Dataflow](https://cloud.google.com/dataflow/).

## Example Usage

```hcl
resource "google_dataflow_flex_template_job" "big_data_job" {
  provider                = google-beta
  name                    = "dataflow-flextemplates-job"
  container_spec_gcs_path = "gs://my-bucket/templates/template.json"
  parameters = {
    inputSubscription = "messages"
  }
}
```

## Note on "destroy" / "apply"
There are many types of Dataflow jobs.  Some Dataflow jobs run constantly,
getting new data from (e.g.) a GCS bucket, and outputting data continuously.
Some jobs process a set amount of data then terminate. All jobs can fail while
running due to programming errors or other issues. In this way, Dataflow jobs
are different from most other Terraform / Google resources.

The Dataflow resource is considered 'existing' while it is in a nonterminal
state.  If it reaches a terminal state (e.g. 'FAILED', 'COMPLETE',
'CANCELLED'), it will be recreated on the next 'apply'.  This is as expected for
jobs which run continuously, but may surprise users who use this resource for
other kinds of Dataflow jobs.

A Dataflow job which is 'destroyed' may be "cancelled" or "drained".  If
"cancelled", the job terminates - any data written remains where it is, but no
new data will be processed.  If "drained", no new data will enter the pipeline,
but any data currently in the pipeline will finish being processed.  The default
is "cancelled", but if a user sets `on_delete` to `"drain"` in the
configuration, you may experience a long wait for your `terraform destroy` to
complete.

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for the resource, required by Dataflow.

* `container_spec_gcs_path` - (Required) The GCS path to the Dataflow job Flex
Template.

- - -

* `parameters` - (Optional) Key/Value pairs to be passed to the Dataflow job (as
used in the template). Additional [pipeline options](https://cloud.google.com/dataflow/docs/guides/specifying-exec-params#setting-other-cloud-dataflow-pipeline-options)
such as `serviceAccount`, `workerMachineType`, etc can be specified here.

* `labels` - (Optional) User labels to be specified for the job. Keys and values
should follow the restrictions specified in the [labeling restrictions](https://cloud.google.com/compute/docs/labeling-resources#restrictions)
page. **Note**: This field is marked as deprecated in Terraform as the API does not currently
support adding labels.
**NOTE**: Google-provided Dataflow templates often provide default labels
that begin with `goog-dataflow-provided`. Unless explicitly set in config, these
labels will be ignored to prevent diffs on re-apply.

* `on_delete` - (Optional) One of "drain" or "cancel". Specifies behavior of
deletion during `terraform destroy`.  See above note.

* `project` - (Optional) The project in which the resource belongs. If it is not
provided, the provider project is used.

* `region` - (Optional) The region in which the created job should run.

## Attributes Reference
In addition to the arguments listed above, the following computed attributes are exported:

* `job_id` - The unique ID of this job.

* `state` - The current state of the resource, selected from the [JobState enum](https://cloud.google.com/dataflow/docs/reference/rest/v1b3/projects.jobs#Job.JobState)

## Import

This resource does not support import.
