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

func ResourceNetworkConnectivitySpoke() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkConnectivitySpokeCreate,
		Read:   resourceNetworkConnectivitySpokeRead,
		Update: resourceNetworkConnectivitySpokeUpdate,
		Delete: resourceNetworkConnectivitySpokeDelete,

		Importer: &schema.ResourceImporter{
			State: resourceNetworkConnectivitySpokeImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"hub": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "Immutable. The URI of the hub that this spoke is attached to.",
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
				Description: "Immutable. The name of the spoke. Spoke names must be unique.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An optional description of the spoke.",
			},

			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Optional labels in key:value format. For more information about labels, see [Requirements for labels](https://cloud.google.com/resource-manager/docs/creating-managing-labels#requirements).",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"linked_interconnect_attachments": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				Description:   "A collection of VLAN attachment resources. These resources should be redundant attachments that all advertise the same prefixes to Google Cloud. Alternatively, in active/passive configurations, all attachments should be capable of advertising the same prefixes.",
				MaxItems:      1,
				Elem:          NetworkConnectivitySpokeLinkedInterconnectAttachmentsSchema(),
				ConflictsWith: []string{"linked_vpn_tunnels", "linked_router_appliance_instances"},
			},

			"linked_router_appliance_instances": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				Description:   "The URIs of linked Router appliance resources",
				MaxItems:      1,
				Elem:          NetworkConnectivitySpokeLinkedRouterApplianceInstancesSchema(),
				ConflictsWith: []string{"linked_vpn_tunnels", "linked_interconnect_attachments"},
			},

			"linked_vpn_tunnels": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				Description:   "The URIs of linked VPN tunnel resources",
				MaxItems:      1,
				Elem:          NetworkConnectivitySpokeLinkedVpnTunnelsSchema(),
				ConflictsWith: []string{"linked_interconnect_attachments", "linked_router_appliance_instances"},
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
				Description: "Output only. The time the spoke was created.",
			},

			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The current lifecycle state of this spoke. Possible values: STATE_UNSPECIFIED, CREATING, ACTIVE, DELETING",
			},

			"unique_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The Google-generated UUID for the spoke. This value is unique across all spoke resources. If a spoke is deleted and another with the same name is created, the new spoke is assigned a different unique_id.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The time the spoke was last updated.",
			},
		},
	}
}

func NetworkConnectivitySpokeLinkedInterconnectAttachmentsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"site_to_site_data_transfer": {
				Type:        schema.TypeBool,
				Required:    true,
				ForceNew:    true,
				Description: "A value that controls whether site-to-site data transfer is enabled for these resources. Note that data transfer is available only in supported locations.",
			},

			"uris": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "The URIs of linked interconnect attachment resources",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func NetworkConnectivitySpokeLinkedRouterApplianceInstancesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instances": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "The list of router appliance instances",
				Elem:        NetworkConnectivitySpokeLinkedRouterApplianceInstancesInstancesSchema(),
			},

			"site_to_site_data_transfer": {
				Type:        schema.TypeBool,
				Required:    true,
				ForceNew:    true,
				Description: "A value that controls whether site-to-site data transfer is enabled for these resources. Note that data transfer is available only in supported locations.",
			},
		},
	}
}

func NetworkConnectivitySpokeLinkedRouterApplianceInstancesInstancesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The IP address on the VM to use for peering.",
			},

			"virtual_machine": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "The URI of the virtual machine resource",
			},
		},
	}
}

func NetworkConnectivitySpokeLinkedVpnTunnelsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"site_to_site_data_transfer": {
				Type:        schema.TypeBool,
				Required:    true,
				ForceNew:    true,
				Description: "A value that controls whether site-to-site data transfer is enabled for these resources. Note that data transfer is available only in supported locations.",
			},

			"uris": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "The URIs of linked VPN tunnel resources.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceNetworkConnectivitySpokeCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &networkconnectivity.Spoke{
		Hub:                            dcl.String(d.Get("hub").(string)),
		Location:                       dcl.String(d.Get("location").(string)),
		Name:                           dcl.String(d.Get("name").(string)),
		Description:                    dcl.String(d.Get("description").(string)),
		Labels:                         tpgresource.CheckStringMap(d.Get("labels")),
		LinkedInterconnectAttachments:  expandNetworkConnectivitySpokeLinkedInterconnectAttachments(d.Get("linked_interconnect_attachments")),
		LinkedRouterApplianceInstances: expandNetworkConnectivitySpokeLinkedRouterApplianceInstances(d.Get("linked_router_appliance_instances")),
		LinkedVpnTunnels:               expandNetworkConnectivitySpokeLinkedVpnTunnels(d.Get("linked_vpn_tunnels")),
		Project:                        dcl.String(project),
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
	res, err := client.ApplySpoke(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating Spoke: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Spoke %q: %#v", d.Id(), res)

	return resourceNetworkConnectivitySpokeRead(d, meta)
}

func resourceNetworkConnectivitySpokeRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &networkconnectivity.Spoke{
		Hub:                            dcl.String(d.Get("hub").(string)),
		Location:                       dcl.String(d.Get("location").(string)),
		Name:                           dcl.String(d.Get("name").(string)),
		Description:                    dcl.String(d.Get("description").(string)),
		Labels:                         tpgresource.CheckStringMap(d.Get("labels")),
		LinkedInterconnectAttachments:  expandNetworkConnectivitySpokeLinkedInterconnectAttachments(d.Get("linked_interconnect_attachments")),
		LinkedRouterApplianceInstances: expandNetworkConnectivitySpokeLinkedRouterApplianceInstances(d.Get("linked_router_appliance_instances")),
		LinkedVpnTunnels:               expandNetworkConnectivitySpokeLinkedVpnTunnels(d.Get("linked_vpn_tunnels")),
		Project:                        dcl.String(project),
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
	res, err := client.GetSpoke(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("NetworkConnectivitySpoke %q", d.Id())
		return tpgdclresource.HandleNotFoundDCLError(err, d, resourceName)
	}

	if err = d.Set("hub", res.Hub); err != nil {
		return fmt.Errorf("error setting hub in state: %s", err)
	}
	if err = d.Set("location", res.Location); err != nil {
		return fmt.Errorf("error setting location in state: %s", err)
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
	if err = d.Set("linked_interconnect_attachments", flattenNetworkConnectivitySpokeLinkedInterconnectAttachments(res.LinkedInterconnectAttachments)); err != nil {
		return fmt.Errorf("error setting linked_interconnect_attachments in state: %s", err)
	}
	if err = d.Set("linked_router_appliance_instances", flattenNetworkConnectivitySpokeLinkedRouterApplianceInstances(res.LinkedRouterApplianceInstances)); err != nil {
		return fmt.Errorf("error setting linked_router_appliance_instances in state: %s", err)
	}
	if err = d.Set("linked_vpn_tunnels", flattenNetworkConnectivitySpokeLinkedVpnTunnels(res.LinkedVpnTunnels)); err != nil {
		return fmt.Errorf("error setting linked_vpn_tunnels in state: %s", err)
	}
	if err = d.Set("project", res.Project); err != nil {
		return fmt.Errorf("error setting project in state: %s", err)
	}
	if err = d.Set("create_time", res.CreateTime); err != nil {
		return fmt.Errorf("error setting create_time in state: %s", err)
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
func resourceNetworkConnectivitySpokeUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &networkconnectivity.Spoke{
		Hub:                            dcl.String(d.Get("hub").(string)),
		Location:                       dcl.String(d.Get("location").(string)),
		Name:                           dcl.String(d.Get("name").(string)),
		Description:                    dcl.String(d.Get("description").(string)),
		Labels:                         tpgresource.CheckStringMap(d.Get("labels")),
		LinkedInterconnectAttachments:  expandNetworkConnectivitySpokeLinkedInterconnectAttachments(d.Get("linked_interconnect_attachments")),
		LinkedRouterApplianceInstances: expandNetworkConnectivitySpokeLinkedRouterApplianceInstances(d.Get("linked_router_appliance_instances")),
		LinkedVpnTunnels:               expandNetworkConnectivitySpokeLinkedVpnTunnels(d.Get("linked_vpn_tunnels")),
		Project:                        dcl.String(project),
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
	res, err := client.ApplySpoke(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error updating Spoke: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Spoke %q: %#v", d.Id(), res)

	return resourceNetworkConnectivitySpokeRead(d, meta)
}

func resourceNetworkConnectivitySpokeDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &networkconnectivity.Spoke{
		Hub:                            dcl.String(d.Get("hub").(string)),
		Location:                       dcl.String(d.Get("location").(string)),
		Name:                           dcl.String(d.Get("name").(string)),
		Description:                    dcl.String(d.Get("description").(string)),
		Labels:                         tpgresource.CheckStringMap(d.Get("labels")),
		LinkedInterconnectAttachments:  expandNetworkConnectivitySpokeLinkedInterconnectAttachments(d.Get("linked_interconnect_attachments")),
		LinkedRouterApplianceInstances: expandNetworkConnectivitySpokeLinkedRouterApplianceInstances(d.Get("linked_router_appliance_instances")),
		LinkedVpnTunnels:               expandNetworkConnectivitySpokeLinkedVpnTunnels(d.Get("linked_vpn_tunnels")),
		Project:                        dcl.String(project),
	}

	log.Printf("[DEBUG] Deleting Spoke %q", d.Id())
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
	if err := client.DeleteSpoke(context.Background(), obj); err != nil {
		return fmt.Errorf("Error deleting Spoke: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting Spoke %q", d.Id())
	return nil
}

func resourceNetworkConnectivitySpokeImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)

	if err := tpgresource.ParseImportId([]string{
		"projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/spokes/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<name>[^/]+)",
		"(?P<location>[^/]+)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVarsForId(d, config, "projects/{{project}}/locations/{{location}}/spokes/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func expandNetworkConnectivitySpokeLinkedInterconnectAttachments(o interface{}) *networkconnectivity.SpokeLinkedInterconnectAttachments {
	if o == nil {
		return networkconnectivity.EmptySpokeLinkedInterconnectAttachments
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return networkconnectivity.EmptySpokeLinkedInterconnectAttachments
	}
	obj := objArr[0].(map[string]interface{})
	return &networkconnectivity.SpokeLinkedInterconnectAttachments{
		SiteToSiteDataTransfer: dcl.Bool(obj["site_to_site_data_transfer"].(bool)),
		Uris:                   tpgdclresource.ExpandStringArray(obj["uris"]),
	}
}

func flattenNetworkConnectivitySpokeLinkedInterconnectAttachments(obj *networkconnectivity.SpokeLinkedInterconnectAttachments) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"site_to_site_data_transfer": obj.SiteToSiteDataTransfer,
		"uris":                       obj.Uris,
	}

	return []interface{}{transformed}

}

func expandNetworkConnectivitySpokeLinkedRouterApplianceInstances(o interface{}) *networkconnectivity.SpokeLinkedRouterApplianceInstances {
	if o == nil {
		return networkconnectivity.EmptySpokeLinkedRouterApplianceInstances
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return networkconnectivity.EmptySpokeLinkedRouterApplianceInstances
	}
	obj := objArr[0].(map[string]interface{})
	return &networkconnectivity.SpokeLinkedRouterApplianceInstances{
		Instances:              expandNetworkConnectivitySpokeLinkedRouterApplianceInstancesInstancesArray(obj["instances"]),
		SiteToSiteDataTransfer: dcl.Bool(obj["site_to_site_data_transfer"].(bool)),
	}
}

func flattenNetworkConnectivitySpokeLinkedRouterApplianceInstances(obj *networkconnectivity.SpokeLinkedRouterApplianceInstances) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"instances":                  flattenNetworkConnectivitySpokeLinkedRouterApplianceInstancesInstancesArray(obj.Instances),
		"site_to_site_data_transfer": obj.SiteToSiteDataTransfer,
	}

	return []interface{}{transformed}

}
func expandNetworkConnectivitySpokeLinkedRouterApplianceInstancesInstancesArray(o interface{}) []networkconnectivity.SpokeLinkedRouterApplianceInstancesInstances {
	if o == nil {
		return make([]networkconnectivity.SpokeLinkedRouterApplianceInstancesInstances, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make([]networkconnectivity.SpokeLinkedRouterApplianceInstancesInstances, 0)
	}

	items := make([]networkconnectivity.SpokeLinkedRouterApplianceInstancesInstances, 0, len(objs))
	for _, item := range objs {
		i := expandNetworkConnectivitySpokeLinkedRouterApplianceInstancesInstances(item)
		items = append(items, *i)
	}

	return items
}

func expandNetworkConnectivitySpokeLinkedRouterApplianceInstancesInstances(o interface{}) *networkconnectivity.SpokeLinkedRouterApplianceInstancesInstances {
	if o == nil {
		return networkconnectivity.EmptySpokeLinkedRouterApplianceInstancesInstances
	}

	obj := o.(map[string]interface{})
	return &networkconnectivity.SpokeLinkedRouterApplianceInstancesInstances{
		IPAddress:      dcl.String(obj["ip_address"].(string)),
		VirtualMachine: dcl.String(obj["virtual_machine"].(string)),
	}
}

func flattenNetworkConnectivitySpokeLinkedRouterApplianceInstancesInstancesArray(objs []networkconnectivity.SpokeLinkedRouterApplianceInstancesInstances) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenNetworkConnectivitySpokeLinkedRouterApplianceInstancesInstances(&item)
		items = append(items, i)
	}

	return items
}

func flattenNetworkConnectivitySpokeLinkedRouterApplianceInstancesInstances(obj *networkconnectivity.SpokeLinkedRouterApplianceInstancesInstances) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"ip_address":      obj.IPAddress,
		"virtual_machine": obj.VirtualMachine,
	}

	return transformed

}

func expandNetworkConnectivitySpokeLinkedVpnTunnels(o interface{}) *networkconnectivity.SpokeLinkedVpnTunnels {
	if o == nil {
		return networkconnectivity.EmptySpokeLinkedVpnTunnels
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return networkconnectivity.EmptySpokeLinkedVpnTunnels
	}
	obj := objArr[0].(map[string]interface{})
	return &networkconnectivity.SpokeLinkedVpnTunnels{
		SiteToSiteDataTransfer: dcl.Bool(obj["site_to_site_data_transfer"].(bool)),
		Uris:                   tpgdclresource.ExpandStringArray(obj["uris"]),
	}
}

func flattenNetworkConnectivitySpokeLinkedVpnTunnels(obj *networkconnectivity.SpokeLinkedVpnTunnels) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"site_to_site_data_transfer": obj.SiteToSiteDataTransfer,
		"uris":                       obj.Uris,
	}

	return []interface{}{transformed}

}
