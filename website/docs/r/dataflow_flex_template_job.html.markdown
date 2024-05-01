---
subcategory: "Dataflow"
description: |-
  Creates a job in Dataflow based on a Flex Template.
---

# google\_dataflow\_flex\_template\_job

Creates a [Flex Template](https://cloud.google.com/dataflow/docs/guides/templates/using-flex-templates)
job on Dataflow, which is an implementation of Apache Beam running on Google
Compute Engine. For more information see the official documentation for [Beam](https://beam.apache.org)
and [Dataflow](https://cloud.google.com/dataflow/).

~> **Warning:** This resource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta resources.

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

You can potentially short-circuit the wait by setting `skip_wait_on_job_termination`
to `true`, but beware that unless you take active steps to ensure that the job
`name` parameter changes between instances, the name will conflict and the launch
of the new job will fail. One way to do this is with a
[random_id](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/id)
resource, for example:

```hcl
variable "big_data_job_subscription_id" {
  type    = string
  default = "projects/myproject/subscriptions/messages"
}

resource "random_id" "big_data_job_name_suffix" {
  byte_length = 4
  keepers = {
    region          = var.region
    subscription_id = var.big_data_job_subscription_id
  }
}
resource "google_dataflow_flex_template_job" "big_data_job" {
  provider                      = google-beta
  name                          = "dataflow-flextemplates-job-${random_id.big_data_job_name_suffix.dec}"
  region                        = var.region
  container_spec_gcs_path       = "gs://my-bucket/templates/template.json"
  skip_wait_on_job_termination = true
  parameters = {
    inputSubscription = var.big_data_job_subscription_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Immutable. A unique name for the resource, required by Dataflow.

* `container_spec_gcs_path` - (Required) The GCS path to the Dataflow job Flex
Template.

- - -

* `additional_experiments` - (Optional) List of experiments that should be used by the job. An example value is `["enable_stackdriver_agent_metrics"]`.

* `autoscaling_algorithm` - (Optional) The algorithm to use for autoscaling.

* `parameters` - **Template specific** Key/Value pairs to be forwarded to the pipeline's options; keys are
  case-sensitive based on the language on which the pipeline is coded, mostly Java.
  **Note**: do not configure Dataflow options here in parameters.

* `enable_streaming_engine` - (Optional) Immutable. Indicates if the job should use the streaming engine feature.

* `ip_configuration` - (Optional) The configuration for VM IPs.  Options are `"WORKER_IP_PUBLIC"` or `"WORKER_IP_PRIVATE"`.

* `kms_key_name` - (Optional) The name for the Cloud KMS key for the job. Key format is: `projects/PROJECT_ID/locations/LOCATION/keyRings/KEY_RING/cryptoKeys/KEY`

* `labels` - (Optional) User labels to be specified for the job. Keys and values
should follow the restrictions specified in the [labeling restrictions](https://cloud.google.com/compute/docs/labeling-resources#restrictions)
page. 
**Note**: This field is non-authoritative, and will only manage the labels present in your configuration. Please refer to the field `effective_labels` for all of the labels present on the resource.

* `terraform_labels` -
  The combination of labels configured directly on the resource and default labels configured on the provider.

* `effective_labels` -
  All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.

* `launcher_machine_type` - (Optional) The machine type to use for launching the job. The default is n1-standard-1.

* `machine_type` - (Optional) The machine type to use for the job.

* `max_workers` - (Optional) Immutable. The maximum number of Google Compute Engine instances to be made available to your pipeline during execution, from 1 to 1000.

* `network` - (Optional) The network to which VMs will be assigned. If it is not provided, "default" will be used.

* `num_workers` - (Optional) Immutable. The initial number of Google Compute Engine instances for the job.

* `on_delete` - (Optional) One of "drain" or "cancel". Specifies behavior of
deletion during `terraform destroy`.  See above note.

* `project` - (Optional) The project in which the resource belongs. If it is not
provided, the provider project is used.

* `region` - (Optional) Immutable. The region in which the created job should run.

* `sdk_container_image` - (Optional) Docker registry location of container image to use for the 'worker harness. Default is the container for the version of the SDK. Note this field is only valid for portable pipelines.

* `service_account_email` - (Optional) Service account email to run the workers as.

* `skip_wait_on_job_termination` - (Optional)  If set to `true`, terraform will
treat `DRAINING` and `CANCELLING` as terminal states when deleting the resource,
and will remove the resource from terraform state and move on.  See above note.

* `staging_location` - (Optional) The Cloud Storage path to use for staging files. Must be a valid Cloud Storage URL, beginning with gs://.

* `subnetwork` - (Optional) The subnetwork to which VMs will be assigned. Should be of the form "regions/REGION/subnetworks/SUBNETWORK".

* `temp_location` - (Optional) The Cloud Storage path to use for temporary files. Must be a valid Cloud Storage URL, beginning with gs://.

* `transform_name_mapping` - (Optional) Only applicable when updating a pipeline. Map of transform name prefixes of the job to be replaced with the corresponding name prefixes of the new job.Only applicable when updating a pipeline. Map of transform name prefixes of the job to be replaced with the corresponding name prefixes of the new job.

## Attributes Reference
In addition to the arguments listed above, the following computed attributes are exported:

* `job_id` - The unique ID of this job.

* `state` - The current state of the resource, selected from the [JobState enum](https://cloud.google.com/dataflow/docs/reference/rest/v1b3/projects.jobs#Job.JobState)

## Import

This resource does not support import.
