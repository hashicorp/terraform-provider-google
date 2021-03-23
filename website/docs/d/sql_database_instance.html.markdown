---
subcategory: "Cloud SQL"
layout: "google"
page_title: "Google: google_sql_database_instance"
sidebar_current: "docs-google-datasource-sql-database-instance"
description: |-
  Get a  SQL database instance in Google Cloud SQL.
---

# google\_sql\_database\_instance

Use this data source to get information about a Cloud SQL instance

## Example Usage 


```hcl
data "google_sql_database_instance" "qa" {
    name = "test-sql-instance"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (required) The name of the instance.

* `project` - (optional) The ID of the project in which the resource belongs.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:
    
* `settings` -  The settings to use for the database. The
    configuration is detailed below.

* `database_version` - The MySQL, PostgreSQL or SQL Server (beta) version to use.

* `master_instance_name` - The name of the instance that will act as
    the master in the replication setup.

* `replica_configuration` - The configuration for replication. The
    configuration is detailed below.
    
* `root_password` - Initial root password. Required for MS SQL Server, ignored by MySQL and PostgreSQL.

* `encryption_key_name` - [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html))
    The full path to the encryption key used for the CMEK disk encryption.
    
The `settings` block contains:

* `tier` - The machine type to use.
    
* `activation_policy` - This specifies when the instance should be
    active. Can be either `ALWAYS`, `NEVER` or `ON_DEMAND`.

* `authorized_gae_applications` - (Deprecated) This property is only applicable to First Generation instances.
    First Generation instances are now deprecated, see [here](https://cloud.google.com/sql/docs/mysql/upgrade-2nd-gen)
    for information on how to upgrade to Second Generation instances.
    A list of Google App Engine (GAE) project names that are allowed to access this instance.

* `availability_type` - The availability type of the Cloud SQL
instance, high availability (`REGIONAL`) or single zone (`ZONAL`).

* `crash_safe_replication` - (Deprecated) This property is only applicable to First Generation instances.
    First Generation instances are now deprecated, see [here](https://cloud.google.com/sql/docs/mysql/upgrade-2nd-gen)

* `disk_autoresize` - Configuration to increase storage size automatically.

* `disk_size` - The size of data disk, in GB.

* `disk_type` - The type of data disk.

* `pricing_plan` - Pricing plan for this instance.

* `replication_type` - This property is only applicable to First Generation instances.
    First Generation instances are now deprecated, see [here](https://cloud.google.com/sql/docs/mysql/upgrade-2nd-gen)

* `user_labels` - A set of key/value user label pairs to assign to the instance.

The `settings.database_flags` sublist contains:

* `name` - Name of the flag.

* `value` - Value of the flag.

The `settings.backup_configuration` subblock contains:

* `binary_log_enabled` - True if binary logging is enabled.

* `enabled` - True if backup configuration is enabled.

* `start_time` - `HH:MM` format time indicating when backup configuration starts.

The `settings.ip_configuration` subblock contains:

* `ipv4_enabled` - Whether this Cloud SQL instance should be assigned a public IPV4 address. 

* `private_network` - The VPC network from which the Cloud SQL instance is accessible for private IP.

* `require_ssl` - True if mysqld default to `REQUIRE X509` for users connecting over IP.

The `settings.ip_configuration.authorized_networks[]` sublist contains:

* `expiration_time` - The [RFC 3339](https://tools.ietf.org/html/rfc3339)
  formatted date time string indicating when this whitelist expires.

* `name` - A name for this whitelist entry.

* `value` - A CIDR notation IPv4 or IPv6 address that is allowed to access this instance.

The `settings.location_preference` subblock contains:

* `follow_gae_application` - A GAE application whose zone to remain in.

* `zone` - The preferred compute engine.

The `settings.maintenance_window` subblock for instances declares a one-hour
[maintenance window](https://cloud.google.com/sql/docs/instance-settings?hl=en#maintenance-window-2ndgen)
when an Instance can automatically restart to apply updates. The maintenance window is specified in UTC time. It contains:

* `day` - Day of week (`1-7`), starting on Monday.

* `hour` - Hour of day (`0-23`), ignored if `day` not set.

* `update_track` - Receive updates earlier (`canary`) or later (`stable`).

The `settings.insights_config` subblock for instances declares [Query Insights](https://cloud.google.com/sql/docs/postgres/insights-overview) configuration. It contains:

* `query_insights_enabled` - True if Query Insights feature is enabled.

* `query_string_length` - Maximum query length stored in bytes. Between 256 and 4500. Default to 1024.

* `record_application_tags` - True if Query Insights will record application tags from query when enabled.

* `record_client_address` - True if Query Insights will record client address when enabled.

The `replica_configuration` block contains:

* `ca_certificate` - PEM representation of the trusted CA's x509 certificate.

* `client_certificate` - PEM representation of the replica's x509 certificate.

* `client_key` - PEM representation of the replica's private key.

* `connect_retry_interval` - The number of seconds between connect retries.

* `dump_file_path` - Path to a SQL file in GCS from which replica instances are created. 

* `failover_target` - Specifies if the replica is the failover target.

* `master_heartbeat_period` - Time in ms between replication heartbeats.

* `password` - Password for the replication connection.

* `sslCipher` - Permissible ciphers for use in SSL encryption.

* `username` - Username for replication connection.

* `verify_server_certificate` - True if the master's common name value is checked during the SSL handshake.

* `self_link` - The URI of the created resource.

* `connection_name` - The connection name of the instance to be used in connection strings.

* `service_account_email_address` - The service account email address assigned to the instance.

* `ip_address.0.ip_address` - The IPv4 address assigned.

* `ip_address.0.time_to_retire` - The time this IP address will be retired, in RFC 3339 format.

* `ip_address.0.type` - The type of this IP address.

* `first_ip_address` - The first IPv4 address of any type assigned.

* `public_ip_address` - The first public (`PRIMARY`) IPv4 address assigned.

* `private_ip_address` - The first private (`PRIVATE`) IPv4 address assigned.

* `settings.version` - Used to make sure changes to the `settings` block are atomic.

* `server_ca_cert.0.cert` - The CA Certificate used to connect to the SQL Instance via SSL.

* `server_ca_cert.0.common_name` - The CN valid for the CA Cert.

* `server_ca_cert.0.create_time` - Creation time of the CA Cert.

* `server_ca_cert.0.expiration_time` - Expiration time of the CA Cert.

* `server_ca_cert.0.sha1_fingerprint` - SHA Fingerprint of the CA Cert.
