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

func resourceComputeFirewallPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeFirewallPolicyCreate,
		Read:   resourceComputeFirewallPolicyRead,
		Update: resourceComputeFirewallPolicyUpdate,
		Delete: resourceComputeFirewallPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceComputeFirewallPolicyImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"parent": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      ``,
			},

			"short_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: ``,
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: ``,
			},

			"creation_timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: ``,
			},

			"fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: ``,
			},

			"firewall_policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: ``,
			},

			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: ``,
			},

			"rule_tuple_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: ``,
			},

			"self_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: ``,
			},

			"self_link_with_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: ``,
			},
		},
	}
}

func resourceComputeFirewallPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	obj := &compute.FirewallPolicy{
		Parent:      dcl.String(d.Get("parent").(string)),
		ShortName:   dcl.String(d.Get("short_name").(string)),
		Description: dcl.String(d.Get("description").(string)),
	}

	id, err := replaceVars(d, config, "locations/global/firewallPolicies/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
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
	res, err := client.ApplyFirewallPolicy(context.Background(), obj, createDirective...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating FirewallPolicy: %s", err)
	}

	log.Printf("[DEBUG] Finished creating FirewallPolicy %q: %#v", d.Id(), res)

	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}
	// Id has a server-generated value, set again after creation
	id, err = replaceVars(d, config, "locations/global/firewallPolicies/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return resourceComputeFirewallPolicyRead(d, meta)
}

func resourceComputeFirewallPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	obj := &compute.FirewallPolicy{
		Parent:      dcl.String(d.Get("parent").(string)),
		ShortName:   dcl.String(d.Get("short_name").(string)),
		Description: dcl.String(d.Get("description").(string)),
		Name:        dcl.StringOrNil(d.Get("name").(string)),
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
	res, err := client.GetFirewallPolicy(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("ComputeFirewallPolicy %q", d.Id())
		return handleNotFoundDCLError(err, d, resourceName)
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
	config := meta.(*Config)

	obj := &compute.FirewallPolicy{
		Parent:      dcl.String(d.Get("parent").(string)),
		ShortName:   dcl.String(d.Get("short_name").(string)),
		Description: dcl.String(d.Get("description").(string)),
		Name:        dcl.StringOrNil(d.Get("name").(string)),
	}
	directive := UpdateDirective
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
	config := meta.(*Config)

	obj := &compute.FirewallPolicy{
		Parent:      dcl.String(d.Get("parent").(string)),
		ShortName:   dcl.String(d.Get("short_name").(string)),
		Description: dcl.String(d.Get("description").(string)),
		Name:        dcl.StringOrNil(d.Get("name").(string)),
	}

	log.Printf("[DEBUG] Deleting FirewallPolicy %q", d.Id())
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
	if err := client.DeleteFirewallPolicy(context.Background(), obj); err != nil {
		return fmt.Errorf("Error deleting FirewallPolicy: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting FirewallPolicy %q", d.Id())
	return nil
}

func resourceComputeFirewallPolicyImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"locations/global/firewallPolicies/(?P<name>[^/]+)",
		"(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "locations/global/firewallPolicies/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
