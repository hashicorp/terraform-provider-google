// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: DCL     ***
//
// ----------------------------------------------------------------------------
//
//     This file is managed by Magic Modules (https://github.com/GoogleCloudPlatform/magic-modules)
//     and is based on the DCL (https://github.com/GoogleCloudPlatform/declarative-resource-client-library).
//     Changes will need to be made to the DCL or Magic Modules instead of here.
//
//     We are not currently able to accept contributions to this file. If changes
//     are required, please file an issue at https://github.com/hashicorp/terraform-provider-google/issues/new/choose
//
// ----------------------------------------------------------------------------

package containerazure

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	containerazure "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/containerazure"

	"github.com/hashicorp/terraform-provider-google/google/tpgdclresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceContainerAzureCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceContainerAzureClusterCreate,
		Read:   resourceContainerAzureClusterRead,
		Update: resourceContainerAzureClusterUpdate,
		Delete: resourceContainerAzureClusterDelete,

		Importer: &schema.ResourceImporter{
			State: resourceContainerAzureClusterImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"authorization": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Configuration related to the cluster RBAC settings.",
				MaxItems:    1,
				Elem:        ContainerAzureClusterAuthorizationSchema(),
			},

			"azure_region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Azure region where the cluster runs. Each Google Cloud region supports a subset of nearby Azure regions. You can call to list all supported Azure regions within a given Google Cloud region.",
			},

			"control_plane": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Configuration related to the cluster control plane.",
				MaxItems:    1,
				Elem:        ContainerAzureClusterControlPlaneSchema(),
			},

			"fleet": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "Fleet configuration.",
				MaxItems:    1,
				Elem:        ContainerAzureClusterFleetSchema(),
			},

			"location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The location for the resource",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of this resource.",
			},

			"networking": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "Cluster-wide networking configuration.",
				MaxItems:    1,
				Elem:        ContainerAzureClusterNetworkingSchema(),
			},

			"resource_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ARM ID of the resource group where the cluster resources are deployed. For example: `/subscriptions/*/resourceGroups/*`",
			},

			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Annotations on the cluster. This field has the same restrictions as Kubernetes annotations. The total size of all keys and values combined is limited to 256k. Keys can have 2 segments: prefix (optional) and name (required), separated by a slash (/). Prefix must be a DNS subdomain. Name must be 63 characters or less, begin and end with alphanumerics, with dashes (-), underscores (_), dots (.), and alphanumerics between.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"azure_services_authentication": {
				Type:          schema.TypeList,
				Optional:      true,
				Description:   "Azure authentication configuration for management of Azure resources",
				MaxItems:      1,
				Elem:          ContainerAzureClusterAzureServicesAuthenticationSchema(),
				ConflictsWith: []string{"client"},
			},

			"client": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "Name of the AzureClient. The `AzureClient` resource must reside on the same GCP project and region as the `AzureCluster`. `AzureClient` names are formatted as `projects/<project-number>/locations/<region>/azureClients/<client-id>`. See Resource Names (https:cloud.google.com/apis/design/resource_names) for more details on Google Cloud resource names.",
				ConflictsWith:    []string{"azure_services_authentication"},
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional. A human readable description of this cluster. Cannot be longer than 255 UTF-8 encoded bytes.",
			},

			"project": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "The project for the resource",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The time at which this cluster was created.",
			},

			"endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The endpoint of the cluster's API server.",
			},

			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Allows clients to perform consistent read-modify-writes through optimistic concurrency control. May be sent on update and delete requests to ensure the client has an up-to-date value before proceeding.",
			},

			"reconciling": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Output only. If set, there are currently changes in flight to the cluster.",
			},

			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The current state of the cluster. Possible values: STATE_UNSPECIFIED, PROVISIONING, RUNNING, RECONCILING, STOPPING, ERROR, DEGRADED",
			},

			"uid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. A globally unique identifier for the cluster.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The time at which this cluster was last updated.",
			},

			"workload_identity_config": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Output only. Workload Identity settings.",
				Elem:        ContainerAzureClusterWorkloadIdentityConfigSchema(),
			},
		},
	}
}

func ContainerAzureClusterAuthorizationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"admin_users": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Users that can perform operations as a cluster admin. A new ClusterRoleBinding will be created to grant the cluster-admin ClusterRole to the users. Up to ten admin users can be provided. For more info on RBAC, see https://kubernetes.io/docs/reference/access-authn-authz/rbac/#user-facing-roles",
				Elem:        ContainerAzureClusterAuthorizationAdminUsersSchema(),
			},
		},
	}
}

func ContainerAzureClusterAuthorizationAdminUsersSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the user, e.g. `my-gcp-id@gmail.com`.",
			},
		},
	}
}

func ContainerAzureClusterControlPlaneSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ssh_config": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "SSH configuration for how to access the underlying control plane machines.",
				MaxItems:    1,
				Elem:        ContainerAzureClusterControlPlaneSshConfigSchema(),
			},

			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ARM ID of the subnet where the control plane VMs are deployed. Example: `/subscriptions//resourceGroups//providers/Microsoft.Network/virtualNetworks//subnets/default`.",
			},

			"version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Kubernetes version to run on control plane replicas (e.g. `1.19.10-gke.1000`). You can list all supported versions on a given Google Cloud region by calling GetAzureServerConfig.",
			},

			"database_encryption": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Configuration related to application-layer secrets encryption.",
				MaxItems:    1,
				Elem:        ContainerAzureClusterControlPlaneDatabaseEncryptionSchema(),
			},

			"main_volume": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Configuration related to the main volume provisioned for each control plane replica. The main volume is in charge of storing all of the cluster's etcd state. When unspecified, it defaults to a 8-GiB Azure Disk.",
				MaxItems:    1,
				Elem:        ContainerAzureClusterControlPlaneMainVolumeSchema(),
			},

			"proxy_config": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Proxy configuration for outbound HTTP(S) traffic.",
				MaxItems:    1,
				Elem:        ContainerAzureClusterControlPlaneProxyConfigSchema(),
			},

			"replica_placements": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Configuration for where to place the control plane replicas. Up to three replica placement instances can be specified. If replica_placements is set, the replica placement instances will be applied to the three control plane replicas as evenly as possible.",
				Elem:        ContainerAzureClusterControlPlaneReplicaPlacementsSchema(),
			},

			"root_volume": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Configuration related to the root volume provisioned for each control plane replica. When unspecified, it defaults to 32-GiB Azure Disk.",
				MaxItems:    1,
				Elem:        ContainerAzureClusterControlPlaneRootVolumeSchema(),
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. A set of tags to apply to all underlying control plane Azure resources.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"vm_size": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Optional. The Azure VM size name. Example: `Standard_DS2_v2`. For available VM sizes, see https://docs.microsoft.com/en-us/azure/virtual-machines/vm-naming-conventions. When unspecified, it defaults to `Standard_DS2_v2`.",
			},
		},
	}
}

func ContainerAzureClusterControlPlaneSshConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"authorized_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The SSH public key data for VMs managed by Anthos. This accepts the authorized_keys file format used in OpenSSH according to the sshd(8) manual page.",
			},
		},
	}
}

func ContainerAzureClusterControlPlaneDatabaseEncryptionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ARM ID of the Azure Key Vault key to encrypt / decrypt data. For example: `/subscriptions/<subscription-id>/resourceGroups/<resource-group-id>/providers/Microsoft.KeyVault/vaults/<key-vault-id>/keys/<key-name>` Encryption will always take the latest version of the key and hence specific version is not supported.",
			},
		},
	}
}

func ContainerAzureClusterControlPlaneMainVolumeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"size_gib": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The size of the disk, in GiBs. When unspecified, a default value is provided. See the specific reference in the parent resource.",
			},
		},
	}
}

func ContainerAzureClusterControlPlaneProxyConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ARM ID the of the resource group containing proxy keyvault. Resource group ids are formatted as `/subscriptions/<subscription-id>/resourceGroups/<resource-group-name>`",
			},

			"secret_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The URL the of the proxy setting secret with its version. Secret ids are formatted as `https:<key-vault-name>.vault.azure.net/secrets/<secret-name>/<secret-version>`.",
			},
		},
	}
}

func ContainerAzureClusterControlPlaneReplicaPlacementsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"azure_availability_zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "For a given replica, the Azure availability zone where to provision the control plane VM and the ETCD disk.",
			},

			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "For a given replica, the ARM ID of the subnet where the control plane VM is deployed. Make sure it's a subnet under the virtual network in the cluster configuration.",
			},
		},
	}
}

func ContainerAzureClusterControlPlaneRootVolumeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"size_gib": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The size of the disk, in GiBs. When unspecified, a default value is provided. See the specific reference in the parent resource.",
			},
		},
	}
}

func ContainerAzureClusterFleetSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "The number of the Fleet host project where this cluster will be registered.",
			},

			"membership": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the managed Hub Membership resource associated to this cluster. Membership names are formatted as projects/<project-number>/locations/global/membership/<cluster-id>.",
			},
		},
	}
}

func ContainerAzureClusterNetworkingSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"pod_address_cidr_blocks": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "The IP address range of the pods in this cluster, in CIDR notation (e.g. `10.96.0.0/14`). All pods in the cluster get assigned a unique RFC1918 IPv4 address from these ranges. Only a single range is supported. This field cannot be changed after creation.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"service_address_cidr_blocks": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "The IP address range for services in this cluster, in CIDR notation (e.g. `10.96.0.0/14`). All services in the cluster get assigned a unique RFC1918 IPv4 address from these ranges. Only a single range is supported. This field cannot be changed after creating a cluster.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"virtual_network_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Azure Resource Manager (ARM) ID of the VNet associated with your cluster. All components in the cluster (i.e. control plane and node pools) run on a single VNet. Example: `/subscriptions/*/resourceGroups/*/providers/Microsoft.Network/virtualNetworks/*` This field cannot be changed after creation.",
			},
		},
	}
}

func ContainerAzureClusterAzureServicesAuthenticationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"application_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Azure Active Directory Application ID for Authentication configuration.",
			},

			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Azure Active Directory Tenant ID for Authentication configuration.",
			},
		},
	}
}

func ContainerAzureClusterWorkloadIdentityConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"identity_provider": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the OIDC Identity Provider (IdP) associated to the Workload Identity Pool.",
			},

			"issuer_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The OIDC issuer URL for this cluster.",
			},

			"workload_pool": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Workload Identity Pool associated to the cluster.",
			},
		},
	}
}

func resourceContainerAzureClusterCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &containerazure.Cluster{
		Authorization:               expandContainerAzureClusterAuthorization(d.Get("authorization")),
		AzureRegion:                 dcl.String(d.Get("azure_region").(string)),
		ControlPlane:                expandContainerAzureClusterControlPlane(d.Get("control_plane")),
		Fleet:                       expandContainerAzureClusterFleet(d.Get("fleet")),
		Location:                    dcl.String(d.Get("location").(string)),
		Name:                        dcl.String(d.Get("name").(string)),
		Networking:                  expandContainerAzureClusterNetworking(d.Get("networking")),
		ResourceGroupId:             dcl.String(d.Get("resource_group_id").(string)),
		Annotations:                 tpgresource.CheckStringMap(d.Get("annotations")),
		AzureServicesAuthentication: expandContainerAzureClusterAzureServicesAuthentication(d.Get("azure_services_authentication")),
		Client:                      dcl.String(d.Get("client").(string)),
		Description:                 dcl.String(d.Get("description").(string)),
		Project:                     dcl.String(project),
	}

	id, err := obj.ID()
	if err != nil {
		return fmt.Errorf("error constructing id: %s", err)
	}
	d.SetId(id)
	directive := tpgdclresource.CreateDirective
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLContainerAzureClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutCreate))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyCluster(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating Cluster: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Cluster %q: %#v", d.Id(), res)

	return resourceContainerAzureClusterRead(d, meta)
}

func resourceContainerAzureClusterRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &containerazure.Cluster{
		Authorization:               expandContainerAzureClusterAuthorization(d.Get("authorization")),
		AzureRegion:                 dcl.String(d.Get("azure_region").(string)),
		ControlPlane:                expandContainerAzureClusterControlPlane(d.Get("control_plane")),
		Fleet:                       expandContainerAzureClusterFleet(d.Get("fleet")),
		Location:                    dcl.String(d.Get("location").(string)),
		Name:                        dcl.String(d.Get("name").(string)),
		Networking:                  expandContainerAzureClusterNetworking(d.Get("networking")),
		ResourceGroupId:             dcl.String(d.Get("resource_group_id").(string)),
		Annotations:                 tpgresource.CheckStringMap(d.Get("annotations")),
		AzureServicesAuthentication: expandContainerAzureClusterAzureServicesAuthentication(d.Get("azure_services_authentication")),
		Client:                      dcl.String(d.Get("client").(string)),
		Description:                 dcl.String(d.Get("description").(string)),
		Project:                     dcl.String(project),
	}

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLContainerAzureClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutRead))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.GetCluster(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("ContainerAzureCluster %q", d.Id())
		return tpgdclresource.HandleNotFoundDCLError(err, d, resourceName)
	}

	if err = d.Set("authorization", flattenContainerAzureClusterAuthorization(res.Authorization)); err != nil {
		return fmt.Errorf("error setting authorization in state: %s", err)
	}
	if err = d.Set("azure_region", res.AzureRegion); err != nil {
		return fmt.Errorf("error setting azure_region in state: %s", err)
	}
	if err = d.Set("control_plane", flattenContainerAzureClusterControlPlane(res.ControlPlane)); err != nil {
		return fmt.Errorf("error setting control_plane in state: %s", err)
	}
	if err = d.Set("fleet", flattenContainerAzureClusterFleet(res.Fleet)); err != nil {
		return fmt.Errorf("error setting fleet in state: %s", err)
	}
	if err = d.Set("location", res.Location); err != nil {
		return fmt.Errorf("error setting location in state: %s", err)
	}
	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}
	if err = d.Set("networking", flattenContainerAzureClusterNetworking(res.Networking)); err != nil {
		return fmt.Errorf("error setting networking in state: %s", err)
	}
	if err = d.Set("resource_group_id", res.ResourceGroupId); err != nil {
		return fmt.Errorf("error setting resource_group_id in state: %s", err)
	}
	if err = d.Set("annotations", res.Annotations); err != nil {
		return fmt.Errorf("error setting annotations in state: %s", err)
	}
	if err = d.Set("azure_services_authentication", flattenContainerAzureClusterAzureServicesAuthentication(res.AzureServicesAuthentication)); err != nil {
		return fmt.Errorf("error setting azure_services_authentication in state: %s", err)
	}
	if err = d.Set("client", res.Client); err != nil {
		return fmt.Errorf("error setting client in state: %s", err)
	}
	if err = d.Set("description", res.Description); err != nil {
		return fmt.Errorf("error setting description in state: %s", err)
	}
	if err = d.Set("project", res.Project); err != nil {
		return fmt.Errorf("error setting project in state: %s", err)
	}
	if err = d.Set("create_time", res.CreateTime); err != nil {
		return fmt.Errorf("error setting create_time in state: %s", err)
	}
	if err = d.Set("endpoint", res.Endpoint); err != nil {
		return fmt.Errorf("error setting endpoint in state: %s", err)
	}
	if err = d.Set("etag", res.Etag); err != nil {
		return fmt.Errorf("error setting etag in state: %s", err)
	}
	if err = d.Set("reconciling", res.Reconciling); err != nil {
		return fmt.Errorf("error setting reconciling in state: %s", err)
	}
	if err = d.Set("state", res.State); err != nil {
		return fmt.Errorf("error setting state in state: %s", err)
	}
	if err = d.Set("uid", res.Uid); err != nil {
		return fmt.Errorf("error setting uid in state: %s", err)
	}
	if err = d.Set("update_time", res.UpdateTime); err != nil {
		return fmt.Errorf("error setting update_time in state: %s", err)
	}
	if err = d.Set("workload_identity_config", flattenContainerAzureClusterWorkloadIdentityConfig(res.WorkloadIdentityConfig)); err != nil {
		return fmt.Errorf("error setting workload_identity_config in state: %s", err)
	}

	return nil
}
func resourceContainerAzureClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &containerazure.Cluster{
		Authorization:               expandContainerAzureClusterAuthorization(d.Get("authorization")),
		AzureRegion:                 dcl.String(d.Get("azure_region").(string)),
		ControlPlane:                expandContainerAzureClusterControlPlane(d.Get("control_plane")),
		Fleet:                       expandContainerAzureClusterFleet(d.Get("fleet")),
		Location:                    dcl.String(d.Get("location").(string)),
		Name:                        dcl.String(d.Get("name").(string)),
		Networking:                  expandContainerAzureClusterNetworking(d.Get("networking")),
		ResourceGroupId:             dcl.String(d.Get("resource_group_id").(string)),
		Annotations:                 tpgresource.CheckStringMap(d.Get("annotations")),
		AzureServicesAuthentication: expandContainerAzureClusterAzureServicesAuthentication(d.Get("azure_services_authentication")),
		Client:                      dcl.String(d.Get("client").(string)),
		Description:                 dcl.String(d.Get("description").(string)),
		Project:                     dcl.String(project),
	}
	directive := tpgdclresource.UpdateDirective
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLContainerAzureClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutUpdate))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyCluster(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error updating Cluster: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Cluster %q: %#v", d.Id(), res)

	return resourceContainerAzureClusterRead(d, meta)
}

func resourceContainerAzureClusterDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &containerazure.Cluster{
		Authorization:               expandContainerAzureClusterAuthorization(d.Get("authorization")),
		AzureRegion:                 dcl.String(d.Get("azure_region").(string)),
		ControlPlane:                expandContainerAzureClusterControlPlane(d.Get("control_plane")),
		Fleet:                       expandContainerAzureClusterFleet(d.Get("fleet")),
		Location:                    dcl.String(d.Get("location").(string)),
		Name:                        dcl.String(d.Get("name").(string)),
		Networking:                  expandContainerAzureClusterNetworking(d.Get("networking")),
		ResourceGroupId:             dcl.String(d.Get("resource_group_id").(string)),
		Annotations:                 tpgresource.CheckStringMap(d.Get("annotations")),
		AzureServicesAuthentication: expandContainerAzureClusterAzureServicesAuthentication(d.Get("azure_services_authentication")),
		Client:                      dcl.String(d.Get("client").(string)),
		Description:                 dcl.String(d.Get("description").(string)),
		Project:                     dcl.String(project),
	}

	log.Printf("[DEBUG] Deleting Cluster %q", d.Id())
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLContainerAzureClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutDelete))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	if err := client.DeleteCluster(context.Background(), obj); err != nil {
		return fmt.Errorf("Error deleting Cluster: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting Cluster %q", d.Id())
	return nil
}

func resourceContainerAzureClusterImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)

	if err := tpgresource.ParseImportId([]string{
		"projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/azureClusters/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<name>[^/]+)",
		"(?P<location>[^/]+)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVarsForId(d, config, "projects/{{project}}/locations/{{location}}/azureClusters/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func expandContainerAzureClusterAuthorization(o interface{}) *containerazure.ClusterAuthorization {
	if o == nil {
		return containerazure.EmptyClusterAuthorization
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return containerazure.EmptyClusterAuthorization
	}
	obj := objArr[0].(map[string]interface{})
	return &containerazure.ClusterAuthorization{
		AdminUsers: expandContainerAzureClusterAuthorizationAdminUsersArray(obj["admin_users"]),
	}
}

func flattenContainerAzureClusterAuthorization(obj *containerazure.ClusterAuthorization) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"admin_users": flattenContainerAzureClusterAuthorizationAdminUsersArray(obj.AdminUsers),
	}

	return []interface{}{transformed}

}
func expandContainerAzureClusterAuthorizationAdminUsersArray(o interface{}) []containerazure.ClusterAuthorizationAdminUsers {
	if o == nil {
		return make([]containerazure.ClusterAuthorizationAdminUsers, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make([]containerazure.ClusterAuthorizationAdminUsers, 0)
	}

	items := make([]containerazure.ClusterAuthorizationAdminUsers, 0, len(objs))
	for _, item := range objs {
		i := expandContainerAzureClusterAuthorizationAdminUsers(item)
		items = append(items, *i)
	}

	return items
}

func expandContainerAzureClusterAuthorizationAdminUsers(o interface{}) *containerazure.ClusterAuthorizationAdminUsers {
	if o == nil {
		return containerazure.EmptyClusterAuthorizationAdminUsers
	}

	obj := o.(map[string]interface{})
	return &containerazure.ClusterAuthorizationAdminUsers{
		Username: dcl.String(obj["username"].(string)),
	}
}

func flattenContainerAzureClusterAuthorizationAdminUsersArray(objs []containerazure.ClusterAuthorizationAdminUsers) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenContainerAzureClusterAuthorizationAdminUsers(&item)
		items = append(items, i)
	}

	return items
}

func flattenContainerAzureClusterAuthorizationAdminUsers(obj *containerazure.ClusterAuthorizationAdminUsers) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"username": obj.Username,
	}

	return transformed

}

func expandContainerAzureClusterControlPlane(o interface{}) *containerazure.ClusterControlPlane {
	if o == nil {
		return containerazure.EmptyClusterControlPlane
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return containerazure.EmptyClusterControlPlane
	}
	obj := objArr[0].(map[string]interface{})
	return &containerazure.ClusterControlPlane{
		SshConfig:          expandContainerAzureClusterControlPlaneSshConfig(obj["ssh_config"]),
		SubnetId:           dcl.String(obj["subnet_id"].(string)),
		Version:            dcl.String(obj["version"].(string)),
		DatabaseEncryption: expandContainerAzureClusterControlPlaneDatabaseEncryption(obj["database_encryption"]),
		MainVolume:         expandContainerAzureClusterControlPlaneMainVolume(obj["main_volume"]),
		ProxyConfig:        expandContainerAzureClusterControlPlaneProxyConfig(obj["proxy_config"]),
		ReplicaPlacements:  expandContainerAzureClusterControlPlaneReplicaPlacementsArray(obj["replica_placements"]),
		RootVolume:         expandContainerAzureClusterControlPlaneRootVolume(obj["root_volume"]),
		Tags:               tpgresource.CheckStringMap(obj["tags"]),
		VmSize:             dcl.StringOrNil(obj["vm_size"].(string)),
	}
}

func flattenContainerAzureClusterControlPlane(obj *containerazure.ClusterControlPlane) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"ssh_config":          flattenContainerAzureClusterControlPlaneSshConfig(obj.SshConfig),
		"subnet_id":           obj.SubnetId,
		"version":             obj.Version,
		"database_encryption": flattenContainerAzureClusterControlPlaneDatabaseEncryption(obj.DatabaseEncryption),
		"main_volume":         flattenContainerAzureClusterControlPlaneMainVolume(obj.MainVolume),
		"proxy_config":        flattenContainerAzureClusterControlPlaneProxyConfig(obj.ProxyConfig),
		"replica_placements":  flattenContainerAzureClusterControlPlaneReplicaPlacementsArray(obj.ReplicaPlacements),
		"root_volume":         flattenContainerAzureClusterControlPlaneRootVolume(obj.RootVolume),
		"tags":                obj.Tags,
		"vm_size":             obj.VmSize,
	}

	return []interface{}{transformed}

}

func expandContainerAzureClusterControlPlaneSshConfig(o interface{}) *containerazure.ClusterControlPlaneSshConfig {
	if o == nil {
		return containerazure.EmptyClusterControlPlaneSshConfig
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return containerazure.EmptyClusterControlPlaneSshConfig
	}
	obj := objArr[0].(map[string]interface{})
	return &containerazure.ClusterControlPlaneSshConfig{
		AuthorizedKey: dcl.String(obj["authorized_key"].(string)),
	}
}

func flattenContainerAzureClusterControlPlaneSshConfig(obj *containerazure.ClusterControlPlaneSshConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"authorized_key": obj.AuthorizedKey,
	}

	return []interface{}{transformed}

}

func expandContainerAzureClusterControlPlaneDatabaseEncryption(o interface{}) *containerazure.ClusterControlPlaneDatabaseEncryption {
	if o == nil {
		return containerazure.EmptyClusterControlPlaneDatabaseEncryption
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return containerazure.EmptyClusterControlPlaneDatabaseEncryption
	}
	obj := objArr[0].(map[string]interface{})
	return &containerazure.ClusterControlPlaneDatabaseEncryption{
		KeyId: dcl.String(obj["key_id"].(string)),
	}
}

func flattenContainerAzureClusterControlPlaneDatabaseEncryption(obj *containerazure.ClusterControlPlaneDatabaseEncryption) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"key_id": obj.KeyId,
	}

	return []interface{}{transformed}

}

func expandContainerAzureClusterControlPlaneMainVolume(o interface{}) *containerazure.ClusterControlPlaneMainVolume {
	if o == nil {
		return nil
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return nil
	}
	obj := objArr[0].(map[string]interface{})
	return &containerazure.ClusterControlPlaneMainVolume{
		SizeGib: dcl.Int64OrNil(int64(obj["size_gib"].(int))),
	}
}

func flattenContainerAzureClusterControlPlaneMainVolume(obj *containerazure.ClusterControlPlaneMainVolume) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"size_gib": obj.SizeGib,
	}

	return []interface{}{transformed}

}

func expandContainerAzureClusterControlPlaneProxyConfig(o interface{}) *containerazure.ClusterControlPlaneProxyConfig {
	if o == nil {
		return containerazure.EmptyClusterControlPlaneProxyConfig
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return containerazure.EmptyClusterControlPlaneProxyConfig
	}
	obj := objArr[0].(map[string]interface{})
	return &containerazure.ClusterControlPlaneProxyConfig{
		ResourceGroupId: dcl.String(obj["resource_group_id"].(string)),
		SecretId:        dcl.String(obj["secret_id"].(string)),
	}
}

func flattenContainerAzureClusterControlPlaneProxyConfig(obj *containerazure.ClusterControlPlaneProxyConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"resource_group_id": obj.ResourceGroupId,
		"secret_id":         obj.SecretId,
	}

	return []interface{}{transformed}

}
func expandContainerAzureClusterControlPlaneReplicaPlacementsArray(o interface{}) []containerazure.ClusterControlPlaneReplicaPlacements {
	if o == nil {
		return make([]containerazure.ClusterControlPlaneReplicaPlacements, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make([]containerazure.ClusterControlPlaneReplicaPlacements, 0)
	}

	items := make([]containerazure.ClusterControlPlaneReplicaPlacements, 0, len(objs))
	for _, item := range objs {
		i := expandContainerAzureClusterControlPlaneReplicaPlacements(item)
		items = append(items, *i)
	}

	return items
}

func expandContainerAzureClusterControlPlaneReplicaPlacements(o interface{}) *containerazure.ClusterControlPlaneReplicaPlacements {
	if o == nil {
		return containerazure.EmptyClusterControlPlaneReplicaPlacements
	}

	obj := o.(map[string]interface{})
	return &containerazure.ClusterControlPlaneReplicaPlacements{
		AzureAvailabilityZone: dcl.String(obj["azure_availability_zone"].(string)),
		SubnetId:              dcl.String(obj["subnet_id"].(string)),
	}
}

func flattenContainerAzureClusterControlPlaneReplicaPlacementsArray(objs []containerazure.ClusterControlPlaneReplicaPlacements) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenContainerAzureClusterControlPlaneReplicaPlacements(&item)
		items = append(items, i)
	}

	return items
}

func flattenContainerAzureClusterControlPlaneReplicaPlacements(obj *containerazure.ClusterControlPlaneReplicaPlacements) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"azure_availability_zone": obj.AzureAvailabilityZone,
		"subnet_id":               obj.SubnetId,
	}

	return transformed

}

func expandContainerAzureClusterControlPlaneRootVolume(o interface{}) *containerazure.ClusterControlPlaneRootVolume {
	if o == nil {
		return nil
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return nil
	}
	obj := objArr[0].(map[string]interface{})
	return &containerazure.ClusterControlPlaneRootVolume{
		SizeGib: dcl.Int64OrNil(int64(obj["size_gib"].(int))),
	}
}

func flattenContainerAzureClusterControlPlaneRootVolume(obj *containerazure.ClusterControlPlaneRootVolume) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"size_gib": obj.SizeGib,
	}

	return []interface{}{transformed}

}

func expandContainerAzureClusterFleet(o interface{}) *containerazure.ClusterFleet {
	if o == nil {
		return containerazure.EmptyClusterFleet
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return containerazure.EmptyClusterFleet
	}
	obj := objArr[0].(map[string]interface{})
	return &containerazure.ClusterFleet{
		Project: dcl.StringOrNil(obj["project"].(string)),
	}
}

func flattenContainerAzureClusterFleet(obj *containerazure.ClusterFleet) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"project":    obj.Project,
		"membership": obj.Membership,
	}

	return []interface{}{transformed}

}

func expandContainerAzureClusterNetworking(o interface{}) *containerazure.ClusterNetworking {
	if o == nil {
		return containerazure.EmptyClusterNetworking
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return containerazure.EmptyClusterNetworking
	}
	obj := objArr[0].(map[string]interface{})
	return &containerazure.ClusterNetworking{
		PodAddressCidrBlocks:     tpgdclresource.ExpandStringArray(obj["pod_address_cidr_blocks"]),
		ServiceAddressCidrBlocks: tpgdclresource.ExpandStringArray(obj["service_address_cidr_blocks"]),
		VirtualNetworkId:         dcl.String(obj["virtual_network_id"].(string)),
	}
}

func flattenContainerAzureClusterNetworking(obj *containerazure.ClusterNetworking) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"pod_address_cidr_blocks":     obj.PodAddressCidrBlocks,
		"service_address_cidr_blocks": obj.ServiceAddressCidrBlocks,
		"virtual_network_id":          obj.VirtualNetworkId,
	}

	return []interface{}{transformed}

}

func expandContainerAzureClusterAzureServicesAuthentication(o interface{}) *containerazure.ClusterAzureServicesAuthentication {
	if o == nil {
		return containerazure.EmptyClusterAzureServicesAuthentication
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return containerazure.EmptyClusterAzureServicesAuthentication
	}
	obj := objArr[0].(map[string]interface{})
	return &containerazure.ClusterAzureServicesAuthentication{
		ApplicationId: dcl.String(obj["application_id"].(string)),
		TenantId:      dcl.String(obj["tenant_id"].(string)),
	}
}

func flattenContainerAzureClusterAzureServicesAuthentication(obj *containerazure.ClusterAzureServicesAuthentication) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"application_id": obj.ApplicationId,
		"tenant_id":      obj.TenantId,
	}

	return []interface{}{transformed}

}

func flattenContainerAzureClusterWorkloadIdentityConfig(obj *containerazure.ClusterWorkloadIdentityConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"identity_provider": obj.IdentityProvider,
		"issuer_uri":        obj.IssuerUri,
		"workload_pool":     obj.WorkloadPool,
	}

	return []interface{}{transformed}

}
