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

func ResourceContainerAzureNodePool() *schema.Resource {
	return &schema.Resource{
		Create: resourceContainerAzureNodePoolCreate,
		Read:   resourceContainerAzureNodePoolRead,
		Update: resourceContainerAzureNodePoolUpdate,
		Delete: resourceContainerAzureNodePoolDelete,

		Importer: &schema.ResourceImporter{
			State: resourceContainerAzureNodePoolImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"autoscaling": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Autoscaler configuration for this node pool.",
				MaxItems:    1,
				Elem:        ContainerAzureNodePoolAutoscalingSchema(),
			},

			"cluster": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "The azureCluster for the resource",
			},

			"config": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The node configuration of the node pool.",
				MaxItems:    1,
				Elem:        ContainerAzureNodePoolConfigSchema(),
			},

			"location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The location for the resource",
			},

			"max_pods_constraint": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "The constraint on the maximum number of pods that can be run simultaneously on a node in the node pool.",
				MaxItems:    1,
				Elem:        ContainerAzureNodePoolMaxPodsConstraintSchema(),
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of this resource.",
			},

			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ARM ID of the subnet where the node pool VMs run. Make sure it's a subnet under the virtual network in the cluster configuration.",
			},

			"version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Kubernetes version (e.g. `1.19.10-gke.1000`) running on this node pool.",
			},

			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Optional. Annotations on the node pool. This field has the same restrictions as Kubernetes annotations. The total size of all keys and values combined is limited to 256k. Keys can have 2 segments: prefix (optional) and name (required), separated by a slash (/). Prefix must be a DNS subdomain. Name must be 63 characters or less, begin and end with alphanumerics, with dashes (-), underscores (_), dots (.), and alphanumerics between.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"azure_availability_zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The Azure availability zone of the nodes in this nodepool. When unspecified, it defaults to `1`.",
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
				Description: "Output only. The time at which this node pool was created.",
			},

			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Allows clients to perform consistent read-modify-writes through optimistic concurrency control. May be sent on update and delete requests to ensure the client has an up-to-date value before proceeding.",
			},

			"reconciling": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Output only. If set, there are currently pending changes to the node pool.",
			},

			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The current state of the node pool. Possible values: STATE_UNSPECIFIED, PROVISIONING, RUNNING, RECONCILING, STOPPING, ERROR, DEGRADED",
			},

			"uid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. A globally unique identifier for the node pool.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The time at which this node pool was last updated.",
			},
		},
	}
}

func ContainerAzureNodePoolAutoscalingSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"max_node_count": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Maximum number of nodes in the node pool. Must be >= min_node_count.",
			},

			"min_node_count": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Minimum number of nodes in the node pool. Must be >= 1 and <= max_node_count.",
			},
		},
	}
}

func ContainerAzureNodePoolConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ssh_config": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "SSH configuration for how to access the node pool machines.",
				MaxItems:    1,
				Elem:        ContainerAzureNodePoolConfigSshConfigSchema(),
			},

			"proxy_config": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Proxy configuration for outbound HTTP(S) traffic.",
				MaxItems:    1,
				Elem:        ContainerAzureNodePoolConfigProxyConfigSchema(),
			},

			"root_volume": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Configuration related to the root volume provisioned for each node pool machine. When unspecified, it defaults to a 32-GiB Azure Disk.",
				MaxItems:    1,
				Elem:        ContainerAzureNodePoolConfigRootVolumeSchema(),
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. A set of tags to apply to all underlying Azure resources for this node pool. This currently only includes Virtual Machine Scale Sets. Specify at most 50 pairs containing alphanumerics, spaces, and symbols (.+-=_:@/). Keys can be up to 127 Unicode characters. Values can be up to 255 Unicode characters.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"vm_size": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The Azure VM size name. Example: `Standard_DS2_v2`. See (/anthos/clusters/docs/azure/reference/supported-vms) for options. When unspecified, it defaults to `Standard_DS2_v2`.",
			},
		},
	}
}

func ContainerAzureNodePoolConfigSshConfigSchema() *schema.Resource {
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

func ContainerAzureNodePoolConfigProxyConfigSchema() *schema.Resource {
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

func ContainerAzureNodePoolConfigRootVolumeSchema() *schema.Resource {
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

func ContainerAzureNodePoolMaxPodsConstraintSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"max_pods_per_node": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The maximum number of pods to schedule on a single node.",
			},
		},
	}
}

func resourceContainerAzureNodePoolCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &containerazure.NodePool{
		Autoscaling:           expandContainerAzureNodePoolAutoscaling(d.Get("autoscaling")),
		Cluster:               dcl.String(d.Get("cluster").(string)),
		Config:                expandContainerAzureNodePoolConfig(d.Get("config")),
		Location:              dcl.String(d.Get("location").(string)),
		MaxPodsConstraint:     expandContainerAzureNodePoolMaxPodsConstraint(d.Get("max_pods_constraint")),
		Name:                  dcl.String(d.Get("name").(string)),
		SubnetId:              dcl.String(d.Get("subnet_id").(string)),
		Version:               dcl.String(d.Get("version").(string)),
		Annotations:           tpgresource.CheckStringMap(d.Get("annotations")),
		AzureAvailabilityZone: dcl.StringOrNil(d.Get("azure_availability_zone").(string)),
		Project:               dcl.String(project),
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
	res, err := client.ApplyNodePool(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating NodePool: %s", err)
	}

	log.Printf("[DEBUG] Finished creating NodePool %q: %#v", d.Id(), res)

	return resourceContainerAzureNodePoolRead(d, meta)
}

func resourceContainerAzureNodePoolRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &containerazure.NodePool{
		Autoscaling:           expandContainerAzureNodePoolAutoscaling(d.Get("autoscaling")),
		Cluster:               dcl.String(d.Get("cluster").(string)),
		Config:                expandContainerAzureNodePoolConfig(d.Get("config")),
		Location:              dcl.String(d.Get("location").(string)),
		MaxPodsConstraint:     expandContainerAzureNodePoolMaxPodsConstraint(d.Get("max_pods_constraint")),
		Name:                  dcl.String(d.Get("name").(string)),
		SubnetId:              dcl.String(d.Get("subnet_id").(string)),
		Version:               dcl.String(d.Get("version").(string)),
		Annotations:           tpgresource.CheckStringMap(d.Get("annotations")),
		AzureAvailabilityZone: dcl.StringOrNil(d.Get("azure_availability_zone").(string)),
		Project:               dcl.String(project),
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
	res, err := client.GetNodePool(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("ContainerAzureNodePool %q", d.Id())
		return tpgdclresource.HandleNotFoundDCLError(err, d, resourceName)
	}

	if err = d.Set("autoscaling", flattenContainerAzureNodePoolAutoscaling(res.Autoscaling)); err != nil {
		return fmt.Errorf("error setting autoscaling in state: %s", err)
	}
	if err = d.Set("cluster", res.Cluster); err != nil {
		return fmt.Errorf("error setting cluster in state: %s", err)
	}
	if err = d.Set("config", flattenContainerAzureNodePoolConfig(res.Config)); err != nil {
		return fmt.Errorf("error setting config in state: %s", err)
	}
	if err = d.Set("location", res.Location); err != nil {
		return fmt.Errorf("error setting location in state: %s", err)
	}
	if err = d.Set("max_pods_constraint", flattenContainerAzureNodePoolMaxPodsConstraint(res.MaxPodsConstraint)); err != nil {
		return fmt.Errorf("error setting max_pods_constraint in state: %s", err)
	}
	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}
	if err = d.Set("subnet_id", res.SubnetId); err != nil {
		return fmt.Errorf("error setting subnet_id in state: %s", err)
	}
	if err = d.Set("version", res.Version); err != nil {
		return fmt.Errorf("error setting version in state: %s", err)
	}
	if err = d.Set("annotations", res.Annotations); err != nil {
		return fmt.Errorf("error setting annotations in state: %s", err)
	}
	if err = d.Set("azure_availability_zone", res.AzureAvailabilityZone); err != nil {
		return fmt.Errorf("error setting azure_availability_zone in state: %s", err)
	}
	if err = d.Set("project", res.Project); err != nil {
		return fmt.Errorf("error setting project in state: %s", err)
	}
	if err = d.Set("create_time", res.CreateTime); err != nil {
		return fmt.Errorf("error setting create_time in state: %s", err)
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

	return nil
}
func resourceContainerAzureNodePoolUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &containerazure.NodePool{
		Autoscaling:           expandContainerAzureNodePoolAutoscaling(d.Get("autoscaling")),
		Cluster:               dcl.String(d.Get("cluster").(string)),
		Config:                expandContainerAzureNodePoolConfig(d.Get("config")),
		Location:              dcl.String(d.Get("location").(string)),
		MaxPodsConstraint:     expandContainerAzureNodePoolMaxPodsConstraint(d.Get("max_pods_constraint")),
		Name:                  dcl.String(d.Get("name").(string)),
		SubnetId:              dcl.String(d.Get("subnet_id").(string)),
		Version:               dcl.String(d.Get("version").(string)),
		Annotations:           tpgresource.CheckStringMap(d.Get("annotations")),
		AzureAvailabilityZone: dcl.StringOrNil(d.Get("azure_availability_zone").(string)),
		Project:               dcl.String(project),
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
	res, err := client.ApplyNodePool(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error updating NodePool: %s", err)
	}

	log.Printf("[DEBUG] Finished creating NodePool %q: %#v", d.Id(), res)

	return resourceContainerAzureNodePoolRead(d, meta)
}

func resourceContainerAzureNodePoolDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &containerazure.NodePool{
		Autoscaling:           expandContainerAzureNodePoolAutoscaling(d.Get("autoscaling")),
		Cluster:               dcl.String(d.Get("cluster").(string)),
		Config:                expandContainerAzureNodePoolConfig(d.Get("config")),
		Location:              dcl.String(d.Get("location").(string)),
		MaxPodsConstraint:     expandContainerAzureNodePoolMaxPodsConstraint(d.Get("max_pods_constraint")),
		Name:                  dcl.String(d.Get("name").(string)),
		SubnetId:              dcl.String(d.Get("subnet_id").(string)),
		Version:               dcl.String(d.Get("version").(string)),
		Annotations:           tpgresource.CheckStringMap(d.Get("annotations")),
		AzureAvailabilityZone: dcl.StringOrNil(d.Get("azure_availability_zone").(string)),
		Project:               dcl.String(project),
	}

	log.Printf("[DEBUG] Deleting NodePool %q", d.Id())
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
	if err := client.DeleteNodePool(context.Background(), obj); err != nil {
		return fmt.Errorf("Error deleting NodePool: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting NodePool %q", d.Id())
	return nil
}

func resourceContainerAzureNodePoolImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)

	if err := tpgresource.ParseImportId([]string{
		"projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/azureClusters/(?P<cluster>[^/]+)/azureNodePools/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<cluster>[^/]+)/(?P<name>[^/]+)",
		"(?P<location>[^/]+)/(?P<cluster>[^/]+)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVarsForId(d, config, "projects/{{project}}/locations/{{location}}/azureClusters/{{cluster}}/azureNodePools/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func expandContainerAzureNodePoolAutoscaling(o interface{}) *containerazure.NodePoolAutoscaling {
	if o == nil {
		return containerazure.EmptyNodePoolAutoscaling
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return containerazure.EmptyNodePoolAutoscaling
	}
	obj := objArr[0].(map[string]interface{})
	return &containerazure.NodePoolAutoscaling{
		MaxNodeCount: dcl.Int64(int64(obj["max_node_count"].(int))),
		MinNodeCount: dcl.Int64(int64(obj["min_node_count"].(int))),
	}
}

func flattenContainerAzureNodePoolAutoscaling(obj *containerazure.NodePoolAutoscaling) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"max_node_count": obj.MaxNodeCount,
		"min_node_count": obj.MinNodeCount,
	}

	return []interface{}{transformed}

}

func expandContainerAzureNodePoolConfig(o interface{}) *containerazure.NodePoolConfig {
	if o == nil {
		return containerazure.EmptyNodePoolConfig
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return containerazure.EmptyNodePoolConfig
	}
	obj := objArr[0].(map[string]interface{})
	return &containerazure.NodePoolConfig{
		SshConfig:   expandContainerAzureNodePoolConfigSshConfig(obj["ssh_config"]),
		ProxyConfig: expandContainerAzureNodePoolConfigProxyConfig(obj["proxy_config"]),
		RootVolume:  expandContainerAzureNodePoolConfigRootVolume(obj["root_volume"]),
		Tags:        tpgresource.CheckStringMap(obj["tags"]),
		VmSize:      dcl.StringOrNil(obj["vm_size"].(string)),
	}
}

func flattenContainerAzureNodePoolConfig(obj *containerazure.NodePoolConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"ssh_config":   flattenContainerAzureNodePoolConfigSshConfig(obj.SshConfig),
		"proxy_config": flattenContainerAzureNodePoolConfigProxyConfig(obj.ProxyConfig),
		"root_volume":  flattenContainerAzureNodePoolConfigRootVolume(obj.RootVolume),
		"tags":         obj.Tags,
		"vm_size":      obj.VmSize,
	}

	return []interface{}{transformed}

}

func expandContainerAzureNodePoolConfigSshConfig(o interface{}) *containerazure.NodePoolConfigSshConfig {
	if o == nil {
		return containerazure.EmptyNodePoolConfigSshConfig
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return containerazure.EmptyNodePoolConfigSshConfig
	}
	obj := objArr[0].(map[string]interface{})
	return &containerazure.NodePoolConfigSshConfig{
		AuthorizedKey: dcl.String(obj["authorized_key"].(string)),
	}
}

func flattenContainerAzureNodePoolConfigSshConfig(obj *containerazure.NodePoolConfigSshConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"authorized_key": obj.AuthorizedKey,
	}

	return []interface{}{transformed}

}

func expandContainerAzureNodePoolConfigProxyConfig(o interface{}) *containerazure.NodePoolConfigProxyConfig {
	if o == nil {
		return containerazure.EmptyNodePoolConfigProxyConfig
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return containerazure.EmptyNodePoolConfigProxyConfig
	}
	obj := objArr[0].(map[string]interface{})
	return &containerazure.NodePoolConfigProxyConfig{
		ResourceGroupId: dcl.String(obj["resource_group_id"].(string)),
		SecretId:        dcl.String(obj["secret_id"].(string)),
	}
}

func flattenContainerAzureNodePoolConfigProxyConfig(obj *containerazure.NodePoolConfigProxyConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"resource_group_id": obj.ResourceGroupId,
		"secret_id":         obj.SecretId,
	}

	return []interface{}{transformed}

}

func expandContainerAzureNodePoolConfigRootVolume(o interface{}) *containerazure.NodePoolConfigRootVolume {
	if o == nil {
		return nil
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return nil
	}
	obj := objArr[0].(map[string]interface{})
	return &containerazure.NodePoolConfigRootVolume{
		SizeGib: dcl.Int64OrNil(int64(obj["size_gib"].(int))),
	}
}

func flattenContainerAzureNodePoolConfigRootVolume(obj *containerazure.NodePoolConfigRootVolume) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"size_gib": obj.SizeGib,
	}

	return []interface{}{transformed}

}

func expandContainerAzureNodePoolMaxPodsConstraint(o interface{}) *containerazure.NodePoolMaxPodsConstraint {
	if o == nil {
		return containerazure.EmptyNodePoolMaxPodsConstraint
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return containerazure.EmptyNodePoolMaxPodsConstraint
	}
	obj := objArr[0].(map[string]interface{})
	return &containerazure.NodePoolMaxPodsConstraint{
		MaxPodsPerNode: dcl.Int64(int64(obj["max_pods_per_node"].(int))),
	}
}

func flattenContainerAzureNodePoolMaxPodsConstraint(obj *containerazure.NodePoolMaxPodsConstraint) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"max_pods_per_node": obj.MaxPodsPerNode,
	}

	return []interface{}{transformed}

}
