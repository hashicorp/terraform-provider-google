package google

import (
	"bytes"
	"fmt"
	"net"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"

	"google.golang.org/api/compute/v1"
)

func resourceComputeVpnTunnel() *schema.Resource {
	return &schema.Resource{
		// Unfortunately, the VPNTunnelService does not support update
		// operations. This is why everything is marked forcenew
		Create: resourceComputeVpnTunnelCreate,
		Read:   resourceComputeVpnTunnelRead,
		Delete: resourceComputeVpnTunnelDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"peer_ip": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validatePeerAddr,
			},

			"shared_secret": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"target_vpn_gateway": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"detailed_status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"ike_version": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2,
				ForceNew: true,
			},

			"local_traffic_selector": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"remote_traffic_selector": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"router": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"self_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceComputeVpnTunnelCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	peerIp := d.Get("peer_ip").(string)
	sharedSecret := d.Get("shared_secret").(string)
	targetVpnGateway := d.Get("target_vpn_gateway").(string)
	ikeVersion := d.Get("ike_version").(int)

	if ikeVersion < 1 || ikeVersion > 2 {
		return fmt.Errorf("Only IKE version 1 or 2 supported, not %d", ikeVersion)
	}

	// Build up the list of sources
	var localTrafficSelectors []string
	if v := d.Get("local_traffic_selector").(*schema.Set); v.Len() > 0 {
		localTrafficSelectors = make([]string, v.Len())
		for i, v := range v.List() {
			localTrafficSelectors[i] = v.(string)
		}
	}

	var remoteTrafficSelectors []string
	if v := d.Get("remote_traffic_selector").(*schema.Set); v.Len() > 0 {
		remoteTrafficSelectors = make([]string, v.Len())
		for i, v := range v.List() {
			remoteTrafficSelectors[i] = v.(string)
		}
	}

	vpnTunnelsService := compute.NewVpnTunnelsService(config.clientCompute)

	vpnTunnel := &compute.VpnTunnel{
		Name:                  name,
		PeerIp:                peerIp,
		SharedSecret:          sharedSecret,
		TargetVpnGateway:      targetVpnGateway,
		IkeVersion:            int64(ikeVersion),
		LocalTrafficSelector:  localTrafficSelectors,
		RemoteTrafficSelector: remoteTrafficSelectors,
	}

	if v, ok := d.GetOk("description"); ok {
		vpnTunnel.Description = v.(string)
	}

	if v, ok := d.GetOk("router"); ok {
		routerLink, err := getRouterLink(config, project, region, v.(string))
		if err != nil {
			return err
		}
		vpnTunnel.Router = routerLink
	}

	op, err := vpnTunnelsService.Insert(project, region, vpnTunnel).Do()
	if err != nil {
		return fmt.Errorf("Error Inserting VPN Tunnel %s : %s", name, err)
	}

	err = computeOperationWait(config, op, project, "Inserting VPN Tunnel")
	if err != nil {
		return fmt.Errorf("Error Waiting to Insert VPN Tunnel %s: %s", name, err)
	}

	return resourceComputeVpnTunnelRead(d, meta)
}

func resourceComputeVpnTunnelRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)

	vpnTunnelsService := compute.NewVpnTunnelsService(config.clientCompute)

	vpnTunnel, err := vpnTunnelsService.Get(project, region, name).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("VPN Tunnel %q", d.Get("name").(string)))
	}

	localTrafficSelectors := []string{}
	for _, lts := range vpnTunnel.LocalTrafficSelector {
		localTrafficSelectors = append(localTrafficSelectors, lts)
	}
	d.Set("local_traffic_selector", localTrafficSelectors)

	remoteTrafficSelectors := []string{}
	for _, rts := range vpnTunnel.RemoteTrafficSelector {
		remoteTrafficSelectors = append(remoteTrafficSelectors, rts)
	}
	d.Set("remote_traffic_selector", remoteTrafficSelectors)

	d.Set("detailed_status", vpnTunnel.DetailedStatus)
	d.Set("self_link", vpnTunnel.SelfLink)

	d.SetId(name)

	return nil
}

func resourceComputeVpnTunnelDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)

	vpnTunnelsService := compute.NewVpnTunnelsService(config.clientCompute)

	op, err := vpnTunnelsService.Delete(project, region, name).Do()
	if err != nil {
		return fmt.Errorf("Error Reading VPN Tunnel %s: %s", name, err)
	}

	err = computeOperationWait(config, op, project, "Deleting VPN Tunnel")
	if err != nil {
		return fmt.Errorf("Error Waiting to Delete VPN Tunnel %s: %s", name, err)
	}

	return nil
}

// validatePeerAddr returns false if a tunnel's peer_ip property
// is invalid. Currently, only addresses that collide with RFC
// 5735 (https://tools.ietf.org/html/rfc5735) fail validation.
func validatePeerAddr(i interface{}, val string) ([]string, []error) {
	ip := net.ParseIP(i.(string))
	if ip == nil {
		return nil, []error{fmt.Errorf("could not parse %q to IP address", val)}
	}
	for _, test := range invalidPeerAddrs {
		if bytes.Compare(ip, test.from) >= 0 && bytes.Compare(ip, test.to) <= 0 {
			return nil, []error{fmt.Errorf("address is invalid (is between %q and %q, conflicting with RFC5735)", test.from, test.to)}
		}
	}
	return nil, nil
}

// invalidPeerAddrs is a collection of IP addres ranges that represent
// a conflict with RFC 5735 (https://tools.ietf.org/html/rfc5735#page-3).
// CIDR range notations in the RFC were converted to a (from, to) pair
// for easy checking with bytes.Compare.
var invalidPeerAddrs = []struct {
	from net.IP
	to   net.IP
}{
	{
		from: net.ParseIP("0.0.0.0"),
		to:   net.ParseIP("0.255.255.255"),
	},
	{
		from: net.ParseIP("10.0.0.0"),
		to:   net.ParseIP("10.255.255.255"),
	},
	{
		from: net.ParseIP("127.0.0.0"),
		to:   net.ParseIP("127.255.255.255"),
	},
	{
		from: net.ParseIP("169.254.0.0"),
		to:   net.ParseIP("169.254.255.255"),
	},
	{
		from: net.ParseIP("172.16.0.0"),
		to:   net.ParseIP("172.31.255.255"),
	},
	{
		from: net.ParseIP("192.0.0.0"),
		to:   net.ParseIP("192.0.0.255"),
	},
	{
		from: net.ParseIP("192.0.2.0"),
		to:   net.ParseIP("192.0.2.255"),
	},
	{
		from: net.ParseIP("192.88.99.0"),
		to:   net.ParseIP("192.88.99.255"),
	},
	{
		from: net.ParseIP("192.168.0.0"),
		to:   net.ParseIP("192.168.255.255"),
	},
	{
		from: net.ParseIP("198.18.0.0"),
		to:   net.ParseIP("198.19.255.255"),
	},
	{
		from: net.ParseIP("198.51.100.0"),
		to:   net.ParseIP("198.51.100.255"),
	},
	{
		from: net.ParseIP("203.0.113.0"),
		to:   net.ParseIP("203.0.113.255"),
	},
	{
		from: net.ParseIP("224.0.0.0"),
		to:   net.ParseIP("239.255.255.255"),
	},
	{
		from: net.ParseIP("240.0.0.0"),
		to:   net.ParseIP("255.255.255.255"),
	},
	{
		from: net.ParseIP("255.255.255.255"),
		to:   net.ParseIP("255.255.255.255"),
	},
}

func getVpnTunnelLink(config *Config, project string, region string, tunnel string) (string, error) {

	if !strings.HasPrefix(tunnel, "https://www.googleapis.com/compute/") {
		// Tunnel value provided is just the name, lookup the tunnel SelfLink
		tunnelData, err := config.clientCompute.VpnTunnels.Get(
			project, region, tunnel).Do()
		if err != nil {
			return "", fmt.Errorf("Error reading tunnel: %s", err)
		}
		tunnel = tunnelData.SelfLink
	}

	return tunnel, nil

}

func getVpnTunnelName(vpntunnel string) (string, error) {

	if strings.HasPrefix(vpntunnel, "https://www.googleapis.com/compute/") {
		// extract the VPN tunnel name from SelfLink URL
		vpntunnelName := vpntunnel[strings.LastIndex(vpntunnel, "/")+1:]
		if vpntunnelName == "" {
			return "", fmt.Errorf("VPN tunnel url not valid")
		}
		return vpntunnelName, nil
	}

	return vpntunnel, nil
}
