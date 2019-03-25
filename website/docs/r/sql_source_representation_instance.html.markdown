---
layout: "google"
page_title: "Google: google_sql_source_representation_instance"
sidebar_current: "docs-google-sql-source-representation-instance"
description: |-
  Creates a new SQL Source Representation Instance, which holds metadata about
  on-premises MySQL instances. This enables Google CloudSQL MySQL instances to
  become replicas of on-premises databases.

  Replicas are configured by setting the source representation instance name in
  their `master_instance_name` value.

---

# google\_sql\_database\_instance

Creates a new source representation Instance for Google CloudSQL.

For more information, see the [official documentation](https://cloud.google.com/sql/),
or the [JSON API](https://cloud.google.com/sql/docs/admin-api/v1beta4/instances).

~> **NOTE on `google_sql_source_representation_instance`:** - This feature is
currently only available for MySQL databases.

## Example Usage

```hcl
resource "google_sql_source_representation_instance" "on_premises_master" {
  name = "master-instance"
  database_version = "MYSQL_5_6"
  region = "us-central1"
  host = "1.2.3.4"
  port = "3306"
}
```

## Argument Reference

The following arguments are supported:

* `host` - (Required) The host of the on-premises instance.

* `region` - (Required) The region the instances replicating this database will sit in.
    Note, Cloud SQL is not available in all regions - choose from one of the
    options listed [here](https://cloud.google.com/sql/docs/mysql/instance-locations).
    A valid region must be provided to use this resource. If a region is not
    provided in the resource definition, the provider region will be used
    instead, but this will be an apply-time error if the provider region is not
    supported with Cloud SQL.
    If you choose not to provide the `region` argument for this resource, make
    sure you understand this.

- - -

* `database_version` - (Optional, Default: `MYSQL_5_7`) The MySQL version to use.
    Can be `MYSQL_5_6` or `MYSQL_5_7`.

* `name` - (Optional, Computed) The name of the instance. If the name is left
    blank, Terraform will randomly generate one when the instance is first
    created. Note that unlike regular CloudSQL instances, there are no
    limitations on source representation instance name reuse.

* `port` - (Optional, Default: `3306`) The port of the on-premises instance.

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `self_link` - The URI of the created resource.

## Timeouts

`google_sql_source_representation_instance` provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 10 minutes.
- `delete` - Default is 10 minutes.

## Import

Database instances can be imported using one of any of these accepted formats:

```
$ terraform import google_sql_source_representation_instance.master projects/{{project}}/instances/{{name}}
$ terraform import google_sql_source_representation_instance.master {{project}}/{{name}}
$ terraform import google_sql_source_representation_instance.master {{name}}
```
