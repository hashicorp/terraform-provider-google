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
  The NetworkConnectivity Hub resource
---

# google_network_connectivity_hub

The NetworkConnectivity Hub resource

## Example Usage - basic_hub
A basic test of a networkconnectivity hub
```hcl
resource "google_network_connectivity_hub" "primary" {
  name        = "hub"
  description = "A sample hub"

  labels = {
    label-one = "value-one"
  }

  project = "my-project-name"
}


```

## Argument Reference

The following arguments are supported:

* `name` -
  (Required)
  Immutable. The name of the hub. Hub names must be unique. They use the following form: `projects/{project_number}/locations/global/hubs/{hub_id}`
  


- - -

* `description` -
  (Optional)
  An optional description of the hub.
  
* `labels` -
  (Optional)
  Optional labels in key:value format. For more information about labels, see [Requirements for labels](https://cloud.google.com/resource-manager/docs/creating-managing-labels#requirements).

**Note**: This field is non-authoritative, and will only manage the labels present in your configuration. Please refer to the field `effective_labels` for all of the labels present on the resource.
  
* `project` -
  (Optional)
  The project for the resource
  


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/global/hubs/{{name}}`

* `create_time` -
  Output only. The time the hub was created.
  
* `effective_labels` -
  All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.
  
* `routing_vpcs` -
  The VPC network associated with this hub's spokes. All of the VPN tunnels, VLAN attachments, and router appliance instances referenced by this hub's spokes must belong to this VPC network. This field is read-only. Network Connectivity Center automatically populates it based on the set of spokes attached to the hub.
  
* `state` -
  Output only. The current lifecycle state of this hub. Possible values: STATE_UNSPECIFIED, CREATING, ACTIVE, DELETING
  
* `unique_id` -
  Output only. The Google-generated UUID for the hub. This value is unique across all hub resources. If a hub is deleted and another with the same name is created, the new hub is assigned a different unique_id.
  
* `update_time` -
  Output only. The time the hub was last updated.
  
## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

Hub can be imported using any of these accepted formats:

```
$ terraform import google_network_connectivity_hub.default projects/{{project}}/locations/global/hubs/{{name}}
$ terraform import google_network_connectivity_hub.default {{project}}/{{name}}
$ terraform import google_network_connectivity_hub.default {{name}}
```



