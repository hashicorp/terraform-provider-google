---
layout: "google"
page_title: "Google: google_sql_database_instance"
sidebar_current: "docs-google-sql-database-instance"
description: |-
  Creates a new SQL database instance in Google Cloud SQL.
---

# google\_sql\_database\_instance

Creates a new Google SQL Database Instance. For more information, see the [official documentation](https://cloud.google.com/sql/),
or the [JSON API](https://cloud.google.com/sql/docs/admin-api/v1beta4/instances).

~> **NOTE on `google_sql_database_instance`:** - Second-generation instances include a
default 'root'@'%' user with no password. This user will be deleted by Terraform on
instance creation. You should use `google_sql_user` to define a custom user with
a restricted host and strong password.

## Example Usage

### SQL First Generation

```hcl
resource "google_sql_database_instance" "master" {
  name = "master-instance"
  database_version = "MYSQL_5_6"
  # First-generation instance regions are not the conventional
  # Google Compute Engine regions. See argument reference below.
  region = "us-central"

  settings {
    tier = "D0"
  }
}
```

### SQL Second generation

```hcl
resource "google_sql_database_instance" "master" {
  name = "master-instance"
  database_version = "POSTGRES_9_6"
  region = "us-central1"

  settings {
    # Second-generation instance tiers are based on the machine
    # type. See argument reference below.
    tier = "db-f1-micro"
  }
}
```

### Granular restriction of network access

```hcl
resource "google_compute_instance" "apps" {
  count        = 8
  name         = "apps-${count.index + 1}"
  machine_type = "f1-micro"
  
  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-1804-lts"
    }
  }

  network_interface {
    network = "default"

    access_config {
      // Ephemeral IP
    }
  }
}

data "null_data_source" "auth_netw_postgres_allowed_1" {
  count = "${length(google_compute_instance.apps.*.self_link)}"

  inputs = {
    name  = "apps-${count.index + 1}"
    value = "${element(google_compute_instance.apps.*.network_interface.0.access_config.0.nat_ip, count.index)}"
  }
}

data "null_data_source" "auth_netw_postgres_allowed_2" {
  count = 2

  inputs = {
    name  = "onprem-${count.index + 1}"
    value = "${element(list("192.168.1.2", "192.168.2.3"), count.index)}"
  }
}

resource "google_sql_database_instance" "postgres" {
  name = "postgres-instance"
  database_version = "POSTGRES_9_6"

  settings {
    tier = "db-f1-micro"
    
    ip_configuration {
      authorized_networks = [
        "${data.null_data_source.auth_netw_postgres_allowed_1.*.outputs}",
        "${data.null_data_source.auth_netw_postgres_allowed_2.*.outputs}",
      ]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required) The region the instance will sit in. Note, first-generation Cloud SQL instance
    regions do not line up with the Google Compute Engine (GCE) regions, and Cloud SQL is not
    available in all regions - choose from one of the options listed [here](https://cloud.google.com/sql/docs/mysql/instance-locations).
    A valid region must be provided to use this resource. If a region is not provided in the resource definition,
    the provider region will be used instead, but this will be an apply-time error for all first-generation
    instances *and* for second-generation instances if the provider region is not supported with Cloud SQL.
    If you choose not to provide the `region` argument for this resource, make sure you understand this.

* `settings` - (Required) The settings to use for the database. The
    configuration is detailed below.

- - -

* `database_version` - (Optional, Default: `MYSQL_5_6`) The MySQL version to
    use. Can be `MYSQL_5_6`, `MYSQL_5_7` or `POSTGRES_9_6` for second-generation
    instances, or `MYSQL_5_5` or `MYSQL_5_6` for first-generation instances.
    See [Second Generation Capabilities](https://cloud.google.com/sql/docs/1st-2nd-gen-differences)
    for more information. `POSTGRES_9_6` support is in beta.

* `name` - (Optional, Computed) The name of the instance. If the name is left
    blank, Terraform will randomly generate one when the instance is first
    created. This is done because after a name is used, it cannot be reused for
    up to [one week](https://cloud.google.com/sql/docs/delete-instance).

* `master_instance_name` - (Optional) The name of the instance that will act as
    the master in the replication setup. Note, this requires the master to have
    `binary_log_enabled` set, as well as existing backups.

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

* `replica_configuration` - (Optional) The configuration for replication. The
    configuration is detailed below.

The required `settings` block supports:

* `tier` - (Required) The machine tier (First Generation) or type (Second Generation) to use. See
    [tiers](https://cloud.google.com/sql/docs/admin-api/v1beta4/tiers) for more details and
    supported versions. Postgres supports only shared-core machine types such as `db-f1-micro`, and custom
    machine types such as `db-custom-2-13312`. See the
    [Custom Machine Type Documentation](https://cloud.google.com/compute/docs/instances/creating-instance-with-custom-machine-type#create)
    to learn about specifying custom machine types.

* `activation_policy` - (Optional) This specifies when the instance should be
    active. Can be either `ALWAYS`, `NEVER` or `ON_DEMAND`.

* `authorized_gae_applications` - (Optional) A list of Google App Engine (GAE)
    project names that are allowed to access this instance.

* `availability_type` - (Optional) This specifies whether a PostgreSQL instance
    should be set up for high availability (`REGIONAL`) or single zone (`ZONAL`).

* `crash_safe_replication` - (Optional) Specific to read instances, indicates
    when crash-safe replication flags are enabled.

* `disk_autoresize` - (Optional, Second Generation, Default: `true`) Configuration to increase storage size automatically.

* `disk_size` - (Optional, Second Generation, Default: `10`) The size of data disk, in GB. Size of a running instance cannot be reduced but can be increased.

* `disk_type` - (Optional, Second Generation, Default: `PD_SSD`) The type of data disk: PD_SSD or PD_HDD.

* `pricing_plan` - (Optional, First Generation) Pricing plan for this instance, can be one of
    `PER_USE` or `PACKAGE`.

* `replication_type` - (Optional) Replication type for this instance, can be one
    of `ASYNCHRONOUS` or `SYNCHRONOUS`.

* `user_labels` - (Optional) A set of key/value user label pairs to assign to the instance.

The optional `settings.database_flags` sublist supports:

* `name` - (Optional) Name of the flag.

* `value` - (Optional) Value of the flag.

The optional `settings.backup_configuration` subblock supports:

* `binary_log_enabled` - (Optional) True if binary logging is enabled. If
    `logging` is false, this must be as well. Cannot be used with Postgres.

* `enabled` - (Optional) True if backup configuration is enabled.

* `start_time` - (Optional) `HH:MM` format time indicating when backup
    configuration starts.

The optional `settings.ip_configuration` subblock supports:

* `ipv4_enabled` - (Optional) True if the instance should be assigned an IP
    address. The IPv4 address cannot be disabled for Second Generation instances.

* `require_ssl` - (Optional) True if mysqld should default to `REQUIRE X509`
    for users connecting over IP.

The optional `settings.ip_configuration.authorized_networks[]` sublist supports:

* `expiration_time` - (Optional) The [RFC 3339](https://tools.ietf.org/html/rfc3339)
  formatted date time string indicating when this whitelist expires.

* `name` - (Optional) A name for this whitelist entry.

* `value` - (Optional) A CIDR notation IPv4 or IPv6 address that is allowed to
    access this instance. Must be set even if other two attributes are not for
    the whitelist to become active.

The optional `settings.location_preference` subblock supports:

* `follow_gae_application` - (Optional) A GAE application whose zone to remain
    in. Must be in the same region as this instance.

* `zone` - (Optional) The preferred compute engine
    [zone](https://cloud.google.com/compute/docs/zones?hl=en).

The optional `settings.maintenance_window` subblock for Second Generation
instances declares a one-hour [maintenance window](https://cloud.google.com/sql/docs/instance-settings?hl=en#maintenance-window-2ndgen)
when an Instance can automatically restart to apply updates. The maintenance window is specified in UTC time. It supports:

* `day` - (Optional) Day of week (`1-7`), starting on Monday

* `hour` - (Optional) Hour of day (`0-23`), ignored if `day` not set

* `update_track` - (Optional) Receive updates earlier (`canary`) or later
(`stable`)

The optional `replica_configuration` block must have `master_instance_name` set
to work, cannot be updated, and supports:

* `ca_certificate` - (Optional) PEM representation of the trusted CA's x509
    certificate.

* `client_certificate` - (Optional) PEM representation of the slave's x509
    certificate.

* `client_key` - (Optional) PEM representation of the slave's private key. The
    corresponding public key in encoded in the `client_certificate`.

* `connect_retry_interval` - (Optional, Default: 60) The number of seconds
    between connect retries.

* `dump_file_path` - (Optional) Path to a SQL file in GCS from which slave
    instances are created. Format is `gs://bucket/filename`.

* `failover_target` - (Optional) Specifies if the replica is the failover target.
    If the field is set to true the replica will be designated as a failover replica.
    If the master instance fails, the replica instance will be promoted as
    the new master instance.

* `master_heartbeat_period` - (Optional) Time in ms between replication
    heartbeats.

* `password` - (Optional) Password for the replication connection.

* `sslCipher` - (Optional) Permissible ciphers for use in SSL encryption.

* `username` - (Optional) Username for replication connection.

* `verify_server_certificate` - (Optional) True if the master's common name
    value is checked during the SSL handshake.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `first_ip_address` - The first IPv4 address of the addresses assigned. This is
is to support accessing the [first address in the list in a terraform output](https://github.com/terraform-providers/terraform-provider-google/issues/912)
when the resource is configured with a `count`.

* `connection_name` - The connection name of the instance to be used in connection strings.

* `ip_address.0.ip_address` - The IPv4 address assigned.

* `ip_address.0.time_to_retire` - The time this IP address will be retired, in RFC
    3339 format.

* `self_link` - The URI of the created resource.

* `settings.version` - Used to make sure changes to the `settings` block are
    atomic.
    
* `server_ca_cert.0.cert` - The CA Certificate used to connect to the SQL Instance via SSL.

* `server_ca_cert.0.common_name` - The CN valid for the CA Cert.

* `server_ca_cert.0.create_time` - Creation time of the CA Cert.

* `server_ca_cert.0.expiration_time` - Expiration time of the CA Cert.

* `server_ca_cert.0.sha1_fingerprint` - SHA Fingerprint of the CA Cert.

* `service_account_email_address` - The service account email address assigned to the
instance. This property is applicable only to Second Generation instances.

## Timeouts

`google_sql_database_instance` provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 10 minutes.
- `update` - Default is 10 minutes.
- `delete` - Default is 10 minutes.

## Import

Database instances can be imported using one of any of these accepted formats:

```
$ terraform import google_sql_database_instance.master projects/{{project}}/instances/{{name}}
$ terraform import google_sql_database_instance.master {{project}}/{{name}}
$ terraform import google_sql_database_instance.master {{name}}

```

~> **NOTE:** Some fields (such as `replica_configuration`) won't show a diff if they are unset in
config and set on the server.
When importing, double-check that your config has all the fields set that you expect- just seeing
no diff isn't sufficient to know that your config could reproduce the imported resource.
