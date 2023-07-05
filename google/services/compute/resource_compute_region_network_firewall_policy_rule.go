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

func ResourceComputeRegionNetworkFirewallPolicyRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeRegionNetworkFirewallPolicyRuleCreate,
		Read:   resourceComputeRegionNetworkFirewallPolicyRuleRead,
		Update: resourceComputeRegionNetworkFirewallPolicyRuleUpdate,
		Delete: resourceComputeRegionNetworkFirewallPolicyRuleDelete,

		Importer: &schema.ResourceImporter{
			State: resourceComputeRegionNetworkFirewallPolicyRuleImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Action to perform when the client connection triggers the rule. Valid actions are \"allow\", \"deny\" and \"goto_next\".",
			},

			"direction": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The direction in which this rule applies. Possible values: INGRESS, EGRESS",
			},

			"firewall_policy": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "The firewall policy of the resource.",
			},

			"match": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "A match condition that incoming traffic is evaluated against. If it evaluates to true, the corresponding 'action' is enforced.",
				MaxItems:    1,
				Elem:        ComputeRegionNetworkFirewallPolicyRuleMatchSchema(),
			},

			"priority": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "An integer indicating the priority of a rule in the list. The priority must be a positive value between 0 and 2147483647. Rules are evaluated from highest to lowest priority where 0 is the highest priority and 2147483647 is the lowest prority.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An optional description for this resource.",
			},

			"disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Denotes whether the firewall policy rule is disabled. When set to true, the firewall policy rule is not enforced and traffic behaves as if it did not exist. If this is unspecified, the firewall policy rule will be enabled.",
			},

			"enable_logging": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Denotes whether to enable logging for a particular rule. If logging is enabled, logs will be exported to the configured export destination in Stackdriver. Logs may be exported to BigQuery or Pub/Sub. Note: you cannot enable logging on \"goto_next\" rules.",
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

			"rule_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An optional name for the rule. This field is not a unique identifier and can be updated.",
			},

			"target_secure_tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of secure tags that controls which instances the firewall rule applies to. If <code>targetSecureTag</code> are specified, then the firewall rule applies only to instances in the VPC network that have one of those EFFECTIVE secure tags, if all the target_secure_tag are in INEFFECTIVE state, then this rule will be ignored. <code>targetSecureTag</code> may not be set at the same time as <code>targetServiceAccounts</code>. If neither <code>targetServiceAccounts</code> nor <code>targetSecureTag</code> are specified, the firewall rule applies to all instances on the specified network. Maximum number of target label tags allowed is 256.",
				Elem:        ComputeRegionNetworkFirewallPolicyRuleTargetSecureTagsSchema(),
			},

			"target_service_accounts": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of service accounts indicating the sets of instances that are applied with this rule.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"kind": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the resource. Always `compute#firewallPolicyRule` for firewall policy rules",
			},

			"rule_tuple_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Calculation of the complexity of a single firewall policy rule.",
			},
		},
	}
}

func ComputeRegionNetworkFirewallPolicyRuleMatchSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"layer4_configs": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Pairs of IP protocols and ports that the rule should match.",
				Elem:        ComputeRegionNetworkFirewallPolicyRuleMatchLayer4ConfigsSchema(),
			},

			"dest_address_groups": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Address groups which should be matched against the traffic destination. Maximum number of destination address groups is 10. Destination address groups is only supported in Egress rules.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"dest_fqdns": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Domain names that will be used to match against the resolved domain name of destination of traffic. Can only be specified if DIRECTION is egress.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"dest_ip_ranges": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "CIDR IP address range. Maximum number of destination CIDR IP ranges allowed is 5000.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"dest_region_codes": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Unicode country codes whose IP addresses will be used to match against the source of traffic. Can only be specified if DIRECTION is egress.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"dest_threat_intelligences": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Name of the Google Cloud Threat Intelligence list.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"src_address_groups": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Address groups which should be matched against the traffic source. Maximum number of source address groups is 10. Source address groups is only supported in Ingress rules.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"src_fqdns": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Domain names that will be used to match against the resolved domain name of source of traffic. Can only be specified if DIRECTION is ingress.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"src_ip_ranges": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "CIDR IP address range. Maximum number of source CIDR IP ranges allowed is 5000.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"src_region_codes": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Unicode country codes whose IP addresses will be used to match against the source of traffic. Can only be specified if DIRECTION is ingress.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"src_secure_tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of secure tag values, which should be matched at the source of the traffic. For INGRESS rule, if all the <code>srcSecureTag</code> are INEFFECTIVE, and there is no <code>srcIpRange</code>, this rule will be ignored. Maximum number of source tag values allowed is 256.",
				Elem:        ComputeRegionNetworkFirewallPolicyRuleMatchSrcSecureTagsSchema(),
			},

			"src_threat_intelligences": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Name of the Google Cloud Threat Intelligence list.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func ComputeRegionNetworkFirewallPolicyRuleMatchLayer4ConfigsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ip_protocol": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The IP protocol to which this rule applies. The protocol type is required when creating a firewall rule. This value can either be one of the following well known protocol strings (`tcp`, `udp`, `icmp`, `esp`, `ah`, `ipip`, `sctp`), or the IP protocol number.",
			},

			"ports": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "An optional list of ports to which this rule applies. This field is only applicable for UDP or TCP protocol. Each entry must be either an integer or a range. If not specified, this rule applies to connections through any port. Example inputs include: ``.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func ComputeRegionNetworkFirewallPolicyRuleMatchSrcSecureTagsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "Name of the secure tag, created with TagManager's TagValue API. @pattern tagValues/[0-9]+",
			},

			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "[Output Only] State of the secure tag, either `EFFECTIVE` or `INEFFECTIVE`. A secure tag is `INEFFECTIVE` when it is deleted or its network is deleted.",
			},
		},
	}
}

func ComputeRegionNetworkFirewallPolicyRuleTargetSecureTagsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "Name of the secure tag, created with TagManager's TagValue API. @pattern tagValues/[0-9]+",
			},

			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "[Output Only] State of the secure tag, either `EFFECTIVE` or `INEFFECTIVE`. A secure tag is `INEFFECTIVE` when it is deleted or its network is deleted.",
			},
		},
	}
}

func resourceComputeRegionNetworkFirewallPolicyRuleCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	region, err := tpgresource.GetRegion(d, config)
	if err != nil {
		return err
	}

	obj := &compute.NetworkFirewallPolicyRule{
		Action:                dcl.String(d.Get("action").(string)),
		Direction:             compute.NetworkFirewallPolicyRuleDirectionEnumRef(d.Get("direction").(string)),
		FirewallPolicy:        dcl.String(d.Get("firewall_policy").(string)),
		Match:                 expandComputeRegionNetworkFirewallPolicyRuleMatch(d.Get("match")),
		Priority:              dcl.Int64(int64(d.Get("priority").(int))),
		Description:           dcl.String(d.Get("description").(string)),
		Disabled:              dcl.Bool(d.Get("disabled").(bool)),
		EnableLogging:         dcl.Bool(d.Get("enable_logging").(bool)),
		Project:               dcl.String(project),
		Location:              dcl.String(region),
		RuleName:              dcl.String(d.Get("rule_name").(string)),
		TargetSecureTags:      expandComputeRegionNetworkFirewallPolicyRuleTargetSecureTagsArray(d.Get("target_secure_tags")),
		TargetServiceAccounts: tpgdclresource.ExpandStringArray(d.Get("target_service_accounts")),
	}

	id, err := tpgresource.ReplaceVarsForId(d, config, "projects/{{project}}/regions/{{region}}/firewallPolicies/{{firewall_policy}}/{{priority}}")
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
	res, err := client.ApplyNetworkFirewallPolicyRule(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating NetworkFirewallPolicyRule: %s", err)
	}

	log.Printf("[DEBUG] Finished creating NetworkFirewallPolicyRule %q: %#v", d.Id(), res)

	return resourceComputeRegionNetworkFirewallPolicyRuleRead(d, meta)
}

func resourceComputeRegionNetworkFirewallPolicyRuleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	region, err := tpgresource.GetRegion(d, config)
	if err != nil {
		return err
	}

	obj := &compute.NetworkFirewallPolicyRule{
		Action:                dcl.String(d.Get("action").(string)),
		Direction:             compute.NetworkFirewallPolicyRuleDirectionEnumRef(d.Get("direction").(string)),
		FirewallPolicy:        dcl.String(d.Get("firewall_policy").(string)),
		Match:                 expandComputeRegionNetworkFirewallPolicyRuleMatch(d.Get("match")),
		Priority:              dcl.Int64(int64(d.Get("priority").(int))),
		Description:           dcl.String(d.Get("description").(string)),
		Disabled:              dcl.Bool(d.Get("disabled").(bool)),
		EnableLogging:         dcl.Bool(d.Get("enable_logging").(bool)),
		Project:               dcl.String(project),
		Location:              dcl.String(region),
		RuleName:              dcl.String(d.Get("rule_name").(string)),
		TargetSecureTags:      expandComputeRegionNetworkFirewallPolicyRuleTargetSecureTagsArray(d.Get("target_secure_tags")),
		TargetServiceAccounts: tpgdclresource.ExpandStringArray(d.Get("target_service_accounts")),
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
	res, err := client.GetNetworkFirewallPolicyRule(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("ComputeRegionNetworkFirewallPolicyRule %q", d.Id())
		return tpgdclresource.HandleNotFoundDCLError(err, d, resourceName)
	}

	if err = d.Set("action", res.Action); err != nil {
		return fmt.Errorf("error setting action in state: %s", err)
	}
	if err = d.Set("direction", res.Direction); err != nil {
		return fmt.Errorf("error setting direction in state: %s", err)
	}
	if err = d.Set("firewall_policy", res.FirewallPolicy); err != nil {
		return fmt.Errorf("error setting firewall_policy in state: %s", err)
	}
	if err = d.Set("match", flattenComputeRegionNetworkFirewallPolicyRuleMatch(res.Match)); err != nil {
		return fmt.Errorf("error setting match in state: %s", err)
	}
	if err = d.Set("priority", res.Priority); err != nil {
		return fmt.Errorf("error setting priority in state: %s", err)
	}
	if err = d.Set("description", res.Description); err != nil {
		return fmt.Errorf("error setting description in state: %s", err)
	}
	if err = d.Set("disabled", res.Disabled); err != nil {
		return fmt.Errorf("error setting disabled in state: %s", err)
	}
	if err = d.Set("enable_logging", res.EnableLogging); err != nil {
		return fmt.Errorf("error setting enable_logging in state: %s", err)
	}
	if err = d.Set("project", res.Project); err != nil {
		return fmt.Errorf("error setting project in state: %s", err)
	}
	if err = d.Set("region", res.Location); err != nil {
		return fmt.Errorf("error setting region in state: %s", err)
	}
	if err = d.Set("rule_name", res.RuleName); err != nil {
		return fmt.Errorf("error setting rule_name in state: %s", err)
	}
	if err = d.Set("target_secure_tags", flattenComputeRegionNetworkFirewallPolicyRuleTargetSecureTagsArray(res.TargetSecureTags)); err != nil {
		return fmt.Errorf("error setting target_secure_tags in state: %s", err)
	}
	if err = d.Set("target_service_accounts", res.TargetServiceAccounts); err != nil {
		return fmt.Errorf("error setting target_service_accounts in state: %s", err)
	}
	if err = d.Set("kind", res.Kind); err != nil {
		return fmt.Errorf("error setting kind in state: %s", err)
	}
	if err = d.Set("rule_tuple_count", res.RuleTupleCount); err != nil {
		return fmt.Errorf("error setting rule_tuple_count in state: %s", err)
	}

	return nil
}
func resourceComputeRegionNetworkFirewallPolicyRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	region, err := tpgresource.GetRegion(d, config)
	if err != nil {
		return err
	}

	obj := &compute.NetworkFirewallPolicyRule{
		Action:                dcl.String(d.Get("action").(string)),
		Direction:             compute.NetworkFirewallPolicyRuleDirectionEnumRef(d.Get("direction").(string)),
		FirewallPolicy:        dcl.String(d.Get("firewall_policy").(string)),
		Match:                 expandComputeRegionNetworkFirewallPolicyRuleMatch(d.Get("match")),
		Priority:              dcl.Int64(int64(d.Get("priority").(int))),
		Description:           dcl.String(d.Get("description").(string)),
		Disabled:              dcl.Bool(d.Get("disabled").(bool)),
		EnableLogging:         dcl.Bool(d.Get("enable_logging").(bool)),
		Project:               dcl.String(project),
		Location:              dcl.String(region),
		RuleName:              dcl.String(d.Get("rule_name").(string)),
		TargetSecureTags:      expandComputeRegionNetworkFirewallPolicyRuleTargetSecureTagsArray(d.Get("target_secure_tags")),
		TargetServiceAccounts: tpgdclresource.ExpandStringArray(d.Get("target_service_accounts")),
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
	res, err := client.ApplyNetworkFirewallPolicyRule(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error updating NetworkFirewallPolicyRule: %s", err)
	}

	log.Printf("[DEBUG] Finished creating NetworkFirewallPolicyRule %q: %#v", d.Id(), res)

	return resourceComputeRegionNetworkFirewallPolicyRuleRead(d, meta)
}

func resourceComputeRegionNetworkFirewallPolicyRuleDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	region, err := tpgresource.GetRegion(d, config)
	if err != nil {
		return err
	}

	obj := &compute.NetworkFirewallPolicyRule{
		Action:                dcl.String(d.Get("action").(string)),
		Direction:             compute.NetworkFirewallPolicyRuleDirectionEnumRef(d.Get("direction").(string)),
		FirewallPolicy:        dcl.String(d.Get("firewall_policy").(string)),
		Match:                 expandComputeRegionNetworkFirewallPolicyRuleMatch(d.Get("match")),
		Priority:              dcl.Int64(int64(d.Get("priority").(int))),
		Description:           dcl.String(d.Get("description").(string)),
		Disabled:              dcl.Bool(d.Get("disabled").(bool)),
		EnableLogging:         dcl.Bool(d.Get("enable_logging").(bool)),
		Project:               dcl.String(project),
		Location:              dcl.String(region),
		RuleName:              dcl.String(d.Get("rule_name").(string)),
		TargetSecureTags:      expandComputeRegionNetworkFirewallPolicyRuleTargetSecureTagsArray(d.Get("target_secure_tags")),
		TargetServiceAccounts: tpgdclresource.ExpandStringArray(d.Get("target_service_accounts")),
	}

	log.Printf("[DEBUG] Deleting NetworkFirewallPolicyRule %q", d.Id())
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
	if err := client.DeleteNetworkFirewallPolicyRule(context.Background(), obj); err != nil {
		return fmt.Errorf("Error deleting NetworkFirewallPolicyRule: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting NetworkFirewallPolicyRule %q", d.Id())
	return nil
}

func resourceComputeRegionNetworkFirewallPolicyRuleImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)

	if err := tpgresource.ParseImportId([]string{
		"projects/(?P<project>[^/]+)/regions/(?P<region>[^/]+)/firewallPolicies/(?P<firewall_policy>[^/]+)/(?P<priority>[^/]+)",
		"(?P<project>[^/]+)/(?P<region>[^/]+)/(?P<firewall_policy>[^/]+)/(?P<priority>[^/]+)",
		"(?P<region>[^/]+)/(?P<firewall_policy>[^/]+)/(?P<priority>[^/]+)",
		"(?P<firewall_policy>[^/]+)/(?P<priority>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVarsForId(d, config, "projects/{{project}}/regions/{{region}}/firewallPolicies/{{firewall_policy}}/{{priority}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func expandComputeRegionNetworkFirewallPolicyRuleMatch(o interface{}) *compute.NetworkFirewallPolicyRuleMatch {
	if o == nil {
		return compute.EmptyNetworkFirewallPolicyRuleMatch
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return compute.EmptyNetworkFirewallPolicyRuleMatch
	}
	obj := objArr[0].(map[string]interface{})
	return &compute.NetworkFirewallPolicyRuleMatch{
		Layer4Configs:           expandComputeRegionNetworkFirewallPolicyRuleMatchLayer4ConfigsArray(obj["layer4_configs"]),
		DestAddressGroups:       tpgdclresource.ExpandStringArray(obj["dest_address_groups"]),
		DestFqdns:               tpgdclresource.ExpandStringArray(obj["dest_fqdns"]),
		DestIPRanges:            tpgdclresource.ExpandStringArray(obj["dest_ip_ranges"]),
		DestRegionCodes:         tpgdclresource.ExpandStringArray(obj["dest_region_codes"]),
		DestThreatIntelligences: tpgdclresource.ExpandStringArray(obj["dest_threat_intelligences"]),
		SrcAddressGroups:        tpgdclresource.ExpandStringArray(obj["src_address_groups"]),
		SrcFqdns:                tpgdclresource.ExpandStringArray(obj["src_fqdns"]),
		SrcIPRanges:             tpgdclresource.ExpandStringArray(obj["src_ip_ranges"]),
		SrcRegionCodes:          tpgdclresource.ExpandStringArray(obj["src_region_codes"]),
		SrcSecureTags:           expandComputeRegionNetworkFirewallPolicyRuleMatchSrcSecureTagsArray(obj["src_secure_tags"]),
		SrcThreatIntelligences:  tpgdclresource.ExpandStringArray(obj["src_threat_intelligences"]),
	}
}

func flattenComputeRegionNetworkFirewallPolicyRuleMatch(obj *compute.NetworkFirewallPolicyRuleMatch) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"layer4_configs":            flattenComputeRegionNetworkFirewallPolicyRuleMatchLayer4ConfigsArray(obj.Layer4Configs),
		"dest_address_groups":       obj.DestAddressGroups,
		"dest_fqdns":                obj.DestFqdns,
		"dest_ip_ranges":            obj.DestIPRanges,
		"dest_region_codes":         obj.DestRegionCodes,
		"dest_threat_intelligences": obj.DestThreatIntelligences,
		"src_address_groups":        obj.SrcAddressGroups,
		"src_fqdns":                 obj.SrcFqdns,
		"src_ip_ranges":             obj.SrcIPRanges,
		"src_region_codes":          obj.SrcRegionCodes,
		"src_secure_tags":           flattenComputeRegionNetworkFirewallPolicyRuleMatchSrcSecureTagsArray(obj.SrcSecureTags),
		"src_threat_intelligences":  obj.SrcThreatIntelligences,
	}

	return []interface{}{transformed}

}
func expandComputeRegionNetworkFirewallPolicyRuleMatchLayer4ConfigsArray(o interface{}) []compute.NetworkFirewallPolicyRuleMatchLayer4Configs {
	if o == nil {
		return make([]compute.NetworkFirewallPolicyRuleMatchLayer4Configs, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make([]compute.NetworkFirewallPolicyRuleMatchLayer4Configs, 0)
	}

	items := make([]compute.NetworkFirewallPolicyRuleMatchLayer4Configs, 0, len(objs))
	for _, item := range objs {
		i := expandComputeRegionNetworkFirewallPolicyRuleMatchLayer4Configs(item)
		items = append(items, *i)
	}

	return items
}

func expandComputeRegionNetworkFirewallPolicyRuleMatchLayer4Configs(o interface{}) *compute.NetworkFirewallPolicyRuleMatchLayer4Configs {
	if o == nil {
		return compute.EmptyNetworkFirewallPolicyRuleMatchLayer4Configs
	}

	obj := o.(map[string]interface{})
	return &compute.NetworkFirewallPolicyRuleMatchLayer4Configs{
		IPProtocol: dcl.String(obj["ip_protocol"].(string)),
		Ports:      tpgdclresource.ExpandStringArray(obj["ports"]),
	}
}

func flattenComputeRegionNetworkFirewallPolicyRuleMatchLayer4ConfigsArray(objs []compute.NetworkFirewallPolicyRuleMatchLayer4Configs) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenComputeRegionNetworkFirewallPolicyRuleMatchLayer4Configs(&item)
		items = append(items, i)
	}

	return items
}

func flattenComputeRegionNetworkFirewallPolicyRuleMatchLayer4Configs(obj *compute.NetworkFirewallPolicyRuleMatchLayer4Configs) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"ip_protocol": obj.IPProtocol,
		"ports":       obj.Ports,
	}

	return transformed

}
func expandComputeRegionNetworkFirewallPolicyRuleMatchSrcSecureTagsArray(o interface{}) []compute.NetworkFirewallPolicyRuleMatchSrcSecureTags {
	if o == nil {
		return make([]compute.NetworkFirewallPolicyRuleMatchSrcSecureTags, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make([]compute.NetworkFirewallPolicyRuleMatchSrcSecureTags, 0)
	}

	items := make([]compute.NetworkFirewallPolicyRuleMatchSrcSecureTags, 0, len(objs))
	for _, item := range objs {
		i := expandComputeRegionNetworkFirewallPolicyRuleMatchSrcSecureTags(item)
		items = append(items, *i)
	}

	return items
}

func expandComputeRegionNetworkFirewallPolicyRuleMatchSrcSecureTags(o interface{}) *compute.NetworkFirewallPolicyRuleMatchSrcSecureTags {
	if o == nil {
		return compute.EmptyNetworkFirewallPolicyRuleMatchSrcSecureTags
	}

	obj := o.(map[string]interface{})
	return &compute.NetworkFirewallPolicyRuleMatchSrcSecureTags{
		Name: dcl.String(obj["name"].(string)),
	}
}

func flattenComputeRegionNetworkFirewallPolicyRuleMatchSrcSecureTagsArray(objs []compute.NetworkFirewallPolicyRuleMatchSrcSecureTags) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenComputeRegionNetworkFirewallPolicyRuleMatchSrcSecureTags(&item)
		items = append(items, i)
	}

	return items
}

func flattenComputeRegionNetworkFirewallPolicyRuleMatchSrcSecureTags(obj *compute.NetworkFirewallPolicyRuleMatchSrcSecureTags) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"name":  obj.Name,
		"state": obj.State,
	}

	return transformed

}
func expandComputeRegionNetworkFirewallPolicyRuleTargetSecureTagsArray(o interface{}) []compute.NetworkFirewallPolicyRuleTargetSecureTags {
	if o == nil {
		return make([]compute.NetworkFirewallPolicyRuleTargetSecureTags, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make([]compute.NetworkFirewallPolicyRuleTargetSecureTags, 0)
	}

	items := make([]compute.NetworkFirewallPolicyRuleTargetSecureTags, 0, len(objs))
	for _, item := range objs {
		i := expandComputeRegionNetworkFirewallPolicyRuleTargetSecureTags(item)
		items = append(items, *i)
	}

	return items
}

func expandComputeRegionNetworkFirewallPolicyRuleTargetSecureTags(o interface{}) *compute.NetworkFirewallPolicyRuleTargetSecureTags {
	if o == nil {
		return compute.EmptyNetworkFirewallPolicyRuleTargetSecureTags
	}

	obj := o.(map[string]interface{})
	return &compute.NetworkFirewallPolicyRuleTargetSecureTags{
		Name: dcl.String(obj["name"].(string)),
	}
}

func flattenComputeRegionNetworkFirewallPolicyRuleTargetSecureTagsArray(objs []compute.NetworkFirewallPolicyRuleTargetSecureTags) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenComputeRegionNetworkFirewallPolicyRuleTargetSecureTags(&item)
		items = append(items, i)
	}

	return items
}

func flattenComputeRegionNetworkFirewallPolicyRuleTargetSecureTags(obj *compute.NetworkFirewallPolicyRuleTargetSecureTags) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"name":  obj.Name,
		"state": obj.State,
	}

	return transformed

}
