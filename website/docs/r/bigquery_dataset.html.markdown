---
layout: "google"
page_title: "Google: google_bigquery_dataset"
sidebar_current: "docs-google-bigquery-dataset"
description: |-
  Creates a dataset resource for Google BigQuery.
---

# google_bigquery_dataset

Creates a dataset resource for Google BigQuery. For more information see
[the official documentation](https://cloud.google.com/bigquery/docs/) and
[API](https://cloud.google.com/bigquery/docs/reference/rest/v2/datasets).


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

  access {
    role   = "READER"
    domain = "example.com"
  }
  access {
    role           = "WRITER"
    group_by_email = "writers@example.com"
  }
}
```

## Argument Reference

The following arguments are supported:

* `dataset_id` - (Required) A unique ID for the resource.
    Changing this forces a new resource to be created.

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

* `friendly_name` - (Optional) A descriptive name for the dataset.

* `description` - (Optional) A user-friendly description of the dataset.

* `delete_contents_on_destroy` - (Optional) If set to `true`, delete all the
    tables in the dataset when destroying the resource; otherwise, destroying
    the resource will fail if tables are present.

* `location` - (Optional) The geographic location where the dataset should reside.
    See [official docs](https://cloud.google.com/bigquery/docs/dataset-locations).

    There are two types of locations, regional or multi-regional.
    A regional location is a specific geographic place, such as Tokyo, and a
    multi-regional location is a large geographic area, such as the United States,
    that contains at least two geographic places

    Possible regional values include: `asia-east1`, `asia-northeast1`, `asia-southeast1`
     `australia-southeast1`, `europe-north1`, `europe-west2` and `us-east4`.

    Possible multi-regional values:`EU` and `US`.

    The default value is multi-regional location `US`.
    Changing this forces a new resource to be created.

* `default_partition_expiration_ms` - (Optional) The default partition expiration
    for all partitioned tables in the dataset, in milliseconds.

    Once this property is set, all newly-created partitioned tables in the dataset
    will have an expirationMs property in the timePartitioning settings set to this
    value, and changing the value will only affect new tables, not existing ones.
    The storage in a partition will have an expiration time of its partition time
    plus this value. Setting this property overrides the use of
    `defaultTableExpirationMs` for partitioned tables: only one of
    `defaultTableExpirationMs` and `defaultPartitionExpirationMs` will be used for
    any new partitioned table. If you provide an explicit
    `timePartitioning.expirationMs` when creating or updating a partitioned table,
    that value takes precedence over the default partition expiration time
    indicated by this property.

* `default_table_expiration_ms` - (Optional) The default lifetime of all
    tables in the dataset, in milliseconds. The minimum value is 3600000
    milliseconds (one hour).

    Once this property is set, all newly-created
    tables in the dataset will have an expirationTime property set to the
    creation time plus the value in this property, and changing the value
    will only affect new tables, not existing ones. When the
    expirationTime for a given table is reached, that table will be
    deleted automatically. If a table's expirationTime is modified or
    removed before the table expires, or if you provide an explicit
    expirationTime when creating a table, that value takes precedence
    over the default expiration time indicated by this property.

* `labels` - (Optional) A mapping of labels to assign to the resource.

* `access` - (Optional) An array of objects that define dataset access for
    one or more entities. Structure is documented below.

The `access` block supports the following fields (exactly one of `domain`,
`group_by_email`, `special_group`, `user_by_email`, or `view` must be set,
even though they are marked optional):

* `role` - (Required unless `view` is set) Describes the rights granted to
    the user specified by the other member of the access object. 
    Primitive, Predefined and custom roles are supported.
    Predefined roles that have equivalent primitive roles are swapped 
    by the API to their Primitive counterparts, and will show a diff post-create. 
    See [official docs](https://cloud.google.com/bigquery/docs/access-control).

* `domain` - (Optional) A domain to grant access to.

* `group_by_email` - (Optional) An email address of a Google Group to grant
    access to.

* `special_group` - (Optional) A special group to grant access to.
  Possible values include:
  * `projectOwners`: Owners of the enclosing project.
  * `projectReaders`: Readers of the enclosing project.
  * `projectWriters`: Writers of the enclosing project.
  * `allAuthenticatedUsers`: All authenticated BigQuery users.

* `user_by_email` - (Optional) An email address of a user to grant access to.

* `view` - (Optional) A view from a different dataset to grant access to.
    Queries executed against that view will have read access to tables in this
    dataset. The role field is not required when this field is set. If that
    view is updated by any user, access to the view needs to be granted again
    via an update operation. Structure is documented below.

The `access.view` block supports:

* `dataset_id` - (Required) The ID of the dataset containing this table.

* `project_id` - (Required) The ID of the project containing this table.

* `table_id` - (Required) The ID of the table.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `self_link` - The URI of the created resource.

* `etag` - A hash of the resource.

* `creation_time` - The time when this dataset was created, in milliseconds since the epoch.

* `last_modified_time` -  The date when this dataset or any of its tables was last modified,
  in milliseconds since the epoch.

## Import

BigQuery datasets can be imported using the `project` and `dataset_id`, e.g.

```
$ terraform import google_bigquery_dataset.default gcp-project:foo
```
