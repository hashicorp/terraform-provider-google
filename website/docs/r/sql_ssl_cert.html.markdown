---
subcategory: "Cloud SQL"
layout: "google"
page_title: "Google: google_sql_ssl_cert"
sidebar_current: "docs-google-sql-ssl-cert"
description: |-
  Creates a new SQL Ssl Cert in Google Cloud SQL.
---

# google\_sql\_ssl\_cert

Creates a new Google SQL SSL Cert on a Google SQL Instance. For more information, see the [official documentation](https://cloud.google.com/sql/), or the [JSON API](https://cloud.google.com/sql/docs/mysql/admin-api/v1beta4/sslCerts).

~> **Note:** All arguments including the private key will be stored in the raw state as plain-text.
[Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

Example creating a SQL Client Certificate.

```hcl
resource "random_id" "db_name_suffix" {
  byte_length = 4
}

resource "google_sql_database_instance" "master" {
  name = "master-instance-${random_id.db_name_suffix.hex}"

  settings {
    tier = "db-f1-micro"
  }
}

resource "google_sql_ssl_cert" "client_cert" {
  common_name = "client-name"
  instance    = google_sql_database_instance.master.name
}
```

## Argument Reference

The following arguments are supported:

* `instance` - (Required) The name of the Cloud SQL instance. Changing this
    forces a new resource to be created.

* `common_name` - (Required) The common name to be used in the certificate to identify the
    client. Constrained to [a-zA-Z.-_ ]+. Changing this forces a new resource to be created.

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `sha1_fingerprint` - The SHA1 Fingerprint of the certificate.
* `private_key` - The private key associated with the client certificate.
* `server_ca_cert` - The CA cert of the server this client cert was generated from.
* `cert` - The actual certificate data for this client certificate.
* `cert_serial_number` - The serial number extracted from the certificate data.
* `create_time` - The time when the certificate was created in RFC 3339 format,
    for example 2012-11-15T16:19:00.094Z.
* `expiration_time` - The time when the certificate expires in RFC 3339 format,
    for example 2012-11-15T16:19:00.094Z.

## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 10 minutes.
- `delete` - Default is 10 minutes.

## Import

Since the contents of the certificate cannot be accessed after its creation, this resource cannot be imported.
