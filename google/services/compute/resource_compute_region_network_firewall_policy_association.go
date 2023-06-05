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

package compute

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	compute "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/compute"

	"github.com/hashicorp/terraform-provider-google/google/tpgdclresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceComputeRegionNetworkFirewallPolicyAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeRegionNetworkFirewallPolicyAssociationCreate,
		Read:   resourceComputeRegionNetworkFirewallPolicyAssociationRead,
		Delete: resourceComputeRegionNetworkFirewallPolicyAssociationDelete,

		Importer: &schema.ResourceImporter{
			State: resourceComputeRegionNetworkFirewallPolicyAssociationImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"attachment_target": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "The target that the firewall policy is attached to.",
			},

			"firewall_policy": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "The firewall policy ID of the association.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name for an association.",
			},

			"project": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "The project for the resource",
			},

			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "The location of this resource.",
			},

			"short_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The short name of the firewall policy of the association.",
			},
		},
	}
}

func resourceComputeRegionNetworkFirewallPolicyAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	region, err := tpgresource.GetRegion(d, config)
	if err != nil {
		return err
	}

	obj := &compute.NetworkFirewallPolicyAssociation{
		AttachmentTarget: dcl.String(d.Get("attachment_target").(string)),
		FirewallPolicy:   dcl.String(d.Get("firewall_policy").(string)),
		Name:             dcl.String(d.Get("name").(string)),
		Project:          dcl.String(project),
		Location:         dcl.String(region),
	}

	id, err := tpgresource.ReplaceVarsForId(d, config, "projects/{{project}}/regions/{{region}}/firewallPolicies/{{firewall_policy}}/associations/{{name}}")
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
	client := transport_tpg.NewDCLComputeClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutCreate))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyNetworkFirewallPolicyAssociation(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating NetworkFirewallPolicyAssociation: %s", err)
	}

	log.Printf("[DEBUG] Finished creating NetworkFirewallPolicyAssociation %q: %#v", d.Id(), res)

	return resourceComputeRegionNetworkFirewallPolicyAssociationRead(d, meta)
}

func resourceComputeRegionNetworkFirewallPolicyAssociationRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	region, err := tpgresource.GetRegion(d, config)
	if err != nil {
		return err
	}

	obj := &compute.NetworkFirewallPolicyAssociation{
		AttachmentTarget: dcl.String(d.Get("attachment_target").(string)),
		FirewallPolicy:   dcl.String(d.Get("firewall_policy").(string)),
		Name:             dcl.String(d.Get("name").(string)),
		Project:          dcl.String(project),
		Location:         dcl.String(region),
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
	client := transport_tpg.NewDCLComputeClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutRead))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.GetNetworkFirewallPolicyAssociation(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("ComputeRegionNetworkFirewallPolicyAssociation %q", d.Id())
		return tpgdclresource.HandleNotFoundDCLError(err, d, resourceName)
	}

	if err = d.Set("attachment_target", res.AttachmentTarget); err != nil {
		return fmt.Errorf("error setting attachment_target in state: %s", err)
	}
	if err = d.Set("firewall_policy", res.FirewallPolicy); err != nil {
		return fmt.Errorf("error setting firewall_policy in state: %s", err)
	}
	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}
	if err = d.Set("project", res.Project); err != nil {
		return fmt.Errorf("error setting project in state: %s", err)
	}
	if err = d.Set("region", res.Location); err != nil {
		return fmt.Errorf("error setting region in state: %s", err)
	}
	if err = d.Set("short_name", res.ShortName); err != nil {
		return fmt.Errorf("error setting short_name in state: %s", err)
	}

	return nil
}

func resourceComputeRegionNetworkFirewallPolicyAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	region, err := tpgresource.GetRegion(d, config)
	if err != nil {
		return err
	}

	obj := &compute.NetworkFirewallPolicyAssociation{
		AttachmentTarget: dcl.String(d.Get("attachment_target").(string)),
		FirewallPolicy:   dcl.String(d.Get("firewall_policy").(string)),
		Name:             dcl.String(d.Get("name").(string)),
		Project:          dcl.String(project),
		Location:         dcl.String(region),
	}

	log.Printf("[DEBUG] Deleting NetworkFirewallPolicyAssociation %q", d.Id())
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLComputeClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutDelete))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	if err := client.DeleteNetworkFirewallPolicyAssociation(context.Background(), obj); err != nil {
		return fmt.Errorf("Error deleting NetworkFirewallPolicyAssociation: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting NetworkFirewallPolicyAssociation %q", d.Id())
	return nil
}

func resourceComputeRegionNetworkFirewallPolicyAssociationImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)

	if err := tpgresource.ParseImportId([]string{
		"projects/(?P<project>[^/]+)/regions/(?P<region>[^/]+)/firewallPolicies/(?P<firewall_policy>[^/]+)/associations/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<region>[^/]+)/(?P<firewall_policy>[^/]+)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVarsForId(d, config, "projects/{{project}}/regions/{{region}}/firewallPolicies/{{firewall_policy}}/associations/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
