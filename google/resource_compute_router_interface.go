package google

import (
	"fmt"
	"log"
	"time"

	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/googleapi"

	"google.golang.org/api/compute/v1"
)

func resourceComputeRouterInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeRouterInterfaceCreate,
		Read:   resourceComputeRouterInterfaceRead,
		Delete: resourceComputeRouterInterfaceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceComputeRouterInterfaceImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `A unique name for the interface, required by GCE. Changing this forces a new interface to be created.`,
			},
			"router": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the router this interface will be attached to. Changing this forces a new interface to be created.`,
			},
			"vpn_tunnel": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				AtLeastOneOf:     []string{"ip_range", "interconnect_attachment", "subnetwork", "vpn_tunnel"},
				ConflictsWith:    []string{"interconnect_attachment", "subnetwork"},
				Description:      `The name or resource link to the VPN tunnel this interface will be linked to. Changing this forces a new interface to be created. Only one of vpn_tunnel, interconnect_attachment or subnetwork can be specified.`,
			},
			"interconnect_attachment": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				AtLeastOneOf:     []string{"ip_range", "interconnect_attachment", "subnetwork", "vpn_tunnel"},
				ConflictsWith:    []string{"subnetwork", "vpn_tunnel"},
				Description:      `The name or resource link to the VLAN interconnect for this interface. Changing this forces a new interface to be created. Only one of interconnect_attachment, subnetwork or vpn_tunnel can be specified.`,
			},
			"ip_range": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				AtLeastOneOf: []string{"ip_range", "interconnect_attachment", "subnetwork", "vpn_tunnel"},
				Description:  `The IP address and range of the interface. The IP range must be in the RFC3927 link-local IP space. Changing this forces a new interface to be created.`,
			},
			"private_ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The regional private internal IP address that is used to establish BGP sessions to a VM instance acting as a third-party Router Appliance. Changing this forces a new interface to be created.`,
			},
			"subnetwork": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				AtLeastOneOf:     []string{"ip_range", "interconnect_attachment", "subnetwork", "vpn_tunnel"},
				ConflictsWith:    []string{"interconnect_attachment", "vpn_tunnel"},
				Description:      `The URI of the subnetwork resource that this interface belongs to, which must be in the same region as the Cloud Router. Changing this forces a new interface to be created. Only one of subnetwork, interconnect_attachment or vpn_tunnel can be specified.`,
			},
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The ID of the project in which this interface's router belongs. If it is not provided, the provider project is used. Changing this forces a new interface to be created.`,
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region this interface's router sits in. If not specified, the project region will be used. Changing this forces a new interface to be created.`,
			},

			"redundant_interface": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The name of the interface that is redundant to this interface. Changing this forces a new interface to be created.`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceComputeRouterInterfaceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	routerName := d.Get("router").(string)
	ifaceName := d.Get("name").(string)

	routerLock := getRouterLockName(region, routerName)
	mutexKV.Lock(routerLock)
	defer mutexKV.Unlock(routerLock)

	routersService := config.NewComputeClient(userAgent).Routers
	router, err := routersService.Get(project, region, routerName).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			log.Printf("[WARN] Removing router interface %s because its router %s/%s is gone", ifaceName, region, routerName)
			d.SetId("")

			return nil
		}

		return fmt.Errorf("Error Reading router %s/%s: %s", region, routerName, err)
	}

	ifaces := router.Interfaces

	for _, iface := range ifaces {
		if iface.Name == ifaceName {
			d.SetId("")
			return fmt.Errorf("Router %s has interface %s already", routerName, ifaceName)
		}
	}

	iface := &compute.RouterInterface{Name: ifaceName}

	if riVal, ok := d.GetOk("redundant_interface"); ok {
		iface.RedundantInterface = riVal.(string)
	}

	if ipRangeVal, ok := d.GetOk("ip_range"); ok {
		iface.IpRange = ipRangeVal.(string)
	}

	if privateIpVal, ok := d.GetOk("private_ip_address"); ok {
		iface.PrivateIpAddress = privateIpVal.(string)
	}

	if vpnVal, ok := d.GetOk("vpn_tunnel"); ok {
		vpnTunnel, err := getVpnTunnelLink(config, project, region, vpnVal.(string), userAgent)
		if err != nil {
			return err
		}
		iface.LinkedVpnTunnel = vpnTunnel
	}

	if icVal, ok := d.GetOk("interconnect_attachment"); ok {
		interconnectAttachment, err := getInterconnectAttachmentLink(config, project, region, icVal.(string), userAgent)
		if err != nil {
			return err
		}
		iface.LinkedInterconnectAttachment = interconnectAttachment
	}

	if subVal, ok := d.GetOk("subnetwork"); ok {
		iface.Subnetwork = subVal.(string)
	}

	log.Printf("[INFO] Adding interface %s", ifaceName)
	ifaces = append(ifaces, iface)
	patchRouter := &compute.Router{
		Interfaces: ifaces,
	}

	log.Printf("[DEBUG] Updating router %s/%s with interfaces: %+v", region, routerName, ifaces)
	op, err := routersService.Patch(project, region, router.Name, patchRouter).Do()
	if err != nil {
		return fmt.Errorf("Error patching router %s/%s: %s", region, routerName, err)
	}
	d.SetId(fmt.Sprintf("%s/%s/%s", region, routerName, ifaceName))
	err = computeOperationWaitTime(config, op, project, "Patching router", userAgent, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error waiting to patch router %s/%s: %s", region, routerName, err)
	}

	return resourceComputeRouterInterfaceRead(d, meta)
}

func resourceComputeRouterInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	routerName := d.Get("router").(string)
	ifaceName := d.Get("name").(string)

	routersService := config.NewComputeClient(userAgent).Routers
	router, err := routersService.Get(project, region, routerName).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			log.Printf("[WARN] Removing router interface %s because its router %s/%s is gone", ifaceName, region, routerName)
			d.SetId("")

			return nil
		}

		return fmt.Errorf("Error Reading router %s/%s: %s", region, routerName, err)
	}

	for _, iface := range router.Interfaces {

		if iface.Name == ifaceName {
			d.SetId(fmt.Sprintf("%s/%s/%s", region, routerName, ifaceName))
			if err := d.Set("vpn_tunnel", iface.LinkedVpnTunnel); err != nil {
				return fmt.Errorf("Error setting vpn_tunnel: %s", err)
			}
			if err := d.Set("interconnect_attachment", iface.LinkedInterconnectAttachment); err != nil {
				return fmt.Errorf("Error setting interconnect_attachment: %s", err)
			}
			if err := d.Set("ip_range", iface.IpRange); err != nil {
				return fmt.Errorf("Error setting ip_range: %s", err)
			}
			if err := d.Set("private_ip_address", iface.PrivateIpAddress); err != nil {
				return fmt.Errorf("Error setting private_ip_address: %s", err)
			}
			if err := d.Set("subnetwork", iface.Subnetwork); err != nil {
				return fmt.Errorf("Error setting subnetwork: %s", err)
			}
			if err := d.Set("region", region); err != nil {
				return fmt.Errorf("Error setting region: %s", err)
			}
			if err := d.Set("project", project); err != nil {
				return fmt.Errorf("Error setting project: %s", err)
			}
			if err := d.Set("redundant_interface", iface.RedundantInterface); err != nil {
				return fmt.Errorf("Error setting redundant interface: %s", err)
			}
			return nil
		}
	}

	log.Printf("[WARN] Removing router interface %s/%s/%s because it is gone", region, routerName, ifaceName)
	d.SetId("")
	return nil
}

func resourceComputeRouterInterfaceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	routerName := d.Get("router").(string)
	ifaceName := d.Get("name").(string)

	routerLock := getRouterLockName(region, routerName)
	mutexKV.Lock(routerLock)
	defer mutexKV.Unlock(routerLock)

	routersService := config.NewComputeClient(userAgent).Routers
	router, err := routersService.Get(project, region, routerName).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			log.Printf("[WARN] Removing router interface %s because its router %s/%s is gone", ifaceName, region, routerName)

			return nil
		}

		return fmt.Errorf("Error Reading Router %s: %s", routerName, err)
	}

	var ifaceFound bool

	newIfaces := make([]*compute.RouterInterface, 0, len(router.Interfaces))
	for _, iface := range router.Interfaces {

		if iface.Name == ifaceName {
			ifaceFound = true
			continue
		} else {
			newIfaces = append(newIfaces, iface)
		}
	}

	if !ifaceFound {
		log.Printf("[DEBUG] Router %s/%s had no interface %s already", region, routerName, ifaceName)
		d.SetId("")
		return nil
	}

	log.Printf(
		"[INFO] Removing interface %s from router %s/%s", ifaceName, region, routerName)
	patchRouter := &compute.Router{
		Interfaces: newIfaces,
	}

	if len(newIfaces) == 0 {
		patchRouter.ForceSendFields = append(patchRouter.ForceSendFields, "Interfaces")
	}

	log.Printf("[DEBUG] Updating router %s/%s with interfaces: %+v", region, routerName, newIfaces)
	op, err := routersService.Patch(project, region, router.Name, patchRouter).Do()
	if err != nil {
		return fmt.Errorf("Error patching router %s/%s: %s", region, routerName, err)
	}

	err = computeOperationWaitTime(config, op, project, "Patching router", userAgent, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf("Error waiting to patch router %s/%s: %s", region, routerName, err)
	}

	d.SetId("")
	return nil
}

func resourceComputeRouterInterfaceImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid router interface specifier. Expecting {region}/{router}/{interface}")
	}

	if err := d.Set("region", parts[0]); err != nil {
		return nil, fmt.Errorf("Error setting region: %s", err)
	}
	if err := d.Set("router", parts[1]); err != nil {
		return nil, fmt.Errorf("Error setting router: %s", err)
	}
	if err := d.Set("name", parts[2]); err != nil {
		return nil, fmt.Errorf("Error setting name: %s", err)
	}

	return []*schema.ResourceData{d}, nil
}
