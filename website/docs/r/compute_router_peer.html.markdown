---
subcategory: "Compute Engine"
description: |-
  BGP information that must be configured into the routing stack to
  establish BGP peering.
---

# google_compute_router_peer

BGP information that must be configured into the routing stack to
establish BGP peering. This information must specify the peer ASN
and either the interface name, IP address, or peer IP address.
Please refer to RFC4273.


To get more information about RouterBgpPeer, see:

* [API documentation](https://cloud.google.com/compute/docs/reference/rest/v1/routers)
* How-to Guides
    * [Google Cloud Router](https://cloud.google.com/router/docs/)

## Example Usage - Router Peer Basic


```hcl
resource "google_compute_router_peer" "peer" {
  name                      = "my-router-peer"
  router                    = "my-router"
  region                    = "us-central1"
  peer_asn                  = 65513
  advertised_route_priority = 100
  interface                 = "interface-1"
}
```
## Example Usage - Router Peer Disabled


```hcl
resource "google_compute_router_peer" "peer" {
  name                      = "my-router-peer"
  router                    = "my-router"
  region                    = "us-central1"
  peer_ip_address           = "169.254.1.2"
  peer_asn                  = 65513
  advertised_route_priority = 100
  interface                 = "interface-1"
  enable                    = false
}
```
## Example Usage - Router Peer Bfd


```hcl
resource "google_compute_router_peer" "peer" {
  name                      = "my-router-peer"
  router                    = "my-router"
  region                    = "us-central1"
  peer_ip_address           = "169.254.1.2"
  peer_asn                  = 65513
  advertised_route_priority = 100
  interface                 = "interface-1"

  bfd {
    min_receive_interval        = 1000
    min_transmit_interval       = 1000
    multiplier                  = 5
    session_initialization_mode = "ACTIVE"
  }
}
```
<div class = "oics-button" style="float: right; margin: 0 0 -15px">
  <a href="https://console.cloud.google.com/cloudshell/open?cloudshell_git_repo=https%3A%2F%2Fgithub.com%2Fterraform-google-modules%2Fdocs-examples.git&cloudshell_working_dir=router_peer_router_appliance&cloudshell_image=gcr.io%2Fcloudshell-images%2Fcloudshell%3Alatest&open_in_editor=main.tf&cloudshell_print=.%2Fmotd&cloudshell_tutorial=.%2Ftutorial.md" target="_blank">
    <img alt="Open in Cloud Shell" src="//gstatic.com/cloudssh/images/open-btn.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>
## Example Usage - Router Peer Router Appliance


```hcl
resource "google_compute_network" "network" {
  name                    = "my-router-net"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnetwork" {
  name          = "my-router-sub"
  network       = google_compute_network.network.self_link
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
}

resource "google_compute_address" "addr_intf" {
  name         = "my-router-addr-intf"
  region       = google_compute_subnetwork.subnetwork.region
  subnetwork   = google_compute_subnetwork.subnetwork.id
  address_type = "INTERNAL"
}

resource "google_compute_address" "addr_intf_redundant" {
  name         = "my-router-addr-intf-red"
  region       = google_compute_subnetwork.subnetwork.region
  subnetwork   = google_compute_subnetwork.subnetwork.id
  address_type = "INTERNAL"
}

resource "google_compute_address" "addr_peer" {
  name         = "my-router-addr-peer"
  region       = google_compute_subnetwork.subnetwork.region
  subnetwork   = google_compute_subnetwork.subnetwork.id
  address_type = "INTERNAL"
}

resource "google_compute_instance" "instance" {
  name           = "router-appliance"
  zone           = "us-central1-a"
  machine_type   = "e2-medium"
  can_ip_forward = true

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }

  network_interface {
    network_ip = google_compute_address.addr_peer.address
    subnetwork = google_compute_subnetwork.subnetwork.self_link
  }
}

resource "google_network_connectivity_hub" "hub" {
  name = "my-router-hub"
}

resource "google_network_connectivity_spoke" "spoke" {
  name     = "my-router-spoke"
  location = google_compute_subnetwork.subnetwork.region
  hub      = google_network_connectivity_hub.hub.id

  linked_router_appliance_instances {
    instances {
      virtual_machine = google_compute_instance.instance.self_link
      ip_address      = google_compute_address.addr_peer.address
    }
    site_to_site_data_transfer = false
  }
}

resource "google_compute_router" "router" {
  name    = "my-router-router"
  region  = google_compute_subnetwork.subnetwork.region
  network = google_compute_network.network.self_link
  bgp {
    asn = 64514
  }
}

resource "google_compute_router_interface" "interface_redundant" {
  name               = "my-router-intf-red"
  region             = google_compute_router.router.region
  router             = google_compute_router.router.name
  subnetwork         = google_compute_subnetwork.subnetwork.self_link
  private_ip_address = google_compute_address.addr_intf_redundant.address
}

resource "google_compute_router_interface" "interface" {
  name                = "my-router-intf"
  region              = google_compute_router.router.region
  router              = google_compute_router.router.name
  subnetwork          = google_compute_subnetwork.subnetwork.self_link
  private_ip_address  = google_compute_address.addr_intf.address
  redundant_interface = google_compute_router_interface.interface_redundant.name
}

resource "google_compute_router_peer" "peer" {
  name                      = "my-router-peer"
  router                    = google_compute_router.router.name
  region                    = google_compute_router.router.region
  interface                 = google_compute_router_interface.interface.name
  router_appliance_instance = google_compute_instance.instance.self_link
  peer_asn                  = 65513
  peer_ip_address           = google_compute_address.addr_peer.address
}
```

## Example Usage - Router Peer md5 authentication key


```hcl
  resource "google_compute_router_peer" "foobar" {
    name                      = "%s-peer"
    router                    = google_compute_router.foobar.name
    region                    = google_compute_router.foobar.region
    peer_asn                  = 65515
    advertised_route_priority = 100
    interface                 = google_compute_router_interface.foobar.name
    peer_ip_address           = "169.254.3.2"
    md5_authentication_key {
      name = "%s-peer-key"
      key = "%s-peer-key-value"
    }
  }
```

## Example Usage - Router peer export and import policies

```hcl
  resource "google_compute_network" "network" {
  provider = google-beta
  name = "my-router-net"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnetwork" {
  provider = google-beta
  name          = "my-router-subnet"
  network       = google_compute_network.network.self_link
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
}

resource "google_compute_address" "address" {
  provider = google-beta
  name   = "my-router"
  region = google_compute_subnetwork.subnetwork.region
}

resource "google_compute_ha_vpn_gateway" "vpn_gateway" {
  provider = google-beta
  name    = "my-router-gateway"
  network = google_compute_network.network.self_link
  region  = google_compute_subnetwork.subnetwork.region
}

resource "google_compute_external_vpn_gateway" "external_gateway" {
  provider = google-beta
  name            = "my-router-external-gateway"
  redundancy_type = "SINGLE_IP_INTERNALLY_REDUNDANT"
  description     = "An externally managed VPN gateway"
  interface {
    id         = 0
    ip_address = "8.8.8.8"
  }
}

resource "google_compute_router" "router" {
  provider = google-beta
  name    = "my-router"
  region  = google_compute_subnetwork.subnetwork.region
  network = google_compute_network.network.self_link
  bgp {
    asn = 64514
  }
}

resource "google_compute_vpn_tunnel" "vpn_tunnel" {
  provider = google-beta
  name               = "my-router"
  region             = google_compute_subnetwork.subnetwork.region
  vpn_gateway = google_compute_ha_vpn_gateway.vpn_gateway.id
  peer_external_gateway           = google_compute_external_vpn_gateway.external_gateway.id
  peer_external_gateway_interface = 0  
  shared_secret      = "unguessable"
  router             = google_compute_router.router.name
  vpn_gateway_interface           = 0
}

resource "google_compute_router_interface" "router_interface" {
  provider = google-beta
  name       = "my-router"
  router     = google_compute_router.router.name
  region     = google_compute_router.router.region
  vpn_tunnel = google_compute_vpn_tunnel.vpn_tunnel.name
}

resource "google_compute_router_route_policy" "rp-export" {
  provider = google-beta
	name = "my-router-rp-export"
  router = google_compute_router.router.name
  region = google_compute_router.router.region
  type = "ROUTE_POLICY_TYPE_EXPORT"
	terms {
    priority = 2
    match {
      expression = "destination == '10.0.0.0/12'"
      title      = "export_expression"
      description = "acceptance expression for export"
    }
    actions {
      expression = "accept()"
    }
  }
  depends_on = [
    google_compute_router_interface.router_interface
  ]
}

resource "google_compute_router_route_policy" "rp-import" {
  provider = google-beta
  name = "my-router-rp-import"
  router = google_compute_router.router.name
  region = google_compute_router.router.region
	type = "ROUTE_POLICY_TYPE_IMPORT"
  terms {
    priority = 1
    match {
      expression = "destination == '10.0.0.0/12'"
      title      = "import_expression"
      description = "acceptance expression for import"
	  }
    actions {
      expression = "accept()"
    }
  }
  depends_on = [
    google_compute_router_interface.router_interface, google_compute_router_route_policy.rp-export
  ]
}

resource "google_compute_router_peer" "router_peer" {
  provider = google-beta
  name                      = "my-router-peer"
  router                    = google_compute_router.router.name
  region                    = google_compute_router.router.region
  peer_asn                  = 65515
  advertised_route_priority = 100
  interface                 = google_compute_router_interface.router_interface.name
  md5_authentication_key {
    name = "my-router-peer-key"
    key = "my-router-peer-key-value"
  }
  import_policies           = [google_compute_router_route_policy.rp-import.name]
  export_policies           = [google_compute_router_route_policy.rp-export.name]
  depends_on = [
    google_compute_router_route_policy.rp-export, google_compute_router_route_policy.rp-import, google_compute_router_interface.router_interface
  ]
}
```

## Argument Reference

The following arguments are supported:


* `name` -
  (Required)
  Name of this BGP peer. The name must be 1-63 characters long,
  and comply with RFC1035. Specifically, the name must be 1-63 characters
  long and match the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` which
  means the first character must be a lowercase letter, and all
  following characters must be a dash, lowercase letter, or digit,
  except the last character, which cannot be a dash.

* `interface` -
  (Required)
  Name of the interface the BGP peer is associated with.

* `peer_asn` -
  (Required)
  Peer BGP Autonomous System Number (ASN).
  Each BGP interface may use a different value.

* `router` -
  (Required)
  The name of the Cloud Router in which this BgpPeer will be configured.


- - -


* `ip_address` -
  (Optional)
  IP address of the interface inside Google Cloud Platform.
  Only IPv4 is supported.

* `peer_ip_address` -
  (Optional)
  IP address of the BGP interface outside Google Cloud Platform.
  Only IPv4 is supported. Required if `ip_address` is set.

* `advertised_route_priority` -
  (Optional)
  The priority of routes advertised to this BGP peer.
  Where there is more than one matching route of maximum
  length, the routes with the lowest priority value win.

* `advertise_mode` -
  (Optional)
  User-specified flag to indicate which mode to use for advertisement.
  Valid values of this enum field are: `DEFAULT`, `CUSTOM`
  Default value is `DEFAULT`.
  Possible values are: `DEFAULT`, `CUSTOM`.

* `advertised_groups` -
  (Optional)
  User-specified list of prefix groups to advertise in custom
  mode, which currently supports the following option:
  * `ALL_SUBNETS`: Advertises all of the router's own VPC subnets.
  This excludes any routes learned for subnets that use VPC Network
  Peering.

  Note that this field can only be populated if advertiseMode is `CUSTOM`
  and overrides the list defined for the router (in the "bgp" message).
  These groups are advertised in addition to any specified prefixes.
  Leave this field blank to advertise no custom groups.

* `advertised_ip_ranges` -
  (Optional)
  User-specified list of individual IP ranges to advertise in
  custom mode. This field can only be populated if advertiseMode
  is `CUSTOM` and is advertised to all peers of the router. These IP
  ranges will be advertised in addition to any specified groups.
  Leave this field blank to advertise no custom IP ranges.
  Structure is [documented below](#nested_advertised_ip_ranges).

* `bfd` -
  (Optional)
  BFD configuration for the BGP peering.
  Structure is [documented below](#nested_bfd).

* `enable` -
  (Optional)
  The status of the BGP peer connection. If set to false, any active session
  with the peer is terminated and all associated routing information is removed.
  If set to true, the peer connection can be established with routing information.
  The default is true.

* `router_appliance_instance` -
  (Optional)
  The URI of the VM instance that is used as third-party router appliances
  such as Next Gen Firewalls, Virtual Routers, or Router Appliances.
  The VM instance must be located in zones contained in the same region as
  this Cloud Router. The VM instance is the peer side of the BGP session.

* `enable_ipv6` -
  (Optional)
  Enable IPv6 traffic over BGP Peer. If not specified, it is disabled by default.

* `enable_ipv4` -
  (Optional)
  Enable IPv4 traffic over BGP Peer. It is enabled by default if the peerIpAddress is version 4.

* `ipv6_nexthop_address` -
  (Optional)
  IPv6 address of the interface inside Google Cloud Platform.
  The address must be in the range 2600:2d00:0:2::/64 or 2600:2d00:0:3::/64.
  If you do not specify the next hop addresses, Google Cloud automatically
  assigns unused addresses from the 2600:2d00:0:2::/64 or 2600:2d00:0:3::/64 range for you.

* `ipv4_nexthop_address` -
  (Optional)
  IPv4 address of the interface inside Google Cloud Platform.

* `peer_ipv6_nexthop_address` -
  (Optional)
  IPv6 address of the BGP interface outside Google Cloud Platform.
  The address must be in the range 2600:2d00:0:2::/64 or 2600:2d00:0:3::/64.
  If you do not specify the next hop addresses, Google Cloud automatically
  assigns unused addresses from the 2600:2d00:0:2::/64 or 2600:2d00:0:3::/64 range for you.

* `peer_ipv4_nexthop_address` -
  (Optional)
  IPv4 address of the BGP interface outside Google Cloud Platform.

*  `export_policies` -
  (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html)) 
  routers.list of export policies applied to this peer, in the order they must be evaluated. 
  The name must correspond to an existing policy that has ROUTE_POLICY_TYPE_EXPORT type.

*  `import_policies` -
  (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html)) 
  routers.list of import policies applied to this peer, in the order they must be evaluated. 
  The name must correspond to an existing policy that has ROUTE_POLICY_TYPE_IMPORT type.

* `region` -
  (Optional)
  Region where the router and BgpPeer reside.
  If it is not provided, the provider region is used.

* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.

* `md5_authentication_key` - (Optional) Configuration for MD5 authentication on the BGP session.
  Structure is [documented below](#nested_md5_authentication_key).

<a name="nested_advertised_ip_ranges"></a>The `advertised_ip_ranges` block supports:

* `range` -
  (Required)
  The IP range to advertise. The value must be a
  CIDR-formatted string.

* `description` -
  (Optional)
  User-specified description for the IP range.

<a name="nested_bfd"></a>The `bfd` block supports:

* `session_initialization_mode` -
  (Required)
  The BFD session initialization mode for this BGP peer.
  If set to `ACTIVE`, the Cloud Router will initiate the BFD session
  for this BGP peer. If set to `PASSIVE`, the Cloud Router will wait
  for the peer router to initiate the BFD session for this BGP peer.
  If set to `DISABLED`, BFD is disabled for this BGP peer.
  Possible values are: `ACTIVE`, `DISABLED`, `PASSIVE`.

* `min_transmit_interval` -
  (Optional)
  The minimum interval, in milliseconds, between BFD control packets
  transmitted to the peer router. The actual value is negotiated
  between the two routers and is equal to the greater of this value
  and the corresponding receive interval of the other router. If set,
  this value must be between 1000 and 30000.

* `min_receive_interval` -
  (Optional)
  The minimum interval, in milliseconds, between BFD control packets
  received from the peer router. The actual value is negotiated
  between the two routers and is equal to the greater of this value
  and the transmit interval of the other router. If set, this value
  must be between 1000 and 30000.

* `multiplier` -
  (Optional)
  The number of consecutive BFD packets that must be missed before
  BFD declares that a peer is unavailable. If set, the value must
  be a value between 5 and 16.

<a name="nested_md5_authentication_key"></a>The `md5_authentication_key` block supports:

* `name` -
  (Required)
  Name used to identify the key. Must be unique within a router. Must comply with RFC1035.

* `key` -
  (Required, Input Only)
  The MD5 authentication key for this BGP peer. Maximum length is 80 characters. Can only contain printable ASCII characters

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/regions/{{region}}/routers/{{router}}/{{name}}`

* `management_type` -
  The resource that configures and manages this BGP peer.
  * `MANAGED_BY_USER` is the default value and can be managed by
  you or other users
  * `MANAGED_BY_ATTACHMENT` is a BGP peer that is configured and
  managed by Cloud Interconnect, specifically by an
  InterconnectAttachment of type PARTNER. Google automatically
  creates, updates, and deletes this type of BGP peer when the
  PARTNER InterconnectAttachment is created, updated,
  or deleted.


## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import


RouterBgpPeer can be imported using any of these accepted formats:

* `projects/{{project}}/regions/{{region}}/routers/{{router}}/{{name}}`
* `{{project}}/{{region}}/{{router}}/{{name}}`
* `{{region}}/{{router}}/{{name}}`
* `{{router}}/{{name}}`


In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import RouterBgpPeer using one of the formats above. For example:

```tf
import {
  id = "projects/{{project}}/regions/{{region}}/routers/{{router}}/{{name}}"
  to = google_compute_router_peer.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), RouterBgpPeer can be imported using one of the formats above. For example:

```
$ terraform import google_compute_router_peer.default projects/{{project}}/regions/{{region}}/routers/{{router}}/{{name}}
$ terraform import google_compute_router_peer.default {{project}}/{{region}}/{{router}}/{{name}}
$ terraform import google_compute_router_peer.default {{region}}/{{router}}/{{name}}
$ terraform import google_compute_router_peer.default {{router}}/{{name}}
```

## User Project Overrides

This resource supports [User Project Overrides](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/provider_reference#user_project_override).
