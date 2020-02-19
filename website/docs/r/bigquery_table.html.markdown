---
subcategory: "BigQuery"
layout: "google"
page_title: "Google: google_bigquery_table"
sidebar_current: "docs-google-bigquery-table"
description: |-
  Creates a table resource in a dataset for Google BigQuery.
---

# google_bigquery_table

Creates a table resource in a dataset for Google BigQuery. For more information see
[the official documentation](https://cloud.google.com/bigquery/docs/) and
[API](https://cloud.google.com/bigquery/docs/reference/rest/v2/tables).


## Example Usage

```hcl
resource "google_bigquery_dataset" "default" {
  dataset_id                  = "foo"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "EU"
  default_table_expiration_ms = 3600000

  labels = {
    env = "default"
  }
}

resource "google_bigquery_table" "default" {
  dataset_id = google_bigquery_dataset.default.dataset_id
  table_id   = "bar"

  time_partitioning {
    type = "DAY"
  }

  labels = {
    env = "default"
  }

  schema = <<EOF
[
  {
    "name": "permalink",
    "type": "STRING",
    "mode": "NULLABLE",
    "description": "The Permalink"
  },
  {
    "name": "state",
    "type": "STRING",
    "mode": "NULLABLE",
    "description": "State where the head office is located"
  }
]
EOF

}

resource "google_bigquery_table" "sheet" {
  dataset_id = google_bigquery_dataset.default.dataset_id
  table_id   = "sheet"

  external_data_configuration {
    autodetect    = true
    source_format = "GOOGLE_SHEETS"

    google_sheets_options {
      skip_leading_rows = 1
    }

    source_uris = [
      "https://docs.google.com/spreadsheets/d/123456789012345",
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `dataset_id` - (Required) The dataset ID to create the table in.
    Changing this forces a new resource to be created.

* `table_id` - (Required) A unique ID for the resource.
    Changing this forces a new resource to be created.

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

* `description` - (Optional) The field description.

* `expiration_time` - (Optional) The time when this table expires, in
    milliseconds since the epoch. If not present, the table will persist
    indefinitely. Expired tables will be deleted and their storage
    reclaimed.

* `external_data_configuration` - (Optional) Describes the data format,
    location, and other properties of a table stored outside of BigQuery.
    By defining these properties, the data source can then be queried as
    if it were a standard BigQuery table. Structure is documented below.

* `friendly_name` - (Optional) A descriptive name for the table.

* `encryption_configuration` - (Optional) Specifies how the table should be encrypted.
    If left blank, the table will be encrypted with a Google-managed key; that process
    is transparent to the user.  Structure is documented below.

* `labels` - (Optional) A mapping of labels to assign to the resource.

* `schema` - (Optional) A JSON schema for the table. Schema is required
    for CSV and JSON formats and is disallowed for Google Cloud
    Bigtable, Cloud Datastore backups, and Avro formats when using
    external tables. For more information see the
    [BigQuery API documentation](https://cloud.google.com/bigquery/docs/reference/rest/v2/tables#resource).
    ~>**NOTE**: Because this field expects a JSON string, any changes to the
    string will create a diff, even if the JSON itself hasn't changed.
    If the API returns a different value for the same schema, e.g. it
    switched the order of values or replaced `STRUCT` field type with `RECORD`
    field type, we currently cannot suppress the recurring diff this causes.
    As a workaround, we recommend using the schema as returned by the API.

* `time_partitioning` - (Optional) If specified, configures time-based
    partitioning for this table. Structure is documented below.

* `clustering` - (Optional) Specifies column names to use for data clustering.
    Up to four top-level columns are allowed, and should be specified in
    descending priority order.

* `view` - (Optional) If specified, configures this table as a view.
    Structure is documented below.

The `external_data_configuration` block supports:

* `autodetect` - (Required) - Let BigQuery try to autodetect the schema
    and format of the table.

* `compression` (Optional) - The compression type of the data source.
    Valid values are "NONE" or "GZIP".

* `csv_options` (Optional) - Additional properties to set if
    `source_format` is set to "CSV". Structure is documented below.

* `google_sheets_options` (Optional) - Additional options if
    `source_format` is set to "GOOGLE_SHEETS". Structure is
    documented below.

* `ignore_unknown_values` (Optional) - Indicates if BigQuery should
    allow extra values that are not represented in the table schema.
    If true, the extra values are ignored. If false, records with
    extra columns are treated as bad records, and if there are too
    many bad records, an invalid error is returned in the job result.
    The default value is false.

* `max_bad_records` (Optional) - The maximum number of bad records that
    BigQuery can ignore when reading data.

* `source_format` (Required) - The data format. Supported values are:
    "CSV", "GOOGLE_SHEETS", "NEWLINE_DELIMITED_JSON", "AVRO", "PARQUET",
    and "DATSTORE_BACKUP". To use "GOOGLE_SHEETS"
    the `scopes` must include
    "https://www.googleapis.com/auth/drive.readonly".

* `source_uris` - (Required) A list of the fully-qualified URIs that point to
    your data in Google Cloud.

The `csv_options` block supports:

* `quote` (Required) - The value that is used to quote data sections in a
    CSV file. If your data does not contain quoted sections, set the
    property value to an empty string. If your data contains quoted newline
    characters, you must also set the `allow_quoted_newlines` property to true.
    The API-side default is `"`, specified in Terraform escaped as `\"`. Due to
    limitations with Terraform default values, this value is required to be
    explicitly set. 

* `allow_jagged_rows` (Optional) - Indicates if BigQuery should accept rows
    that are missing trailing optional columns.

* `allow_quoted_newlines` (Optional) - Indicates if BigQuery should allow
    quoted data sections that contain newline characters in a CSV file.
    The default value is false.

* `encoding` (Optional) - The character encoding of the data. The supported
    values are UTF-8 or ISO-8859-1.

* `field_delimiter` (Optional) - The separator for fields in a CSV file.

* `skip_leading_rows` (Optional) - The number of rows at the top of a CSV
    file that BigQuery will skip when reading the data.

The `google_sheets_options` block supports:

* `range` (Optional) - Range of a sheet to query from. Only used when
    non-empty. At least one of `range` or `skip_leading_rows` must be set.
    Typical format: "sheet_name!top_left_cell_id:bottom_right_cell_id"
    For example: "sheet1!A1:B20"

* `skip_leading_rows` (Optional) - The number of rows at the top of the sheet
    that BigQuery will skip when reading the data. At least one of `range` or
    `skip_leading_rows` must be set.

The `time_partitioning` block supports:

* `expiration_ms` -  (Optional) Number of milliseconds for which to keep the
    storage for a partition.

* `field` - (Optional) The field used to determine how to create a time-based
    partition. If time-based partitioning is enabled without this value, the
    table is partitioned based on the load time.

* `type` - (Required) The only type supported is DAY, which will generate
    one partition per day based on data loading time.

* `require_partition_filter` - (Optional) If set to true, queries over this table
    require a partition filter that can be used for partition elimination to be
    specified.

The `view` block supports:

* `query` - (Required) A query that BigQuery executes when the view is referenced.

* `use_legacy_sql` - (Optional) Specifies whether to use BigQuery's legacy SQL for this view.
    The default value is true. If set to false, the view will use BigQuery's standard SQL.

The `encryption_configuration` block supports the following arguments:

* `kms_key_name` - (Required) The self link or full name of a key which should be used to
    encrypt this table.  Note that the default bigquery service account will need to have
    encrypt/decrypt permissions on this key - you may want to see the
    `google_bigquery_default_service_account` datasource and the
    `google_kms_crypto_key_iam_binding` resource.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `creation_time` - The time when this table was created, in milliseconds since the epoch.

* `etag` - A hash of the resource.

* `last_modified_time` - The time when this table was last modified, in milliseconds since the epoch.

* `location` - The geographic location where the table resides. This value is inherited from the dataset.

* `num_bytes` - The size of this table in bytes, excluding any data in the streaming buffer.

* `num_long_term_bytes` - The number of bytes in the table that are considered "long-term storage".

* `num_rows` - The number of rows of data in this table, excluding any data in the streaming buffer.

* `self_link` - The URI of the created resource.

* `type` - Describes the table type.

## Import

BigQuery tables can be imported using the `project`, `dataset_id`, and `table_id`, e.g.

```
$ terraform import google_bigquery_table.default gcp-project/foo/bar
```
