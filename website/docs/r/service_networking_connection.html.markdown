---
subcategory: "Service Networking"
layout: "google"
page_title: "Google: google_service_networking_connection"
sidebar_current: "docs-google-service-networking-connection"
description: |-
  Manages creating a private VPC connection to a service provider.
---

# google\_service\_networking\_connection

Manages a private VPC connection with a GCP service provider. For more information see
[the official documentation](https://cloud.google.com/vpc/docs/configure-private-services-access#creating-connection)
and
[API](https://cloud.google.com/service-infrastructure/docs/service-networking/reference/rest/v1/services.connections).

## Example usage

```hcl
resource "google_compute_network" "peering_network" {
  name = "peering-network"
}

resource "google_compute_global_address" "private_ip_alloc" {
  name          = "private-ip-alloc"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.peering_network.id
}

resource "google_service_networking_connection" "foobar" {
  network                 = google_compute_network.peering_network.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_alloc.name]
}
```

## Argument Reference

The following arguments are supported:

* `network` - (Required) Name of VPC network connected with service producers using VPC peering.

* `service` - (Required) Provider peering service that is managing peering connectivity for a
  service provider organization. For Google services that support this functionality it is
  'servicenetworking.googleapis.com'.

* `reserved_peering_ranges` - (Required) Named IP address range(s) of PEERING type reserved for
  this service provider. Note that invoking this method with a different range when connection
  is already established will not reallocate already provisioned service producer subnetworks.
