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
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"google.golang.org/api/compute/v1"
)

func resourceComputeAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeAddressCreate,
		Read:   resourceComputeAddressRead,
		Delete: resourceComputeAddressDelete,

		Importer: &schema.ResourceImporter{
			State: resourceComputeAddressImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(240 * time.Second),
			Delete: schema.DefaultTimeout(240 * time.Second),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateRegexp(`^(?:[a-z](?:[-a-z0-9]{0,61}[a-z0-9])?)$`),
			},
			"address": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
			"address_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"INTERNAL", "EXTERNAL", ""}, false),
				Default:      "EXTERNAL",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"network_tier": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PREMIUM", "STANDARD", ""}, false),
			},
			"region": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
			},
			"subnetwork": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
			},
			"creation_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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

func resourceComputeAddressCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	obj := make(map[string]interface{})
	addressProp, err := expandComputeAddressAddress(d.Get("address"), d, config)
	if err != nil {
		return err
	} else if !isEmptyValue(reflect.ValueOf(addressProp)) {
		obj["address"] = addressProp
	}
	addressTypeProp, err := expandComputeAddressAddressType(d.Get("address_type"), d, config)
	if err != nil {
		return err
	} else if !isEmptyValue(reflect.ValueOf(addressTypeProp)) {
		obj["addressType"] = addressTypeProp
	}
	descriptionProp, err := expandComputeAddressDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if !isEmptyValue(reflect.ValueOf(descriptionProp)) {
		obj["description"] = descriptionProp
	}
	nameProp, err := expandComputeAddressName(d.Get("name"), d, config)
	if err != nil {
		return err
	} else if !isEmptyValue(reflect.ValueOf(nameProp)) {
		obj["name"] = nameProp
	}
	networkTierProp, err := expandComputeAddressNetworkTier(d.Get("network_tier"), d, config)
	if err != nil {
		return err
	} else if !isEmptyValue(reflect.ValueOf(networkTierProp)) {
		obj["networkTier"] = networkTierProp
	}
	subnetworkProp, err := expandComputeAddressSubnetwork(d.Get("subnetwork"), d, config)
	if err != nil {
		return err
	} else if !isEmptyValue(reflect.ValueOf(subnetworkProp)) {
		obj["subnetwork"] = subnetworkProp
	}
	regionProp, err := expandComputeAddressRegion(d.Get("region"), d, config)
	if err != nil {
		return err
	} else if !isEmptyValue(reflect.ValueOf(regionProp)) {
		obj["region"] = regionProp
	}

	url, err := replaceVars(d, config, "https://www.googleapis.com/compute/v1/projects/{{project}}/regions/{{region}}/addresses")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new Address: %#v", obj)
	res, err := sendRequestWithTimeout(config, "POST", url, obj, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating Address: %s", err)
	}

	// Store the ID now
	id, err := replaceVars(d, config, "{{project}}/{{region}}/{{name}}")
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
		config.clientCompute, op, project, "Creating Address",
		int(d.Timeout(schema.TimeoutCreate).Minutes()))

	if waitErr != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error waiting to create Address: %s", waitErr)
	}

	log.Printf("[DEBUG] Finished creating Address %q: %#v", d.Id(), res)

	return resourceComputeAddressRead(d, meta)
}

func resourceComputeAddressRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	url, err := replaceVars(d, config, "https://www.googleapis.com/compute/v1/projects/{{project}}/regions/{{region}}/addresses/{{name}}")
	if err != nil {
		return err
	}

	res, err := sendRequest(config, "GET", url, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("ComputeAddress %q", d.Id()))
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading Address: %s", err)
	}

	if err := d.Set("address", flattenComputeAddressAddress(res["address"], d)); err != nil {
		return fmt.Errorf("Error reading Address: %s", err)
	}
	if err := d.Set("address_type", flattenComputeAddressAddressType(res["addressType"], d)); err != nil {
		return fmt.Errorf("Error reading Address: %s", err)
	}
	if err := d.Set("creation_timestamp", flattenComputeAddressCreationTimestamp(res["creationTimestamp"], d)); err != nil {
		return fmt.Errorf("Error reading Address: %s", err)
	}
	if err := d.Set("description", flattenComputeAddressDescription(res["description"], d)); err != nil {
		return fmt.Errorf("Error reading Address: %s", err)
	}
	if err := d.Set("name", flattenComputeAddressName(res["name"], d)); err != nil {
		return fmt.Errorf("Error reading Address: %s", err)
	}
	if err := d.Set("network_tier", flattenComputeAddressNetworkTier(res["networkTier"], d)); err != nil {
		return fmt.Errorf("Error reading Address: %s", err)
	}
	if err := d.Set("subnetwork", flattenComputeAddressSubnetwork(res["subnetwork"], d)); err != nil {
		return fmt.Errorf("Error reading Address: %s", err)
	}
	if err := d.Set("users", flattenComputeAddressUsers(res["users"], d)); err != nil {
		return fmt.Errorf("Error reading Address: %s", err)
	}
	if err := d.Set("region", flattenComputeAddressRegion(res["region"], d)); err != nil {
		return fmt.Errorf("Error reading Address: %s", err)
	}
	if err := d.Set("self_link", ConvertSelfLinkToV1(res["selfLink"].(string))); err != nil {
		return fmt.Errorf("Error reading Address: %s", err)
	}

	return nil
}

func resourceComputeAddressDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	url, err := replaceVars(d, config, "https://www.googleapis.com/compute/v1/projects/{{project}}/regions/{{region}}/addresses/{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting Address %q", d.Id())
	res, err := sendRequestWithTimeout(config, "DELETE", url, obj, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return handleNotFoundError(err, d, "Address")
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
		config.clientCompute, op, project, "Deleting Address",
		int(d.Timeout(schema.TimeoutDelete).Minutes()))

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting Address %q: %#v", d.Id(), res)
	return nil
}

func resourceComputeAddressImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{"projects/(?P<project>[^/]+)/regions/(?P<region>[^/]+)/addresses/(?P<name>[^/]+)", "(?P<project>[^/]+)/(?P<region>[^/]+)/(?P<name>[^/]+)", "(?P<name>[^/]+)"}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "{{project}}/{{region}}/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenComputeAddressAddress(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenComputeAddressAddressType(v interface{}, d *schema.ResourceData) interface{} {
	if v == nil || v.(string) == "" {
		return "EXTERNAL"
	}
	return v
}

func flattenComputeAddressCreationTimestamp(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenComputeAddressDescription(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenComputeAddressName(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenComputeAddressNetworkTier(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenComputeAddressSubnetwork(v interface{}, d *schema.ResourceData) interface{} {
	if v == nil {
		return v
	}
	return ConvertSelfLinkToV1(v.(string))
}

func flattenComputeAddressUsers(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenComputeAddressRegion(v interface{}, d *schema.ResourceData) interface{} {
	if v == nil {
		return v
	}
	return NameFromSelfLinkStateFunc(v)
}

func expandComputeAddressAddress(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeAddressAddressType(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeAddressDescription(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeAddressName(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeAddressNetworkTier(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeAddressSubnetwork(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	f, err := parseRegionalFieldValue("subnetworks", v.(string), "project", "region", "zone", d, config, true)
	if err != nil {
		return nil, fmt.Errorf("Invalid value for subnetwork: %s", err)
	}
	return f.RelativeLink(), nil
}

func expandComputeAddressRegion(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	f, err := parseGlobalFieldValue("regions", v.(string), "project", d, config, true)
	if err != nil {
		return nil, fmt.Errorf("Invalid value for region: %s", err)
	}
	return f.RelativeLink(), nil
}
