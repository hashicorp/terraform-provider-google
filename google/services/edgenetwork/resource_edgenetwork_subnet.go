// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
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

package edgenetwork

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceEdgenetworkSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceEdgenetworkSubnetCreate,
		Read:   resourceEdgenetworkSubnetRead,
		Delete: resourceEdgenetworkSubnetDelete,

		Importer: &schema.ResourceImporter{
			State: resourceEdgenetworkSubnetImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The Google Cloud region to which the target Distributed Cloud Edge zone belongs.`,
			},
			"network": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description: `The ID of the network to which this router belongs.
Must be of the form: 'projects/{{project}}/locations/{{location}}/zones/{{zone}}/networks/{{network_id}}'`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `A unique ID that identifies this subnet.`,
			},
			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the target Distributed Cloud Edge zone.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `A free-text description of the resource. Max length 1024 characters.`,
			},
			"ipv4_cidr": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: `The ranges of ipv4 addresses that are owned by this subnetwork, in CIDR format.`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ipv6_cidr": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: `The ranges of ipv6 addresses that are owned by this subnetwork, in CIDR format.`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: `Labels associated with this resource.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"vlan_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: `VLAN ID for this subnetwork. If not specified, one is assigned automatically.`,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The time when the subnet was created.
A timestamp in RFC3339 UTC "Zulu" format, with nanosecond resolution and up to nine
fractional digits. Examples: '2014-10-02T15:01:23Z' and '2014-10-02T15:01:23.045123456Z'.`,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The canonical name of this resource, with format
'projects/{{project}}/locations/{{location}}/zones/{{zone}}/subnets/{{subnet_id}}'`,
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Current stage of the resource to the device by config push.`,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The time when the subnet was last updated.
A timestamp in RFC3339 UTC "Zulu" format, with nanosecond resolution and up to nine
fractional digits. Examples: '2014-10-02T15:01:23Z' and '2014-10-02T15:01:23.045123456Z'.`,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceEdgenetworkSubnetCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	labelsProp, err := expandEdgenetworkSubnetLabels(d.Get("labels"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("labels"); !tpgresource.IsEmptyValue(reflect.ValueOf(labelsProp)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}
	descriptionProp, err := expandEdgenetworkSubnetDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	networkProp, err := expandEdgenetworkSubnetNetwork(d.Get("network"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("network"); !tpgresource.IsEmptyValue(reflect.ValueOf(networkProp)) && (ok || !reflect.DeepEqual(v, networkProp)) {
		obj["network"] = networkProp
	}
	ipv4CidrProp, err := expandEdgenetworkSubnetIpv4Cidr(d.Get("ipv4_cidr"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("ipv4_cidr"); !tpgresource.IsEmptyValue(reflect.ValueOf(ipv4CidrProp)) && (ok || !reflect.DeepEqual(v, ipv4CidrProp)) {
		obj["ipv4Cidr"] = ipv4CidrProp
	}
	ipv6CidrProp, err := expandEdgenetworkSubnetIpv6Cidr(d.Get("ipv6_cidr"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("ipv6_cidr"); !tpgresource.IsEmptyValue(reflect.ValueOf(ipv6CidrProp)) && (ok || !reflect.DeepEqual(v, ipv6CidrProp)) {
		obj["ipv6Cidr"] = ipv6CidrProp
	}
	vlanIdProp, err := expandEdgenetworkSubnetVlanId(d.Get("vlan_id"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("vlan_id"); !tpgresource.IsEmptyValue(reflect.ValueOf(vlanIdProp)) && (ok || !reflect.DeepEqual(v, vlanIdProp)) {
		obj["vlanId"] = vlanIdProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{EdgenetworkBasePath}}projects/{{project}}/locations/{{location}}/zones/{{zone}}/subnets?subnetId={{subnet_id}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new Subnet: %#v", obj)
	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Subnet: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutCreate),
	})
	if err != nil {
		return fmt.Errorf("Error creating Subnet: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/zones/{{zone}}/subnets/{{subnet_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = EdgenetworkOperationWaitTime(
		config, res, project, "Creating Subnet", userAgent,
		d.Timeout(schema.TimeoutCreate))

	if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error waiting to create Subnet: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Subnet %q: %#v", d.Id(), res)

	return resourceEdgenetworkSubnetRead(d, meta)
}

func resourceEdgenetworkSubnetRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{EdgenetworkBasePath}}projects/{{project}}/locations/{{location}}/zones/{{zone}}/subnets/{{subnet_id}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Subnet: %s", err)
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
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("EdgenetworkSubnet %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading Subnet: %s", err)
	}

	if err := d.Set("name", flattenEdgenetworkSubnetName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading Subnet: %s", err)
	}
	if err := d.Set("labels", flattenEdgenetworkSubnetLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading Subnet: %s", err)
	}
	if err := d.Set("description", flattenEdgenetworkSubnetDescription(res["description"], d, config)); err != nil {
		return fmt.Errorf("Error reading Subnet: %s", err)
	}
	if err := d.Set("create_time", flattenEdgenetworkSubnetCreateTime(res["createTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading Subnet: %s", err)
	}
	if err := d.Set("update_time", flattenEdgenetworkSubnetUpdateTime(res["updateTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading Subnet: %s", err)
	}
	if err := d.Set("network", flattenEdgenetworkSubnetNetwork(res["network"], d, config)); err != nil {
		return fmt.Errorf("Error reading Subnet: %s", err)
	}
	if err := d.Set("ipv4_cidr", flattenEdgenetworkSubnetIpv4Cidr(res["ipv4Cidr"], d, config)); err != nil {
		return fmt.Errorf("Error reading Subnet: %s", err)
	}
	if err := d.Set("ipv6_cidr", flattenEdgenetworkSubnetIpv6Cidr(res["ipv6Cidr"], d, config)); err != nil {
		return fmt.Errorf("Error reading Subnet: %s", err)
	}
	if err := d.Set("vlan_id", flattenEdgenetworkSubnetVlanId(res["vlanId"], d, config)); err != nil {
		return fmt.Errorf("Error reading Subnet: %s", err)
	}
	if err := d.Set("state", flattenEdgenetworkSubnetState(res["state"], d, config)); err != nil {
		return fmt.Errorf("Error reading Subnet: %s", err)
	}

	return nil
}

func resourceEdgenetworkSubnetDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Subnet: %s", err)
	}
	billingProject = project

	url, err := tpgresource.ReplaceVars(d, config, "{{EdgenetworkBasePath}}projects/{{project}}/locations/{{location}}/zones/{{zone}}/subnets/{{subnet_id}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting Subnet %q", d.Id())

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "DELETE",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutDelete),
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, "Subnet")
	}

	err = EdgenetworkOperationWaitTime(
		config, res, project, "Deleting Subnet", userAgent,
		d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting Subnet %q: %#v", d.Id(), res)
	return nil
}

func resourceEdgenetworkSubnetImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/zones/(?P<zone>[^/]+)/subnets/(?P<subnet_id>[^/]+)",
		"(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<zone>[^/]+)/(?P<subnet_id>[^/]+)",
		"(?P<location>[^/]+)/(?P<zone>[^/]+)/(?P<subnet_id>[^/]+)",
		"(?P<location>[^/]+)/(?P<subnet_id>[^/]+)",
		"(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/zones/{{zone}}/subnets/{{subnet_id}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenEdgenetworkSubnetName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenEdgenetworkSubnetLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenEdgenetworkSubnetDescription(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenEdgenetworkSubnetCreateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenEdgenetworkSubnetUpdateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenEdgenetworkSubnetNetwork(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	return tpgresource.ConvertSelfLinkToV1(v.(string))
}

func flattenEdgenetworkSubnetIpv4Cidr(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenEdgenetworkSubnetIpv6Cidr(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenEdgenetworkSubnetVlanId(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenEdgenetworkSubnetState(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandEdgenetworkSubnetLabels(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandEdgenetworkSubnetDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandEdgenetworkSubnetNetwork(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	f, err := tpgresource.ParseZonalFieldValue("networks", v.(string), "project", "zone", d, config, true)
	if err != nil {
		return nil, fmt.Errorf("Invalid value for network: %s", err)
	}
	return f.RelativeLink(), nil
}

func expandEdgenetworkSubnetIpv4Cidr(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandEdgenetworkSubnetIpv6Cidr(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandEdgenetworkSubnetVlanId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
