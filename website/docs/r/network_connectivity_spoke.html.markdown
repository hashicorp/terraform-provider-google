---
# ----------------------------------------------------------------------------
#
#     ***     AUTO GENERATED CODE    ***    Type: DCL     ***
#
# ----------------------------------------------------------------------------
#
#     This file is managed by Magic Modules (https:#github.com/GoogleCloudPlatform/magic-modules)
#     and is based on the DCL (https:#github.com/GoogleCloudPlatform/declarative-resource-client-library).
#     Changes will need to be made to the DCL or Magic Modules instead of here.
#
#     We are not currently able to accept contributions to this file. If changes
#     are required, please file an issue at https:#github.com/hashicorp/terraform-provider-google/issues/new/choose
#
# ----------------------------------------------------------------------------
subcategory: "NetworkConnectivity"
description: |-
  The NetworkConnectivity Spoke resource
---

# google_network_connectivity_spoke

The NetworkConnectivity Spoke resource

## Example Usage - linked_vpc_network
```hcl

resource "google_compute_network" "network" {
  name                    = "network"
  auto_create_subnetworks = false
}

resource "google_network_connectivity_hub" "basic_hub" {
  name        = "hub"
  description = "A sample hub"
  labels = {
    label-two = "value-one"
  }
}

resource "google_network_connectivity_spoke" "primary" {
  name = "name"
  location = "global"
  description = "A sample spoke with a linked routher appliance instance"
  labels = {
    label-one = "value-one"
  }
  hub = google_network_connectivity_hub.basic_hub.id
  linked_vpc_network {
    exclude_export_ranges = [
      "198.51.100.0/24",
      "10.10.0.0/16"
    ]
    uri = google_compute_network.network.self_link
  }
}
```
## Example Usage - router_appliance
```hcl

resource "google_compute_network" "network" {
  name                    = "network"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnetwork" {
  name          = "subnet"
  ip_cidr_range = "10.0.0.0/28"
  region        = "us-west1"
  network       = google_compute_network.network.self_link
}

resource "google_compute_instance" "instance" {
  name         = "instance"
  machine_type = "e2-medium"
  can_ip_forward = true
  zone         = "us-west1-a"

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
  name        = "hub"
  description = "A sample hub"
  labels = {
    label-two = "value-one"
  }
}

resource "google_network_connectivity_spoke" "primary" {
  name = "name"
  location = "us-west1"
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
  }
}
```

## Argument Reference

The following arguments are supported:

* `hub` -
  (Required)
  Immutable. The URI of the hub that this spoke is attached to.
  
* `location` -
  (Required)
  The location for the resource
  
* `name` -
  (Required)
  Immutable. The name of the spoke. Spoke names must be unique.
  


The `instances` block supports:
    
* `ip_address` -
  (Optional)
  The IP address on the VM to use for peering.
    
* `virtual_machine` -
  (Optional)
  The URI of the virtual machine resource
    
- - -

* `description` -
  (Optional)
  An optional description of the spoke.
  
* `labels` -
  (Optional)
  Optional labels in key:value format. For more information about labels, see [Requirements for labels](https://cloud.google.com/resource-manager/docs/creating-managing-labels#requirements).

**Note**: This field is non-authoritative, and will only manage the labels present in your configuration. Please refer to the field `effective_labels` for all of the labels present on the resource.
  
* `linked_interconnect_attachments` -
  (Optional)
  A collection of VLAN attachment resources. These resources should be redundant attachments that all advertise the same prefixes to Google Cloud. Alternatively, in active/passive configurations, all attachments should be capable of advertising the same prefixes.
  
* `linked_router_appliance_instances` -
  (Optional)
  The URIs of linked Router appliance resources
  
* `linked_vpc_network` -
  (Optional)
  VPC network that is associated with the spoke.
  
* `linked_vpn_tunnels` -
  (Optional)
  The URIs of linked VPN tunnel resources
  
* `project` -
  (Optional)
  The project for the resource
  


The `linked_interconnect_attachments` block supports:
    
* `site_to_site_data_transfer` -
  (Required)
  A value that controls whether site-to-site data transfer is enabled for these resources. Note that data transfer is available only in supported locations.
    
* `uris` -
  (Required)
  The URIs of linked interconnect attachment resources
    
The `linked_router_appliance_instances` block supports:
    
* `instances` -
  (Required)
  The list of router appliance instances
    
* `site_to_site_data_transfer` -
  (Required)
  A value that controls whether site-to-site data transfer is enabled for these resources. Note that data transfer is available only in supported locations.
    
The `linked_vpc_network` block supports:
    
* `exclude_export_ranges` -
  (Optional)
  IP ranges encompassing the subnets to be excluded from peering.
    
* `uri` -
  (Required)
  The URI of the VPC network resource.
    
The `linked_vpn_tunnels` block supports:
    
* `site_to_site_data_transfer` -
  (Required)
  A value that controls whether site-to-site data transfer is enabled for these resources. Note that data transfer is available only in supported locations.
    
* `uris` -
  (Required)
  The URIs of linked VPN tunnel resources.
    
## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{location}}/spokes/{{name}}`

* `create_time` -
  Output only. The time the spoke was created.
  
* `effective_labels` -
  All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.
  
* `state` -
  Output only. The current lifecycle state of this spoke. Possible values: STATE_UNSPECIFIED, CREATING, ACTIVE, DELETING
  
* `terraform_labels` -
  The combination of labels configured directly on the resource and default labels configured on the provider.
  
* `unique_id` -
  Output only. The Google-generated UUID for the spoke. This value is unique across all spoke resources. If a spoke is deleted and another with the same name is created, the new spoke is assigned a different unique_id.
  
* `update_time` -
  Output only. The time the spoke was last updated.
  
## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

Spoke can be imported using any of these accepted formats:

```
$ terraform import google_network_connectivity_spoke.default projects/{{project}}/locations/{{location}}/spokes/{{name}}
$ terraform import google_network_connectivity_spoke.default {{project}}/{{location}}/{{name}}
$ terraform import google_network_connectivity_spoke.default {{location}}/{{name}}
```



