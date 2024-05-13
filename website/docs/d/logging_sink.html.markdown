---
subcategory: "Cloud (Stackdriver) Logging"
description: |-
  Get information about a Google Cloud Logging Sink.
---

# google_logging_sink

Use this data source to get a project, folder, organization or billing account logging sink details.
To get more information about Service, see:

[API documentation](https://cloud.google.com/logging/docs/reference/v2/rest/v2/sinks)

## Example Usage - Retrieve Project Logging Sink Basic


```hcl
data google_logging_sink "project-sink" {
	id = "projects/0123456789/sinks/my-sink-name"
}
```

## Argument Reference

The following arguments are supported:



- - -

* `id` - (Required) The identifier for the resource. 
    Examples:

    - `projects/[PROJECT_ID]/sinks/[SINK_NAME]`
    - `organizations/[ORGANIZATION_ID]/sinks/[SINK_NAME]`
    -  `billingAccounts/[BILLING_ACCOUNT_ID]/sinks/[SINK_NAME]`
    - `folders/[FOLDER_ID]/sinks/[SINK_NAME]`


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:


* `name` - The name of the logging sink.

* `destination` - The destination of the sink (or, in other words, where logs are written to).

* `filter` - The filter which is applied when exporting logs. Only log entries that match the filter are exported.

* `description` - The description of this sink.

* `disabled` - Whether this sink is disabled and it does not export any log entries.

* `writer_identity` - The identity associated with this sink. This identity must be granted write access to the configured `destination`.

* `bigquery_options` - Options that affect sinks exporting data to BigQuery. Structure is [documented below](#nested_bigquery_options).

* `exclusions` - Log entries that match any of the exclusion filters are not exported. Structure is [documented below](#nested_exclusions).

<a name="nested_bigquery_options"></a>The `bigquery_options` block supports:

* `use_partitioned_tables` - Whether [BigQuery's partition tables](https://cloud.google.com/bigquery/docs/partitioned-tables) are used.

<a name="nested_exclusions"></a>The `exclusions` block supports:

* `name` - A client-assigned identifier, such as `load-balancer-exclusion`. 
* `description` - A description of this exclusion.
* `filter` - An advanced logs filter that matches the log entries to be excluded. 
* `disabled` - Whether this exclusion is disabled and it does not exclude any log entries.
