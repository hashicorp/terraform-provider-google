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
subcategory: "ContainerAzure"
layout: "google"
page_title: "Google: google_container_azure_cluster"
description: |-
  An Anthos cluster running on Azure.
---

# google_container_azure_cluster

An Anthos cluster running on Azure.

For more information, see:
* [Multicloud overview](https://cloud.google.com/anthos/clusters/docs/multi-cloud)
## Example Usage - basic_azure_cluster
A basic example of a containerazure azure cluster
```hcl
data "google_container_azure_versions" "versions" {
  project = "my-project-name"
  location = "us-west1"
}

resource "google_container_azure_cluster" "primary" {
  authorization {
    admin_users {
      username = "mmv2@google.com"
    }
  }

  azure_region = "westus2"
  client       = "projects/my-project-number/locations/us-west1/azureClients/${google_container_azure_client.basic.name}"

  control_plane {
    ssh_config {
      authorized_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC8yaayO6lnb2v+SedxUMa2c8vtIEzCzBjM3EJJsv8Vm9zUDWR7dXWKoNGARUb2mNGXASvI6mFIDXTIlkQ0poDEPpMaXR0g2cb5xT8jAAJq7fqXL3+0rcJhY/uigQ+MrT6s+ub0BFVbsmGHNrMQttXX9gtmwkeAEvj3mra9e5pkNf90qlKnZz6U0SVArxVsLx07vHPHDIYrl0OPG4zUREF52igbBPiNrHJFDQJT/4YlDMJmo/QT/A1D6n9ocemvZSzhRx15/Arjowhr+VVKSbaxzPtEfY0oIg2SrqJnnr/l3Du5qIefwh5VmCZe4xopPUaDDoOIEFriZ88sB+3zz8ib8sk8zJJQCgeP78tQvXCgS+4e5W3TUg9mxjB6KjXTyHIVhDZqhqde0OI3Fy1UuVzRUwnBaLjBnAwP5EoFQGRmDYk/rEYe7HTmovLeEBUDQocBQKT4Ripm/xJkkWY7B07K/tfo56dGUCkvyIVXKBInCh+dLK7gZapnd4UWkY0xBYcwo1geMLRq58iFTLA2j/JmpmHXp7m0l7jJii7d44uD3tTIFYThn7NlOnvhLim/YcBK07GMGIN7XwrrKZKmxXaspw6KBWVhzuw1UPxctxshYEaMLfFg/bwOw8HvMPr9VtrElpSB7oiOh91PDIPdPBgHCi7N2QgQ5l/ZDBHieSpNrQ== thomasrodgers"
    }

    subnet_id = "/subscriptions/12345678-1234-1234-1234-123456789111/resourceGroups/my--dev-byo/providers/Microsoft.Network/virtualNetworks/my--dev-vnet/subnets/default"
    version   = "${data.google_container_azure_versions.versions.valid_versions[0]}"
  }

  fleet {
    project = "my-project-number"
  }

  location = "us-west1"
  name     = "name"

  networking {
    pod_address_cidr_blocks     = ["10.200.0.0/16"]
    service_address_cidr_blocks = ["10.32.0.0/24"]
    virtual_network_id          = "/subscriptions/12345678-1234-1234-1234-123456789111/resourceGroups/my--dev-byo/providers/Microsoft.Network/virtualNetworks/my--dev-vnet"
  }

  resource_group_id = "/subscriptions/12345678-1234-1234-1234-123456789111/resourceGroups/my--dev-cluster"
  project           = "my-project-name"
}

resource "google_container_azure_client" "basic" {
  application_id = "12345678-1234-1234-1234-123456789111"
  location       = "us-west1"
  name           = "client-name"
  tenant_id      = "12345678-1234-1234-1234-123456789111"
  project        = "my-project-name"
}


```

## Argument Reference

The following arguments are supported:

* `authorization` -
  (Required)
  Configuration related to the cluster RBAC settings.
  
* `azure_region` -
  (Required)
  The Azure region where the cluster runs. Each Google Cloud region supports a subset of nearby Azure regions. You can call to list all supported Azure regions within a given Google Cloud region.
  
* `client` -
  (Required)
  Name of the AzureClient. The `AzureClient` resource must reside on the same GCP project and region as the `AzureCluster`. `AzureClient` names are formatted as `projects/<project-number>/locations/<region>/azureClients/<client-id>`. See Resource Names (https:cloud.google.com/apis/design/resource_names) for more details on Google Cloud resource names.
  
* `control_plane` -
  (Required)
  Configuration related to the cluster control plane.
  
* `fleet` -
  (Required)
  Fleet configuration.
  
* `location` -
  (Required)
  The location for the resource
  
* `name` -
  (Required)
  The name of this resource.
  
* `networking` -
  (Required)
  Cluster-wide networking configuration.
  
* `resource_group_id` -
  (Required)
  The ARM ID of the resource group where the cluster resources are deployed. For example: `/subscriptions/*/resourceGroups/*`
  


The `authorization` block supports:
    
* `admin_users` -
  (Required)
  Users that can perform operations as a cluster admin. A new ClusterRoleBinding will be created to grant the cluster-admin ClusterRole to the users. At most one user can be specified. For more info on RBAC, see https://kubernetes.io/docs/reference/access-authn-authz/rbac/#user-facing-roles
    
The `admin_users` block supports:
    
* `username` -
  (Required)
  The name of the user, e.g. `my-gcp-id@gmail.com`.
    
The `control_plane` block supports:
    
* `database_encryption` -
  (Optional)
  Optional. Configuration related to application-layer secrets encryption.
    
* `main_volume` -
  (Optional)
  Optional. Configuration related to the main volume provisioned for each control plane replica. The main volume is in charge of storing all of the cluster's etcd state. When unspecified, it defaults to a 8-GiB Azure Disk.
    
* `proxy_config` -
  (Optional)
  Proxy configuration for outbound HTTP(S) traffic.
    
* `replica_placements` -
  (Optional)
  Configuration for where to place the control plane replicas. Up to three replica placement instances can be specified. If replica_placements is set, the replica placement instances will be applied to the three control plane replicas as evenly as possible.
    
* `root_volume` -
  (Optional)
  Optional. Configuration related to the root volume provisioned for each control plane replica. When unspecified, it defaults to 32-GiB Azure Disk.
    
* `ssh_config` -
  (Required)
  SSH configuration for how to access the underlying control plane machines.
    
* `subnet_id` -
  (Required)
  The ARM ID of the subnet where the control plane VMs are deployed. Example: `/subscriptions//resourceGroups//providers/Microsoft.Network/virtualNetworks//subnets/default`.
    
* `tags` -
  (Optional)
  Optional. A set of tags to apply to all underlying control plane Azure resources.
    
* `version` -
  (Required)
  The Kubernetes version to run on control plane replicas (e.g. `1.19.10-gke.1000`). You can list all supported versions on a given Google Cloud region by calling GetAzureServerConfig.
    
* `vm_size` -
  (Optional)
  Optional. The Azure VM size name. Example: `Standard_DS2_v2`. For available VM sizes, see https://docs.microsoft.com/en-us/azure/virtual-machines/vm-naming-conventions. When unspecified, it defaults to `Standard_DS2_v2`.
    
The `ssh_config` block supports:
    
* `authorized_key` -
  (Required)
  The SSH public key data for VMs managed by Anthos. This accepts the authorized_keys file format used in OpenSSH according to the sshd(8) manual page.
    
The `fleet` block supports:
    
* `membership` -
  The name of the managed Hub Membership resource associated to this cluster. Membership names are formatted as projects/<project-number>/locations/global/membership/<cluster-id>.
    
* `project` -
  (Optional)
  The number of the Fleet host project where this cluster will be registered.
    
The `networking` block supports:
    
* `pod_address_cidr_blocks` -
  (Required)
  The IP address range of the pods in this cluster, in CIDR notation (e.g. `10.96.0.0/14`). All pods in the cluster get assigned a unique RFC1918 IPv4 address from these ranges. Only a single range is supported. This field cannot be changed after creation.
    
* `service_address_cidr_blocks` -
  (Required)
  The IP address range for services in this cluster, in CIDR notation (e.g. `10.96.0.0/14`). All services in the cluster get assigned a unique RFC1918 IPv4 address from these ranges. Only a single range is supported. This field cannot be changed after creating a cluster.
    
* `virtual_network_id` -
  (Required)
  The Azure Resource Manager (ARM) ID of the VNet associated with your cluster. All components in the cluster (i.e. control plane and node pools) run on a single VNet. Example: `/subscriptions/*/resourceGroups/*/providers/Microsoft.Network/virtualNetworks/*` This field cannot be changed after creation.
    
- - -

* `annotations` -
  (Optional)
  Optional. Annotations on the cluster. This field has the same restrictions as Kubernetes annotations. The total size of all keys and values combined is limited to 256k. Keys can have 2 segments: prefix (optional) and name (required), separated by a slash (/). Prefix must be a DNS subdomain. Name must be 63 characters or less, begin and end with alphanumerics, with dashes (-), underscores (_), dots (.), and alphanumerics between.
  
* `description` -
  (Optional)
  Optional. A human readable description of this cluster. Cannot be longer than 255 UTF-8 encoded bytes.
  
* `project` -
  (Optional)
  The project for the resource
  


The `database_encryption` block supports:
    
* `key_id` -
  (Required)
  The ARM ID of the Azure Key Vault key to encrypt / decrypt data. For example: `/subscriptions/<subscription-id>/resourceGroups/<resource-group-id>/providers/Microsoft.KeyVault/vaults/<key-vault-id>/keys/<key-name>` Encryption will always take the latest version of the key and hence specific version is not supported.
    
The `main_volume` block supports:
    
* `size_gib` -
  (Optional)
  Optional. The size of the disk, in GiBs. When unspecified, a default value is provided. See the specific reference in the parent resource.
    
The `proxy_config` block supports:
    
* `resource_group_id` -
  (Required)
  The ARM ID the of the resource group containing proxy keyvault. Resource group ids are formatted as `/subscriptions/<subscription-id>/resourceGroups/<resource-group-name>`
    
* `secret_id` -
  (Required)
  The URL the of the proxy setting secret with its version. Secret ids are formatted as `https:<key-vault-name>.vault.azure.net/secrets/<secret-name>/<secret-version>`.
    
The `replica_placements` block supports:
    
* `azure_availability_zone` -
  (Required)
  For a given replica, the Azure availability zone where to provision the control plane VM and the ETCD disk.
    
* `subnet_id` -
  (Required)
  For a given replica, the ARM ID of the subnet where the control plane VM is deployed. Make sure it's a subnet under the virtual network in the cluster configuration.
    
The `root_volume` block supports:
    
* `size_gib` -
  (Optional)
  Optional. The size of the disk, in GiBs. When unspecified, a default value is provided. See the specific reference in the parent resource.
    
## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{location}}/azureClusters/{{name}}`

* `create_time` -
  Output only. The time at which this cluster was created.
  
* `endpoint` -
  Output only. The endpoint of the cluster's API server.
  
* `etag` -
  Allows clients to perform consistent read-modify-writes through optimistic concurrency control. May be sent on update and delete requests to ensure the client has an up-to-date value before proceeding.
  
* `reconciling` -
  Output only. If set, there are currently changes in flight to the cluster.
  
* `state` -
  Output only. The current state of the cluster. Possible values: STATE_UNSPECIFIED, PROVISIONING, RUNNING, RECONCILING, STOPPING, ERROR, DEGRADED
  
* `uid` -
  Output only. A globally unique identifier for the cluster.
  
* `update_time` -
  Output only. The time at which this cluster was last updated.
  
* `workload_identity_config` -
  Output only. Workload Identity settings.
  
## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

Cluster can be imported using any of these accepted formats:

```
$ terraform import google_container_azure_cluster.default projects/{{project}}/locations/{{location}}/azureClusters/{{name}}
$ terraform import google_container_azure_cluster.default {{project}}/{{location}}/{{name}}
$ terraform import google_container_azure_cluster.default {{location}}/{{name}}
```



