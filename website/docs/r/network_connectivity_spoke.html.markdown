---
# ----------------------------------------------------------------------------
#
#     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
#
# ----------------------------------------------------------------------------
#
#     This code is generated by Magic Modules using the following:
#
#     Configuration: https:#github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products/networkconnectivity/Spoke.yaml
#     Template:      https:#github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/resource.html.markdown.tmpl
#
#     DO NOT EDIT this file directly. Any changes made to this file will be
#     overwritten during the next generation cycle.
#
# ----------------------------------------------------------------------------
subcategory: "Network Connectivity"
description: |-
  The NetworkConnectivity Spoke resource
---

# google_network_connectivity_spoke

The NetworkConnectivity Spoke resource


To get more information about Spoke, see:

* [API documentation](https://cloud.google.com/network-connectivity/docs/reference/networkconnectivity/rest/v1beta/projects.locations.spokes)
* How-to Guides
    * [Official Documentation](https://cloud.google.com/network-connectivity/docs/network-connectivity-center/concepts/overview)

<div class = "oics-button" style="float: right; margin: 0 0 -15px">
  <a href="https://console.cloud.google.com/cloudshell/open?cloudshell_git_repo=https%3A%2F%2Fgithub.com%2Fterraform-google-modules%2Fdocs-examples.git&cloudshell_image=gcr.io%2Fcloudshell-images%2Fcloudshell%3Alatest&cloudshell_print=.%2Fmotd&cloudshell_tutorial=.%2Ftutorial.md&cloudshell_working_dir=network_connectivity_spoke_linked_vpc_network_basic&open_in_editor=main.tf" target="_blank">
    <img alt="Open in Cloud Shell" src="//gstatic.com/cloudssh/images/open-btn.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>
## Example Usage - Network Connectivity Spoke Linked Vpc Network Basic


```hcl
resource "google_compute_network" "network" {
  name                    = "net"
  auto_create_subnetworks = false
}

resource "google_network_connectivity_hub" "basic_hub" {
  name        = "hub1"
  description = "A sample hub"
  labels = {
    label-two = "value-one"
  }
}

resource "google_network_connectivity_spoke" "primary"  {
  name = "spoke1"
  location = "global"
  description = "A sample spoke with a linked router appliance instance"
  labels = {
    label-one = "value-one"
  }
  hub = google_network_connectivity_hub.basic_hub.id
  linked_vpc_network {
    exclude_export_ranges = [
      "198.51.100.0/24",
      "10.10.0.0/16"
    ]
    include_export_ranges = [
      "198.51.100.0/23", 
      "10.0.0.0/8"
    ]
    uri = google_compute_network.network.self_link
  }
}
```
<div class = "oics-button" style="float: right; margin: 0 0 -15px">
  <a href="https://console.cloud.google.com/cloudshell/open?cloudshell_git_repo=https%3A%2F%2Fgithub.com%2Fterraform-google-modules%2Fdocs-examples.git&cloudshell_image=gcr.io%2Fcloudshell-images%2Fcloudshell%3Alatest&cloudshell_print=.%2Fmotd&cloudshell_tutorial=.%2Ftutorial.md&cloudshell_working_dir=network_connectivity_spoke_linked_vpc_network_group&open_in_editor=main.tf" target="_blank">
    <img alt="Open in Cloud Shell" src="//gstatic.com/cloudssh/images/open-btn.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>
## Example Usage - Network Connectivity Spoke Linked Vpc Network Group


```hcl
resource "google_compute_network" "network" {
  name                    = "net-spoke"
  auto_create_subnetworks = false
}

resource "google_network_connectivity_hub" "basic_hub" {
  name        = "hub1-spoke"
  description = "A sample hub"
  labels = {
    label-two = "value-one"
  }
}

resource "google_network_connectivity_group" "default_group"  {
 hub         = google_network_connectivity_hub.basic_hub.id
 name        = "default"
 description = "A sample hub group"
}

resource "google_network_connectivity_spoke" "primary"  {
  name = "group-spoke1"
  location = "global"
  description = "A sample spoke with a linked VPC"
  labels = {
    label-one = "value-one"
  }
  hub = google_network_connectivity_hub.basic_hub.id
  linked_vpc_network {
    exclude_export_ranges = [
      "198.51.100.0/24",
      "10.10.0.0/16"
    ]
    include_export_ranges = [
      "198.51.100.0/23",
      "10.0.0.0/8"
    ]
    uri = google_compute_network.network.self_link
  }
  group = google_network_connectivity_group.default_group.id
}
```
<div class = "oics-button" style="float: right; margin: 0 0 -15px">
  <a href="https://console.cloud.google.com/cloudshell/open?cloudshell_git_repo=https%3A%2F%2Fgithub.com%2Fterraform-google-modules%2Fdocs-examples.git&cloudshell_image=gcr.io%2Fcloudshell-images%2Fcloudshell%3Alatest&cloudshell_print=.%2Fmotd&cloudshell_tutorial=.%2Ftutorial.md&cloudshell_working_dir=network_connectivity_spoke_router_appliance_basic&open_in_editor=main.tf" target="_blank">
    <img alt="Open in Cloud Shell" src="//gstatic.com/cloudssh/images/open-btn.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>
## Example Usage - Network Connectivity Spoke Router Appliance Basic


```hcl
resource "google_compute_network" "network" {
  name                    = "tf-test-network%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnetwork" {
  name          = "tf-test-subnet%{random_suffix}"
  ip_cidr_range = "10.0.0.0/28"
  region        = "us-central1"
  network       = google_compute_network.network.self_link
}

resource "google_compute_instance" "instance" {
  name         = "tf-test-instance%{random_suffix}"
  machine_type = "e2-medium"
  can_ip_forward = true
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "projects/debian-cloud/global/images/debian-10-buster-v20210817"
    }
  }

  network_interface {
    subnetwork = google_compute_subnetwork.subnetwork.name
    network_ip = "10.0.0.2"
    access_config {
        network_tier = "PREMIUM"
    }
  }
}

resource "google_network_connectivity_hub" "basic_hub" {
  name        = "tf-test-hub%{random_suffix}"
  description = "A sample hub"
  labels = {
    label-two = "value-one"
  }
}

resource "google_network_connectivity_spoke" "primary" {
  name = "tf-test-name%{random_suffix}"
  location = "us-central1"
  description = "A sample spoke with a linked routher appliance instance"
  labels = {
    label-one = "value-one"
  }
  hub =  google_network_connectivity_hub.basic_hub.id
  linked_router_appliance_instances {
    instances {
        virtual_machine = google_compute_instance.instance.self_link
        ip_address = "10.0.0.2"
    }
    site_to_site_data_transfer = true
    include_import_ranges = ["ALL_IPV4_RANGES"]
  }
}
```
<div class = "oics-button" style="float: right; margin: 0 0 -15px">
  <a href="https://console.cloud.google.com/cloudshell/open?cloudshell_git_repo=https%3A%2F%2Fgithub.com%2Fterraform-google-modules%2Fdocs-examples.git&cloudshell_image=gcr.io%2Fcloudshell-images%2Fcloudshell%3Alatest&cloudshell_print=.%2Fmotd&cloudshell_tutorial=.%2Ftutorial.md&cloudshell_working_dir=network_connectivity_spoke_vpn_tunnel_basic&open_in_editor=main.tf" target="_blank">
    <img alt="Open in Cloud Shell" src="//gstatic.com/cloudssh/images/open-btn.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>
## Example Usage - Network Connectivity Spoke Vpn Tunnel Basic


```hcl
resource "google_network_connectivity_hub" "basic_hub" {
  name        = "basic-hub1"
  description = "A sample hub"
  labels = {
    label-two = "value-one"
  }
}

resource "google_compute_network" "network" {
  name                    = "basic-network"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnetwork" {
  name          = "basic-subnetwork"
  ip_cidr_range = "10.0.0.0/28"
  region        = "us-central1"
  network       = google_compute_network.network.self_link
}

resource "google_compute_ha_vpn_gateway" "gateway" {
  name    = "vpn-gateway"
  network = google_compute_network.network.id
}

resource "google_compute_external_vpn_gateway" "external_vpn_gw" {
  name            = "external-vpn-gateway"
  redundancy_type = "SINGLE_IP_INTERNALLY_REDUNDANT"
  description     = "An externally managed VPN gateway"
  interface {
    id         = 0
    ip_address = "8.8.8.8"
  }
}

resource "google_compute_router" "router" {
  name    = "external-vpn-gateway"
  region  = "us-central1"
  network = google_compute_network.network.name
  bgp {
    asn = 64514
  }
}

resource "google_compute_vpn_tunnel" "tunnel1" {
  name                            = "tunnel1"
  region                          = "us-central1"
  vpn_gateway                     = google_compute_ha_vpn_gateway.gateway.id
  peer_external_gateway           = google_compute_external_vpn_gateway.external_vpn_gw.id
  peer_external_gateway_interface = 0
  shared_secret                   = "a secret message"
  router                          = google_compute_router.router.id
  vpn_gateway_interface           = 0
}

resource "google_compute_vpn_tunnel" "tunnel2" {
  name                            = "tunnel2"
  region                          = "us-central1"
  vpn_gateway                     = google_compute_ha_vpn_gateway.gateway.id
  peer_external_gateway           = google_compute_external_vpn_gateway.external_vpn_gw.id
  peer_external_gateway_interface = 0
  shared_secret                   = "a secret message"
  router                          = " ${google_compute_router.router.id}"
  vpn_gateway_interface           = 1
}

resource "google_compute_router_interface" "router_interface1" {
  name       = "router-interface1"
  router     = google_compute_router.router.name
  region     = "us-central1"
  ip_range   = "169.254.0.1/30"
  vpn_tunnel = google_compute_vpn_tunnel.tunnel1.name
}

resource "google_compute_router_peer" "router_peer1" {
  name                      = "router-peer1"
  router                    = google_compute_router.router.name
  region                    = "us-central1"
  peer_ip_address           = "169.254.0.2"
  peer_asn                  = 64515
  advertised_route_priority = 100
  interface                 = google_compute_router_interface.router_interface1.name
}

resource "google_compute_router_interface" "router_interface2" {
  name       = "router-interface2"
  router     = google_compute_router.router.name
  region     = "us-central1"
  ip_range   = "169.254.1.1/30"
  vpn_tunnel = google_compute_vpn_tunnel.tunnel2.name
}

resource "google_compute_router_peer" "router_peer2" {
  name                      = "router-peer2"
  router                    = google_compute_router.router.name
  region                    = "us-central1"
  peer_ip_address           = "169.254.1.2"
  peer_asn                  = 64515
  advertised_route_priority = 100
  interface                 = google_compute_router_interface.router_interface2.name
}

resource "google_network_connectivity_spoke" "tunnel1" {
  name        = "vpn-tunnel-1-spoke"
  location    = "us-central1"
  description = "A sample spoke with a linked VPN Tunnel"
  labels = {
    label-one = "value-one"
  }
  hub = google_network_connectivity_hub.basic_hub.id
  linked_vpn_tunnels {
    uris                       = [google_compute_vpn_tunnel.tunnel1.self_link]
    site_to_site_data_transfer = true
    include_import_ranges      = ["ALL_IPV4_RANGES"]
  }
}

resource "google_network_connectivity_spoke" "tunnel2" {
  name        = "vpn-tunnel-2-spoke"
  location    = "us-central1"
  description = "A sample spoke with a linked VPN Tunnel"
  labels = {
    label-one = "value-one"
  }
  hub = google_network_connectivity_hub.basic_hub.id
  linked_vpn_tunnels {
    uris                       = [google_compute_vpn_tunnel.tunnel2.self_link]
    site_to_site_data_transfer = true
    include_import_ranges      = ["ALL_IPV4_RANGES"]
  }
}
```
<div class = "oics-button" style="float: right; margin: 0 0 -15px">
  <a href="https://console.cloud.google.com/cloudshell/open?cloudshell_git_repo=https%3A%2F%2Fgithub.com%2Fterraform-google-modules%2Fdocs-examples.git&cloudshell_image=gcr.io%2Fcloudshell-images%2Fcloudshell%3Alatest&cloudshell_print=.%2Fmotd&cloudshell_tutorial=.%2Ftutorial.md&cloudshell_working_dir=network_connectivity_spoke_interconnect_attachment_basic&open_in_editor=main.tf" target="_blank">
    <img alt="Open in Cloud Shell" src="//gstatic.com/cloudssh/images/open-btn.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>
## Example Usage - Network Connectivity Spoke Interconnect Attachment Basic


```hcl
resource "google_network_connectivity_hub" "basic_hub" {
  name        = "basic-hub1"
  description = "A sample hub"
  labels = {
    label-two = "value-one"
  }
}

resource "google_compute_network" "network" {
  name                    = "basic-network"
  auto_create_subnetworks = false
}

resource "google_compute_router" "router" {
  name    = "external-vpn-gateway"
  region  = "us-central1"
  network = google_compute_network.network.name
  bgp {
    asn = 16550
  }
}

resource "google_compute_interconnect_attachment" "interconnect-attachment" {
  name                     = "partner-interconnect1"
  edge_availability_domain = "AVAILABILITY_DOMAIN_1"
  type                     = "PARTNER"
  router                   = google_compute_router.router.id
  mtu                      = 1500
  region                   = "us-central1"
}

resource "google_network_connectivity_spoke" "primary" {
  name        = "interconnect-attachment-spoke"
  location    = "us-central1"
  description = "A sample spoke with a linked Interconnect Attachment"
  labels = {
    label-one = "value-one"
  }
  hub = google_network_connectivity_hub.basic_hub.id
  linked_interconnect_attachments {
    uris                       = [google_compute_interconnect_attachment.interconnect-attachment.self_link]
    site_to_site_data_transfer = true
    include_import_ranges      = ["ALL_IPV4_RANGES"]
  }
}
```
<div class = "oics-button" style="float: right; margin: 0 0 -15px">
  <a href="https://console.cloud.google.com/cloudshell/open?cloudshell_git_repo=https%3A%2F%2Fgithub.com%2Fterraform-google-modules%2Fdocs-examples.git&cloudshell_image=gcr.io%2Fcloudshell-images%2Fcloudshell%3Alatest&cloudshell_print=.%2Fmotd&cloudshell_tutorial=.%2Ftutorial.md&cloudshell_working_dir=network_connectivity_spoke_linked_producer_vpc_network_basic&open_in_editor=main.tf" target="_blank">
    <img alt="Open in Cloud Shell" src="//gstatic.com/cloudssh/images/open-btn.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>
## Example Usage - Network Connectivity Spoke Linked Producer Vpc Network Basic


```hcl
resource "google_compute_network" "network" {
  name                    = "net-spoke"
  auto_create_subnetworks = false
}

resource "google_compute_global_address" "address" {
  name          = "test-address"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.network.id
}

resource "google_service_networking_connection" "peering" {
  network                 = google_compute_network.network.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.address.name]
}

resource "google_network_connectivity_hub" "basic_hub" {
  name = "hub-basic"
}

resource "google_network_connectivity_spoke" "linked_vpc_spoke"  {
  name     = "vpc-spoke"
  location = "global"
  hub      = google_network_connectivity_hub.basic_hub.id
  linked_vpc_network {
    uri = google_compute_network.network.self_link
  }
}

resource "google_network_connectivity_spoke" "primary"  {
  name        = "producer-spoke"
  location    = "global"
  description = "A sample spoke with a linked router appliance instance"
  labels = {
    label-one = "value-one"
  }
  hub         = google_network_connectivity_hub.basic_hub.id
  linked_producer_vpc_network {
    network = google_compute_network.network.name
    peering = google_service_networking_connection.peering.peering
    exclude_export_ranges = [
    "198.51.100.0/24",
    "10.10.0.0/16"
    ]
  }
  depends_on  = [google_network_connectivity_spoke.linked_vpc_spoke]
}
```
<div class = "oics-button" style="float: right; margin: 0 0 -15px">
  <a href="https://console.cloud.google.com/cloudshell/open?cloudshell_git_repo=https%3A%2F%2Fgithub.com%2Fterraform-google-modules%2Fdocs-examples.git&cloudshell_image=gcr.io%2Fcloudshell-images%2Fcloudshell%3Alatest&cloudshell_print=.%2Fmotd&cloudshell_tutorial=.%2Ftutorial.md&cloudshell_working_dir=network_connectivity_spoke_center_group&open_in_editor=main.tf" target="_blank">
    <img alt="Open in Cloud Shell" src="//gstatic.com/cloudssh/images/open-btn.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>
## Example Usage - Network Connectivity Spoke Center Group


```hcl
resource "google_compute_network" "network" {
  name                    = "tf-net"
  auto_create_subnetworks = false
}

resource "google_network_connectivity_hub" "star_hub" {
  name = "hub-basic"
  preset_topology = "STAR"
}

resource "google_network_connectivity_group" "center_group" { 
  name = "center"  # (default , center , edge)
  hub  = google_network_connectivity_hub.star_hub.id
  auto_accept {
    auto_accept_projects = [
      "foo%{random_suffix}", 
      "bar%{random_suffix}", 
    ]
  }
}

resource "google_network_connectivity_spoke" "primary"  {
  name = "vpc-spoke"
  location = "global"
  description = "A sample spoke"
  labels = {
    label-one = "value-one"
  }
  hub = google_network_connectivity_hub.star_hub.id
  group  = google_network_connectivity_group.center_group.id

  linked_vpc_network {
    uri = google_compute_network.network.self_link
  }
}
```
<div class = "oics-button" style="float: right; margin: 0 0 -15px">
  <a href="https://console.cloud.google.com/cloudshell/open?cloudshell_git_repo=https%3A%2F%2Fgithub.com%2Fterraform-google-modules%2Fdocs-examples.git&cloudshell_image=gcr.io%2Fcloudshell-images%2Fcloudshell%3Alatest&cloudshell_print=.%2Fmotd&cloudshell_tutorial=.%2Ftutorial.md&cloudshell_working_dir=network_connectivity_spoke_linked_vpc_network_ipv6_support&open_in_editor=main.tf" target="_blank">
    <img alt="Open in Cloud Shell" src="//gstatic.com/cloudssh/images/open-btn.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>
## Example Usage - Network Connectivity Spoke Linked Vpc Network Ipv6 Support


```hcl
resource "google_compute_network" "network" {
  name                    = "net"
  auto_create_subnetworks = false
}

resource "google_network_connectivity_hub" "basic_hub" {
  name        = "hub1"
  description = "A sample hub"
  labels = {
    label-two = "value-one"
  }
}

resource "google_network_connectivity_spoke" "primary"  {
  name = "spoke1-ipv6"
  location = "global"
  description = "A sample spoke with a linked VPC that include export ranges of all IPv6"
  labels = {
    label-one = "value-one"
  }
  hub = google_network_connectivity_hub.basic_hub.id
  linked_vpc_network {
    include_export_ranges = [
      "ALL_IPV6_RANGES",
      "ALL_PRIVATE_IPV4_RANGES"
    ]
    uri = google_compute_network.network.self_link
  }
}
```

## Argument Reference

The following arguments are supported:


* `name` -
  (Required)
  Immutable. The name of the spoke. Spoke names must be unique.

* `hub` -
  (Required)
  Immutable. The URI of the hub that this spoke is attached to.

* `location` -
  (Required)
  The location for the resource


* `labels` -
  (Optional)
  Optional labels in key:value format. For more information about labels, see [Requirements for labels](https://cloud.google.com/resource-manager/docs/creating-managing-labels#requirements).
  **Note**: This field is non-authoritative, and will only manage the labels present in your configuration.
  Please refer to the field `effective_labels` for all of the labels present on the resource.

* `description` -
  (Optional)
  An optional description of the spoke.

* `group` -
  (Optional)
  The name of the group that this spoke is associated with.

* `linked_vpn_tunnels` -
  (Optional)
  The URIs of linked VPN tunnel resources
  Structure is [documented below](#nested_linked_vpn_tunnels).

* `linked_interconnect_attachments` -
  (Optional)
  A collection of VLAN attachment resources. These resources should be redundant attachments that all advertise the same prefixes to Google Cloud. Alternatively, in active/passive configurations, all attachments should be capable of advertising the same prefixes.
  Structure is [documented below](#nested_linked_interconnect_attachments).

* `linked_router_appliance_instances` -
  (Optional)
  The URIs of linked Router appliance resources
  Structure is [documented below](#nested_linked_router_appliance_instances).

* `linked_vpc_network` -
  (Optional)
  VPC network that is associated with the spoke.
  Structure is [documented below](#nested_linked_vpc_network).

* `linked_producer_vpc_network` -
  (Optional)
  Producer VPC network that is associated with the spoke.
  Structure is [documented below](#nested_linked_producer_vpc_network).

* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.



<a name="nested_linked_vpn_tunnels"></a>The `linked_vpn_tunnels` block supports:

* `uris` -
  (Required)
  The URIs of linked VPN tunnel resources.

* `site_to_site_data_transfer` -
  (Required)
  A value that controls whether site-to-site data transfer is enabled for these resources. Note that data transfer is available only in supported locations.

* `include_import_ranges` -
  (Optional)
  IP ranges allowed to be included during import from hub (does not control transit connectivity).
  The only allowed value for now is "ALL_IPV4_RANGES".

<a name="nested_linked_interconnect_attachments"></a>The `linked_interconnect_attachments` block supports:

* `uris` -
  (Required)
  The URIs of linked interconnect attachment resources

* `site_to_site_data_transfer` -
  (Required)
  A value that controls whether site-to-site data transfer is enabled for these resources. Note that data transfer is available only in supported locations.

* `include_import_ranges` -
  (Optional)
  IP ranges allowed to be included during import from hub (does not control transit connectivity).
  The only allowed value for now is "ALL_IPV4_RANGES".

<a name="nested_linked_router_appliance_instances"></a>The `linked_router_appliance_instances` block supports:

* `instances` -
  (Required)
  The list of router appliance instances
  Structure is [documented below](#nested_linked_router_appliance_instances_instances).

* `site_to_site_data_transfer` -
  (Required)
  A value that controls whether site-to-site data transfer is enabled for these resources. Note that data transfer is available only in supported locations.

* `include_import_ranges` -
  (Optional)
  IP ranges allowed to be included during import from hub (does not control transit connectivity).
  The only allowed value for now is "ALL_IPV4_RANGES".


<a name="nested_linked_router_appliance_instances_instances"></a>The `instances` block supports:

* `virtual_machine` -
  (Required)
  The URI of the virtual machine resource

* `ip_address` -
  (Required)
  The IP address on the VM to use for peering.

<a name="nested_linked_vpc_network"></a>The `linked_vpc_network` block supports:

* `uri` -
  (Required)
  The URI of the VPC network resource.

* `exclude_export_ranges` -
  (Optional)
  IP ranges encompassing the subnets to be excluded from peering.

* `include_export_ranges` -
  (Optional)
  IP ranges allowed to be included from peering.

<a name="nested_linked_producer_vpc_network"></a>The `linked_producer_vpc_network` block supports:

* `network` -
  (Required)
  The URI of the Service Consumer VPC that the Producer VPC is peered with.

* `peering` -
  (Required)
  The name of the VPC peering between the Service Consumer VPC and the Producer VPC (defined in the Tenant project) which is added to the NCC hub. This peering must be in ACTIVE state.

* `producer_network` -
  (Output)
  The URI of the Producer VPC.

* `include_export_ranges` -
  (Optional)
  IP ranges allowed to be included from peering.

* `exclude_export_ranges` -
  (Optional)
  IP ranges encompassing the subnets to be excluded from peering.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{location}}/spokes/{{name}}`

* `create_time` -
  Output only. The time the spoke was created.

* `update_time` -
  Output only. The time the spoke was last updated.

* `unique_id` -
  Output only. The Google-generated UUID for the spoke. This value is unique across all spoke resources. If a spoke is deleted and another with the same name is created, the new spoke is assigned a different unique_id.

* `state` -
  Output only. The current lifecycle state of this spoke.

* `reasons` -
  The reasons for the current state in the lifecycle
  Structure is [documented below](#nested_reasons).

* `terraform_labels` -
  The combination of labels configured directly on the resource
   and default labels configured on the provider.

* `effective_labels` -
  All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.


<a name="nested_reasons"></a>The `reasons` block contains:

* `code` -
  (Optional)
  The code associated with this reason.

* `message` -
  (Optional)
  Human-readable details about this reason.

* `user_details` -
  (Optional)
  Additional information provided by the user in the RejectSpoke call.

## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import


Spoke can be imported using any of these accepted formats:

* `projects/{{project}}/locations/{{location}}/spokes/{{name}}`
* `{{project}}/{{location}}/{{name}}`
* `{{location}}/{{name}}`


In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Spoke using one of the formats above. For example:

```tf
import {
  id = "projects/{{project}}/locations/{{location}}/spokes/{{name}}"
  to = google_network_connectivity_spoke.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), Spoke can be imported using one of the formats above. For example:

```
$ terraform import google_network_connectivity_spoke.default projects/{{project}}/locations/{{location}}/spokes/{{name}}
$ terraform import google_network_connectivity_spoke.default {{project}}/{{location}}/{{name}}
$ terraform import google_network_connectivity_spoke.default {{location}}/{{name}}
```

## User Project Overrides

This resource supports [User Project Overrides](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/provider_reference#user_project_override).
