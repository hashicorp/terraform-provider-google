---
layout: "google"
page_title: "Google: google_dns_record_set"
sidebar_current: "docs-google-dns-record-set"
description: |-
  Manages a set of DNS records within Google Cloud DNS.
---

# google\_dns\_record\_set

Manages a set of DNS records within Google Cloud DNS. For more information see [the official documentation](https://cloud.google.com/dns/records/) and
[API](https://cloud.google.com/dns/api/v1/resourceRecordSets).

~> **Note:** The Google Cloud DNS API requires NS records be present at all
times. To accommodate this, when creating NS records, the default records
Google automatically creates will be silently overwritten.  Also, when
destroying NS records, Terraform will not actually remove NS records, but will
report that it did.

## Example Usage

### Binding a DNS name to the ephemeral IP of a new instance:

```hcl
resource "google_dns_record_set" "frontend" {
  name = "frontend.${google_dns_managed_zone.prod.dns_name}"
  type = "A"
  ttl  = 300

  managed_zone = "${google_dns_managed_zone.prod.name}"

  rrdatas = ["${google_compute_instance.frontend.network_interface.0.access_config.0.nat_ip}"]
}

resource "google_compute_instance" "frontend" {
  name         = "frontend"
  machine_type = "g1-small"
  zone         = "us-central1-b"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }

  network_interface {
    network       = "default"
    access_config = {}
  }
}

resource "google_dns_managed_zone" "prod" {
  name     = "prod-zone"
  dns_name = "prod.mydomain.com."
}
```

### Adding an A record

```hcl
resource "google_dns_record_set" "a" {
  name = "backend.${google_dns_managed_zone.prod.dns_name}"
  managed_zone = "${google_dns_managed_zone.prod.name}"
  type = "A"
  ttl  = 300

  rrdatas = ["8.8.8.8"]
}

resource "google_dns_managed_zone" "prod" {
  name     = "prod-zone"
  dns_name = "prod.mydomain.com."
}
```

### Adding an MX record

```hcl
resource "google_dns_record_set" "mx" {
  name = "${google_dns_managed_zone.prod.dns_name}"
  managed_zone = "${google_dns_managed_zone.prod.name}"
  type = "MX"
  ttl  = 3600

  rrdatas = [
    "1 aspmx.l.google.com.",
    "5 alt1.aspmx.l.google.com.",
    "5 alt2.aspmx.l.google.com.",
    "10 alt3.aspmx.l.google.com.",
    "10 alt4.aspmx.l.google.com."
  ]
}

resource "google_dns_managed_zone" "prod" {
  name     = "prod-zone"
  dns_name = "prod.mydomain.com."
}
```

### Adding an SPF record

Quotes (`""`) must be added around your `rrdatas` for a SPF record. Otherwise `rrdatas` string gets split on spaces.

```hcl
resource "google_dns_record_set" "spf" {
  name = "frontend.${google_dns_managed_zone.prod.dns_name}"
  managed_zone = "${google_dns_managed_zone.prod.name}"
  type = "TXT"
  ttl  = 300

  rrdatas = ["\"v=spf1 ip4:111.111.111.111 include:backoff.email-example.com -all\""]
}

resource "google_dns_managed_zone" "prod" {
  name     = "prod-zone"
  dns_name = "prod.mydomain.com."
}
```

### Adding a CNAME record
 The list of `rrdatas` should only contain a single string corresponding to the Canonical Name intended.
 ```hcl
resource "google_dns_record_set" "cname" {
  name = "frontend.${google_dns_managed_zone.prod.dns_name}"
  managed_zone = "${google_dns_managed_zone.prod.name}"
  type = "CNAME"
  ttl  = 300
  rrdatas = ["frontend.mydomain.com."]
}

resource "google_dns_managed_zone" "prod" {
  name        = "prod-zone"
  dns_name    = "prod.mydomain.com."
}
```

## Argument Reference

The following arguments are supported:

* `managed_zone` - (Required) The name of the zone in which this record set will
    reside.

* `name` - (Required) The DNS name this record set will apply to.

* `rrdatas` - (Required) The string data for the records in this record set
    whose meaning depends on the DNS type. For TXT record, if the string data contains spaces, add surrounding `\"` if you don't want your string to get split on spaces.

* `ttl` - (Required) The time-to-live of this record set (seconds).

* `type` - (Required) The DNS record set type.

- - -

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

Only the arguments listed above are exposed as attributes.

## Import

DNS record set can be imported using the `zone name`, `record name` and record `type`, e.g.

```
$ terraform import google_dns_record_set.frontend prod-zone/frontend.prod.mydomain.com./A
```

Note: The record name must include the trailing dot at the end.
