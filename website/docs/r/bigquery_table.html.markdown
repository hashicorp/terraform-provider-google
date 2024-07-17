---
subcategory: "BigQuery"
description: |-
  Creates a table resource in a dataset for Google BigQuery.
---

# google_bigquery_table

Creates a table resource in a dataset for Google BigQuery. For more information see
[the official documentation](https://cloud.google.com/bigquery/docs/) and
[API](https://cloud.google.com/bigquery/docs/reference/rest/v2/tables).

-> **Note**: On newer versions of the provider, you must explicitly set `deletion_protection=false`
(and run `terraform apply` to write the field to state) in order to destroy an instance.
It is recommended to not set this field (or set it to true) until you're ready to destroy.

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
    if it were a standard BigQuery table. Structure is [documented below](#nested_external_data_configuration).

* `friendly_name` - (Optional) A descriptive name for the table.

* `max_staleness`: (Optional) The maximum staleness of data that could be
  returned when the table (or stale MV) is queried. Staleness encoded as a
  string encoding of [SQL IntervalValue
  type](https://cloud.google.com/bigquery/docs/reference/standard-sql/data-types#interval_type).

* `encryption_configuration` - (Optional) Specifies how the table should be encrypted.
    If left blank, the table will be encrypted with a Google-managed key; that process
    is transparent to the user.  Structure is [documented below](#nested_encryption_configuration).

* `labels` - (Optional) A mapping of labels to assign to the resource.

    **Note**: This field is non-authoritative, and will only manage the labels present in your configuration.
    Please refer to the field 'effective_labels' for all of the labels present on the resource.

* `terraform_labels` -
  The combination of labels configured directly on the resource and default labels configured on the provider.

* `effective_labels` -
  All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.

* <a name="schema"></a>`schema` - (Optional) A JSON schema for the table.

    ~>**NOTE:** Because this field expects a JSON string, any changes to the
    string will create a diff, even if the JSON itself hasn't changed.
    If the API returns a different value for the same schema, e.g. it
    switched the order of values or replaced `STRUCT` field type with `RECORD`
    field type, we currently cannot suppress the recurring diff this causes.
    As a workaround, we recommend using the schema as returned by the API.

    ~>**NOTE:**  If you use `external_data_configuration`
    [documented below](#nested_external_data_configuration) and do **not** set
    `external_data_configuration.connection_id`, schemas must be specified
    with `external_data_configuration.schema`. Otherwise, schemas must be
    specified with this top-level field.

* `time_partitioning` - (Optional) If specified, configures time-based
    partitioning for this table. Structure is [documented below](#nested_time_partitioning).

* `range_partitioning` - (Optional) If specified, configures range-based
    partitioning for this table. Structure is [documented below](#nested_range_partitioning).

* `require_partition_filter` - (Optional) If set to true, queries over this table
    require a partition filter that can be used for partition elimination to be
    specified.

* `clustering` - (Optional) Specifies column names to use for data clustering.
    Up to four top-level columns are allowed, and should be specified in
    descending priority order.

* `view` - (Optional) If specified, configures this table as a view.
    Structure is [documented below](#nested_view).

* `materialized_view` - (Optional) If specified, configures this table as a materialized view.
    Structure is [documented below](#nested_materialized_view).

* `deletion_protection` - (Optional) Whether Terraform will be prevented from destroying the table.
    When the field is set to true or unset in Terraform state, a `terraform apply`
    or `terraform destroy` that would delete the table will fail.
    When the field is set to false, deleting the table is allowed..

* `table_constraints` - (Optional) Defines the primary key and foreign keys. 
    Structure is [documented below](#nested_table_constraints).

* `table_replication_info` - (Optional) Replication info of a table created
    using "AS REPLICA" DDL like:
    `CREATE MATERIALIZED VIEW mv1 AS REPLICA OF src_mv`.
    Structure is [documented below](#nested_table_replication_info).

* `resource_tags` - (Optional) The tags attached to this table. Tag keys are
    globally unique. Tag key is expected to be in the namespaced format, for
    example "123456789012/environment" where 123456789012 is the ID of the
    parent organization or project resource for this tag key. Tag value is
    expected to be the short name, for example "Production".

* `allow_resource_tags_on_deletion` - (Optional) If set to true, it allows table
    deletion when there are still resource tags attached. The default value is
    false.

<a name="nested_external_data_configuration"></a>The `external_data_configuration` block supports:

* `autodetect` - (Required) - Let BigQuery try to autodetect the schema
    and format of the table.

* `compression` (Optional) - The compression type of the data source.
    Valid values are "NONE" or "GZIP".

* `connection_id` (Optional) - The connection specifying the credentials to be used to read
    external storage, such as Azure Blob, Cloud Storage, or S3. The `connection_id` can have
    the form `{{project}}.{{location}}.{{connection_id}}`
    or `projects/{{project}}/locations/{{location}}/connections/{{connection_id}}`.

    ~>**NOTE:** If you set `external_data_configuration.connection_id`, the
    table schema must be specified using the top-level `schema` field
    [documented above](#schema).

* `csv_options` (Optional) - Additional properties to set if
    `source_format` is set to "CSV". Structure is [documented below](#nested_csv_options).

* `bigtable_options` (Optional) - Additional properties to set if
    `source_format` is set to "BIGTABLE". Structure is [documented below](#nested_bigtable_options).

* `json_options` (Optional) - Additional properties to set if
    `source_format` is set to "JSON". Structure is [documented below](#nested_json_options).

* `json_extension` (Optional) - Used to indicate that a JSON variant, rather than normal JSON, is being used as the sourceFormat. This should only be used in combination with the `JSON` source format. Valid values are: `GEOJSON`.

* `parquet_options` (Optional) - Additional properties to set if
    `source_format` is set to "PARQUET". Structure is [documented below](#nested_parquet_options).

* `google_sheets_options` (Optional) - Additional options if
    `source_format` is set to "GOOGLE_SHEETS". Structure is
    [documented below](#nested_google_sheets_options).

* `hive_partitioning_options` (Optional) - When set, configures hive partitioning
    support. Not all storage formats support hive partitioning -- requesting hive
    partitioning on an unsupported format will lead to an error, as will providing
    an invalid specification. Structure is [documented below](#nested_hive_partitioning_options).

* `avro_options` (Optional) - Additional options if `source_format` is set to
    "AVRO".  Structure is [documented below](#nested_avro_options).

* `ignore_unknown_values` (Optional) - Indicates if BigQuery should
    allow extra values that are not represented in the table schema.
    If true, the extra values are ignored. If false, records with
    extra columns are treated as bad records, and if there are too
    many bad records, an invalid error is returned in the job result.
    The default value is false.

* `max_bad_records` (Optional) - The maximum number of bad records that
    BigQuery can ignore when reading data.

* `schema` - (Optional) A JSON schema for the external table. Schema is required
    for CSV and JSON formats if autodetect is not on. Schema is disallowed
    for Google Cloud Bigtable, Cloud Datastore backups, Avro, Iceberg, ORC and Parquet formats.
    ~>**NOTE:** Because this field expects a JSON string, any changes to the
    string will create a diff, even if the JSON itself hasn't changed.
    Furthermore drift for this field cannot not be detected because BigQuery
    only uses this schema to compute the effective schema for the table, therefore
    any changes on the configured value will force the table to be recreated.
    This schema is effectively only applied when creating a table from an external
    datasource, after creation the computed schema will be stored in
    `google_bigquery_table.schema`

    ~>**NOTE:** If you set `external_data_configuration.connection_id`, the
    table schema must be specified using the top-level `schema` field
    [documented above](#schema).

* `source_format` (Optional) - The data format. Please see sourceFormat under
    [ExternalDataConfiguration](https://cloud.google.com/bigquery/docs/reference/rest/v2/tables#externaldataconfiguration)
    in Bigquery's public API documentation for supported formats. To use "GOOGLE_SHEETS"
    the `scopes` must include "https://www.googleapis.com/auth/drive.readonly".

* `source_uris` - (Required) A list of the fully-qualified URIs that point to
    your data in Google Cloud.

* `file_set_spec_type` - (Optional) Specifies how source URIs are interpreted for constructing the file set to load.
    By default source URIs are expanded against the underlying storage.
    Other options include specifying manifest files. Only applicable to object storage systems. [Docs](cloud/bigquery/docs/reference/rest/v2/tables#filesetspectype)

* `reference_file_schema_uri` - (Optional) When creating an external table, the user can provide a reference file with the table schema. This is enabled for the following formats: AVRO, PARQUET, ORC.

* `metadata_cache_mode` - (Optional) Metadata Cache Mode for the table. Set this to enable caching of metadata from external data source. Valid values are `AUTOMATIC` and `MANUAL`.

* `object_metadata` - (Optional) Object Metadata is used to create Object Tables. Object Tables contain a listing of objects (with their metadata) found at the sourceUris. If `object_metadata` is set, `source_format` should be omitted.

<a name="nested_csv_options"></a>The `csv_options` block supports:

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

<a name="nested_bigtable_options"></a>The `bigtable_options` block supports:

* `column_family` (Optional) - A list of column families to expose in the table schema along with their types. This list restricts the column families that can be referenced in queries and specifies their value types. You can use this list to do type conversions - see the 'type' field for more details. If you leave this list empty, all column families are present in the table schema and their values are read as BYTES. During a query only the column families referenced in that query are read from Bigtable.  Structure is [documented below](#nested_column_family).
* `ignore_unspecified_column_families` (Optional) - If field is true, then the column families that are not specified in columnFamilies list are not exposed in the table schema. Otherwise, they are read with BYTES type values. The default value is false.
* `read_rowkey_as_string` (Optional) - If field is true, then the rowkey column families will be read and converted to string. Otherwise they are read with BYTES type values and users need to manually cast them with CAST if necessary. The default value is false.
* `output_column_families_as_json` (Optional) - If field is true, then each column family will be read as a single JSON column. Otherwise they are read as a repeated cell structure containing timestamp/value tuples. The default value is false.

<a name="nested_column_family"></a>The `column_family` block supports:

* `column` (Optional) - A List of columns that should be exposed as individual fields as opposed to a list of (column name, value) pairs. All columns whose qualifier matches a qualifier in this list can be accessed as Other columns can be accessed as a list through column field.  Structure is [documented below](#nested_column).
* `family_id` (Optional) - Identifier of the column family.
* `type` (Optional) - The type to convert the value in cells of this column family. The values are expected to be encoded using HBase Bytes.toBytes function when using the BINARY encoding value. Following BigQuery types are allowed (case-sensitive): "BYTES", "STRING", "INTEGER", "FLOAT", "BOOLEAN", "JSON". Default type is BYTES. This can be overridden for a specific column by listing that column in 'columns' and specifying a type for it.
* `encoding` (Optional) - The encoding of the values when the type is not STRING. Acceptable encoding values are: TEXT - indicates values are alphanumeric text strings. BINARY - indicates values are encoded using HBase Bytes.toBytes family of functions. This can be overridden for a specific column by listing that column in 'columns' and specifying an encoding for it.
* `only_read_latest` (Optional) - If this is set only the latest version of value are exposed for all columns in this column family. This can be overridden for a specific column by listing that column in 'columns' and specifying a different setting for that column.

<a name="nested_column"></a>The `column` block supports:

* `qualifier_encoded` (Optional) - Qualifier of the column. Columns in the parent column family that has this exact qualifier are exposed as . field. If the qualifier is valid UTF-8 string, it can be specified in the qualifierString field. Otherwise, a base-64 encoded value must be set to qualifierEncoded. The column field name is the same as the column qualifier. However, if the qualifier is not a valid BigQuery field identifier i.e. does not match [a-zA-Z][a-zA-Z0-9_]*, a valid identifier must be provided as fieldName.
* `qualifier_string` (Optional) - Qualifier string.
* `field_name` (Optional) - If the qualifier is not a valid BigQuery field identifier i.e. does not match [a-zA-Z][a-zA-Z0-9_]*, a valid identifier must be provided as the column field name and is used as field name in queries.
* `type` (Optional) - The type to convert the value in cells of this column. The values are expected to be encoded using HBase Bytes.toBytes function when using the BINARY encoding value. Following BigQuery types are allowed (case-sensitive): "BYTES", "STRING", "INTEGER", "FLOAT", "BOOLEAN", "JSON", Default type is "BYTES". 'type' can also be set at the column family level. However, the setting at this level takes precedence if 'type' is set at both levels.
* `encoding` (Optional) - The encoding of the values when the type is not STRING. Acceptable encoding values are: TEXT - indicates values are alphanumeric text strings. BINARY - indicates values are encoded using HBase Bytes.toBytes family of functions. 'encoding' can also be set at the column family level. However, the setting at this level takes precedence if 'encoding' is set at both levels.
* `only_read_latest` (Optional) - If this is set, only the latest version of value in this column are exposed. 'onlyReadLatest' can also be set at the column family level. However, the setting at this level takes precedence if 'onlyReadLatest' is set at both levels.

<a name="nested_json_options"></a>The `json_options` block supports:

* `encoding` (Optional) - The character encoding of the data. The supported values are UTF-8, UTF-16BE, UTF-16LE, UTF-32BE, and UTF-32LE. The default value is UTF-8.

<a name="nested_google_sheets_options"></a>The `google_sheets_options` block supports:

* `range` (Optional) - Range of a sheet to query from. Only used when
    non-empty. At least one of `range` or `skip_leading_rows` must be set.
    Typical format: "sheet_name!top_left_cell_id:bottom_right_cell_id"
    For example: "sheet1!A1:B20"

* `skip_leading_rows` (Optional) - The number of rows at the top of the sheet
    that BigQuery will skip when reading the data. At least one of `range` or
    `skip_leading_rows` must be set.

<a name="nested_hive_partitioning_options"></a>The `hive_partitioning_options` block supports:

* `mode` (Optional) - When set, what mode of hive partitioning to use when
    reading data. The following modes are supported.
    * AUTO: automatically infer partition key name(s) and type(s).
    * STRINGS: automatically infer partition key name(s). All types are
      Not all storage formats support hive partitioning. Requesting hive
      partitioning on an unsupported format will lead to an error.
      Currently supported formats are: JSON, CSV, ORC, Avro and Parquet.
    * CUSTOM: when set to `CUSTOM`, you must encode the partition key schema within the `source_uri_prefix` by setting `source_uri_prefix` to `gs://bucket/path_to_table/{key1:TYPE1}/{key2:TYPE2}/{key3:TYPE3}`.

* `require_partition_filter` - (Optional) If set to true, queries over this table
    require a partition filter that can be used for partition elimination to be
    specified.

* `source_uri_prefix` (Optional) - When hive partition detection is requested,
    a common for all source uris must be required. The prefix must end immediately
    before the partition key encoding begins. For example, consider files following
    this data layout. `gs://bucket/path_to_table/dt=2019-06-01/country=USA/id=7/file.avro`
    `gs://bucket/path_to_table/dt=2019-05-31/country=CA/id=3/file.avro` When hive
    partitioning is requested with either AUTO or STRINGS detection, the common prefix
    can be either of `gs://bucket/path_to_table` or `gs://bucket/path_to_table/`.
    Note that when `mode` is set to `CUSTOM`, you must encode the partition key schema within the `source_uri_prefix` by setting `source_uri_prefix` to `gs://bucket/path_to_table/{key1:TYPE1}/{key2:TYPE2}/{key3:TYPE3}`.

<a name="nested_avro_options"></a>The `avro_options` block supports:

* `use_avro_logical_types` (Optional) - If is set to true, indicates whether
    to interpret logical types as the corresponding BigQuery data type
    (for example, TIMESTAMP), instead of using the raw type (for example, INTEGER).

<a name="nested_parquet_options"></a>The `parquet_options` block supports:

* `enum_as_string` (Optional) - Indicates whether to infer Parquet ENUM logical type as STRING instead of BYTES by default.

* `enable_list_inference` (Optional) - Indicates whether to use schema inference specifically for Parquet LIST logical type.

<a name="nested_time_partitioning"></a>The `time_partitioning` block supports:

* `expiration_ms` -  (Optional) Number of milliseconds for which to keep the
    storage for a partition.

* `field` - (Optional) The field used to determine how to create a time-based
    partition. If time-based partitioning is enabled without this value, the
    table is partitioned based on the load time.

* `type` - (Required) The supported types are DAY, HOUR, MONTH, and YEAR,
    which will generate one partition per day, hour, month, and year, respectively.

* `require_partition_filter` - (Optional) If set to true, queries over this table
    require a partition filter that can be used for partition elimination to be
    specified. `require_partition_filter` is deprecated and will be removed in
    a future major release. Use the top level field with the same name instead.

<a name="nested_range_partitioning"></a>The `range_partitioning` block supports:

* `field` - (Required) The field used to determine how to create a range-based
    partition.

* `range` - (Required) Information required to partition based on ranges.
    Structure is [documented below](#nested_range).

<a name="nested_range"></a>The `range` block supports:

* `start` - (Required) Start of the range partitioning, inclusive.

* `end` - (Required) End of the range partitioning, exclusive.

* `interval` - (Required) The width of each range within the partition.

<a name="nested_view"></a>The `view` block supports:

* `query` - (Required) A query that BigQuery executes when the view is referenced.

* `use_legacy_sql` - (Optional) Specifies whether to use BigQuery's legacy SQL for this view.
    The default value is true. If set to false, the view will use BigQuery's standard SQL.

<a name="nested_materialized_view"></a>The `materialized_view` block supports:

* `query` - (Required) A query whose result is persisted.

* `enable_refresh` - (Optional) Specifies whether to use BigQuery's automatic refresh for this materialized view when the base table is updated.
    The default value is true.

* `refresh_interval_ms` - (Optional) The maximum frequency at which this materialized view will be refreshed.
    The default value is 1800000

* `allow_non_incremental_definition` - (Optional) Allow non incremental materialized view definition.
    The default value is false.

<a name="nested_encryption_configuration"></a>The `encryption_configuration` block supports the following arguments:

* `kms_key_name` - (Required) The self link or full name of a key which should be used to
    encrypt this table.  Note that the default bigquery service account will need to have
    encrypt/decrypt permissions on this key - you may want to see the
    `google_bigquery_default_service_account` datasource and the
    `google_kms_crypto_key_iam_binding` resource.

<a name="nested_table_constraints"></a>The `table_constraints` block supports:

* `primary_key` - (Optional) Represents the primary key constraint
    on a table's columns. Present only if the table has a primary key.
    The primary key is not enforced.
    Structure is [documented below](#nested_primary_key).

* `foreign_keys` - (Optional) Present only if the table has a foreign key.
    The foreign key is not enforced.
    Structure is [documented below](#nested_foreign_keys).

<a name="nested_primary_key"></a>The `primary_key` block supports:

* `columns`: (Required) The columns that are composed of the primary key constraint.

<a name="nested_foreign_keys"></a>The `foreign_keys` block supports:

* `name`: (Optional) Set only if the foreign key constraint is named.

* `referenced_table`: (Required) The table that holds the primary key
    and is referenced by this foreign key.
    Structure is [documented below](#nested_referenced_table).

* `column_references`: (Required) The pair of the foreign key column and primary key column.
    Structure is [documented below](#nested_column_references).

<a name="nested_referenced_table"></a>The `referenced_table` block supports:

* `project_id`: (Required) The ID of the project containing this table.

* `dataset_id`: (Required) The ID of the dataset containing this table.

* `table_id`: (Required) The ID of the table. The ID must contain only
	letters (a-z, A-Z), numbers (0-9), or underscores (_). The maximum
	length is 1,024 characters. Certain operations allow suffixing of
	the table ID with a partition decorator, such as
	sample_table$20190123.

<a name="nested_column_references"></a>The `column_references` block supports:

* `referencing_column`: (Required) The column that composes the foreign key.

* `referenced_column`: (Required) The column in the primary key that are
    referenced by the referencingColumn

<a name="nested_table_replication_info"></a>The `table_replication_info` block supports:

* `source_project_id` (Required) - The ID of the source project.

* `source_dataset_id` (Required) - The ID of the source dataset.

* `source_table_id` (Required) - The ID of the source materialized view.

* `replication_interval_ms` (Optional) - The interval at which the source
    materialized view is polled for updates. The default is 300000.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `id` - an identifier for the resource with format `projects/{{project}}/datasets/{{dataset}}/tables/{{name}}`

* `creation_time` - The time when this table was created, in milliseconds since the epoch.

* `etag` - A hash of the resource.

* `kms_key_version` - The self link or full name of the kms key version used to encrypt this table.

* `last_modified_time` - The time when this table was last modified, in milliseconds since the epoch.

* `location` - The geographic location where the table resides. This value is inherited from the dataset.

* `num_bytes` - The size of this table in bytes, excluding any data in the streaming buffer.

* `num_long_term_bytes` - The number of bytes in the table that are considered "long-term storage".

* `num_rows` - The number of rows of data in this table, excluding any data in the streaming buffer.

* `self_link` - The URI of the created resource.

* `type` - Describes the table type.

## Import

BigQuery tables can be imported using any of these accepted formats:

* `projects/{{project}}/datasets/{{dataset_id}}/tables/{{table_id}}`
* `{{project}}/{{dataset_id}}/{{table_id}}`
* `{{dataset_id}}/{{table_id}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import BigQuery tables using one of the formats above. For example:

```tf
import {
  id = "projects/{{project}}/datasets/{{dataset_id}}/tables/{{table_id}}"
  to = google_bigquery_table.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), BigQuery tables can be imported using one of the formats above. For example:

```
$ terraform import google_bigquery_table.default projects/{{project}}/datasets/{{dataset_id}}/tables/{{table_id}}
$ terraform import google_bigquery_table.default {{project}}/{{dataset_id}}/{{table_id}}
$ terraform import google_bigquery_table.default {{dataset_id}}/{{table_id}}
```
