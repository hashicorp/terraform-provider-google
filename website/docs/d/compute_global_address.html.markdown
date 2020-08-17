---
subcategory: "Compute Engine"
layout: "google"
page_title: "Google: google_compute_global_address"
sidebar_current: "docs-google-datasource-compute-global-address"
description: |-
  Get the IP address from a static address reserved for a Global Forwarding Rule.
---

# google\_compute\_global\_address

Get the IP address from a static address reserved for a Global Forwarding Rule which are only used for HTTP load balancing. For more information see
the official [API](https://cloud.google.com/compute/docs/reference/latest/globalAddresses) documentation.

## Example Usage

```hcl
data "google_compute_global_address" "my_address" {
  name = "foobar"
}

resource "google_dns_record_set" "frontend" {
  name = "lb.${google_dns_managed_zone.prod.dns_name}"
  type = "A"
  ttl  = 300

  managed_zone = google_dns_managed_zone.prod.name

  rrdatas = [data.google_compute_global_address.my_address.address]
}

resource "google_dns_managed_zone" "prod" {
  name     = "prod-zone"
  dns_name = "prod.mydomain.com."
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for the resource, required by GCE.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `self_link` - The URI of the created resource.
* `address` - The IP of the created resource.
* `status` - Indicates if the address is used. Possible values are: RESERVED or IN_USE.
