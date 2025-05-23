---
# ----------------------------------------------------------------------------
#
#     ***     AUTO GENERATED CODE    ***    Type: Handwritten     ***
#
# ----------------------------------------------------------------------------
#
#     This code is generated by Magic Modules using the following:
#
#     Source file: https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/website/docs/r/service_networking_connection.html.markdown
#
#     DO NOT EDIT this file directly. Any changes made to this file will be
#     overwritten during the next generation cycle.
#
# ----------------------------------------------------------------------------
subcategory: "Service Networking"
description: |-
  Manages creating a private VPC connection to a service provider.
---

# google_service_networking_connection

Manages a private VPC connection with a GCP service provider. For more information see
[the official documentation](https://cloud.google.com/vpc/docs/configure-private-services-access#creating-connection)
and
[API](https://cloud.google.com/service-infrastructure/docs/service-networking/reference/rest/v1/services.connections).

## Example usage

```hcl
# Create a VPC network
resource "google_compute_network" "peering_network" {
  name = "peering-network"
}

# Create an IP address
resource "google_compute_global_address" "private_ip_alloc" {
  name          = "private-ip-alloc"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.peering_network.id
}

# Create a private connection
resource "google_service_networking_connection" "default" {
  network                 = google_compute_network.peering_network.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_alloc.name]
}

# (Optional) Import or export custom routes
resource "google_compute_network_peering_routes_config" "peering_routes" {
  peering = google_service_networking_connection.default.peering
  network = google_compute_network.peering_network.name

  import_custom_routes = true
  export_custom_routes = true
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

* `deletion_policy` - (Optional) The deletion policy for the service networking connection. Setting to ABANDON allows the resource to be abandoned rather than deleted. This will enable a successful terraform destroy when destroying CloudSQL instances. Use with care as it can lead to dangling resources.

* `update_on_creation_fail` - (Optional) When set to true, enforce an update of the reserved peering ranges on the existing service networking connection in case of a new connection creation failure.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `peering` - (Computed) The name of the VPC Network Peering connection that was created by the service producer.


## Import

ServiceNetworkingConnection can be imported using any of these accepted formats

* `{{peering-network}}:{{service}}`
* `projects/{{project}}/global/networks/{{peering-network}}:{{service}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import NAME_HERE using one of the formats above. For example:

```tf
import {
  id = "projects/{{project}}/global/networks/{{peering-network}}:{{service}}"
  to = google_service_networking_connection.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), NAME_HERE can be imported using one of the formats above. For example:

```
$ terraform import google_service_networking_connection.default {{peering-network}}:{{service}}
$ terraform import google_service_networking_connection.default /projects/{{project}}/global/networks/{{peering-network}}:{{service}}
```


## User Project Overrides

This resource supports [User Project Overrides](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/provider_reference#user_project_override).
