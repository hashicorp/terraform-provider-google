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

package google

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	compute "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/compute"
)

func resourceComputeFirewallPolicyAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeFirewallPolicyAssociationCreate,
		Read:   resourceComputeFirewallPolicyAssociationRead,
		Delete: resourceComputeFirewallPolicyAssociationDelete,

		Importer: &schema.ResourceImporter{
			State: resourceComputeFirewallPolicyAssociationImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"attachment_target": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      "The target that the firewall policy is attached to.",
			},

			"firewall_policy": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      "The firewall policy ID of the association.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name for an association.",
			},

			"short_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The short name of the firewall policy of the association.",
			},
		},
	}
}

func resourceComputeFirewallPolicyAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	obj := &compute.FirewallPolicyAssociation{
		AttachmentTarget: dcl.String(d.Get("attachment_target").(string)),
		FirewallPolicy:   dcl.String(d.Get("firewall_policy").(string)),
		Name:             dcl.String(d.Get("name").(string)),
	}

	id, err := replaceVarsForId(d, config, "locations/global/firewallPolicies/{{firewall_policy}}/associations/{{name}}")
	if err != nil {
		return fmt.Errorf("error constructing id: %s", err)
	}
	d.SetId(id)
	createDirective := CreateDirective
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}
	billingProject := ""
	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := NewDCLComputeClient(config, userAgent, billingProject)
	res, err := client.ApplyFirewallPolicyAssociation(context.Background(), obj, createDirective...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating FirewallPolicyAssociation: %s", err)
	}

	log.Printf("[DEBUG] Finished creating FirewallPolicyAssociation %q: %#v", d.Id(), res)

	return resourceComputeFirewallPolicyAssociationRead(d, meta)
}

func resourceComputeFirewallPolicyAssociationRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	obj := &compute.FirewallPolicyAssociation{
		AttachmentTarget: dcl.String(d.Get("attachment_target").(string)),
		FirewallPolicy:   dcl.String(d.Get("firewall_policy").(string)),
		Name:             dcl.String(d.Get("name").(string)),
	}

	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}
	billingProject := ""
	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := NewDCLComputeClient(config, userAgent, billingProject)
	res, err := client.GetFirewallPolicyAssociation(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("ComputeFirewallPolicyAssociation %q", d.Id())
		return handleNotFoundDCLError(err, d, resourceName)
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
	if err = d.Set("short_name", res.ShortName); err != nil {
		return fmt.Errorf("error setting short_name in state: %s", err)
	}

	return nil
}

func resourceComputeFirewallPolicyAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	obj := &compute.FirewallPolicyAssociation{
		AttachmentTarget: dcl.String(d.Get("attachment_target").(string)),
		FirewallPolicy:   dcl.String(d.Get("firewall_policy").(string)),
		Name:             dcl.String(d.Get("name").(string)),
	}

	log.Printf("[DEBUG] Deleting FirewallPolicyAssociation %q", d.Id())
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}
	billingProject := ""
	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := NewDCLComputeClient(config, userAgent, billingProject)
	if err := client.DeleteFirewallPolicyAssociation(context.Background(), obj); err != nil {
		return fmt.Errorf("Error deleting FirewallPolicyAssociation: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting FirewallPolicyAssociation %q", d.Id())
	return nil
}

func resourceComputeFirewallPolicyAssociationImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"locations/global/firewallPolicies/(?P<firewall_policy>[^/]+)/associations/(?P<name>[^/]+)",
		"(?P<firewall_policy>[^/]+)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVarsForId(d, config, "locations/global/firewallPolicies/{{firewall_policy}}/associations/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
