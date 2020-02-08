package google

import (
	"fmt"
	"log"

	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

func resourceComputeRouterInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeRouterInterfaceCreate,
		Read:   resourceComputeRouterInterfaceRead,
		Delete: resourceComputeRouterInterfaceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceComputeRouterInterfaceImportState,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"router": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpn_tunnel": {
				Type:             schema.TypeString,
				ConflictsWith:    []string{"interconnect_attachment"},
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				AtLeastOneOf:     []string{"vpn_tunnel", "interconnect_attachment", "ip_range"},
			},
			"interconnect_attachment": {
				Type:             schema.TypeString,
				ConflictsWith:    []string{"vpn_tunnel"},
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				AtLeastOneOf:     []string{"vpn_tunnel", "interconnect_attachment", "ip_range"},
			},
			"ip_range": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"vpn_tunnel", "interconnect_attachment", "ip_range"},
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceComputeRouterInterfaceCreate(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)

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

	routersService := config.clientCompute.Routers
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

	if ipVal, ok := d.GetOk("ip_range"); ok {
		iface.IpRange = ipVal.(string)
	}

	if vpnVal, ok := d.GetOk("vpn_tunnel"); ok {
		vpnTunnel, err := getVpnTunnelLink(config, project, region, vpnVal.(string))
		if err != nil {
			return err
		}
		iface.LinkedVpnTunnel = vpnTunnel
	}

	if icVal, ok := d.GetOk("interconnect_attachment"); ok {
		interconnectAttachment, err := getInterconnectAttachmentLink(config, project, region, icVal.(string))
		if err != nil {
			return err
		}
		iface.LinkedInterconnectAttachment = interconnectAttachment
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
	err = computeOperationWait(config, op, project, "Patching router")
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error waiting to patch router %s/%s: %s", region, routerName, err)
	}

	return resourceComputeRouterInterfaceRead(d, meta)
}

func resourceComputeRouterInterfaceRead(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)

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

	routersService := config.clientCompute.Routers
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
			d.Set("vpn_tunnel", iface.LinkedVpnTunnel)
			d.Set("interconnect_attachment", iface.LinkedInterconnectAttachment)
			d.Set("ip_range", iface.IpRange)
			d.Set("region", region)
			d.Set("project", project)
			return nil
		}
	}

	log.Printf("[WARN] Removing router interface %s/%s/%s because it is gone", region, routerName, ifaceName)
	d.SetId("")
	return nil
}

func resourceComputeRouterInterfaceDelete(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)

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

	routersService := config.clientCompute.Routers
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

	err = computeOperationWait(config, op, project, "Patching router")
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

	d.Set("region", parts[0])
	d.Set("router", parts[1])
	d.Set("name", parts[2])

	return []*schema.ResourceData{d}, nil
}
