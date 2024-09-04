// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"fmt"
	"log"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"
)

func ipv6RepresentationDiffSuppress(_, old, new string, d *schema.ResourceData) bool {
	//Diff suppress any equal IPV6 address in different representations
	//An IPV6 address can have long or short representations
	//E.g 2001:0cb0:0000:0000:0fc0:0000:0000:0abc, after compression:
	//A) 2001:0cb0::0fc0:0000:0000:0abc (Omit groups of all zeros)
	//B) 2001:cb0:0:0:fc0::abc (Omit leading zeros)
	//C) 2001:cb0::fc0:0:0:abc (Combining A and B)
	//The GCP API follows rule B) for normalzation

	oldIp := net.ParseIP(old)
	newIp := net.ParseIP(new)
	return oldIp.Equal(newIp)
}

func ResourceComputeRouterBgpPeer() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeRouterBgpPeerCreate,
		Read:   resourceComputeRouterBgpPeerRead,
		Update: resourceComputeRouterBgpPeerUpdate,
		Delete: resourceComputeRouterBgpPeerDelete,

		Importer: &schema.ResourceImporter{
			State: resourceComputeRouterBgpPeerImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			tpgresource.DefaultProviderProject,
		),

		Schema: map[string]*schema.Schema{
			"interface": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Name of the interface the BGP peer is associated with.`,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidateRFC1035Name(2, 63),
				Description: `Name of this BGP peer. The name must be 1-63 characters long,
and comply with RFC1035. Specifically, the name must be 1-63 characters
long and match the regular expression '[a-z]([-a-z0-9]*[a-z0-9])?' which
means the first character must be a lowercase letter, and all
following characters must be a dash, lowercase letter, or digit,
except the last character, which cannot be a dash.`,
			},
			"peer_asn": {
				Type:     schema.TypeInt,
				Required: true,
				Description: `Peer BGP Autonomous System Number (ASN).
Each BGP interface may use a different value.`,
			},
			"router": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      `The name of the Cloud Router in which this BgpPeer will be configured.`,
			},
			"advertise_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: verify.ValidateEnum([]string{"DEFAULT", "CUSTOM", ""}),
				Description: `User-specified flag to indicate which mode to use for advertisement.
Valid values of this enum field are: 'DEFAULT', 'CUSTOM' Default value: "DEFAULT" Possible values: ["DEFAULT", "CUSTOM"]`,
				Default: "DEFAULT",
			},
			"advertised_groups": {
				Type:     schema.TypeList,
				Optional: true,
				Description: `User-specified list of prefix groups to advertise in custom
mode, which currently supports the following option:

* 'ALL_SUBNETS': Advertises all of the router's own VPC subnets.
This excludes any routes learned for subnets that use VPC Network
Peering.


Note that this field can only be populated if advertiseMode is 'CUSTOM'
and overrides the list defined for the router (in the "bgp" message).
These groups are advertised in addition to any specified prefixes.
Leave this field blank to advertise no custom groups.`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"advertised_ip_ranges": {
				Type:     schema.TypeList,
				Optional: true,
				Description: `User-specified list of individual IP ranges to advertise in
custom mode. This field can only be populated if advertiseMode
is 'CUSTOM' and is advertised to all peers of the router. These IP
ranges will be advertised in addition to any specified groups.
Leave this field blank to advertise no custom IP ranges.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"range": {
							Type:     schema.TypeString,
							Required: true,
							Description: `The IP range to advertise. The value must be a
CIDR-formatted string.`,
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `User-specified description for the IP range.`,
						},
					},
				},
			},
			"advertised_route_priority": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: `The priority of routes advertised to this BGP peer.
Where there is more than one matching route of maximum
length, the routes with the lowest priority value win.`,
			},
			"custom_learned_ip_ranges": {
				Type:     schema.TypeList,
				Optional: true,
				Description: `The custom learned route IP address range. Must be a valid CIDR-formatted prefix. If an 
IP address is provided without a subnet mask, it is interpreted as, for IPv4, a /32 singular IP address range, and, for IPv6, /128.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"range": {
							Type:     schema.TypeString,
							Required: true,
							Description: `The IP range to advertise. The value must be a
CIDR-formatted string.`,
						},
					},
				},
			},
			"custom_learned_route_priority": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: `The user-defined custom learned route priority for a BGP session.
This value is applied to all custom learned route ranges for the session. You can choose a value
from 0 to 65335. If you don't provide a value, Google Cloud assigns a priority of 100 to the ranges.`,
			},

			"bfd": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: `BFD configuration for the BGP peering.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"session_initialization_mode": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: verify.ValidateEnum([]string{"ACTIVE", "DISABLED", "PASSIVE"}),
							Description: `The BFD session initialization mode for this BGP peer.
If set to 'ACTIVE', the Cloud Router will initiate the BFD session
for this BGP peer. If set to 'PASSIVE', the Cloud Router will wait
for the peer router to initiate the BFD session for this BGP peer.
If set to 'DISABLED', BFD is disabled for this BGP peer. Possible values: ["ACTIVE", "DISABLED", "PASSIVE"]`,
						},
						"min_receive_interval": {
							Type:     schema.TypeInt,
							Optional: true,
							Description: `The minimum interval, in milliseconds, between BFD control packets
received from the peer router. The actual value is negotiated
between the two routers and is equal to the greater of this value
and the transmit interval of the other router. If set, this value
must be between 1000 and 30000.`,
							Default: 1000,
						},
						"min_transmit_interval": {
							Type:     schema.TypeInt,
							Optional: true,
							Description: `The minimum interval, in milliseconds, between BFD control packets
transmitted to the peer router. The actual value is negotiated
between the two routers and is equal to the greater of this value
and the corresponding receive interval of the other router. If set,
this value must be between 1000 and 30000.`,
							Default: 1000,
						},
						"multiplier": {
							Type:     schema.TypeInt,
							Optional: true,
							Description: `The number of consecutive BFD packets that must be missed before
BFD declares that a peer is unavailable. If set, the value must
be a value between 5 and 16.`,
							Default: 5,
						},
					},
				},
			},
			"enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Description: `The status of the BGP peer connection. If set to false, any active session
with the peer is terminated and all associated routing information is removed.
If set to true, the peer connection can be established with routing information.
The default is true.`,
				Default: true,
			},
			"enable_ipv6": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Enable IPv6 traffic over BGP Peer. If not specified, it is disabled by default.`,
				Default:     false,
			},
			"enable_ipv4": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Enable IPv4 traffic over BGP Peer. It is enabled by default if the peerIpAddress is version 4.`,
				Computed:    true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				Description: `IP address of the interface inside Google Cloud Platform.
Only IPv4 is supported.`,
			},
			"ipv6_nexthop_address": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ValidateFunc:     verify.ValidateIpAddress,
				DiffSuppressFunc: ipv6RepresentationDiffSuppress,
				Description: `IPv6 address of the interface inside Google Cloud Platform.
The address must be in the range 2600:2d00:0:2::/64 or 2600:2d00:0:3::/64.
If you do not specify the next hop addresses, Google Cloud automatically
assigns unused addresses from the 2600:2d00:0:2::/64 or 2600:2d00:0:3::/64 range for you.`,
			},
			"ipv4_nexthop_address": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: verify.ValidateIpAddress,
				Description:  `IPv4 address of the interface inside Google Cloud Platform.`,
			},
			"peer_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				Description: `IP address of the BGP interface outside Google Cloud Platform.
Only IPv4 is supported. Required if 'ip_address' is set.`,
			},
			"peer_ipv6_nexthop_address": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ValidateFunc:     verify.ValidateIpAddress,
				DiffSuppressFunc: ipv6RepresentationDiffSuppress,
				Description: `IPv6 address of the BGP interface outside Google Cloud Platform.
The address must be in the range 2600:2d00:0:2::/64 or 2600:2d00:0:3::/64.
If you do not specify the next hop addresses, Google Cloud automatically
assigns unused addresses from the 2600:2d00:0:2::/64 or 2600:2d00:0:3::/64 range for you.`,
			},
			"peer_ipv4_nexthop_address": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: verify.ValidateIpAddress,
				Description:  `IPv4 address of the BGP interface outside Google Cloud Platform.`,
			},
			"region": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description: `Region where the router and BgpPeer reside.
If it is not provided, the provider region is used.`,
			},
			"router_appliance_instance": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description: `The URI of the VM instance that is used as third-party router appliances
such as Next Gen Firewalls, Virtual Routers, or Router Appliances.
The VM instance must be located in zones contained in the same region as
this Cloud Router. The VM instance is the peer side of the BGP session.`,
			},
			"management_type": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The resource that configures and manages this BGP peer.

* 'MANAGED_BY_USER' is the default value and can be managed by
you or other users
* 'MANAGED_BY_ATTACHMENT' is a BGP peer that is configured and
managed by Cloud Interconnect, specifically by an
InterconnectAttachment of type PARTNER. Google automatically
creates, updates, and deletes this type of BGP peer when the
PARTNER InterconnectAttachment is created, updated,
or deleted.`,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"md5_authentication_key": {
				Type:     schema.TypeList,
				Optional: true,
				Description: `Present if MD5 authentication is enabled for the peering. Must be the name
of one of the entries in the Router.md5_authentication_keys. The field must comply with RFC1035.`,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							Description: `[REQUIRED] Name used to identify the key.
Must be unique within a router. Must be referenced by exactly one bgpPeer. Must comply with RFC1035.`,
						},
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Value of the key.`,
							Sensitive:   true,
						},
					},
				},
			},
		},
		UseJSONNumber: true,
	}
}

func resourceComputeRouterBgpPeerCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	nameProp, err := expandNestedComputeRouterBgpPeerName(d.Get("name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	interfaceNameProp, err := expandNestedComputeRouterBgpPeerInterface(d.Get("interface"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("interface"); !tpgresource.IsEmptyValue(reflect.ValueOf(interfaceNameProp)) && (ok || !reflect.DeepEqual(v, interfaceNameProp)) {
		obj["interfaceName"] = interfaceNameProp
	}
	ipAddressProp, err := expandNestedComputeRouterBgpPeerIpAddress(d.Get("ip_address"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("ip_address"); !tpgresource.IsEmptyValue(reflect.ValueOf(ipAddressProp)) && (ok || !reflect.DeepEqual(v, ipAddressProp)) {
		obj["ipAddress"] = ipAddressProp
	}
	peerIpAddressProp, err := expandNestedComputeRouterBgpPeerPeerIpAddress(d.Get("peer_ip_address"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("peer_ip_address"); !tpgresource.IsEmptyValue(reflect.ValueOf(peerIpAddressProp)) && (ok || !reflect.DeepEqual(v, peerIpAddressProp)) {
		obj["peerIpAddress"] = peerIpAddressProp
	}
	peerAsnProp, err := expandNestedComputeRouterBgpPeerPeerAsn(d.Get("peer_asn"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("peer_asn"); !tpgresource.IsEmptyValue(reflect.ValueOf(peerAsnProp)) && (ok || !reflect.DeepEqual(v, peerAsnProp)) {
		obj["peerAsn"] = peerAsnProp
	}
	advertisedRoutePriorityProp, err := expandNestedComputeRouterBgpPeerAdvertisedRoutePriority(d.Get("advertised_route_priority"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOk("advertised_route_priority"); ok || !reflect.DeepEqual(v, advertisedRoutePriorityProp) {
		obj["advertisedRoutePriority"] = advertisedRoutePriorityProp
	}
	advertiseModeProp, err := expandNestedComputeRouterBgpPeerAdvertiseMode(d.Get("advertise_mode"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("advertise_mode"); !tpgresource.IsEmptyValue(reflect.ValueOf(advertiseModeProp)) && (ok || !reflect.DeepEqual(v, advertiseModeProp)) {
		obj["advertiseMode"] = advertiseModeProp
	}
	advertisedGroupsProp, err := expandNestedComputeRouterBgpPeerAdvertisedGroups(d.Get("advertised_groups"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("advertised_groups"); ok || !reflect.DeepEqual(v, advertisedGroupsProp) {
		obj["advertisedGroups"] = advertisedGroupsProp
	}
	advertisedIpRangesProp, err := expandNestedComputeRouterBgpPeerAdvertisedIpRanges(d.Get("advertised_ip_ranges"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("advertised_ip_ranges"); ok || !reflect.DeepEqual(v, advertisedIpRangesProp) {
		obj["advertisedIpRanges"] = advertisedIpRangesProp
	}
	customLearnedIpRangesProp, err := expandNestedComputeRouterBgpPeerCustomLearnedIpRanges(d.Get("custom_learned_ip_ranges"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("custom_learned_ip_ranges"); ok || !reflect.DeepEqual(v, customLearnedIpRangesProp) {
		obj["customLearnedIpRanges"] = customLearnedIpRangesProp
	}
	customLearnedRoutePriorityProp, err := expandNestedComputeRouterBgpPeerCustomLearnedRoutePriority(d.Get("custom_learned_route_priority"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("custom_learned_route_priority"); ok || !reflect.DeepEqual(v, customLearnedRoutePriorityProp) {
		obj["customLearnedRoutePriority"] = customLearnedRoutePriorityProp
	}
	bfdProp, err := expandNestedComputeRouterBgpPeerBfd(d.Get("bfd"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("bfd"); !tpgresource.IsEmptyValue(reflect.ValueOf(bfdProp)) && (ok || !reflect.DeepEqual(v, bfdProp)) {
		obj["bfd"] = bfdProp
	}
	enableProp, err := expandNestedComputeRouterBgpPeerEnable(d.Get("enable"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("enable"); ok || !reflect.DeepEqual(v, enableProp) {
		obj["enable"] = enableProp
	}
	routerApplianceInstanceProp, err := expandNestedComputeRouterBgpPeerRouterApplianceInstance(d.Get("router_appliance_instance"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("router_appliance_instance"); !tpgresource.IsEmptyValue(reflect.ValueOf(routerApplianceInstanceProp)) && (ok || !reflect.DeepEqual(v, routerApplianceInstanceProp)) {
		obj["routerApplianceInstance"] = routerApplianceInstanceProp
	}
	enableIpv6Prop, err := expandNestedComputeRouterBgpPeerEnableIpv6(d.Get("enable_ipv6"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("enable_ipv6"); ok || !reflect.DeepEqual(v, enableIpv6Prop) {
		obj["enableIpv6"] = enableIpv6Prop
	}
	enableIpv4Prop, err := expandNestedComputeRouterBgpPeerEnableIpv4(d.Get("enable_ipv4"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("enable_ipv4"); ok || !reflect.DeepEqual(v, enableIpv4Prop) {
		obj["enableIpv4"] = enableIpv4Prop
	}
	ipv4NexthopAddressProp, err := expandNestedComputeRouterBgpPeerIpv4NexthopAddress(d.Get("ipv4_nexthop_address"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("ipv4_nexthop_address"); !tpgresource.IsEmptyValue(reflect.ValueOf(ipv4NexthopAddressProp)) && (ok || !reflect.DeepEqual(v, ipv4NexthopAddressProp)) {
		obj["ipv4NexthopAddress"] = ipv4NexthopAddressProp
	}
	peerIpv4NexthopAddressProp, err := expandNestedComputeRouterBgpPeerPeerIpv4NexthopAddress(d.Get("peer_ipv4_nexthop_address"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("peer_ipv6_nexthop_address"); !tpgresource.IsEmptyValue(reflect.ValueOf(peerIpv4NexthopAddressProp)) && (ok || !reflect.DeepEqual(v, peerIpv4NexthopAddressProp)) {
		obj["peerIpv4NexthopAddress"] = peerIpv4NexthopAddressProp
	}
	ipv6NexthopAddressProp, err := expandNestedComputeRouterBgpPeerIpv6NexthopAddress(d.Get("ipv6_nexthop_address"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("ipv6_nexthop_address"); !tpgresource.IsEmptyValue(reflect.ValueOf(ipv6NexthopAddressProp)) && (ok || !reflect.DeepEqual(v, ipv6NexthopAddressProp)) {
		obj["ipv6NexthopAddress"] = ipv6NexthopAddressProp
	}
	peerIpv6NexthopAddressProp, err := expandNestedComputeRouterBgpPeerPeerIpv6NexthopAddress(d.Get("peer_ipv6_nexthop_address"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("peer_ipv6_nexthop_address"); !tpgresource.IsEmptyValue(reflect.ValueOf(peerIpv6NexthopAddressProp)) && (ok || !reflect.DeepEqual(v, peerIpv6NexthopAddressProp)) {
		obj["peerIpv6NexthopAddress"] = peerIpv6NexthopAddressProp
	}
	md5AuthenticationKeyProp, err := expandNestedComputeRouterBgpPeerMd5AuthenticationKey(d.Get("md5_authentication_key"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("md5_authentication_key"); !tpgresource.IsEmptyValue(reflect.ValueOf(md5AuthenticationKeyProp)) && (ok || !reflect.DeepEqual(v, md5AuthenticationKeyProp)) {
		/*some manual handling is required here as the parent cloud router object has a different layout for keyName and keyValue.
		bgpPeer blocks in cloud router only specify the keyName to be used and the cloudRouter object has another block called
		md5AuthenticationKeys which is an array which specify all the keys (name and value). The constraint here is that a key must
		be used by exactly one bgpPeer to be considered valid.
		*/
		md5AuthenticationKeyName := md5AuthenticationKeyProp.(map[string]interface{})["name"]
		obj["md5AuthenticationKeyName"] = md5AuthenticationKeyName
		obj["md5AuthenticationKey"] = md5AuthenticationKeyProp
	}

	lockName, err := tpgresource.ReplaceVars(d, config, "router/{{region}}/{{router}}")
	if err != nil {
		return err
	}
	transport_tpg.MutexStore.Lock(lockName)
	defer transport_tpg.MutexStore.Unlock(lockName)

	url, err := tpgresource.ReplaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/regions/{{region}}/routers/{{router}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new RouterBgpPeer: %#v", obj)

	obj, err = resourceComputeRouterBgpPeerPatchCreateEncoder(d, meta, obj)
	if err != nil {
		return err
	}
	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for RouterBgpPeer: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "PATCH",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutCreate),
	})
	if err != nil {
		return fmt.Errorf("Error creating RouterBgpPeer: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/regions/{{region}}/routers/{{router}}/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = ComputeOperationWaitTime(
		config, res, project, "Creating RouterBgpPeer", userAgent,
		d.Timeout(schema.TimeoutCreate))

	if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error waiting to create RouterBgpPeer: %s", err)
	}

	log.Printf("[DEBUG] Finished creating RouterBgpPeer %q: %#v", d.Id(), res)

	err = d.Set("md5_authentication_key", []interface{}{md5AuthenticationKeyProp})
	if err != nil {
		return err
	}

	return resourceComputeRouterBgpPeerRead(d, meta)
}

func resourceComputeRouterBgpPeerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/regions/{{region}}/routers/{{router}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for RouterBgpPeer: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("ComputeRouterBgpPeer %q", d.Id()))
	}

	res, err = flattenNestedComputeRouterBgpPeer(d, meta, res)
	if err != nil {
		return err
	}

	if res == nil {
		// Object isn't there any more - remove it from the state.
		log.Printf("[DEBUG] Removing ComputeRouterBgpPeer because it couldn't be matched.")
		d.SetId("")
		return nil
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading RouterBgpPeer: %s", err)
	}

	if err := d.Set("name", flattenNestedComputeRouterBgpPeerName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading RouterBgpPeer: %s", err)
	}
	if err := d.Set("interface", flattenNestedComputeRouterBgpPeerInterface(res["interfaceName"], d, config)); err != nil {
		return fmt.Errorf("Error reading RouterBgpPeer: %s", err)
	}
	if err := d.Set("ip_address", flattenNestedComputeRouterBgpPeerIpAddress(res["ipAddress"], d, config)); err != nil {
		return fmt.Errorf("Error reading RouterBgpPeer: %s", err)
	}
	if err := d.Set("peer_ip_address", flattenNestedComputeRouterBgpPeerPeerIpAddress(res["peerIpAddress"], d, config)); err != nil {
		return fmt.Errorf("Error reading RouterBgpPeer: %s", err)
	}
	if err := d.Set("peer_asn", flattenNestedComputeRouterBgpPeerPeerAsn(res["peerAsn"], d, config)); err != nil {
		return fmt.Errorf("Error reading RouterBgpPeer: %s", err)
	}
	if err := d.Set("advertised_route_priority", flattenNestedComputeRouterBgpPeerAdvertisedRoutePriority(res["advertisedRoutePriority"], d, config)); err != nil {
		return fmt.Errorf("Error reading RouterBgpPeer: %s", err)
	}
	if err := d.Set("advertise_mode", flattenNestedComputeRouterBgpPeerAdvertiseMode(res["advertiseMode"], d, config)); err != nil {
		return fmt.Errorf("Error reading RouterBgpPeer: %s", err)
	}
	if err := d.Set("advertised_groups", flattenNestedComputeRouterBgpPeerAdvertisedGroups(res["advertisedGroups"], d, config)); err != nil {
		return fmt.Errorf("Error reading RouterBgpPeer: %s", err)
	}
	if err := d.Set("advertised_ip_ranges", flattenNestedComputeRouterBgpPeerAdvertisedIpRanges(res["advertisedIpRanges"], d, config)); err != nil {
		return fmt.Errorf("Error reading RouterBgpPeer: %s", err)
	}
	if err := d.Set("custom_learned_ip_ranges", flattenNestedComputeRouterBgpPeerCustomLearnedIpRanges(res["customLearnedIpRanges"], d, config)); err != nil {
		return fmt.Errorf("Error reading RouterBgpPeer: %s", err)
	}
	if err := d.Set("custom_learned_route_priority", flattenNestedComputeRouterBgpPeerCustomLearnedRoutePriority(res["customLearnedRoutePriority"], d, config)); err != nil {
		return fmt.Errorf("Error reading RouterBgpPeer: %s", err)
	}
	if err := d.Set("management_type", flattenNestedComputeRouterBgpPeerManagementType(res["managementType"], d, config)); err != nil {
		return fmt.Errorf("Error reading RouterBgpPeer: %s", err)
	}
	if err := d.Set("bfd", flattenNestedComputeRouterBgpPeerBfd(res["bfd"], d, config)); err != nil {
		return fmt.Errorf("Error reading RouterBgpPeer: %s", err)
	}
	if err := d.Set("enable", flattenNestedComputeRouterBgpPeerEnable(res["enable"], d, config)); err != nil {
		return fmt.Errorf("Error reading RouterBgpPeer: %s", err)
	}
	if err := d.Set("router_appliance_instance", flattenNestedComputeRouterBgpPeerRouterApplianceInstance(res["routerApplianceInstance"], d, config)); err != nil {
		return fmt.Errorf("Error reading RouterBgpPeer: %s", err)
	}
	if err := d.Set("enable_ipv6", flattenNestedComputeRouterBgpPeerEnableIpv6(res["enableIpv6"], d, config)); err != nil {
		return fmt.Errorf("Error reading RouterBgpPeer: %s", err)
	}
	if err := d.Set("enable_ipv4", flattenNestedComputeRouterBgpPeerEnableIpv4(res["enableIpv4"], d, config)); err != nil {
		return fmt.Errorf("Error reading RouterBgpPeer: %s", err)
	}
	if err := d.Set("ipv4_nexthop_address", flattenNestedComputeRouterBgpPeerIpv4NexthopAddress(res["ipv4NexthopAddress"], d, config)); err != nil {
		return fmt.Errorf("Error reading RouterBgpPeer: %s", err)
	}
	if err := d.Set("peer_ipv4_nexthop_address", flattenNestedComputeRouterBgpPeerPeerIpv4NexthopAddress(res["peerIpv4NexthopAddress"], d, config)); err != nil {
		return fmt.Errorf("Error reading RouterBgpPeer: %s", err)
	}
	if err := d.Set("ipv6_nexthop_address", flattenNestedComputeRouterBgpPeerIpv6NexthopAddress(res["ipv6NexthopAddress"], d, config)); err != nil {
		return fmt.Errorf("Error reading RouterBgpPeer: %s", err)
	}
	if err := d.Set("peer_ipv6_nexthop_address", flattenNestedComputeRouterBgpPeerPeerIpv6NexthopAddress(res["peerIpv6NexthopAddress"], d, config)); err != nil {
		return fmt.Errorf("Error reading RouterBgpPeer: %s", err)
	}
	if err := d.Set("md5_authentication_key", flattenNestedComputeRouterBgpPeerMd5AuthenticationKey(res["md5AuthenticationKeyName"], d, config)); err != nil {
		return fmt.Errorf("Error reading RouterBgpPeer: %s", err)
	}

	return nil
}

func resourceComputeRouterBgpPeerUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for RouterBgpPeer: %s", err)
	}
	billingProject = project

	obj := make(map[string]interface{})
	ipAddressProp, err := expandNestedComputeRouterBgpPeerIpAddress(d.Get("ip_address"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("ip_address"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, ipAddressProp)) {
		obj["ipAddress"] = ipAddressProp
	}
	peerIpAddressProp, err := expandNestedComputeRouterBgpPeerPeerIpAddress(d.Get("peer_ip_address"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("peer_ip_address"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, peerIpAddressProp)) {
		obj["peerIpAddress"] = peerIpAddressProp
	}
	peerAsnProp, err := expandNestedComputeRouterBgpPeerPeerAsn(d.Get("peer_asn"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("peer_asn"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, peerAsnProp)) {
		obj["peerAsn"] = peerAsnProp
	}
	advertisedRoutePriorityProp, err := expandNestedComputeRouterBgpPeerAdvertisedRoutePriority(d.Get("advertised_route_priority"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOk("advertised_route_priority"); ok || !reflect.DeepEqual(v, advertisedRoutePriorityProp) {
		obj["advertisedRoutePriority"] = advertisedRoutePriorityProp
	}
	advertiseModeProp, err := expandNestedComputeRouterBgpPeerAdvertiseMode(d.Get("advertise_mode"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("advertise_mode"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, advertiseModeProp)) {
		obj["advertiseMode"] = advertiseModeProp
	}
	advertisedGroupsProp, err := expandNestedComputeRouterBgpPeerAdvertisedGroups(d.Get("advertised_groups"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("advertised_groups"); ok || !reflect.DeepEqual(v, advertisedGroupsProp) {
		obj["advertisedGroups"] = advertisedGroupsProp
	}
	advertisedIpRangesProp, err := expandNestedComputeRouterBgpPeerAdvertisedIpRanges(d.Get("advertised_ip_ranges"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("advertised_ip_ranges"); ok || !reflect.DeepEqual(v, advertisedIpRangesProp) {
		obj["advertisedIpRanges"] = advertisedIpRangesProp
	}
	customLearnedIpRangesProp, err := expandNestedComputeRouterBgpPeerCustomLearnedIpRanges(d.Get("custom_learned_ip_ranges"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("custom_learned_ip_ranges"); ok || !reflect.DeepEqual(v, customLearnedIpRangesProp) {
		obj["customLearnedIpRanges"] = customLearnedIpRangesProp
	}
	customLearnedRoutePriorityProp, err := expandNestedComputeRouterBgpPeerCustomLearnedRoutePriority(d.Get("custom_learned_route_priority"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("custom_learned_route_priority"); ok || !reflect.DeepEqual(v, customLearnedRoutePriorityProp) {
		obj["customLearnedRoutePriority"] = customLearnedRoutePriorityProp
	}
	bfdProp, err := expandNestedComputeRouterBgpPeerBfd(d.Get("bfd"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("bfd"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, bfdProp)) {
		obj["bfd"] = bfdProp
	}
	enableProp, err := expandNestedComputeRouterBgpPeerEnable(d.Get("enable"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("enable"); ok || !reflect.DeepEqual(v, enableProp) {
		obj["enable"] = enableProp
	}
	routerApplianceInstanceProp, err := expandNestedComputeRouterBgpPeerRouterApplianceInstance(d.Get("router_appliance_instance"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("router_appliance_instance"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, routerApplianceInstanceProp)) {
		obj["routerApplianceInstance"] = routerApplianceInstanceProp
	}
	enableIpv6Prop, err := expandNestedComputeRouterBgpPeerEnableIpv6(d.Get("enable_ipv6"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("enable_ipv6"); ok || !reflect.DeepEqual(v, enableIpv6Prop) {
		obj["enableIpv6"] = enableIpv6Prop
	}
	enableIpv4Prop, err := expandNestedComputeRouterBgpPeerEnableIpv4(d.Get("enable_ipv4"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("enable_ipv4"); ok || !reflect.DeepEqual(v, enableIpv4Prop) {
		obj["enableIpv4"] = enableIpv4Prop
	}
	ipv4NexthopAddressProp, err := expandNestedComputeRouterBgpPeerIpv4NexthopAddress(d.Get("ipv4_nexthop_address"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("ipv4_nexthop_address"); !tpgresource.IsEmptyValue(reflect.ValueOf(ipv4NexthopAddressProp)) && (ok || !reflect.DeepEqual(v, ipv4NexthopAddressProp)) {
		obj["ipv4NexthopAddress"] = ipv4NexthopAddressProp
	}
	peerIpv4NexthopAddressProp, err := expandNestedComputeRouterBgpPeerPeerIpv4NexthopAddress(d.Get("peer_ipv4_nexthop_address"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("peer_ipv4_nexthop_address"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, peerIpv4NexthopAddressProp)) {
		obj["peerIpv4NexthopAddress"] = peerIpv4NexthopAddressProp
	}
	ipv6NexthopAddressProp, err := expandNestedComputeRouterBgpPeerIpv6NexthopAddress(d.Get("ipv6_nexthop_address"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("ipv6_nexthop_address"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, ipv6NexthopAddressProp)) {
		obj["ipv6NexthopAddress"] = ipv6NexthopAddressProp
	}
	peerIpv6NexthopAddressProp, err := expandNestedComputeRouterBgpPeerPeerIpv6NexthopAddress(d.Get("peer_ipv6_nexthop_address"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("peer_ipv6_nexthop_address"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, peerIpv6NexthopAddressProp)) {
		obj["peerIpv6NexthopAddress"] = peerIpv6NexthopAddressProp
	}
	md5AuthenticationKeyProp, err := expandNestedComputeRouterBgpPeerMd5AuthenticationKey(d.Get("md5_authentication_key"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("md5_authentication_key"); !tpgresource.IsEmptyValue(reflect.ValueOf(md5AuthenticationKeyProp)) && (ok || !reflect.DeepEqual(v, md5AuthenticationKeyProp)) {
		/*some manual handling is required here as the parent cloud router object has a different layout for keyName and keyValue.
		bgpPeer blocks in cloud router only specify the keyName to be used and the cloudRouter object has another block called
		md5AuthenticationKeys which is an array which specify all the keys (name and value). The constraint here is that a key must
		be used by exactly one bgpPeer to be considered valid.
		*/
		md5AuthenticationKeyName := md5AuthenticationKeyProp.(map[string]interface{})["name"]
		obj["md5AuthenticationKeyName"] = md5AuthenticationKeyName
		obj["md5AuthenticationKey"] = md5AuthenticationKeyProp
	}

	lockName, err := tpgresource.ReplaceVars(d, config, "router/{{region}}/{{router}}")
	if err != nil {
		return err
	}
	transport_tpg.MutexStore.Lock(lockName)
	defer transport_tpg.MutexStore.Unlock(lockName)

	url, err := tpgresource.ReplaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/regions/{{region}}/routers/{{router}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating RouterBgpPeer %q: %#v", d.Id(), obj)

	obj, err = resourceComputeRouterBgpPeerPatchUpdateEncoder(d, meta, obj)
	if err != nil {
		return err
	}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "PATCH",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutUpdate),
	})

	if err != nil {
		return fmt.Errorf("Error updating RouterBgpPeer %q: %s", d.Id(), err)
	} else {
		log.Printf("[DEBUG] Finished updating RouterBgpPeer %q: %#v", d.Id(), res)
	}

	err = ComputeOperationWaitTime(
		config, res, project, "Updating RouterBgpPeer", userAgent,
		d.Timeout(schema.TimeoutUpdate))

	if err != nil {
		return err
	}

	err = d.Set("md5_authentication_key", []interface{}{md5AuthenticationKeyProp})
	if err != nil {
		return err
	}

	return resourceComputeRouterBgpPeerRead(d, meta)
}

func resourceComputeRouterBgpPeerDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for RouterBgpPeer: %s", err)
	}
	billingProject = project

	lockName, err := tpgresource.ReplaceVars(d, config, "router/{{region}}/{{router}}")
	if err != nil {
		return err
	}
	transport_tpg.MutexStore.Lock(lockName)
	defer transport_tpg.MutexStore.Unlock(lockName)

	url, err := tpgresource.ReplaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/regions/{{region}}/routers/{{router}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}

	obj, err = resourceComputeRouterBgpPeerPatchDeleteEncoder(d, meta, obj)
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, "RouterBgpPeer")
	}
	log.Printf("[DEBUG] Deleting RouterBgpPeer %q", d.Id())

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "PATCH",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutDelete),
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, "RouterBgpPeer")
	}

	err = ComputeOperationWaitTime(
		config, res, project, "Deleting RouterBgpPeer", userAgent,
		d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting RouterBgpPeer %q: %#v", d.Id(), res)
	return nil
}

func resourceComputeRouterBgpPeerImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"^projects/(?P<project>[^/]+)/regions/(?P<region>[^/]+)/routers/(?P<router>[^/]+)/(?P<name>[^/]+)$",
		"^(?P<project>[^/]+)/(?P<region>[^/]+)/(?P<router>[^/]+)/(?P<name>[^/]+)$",
		"^(?P<region>[^/]+)/(?P<router>[^/]+)/(?P<name>[^/]+)$",
		"^(?P<router>[^/]+)/(?P<name>[^/]+)$",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/regions/{{region}}/routers/{{router}}/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenNestedComputeRouterBgpPeerName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNestedComputeRouterBgpPeerInterface(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNestedComputeRouterBgpPeerIpAddress(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNestedComputeRouterBgpPeerPeerIpAddress(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNestedComputeRouterBgpPeerPeerAsn(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenNestedComputeRouterBgpPeerAdvertisedRoutePriority(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenNestedComputeRouterBgpPeerAdvertiseMode(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil || tpgresource.IsEmptyValue(reflect.ValueOf(v)) {
		return "DEFAULT"
	}

	return v
}

func flattenNestedComputeRouterBgpPeerAdvertisedGroups(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNestedComputeRouterBgpPeerAdvertisedIpRanges(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"range":       flattenNestedComputeRouterBgpPeerAdvertisedIpRangesRange(original["range"], d, config),
			"description": flattenNestedComputeRouterBgpPeerAdvertisedIpRangesDescription(original["description"], d, config),
		})
	}
	return transformed
}
func flattenNestedComputeRouterBgpPeerAdvertisedIpRangesRange(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNestedComputeRouterBgpPeerAdvertisedIpRangesDescription(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}
func flattenNestedComputeRouterBgpPeerCustomLearnedIpRanges(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"range": flattenNestedComputeRouterBgpPeerCustomLearnedIpRangesRange(original["range"], d, config),
		})
	}
	return transformed
}
func flattenNestedComputeRouterBgpPeerCustomLearnedIpRangesRange(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNestedComputeRouterBgpPeerCustomLearnedRoutePriority(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}
func flattenNestedComputeRouterBgpPeerManagementType(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNestedComputeRouterBgpPeerMd5AuthenticationKey(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	originalKeyValue := d.Get("md5_authentication_key").([]interface{})
	transformed := make(map[string]interface{})
	transformed["name"] = v
	//key value is not returned as it is a sensitive field
	if len(originalKeyValue) != 0 {
		transformed["key"] = originalKeyValue[0].(map[string]interface{})["key"]
	}
	return []interface{}{transformed}
}

func flattenNestedComputeRouterBgpPeerBfd(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["session_initialization_mode"] =
		flattenNestedComputeRouterBgpPeerBfdSessionInitializationMode(original["sessionInitializationMode"], d, config)
	transformed["min_transmit_interval"] =
		flattenNestedComputeRouterBgpPeerBfdMinTransmitInterval(original["minTransmitInterval"], d, config)
	transformed["min_receive_interval"] =
		flattenNestedComputeRouterBgpPeerBfdMinReceiveInterval(original["minReceiveInterval"], d, config)
	transformed["multiplier"] =
		flattenNestedComputeRouterBgpPeerBfdMultiplier(original["multiplier"], d, config)
	return []interface{}{transformed}
}

func flattenNestedComputeRouterBgpPeerBfdSessionInitializationMode(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNestedComputeRouterBgpPeerBfdMinTransmitInterval(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenNestedComputeRouterBgpPeerBfdMinReceiveInterval(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenNestedComputeRouterBgpPeerBfdMultiplier(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenNestedComputeRouterBgpPeerEnable(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return true
	}
	b, err := strconv.ParseBool(v.(string))
	if err != nil {
		// If we can't convert it into a bool return value as is and let caller handle it
		return v
	}
	return b
}

func flattenNestedComputeRouterBgpPeerRouterApplianceInstance(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	return tpgresource.ConvertSelfLinkToV1(v.(string))
}

func flattenNestedComputeRouterBgpPeerEnableIpv6(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNestedComputeRouterBgpPeerEnableIpv4(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNestedComputeRouterBgpPeerIpv4NexthopAddress(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNestedComputeRouterBgpPeerPeerIpv4NexthopAddress(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNestedComputeRouterBgpPeerIpv6NexthopAddress(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNestedComputeRouterBgpPeerPeerIpv6NexthopAddress(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandNestedComputeRouterBgpPeerName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNestedComputeRouterBgpPeerInterface(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNestedComputeRouterBgpPeerIpAddress(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNestedComputeRouterBgpPeerPeerIpAddress(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNestedComputeRouterBgpPeerPeerAsn(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNestedComputeRouterBgpPeerAdvertisedRoutePriority(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNestedComputeRouterBgpPeerAdvertiseMode(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNestedComputeRouterBgpPeerAdvertisedGroups(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNestedComputeRouterBgpPeerAdvertisedIpRanges(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedRange, err := expandNestedComputeRouterBgpPeerAdvertisedIpRangesRange(original["range"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedRange); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["range"] = transformedRange
		}

		transformedDescription, err := expandNestedComputeRouterBgpPeerAdvertisedIpRangesDescription(original["description"], d, config)
		if err != nil {
			return nil, err
		} else {
			transformed["description"] = transformedDescription
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandNestedComputeRouterBgpPeerAdvertisedIpRangesRange(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNestedComputeRouterBgpPeerAdvertisedIpRangesDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
func expandNestedComputeRouterBgpPeerCustomLearnedIpRanges(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedRange, err := expandNestedComputeRouterBgpPeerCustomLearnedIpRangesRange(original["range"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedRange); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["range"] = transformedRange
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandNestedComputeRouterBgpPeerCustomLearnedIpRangesRange(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNestedComputeRouterBgpPeerCustomLearnedRoutePriority(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
func expandNestedComputeRouterBgpPeerBfd(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedSessionInitializationMode, err := expandNestedComputeRouterBgpPeerBfdSessionInitializationMode(original["session_initialization_mode"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedSessionInitializationMode); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["sessionInitializationMode"] = transformedSessionInitializationMode
	}

	transformedMinTransmitInterval, err := expandNestedComputeRouterBgpPeerBfdMinTransmitInterval(original["min_transmit_interval"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMinTransmitInterval); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["minTransmitInterval"] = transformedMinTransmitInterval
	}

	transformedMinReceiveInterval, err := expandNestedComputeRouterBgpPeerBfdMinReceiveInterval(original["min_receive_interval"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMinReceiveInterval); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["minReceiveInterval"] = transformedMinReceiveInterval
	}

	transformedMultiplier, err := expandNestedComputeRouterBgpPeerBfdMultiplier(original["multiplier"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMultiplier); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["multiplier"] = transformedMultiplier
	}

	return transformed, nil
}

func expandNestedComputeRouterBgpPeerMd5AuthenticationKey(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedMd5AuthenticationKeyName, err := expandNestedComputeRouterBgpPeerMd5AuthenticationKeyMd5AuthenticationKeyName(original["name"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMd5AuthenticationKeyName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["name"] = transformedMd5AuthenticationKeyName
	}

	transformedMd5AuthenticationKeyValue, err := expandNestedComputeRouterBgpPeerMd5AuthenticationKeyMd5AuthenticationKeyValue(original["key"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMd5AuthenticationKeyValue); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["key"] = transformedMd5AuthenticationKeyValue
	}

	return transformed, nil
}

func expandNestedComputeRouterBgpPeerMd5AuthenticationKeyMd5AuthenticationKeyName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNestedComputeRouterBgpPeerMd5AuthenticationKeyMd5AuthenticationKeyValue(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNestedComputeRouterBgpPeerBfdSessionInitializationMode(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNestedComputeRouterBgpPeerBfdMinTransmitInterval(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNestedComputeRouterBgpPeerBfdMinReceiveInterval(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNestedComputeRouterBgpPeerBfdMultiplier(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNestedComputeRouterBgpPeerEnable(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	if v == nil {
		return nil, nil
	}

	return strings.ToUpper(strconv.FormatBool(v.(bool))), nil
}

func expandNestedComputeRouterBgpPeerRouterApplianceInstance(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	f, err := tpgresource.ParseZonalFieldValue("instances", v.(string), "project", "zone", d, config, true)
	if err != nil {
		return nil, fmt.Errorf("Invalid value for router_appliance_instance: %s", err)
	}
	return f.RelativeLink(), nil
}

func expandNestedComputeRouterBgpPeerEnableIpv6(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNestedComputeRouterBgpPeerEnableIpv4(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNestedComputeRouterBgpPeerIpv4NexthopAddress(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNestedComputeRouterBgpPeerPeerIpv4NexthopAddress(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNestedComputeRouterBgpPeerIpv6NexthopAddress(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNestedComputeRouterBgpPeerPeerIpv6NexthopAddress(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func flattenNestedComputeRouterBgpPeer(d *schema.ResourceData, meta interface{}, res map[string]interface{}) (map[string]interface{}, error) {
	var v interface{}
	var ok bool

	v, ok = res["bgpPeers"]
	if !ok || v == nil {
		return nil, nil
	}

	switch v.(type) {
	case []interface{}:
		break
	case map[string]interface{}:
		// Construct list out of single nested resource
		v = []interface{}{v}
	default:
		return nil, fmt.Errorf("expected list or map for value bgpPeers. Actual value: %v", v)
	}

	_, item, err := resourceComputeRouterBgpPeerFindNestedObjectInList(d, meta, v.([]interface{}))
	if err != nil {
		return nil, err
	}
	return item, nil
}

func resourceComputeRouterBgpPeerFindNestedObjectInList(d *schema.ResourceData, meta interface{}, items []interface{}) (index int, item map[string]interface{}, err error) {
	expectedName, err := expandNestedComputeRouterBgpPeerName(d.Get("name"), d, meta.(*transport_tpg.Config))
	if err != nil {
		return -1, nil, err
	}
	expectedFlattenedName := flattenNestedComputeRouterBgpPeerName(expectedName, d, meta.(*transport_tpg.Config))

	// Search list for this resource.
	for idx, itemRaw := range items {
		if itemRaw == nil {
			continue
		}
		item := itemRaw.(map[string]interface{})

		itemName := flattenNestedComputeRouterBgpPeerName(item["name"], d, meta.(*transport_tpg.Config))
		// IsEmptyValue check so that if one is nil and the other is "", that's considered a match
		if !(tpgresource.IsEmptyValue(reflect.ValueOf(itemName)) && tpgresource.IsEmptyValue(reflect.ValueOf(expectedFlattenedName))) && !reflect.DeepEqual(itemName, expectedFlattenedName) {
			log.Printf("[DEBUG] Skipping item with name= %#v, looking for %#v)", itemName, expectedFlattenedName)
			continue
		}
		log.Printf("[DEBUG] Found item for resource %q: %#v)", d.Id(), item)
		return idx, item, nil
	}
	return -1, nil, nil
}

// PatchCreateEncoder handles creating request data to PATCH parent resource
// with list including new object.
func resourceComputeRouterBgpPeerPatchCreateEncoder(d *schema.ResourceData, meta interface{}, obj map[string]interface{}) (map[string]interface{}, error) {
	currBgpPeerItems, currMd5AuthenticationKeys, err := resourceComputeRouterBgpPeerListForPatch(d, meta)
	if err != nil {
		return nil, err
	}

	_, found, err := resourceComputeRouterBgpPeerFindNestedObjectInList(d, meta, currBgpPeerItems)
	if err != nil {
		return nil, err
	}

	// Return error if item already created.
	if found != nil {
		return nil, fmt.Errorf("Unable to create RouterBgpPeer, existing object already found: %+v", found)
	}

	var res map[string]interface{}

	// Return list with the resource to create appended
	val, ok := obj["md5AuthenticationKey"]

	if ok {
		kvp := val.(map[string]interface{})
		res = map[string]interface{}{
			"bgpPeers":              append(currBgpPeerItems, obj),
			"md5AuthenticationKeys": append(currMd5AuthenticationKeys, kvp),
		}

		//we need to remove this key from the object as it not a part of bgpRouterPeer
		delete(obj, "md5AuthenticationKey")
	} else {
		res = map[string]interface{}{
			"bgpPeers": append(currBgpPeerItems, obj),
		}
	}

	return res, nil
}

// PatchUpdateEncoder handles creating request data to PATCH parent resource
// with list including updated object.
func resourceComputeRouterBgpPeerPatchUpdateEncoder(d *schema.ResourceData, meta interface{}, obj map[string]interface{}) (map[string]interface{}, error) {
	bgpPeerItems, md5AuthenticationKeys, err := resourceComputeRouterBgpPeerListForPatch(d, meta)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] inside UpdateEncoder - bgpPeerItems:  %+v, md5AuthenticationKeys - %+v", bgpPeerItems, md5AuthenticationKeys)
	log.Printf("[DEBUG] inside UpdateEncoder - obj:  %+v", obj)

	idx, item, err := resourceComputeRouterBgpPeerFindNestedObjectInList(d, meta, bgpPeerItems)
	if err != nil {
		return nil, err
	}

	// Return error if item to update does not exist.
	if item == nil {
		return nil, fmt.Errorf("Unable to update RouterBgpPeer %q - not found in list", d.Id())
	}

	var md5AuthenticationKey map[string]interface{}
	var deletedKeyName interface{}
	var wasPresent bool
	val, ok := obj["md5AuthenticationKey"]
	if ok {
		md5AuthenticationKey = val.(map[string]interface{})
		//remove key from this map as it not needed here
		delete(obj, "md5AuthenticationKey")
	} else {
		//check if key used to be present
		deletedKeyName, wasPresent = item["md5AuthenticationKeyName"]
		if wasPresent {
			delete(item, "md5AuthenticationKeyName")
		}
	}

	//merging the bgpRouterPeer objects
	for k, v := range obj {
		item[k] = v
	}
	log.Printf("[DEBUG] UpdateEncoder - sending new object to be updated %#v", item)

	//merging the md5AuthenticationKeys objects
	isKeyNew := true
	log.Printf("[DEBUG] UpdateEncoder - currentMd5AuthenticationKeys %#v", md5AuthenticationKeys)
	for i, val := range md5AuthenticationKeys {
		key := val.(map[string]interface{})
		if key["name"] == md5AuthenticationKey["name"] {
			key = md5AuthenticationKey
			md5AuthenticationKeys[i] = key
			isKeyNew = false
		}

		if key["name"] == deletedKeyName {
			//if the key was deleted, then remove it from the parent router object as well
			md5AuthenticationKeys = append(md5AuthenticationKeys[:i], md5AuthenticationKeys[i+1:]...)
			log.Printf("[DEBUG] deleting unused key from parent object ,md5AuthenticationKeys - %+v", md5AuthenticationKeys)
		}

	}
	bgpPeerItems[idx] = item
	if isKeyNew {
		md5AuthenticationKeys = append(md5AuthenticationKeys, md5AuthenticationKey)
	}

	res := map[string]interface{}{
		"bgpPeers":              bgpPeerItems,
		"md5AuthenticationKeys": md5AuthenticationKeys,
	}

	return res, nil
}

// PatchDeleteEncoder handles creating request data to PATCH parent resource
// with list excluding object to delete.
func resourceComputeRouterBgpPeerPatchDeleteEncoder(d *schema.ResourceData, meta interface{}, obj map[string]interface{}) (map[string]interface{}, error) {
	currItems, md5AuthenticationKeys, err := resourceComputeRouterBgpPeerListForPatch(d, meta)
	if err != nil {
		return nil, err
	}

	idx, item, err := resourceComputeRouterBgpPeerFindNestedObjectInList(d, meta, currItems)
	if err != nil {
		return nil, err
	}
	if item == nil {
		// Spoof 404 error for proper handling by Delete (i.e. no-op)
		return nil, tpgresource.Fake404("nested", "ComputeRouterBgpPeer")
	}

	//if the removed bgp peer has some md5AuthKey associated with it, then remove the key from the router parent object as well
	keyName := item["md5AuthenticationKeyName"]
	for i, val := range md5AuthenticationKeys {
		key := val.(map[string]interface{})
		if key["name"] == keyName {
			md5AuthenticationKeys = append(md5AuthenticationKeys[:i], md5AuthenticationKeys[i+1:]...)
		}
	}

	updatedItems := append(currItems[:idx], currItems[idx+1:]...)
	res := map[string]interface{}{
		"bgpPeers":              updatedItems,
		"md5AuthenticationKeys": md5AuthenticationKeys,
	}

	return res, nil
}

// ListForPatch handles making API request to get parent resource and
// extracting list of objects.
func resourceComputeRouterBgpPeerListForPatch(d *schema.ResourceData, meta interface{}) ([]interface{}, []interface{}, error) {
	config := meta.(*transport_tpg.Config)
	url, err := tpgresource.ReplaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/regions/{{region}}/routers/{{router}}")
	if err != nil {
		return nil, nil, err
	}
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return nil, nil, err
	}

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return nil, nil, err
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return nil, nil, err
	}

	var v interface{}
	var ok bool
	var ls, keys []interface{}
	var lsOk, keysOk bool

	v, ok = res["bgpPeers"]
	if ok && v != nil {
		ls, lsOk = v.([]interface{})
		if !lsOk {
			return nil, nil, fmt.Errorf(`expected list for nested field "bgpPeers"`)
		}
	}
	v, ok = res["md5AuthenticationKeys"]
	if ok && v != nil {
		keys, keysOk = v.([]interface{})
		if !keysOk {
			return nil, nil, fmt.Errorf(`expected list for nested field "md5AuthenticationKeys"`)
		}
	}

	if lsOk && keysOk {
		return ls, keys, nil
	} else if !lsOk && keysOk {
		return nil, keys, nil
	} else if lsOk && !keysOk {
		return ls, nil, nil
	} else {
		return nil, nil, nil
	}
}
