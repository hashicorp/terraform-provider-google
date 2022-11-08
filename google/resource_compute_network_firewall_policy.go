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

func resourceComputeNetworkFirewallPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeNetworkFirewallPolicyCreate,
		Read:   resourceComputeNetworkFirewallPolicyRead,
		Update: resourceComputeNetworkFirewallPolicyUpdate,
		Delete: resourceComputeNetworkFirewallPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceComputeNetworkFirewallPolicyImport,
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
				Description: "User-provided name of the Network firewall policy. The name should be unique in the project in which the firewall policy is created. The name must be 1-63 characters long, and comply with RFC1035. Specifically, the name must be 1-63 characters long and match the regular expression [a-z]([-a-z0-9]*[a-z0-9])? which means the first character must be a lowercase letter, and all following characters must be a dash, lowercase letter, or digit, except the last character, which cannot be a dash.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An optional description of this resource. Provide this property when you create the resource.",
			},

			"project": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      "The project for the resource",
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

			"network_firewall_policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for the resource. This identifier is defined by the server.",
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

func resourceComputeNetworkFirewallPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	obj := &compute.NetworkFirewallPolicy{
		Name:        dcl.String(d.Get("name").(string)),
		Description: dcl.String(d.Get("description").(string)),
		Project:     dcl.String(project),
	}

	id, err := obj.ID()
	if err != nil {
		return fmt.Errorf("error constructing id: %s", err)
	}
	d.SetId(id)
	directive := CreateDirective
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := NewDCLComputeClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutCreate))
	if bp, err := replaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyNetworkFirewallPolicy(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating NetworkFirewallPolicy: %s", err)
	}

	log.Printf("[DEBUG] Finished creating NetworkFirewallPolicy %q: %#v", d.Id(), res)

	return resourceComputeNetworkFirewallPolicyRead(d, meta)
}

func resourceComputeNetworkFirewallPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	obj := &compute.NetworkFirewallPolicy{
		Name:        dcl.String(d.Get("name").(string)),
		Description: dcl.String(d.Get("description").(string)),
		Project:     dcl.String(project),
	}

	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := NewDCLComputeClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutRead))
	if bp, err := replaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.GetNetworkFirewallPolicy(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("ComputeNetworkFirewallPolicy %q", d.Id())
		return handleNotFoundDCLError(err, d, resourceName)
	}

	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}
	if err = d.Set("description", res.Description); err != nil {
		return fmt.Errorf("error setting description in state: %s", err)
	}
	if err = d.Set("project", res.Project); err != nil {
		return fmt.Errorf("error setting project in state: %s", err)
	}
	if err = d.Set("creation_timestamp", res.CreationTimestamp); err != nil {
		return fmt.Errorf("error setting creation_timestamp in state: %s", err)
	}
	if err = d.Set("fingerprint", res.Fingerprint); err != nil {
		return fmt.Errorf("error setting fingerprint in state: %s", err)
	}
	if err = d.Set("network_firewall_policy_id", res.Id); err != nil {
		return fmt.Errorf("error setting network_firewall_policy_id in state: %s", err)
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
func resourceComputeNetworkFirewallPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	obj := &compute.NetworkFirewallPolicy{
		Name:        dcl.String(d.Get("name").(string)),
		Description: dcl.String(d.Get("description").(string)),
		Project:     dcl.String(project),
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
	client := NewDCLComputeClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutUpdate))
	if bp, err := replaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyNetworkFirewallPolicy(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error updating NetworkFirewallPolicy: %s", err)
	}

	log.Printf("[DEBUG] Finished creating NetworkFirewallPolicy %q: %#v", d.Id(), res)

	return resourceComputeNetworkFirewallPolicyRead(d, meta)
}

func resourceComputeNetworkFirewallPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	obj := &compute.NetworkFirewallPolicy{
		Name:        dcl.String(d.Get("name").(string)),
		Description: dcl.String(d.Get("description").(string)),
		Project:     dcl.String(project),
	}

	log.Printf("[DEBUG] Deleting NetworkFirewallPolicy %q", d.Id())
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := NewDCLComputeClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutDelete))
	if bp, err := replaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	if err := client.DeleteNetworkFirewallPolicy(context.Background(), obj); err != nil {
		return fmt.Errorf("Error deleting NetworkFirewallPolicy: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting NetworkFirewallPolicy %q", d.Id())
	return nil
}

func resourceComputeNetworkFirewallPolicyImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/global/firewallPolicies/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<name>[^/]+)",
		"(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVarsForId(d, config, "projects/{{project}}/global/firewallPolicies/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
