---
subcategory: "Cloud DNS"
page_title: "Google: google_dns_record_set"
description: |-
  Manages a set of DNS records within Google Cloud DNS.
---

# google\_dns\_record\_set

Manages a set of DNS records within Google Cloud DNS. For more information see [the official documentation](https://cloud.google.com/dns/records/) and
[API](https://cloud.google.com/dns/api/v1/resourceRecordSets).

~> **Note:** The provider treats this resource as an authoritative record set. This means existing records (including the default records) for the given type will be overwritten when you create this resource in Terraform. In addition, the Google Cloud DNS API requires NS records to be present at all times, so Terraform will not actually remove NS records during destroy but will report that it did.

## Example Usage

### Binding a DNS name to the ephemeral IP of a new instance:

```hcl
resource "google_dns_record_set" "frontend" {
  name = "frontend.${google_dns_managed_zone.prod.dns_name}"
  type = "A"
  ttl  = 300

  managed_zone = google_dns_managed_zone.prod.name

  rrdatas = [google_compute_instance.frontend.network_interface[0].access_config[0].nat_ip]
}

resource "google_compute_instance" "frontend" {
  name         = "frontend"
  machine_type = "g1-small"
  zone         = "us-central1-b"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }

  network_interface {
    network = "default"
    access_config {
    }
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
  name         = "backend.${google_dns_managed_zone.prod.dns_name}"
  managed_zone = google_dns_managed_zone.prod.name
  type         = "A"
  ttl          = 300

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
  name         = google_dns_managed_zone.prod.dns_name
  managed_zone = google_dns_managed_zone.prod.name
  type         = "MX"
  ttl          = 3600

  rrdatas = [
    "1 aspmx.l.google.com.",
    "5 alt1.aspmx.l.google.com.",
    "5 alt2.aspmx.l.google.com.",
    "10 alt3.aspmx.l.google.com.",
    "10 alt4.aspmx.l.google.com.",
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
  name         = "frontend.${google_dns_managed_zone.prod.dns_name}"
  managed_zone = google_dns_managed_zone.prod.name
  type         = "TXT"
  ttl          = 300

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
  name         = "frontend.${google_dns_managed_zone.prod.dns_name}"
  managed_zone = google_dns_managed_zone.prod.name
  type         = "CNAME"
  ttl          = 300
  rrdatas      = ["frontend.mydomain.com."]
}

resource "google_dns_managed_zone" "prod" {
  name     = "prod-zone"
  dns_name = "prod.mydomain.com."
}
```

### Setting Routing Policy instead of using rrdatas
#### Weighted Round Robin

```hcl
resource "google_dns_record_set" "wrr" {
  name         = "backend.${google_dns_managed_zone.prod.dns_name}"
  managed_zone = google_dns_managed_zone.prod.name
  type         = "A"
  ttl          = 300

  routing_policy {
    wrr {
      weight  = 0.8
      rrdatas =  ["10.128.1.1"]
    }

    wrr {
      weight  = 0.2
      rrdatas =  ["10.130.1.1"]
    }
  }
```

#### Geolocation

```hcl
resource "google_dns_record_set" "geo" {
  name         = "backend.${google_dns_managed_zone.prod.dns_name}"
  managed_zone = google_dns_managed_zone.prod.name
  type         = "A"
  ttl          = 300

  routing_policy {
    geo {
      location = "asia-east1"
      rrdatas  =  ["10.128.1.1"]
    }

    geo {
      location = "us-central1"
      rrdatas  =  ["10.130.1.1"]
    }
  }
}
```

#### Primary-Backup

```hcl
resource "google_dns_record_set" "a" {
  name         = "backend.${google_dns_managed_zone.prod.dns_name}"
  managed_zone = google_dns_managed_zone.prod.name
  type         = "A"
  ttl          = 300

  routing_policy {
    primary_backup {
      trickle_ratio = 0.1

      primary {
        internal_load_balancers {
          load_balancer_type = "regionalL4ilb"
          ip_address         = google_compute_forwarding_rule.prod.ip_address
          port               = "80"
          ip_protocol        = "tcp"
          network_url        = google_compute_network.prod.id
          project            = google_compute_forwarding_rule.prod.project
          region             = google_compute_forwarding_rule.prod.region
        }
      }

      backup_geo {
        location = "asia-east1"
        rrdatas  = ["10.128.1.1"]
      }

      backup_geo {
        location = "us-west1"
        rrdatas  = ["10.130.1.1"]
      }
    }
  }
}

resource "google_dns_managed_zone" "prod" {
  name     = "prod-zone"
  dns_name = "prod.mydomain.com."
}

resource "google_compute_forwarding_rule" "prod" {
  name   = "prod-ilb"
  region = "us-central1"

  load_balancing_scheme = "INTERNAL"
  backend_service       = google_compute_region_backend_service.prod.id
  all_ports             = true
  network               = google_compute_network.prod.name
}

resource "google_compute_region_backend_service" "prod" {
  name   = "prod-backend"
  region = "us-central1"
}

resource "google_compute_network" "prod" {
  name = "prod-network"
}
```

## Argument Reference

The following arguments are supported:

* `managed_zone` - (Required) The name of the zone in which this record set will
    reside.

* `name` - (Required) The DNS name this record set will apply to.

* `type` - (Required) The DNS record set type.

- - -

* `rrdatas` - (Optional) The string data for the records in this record set
    whose meaning depends on the DNS type. For TXT record, if the string data contains spaces, add surrounding `\"` if you don't want your string to get split on spaces. To specify a single record value longer than 255 characters such as a TXT record for DKIM, add `\" \"` inside the Terraform configuration string (e.g. `"first255characters\" \"morecharacters"`).

* `routing_policy` - (Optional) The configuration for steering traffic based on query.
    Now you can specify either Weighted Round Robin(WRR) type or Geolocation(GEO) type.
    Structure is [documented below](#nested_routing_policy).

* `ttl` - (Optional) The time-to-live of this record set (seconds).

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

<a name="nested_routing_policy"></a>The `routing_policy` block supports:

* `wrr` - (Optional) The configuration for Weighted Round Robin based routing policy.
    Structure is [document below](#nested_wrr).

* `geo` - (Optional) The configuration for Geolocation based routing policy.
    Structure is [document below](#nested_geo).

* `enable_geo_fencing` - (Optional) Specifies whether to enable fencing for geo queries.

* `primary_backup` - (Optional) The configuration for a primary-backup policy with global to regional failover. Queries are responded to with the global primary targets, but if none of the primary targets are healthy, then we fallback to a regional failover policy.
    Structure is [document below](#nested_primary_backup).

<a name="nested_wrr"></a>The `wrr` block supports:

* `weight`  - (Required) The ratio of traffic routed to the target.

* `rrdatas` - (Optional) Same as `rrdatas` above.

* `health_checked_targets` - (Optional) The list of targets to be health checked. Note that if DNSSEC is enabled for this zone, only one of `rrdatas` or `health_checked_targets` can be set.
    Structure is [document below](#nested_health_checked_targets).

<a name="nested_geo"></a>The `geo` block supports:

* `location` - (Required) The location name defined in Google Cloud.

* `rrdatas` - (Optional) Same as `rrdatas` above.

* `health_checked_targets` - (Optional) For A and AAAA types only. The list of targets to be health checked. These can be specified along with `rrdatas` within this item.
    Structure is [document below](#nested_health_checked_targets).

<a name="nested_primary_backup"></a>The `primary_backup` block supports:

* `primary` - (Required) The list of global primary targets to be health checked.
    Structure is [document below](#nested_health_checked_targets).

* `backup_geo` - (Required) The backup geo targets, which provide a regional failover policy for the otherwise global primary targets.
    Structure is [document above](#nested_geo).

* `enable_geo_fencing_for_backups` - (Optional) Specifies whether to enable fencing for backup geo queries.

* `trickle_ratio` - (Optional) Specifies the percentage of traffic to send to the backup targets even when the primary targets are healthy.

<a name="nested_health_checked_targets"></a>The `health_checked_targets` block supports:

* `internal_load_balancers` - (Required) The list of internal load balancers to health check.
    Structure is [document below](#nested_internal_load_balancers).

<a name="nested_internal_load_balancers"></a>The `internal_load_balancers` block supports:

* `load_balancer_type` - (Required) The type of load balancer. This value is case-sensitive. Possible values: ["regionalL4ilb"]

* `ip_address` - (Required) The frontend IP address of the load balancer.

* `port` - (Required) The configured port of the load balancer.

* `ip_protocol` - (Required) The configured IP protocol of the load balancer. This value is case-sensitive. Possible values: ["tcp", "udp"]

* `network_url` - (Required) The fully qualified url of the network in which the load balancer belongs. This should be formatted like `projects/{project}/global/networks/{network}` or `https://www.googleapis.com/compute/v1/projects/{project}/global/networks/{network}`.

* `project` - (Required) The ID of the project in which the load balancer belongs.

* `region` - (Optional) The region of the load balancer. Only needed for regional load balancers.

## Attributes Reference

-In addition to the arguments listed above, the following computed attributes are
-exported:

* `id` - an identifier for the resource with format `projects/{{project}}/managedZones/{{zone}}/rrsets/{{name}}/{{type}}`

## Import

DNS record sets can be imported using either of these accepted formats:

```
$ terraform import google_dns_record_set.frontend projects/{{project}}/managedZones/{{zone}}/rrsets/{{name}}/{{type}}
$ terraform import google_dns_record_set.frontend {{project}}/{{zone}}/{{name}}/{{type}}
$ terraform import google_dns_record_set.frontend {{zone}}/{{name}}/{{type}}
```

Note: The record name must include the trailing dot at the end.
