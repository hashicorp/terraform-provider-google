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

package networkconnectivity

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	networkconnectivity "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/networkconnectivity"

	"github.com/hashicorp/terraform-provider-google/google/tpgdclresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceNetworkConnectivityHub() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkConnectivityHubCreate,
		Read:   resourceNetworkConnectivityHubRead,
		Update: resourceNetworkConnectivityHubUpdate,
		Delete: resourceNetworkConnectivityHubDelete,

		Importer: &schema.ResourceImporter{
			State: resourceNetworkConnectivityHubImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Immutable. The name of the hub. Hub names must be unique. They use the following form: `projects/{project_number}/locations/global/hubs/{hub_id}`",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An optional description of the hub.",
			},

			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Optional labels in key:value format. For more information about labels, see [Requirements for labels](https://cloud.google.com/resource-manager/docs/creating-managing-labels#requirements).",
				Elem:        &schema.Schema{Type: schema.TypeString},
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
				Description: "Output only. The time the hub was created.",
			},

			"routing_vpcs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The VPC network associated with this hub's spokes. All of the VPN tunnels, VLAN attachments, and router appliance instances referenced by this hub's spokes must belong to this VPC network. This field is read-only. Network Connectivity Center automatically populates it based on the set of spokes attached to the hub.",
				Elem:        NetworkConnectivityHubRoutingVpcsSchema(),
			},

			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The current lifecycle state of this hub. Possible values: STATE_UNSPECIFIED, CREATING, ACTIVE, DELETING",
			},

			"unique_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The Google-generated UUID for the hub. This value is unique across all hub resources. If a hub is deleted and another with the same name is created, the new hub is assigned a different unique_id.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The time the hub was last updated.",
			},
		},
	}
}

func NetworkConnectivityHubRoutingVpcsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URI of the VPC network.",
			},
		},
	}
}

func resourceNetworkConnectivityHubCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &networkconnectivity.Hub{
		Name:        dcl.String(d.Get("name").(string)),
		Description: dcl.String(d.Get("description").(string)),
		Labels:      tpgresource.CheckStringMap(d.Get("labels")),
		Project:     dcl.String(project),
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
	client := transport_tpg.NewDCLNetworkConnectivityClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutCreate))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyHub(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating Hub: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Hub %q: %#v", d.Id(), res)

	return resourceNetworkConnectivityHubRead(d, meta)
}

func resourceNetworkConnectivityHubRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &networkconnectivity.Hub{
		Name:        dcl.String(d.Get("name").(string)),
		Description: dcl.String(d.Get("description").(string)),
		Labels:      tpgresource.CheckStringMap(d.Get("labels")),
		Project:     dcl.String(project),
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
	client := transport_tpg.NewDCLNetworkConnectivityClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutRead))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.GetHub(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("NetworkConnectivityHub %q", d.Id())
		return tpgdclresource.HandleNotFoundDCLError(err, d, resourceName)
	}

	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}
	if err = d.Set("description", res.Description); err != nil {
		return fmt.Errorf("error setting description in state: %s", err)
	}
	if err = d.Set("labels", res.Labels); err != nil {
		return fmt.Errorf("error setting labels in state: %s", err)
	}
	if err = d.Set("project", res.Project); err != nil {
		return fmt.Errorf("error setting project in state: %s", err)
	}
	if err = d.Set("create_time", res.CreateTime); err != nil {
		return fmt.Errorf("error setting create_time in state: %s", err)
	}
	if err = d.Set("routing_vpcs", flattenNetworkConnectivityHubRoutingVpcsArray(res.RoutingVpcs)); err != nil {
		return fmt.Errorf("error setting routing_vpcs in state: %s", err)
	}
	if err = d.Set("state", res.State); err != nil {
		return fmt.Errorf("error setting state in state: %s", err)
	}
	if err = d.Set("unique_id", res.UniqueId); err != nil {
		return fmt.Errorf("error setting unique_id in state: %s", err)
	}
	if err = d.Set("update_time", res.UpdateTime); err != nil {
		return fmt.Errorf("error setting update_time in state: %s", err)
	}

	return nil
}
func resourceNetworkConnectivityHubUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &networkconnectivity.Hub{
		Name:        dcl.String(d.Get("name").(string)),
		Description: dcl.String(d.Get("description").(string)),
		Labels:      tpgresource.CheckStringMap(d.Get("labels")),
		Project:     dcl.String(project),
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
	client := transport_tpg.NewDCLNetworkConnectivityClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutUpdate))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyHub(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error updating Hub: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Hub %q: %#v", d.Id(), res)

	return resourceNetworkConnectivityHubRead(d, meta)
}

func resourceNetworkConnectivityHubDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &networkconnectivity.Hub{
		Name:        dcl.String(d.Get("name").(string)),
		Description: dcl.String(d.Get("description").(string)),
		Labels:      tpgresource.CheckStringMap(d.Get("labels")),
		Project:     dcl.String(project),
	}

	log.Printf("[DEBUG] Deleting Hub %q", d.Id())
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLNetworkConnectivityClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutDelete))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	if err := client.DeleteHub(context.Background(), obj); err != nil {
		return fmt.Errorf("Error deleting Hub: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting Hub %q", d.Id())
	return nil
}

func resourceNetworkConnectivityHubImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)

	if err := tpgresource.ParseImportId([]string{
		"projects/(?P<project>[^/]+)/locations/global/hubs/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<name>[^/]+)",
		"(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVarsForId(d, config, "projects/{{project}}/locations/global/hubs/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenNetworkConnectivityHubRoutingVpcsArray(objs []networkconnectivity.HubRoutingVpcs) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenNetworkConnectivityHubRoutingVpcs(&item)
		items = append(items, i)
	}

	return items
}

func flattenNetworkConnectivityHubRoutingVpcs(obj *networkconnectivity.HubRoutingVpcs) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"uri": obj.Uri,
	}

	return transformed

}
