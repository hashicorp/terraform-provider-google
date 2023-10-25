---
subcategory: "BigQuery"
description: |-
  A datasource to retrieve information about a BigQuery dataset.
---

# `google_bigquery_dataset`

Get information about a BigQuery dataset. For more information see
the [official documentation](https://cloud.google.com/bigquery/docs)
and [API](https://cloud.google.com/bigquery/docs/reference/rest/v2/datasets).

## Example Usage

```hcl
data "google_bigquery_dataset" "dataset" {
  dataset_id = "my-bq-dataset"
  project = "my-project"
}
```

## Argument Reference

The following arguments are supported:

* `dataset_id` - (Required) The dataset ID.

* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.

## Attributes Reference

See [google_bigquery_dataset](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/bigquery_dataset) resource for details of the available attributes.
