package google

import (
	"fmt"
	"log"
	"time"

	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	computeBeta "google.golang.org/api/compute/v0.beta"

	"google.golang.org/api/googleapi"
)

var (
	routerNatSubnetworkConfig = &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			// this field is optional with a default in the API, but we
			// don't have the ability to support complex defaults inside
			// nested fields
			"source_ip_ranges_to_nat": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"secondary_ip_range_names": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
)

func resourceComputeRouterNat() *schema.Resource {
	return &schema.Resource{
		// TODO(https://github.com/GoogleCloudPlatform/magic-modules/issues/963): Implement Update
		Create: resourceComputeRouterNatCreate,
		Read:   resourceComputeRouterNatRead,
		Delete: resourceComputeRouterNatDelete,
		Importer: &schema.ResourceImporter{
			State: resourceComputeRouterNatImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateRFC1035Name(2, 63),
			},
			"router": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nat_ip_allocate_option": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"MANUAL_ONLY", "AUTO_ONLY"}, false),
			},
			"nat_ips": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"source_subnetwork_ip_ranges_to_nat": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALL_SUBNETWORKS_ALL_IP_RANGES", "ALL_SUBNETWORKS_ALL_PRIMARY_IP_RANGES", "LIST_OF_SUBNETWORKS"}, false),
			},
			"subnetwork": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     routerNatSubnetworkConfig,
			},
			"min_ports_per_vm": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"udp_idle_timeout_sec": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"icmp_idle_timeout_sec": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"tcp_established_idle_timeout_sec": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"tcp_transitory_idle_timeout_sec": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
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

func resourceComputeRouterNatCreate(d *schema.ResourceData, meta interface{}) error {

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
	natName := d.Get("name").(string)

	routerLock := getRouterLockName(region, routerName)
	mutexKV.Lock(routerLock)
	defer mutexKV.Unlock(routerLock)

	routersService := config.clientComputeBeta.Routers
	router, err := routersService.Get(project, region, routerName).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return fmt.Errorf("Router %s/%s not found", region, routerName)
		}

		return fmt.Errorf("Error Reading router %s/%s: %s", region, routerName, err)
	}

	nats := router.Nats
	for _, nat := range nats {
		if nat.Name == natName {
			return fmt.Errorf("Router %s has nat %s already", routerName, natName)
		}
	}

	nat := &computeBeta.RouterNat{
		Name:                          natName,
		NatIpAllocateOption:           d.Get("nat_ip_allocate_option").(string),
		NatIps:                        convertStringArr(d.Get("nat_ips").(*schema.Set).List()),
		SourceSubnetworkIpRangesToNat: d.Get("source_subnetwork_ip_ranges_to_nat").(string),
		MinPortsPerVm:                 int64(d.Get("min_ports_per_vm").(int)),
		UdpIdleTimeoutSec:             int64(d.Get("udp_idle_timeout_sec").(int)),
		IcmpIdleTimeoutSec:            int64(d.Get("icmp_idle_timeout_sec").(int)),
		TcpEstablishedIdleTimeoutSec:  int64(d.Get("tcp_established_idle_timeout_sec").(int)),
		TcpTransitoryIdleTimeoutSec:   int64(d.Get("tcp_transitory_idle_timeout_sec").(int)),
	}

	if v, ok := d.GetOk("subnetwork"); ok {
		nat.Subnetworks = expandSubnetworks(v.(*schema.Set).List())
	}

	log.Printf("[INFO] Adding nat %s", natName)
	nats = append(nats, nat)
	patchRouter := &computeBeta.Router{
		Nats: nats,
	}

	log.Printf("[DEBUG] Updating router %s/%s with nats: %+v", region, routerName, nats)
	op, err := routersService.Patch(project, region, router.Name, patchRouter).Do()
	if err != nil {
		return fmt.Errorf("Error patching router %s/%s: %s", region, routerName, err)
	}
	d.SetId(fmt.Sprintf("%s/%s/%s", region, routerName, natName))
	err = computeBetaOperationWaitTime(config.clientCompute, op, project, "Patching router", int(d.Timeout(schema.TimeoutCreate).Minutes()))
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error waiting to patch router %s/%s: %s", region, routerName, err)
	}

	return resourceComputeRouterNatRead(d, meta)
}

func resourceComputeRouterNatRead(d *schema.ResourceData, meta interface{}) error {

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
	natName := d.Get("name").(string)

	routersService := config.clientComputeBeta.Routers
	router, err := routersService.Get(project, region, routerName).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			log.Printf("[WARN] Removing router nat %s because its router %s/%s is gone", natName, region, routerName)
			d.SetId("")

			return nil
		}

		return fmt.Errorf("Error Reading router %s/%s: %s", region, routerName, err)
	}

	for _, nat := range router.Nats {

		if nat.Name == natName {
			d.SetId(fmt.Sprintf("%s/%s/%s", region, routerName, natName))
			d.Set("nat_ip_allocate_option", nat.NatIpAllocateOption)
			d.Set("nat_ips", schema.NewSet(schema.HashString, convertStringArrToInterface(convertSelfLinksToV1(nat.NatIps))))
			d.Set("source_subnetwork_ip_ranges_to_nat", nat.SourceSubnetworkIpRangesToNat)
			d.Set("min_ports_per_vm", nat.MinPortsPerVm)
			d.Set("udp_idle_timeout_sec", nat.UdpIdleTimeoutSec)
			d.Set("icmp_idle_timeout_sec", nat.IcmpIdleTimeoutSec)
			d.Set("tcp_established_idle_timeout_sec", nat.TcpEstablishedIdleTimeoutSec)
			d.Set("tcp_transitory_idle_timeout_sec", nat.TcpTransitoryIdleTimeoutSec)
			d.Set("region", region)
			d.Set("project", project)

			if err := d.Set("subnetwork", flattenRouterNatSubnetworkToNatBeta(nat.Subnetworks)); err != nil {
				return fmt.Errorf("Error reading router nat: %s", err)
			}

			return nil
		}
	}

	log.Printf("[WARN] Removing router nat %s/%s/%s because it is gone", region, routerName, natName)
	d.SetId("")
	return nil
}

func resourceComputeRouterNatDelete(d *schema.ResourceData, meta interface{}) error {

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
	natName := d.Get("name").(string)

	routerLock := getRouterLockName(region, routerName)
	mutexKV.Lock(routerLock)
	defer mutexKV.Unlock(routerLock)

	routersService := config.clientComputeBeta.Routers
	router, err := routersService.Get(project, region, routerName).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			log.Printf("[WARN] Removing router nat %s because its router %s/%s is gone", natName, region, routerName)

			return nil
		}

		return fmt.Errorf("Error Reading Router %s: %s", routerName, err)
	}

	var newNats []*computeBeta.RouterNat = make([]*computeBeta.RouterNat, 0, len(router.Nats))
	for _, nat := range router.Nats {
		if nat.Name == natName {
			continue
		} else {
			newNats = append(newNats, nat)
		}
	}

	if len(newNats) == len(router.Nats) {
		log.Printf("[DEBUG] Router %s/%s had no nat %s already", region, routerName, natName)
		d.SetId("")
		return nil
	}

	log.Printf("[INFO] Removing nat %s from router %s/%s", natName, region, routerName)
	patchRouter := &computeBeta.Router{
		Nats: newNats,
	}

	if len(newNats) == 0 {
		patchRouter.ForceSendFields = append(patchRouter.ForceSendFields, "Nats")
	}

	log.Printf("[DEBUG] Updating router %s/%s with nats: %+v", region, routerName, newNats)
	op, err := routersService.Patch(project, region, router.Name, patchRouter).Do()
	if err != nil {
		return fmt.Errorf("Error patching router %s/%s: %s", region, routerName, err)
	}

	err = computeBetaOperationWaitTime(config.clientCompute, op, project, "Patching router", int(d.Timeout(schema.TimeoutDelete).Minutes()))
	if err != nil {
		return fmt.Errorf("Error waiting to patch router %s/%s: %s", region, routerName, err)
	}

	d.SetId("")
	return nil
}

func resourceComputeRouterNatImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid router nat specifier. Expecting {region}/{router}/{nat}")
	}

	d.Set("region", parts[0])
	d.Set("router", parts[1])
	d.Set("name", parts[2])

	return []*schema.ResourceData{d}, nil
}

func expandSubnetworks(subnetworks []interface{}) []*computeBeta.RouterNatSubnetworkToNat {
	result := make([]*computeBeta.RouterNatSubnetworkToNat, 0, len(subnetworks))

	for _, subnetwork := range subnetworks {
		snm := subnetwork.(map[string]interface{})
		subnetworkToNat := computeBeta.RouterNatSubnetworkToNat{
			Name:                snm["name"].(string),
			SourceIpRangesToNat: convertStringSet(snm["source_ip_ranges_to_nat"].(*schema.Set)),
		}
		if v, ok := snm["secondary_ip_range_names"]; ok {
			subnetworkToNat.SecondaryIpRangeNames = convertStringSet(v.(*schema.Set))
		}
		result = append(result, &subnetworkToNat)
	}

	return result
}

func flattenRouterNatSubnetworkToNatBeta(subnetworksToNat []*computeBeta.RouterNatSubnetworkToNat) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(subnetworksToNat))
	for _, subnetworkToNat := range subnetworksToNat {
		stnMap := make(map[string]interface{})
		stnMap["name"] = ConvertSelfLinkToV1(subnetworkToNat.Name)
		stnMap["source_ip_ranges_to_nat"] = schema.NewSet(schema.HashString, convertStringArrToInterface(subnetworkToNat.SourceIpRangesToNat))
		stnMap["secondary_ip_range_names"] = schema.NewSet(schema.HashString, convertStringArrToInterface(subnetworkToNat.SecondaryIpRangeNames))
		result = append(result, stnMap)
	}
	return result
}

func convertSelfLinksToV1(selfLinks []string) []string {
	result := make([]string, 0, len(selfLinks))
	for _, selfLink := range selfLinks {
		result = append(result, ConvertSelfLinkToV1(selfLink))
	}
	return result
}
