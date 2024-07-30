---
subcategory: "Cloud SQL"
description: |-
  Creates a new SQL database instance in Google Cloud SQL.
---

# google_sql_database_instance

Creates a new Google SQL Database Instance. For more information, see the [official documentation](https://cloud.google.com/sql/),
or the [JSON API](https://cloud.google.com/sql/docs/admin-api/v1beta4/instances).

~> **NOTE on `google_sql_database_instance`:** - Second-generation instances include a
default 'root'@'%' user with no password. This user will be deleted by Terraform on
instance creation. You should use `google_sql_user` to define a custom user with
a restricted host and strong password.

-> **Note**: On newer versions of the provider, you must explicitly set `deletion_protection=false`
(and run `terraform apply` to write the field to state) in order to destroy an instance.
It is recommended to not set this field (or set it to true) until you're ready to destroy the instance and its databases.

## Example Usage

### SQL Second Generation Instance

```hcl
resource "google_sql_database_instance" "main" {
  name             = "main-instance"
  database_version = "POSTGRES_15"
  region           = "us-central1"

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

resource "random_id" "db_name_suffix" {
  byte_length = 4
}

locals {
  onprem = ["192.168.1.2", "192.168.2.3"]
}

resource "google_sql_database_instance" "postgres" {
  name             = "postgres-instance-${random_id.db_name_suffix.hex}"
  database_version = "POSTGRES_15"

  settings {
    tier = "db-f1-micro"

    ip_configuration {

      dynamic "authorized_networks" {
        for_each = google_compute_instance.apps
        iterator = apps

        content {
          name  = apps.value.name
          value = apps.value.network_interface.0.access_config.0.nat_ip
        }
      }

      dynamic "authorized_networks" {
        for_each = local.onprem
        iterator = onprem

        content {
          name  = "onprem-${onprem.key}"
          value = onprem.value
        }
      }
    }
  }
}
```

### Private IP Instance
~> **NOTE:** For private IP instance setup, note that the `google_sql_database_instance` does not actually interpolate values from `google_service_networking_connection`. You must explicitly add a `depends_on`reference as shown below.

```hcl
resource "google_compute_network" "private_network" {
  provider = google-beta

  name = "private-network"
}

resource "google_compute_global_address" "private_ip_address" {
  provider = google-beta

  name          = "private-ip-address"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.private_network.id
}

resource "google_service_networking_connection" "private_vpc_connection" {
  provider = google-beta

  network                 = google_compute_network.private_network.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_address.name]
}

resource "random_id" "db_name_suffix" {
  byte_length = 4
}

resource "google_sql_database_instance" "instance" {
  provider = google-beta

  name             = "private-instance-${random_id.db_name_suffix.hex}"
  region           = "us-central1"
  database_version = "MYSQL_5_7"

  depends_on = [google_service_networking_connection.private_vpc_connection]

  settings {
    tier = "db-f1-micro"
    ip_configuration {
      ipv4_enabled                                  = false
      private_network                               = google_compute_network.private_network.self_link
      enable_private_path_for_google_cloud_services = true
    }
  }
}

provider "google-beta" {
  region = "us-central1"
  zone   = "us-central1-a"
}
```

### ENTERPRISE_PLUS Instance with data_cache_config

```hcl
resource "google_sql_database_instance" "main" {
  name             = "enterprise-plus-main-instance"
  database_version = "MYSQL_8_0_31"
  settings {
    tier    = "db-perf-optimized-N-2"
    edition = "ENTERPRISE_PLUS"
    data_cache_config {
        data_cache_enabled = true
    }
  }
}
```

### Cloud SQL Instance with PSC connectivity

```hcl
resource "google_sql_database_instance" "main" {
  name             = "psc-enabled-main-instance"
  database_version = "MYSQL_8_0"
  settings {
    tier    = "db-f1-micro"
    ip_configuration {
      psc_config {
        psc_enabled = true
        allowed_consumer_projects = ["allowed-consumer-project-name"]
      }
      ipv4_enabled = false
    }
    backup_configuration {
      enabled = true
      binary_log_enabled = true
    }
    availability_type = "REGIONAL"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region the instance will sit in. If a region is not provided in the resource definition,
    the provider region will be used instead.

- - -

* `settings` - (Optional) The settings to use for the database. The
    configuration is detailed below. Required if `clone` is not set.

* `database_version` - (Required) The MySQL, PostgreSQL or
SQL Server version to use. Supported values include `MYSQL_5_6`,
`MYSQL_5_7`, `MYSQL_8_0`, `POSTGRES_9_6`,`POSTGRES_10`, `POSTGRES_11`,
`POSTGRES_12`, `POSTGRES_13`, `POSTGRES_14`, `POSTGRES_15`, `SQLSERVER_2017_STANDARD`,
`SQLSERVER_2017_ENTERPRISE`, `SQLSERVER_2017_EXPRESS`, `SQLSERVER_2017_WEB`.
`SQLSERVER_2019_STANDARD`, `SQLSERVER_2019_ENTERPRISE`, `SQLSERVER_2019_EXPRESS`,
`SQLSERVER_2019_WEB`.
[Database Version Policies](https://cloud.google.com/sql/docs/db-versions)
includes an up-to-date reference of supported versions.

* `name` - (Optional, Computed) The name of the instance. If the name is left
    blank, Terraform will randomly generate one when the instance is first
    created. This is done because after a name is used, it cannot be reused for
    up to [one week](https://cloud.google.com/sql/docs/delete-instance).

* `maintenance_version`  - (Optional) The current software version on the instance. This attribute can not be set during creation. Refer to `available_maintenance_versions` attribute to see what `maintenance_version` are available for upgrade. When this attribute gets updated, it will cause an instance restart. Setting a `maintenance_version` value that is older than the current one on the instance will be ignored.

* `master_instance_name` - (Optional) The name of the existing instance that will
    act as the master in the replication setup. Note, this requires the master to
    have `binary_log_enabled` set, as well as existing backups.

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

* `replica_configuration` - (Optional) The configuration for replication. The
    configuration is detailed below. Valid only for MySQL instances.

* `root_password` - (Optional) Initial root password. Can be updated. Required for MS SQL Server.

* `encryption_key_name` - (Optional)
    The full path to the encryption key used for the CMEK disk encryption.  Setting
    up disk encryption currently requires manual steps outside of Terraform.
    The provided key must be in the same region as the SQL instance.  In order
    to use this feature, a special kind of service account must be created and
    granted permission on this key.  This step can currently only be done
    manually, please see [this step](https://cloud.google.com/sql/docs/mysql/configure-cmek#service-account).
    That service account needs the `Cloud KMS > Cloud KMS CryptoKey Encrypter/Decrypter` role on your
    key - please see [this step](https://cloud.google.com/sql/docs/mysql/configure-cmek#grantkey).

* `deletion_protection` - (Optional) Whether Terraform will be prevented from destroying the instance.
    When the field is set to true or unset in Terraform state, a `terraform apply`
    or `terraform destroy` that would delete the instance will fail.
    When the field is set to false, deleting the instance is allowed.

  ~> **NOTE:** This flag only protects instances from deletion within Terraform. To protect your instances from accidental deletion across all surfaces (API, gcloud, Cloud Console and Terraform), use the API flag `settings.deletion_protection_enabled`.

* `restore_backup_context` - (optional) The context needed to restore the database to a backup run. This field will
    cause Terraform to trigger the database to restore from the backup run indicated. The configuration is detailed below.
    **NOTE:** Restoring from a backup is an imperative action and not recommended via Terraform. Adding or modifying this
    block during resource creation/update will trigger the restore action after the resource is created/updated.

* `clone` - (Optional) The context needed to create this instance as a clone of another instance. When this field is set during
    resource creation, Terraform will attempt to clone another instance as indicated in the context. The
    configuration is detailed below.

The `settings` block supports:

* `tier` - (Required) The machine type to use. See [tiers](https://cloud.google.com/sql/docs/admin-api/v1beta4/tiers)
    for more details and supported versions. Postgres supports only shared-core machine types,
    and custom machine types such as `db-custom-2-13312`. See the [Custom Machine Type Documentation](https://cloud.google.com/compute/docs/instances/creating-instance-with-custom-machine-type#create) to learn about specifying custom machine types.

* `edition` - (Optional) The edition of the instance, can be `ENTERPRISE` or `ENTERPRISE_PLUS`.

* `user_labels` - (Optional) A set of key/value user label pairs to assign to the instance.

* `activation_policy` - (Optional) This specifies when the instance should be
    active. Can be either `ALWAYS`, `NEVER` or `ON_DEMAND`.

* `availability_type` - (Optional) The availability type of the Cloud SQL
  instance, high availability (`REGIONAL`) or single zone (`ZONAL`).' For all instances, ensure that
  `settings.backup_configuration.enabled` is set to `true`.
  For MySQL instances, ensure that `settings.backup_configuration.binary_log_enabled` is set to `true`.
  For Postgres and SQL Server instances, ensure that `settings.backup_configuration.point_in_time_recovery_enabled`
  is set to `true`. Defaults to `ZONAL`.

* `collation` - (Optional) The name of server instance collation.

* `connector_enforcement` - (Optional) Specifies if connections must use Cloud SQL connectors.

* `deletion_protection_enabled` - (Optional) Enables deletion protection of an instance at the GCP level. Enabling this protection will guard against accidental deletion across all surfaces (API, gcloud, Cloud Console and Terraform) by enabling the [GCP Cloud SQL instance deletion protection](https://cloud.google.com/sql/docs/postgres/deletion-protection). Terraform provider support was introduced in version 4.48.0. Defaults to `false`.

* `enable_google_ml_integration` - (Optional) Enables [Cloud SQL instances to connect to Vertex AI](https://cloud.google.com/sql/docs/postgres/integrate-cloud-sql-with-vertex-ai) and pass requests for real-time predictions and insights. Defaults to `false`.

* `enable_dataplex_integration` - (Optional) Enables [Cloud SQL instance integration with Dataplex](https://cloud.google.com/sql/docs/mysql/dataplex-catalog-integration). MySQL, Postgres and SQL Server instances are supported for this feature. Defaults to `false`.

* `disk_autoresize` - (Optional) Enables auto-resizing of the storage size. Defaults to `true`. Note that if `disk_size` is set, future `terraform apply` calls will attempt to delete the instance in order to resize the disk to the value specified in disk_size if it has been resized. To avoid this, ensure that `lifecycle.ignore_changes` is applied to `disk_size`.

* `disk_autoresize_limit` - (Optional) The maximum size to which storage capacity can be automatically increased. The default value is 0, which specifies that there is no limit.

* `disk_size` - (Optional) The size of data disk, in GB. Size of a running instance cannot be reduced but can be increased. The minimum value is 10GB. Note that this value will override the resizing from `disk_autoresize` if that feature is enabled. To avoid this, set `lifecycle.ignore_changes` on this field.

* `disk_type` - (Optional) The type of data disk: PD_SSD or PD_HDD. Defaults to `PD_SSD`.

* `pricing_plan` - (Optional) Pricing plan for this instance, can only be `PER_USE`.

* `time_zone` - (Optional) The time_zone to be used by the database engine (supported only for SQL Server), in SQL Server timezone format.

The optional `settings.advanced_machine_features` subblock supports:

* `threads_per_core` - (Optional) The number of threads per core. The value of this flag can be 1 or 2. To disable SMT, set this flag to 1. Only available in Cloud SQL for SQL Server instances. See [smt](https://cloud.google.com/sql/docs/sqlserver/create-instance#smt-create-instance) for more details.

The optional `settings.database_flags` sublist supports:

* `name` - (Required) Name of the flag.

* `value` - (Required) Value of the flag.

The optional `settings.active_directory_config` subblock supports:

* `domain` - (Required) The domain name for the active directory (e.g., mydomain.com).
    Can only be used with SQL Server.

The optional `settings.data_cache_config` subblock supports:

* `data_cache_enabled` - (Optional) Whether data cache is enabled for the instance. Defaults to `false`. Can be used with MYSQL and PostgreSQL only.

The optional `settings.deny_maintenance_period` subblock supports:

* `end_date` - (Required) "deny maintenance period" end date. If the year of the end date is empty, the year of the start date also must be empty. In this case, it means the no maintenance interval recurs every year. The date is in format yyyy-m-dd (the month is without leading zeros)i.e., 2020-1-01, or 2020-11-01, or mm-dd, i.e., 11-01

* `start_date` - (Required) "deny maintenance period" start date. If the year of the start date is empty, the year of the end date also must be empty. In this case, it means the deny maintenance period recurs every year. The date is in format yyyy-m-dd (the month is without leading zeros)i.e., 2020-1-01, or 2020-11-01, or mm-dd, i.e., 11-01

* `time` - (Required) Time in UTC when the "deny maintenance period" starts on startDate and ends on endDate. The time is in format: HH:mm:SS, i.e., 00:00:00

The optional `settings.sql_server_audit_config` subblock supports:

* `bucket` - (Optional) The name of the destination bucket (e.g., gs://mybucket).

* `upload_interval` - (Optional) How often to upload generated audit files. A duration in seconds with up to nine fractional digits, terminated by 's'. Example: "3.5s".

* `retention_interval` - (Optional) How long to keep generated audit files. A duration in seconds with up to nine fractional digits, terminated by 's'. Example: "3.5s".

The optional `settings.backup_configuration` subblock supports:

* `binary_log_enabled` - (Optional) True if binary logging is enabled.
    Can only be used with MySQL.

* `enabled` - (Optional) True if backup configuration is enabled.

* `start_time` - (Optional) `HH:MM` format time indicating when backup
    configuration starts.
* `point_in_time_recovery_enabled` - (Optional) True if Point-in-time recovery is enabled. Will restart database if enabled after instance creation. Valid only for PostgreSQL and SQL Server instances.

* `location` - (Optional) The region where the backup will be stored

* `transaction_log_retention_days` - (Optional) The number of days of transaction logs we retain for point in time restore, from 1-7. For PostgreSQL Enterprise Plus instances, the number of days of retained transaction logs can be set from 1 to 35.

* `backup_retention_settings` - (Optional) Backup retention settings. The configuration is detailed below.

The optional `settings.backup_configuration.backup_retention_settings` subblock supports:

* `retained_backups` - (Optional) Depending on the value of retention_unit, this is used to determine if a backup needs to be deleted. If retention_unit
  is 'COUNT', we will retain this many backups.

* `retention_unit` - (Optional) The unit that 'retained_backups' represents. Defaults to `COUNT`.

The optional `settings.ip_configuration` subblock supports:

* `ipv4_enabled` - (Optional) Whether this Cloud SQL instance should be assigned
a public IPV4 address. At least `ipv4_enabled` must be enabled or a
`private_network` must be configured.

* `private_network` - (Optional) The VPC network from which the Cloud SQL
instance is accessible for private IP. For example, projects/myProject/global/networks/default.
Specifying a network enables private IP.
At least `ipv4_enabled` must be enabled or a `private_network` must be configured.
This setting can be updated, but it cannot be removed after it is set.

* `require_ssl` - (Optional, Deprecated) Whether SSL connections over IP are enforced or not. To change this field, also set the corresponding value in `ssl_mode`. It will be fully deprecated in a future major release. For now, please use `ssl_mode` with a compatible `require_ssl` value instead.

* `ssl_mode` - (Optional) Specify how SSL connection should be enforced in DB connections. This field provides more SSL enforcement options compared to `require_ssl`. To change this field, also set the correspoding value in `require_ssl`.
    * For PostgreSQL instances, the value pairs are listed in the [API reference doc](https://cloud.google.com/sql/docs/postgres/admin-api/rest/v1beta4/instances#ipconfiguration) for `ssl_mode` field.
    * For MySQL instances, use the same value pairs as the PostgreSQL instances.
    * For SQL Server instances, set it to `ALLOW_UNENCRYPTED_AND_ENCRYPTED` when `require_ssl=false` and `ENCRYPTED_ONLY` otherwise.

* `allocated_ip_range` - (Optional) The name of the allocated ip range for the private ip CloudSQL instance. For example: "google-managed-services-default". If set, the instance ip will be created in the allocated range. The range name must comply with [RFC 1035](https://datatracker.ietf.org/doc/html/rfc1035). Specifically, the name must be 1-63 characters long and match the regular expression [a-z]([-a-z0-9]*[a-z0-9])?.

* `enable_private_path_for_google_cloud_services` - (Optional) Whether Google Cloud services such as BigQuery are allowed to access data in this Cloud SQL instance over a private IP connection. SQLSERVER database type is not supported.

The optional `settings.ip_configuration.authorized_networks[]` sublist supports:

* `expiration_time` - (Optional) The [RFC 3339](https://tools.ietf.org/html/rfc3339)
  formatted date time string indicating when this whitelist expires.

* `name` - (Optional) A name for this whitelist entry.

* `value` - (Required) A CIDR notation IPv4 or IPv6 address that is allowed to
    access this instance. Must be set even if other two attributes are not for
    the whitelist to become active.

The optional `settings.ip_configuration.psc_config` sublist supports:

* `psc_enabled` - (Optional) Whether PSC connectivity is enabled for this instance.

* `allowed_consumer_projects` - (Optional) List of consumer projects that are allow-listed for PSC connections to this instance. This instance can be connected to with PSC from any network in these projects. Each consumer project in this list may be represented by a project number (numeric) or by a project id (alphanumeric).

The optional `settings.location_preference` subblock supports:

* `follow_gae_application` - (Optional) A GAE application whose zone to remain
    in. Must be in the same region as this instance.

* `zone` - (Optional) The preferred compute engine
    [zone](https://cloud.google.com/compute/docs/zones?hl=en).

* `secondary_zone` - (Optional) The preferred Compute Engine zone for the secondary/failover.

The optional `settings.maintenance_window` subblock for instances declares a one-hour
[maintenance window](https://cloud.google.com/sql/docs/instance-settings?hl=en#maintenance-window-2ndgen)
when an Instance can automatically restart to apply updates. The maintenance window is specified in UTC time. It supports:

* `day` - (Optional) Day of week (`1-7`), starting on Monday

* `hour` - (Optional) Hour of day (`0-23`), ignored if `day` not set

* `update_track` - (Optional) Receive updates after one week (`canary`) or after two weeks (`stable`) or after five weeks (`week5`) of notification.

The optional `settings.insights_config` subblock for instances declares Query Insights([MySQL](https://cloud.google.com/sql/docs/mysql/using-query-insights), [PostgreSQL](https://cloud.google.com/sql/docs/postgres/using-query-insights)) configuration. It contains:

* `query_insights_enabled` - True if Query Insights feature is enabled.

* `query_string_length` - Maximum query length stored in bytes. Between 256 and 4500. Default to 1024. Higher query lengths are more useful for analytical queries, but they also require more memory. Changing the query length requires you to restart the instance. You can still add tags to queries that exceed the length limit.

* `record_application_tags` - True if Query Insights will record application tags from query when enabled.

* `record_client_address` - True if Query Insights will record client address when enabled.

* `query_plans_per_minute` - Number of query execution plans captured by Insights per minute for all queries combined. Between 0 and 20. Default to 5.

The optional `settings.password_validation_policy` subblock for instances declares [Password Validation Policy](https://cloud.google.com/sql/docs/postgres/built-in-authentication) configuration. It contains:

* `min_length` - Specifies the minimum number of characters that the password must have.

* `complexity` - Checks if the password is a combination of lowercase, uppercase, numeric, and non-alphanumeric characters.

* `reuse_interval` - Specifies the number of previous passwords that you can't reuse.

* `disallow_username_substring` - Prevents the use of the username in the password.

* `password_change_interval` - Specifies the minimum duration after which you can change the password.

* `enable_password_policy` - Enables or disable the password validation policy.

The optional `replica_configuration` block must have `master_instance_name` set
to work, cannot be updated, and supports:

* `ca_certificate` - (Optional) PEM representation of the trusted CA's x509
    certificate.

* `client_certificate` - (Optional) PEM representation of the replica's x509
    certificate.

* `client_key` - (Optional) PEM representation of the replica's private key. The
    corresponding public key in encoded in the `client_certificate`.

* `connect_retry_interval` - (Optional) The number of seconds
    between connect retries. MySQL's default is 60 seconds.

* `dump_file_path` - (Optional) Path to a SQL file in GCS from which replica
    instances are created. Format is `gs://bucket/filename`.

* `failover_target` - (Optional) Specifies if the replica is the failover target.
    If the field is set to true the replica will be designated as a failover replica.
    If the master instance fails, the replica instance will be promoted as
    the new master instance.
  ~> **NOTE:** Not supported for Postgres database.

* `master_heartbeat_period` - (Optional) Time in ms between replication
    heartbeats.

* `password` - (Optional) Password for the replication connection.

* `ssl_cipher` - (Optional) Permissible ciphers for use in SSL encryption.

* `username` - (Optional) Username for replication connection.

* `verify_server_certificate` - (Optional) True if the master's common name
    value is checked during the SSL handshake.

The optional `clone` block supports:

* `source_instance_name` - (Required) Name of the source instance which will be cloned.

* `point_in_time` -  (Optional) The timestamp of the point in time that should be restored.

    A timestamp in RFC3339 UTC "Zulu" format, with nanosecond resolution and up to nine fractional digits. Examples: "2014-10-02T15:01:23Z" and "2014-10-02T15:01:23.045123456Z".

* `preferred_zone` - (Optional) (Point-in-time recovery for PostgreSQL only) Clone to an instance in the specified zone. If no zone is specified, clone to the same zone as the source instance. [clone-unavailable-instance](https://cloud.google.com/sql/docs/postgres/clone-instance#clone-unavailable-instance)

* `database_names` - (Optional) (SQL Server only, use with `point_in_time`) Clone only the specified databases from the source instance. Clone all databases if empty.

* `allocated_ip_range` -  (Optional) The name of the allocated ip range for the private ip CloudSQL instance. For example: "google-managed-services-default". If set, the cloned instance ip will be created in the allocated range. The range name must comply with [RFC 1035](https://tools.ietf.org/html/rfc1035). Specifically, the name must be 1-63 characters long and match the regular expression [a-z]([-a-z0-9]*[a-z0-9])?.

The optional `restore_backup_context` block supports:
**NOTE:** Restoring from a backup is an imperative action and not recommended via Terraform. Adding or modifying this
block during resource creation/update will trigger the restore action after the resource is created/updated.

* `backup_run_id` - (Required) The ID of the backup run to restore from.

* `instance_id` - (Optional) The ID of the instance that the backup was taken from. If left empty,
    this instance's ID will be used.

* `project` - (Optional) The full project ID of the source instance.`

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `self_link` - The URI of the created resource.

* `connection_name` - The connection name of the instance to be used in
connection strings. For example, when connecting with [Cloud SQL Proxy](https://cloud.google.com/sql/docs/mysql/connect-admin-proxy).

* `dsn_name` - The DNS name of the instance. See [Connect to an instance using Private Service Connect](https://cloud.google.com/sql/docs/mysql/configure-private-service-connect#view-summary-information-cloud-sql-instances-psc-enabled) for more details.

* `service_account_email_address` - The service account email address assigned to the
instance.

* `ip_address.0.ip_address` - The IPv4 address assigned.

* `ip_address.0.time_to_retire` - The time this IP address will be retired, in RFC
    3339 format.

* `ip_address.0.type` - The type of this IP address.

  * A `PRIMARY` address is an address that can accept incoming connections.

  * An `OUTGOING` address is the source address of connections originating from the instance, if supported.

  * A `PRIVATE` address is an address for an instance which has been configured to use private networking see: [Private IP](https://cloud.google.com/sql/docs/mysql/private-ip).

* `first_ip_address` - The first IPv4 address of any type assigned. This is to
support accessing the [first address in the list in a terraform output](https://github.com/hashicorp/terraform-provider-google/issues/912)
when the resource is configured with a `count`.

* `available_maintenance_versions`  - The list of all maintenance versions applicable on the instance.

* `public_ip_address` - The first public (`PRIMARY`) IPv4 address assigned. This is
a workaround for an [issue fixed in Terraform 0.12](https://github.com/hashicorp/terraform/issues/17048)
but also provides a convenient way to access an IP of a specific type without
performing filtering in a Terraform config.

* `private_ip_address` - The first private (`PRIVATE`) IPv4 address assigned. This is
a workaround for an [issue fixed in Terraform 0.12](https://github.com/hashicorp/terraform/issues/17048)
but also provides a convenient way to access an IP of a specific type without
performing filtering in a Terraform config.

* `psc_service_attachment_link` - the URI that points to the service attachment of the instance.

* `instance_type` - The type of the instance. The supported values are `SQL_INSTANCE_TYPE_UNSPECIFIED`, `CLOUD_SQL_INSTANCE`, `ON_PREMISES_INSTANCE` and `READ_REPLICA_INSTANCE`.

~> **NOTE:** Users can upgrade a read replica instance to a stand-alone Cloud SQL instance with the help of `instance_type`. To promote, users have to set the `instance_type` property as `CLOUD_SQL_INSTANCE` and remove/unset `master_instance_name` and `replica_configuration` from instance configuration. This operation might cause your instance to restart.

* `settings.version` - Used to make sure changes to the `settings` block are
    atomic.

* `server_ca_cert.0.cert` - The CA Certificate used to connect to the SQL Instance via SSL.

* `server_ca_cert.0.common_name` - The CN valid for the CA Cert.

* `server_ca_cert.0.create_time` - Creation time of the CA Cert.

* `server_ca_cert.0.expiration_time` - Expiration time of the CA Cert.

* `server_ca_cert.0.sha1_fingerprint` - SHA Fingerprint of the CA Cert.

## Timeouts

`google_sql_database_instance` provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `create` - Default is 90 minutes.
- `update` - Default is 90 minutes.
- `delete` - Default is 90 minutes.

## Import

Database instances can be imported using one of any of these accepted formats:

* `projects/{{project}}/instances/{{name}}`
* `{{project}}/{{name}}`
* `{{name}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Database instances using one of the formats above. For example:

```tf
import {
  id = "projects/{{project}}/instances/{{name}}"
  to = google_sql_database_instance.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), Database instances can be imported using one of the formats above. For example:

```
$ terraform import google_sql_database_instance.default projects/{{project}}/instances/{{name}}
$ terraform import google_sql_database_instance.default {{project}}/{{name}}
$ terraform import google_sql_database_instance.default {{name}}
```

~> **NOTE:** Some fields (such as `replica_configuration`) won't show a diff if they are unset in
config and set on the server.
When importing, double-check that your config has all the fields set that you expect- just seeing
no diff isn't sufficient to know that your config could reproduce the imported resource.
