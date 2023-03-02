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

func ResourceComputeGlobalForwardingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeGlobalForwardingRuleCreate,
		Read:   resourceComputeGlobalForwardingRuleRead,
		Update: resourceComputeGlobalForwardingRuleUpdate,
		Delete: resourceComputeGlobalForwardingRuleDelete,

		Importer: &schema.ResourceImporter{
			State: resourceComputeGlobalForwardingRuleImport,
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
				Description: "Name of the resource; provided by the client when the resource is created. The name must be 1-63 characters long, and comply with [RFC1035](https://www.ietf.org/rfc/rfc1035.txt). Specifically, the name must be 1-63 characters long and match the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` which means the first character must be a lowercase letter, and all following characters must be a dash, lowercase letter, or digit, except the last character, which cannot be a dash.",
			},

			"target": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: compareSelfLinkRelativePaths,
				Description:      "The URL of the target resource to receive the matched traffic. For regional forwarding rules, this target must live in the same region as the forwarding rule. For global forwarding rules, this target must be a global load balancing resource. The forwarded traffic must be of a type appropriate to the target object. For `INTERNAL_SELF_MANAGED` load balancing, only `targetHttpProxy` is valid, not `targetHttpsProxy`.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "An optional description of this resource. Provide this property when you create the resource.",
			},

			"ip_address": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: internalIpDiffSuppress,
				Description:      "IP address that this forwarding rule serves. When a client sends traffic to this IP address, the forwarding rule directs the traffic to the target that you specify in the forwarding rule. If you don't specify a reserved IP address, an ephemeral IP address is assigned. Methods for specifying an IP address: * IPv4 dotted decimal, as in `100.1.2.3` * Full URL, as in `https://www.googleapis.com/compute/v1/projects/project_id/regions/region/addresses/address-name` * Partial URL or by name, as in: * `projects/project_id/regions/region/addresses/address-name` * `regions/region/addresses/address-name` * `global/addresses/address-name` * `address-name` The loadBalancingScheme and the forwarding rule's target determine the type of IP address that you can use. For detailed information, refer to [IP address specifications](/load-balancing/docs/forwarding-rule-concepts#ip_address_specifications).",
			},

			"ip_protocol": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: caseDiffSuppress,
				Description:      "The IP protocol to which this rule applies. For protocol forwarding, valid options are `TCP`, `UDP`, `ESP`, `AH`, `SCTP` or `ICMP`. For Internal TCP/UDP Load Balancing, the load balancing scheme is `INTERNAL`, and one of `TCP` or `UDP` are valid. For Traffic Director, the load balancing scheme is `INTERNAL_SELF_MANAGED`, and only `TCP`is valid. For Internal HTTP(S) Load Balancing, the load balancing scheme is `INTERNAL_MANAGED`, and only `TCP` is valid. For HTTP(S), SSL Proxy, and TCP Proxy Load Balancing, the load balancing scheme is `EXTERNAL` and only `TCP` is valid. For Network TCP/UDP Load Balancing, the load balancing scheme is `EXTERNAL`, and one of `TCP` or `UDP` is valid.",
			},

			"ip_version": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The IP Version that will be used by this forwarding rule. Valid options are `IPV4` or `IPV6`. This can only be specified for an external global forwarding rule. Possible values: UNSPECIFIED_VERSION, IPV4, IPV6",
			},

			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels to apply to this rule.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"load_balancing_scheme": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies the forwarding rule type.\n\n*   `EXTERNAL` is used for:\n    *   Classic Cloud VPN gateways\n    *   Protocol forwarding to VMs from an external IP address\n    *   The following load balancers: HTTP(S), SSL Proxy, TCP Proxy, and Network TCP/UDP\n*   `INTERNAL` is used for:\n    *   Protocol forwarding to VMs from an internal IP address\n    *   Internal TCP/UDP load balancers\n*   `INTERNAL_MANAGED` is used for:\n    *   Internal HTTP(S) load balancers\n*   `INTERNAL_SELF_MANAGED` is used for:\n    *   Traffic Director\n*   `EXTERNAL_MANAGED` is used for:\n    *   Global external HTTP(S) load balancers \n\nFor more information about forwarding rules, refer to [Forwarding rule concepts](/load-balancing/docs/forwarding-rule-concepts). Possible values: INVALID, INTERNAL, INTERNAL_MANAGED, INTERNAL_SELF_MANAGED, EXTERNAL, EXTERNAL_MANAGED",
				Default:     "EXTERNAL",
			},

			"metadata_filters": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Opaque filter criteria used by Loadbalancer to restrict routing configuration to a limited set of [xDS](https://github.com/envoyproxy/data-plane-api/blob/master/XDS_PROTOCOL.md) compliant clients. In their xDS requests to Loadbalancer, xDS clients present [node metadata](https://github.com/envoyproxy/data-plane-api/search?q=%22message+Node%22+in%3A%2Fenvoy%2Fapi%2Fv2%2Fcore%2Fbase.proto&). If a match takes place, the relevant configuration is made available to those proxies. Otherwise, all the resources (e.g. `TargetHttpProxy`, `UrlMap`) referenced by the `ForwardingRule` will not be visible to those proxies.\n\nFor each `metadataFilter` in this list, if its `filterMatchCriteria` is set to MATCH_ANY, at least one of the `filterLabel`s must match the corresponding label provided in the metadata. If its `filterMatchCriteria` is set to MATCH_ALL, then all of its `filterLabel`s must match with corresponding labels provided in the metadata.\n\n`metadataFilters` specified here will be applifed before those specified in the `UrlMap` that this `ForwardingRule` references.\n\n`metadataFilters` only applies to Loadbalancers that have their loadBalancingScheme set to `INTERNAL_SELF_MANAGED`.",
				Elem:        ComputeGlobalForwardingRuleMetadataFilterSchema(),
			},

			"network": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      "This field is not used for external load balancing. For `INTERNAL` and `INTERNAL_SELF_MANAGED` load balancing, this field identifies the network that the load balanced IP should belong to for this Forwarding Rule. If this field is not specified, the default network will be used.",
			},

			"port_range": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: portRangeDiffSuppress,
				Description:      "When the load balancing scheme is `EXTERNAL`, `INTERNAL_SELF_MANAGED` and `INTERNAL_MANAGED`, you can specify a `port_range`. Use with a forwarding rule that points to a target proxy or a target pool. Do not use with a forwarding rule that points to a backend service. This field is used along with the `target` field for TargetHttpProxy, TargetHttpsProxy, TargetSslProxy, TargetTcpProxy, TargetVpnGateway, TargetPool, TargetInstance. Applicable only when `IPProtocol` is `TCP`, `UDP`, or `SCTP`, only packets addressed to ports in the specified range will be forwarded to `target`. Forwarding rules with the same `[IPAddress, IPProtocol]` pair must have disjoint port ranges. Some types of forwarding target have constraints on the acceptable ports:\n\n*   TargetHttpProxy: 80, 8080\n*   TargetHttpsProxy: 443\n*   TargetTcpProxy: 25, 43, 110, 143, 195, 443, 465, 587, 700, 993, 995, 1688, 1883, 5222\n*   TargetSslProxy: 25, 43, 110, 143, 195, 443, 465, 587, 700, 993, 995, 1688, 1883, 5222\n*   TargetVpnGateway: 500, 4500\n\n@pattern: d+(?:-d+)?",
			},

			"project": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      "The project this resource belongs in.",
			},

			"label_fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Used internally during label updates.",
			},

			"psc_connection_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The PSC connection id of the PSC Forwarding Rule.",
			},

			"psc_connection_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The PSC connection status of the PSC Forwarding Rule. Possible values: STATUS_UNSPECIFIED, PENDING, ACCEPTED, REJECTED, CLOSED",
			},

			"self_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "[Output Only] Server-defined URL for the resource.",
			},
		},
	}
}

func ComputeGlobalForwardingRuleMetadataFilterSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"filter_labels": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "The list of label value pairs that must match labels in the provided metadata based on `filterMatchCriteria`\n\nThis list must not be empty and can have at the most 64 entries.",
				MaxItems:    64,
				MinItems:    1,
				Elem:        ComputeGlobalForwardingRuleMetadataFilterFilterLabelSchema(),
			},

			"filter_match_criteria": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies how individual `filterLabel` matches within the list of `filterLabel`s contribute towards the overall `metadataFilter` match.\n\nSupported values are:\n\n*   MATCH_ANY: At least one of the `filterLabels` must have a matching label in the provided metadata.\n*   MATCH_ALL: All `filterLabels` must have matching labels in the provided metadata. Possible values: NOT_SET, MATCH_ALL, MATCH_ANY",
			},
		},
	}
}

func ComputeGlobalForwardingRuleMetadataFilterFilterLabelSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of metadata label.\n\nThe name can have a maximum length of 1024 characters and must be at least 1 character long.",
			},

			"value": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The value of the label must match the specified value.\n\nvalue can have a maximum length of 1024 characters.",
			},
		},
	}
}

func resourceComputeGlobalForwardingRuleCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	obj := &compute.ForwardingRule{
		Name:                dcl.String(d.Get("name").(string)),
		Target:              dcl.String(d.Get("target").(string)),
		Description:         dcl.String(d.Get("description").(string)),
		IPAddress:           dcl.StringOrNil(d.Get("ip_address").(string)),
		IPProtocol:          compute.ForwardingRuleIPProtocolEnumRef(d.Get("ip_protocol").(string)),
		IPVersion:           compute.ForwardingRuleIPVersionEnumRef(d.Get("ip_version").(string)),
		Labels:              checkStringMap(d.Get("labels")),
		LoadBalancingScheme: compute.ForwardingRuleLoadBalancingSchemeEnumRef(d.Get("load_balancing_scheme").(string)),
		MetadataFilter:      expandComputeGlobalForwardingRuleMetadataFilterArray(d.Get("metadata_filters")),
		Network:             dcl.StringOrNil(d.Get("network").(string)),
		PortRange:           dcl.String(d.Get("port_range").(string)),
		Project:             dcl.String(project),
	}

	id, err := obj.ID()
	if err != nil {
		return fmt.Errorf("error constructing id: %s", err)
	}
	d.SetId(id)
	directive := CreateDirective
	userAgent, err := generateUserAgentString(d, config.UserAgent)
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
	res, err := client.ApplyForwardingRule(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating ForwardingRule: %s", err)
	}

	log.Printf("[DEBUG] Finished creating ForwardingRule %q: %#v", d.Id(), res)

	return resourceComputeGlobalForwardingRuleRead(d, meta)
}

func resourceComputeGlobalForwardingRuleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	obj := &compute.ForwardingRule{
		Name:                dcl.String(d.Get("name").(string)),
		Target:              dcl.String(d.Get("target").(string)),
		Description:         dcl.String(d.Get("description").(string)),
		IPAddress:           dcl.StringOrNil(d.Get("ip_address").(string)),
		IPProtocol:          compute.ForwardingRuleIPProtocolEnumRef(d.Get("ip_protocol").(string)),
		IPVersion:           compute.ForwardingRuleIPVersionEnumRef(d.Get("ip_version").(string)),
		Labels:              checkStringMap(d.Get("labels")),
		LoadBalancingScheme: compute.ForwardingRuleLoadBalancingSchemeEnumRef(d.Get("load_balancing_scheme").(string)),
		MetadataFilter:      expandComputeGlobalForwardingRuleMetadataFilterArray(d.Get("metadata_filters")),
		Network:             dcl.StringOrNil(d.Get("network").(string)),
		PortRange:           dcl.String(d.Get("port_range").(string)),
		Project:             dcl.String(project),
	}

	userAgent, err := generateUserAgentString(d, config.UserAgent)
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
	res, err := client.GetForwardingRule(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("ComputeGlobalForwardingRule %q", d.Id())
		return handleNotFoundDCLError(err, d, resourceName)
	}

	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}
	if err = d.Set("target", res.Target); err != nil {
		return fmt.Errorf("error setting target in state: %s", err)
	}
	if err = d.Set("description", res.Description); err != nil {
		return fmt.Errorf("error setting description in state: %s", err)
	}
	if err = d.Set("ip_address", res.IPAddress); err != nil {
		return fmt.Errorf("error setting ip_address in state: %s", err)
	}
	if err = d.Set("ip_protocol", res.IPProtocol); err != nil {
		return fmt.Errorf("error setting ip_protocol in state: %s", err)
	}
	if err = d.Set("ip_version", res.IPVersion); err != nil {
		return fmt.Errorf("error setting ip_version in state: %s", err)
	}
	if err = d.Set("labels", res.Labels); err != nil {
		return fmt.Errorf("error setting labels in state: %s", err)
	}
	if err = d.Set("load_balancing_scheme", res.LoadBalancingScheme); err != nil {
		return fmt.Errorf("error setting load_balancing_scheme in state: %s", err)
	}
	if err = d.Set("metadata_filters", flattenComputeGlobalForwardingRuleMetadataFilterArray(res.MetadataFilter)); err != nil {
		return fmt.Errorf("error setting metadata_filters in state: %s", err)
	}
	if err = d.Set("network", res.Network); err != nil {
		return fmt.Errorf("error setting network in state: %s", err)
	}
	if err = d.Set("port_range", res.PortRange); err != nil {
		return fmt.Errorf("error setting port_range in state: %s", err)
	}
	if err = d.Set("project", res.Project); err != nil {
		return fmt.Errorf("error setting project in state: %s", err)
	}
	if err = d.Set("label_fingerprint", res.LabelFingerprint); err != nil {
		return fmt.Errorf("error setting label_fingerprint in state: %s", err)
	}
	if err = d.Set("psc_connection_id", res.PscConnectionId); err != nil {
		return fmt.Errorf("error setting psc_connection_id in state: %s", err)
	}
	if err = d.Set("psc_connection_status", res.PscConnectionStatus); err != nil {
		return fmt.Errorf("error setting psc_connection_status in state: %s", err)
	}
	if err = d.Set("self_link", res.SelfLink); err != nil {
		return fmt.Errorf("error setting self_link in state: %s", err)
	}

	return nil
}
func resourceComputeGlobalForwardingRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	obj := &compute.ForwardingRule{
		Name:                dcl.String(d.Get("name").(string)),
		Target:              dcl.String(d.Get("target").(string)),
		Description:         dcl.String(d.Get("description").(string)),
		IPAddress:           dcl.StringOrNil(d.Get("ip_address").(string)),
		IPProtocol:          compute.ForwardingRuleIPProtocolEnumRef(d.Get("ip_protocol").(string)),
		IPVersion:           compute.ForwardingRuleIPVersionEnumRef(d.Get("ip_version").(string)),
		Labels:              checkStringMap(d.Get("labels")),
		LoadBalancingScheme: compute.ForwardingRuleLoadBalancingSchemeEnumRef(d.Get("load_balancing_scheme").(string)),
		MetadataFilter:      expandComputeGlobalForwardingRuleMetadataFilterArray(d.Get("metadata_filters")),
		Network:             dcl.StringOrNil(d.Get("network").(string)),
		PortRange:           dcl.String(d.Get("port_range").(string)),
		Project:             dcl.String(project),
	}
	directive := UpdateDirective
	userAgent, err := generateUserAgentString(d, config.UserAgent)
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
	res, err := client.ApplyForwardingRule(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error updating ForwardingRule: %s", err)
	}

	log.Printf("[DEBUG] Finished creating ForwardingRule %q: %#v", d.Id(), res)

	return resourceComputeGlobalForwardingRuleRead(d, meta)
}

func resourceComputeGlobalForwardingRuleDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	obj := &compute.ForwardingRule{
		Name:                dcl.String(d.Get("name").(string)),
		Target:              dcl.String(d.Get("target").(string)),
		Description:         dcl.String(d.Get("description").(string)),
		IPAddress:           dcl.StringOrNil(d.Get("ip_address").(string)),
		IPProtocol:          compute.ForwardingRuleIPProtocolEnumRef(d.Get("ip_protocol").(string)),
		IPVersion:           compute.ForwardingRuleIPVersionEnumRef(d.Get("ip_version").(string)),
		Labels:              checkStringMap(d.Get("labels")),
		LoadBalancingScheme: compute.ForwardingRuleLoadBalancingSchemeEnumRef(d.Get("load_balancing_scheme").(string)),
		MetadataFilter:      expandComputeGlobalForwardingRuleMetadataFilterArray(d.Get("metadata_filters")),
		Network:             dcl.StringOrNil(d.Get("network").(string)),
		PortRange:           dcl.String(d.Get("port_range").(string)),
		Project:             dcl.String(project),
	}

	log.Printf("[DEBUG] Deleting ForwardingRule %q", d.Id())
	userAgent, err := generateUserAgentString(d, config.UserAgent)
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
	if err := client.DeleteForwardingRule(context.Background(), obj); err != nil {
		return fmt.Errorf("Error deleting ForwardingRule: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting ForwardingRule %q", d.Id())
	return nil
}

func resourceComputeGlobalForwardingRuleImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/global/forwardingRules/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<name>[^/]+)",
		"(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVarsForId(d, config, "projects/{{project}}/global/forwardingRules/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func expandComputeGlobalForwardingRuleMetadataFilterArray(o interface{}) []compute.ForwardingRuleMetadataFilter {
	if o == nil {
		return make([]compute.ForwardingRuleMetadataFilter, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make([]compute.ForwardingRuleMetadataFilter, 0)
	}

	items := make([]compute.ForwardingRuleMetadataFilter, 0, len(objs))
	for _, item := range objs {
		i := expandComputeGlobalForwardingRuleMetadataFilter(item)
		items = append(items, *i)
	}

	return items
}

func expandComputeGlobalForwardingRuleMetadataFilter(o interface{}) *compute.ForwardingRuleMetadataFilter {
	if o == nil {
		return compute.EmptyForwardingRuleMetadataFilter
	}

	obj := o.(map[string]interface{})
	return &compute.ForwardingRuleMetadataFilter{
		FilterLabel:         expandComputeGlobalForwardingRuleMetadataFilterFilterLabelArray(obj["filter_labels"]),
		FilterMatchCriteria: compute.ForwardingRuleMetadataFilterFilterMatchCriteriaEnumRef(obj["filter_match_criteria"].(string)),
	}
}

func flattenComputeGlobalForwardingRuleMetadataFilterArray(objs []compute.ForwardingRuleMetadataFilter) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenComputeGlobalForwardingRuleMetadataFilter(&item)
		items = append(items, i)
	}

	return items
}

func flattenComputeGlobalForwardingRuleMetadataFilter(obj *compute.ForwardingRuleMetadataFilter) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"filter_labels":         flattenComputeGlobalForwardingRuleMetadataFilterFilterLabelArray(obj.FilterLabel),
		"filter_match_criteria": obj.FilterMatchCriteria,
	}

	return transformed

}
func expandComputeGlobalForwardingRuleMetadataFilterFilterLabelArray(o interface{}) []compute.ForwardingRuleMetadataFilterFilterLabel {
	if o == nil {
		return make([]compute.ForwardingRuleMetadataFilterFilterLabel, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make([]compute.ForwardingRuleMetadataFilterFilterLabel, 0)
	}

	items := make([]compute.ForwardingRuleMetadataFilterFilterLabel, 0, len(objs))
	for _, item := range objs {
		i := expandComputeGlobalForwardingRuleMetadataFilterFilterLabel(item)
		items = append(items, *i)
	}

	return items
}

func expandComputeGlobalForwardingRuleMetadataFilterFilterLabel(o interface{}) *compute.ForwardingRuleMetadataFilterFilterLabel {
	if o == nil {
		return compute.EmptyForwardingRuleMetadataFilterFilterLabel
	}

	obj := o.(map[string]interface{})
	return &compute.ForwardingRuleMetadataFilterFilterLabel{
		Name:  dcl.String(obj["name"].(string)),
		Value: dcl.String(obj["value"].(string)),
	}
}

func flattenComputeGlobalForwardingRuleMetadataFilterFilterLabelArray(objs []compute.ForwardingRuleMetadataFilterFilterLabel) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenComputeGlobalForwardingRuleMetadataFilterFilterLabel(&item)
		items = append(items, i)
	}

	return items
}

func flattenComputeGlobalForwardingRuleMetadataFilterFilterLabel(obj *compute.ForwardingRuleMetadataFilterFilterLabel) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"name":  obj.Name,
		"value": obj.Value,
	}

	return transformed

}
