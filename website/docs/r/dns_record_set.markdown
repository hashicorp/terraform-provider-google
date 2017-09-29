---
layout: "google"
page_title: "Google: google_dns_record_set"
sidebar_current: "docs-google-dns-record-set"
description: |-
  Manages a set of DNS records within Google Cloud DNS.
---

# google\_dns\_record\_set

Manages a set of DNS records within Google Cloud DNS.

## Example Usage

### Binding a DNS name to the ephemeral IP of a new instance:

```hcl
resource "google_dns_record_set" "frontend" {
  name = "frontend.${google_dns_managed_zone.prod.dns_name}"
  type = "A"
  ttl  = 300

  managed_zone = "${google_dns_managed_zone.prod.name}"

  rrdatas = ["${google_compute_instance.frontend.network_interface.0.access_config.0.assigned_nat_ip}"]
}

resource "google_compute_instance" "frontend" {
  name         = "frontend"
  machine_type = "g1-small"
  zone         = "us-central1-b"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-8"
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

### Adding a SPF record

`\"` must be added around your `rrdatas` for a SPF record. Otherwise `rrdatas` string gets split on spaces.

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

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

Only the arguments listed above are exposed as attributes.
