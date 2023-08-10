---
subcategory: "AlloyDB"
description: |-
  Fetches the list of supported alloydb database flags in a location.
---

# google\_alloydb\_supported\_database\_flags

Use this data source to get information about the supported alloydb database flags in a location.

## Example Usage


```hcl
data "google_alloydb_supported_database_flags" "qa" {
    location = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `location` - (required) The canonical id of the location. For example: `us-east1`.

* `project` - (optional) The ID of the project.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `supported_database_flags` - Contains a list of `flag`, which contains the details about a particular flag.

A `flag` object would contain the following fields:-

* `name` - The name of the flag resource, following Google Cloud conventions, e.g.: * projects/{project}/locations/{location}/flags/{flag} This field currently has no semantic meaning.

* `flag_name` - The name of the database flag, e.g. "max_allowed_packets". The is a possibly key for the Instance.database_flags map field.

* `value_type` - ValueType describes the semantic type of the value that the flag accepts. Regardless of the ValueType, the Instance.database_flags field accepts the stringified version of the value, i.e. "20" or "3.14". The supported values are `VALUE_TYPE_UNSPECIFIED`, `STRING`, `INTEGER`, `FLOAT` and `NONE`.

* `accepts_multiple_values` - Whether the database flag accepts multiple values. If true, a comma-separated list of stringified values may be specified.

* `supported_db_versions` - Major database engine versions for which this flag is supported. The supported values are `POSTGRES_14` and `DATABASE_VERSION_UNSPECIFIED`.

* `requires_db_restart` - Whether setting or updating this flag on an Instance requires a database restart. If a flag that requires database restart is set, the backend will automatically restart the database (making sure to satisfy any availability SLO's).

* `string_restrictions` - Restriction on `STRING` type value. The list of allowed values, if bounded. This field will be empty if there is a unbounded number of allowed values.

* `integer_restrictions` - Restriction on `INTEGER` type value. Specifies the minimum value and the maximum value that can be specified, if applicable.

-> **Note** `string_restrictions` and `integer_restrictions` are part of the union field `restrictions`. The restrictions on the flag value per type. `restrictions` can be either `string_restrictions` or `integer_restrictions` but not both.
