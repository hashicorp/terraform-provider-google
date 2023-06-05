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

package cloudbuild

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	cloudbuild "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/cloudbuild"

	"github.com/hashicorp/terraform-provider-google/google/tpgdclresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceCloudbuildWorkerPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudbuildWorkerPoolCreate,
		Read:   resourceCloudbuildWorkerPoolRead,
		Update: resourceCloudbuildWorkerPoolUpdate,
		Delete: resourceCloudbuildWorkerPoolDelete,

		Importer: &schema.ResourceImporter{
			State: resourceCloudbuildWorkerPoolImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
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
				Description: "User-defined name of the `WorkerPool`.",
			},

			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "User specified annotations. See https://google.aip.dev/128#annotations for more details such as format and size limitations.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A user-specified, human-readable name for the `WorkerPool`. If provided, this value must be 1-63 characters.",
			},

			"network_config": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Network configuration for the `WorkerPool`.",
				MaxItems:    1,
				Elem:        CloudbuildWorkerPoolNetworkConfigSchema(),
			},

			"project": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "The project for the resource",
			},

			"worker_config": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: "Configuration to be used for a creating workers in the `WorkerPool`.",
				MaxItems:    1,
				Elem:        CloudbuildWorkerPoolWorkerConfigSchema(),
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. Time at which the request to create the `WorkerPool` was received.",
			},

			"delete_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. Time at which the request to delete the `WorkerPool` was received.",
			},

			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. `WorkerPool` state. Possible values: STATE_UNSPECIFIED, PENDING, APPROVED, REJECTED, CANCELLED",
			},

			"uid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. A unique identifier for the `WorkerPool`.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. Time at which the request to update the `WorkerPool` was received.",
			},
		},
	}
}

func CloudbuildWorkerPoolNetworkConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"peered_network": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareResourceNames,
				Description:      "Required. Immutable. The network definition that the workers are peered to. If this section is left empty, the workers will be peered to `WorkerPool.project_id` on the service producer network. Must be in the format `projects/{project}/global/networks/{network}`, where `{project}` is a project number, such as `12345`, and `{network}` is the name of a VPC network in the project. See [Understanding network configuration options](https://cloud.google.com/cloud-build/docs/custom-workers/set-up-custom-worker-pool-environment#understanding_the_network_configuration_options)",
			},

			"peered_network_ip_range": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Immutable. Subnet IP range within the peered network. This is specified in CIDR notation with a slash and the subnet prefix size. You can optionally specify an IP address before the subnet prefix value. e.g. `192.168.0.0/29` would specify an IP range starting at 192.168.0.0 with a prefix size of 29 bits. `/16` would specify a prefix size of 16 bits, with an automatically determined IP within the peered VPC. If unspecified, a value of `/24` will be used.",
			},
		},
	}
}

func CloudbuildWorkerPoolWorkerConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"disk_size_gb": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Size of the disk attached to the worker, in GB. See [Worker pool config file](https://cloud.google.com/cloud-build/docs/custom-workers/worker-pool-config-file). Specify a value of up to 1000. If `0` is specified, Cloud Build will use a standard disk size.",
			},

			"machine_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Machine type of a worker, such as `n1-standard-1`. See [Worker pool config file](https://cloud.google.com/cloud-build/docs/custom-workers/worker-pool-config-file). If left blank, Cloud Build will use `n1-standard-1`.",
			},

			"no_external_ip": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "If true, workers are created without any public address, which prevents network egress to public IPs.",
			},
		},
	}
}

func resourceCloudbuildWorkerPoolCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &cloudbuild.WorkerPool{
		Location:      dcl.String(d.Get("location").(string)),
		Name:          dcl.String(d.Get("name").(string)),
		Annotations:   tpgresource.CheckStringMap(d.Get("annotations")),
		DisplayName:   dcl.String(d.Get("display_name").(string)),
		NetworkConfig: expandCloudbuildWorkerPoolNetworkConfig(d.Get("network_config")),
		Project:       dcl.String(project),
		WorkerConfig:  expandCloudbuildWorkerPoolWorkerConfig(d.Get("worker_config")),
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
	client := transport_tpg.NewDCLCloudbuildClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutCreate))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyWorkerPool(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating WorkerPool: %s", err)
	}

	log.Printf("[DEBUG] Finished creating WorkerPool %q: %#v", d.Id(), res)

	return resourceCloudbuildWorkerPoolRead(d, meta)
}

func resourceCloudbuildWorkerPoolRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &cloudbuild.WorkerPool{
		Location:      dcl.String(d.Get("location").(string)),
		Name:          dcl.String(d.Get("name").(string)),
		Annotations:   tpgresource.CheckStringMap(d.Get("annotations")),
		DisplayName:   dcl.String(d.Get("display_name").(string)),
		NetworkConfig: expandCloudbuildWorkerPoolNetworkConfig(d.Get("network_config")),
		Project:       dcl.String(project),
		WorkerConfig:  expandCloudbuildWorkerPoolWorkerConfig(d.Get("worker_config")),
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
	client := transport_tpg.NewDCLCloudbuildClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutRead))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.GetWorkerPool(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("CloudbuildWorkerPool %q", d.Id())
		return tpgdclresource.HandleNotFoundDCLError(err, d, resourceName)
	}

	if err = d.Set("location", res.Location); err != nil {
		return fmt.Errorf("error setting location in state: %s", err)
	}
	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}
	if err = d.Set("annotations", res.Annotations); err != nil {
		return fmt.Errorf("error setting annotations in state: %s", err)
	}
	if err = d.Set("display_name", res.DisplayName); err != nil {
		return fmt.Errorf("error setting display_name in state: %s", err)
	}
	if err = d.Set("network_config", flattenCloudbuildWorkerPoolNetworkConfig(res.NetworkConfig)); err != nil {
		return fmt.Errorf("error setting network_config in state: %s", err)
	}
	if err = d.Set("project", res.Project); err != nil {
		return fmt.Errorf("error setting project in state: %s", err)
	}
	if err = d.Set("worker_config", flattenCloudbuildWorkerPoolWorkerConfig(res.WorkerConfig)); err != nil {
		return fmt.Errorf("error setting worker_config in state: %s", err)
	}
	if err = d.Set("create_time", res.CreateTime); err != nil {
		return fmt.Errorf("error setting create_time in state: %s", err)
	}
	if err = d.Set("delete_time", res.DeleteTime); err != nil {
		return fmt.Errorf("error setting delete_time in state: %s", err)
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
func resourceCloudbuildWorkerPoolUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &cloudbuild.WorkerPool{
		Location:      dcl.String(d.Get("location").(string)),
		Name:          dcl.String(d.Get("name").(string)),
		Annotations:   tpgresource.CheckStringMap(d.Get("annotations")),
		DisplayName:   dcl.String(d.Get("display_name").(string)),
		NetworkConfig: expandCloudbuildWorkerPoolNetworkConfig(d.Get("network_config")),
		Project:       dcl.String(project),
		WorkerConfig:  expandCloudbuildWorkerPoolWorkerConfig(d.Get("worker_config")),
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
	client := transport_tpg.NewDCLCloudbuildClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutUpdate))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyWorkerPool(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error updating WorkerPool: %s", err)
	}

	log.Printf("[DEBUG] Finished creating WorkerPool %q: %#v", d.Id(), res)

	return resourceCloudbuildWorkerPoolRead(d, meta)
}

func resourceCloudbuildWorkerPoolDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &cloudbuild.WorkerPool{
		Location:      dcl.String(d.Get("location").(string)),
		Name:          dcl.String(d.Get("name").(string)),
		Annotations:   tpgresource.CheckStringMap(d.Get("annotations")),
		DisplayName:   dcl.String(d.Get("display_name").(string)),
		NetworkConfig: expandCloudbuildWorkerPoolNetworkConfig(d.Get("network_config")),
		Project:       dcl.String(project),
		WorkerConfig:  expandCloudbuildWorkerPoolWorkerConfig(d.Get("worker_config")),
	}

	log.Printf("[DEBUG] Deleting WorkerPool %q", d.Id())
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLCloudbuildClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutDelete))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	if err := client.DeleteWorkerPool(context.Background(), obj); err != nil {
		return fmt.Errorf("Error deleting WorkerPool: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting WorkerPool %q", d.Id())
	return nil
}

func resourceCloudbuildWorkerPoolImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)

	if err := tpgresource.ParseImportId([]string{
		"projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/workerPools/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<name>[^/]+)",
		"(?P<location>[^/]+)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVarsForId(d, config, "projects/{{project}}/locations/{{location}}/workerPools/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func expandCloudbuildWorkerPoolNetworkConfig(o interface{}) *cloudbuild.WorkerPoolNetworkConfig {
	if o == nil {
		return cloudbuild.EmptyWorkerPoolNetworkConfig
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return cloudbuild.EmptyWorkerPoolNetworkConfig
	}
	obj := objArr[0].(map[string]interface{})
	return &cloudbuild.WorkerPoolNetworkConfig{
		PeeredNetwork:        dcl.String(obj["peered_network"].(string)),
		PeeredNetworkIPRange: dcl.String(obj["peered_network_ip_range"].(string)),
	}
}

func flattenCloudbuildWorkerPoolNetworkConfig(obj *cloudbuild.WorkerPoolNetworkConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"peered_network":          obj.PeeredNetwork,
		"peered_network_ip_range": obj.PeeredNetworkIPRange,
	}

	return []interface{}{transformed}

}

func expandCloudbuildWorkerPoolWorkerConfig(o interface{}) *cloudbuild.WorkerPoolWorkerConfig {
	if o == nil {
		return nil
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return nil
	}
	obj := objArr[0].(map[string]interface{})
	return &cloudbuild.WorkerPoolWorkerConfig{
		DiskSizeGb:   dcl.Int64(int64(obj["disk_size_gb"].(int))),
		MachineType:  dcl.String(obj["machine_type"].(string)),
		NoExternalIP: dcl.Bool(obj["no_external_ip"].(bool)),
	}
}

func flattenCloudbuildWorkerPoolWorkerConfig(obj *cloudbuild.WorkerPoolWorkerConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"disk_size_gb":   obj.DiskSizeGb,
		"machine_type":   obj.MachineType,
		"no_external_ip": obj.NoExternalIP,
	}

	return []interface{}{transformed}

}
