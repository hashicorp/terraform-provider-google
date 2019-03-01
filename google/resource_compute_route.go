// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package google

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/compute/v1"
)

func resourceComputeRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeRouteCreate,
		Read:   resourceComputeRouteRead,
		Delete: resourceComputeRouteDelete,

		Importer: &schema.ResourceImporter{
			State: resourceComputeRouteImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(240 * time.Second),
			Delete: schema.DefaultTimeout(240 * time.Second),
		},

		Schema: map[string]*schema.Schema{
			"dest_range": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateRegexp(`^[a-z]([-a-z0-9]*[a-z0-9])?$`),
			},
			"network": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"next_hop_gateway": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
			},
			"next_hop_instance": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
			},
			"next_hop_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"next_hop_vpn_tunnel": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  1000,
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"next_hop_network": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"next_hop_instance_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceComputeRouteCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	obj := make(map[string]interface{})
	destRangeProp, err := expandComputeRouteDestRange(d.Get("dest_range"), d, config)
	if err != nil {
		return err
	} else if !isEmptyValue(reflect.ValueOf(destRangeProp)) {
		obj["destRange"] = destRangeProp
	}
	descriptionProp, err := expandComputeRouteDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if !isEmptyValue(reflect.ValueOf(descriptionProp)) {
		obj["description"] = descriptionProp
	}
	nameProp, err := expandComputeRouteName(d.Get("name"), d, config)
	if err != nil {
		return err
	} else if !isEmptyValue(reflect.ValueOf(nameProp)) {
		obj["name"] = nameProp
	}
	networkProp, err := expandComputeRouteNetwork(d.Get("network"), d, config)
	if err != nil {
		return err
	} else if !isEmptyValue(reflect.ValueOf(networkProp)) {
		obj["network"] = networkProp
	}
	priorityProp, err := expandComputeRoutePriority(d.Get("priority"), d, config)
	if err != nil {
		return err
	} else if !isEmptyValue(reflect.ValueOf(priorityProp)) {
		obj["priority"] = priorityProp
	}
	tagsProp, err := expandComputeRouteTags(d.Get("tags"), d, config)
	if err != nil {
		return err
	} else if !isEmptyValue(reflect.ValueOf(tagsProp)) {
		obj["tags"] = tagsProp
	}
	nextHopGatewayProp, err := expandComputeRouteNextHopGateway(d.Get("next_hop_gateway"), d, config)
	if err != nil {
		return err
	} else if !isEmptyValue(reflect.ValueOf(nextHopGatewayProp)) {
		obj["nextHopGateway"] = nextHopGatewayProp
	}
	nextHopInstanceProp, err := expandComputeRouteNextHopInstance(d.Get("next_hop_instance"), d, config)
	if err != nil {
		return err
	} else if !isEmptyValue(reflect.ValueOf(nextHopInstanceProp)) {
		obj["nextHopInstance"] = nextHopInstanceProp
	}
	nextHopIpProp, err := expandComputeRouteNextHopIp(d.Get("next_hop_ip"), d, config)
	if err != nil {
		return err
	} else if !isEmptyValue(reflect.ValueOf(nextHopIpProp)) {
		obj["nextHopIp"] = nextHopIpProp
	}
	nextHopVpnTunnelProp, err := expandComputeRouteNextHopVpnTunnel(d.Get("next_hop_vpn_tunnel"), d, config)
	if err != nil {
		return err
	} else if !isEmptyValue(reflect.ValueOf(nextHopVpnTunnelProp)) {
		obj["nextHopVpnTunnel"] = nextHopVpnTunnelProp
	}

	url, err := replaceVars(d, config, "https://www.googleapis.com/compute/v1/projects/{{project}}/global/routes")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new Route: %#v", obj)
	res, err := sendRequestWithTimeout(config, "POST", url, obj, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating Route: %s", err)
	}

	// Store the ID now
	id, err := replaceVars(d, config, "{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	op := &compute.Operation{}
	err = Convert(res, op)
	if err != nil {
		return err
	}

	waitErr := computeOperationWaitTime(
		config.clientCompute, op, project, "Creating Route",
		int(d.Timeout(schema.TimeoutCreate).Minutes()))

	if waitErr != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error waiting to create Route: %s", waitErr)
	}

	log.Printf("[DEBUG] Finished creating Route %q: %#v", d.Id(), res)

	return resourceComputeRouteRead(d, meta)
}

func resourceComputeRouteRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	url, err := replaceVars(d, config, "https://www.googleapis.com/compute/v1/projects/{{project}}/global/routes/{{name}}")
	if err != nil {
		return err
	}

	res, err := sendRequest(config, "GET", url, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("ComputeRoute %q", d.Id()))
	}

	res, err = resourceComputeRouteDecoder(d, meta, res)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading Route: %s", err)
	}

	if err := d.Set("dest_range", flattenComputeRouteDestRange(res["destRange"], d)); err != nil {
		return fmt.Errorf("Error reading Route: %s", err)
	}
	if err := d.Set("description", flattenComputeRouteDescription(res["description"], d)); err != nil {
		return fmt.Errorf("Error reading Route: %s", err)
	}
	if err := d.Set("name", flattenComputeRouteName(res["name"], d)); err != nil {
		return fmt.Errorf("Error reading Route: %s", err)
	}
	if err := d.Set("network", flattenComputeRouteNetwork(res["network"], d)); err != nil {
		return fmt.Errorf("Error reading Route: %s", err)
	}
	if err := d.Set("priority", flattenComputeRoutePriority(res["priority"], d)); err != nil {
		return fmt.Errorf("Error reading Route: %s", err)
	}
	if err := d.Set("tags", flattenComputeRouteTags(res["tags"], d)); err != nil {
		return fmt.Errorf("Error reading Route: %s", err)
	}
	if err := d.Set("next_hop_gateway", flattenComputeRouteNextHopGateway(res["nextHopGateway"], d)); err != nil {
		return fmt.Errorf("Error reading Route: %s", err)
	}
	if err := d.Set("next_hop_instance", flattenComputeRouteNextHopInstance(res["nextHopInstance"], d)); err != nil {
		return fmt.Errorf("Error reading Route: %s", err)
	}
	if err := d.Set("next_hop_ip", flattenComputeRouteNextHopIp(res["nextHopIp"], d)); err != nil {
		return fmt.Errorf("Error reading Route: %s", err)
	}
	if err := d.Set("next_hop_vpn_tunnel", flattenComputeRouteNextHopVpnTunnel(res["nextHopVpnTunnel"], d)); err != nil {
		return fmt.Errorf("Error reading Route: %s", err)
	}
	if err := d.Set("next_hop_network", flattenComputeRouteNextHopNetwork(res["nextHopNetwork"], d)); err != nil {
		return fmt.Errorf("Error reading Route: %s", err)
	}
	if err := d.Set("self_link", ConvertSelfLinkToV1(res["selfLink"].(string))); err != nil {
		return fmt.Errorf("Error reading Route: %s", err)
	}

	return nil
}

func resourceComputeRouteDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	url, err := replaceVars(d, config, "https://www.googleapis.com/compute/v1/projects/{{project}}/global/routes/{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting Route %q", d.Id())
	res, err := sendRequestWithTimeout(config, "DELETE", url, obj, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return handleNotFoundError(err, d, "Route")
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	op := &compute.Operation{}
	err = Convert(res, op)
	if err != nil {
		return err
	}

	err = computeOperationWaitTime(
		config.clientCompute, op, project, "Deleting Route",
		int(d.Timeout(schema.TimeoutDelete).Minutes()))

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting Route %q: %#v", d.Id(), res)
	return nil
}

func resourceComputeRouteImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{"projects/(?P<project>[^/]+)/global/routes/(?P<name>[^/]+)", "(?P<project>[^/]+)/(?P<name>[^/]+)", "(?P<name>[^/]+)"}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenComputeRouteDestRange(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenComputeRouteDescription(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenComputeRouteName(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenComputeRouteNetwork(v interface{}, d *schema.ResourceData) interface{} {
	if v == nil {
		return v
	}
	return ConvertSelfLinkToV1(v.(string))
}

func flattenComputeRoutePriority(v interface{}, d *schema.ResourceData) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := strconv.ParseInt(strVal, 10, 64); err == nil {
			return intVal
		} // let terraform core handle it if we can't convert the string to an int.
	}
	return v
}

func flattenComputeRouteTags(v interface{}, d *schema.ResourceData) interface{} {
	if v == nil {
		return v
	}
	return schema.NewSet(schema.HashString, v.([]interface{}))
}

func flattenComputeRouteNextHopGateway(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenComputeRouteNextHopInstance(v interface{}, d *schema.ResourceData) interface{} {
	if v == nil {
		return v
	}
	return ConvertSelfLinkToV1(v.(string))
}

func flattenComputeRouteNextHopIp(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenComputeRouteNextHopVpnTunnel(v interface{}, d *schema.ResourceData) interface{} {
	if v == nil {
		return v
	}
	return ConvertSelfLinkToV1(v.(string))
}

func flattenComputeRouteNextHopNetwork(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandComputeRouteDestRange(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeRouteDescription(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeRouteName(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeRouteNetwork(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	f, err := parseGlobalFieldValue("networks", v.(string), "project", d, config, true)
	if err != nil {
		return nil, fmt.Errorf("Invalid value for network: %s", err)
	}
	return f.RelativeLink(), nil
}

func expandComputeRoutePriority(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeRouteTags(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	return v.(*schema.Set).List(), nil
}

func expandComputeRouteNextHopGateway(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	if v == "default-internet-gateway" {
		project, err := getProject(d, config)
		if err != nil {
			return nil, err
		}
		return fmt.Sprintf("projects/%s/global/gateways/default-internet-gateway", project), nil
	} else {
		return v, nil
	}
}

func expandComputeRouteNextHopInstance(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	if v == "" {
		return v, nil
	}
	val, err := parseZonalFieldValue("instances", v.(string), "project", "next_hop_instance_zone", d, config, true)
	if err != nil {
		return nil, err
	}
	nextInstance, err := config.clientCompute.Instances.Get(val.Project, val.Zone, val.Name).Do()
	if err != nil {
		return nil, err
	}
	return nextInstance.SelfLink, nil
}

func expandComputeRouteNextHopIp(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeRouteNextHopVpnTunnel(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	f, err := parseRegionalFieldValue("vpnTunnels", v.(string), "project", "region", "zone", d, config, true)
	if err != nil {
		return nil, fmt.Errorf("Invalid value for next_hop_vpn_tunnel: %s", err)
	}
	return f.RelativeLink(), nil
}

func resourceComputeRouteDecoder(d *schema.ResourceData, meta interface{}, res map[string]interface{}) (map[string]interface{}, error) {
	if v, ok := res["nextHopInstance"]; ok {
		val, err := parseZonalFieldValue("instances", v.(string), "project", "next_hop_instance_zone", d, meta.(*Config), true)
		if err != nil {
			return nil, err
		}
		d.Set("next_hop_instance_zone", val.Zone)
		res["nextHopInstance"] = val.RelativeLink()
	}

	return res, nil
}
