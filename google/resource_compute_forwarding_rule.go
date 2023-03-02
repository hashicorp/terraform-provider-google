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

func ResourceComputeForwardingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeForwardingRuleCreate,
		Read:   resourceComputeForwardingRuleRead,
		Update: resourceComputeForwardingRuleUpdate,
		Delete: resourceComputeForwardingRuleDelete,

		Importer: &schema.ResourceImporter{
			State: resourceComputeForwardingRuleImport,
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

			"all_ports": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "This field is used along with the `backend_service` field for internal load balancing or with the `target` field for internal TargetInstance. This field cannot be used with `port` or `portRange` fields. When the load balancing scheme is `INTERNAL` and protocol is TCP/UDP, specify this field to allow packets addressed to any ports will be forwarded to the backends configured with this forwarding rule.",
			},

			"allow_global_access": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "This field is used along with the `backend_service` field for internal load balancing or with the `target` field for internal TargetInstance. If the field is set to `TRUE`, clients can access ILB from all regions. Otherwise only allows access from clients in the same region as the internal load balancer.",
			},

			"backend_service": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      "This field is only used for `INTERNAL` load balancing. For internal load balancing, this field identifies the BackendService resource to receive the matched traffic.",
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

			"is_mirroring_collector": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Indicates whether or not this load balancer can be used as a collector for packet mirroring. To prevent mirroring loops, instances behind this load balancer will not have their traffic mirrored even if a `PacketMirroring` rule applies to them. This can only be set to true for load balancers that have their `loadBalancingScheme` set to `INTERNAL`.",
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

			"network": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      "This field is not used for external load balancing. For `INTERNAL` and `INTERNAL_SELF_MANAGED` load balancing, this field identifies the network that the load balanced IP should belong to for this Forwarding Rule. If this field is not specified, the default network will be used.",
			},

			"network_tier": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "This signifies the networking tier used for configuring this load balancer and can only take the following values: `PREMIUM`, `STANDARD`. For regional ForwardingRule, the valid values are `PREMIUM` and `STANDARD`. For GlobalForwardingRule, the valid value is `PREMIUM`. If this field is not specified, it is assumed to be `PREMIUM`. If `IPAddress` is specified, this value must be equal to the networkTier of the Address.",
			},

			"port_range": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: portRangeDiffSuppress,
				Description:      "When the load balancing scheme is `EXTERNAL`, `INTERNAL_SELF_MANAGED` and `INTERNAL_MANAGED`, you can specify a `port_range`. Use with a forwarding rule that points to a target proxy or a target pool. Do not use with a forwarding rule that points to a backend service. This field is used along with the `target` field for TargetHttpProxy, TargetHttpsProxy, TargetSslProxy, TargetTcpProxy, TargetVpnGateway, TargetPool, TargetInstance. Applicable only when `IPProtocol` is `TCP`, `UDP`, or `SCTP`, only packets addressed to ports in the specified range will be forwarded to `target`. Forwarding rules with the same `[IPAddress, IPProtocol]` pair must have disjoint port ranges. Some types of forwarding target have constraints on the acceptable ports:\n\n*   TargetHttpProxy: 80, 8080\n*   TargetHttpsProxy: 443\n*   TargetTcpProxy: 25, 43, 110, 143, 195, 443, 465, 587, 700, 993, 995, 1688, 1883, 5222\n*   TargetSslProxy: 25, 43, 110, 143, 195, 443, 465, 587, 700, 993, 995, 1688, 1883, 5222\n*   TargetVpnGateway: 500, 4500\n\n@pattern: d+(?:-d+)?",
			},

			"ports": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Description: "This field is used along with the `backend_service` field for internal load balancing. When the load balancing scheme is `INTERNAL`, a list of ports can be configured, for example, ['80'], ['8000','9000']. Only packets addressed to these ports are forwarded to the backends configured with the forwarding rule. If the forwarding rule's loadBalancingScheme is INTERNAL, you can specify ports in one of the following ways: * A list of up to five ports, which can be non-contiguous * Keyword `ALL`, which causes the forwarding rule to forward traffic on any port of the forwarding rule's protocol. @pattern: d+(?:-d+)? For more information, refer to [Port specifications](/load-balancing/docs/forwarding-rule-concepts#port_specifications).",
				MaxItems:    5,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
			},

			"project": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      "The project this resource belongs in.",
			},

			"region": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      "The location of this resource.",
			},

			"service_directory_registrations": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Service Directory resources to register this forwarding rule with. Currently, only supports a single Service Directory resource.",
				Elem:        ComputeForwardingRuleServiceDirectoryRegistrationsSchema(),
			},

			"service_label": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "An optional prefix to the service name for this Forwarding Rule. If specified, the prefix is the first label of the fully qualified service name. The label must be 1-63 characters long, and comply with [RFC1035](https://www.ietf.org/rfc/rfc1035.txt). Specifically, the label must be 1-63 characters long and match the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` which means the first character must be a lowercase letter, and all following characters must be a dash, lowercase letter, or digit, except the last character, which cannot be a dash. This field is only used for internal load balancing.",
				ValidateFunc: validateGCEName,
			},

			"subnetwork": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      "This field is only used for `INTERNAL` load balancing. For internal load balancing, this field identifies the subnetwork that the load balanced IP should belong to for this Forwarding Rule. If the network specified is in auto subnet mode, this field is optional. However, if the network is in custom subnet mode, a subnetwork must be specified.",
			},

			"target": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: compareSelfLinkRelativePaths,
				Description:      "The URL of the target resource to receive the matched traffic. For regional forwarding rules, this target must live in the same region as the forwarding rule. For global forwarding rules, this target must be a global load balancing resource. The forwarded traffic must be of a type appropriate to the target object. For `INTERNAL_SELF_MANAGED` load balancing, only `targetHttpProxy` is valid, not `targetHttpsProxy`.",
			},

			"creation_timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "[Output Only] Creation timestamp in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) text format.",
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

			"service_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "[Output Only] The internal fully qualified service name for this Forwarding Rule. This field is only used for internal load balancing.",
			},
		},
	}
}

func ComputeForwardingRuleServiceDirectoryRegistrationsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Service Directory namespace to register the forwarding rule under.",
			},

			"service": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Service Directory service to register the forwarding rule under.",
			},
		},
	}
}

func resourceComputeForwardingRuleCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	obj := &compute.ForwardingRule{
		Name:                          dcl.String(d.Get("name").(string)),
		AllPorts:                      dcl.Bool(d.Get("all_ports").(bool)),
		AllowGlobalAccess:             dcl.Bool(d.Get("allow_global_access").(bool)),
		BackendService:                dcl.String(d.Get("backend_service").(string)),
		Description:                   dcl.String(d.Get("description").(string)),
		IPAddress:                     dcl.StringOrNil(d.Get("ip_address").(string)),
		IPProtocol:                    compute.ForwardingRuleIPProtocolEnumRef(d.Get("ip_protocol").(string)),
		IsMirroringCollector:          dcl.Bool(d.Get("is_mirroring_collector").(bool)),
		Labels:                        checkStringMap(d.Get("labels")),
		LoadBalancingScheme:           compute.ForwardingRuleLoadBalancingSchemeEnumRef(d.Get("load_balancing_scheme").(string)),
		Network:                       dcl.StringOrNil(d.Get("network").(string)),
		NetworkTier:                   compute.ForwardingRuleNetworkTierEnumRef(d.Get("network_tier").(string)),
		PortRange:                     dcl.String(d.Get("port_range").(string)),
		Ports:                         expandStringArray(d.Get("ports")),
		Project:                       dcl.String(project),
		Location:                      dcl.String(region),
		ServiceDirectoryRegistrations: expandComputeForwardingRuleServiceDirectoryRegistrationsArray(d.Get("service_directory_registrations")),
		ServiceLabel:                  dcl.String(d.Get("service_label").(string)),
		Subnetwork:                    dcl.StringOrNil(d.Get("subnetwork").(string)),
		Target:                        dcl.String(d.Get("target").(string)),
	}

	id, err := replaceVarsForId(d, config, "projects/{{project}}/regions/{{region}}/forwardingRules/{{name}}")
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

	return resourceComputeForwardingRuleRead(d, meta)
}

func resourceComputeForwardingRuleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	obj := &compute.ForwardingRule{
		Name:                          dcl.String(d.Get("name").(string)),
		AllPorts:                      dcl.Bool(d.Get("all_ports").(bool)),
		AllowGlobalAccess:             dcl.Bool(d.Get("allow_global_access").(bool)),
		BackendService:                dcl.String(d.Get("backend_service").(string)),
		Description:                   dcl.String(d.Get("description").(string)),
		IPAddress:                     dcl.StringOrNil(d.Get("ip_address").(string)),
		IPProtocol:                    compute.ForwardingRuleIPProtocolEnumRef(d.Get("ip_protocol").(string)),
		IsMirroringCollector:          dcl.Bool(d.Get("is_mirroring_collector").(bool)),
		Labels:                        checkStringMap(d.Get("labels")),
		LoadBalancingScheme:           compute.ForwardingRuleLoadBalancingSchemeEnumRef(d.Get("load_balancing_scheme").(string)),
		Network:                       dcl.StringOrNil(d.Get("network").(string)),
		NetworkTier:                   compute.ForwardingRuleNetworkTierEnumRef(d.Get("network_tier").(string)),
		PortRange:                     dcl.String(d.Get("port_range").(string)),
		Ports:                         expandStringArray(d.Get("ports")),
		Project:                       dcl.String(project),
		Location:                      dcl.String(region),
		ServiceDirectoryRegistrations: expandComputeForwardingRuleServiceDirectoryRegistrationsArray(d.Get("service_directory_registrations")),
		ServiceLabel:                  dcl.String(d.Get("service_label").(string)),
		Subnetwork:                    dcl.StringOrNil(d.Get("subnetwork").(string)),
		Target:                        dcl.String(d.Get("target").(string)),
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
		resourceName := fmt.Sprintf("ComputeForwardingRule %q", d.Id())
		return handleNotFoundDCLError(err, d, resourceName)
	}

	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}
	if err = d.Set("all_ports", res.AllPorts); err != nil {
		return fmt.Errorf("error setting all_ports in state: %s", err)
	}
	if err = d.Set("allow_global_access", res.AllowGlobalAccess); err != nil {
		return fmt.Errorf("error setting allow_global_access in state: %s", err)
	}
	if err = d.Set("backend_service", res.BackendService); err != nil {
		return fmt.Errorf("error setting backend_service in state: %s", err)
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
	if err = d.Set("is_mirroring_collector", res.IsMirroringCollector); err != nil {
		return fmt.Errorf("error setting is_mirroring_collector in state: %s", err)
	}
	if err = d.Set("labels", res.Labels); err != nil {
		return fmt.Errorf("error setting labels in state: %s", err)
	}
	if err = d.Set("load_balancing_scheme", res.LoadBalancingScheme); err != nil {
		return fmt.Errorf("error setting load_balancing_scheme in state: %s", err)
	}
	if err = d.Set("network", res.Network); err != nil {
		return fmt.Errorf("error setting network in state: %s", err)
	}
	if err = d.Set("network_tier", res.NetworkTier); err != nil {
		return fmt.Errorf("error setting network_tier in state: %s", err)
	}
	if err = d.Set("port_range", res.PortRange); err != nil {
		return fmt.Errorf("error setting port_range in state: %s", err)
	}
	if err = d.Set("ports", res.Ports); err != nil {
		return fmt.Errorf("error setting ports in state: %s", err)
	}
	if err = d.Set("project", res.Project); err != nil {
		return fmt.Errorf("error setting project in state: %s", err)
	}
	if err = d.Set("region", res.Location); err != nil {
		return fmt.Errorf("error setting region in state: %s", err)
	}
	if err = d.Set("service_directory_registrations", flattenComputeForwardingRuleServiceDirectoryRegistrationsArray(res.ServiceDirectoryRegistrations)); err != nil {
		return fmt.Errorf("error setting service_directory_registrations in state: %s", err)
	}
	if err = d.Set("service_label", res.ServiceLabel); err != nil {
		return fmt.Errorf("error setting service_label in state: %s", err)
	}
	if err = d.Set("subnetwork", res.Subnetwork); err != nil {
		return fmt.Errorf("error setting subnetwork in state: %s", err)
	}
	if err = d.Set("target", res.Target); err != nil {
		return fmt.Errorf("error setting target in state: %s", err)
	}
	if err = d.Set("creation_timestamp", res.CreationTimestamp); err != nil {
		return fmt.Errorf("error setting creation_timestamp in state: %s", err)
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
	if err = d.Set("service_name", res.ServiceName); err != nil {
		return fmt.Errorf("error setting service_name in state: %s", err)
	}

	return nil
}
func resourceComputeForwardingRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	obj := &compute.ForwardingRule{
		Name:                          dcl.String(d.Get("name").(string)),
		AllPorts:                      dcl.Bool(d.Get("all_ports").(bool)),
		AllowGlobalAccess:             dcl.Bool(d.Get("allow_global_access").(bool)),
		BackendService:                dcl.String(d.Get("backend_service").(string)),
		Description:                   dcl.String(d.Get("description").(string)),
		IPAddress:                     dcl.StringOrNil(d.Get("ip_address").(string)),
		IPProtocol:                    compute.ForwardingRuleIPProtocolEnumRef(d.Get("ip_protocol").(string)),
		IsMirroringCollector:          dcl.Bool(d.Get("is_mirroring_collector").(bool)),
		Labels:                        checkStringMap(d.Get("labels")),
		LoadBalancingScheme:           compute.ForwardingRuleLoadBalancingSchemeEnumRef(d.Get("load_balancing_scheme").(string)),
		Network:                       dcl.StringOrNil(d.Get("network").(string)),
		NetworkTier:                   compute.ForwardingRuleNetworkTierEnumRef(d.Get("network_tier").(string)),
		PortRange:                     dcl.String(d.Get("port_range").(string)),
		Ports:                         expandStringArray(d.Get("ports")),
		Project:                       dcl.String(project),
		Location:                      dcl.String(region),
		ServiceDirectoryRegistrations: expandComputeForwardingRuleServiceDirectoryRegistrationsArray(d.Get("service_directory_registrations")),
		ServiceLabel:                  dcl.String(d.Get("service_label").(string)),
		Subnetwork:                    dcl.StringOrNil(d.Get("subnetwork").(string)),
		Target:                        dcl.String(d.Get("target").(string)),
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

	return resourceComputeForwardingRuleRead(d, meta)
}

func resourceComputeForwardingRuleDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	obj := &compute.ForwardingRule{
		Name:                          dcl.String(d.Get("name").(string)),
		AllPorts:                      dcl.Bool(d.Get("all_ports").(bool)),
		AllowGlobalAccess:             dcl.Bool(d.Get("allow_global_access").(bool)),
		BackendService:                dcl.String(d.Get("backend_service").(string)),
		Description:                   dcl.String(d.Get("description").(string)),
		IPAddress:                     dcl.StringOrNil(d.Get("ip_address").(string)),
		IPProtocol:                    compute.ForwardingRuleIPProtocolEnumRef(d.Get("ip_protocol").(string)),
		IsMirroringCollector:          dcl.Bool(d.Get("is_mirroring_collector").(bool)),
		Labels:                        checkStringMap(d.Get("labels")),
		LoadBalancingScheme:           compute.ForwardingRuleLoadBalancingSchemeEnumRef(d.Get("load_balancing_scheme").(string)),
		Network:                       dcl.StringOrNil(d.Get("network").(string)),
		NetworkTier:                   compute.ForwardingRuleNetworkTierEnumRef(d.Get("network_tier").(string)),
		PortRange:                     dcl.String(d.Get("port_range").(string)),
		Ports:                         expandStringArray(d.Get("ports")),
		Project:                       dcl.String(project),
		Location:                      dcl.String(region),
		ServiceDirectoryRegistrations: expandComputeForwardingRuleServiceDirectoryRegistrationsArray(d.Get("service_directory_registrations")),
		ServiceLabel:                  dcl.String(d.Get("service_label").(string)),
		Subnetwork:                    dcl.StringOrNil(d.Get("subnetwork").(string)),
		Target:                        dcl.String(d.Get("target").(string)),
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

func resourceComputeForwardingRuleImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/regions/(?P<region>[^/]+)/forwardingRules/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<region>[^/]+)/(?P<name>[^/]+)",
		"(?P<region>[^/]+)/(?P<name>[^/]+)",
		"(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVarsForId(d, config, "projects/{{project}}/regions/{{region}}/forwardingRules/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func expandComputeForwardingRuleServiceDirectoryRegistrationsArray(o interface{}) []compute.ForwardingRuleServiceDirectoryRegistrations {
	if o == nil {
		return nil
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return nil
	}

	items := make([]compute.ForwardingRuleServiceDirectoryRegistrations, 0, len(objs))
	for _, item := range objs {
		i := expandComputeForwardingRuleServiceDirectoryRegistrations(item)
		items = append(items, *i)
	}

	return items
}

func expandComputeForwardingRuleServiceDirectoryRegistrations(o interface{}) *compute.ForwardingRuleServiceDirectoryRegistrations {
	if o == nil {
		return nil
	}

	obj := o.(map[string]interface{})
	return &compute.ForwardingRuleServiceDirectoryRegistrations{
		Namespace: dcl.StringOrNil(obj["namespace"].(string)),
		Service:   dcl.String(obj["service"].(string)),
	}
}

func flattenComputeForwardingRuleServiceDirectoryRegistrationsArray(objs []compute.ForwardingRuleServiceDirectoryRegistrations) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenComputeForwardingRuleServiceDirectoryRegistrations(&item)
		items = append(items, i)
	}

	return items
}

func flattenComputeForwardingRuleServiceDirectoryRegistrations(obj *compute.ForwardingRuleServiceDirectoryRegistrations) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"namespace": obj.Namespace,
		"service":   obj.Service,
	}

	return transformed

}
