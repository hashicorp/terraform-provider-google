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

func ResourceComputeFirewallPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeFirewallPolicyCreate,
		Read:   resourceComputeFirewallPolicyRead,
		Update: resourceComputeFirewallPolicyUpdate,
		Delete: resourceComputeFirewallPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceComputeFirewallPolicyImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"parent": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "The parent of the firewall policy.",
			},

			"short_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "User-provided name of the Organization firewall policy. The name should be unique in the organization in which the firewall policy is created. The name must be 1-63 characters long, and comply with RFC1035. Specifically, the name must be 1-63 characters long and match the regular expression [a-z]([-a-z0-9]*[a-z0-9])? which means the first character must be a lowercase letter, and all following characters must be a dash, lowercase letter, or digit, except the last character, which cannot be a dash.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An optional description of this resource. Provide this property when you create the resource.",
			},

			"creation_timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation timestamp in RFC3339 text format.",
			},

			"fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Fingerprint of the resource. This field is used internally during updates of this resource.",
			},

			"firewall_policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for the resource. This identifier is defined by the server.",
			},

			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the resource. It is a numeric ID allocated by GCP which uniquely identifies the Firewall Policy.",
			},

			"rule_tuple_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total count of all firewall policy rule tuples. A firewall policy can not exceed a set number of tuples.",
			},

			"self_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server-defined URL for the resource.",
			},

			"self_link_with_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server-defined URL for this resource with the resource id.",
			},
		},
	}
}

func resourceComputeFirewallPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	obj := &compute.FirewallPolicy{
		Parent:      dcl.String(d.Get("parent").(string)),
		ShortName:   dcl.String(d.Get("short_name").(string)),
		Description: dcl.String(d.Get("description").(string)),
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
	billingProject := ""
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
	res, err := client.ApplyFirewallPolicy(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating FirewallPolicy: %s", err)
	}

	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}
	// ID has a server-generated value, set again after creation.

	id, err = res.ID()
	if err != nil {
		return fmt.Errorf("error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating FirewallPolicy %q: %#v", d.Id(), res)

	return resourceComputeFirewallPolicyRead(d, meta)
}

func resourceComputeFirewallPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	obj := &compute.FirewallPolicy{
		Parent:      dcl.String(d.Get("parent").(string)),
		ShortName:   dcl.String(d.Get("short_name").(string)),
		Description: dcl.String(d.Get("description").(string)),
		Name:        dcl.StringOrNil(d.Get("name").(string)),
	}

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := ""
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
	res, err := client.GetFirewallPolicy(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("ComputeFirewallPolicy %q", d.Id())
		return tpgdclresource.HandleNotFoundDCLError(err, d, resourceName)
	}

	if err = d.Set("parent", res.Parent); err != nil {
		return fmt.Errorf("error setting parent in state: %s", err)
	}
	if err = d.Set("short_name", res.ShortName); err != nil {
		return fmt.Errorf("error setting short_name in state: %s", err)
	}
	if err = d.Set("description", res.Description); err != nil {
		return fmt.Errorf("error setting description in state: %s", err)
	}
	if err = d.Set("creation_timestamp", res.CreationTimestamp); err != nil {
		return fmt.Errorf("error setting creation_timestamp in state: %s", err)
	}
	if err = d.Set("fingerprint", res.Fingerprint); err != nil {
		return fmt.Errorf("error setting fingerprint in state: %s", err)
	}
	if err = d.Set("firewall_policy_id", res.Id); err != nil {
		return fmt.Errorf("error setting firewall_policy_id in state: %s", err)
	}
	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}
	if err = d.Set("rule_tuple_count", res.RuleTupleCount); err != nil {
		return fmt.Errorf("error setting rule_tuple_count in state: %s", err)
	}
	if err = d.Set("self_link", res.SelfLink); err != nil {
		return fmt.Errorf("error setting self_link in state: %s", err)
	}
	if err = d.Set("self_link_with_id", res.SelfLinkWithId); err != nil {
		return fmt.Errorf("error setting self_link_with_id in state: %s", err)
	}

	return nil
}
func resourceComputeFirewallPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	obj := &compute.FirewallPolicy{
		Parent:      dcl.String(d.Get("parent").(string)),
		ShortName:   dcl.String(d.Get("short_name").(string)),
		Description: dcl.String(d.Get("description").(string)),
		Name:        dcl.StringOrNil(d.Get("name").(string)),
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
	client := transport_tpg.NewDCLComputeClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutUpdate))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyFirewallPolicy(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error updating FirewallPolicy: %s", err)
	}

	log.Printf("[DEBUG] Finished creating FirewallPolicy %q: %#v", d.Id(), res)

	return resourceComputeFirewallPolicyRead(d, meta)
}

func resourceComputeFirewallPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	obj := &compute.FirewallPolicy{
		Parent:      dcl.String(d.Get("parent").(string)),
		ShortName:   dcl.String(d.Get("short_name").(string)),
		Description: dcl.String(d.Get("description").(string)),
		Name:        dcl.StringOrNil(d.Get("name").(string)),
	}

	log.Printf("[DEBUG] Deleting FirewallPolicy %q", d.Id())
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := ""
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
	if err := client.DeleteFirewallPolicy(context.Background(), obj); err != nil {
		return fmt.Errorf("Error deleting FirewallPolicy: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting FirewallPolicy %q", d.Id())
	return nil
}

func resourceComputeFirewallPolicyImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)

	if err := tpgresource.ParseImportId([]string{
		"locations/global/firewallPolicies/(?P<name>[^/]+)",
		"(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "locations/global/firewallPolicies/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
