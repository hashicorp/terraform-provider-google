// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This code is generated by Magic Modules using the following:
//
//     Configuration: https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products/compute/ResourcePolicyAttachment.yaml
//     Template:      https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/resource.go.tmpl
//
//     DO NOT EDIT this file directly. Any changes made to this file will be
//     overwritten during the next generation cycle.
//
// ----------------------------------------------------------------------------

package compute

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceComputeResourcePolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeResourcePolicyAttachmentCreate,
		Read:   resourceComputeResourcePolicyAttachmentRead,
		Delete: resourceComputeResourcePolicyAttachmentDelete,

		Importer: &schema.ResourceImporter{
			State: resourceComputeResourcePolicyAttachmentImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			tpgresource.DefaultProviderProject,
			tpgresource.DefaultProviderZone,
		),

		Schema: map[string]*schema.Schema{
			"instance": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      `The name of the instance in which the resource policies are attached to.`,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: `The resource policy to be attached to the instance for scheduling start/stop
operations. Do not specify the self link.`,
			},
			"zone": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      `A reference to the zone where the instance resides.`,
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

func resourceComputeResourcePolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	nameProp, err := expandNestedComputeResourcePolicyAttachmentName(d.Get("name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}

	obj, err = resourceComputeResourcePolicyAttachmentEncoder(d, meta, obj)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/zones/{{zone}}/instances/{{instance}}/addResourcePolicies")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new ResourcePolicyAttachment: %#v", obj)
	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for ResourcePolicyAttachment: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutCreate),
		Headers:   headers,
	})
	if err != nil {
		return fmt.Errorf("Error creating ResourcePolicyAttachment: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "{{project}}/{{zone}}/{{instance}}/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = ComputeOperationWaitTime(
		config, res, project, "Creating ResourcePolicyAttachment", userAgent,
		d.Timeout(schema.TimeoutCreate))

	if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error waiting to create ResourcePolicyAttachment: %s", err)
	}

	log.Printf("[DEBUG] Finished creating ResourcePolicyAttachment %q: %#v", d.Id(), res)

	return resourceComputeResourcePolicyAttachmentRead(d, meta)
}

func resourceComputeResourcePolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/zones/{{zone}}/instances/{{instance}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for ResourcePolicyAttachment: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Headers:   headers,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("ComputeResourcePolicyAttachment %q", d.Id()))
	}

	res, err = flattenNestedComputeResourcePolicyAttachment(d, meta, res)
	if err != nil {
		return err
	}

	if res == nil {
		// Object isn't there any more - remove it from the state.
		log.Printf("[DEBUG] Removing ComputeResourcePolicyAttachment because it couldn't be matched.")
		d.SetId("")
		return nil
	}

	res, err = resourceComputeResourcePolicyAttachmentDecoder(d, meta, res)
	if err != nil {
		return err
	}

	if res == nil {
		// Decoding the object has resulted in it being gone. It may be marked deleted
		log.Printf("[DEBUG] Removing ComputeResourcePolicyAttachment because it no longer exists.")
		d.SetId("")
		return nil
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading ResourcePolicyAttachment: %s", err)
	}

	zone, err := tpgresource.GetZone(d, config)
	if err != nil {
		return err
	}
	if err := d.Set("zone", zone); err != nil {
		return fmt.Errorf("Error reading ResourcePolicyAttachment: %s", err)
	}

	if err := d.Set("name", flattenNestedComputeResourcePolicyAttachmentName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading ResourcePolicyAttachment: %s", err)
	}

	return nil
}

func resourceComputeResourcePolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for ResourcePolicyAttachment: %s", err)
	}
	billingProject = project

	url, err := tpgresource.ReplaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/zones/{{zone}}/instances/{{instance}}/removeResourcePolicies")
	if err != nil {
		return err
	}

	var obj map[string]interface{}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)
	obj = make(map[string]interface{})

	zone, err := tpgresource.GetZone(d, config)
	if err != nil {
		return err
	}
	if zone == "" {
		return fmt.Errorf("zone must be non-empty - set in resource or at provider-level")
	}

	// resourcePolicies are referred to by region but affixed to zonal instances.
	// We construct the regional name from the zone:
	//
	//	projects/{project}/regions/{region}/resourcePolicies/{resourceId}
	region := tpgresource.GetRegionFromZone(zone)
	if region == "" {
		return fmt.Errorf("invalid zone %q, unable to infer region from zone", zone)
	}

	name, err := expandNestedComputeResourcePolicyAttachmentName(d.Get("name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(name)) && (ok || !reflect.DeepEqual(v, name)) {
		obj["resourcePolicies"] = []interface{}{fmt.Sprintf("projects/%s/regions/%s/resourcePolicies/%s", project, region, name)}
	}

	log.Printf("[DEBUG] Deleting ResourcePolicyAttachment %q", d.Id())
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutDelete),
		Headers:   headers,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, "ResourcePolicyAttachment")
	}

	err = ComputeOperationWaitTime(
		config, res, project, "Deleting ResourcePolicyAttachment", userAgent,
		d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting ResourcePolicyAttachment %q: %#v", d.Id(), res)
	return nil
}

func resourceComputeResourcePolicyAttachmentImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"^projects/(?P<project>[^/]+)/zones/(?P<zone>[^/]+)/instances/(?P<instance>[^/]+)/(?P<name>[^/]+)$",
		"^(?P<project>[^/]+)/(?P<zone>[^/]+)/(?P<instance>[^/]+)/(?P<name>[^/]+)$",
		"^(?P<zone>[^/]+)/(?P<instance>[^/]+)/(?P<name>[^/]+)$",
		"^(?P<instance>[^/]+)/(?P<name>[^/]+)$",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "{{project}}/{{zone}}/{{instance}}/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenNestedComputeResourcePolicyAttachmentName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandNestedComputeResourcePolicyAttachmentName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func resourceComputeResourcePolicyAttachmentEncoder(d *schema.ResourceData, meta interface{}, obj map[string]interface{}) (map[string]interface{}, error) {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return nil, err
	}

	zone, err := tpgresource.GetZone(d, config)
	if err != nil {
		return nil, err
	}
	if zone == "" {
		return nil, fmt.Errorf("zone must be non-empty - set in resource or at provider-level")
	}

	// resourcePolicies are referred to by region but affixed to zonal instances.
	// We construct the regional name from the zone:
	//
	//	projects/{project}/regions/{region}/resourcePolicies/{resourceId}
	region := tpgresource.GetRegionFromZone(zone)
	if region == "" {
		return nil, fmt.Errorf("invalid zone %q, unable to infer region from zone", zone)
	}

	obj["resourcePolicies"] = []interface{}{fmt.Sprintf("projects/%s/regions/%s/resourcePolicies/%s", project, region, obj["name"])}
	delete(obj, "name")
	return obj, nil
}

func flattenNestedComputeResourcePolicyAttachment(d *schema.ResourceData, meta interface{}, res map[string]interface{}) (map[string]interface{}, error) {
	var v interface{}
	var ok bool

	v, ok = res["resourcePolicies"]
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
		return nil, fmt.Errorf("expected list or map for value resourcePolicies. Actual value: %v", v)
	}

	_, item, err := resourceComputeResourcePolicyAttachmentFindNestedObjectInList(d, meta, v.([]interface{}))
	if err != nil {
		return nil, err
	}
	return item, nil
}

func resourceComputeResourcePolicyAttachmentFindNestedObjectInList(d *schema.ResourceData, meta interface{}, items []interface{}) (index int, item map[string]interface{}, err error) {
	expectedName, err := expandNestedComputeResourcePolicyAttachmentName(d.Get("name"), d, meta.(*transport_tpg.Config))
	if err != nil {
		return -1, nil, err
	}
	expectedFlattenedName := flattenNestedComputeResourcePolicyAttachmentName(expectedName, d, meta.(*transport_tpg.Config))

	// Search list for this resource.
	for idx, itemRaw := range items {
		if itemRaw == nil {
			continue
		}
		// List response only contains the ID - construct a response object.
		item := map[string]interface{}{
			"name": itemRaw,
		}

		// Decode list item before comparing.
		item, err := resourceComputeResourcePolicyAttachmentDecoder(d, meta, item)
		if err != nil {
			return -1, nil, err
		}

		itemName := flattenNestedComputeResourcePolicyAttachmentName(item["name"], d, meta.(*transport_tpg.Config))
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
func resourceComputeResourcePolicyAttachmentDecoder(d *schema.ResourceData, meta interface{}, res map[string]interface{}) (map[string]interface{}, error) {
	res["name"] = tpgresource.GetResourceNameFromSelfLink(res["name"].(string))
	return res, nil
}
