---
subcategory: "Cloud Functions"
description: |-
  Get information about a Google Cloud Function.
---

# google\_cloudfunctions\_function

Get information about a Google Cloud Function. For more information see
the [official documentation](https://cloud.google.com/functions/docs/)
and [API](https://cloud.google.com/functions/docs/apis).

## Example Usage

```hcl
data "google_cloudfunctions_function" "my-function" {
  name = "function"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of a Cloud Function.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

* `region` - (Optional) The region in which the resource belongs. If it
    is not provided, the provider region is used.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `name` - The name of the Cloud Function.
* `source_archive_bucket` - The GCS bucket containing the zip archive which contains the function.
* `source_archive_object` - The source archive object (file) in archive bucket.
* `description` - Description of the function.
* `available_memory_mb` - Available memory (in MB) to the function.
* `timeout` - Function execution timeout (in seconds).
* `runtime` - The runtime in which the function is running.
* `entry_point` - Name of a JavaScript function that will be executed when the Google Cloud Function is triggered.
* `trigger_http` - If function is triggered by HTTP, this boolean is set.
* `event_trigger` - A source that fires events in response to a condition in another service. Structure is [documented below](#nested_event_trigger).
* `https_trigger_url` - If function is triggered by HTTP, trigger URL is set here.
* `ingress_settings` - Controls what traffic can reach the function.
* `labels` - All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.
* `service_account_email` - The service account email to be assumed by the cloud function.
* `vpc_connector` - The VPC Network Connector that this cloud function can connect to. 
* `vpc_connector_egress_settings` - The egress settings for the connector, controlling what traffic is diverted through it.
* `max_instances` - The limit on the maximum number of function instances that may coexist at a given time. If unset or set to `0`, the API default will be used. 
* `source_repository` - The URL of the Cloud Source Repository that the function is deployed from. Structure is [documented below](#nested_source_repository).

<a name="nested_event_trigger"></a>The `event_trigger` block contains:

* `event_type` - The type of event to observe. For example: `"google.storage.object.finalize"`.
See the documentation on [calling Cloud Functions](https://cloud.google.com/functions/docs/calling/)
for a full reference of accepted triggers.

* `resource` - The name of the resource whose events are being observed, for example, `"myBucket"`

* `failure_policy` - Policy for failed executions. Structure is [documented below](#nested_failure_policy).

<a name="nested_failure_policy"></a>The `failure_policy` block contains:

* `retry` - Whether the function should be retried on failure.


<a name="nested_source_repository"></a>The `source_repository` block contains:

* `url` - The URL pointing to the hosted repository where the function is defined.
