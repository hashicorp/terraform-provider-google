// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package tags

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTagsLocationTagBinding() *schema.Resource {
	return &schema.Resource{
		Create: resourceTagsLocationTagBindingCreate,
		Read:   resourceTagsLocationTagBindingRead,
		Delete: resourceTagsLocationTagBindingDelete,

		Importer: &schema.ResourceImporter{
			State: resourceTagsLocationTagBindingImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"parent": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The full resource name of the resource the TagValue is bound to. E.g. //cloudresourcemanager.googleapis.com/projects/123`,
			},
			"tag_value": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The TagValue of the TagBinding. Must be of the form tagValues/456.`,
			},
			"location": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Description: `The geographic location where the transfer config should reside.
Examples: US, EU, asia-northeast1. The default value is US.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The generated id for the TagBinding. This is a string of the form: 'tagBindings/{full-resource-name}/{tag-value-name}'`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceTagsLocationTagBindingCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	parentProp, err := expandNestedTagsLocationTagBindingParent(d.Get("parent"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("parent"); !tpgresource.IsEmptyValue(reflect.ValueOf(parentProp)) && (ok || !reflect.DeepEqual(v, parentProp)) {
		obj["parent"] = parentProp
	}
	tagValueProp, err := expandNestedTagsLocationTagBindingTagValue(d.Get("tag_value"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("tag_value"); !tpgresource.IsEmptyValue(reflect.ValueOf(tagValueProp)) && (ok || !reflect.DeepEqual(v, tagValueProp)) {
		obj["tagValue"] = tagValueProp
	}

	lockName, err := tpgresource.ReplaceVars(d, config, "tagBindings/{{parent}}")
	if err != nil {
		return err
	}
	transport_tpg.MutexStore.Lock(lockName)
	defer transport_tpg.MutexStore.Unlock(lockName)

	url, err := tpgresource.ReplaceVars(d, config, "{{TagsLocationBasePath}}tagBindings")
	log.Printf("url for TagsLocation: %s", url)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new LocationTagBinding: %#v", obj)

	billingProject := ""

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
		return fmt.Errorf("Error creating LocationTagBinding: %s", err)
	}

	// Use the resource in the operation response to populate
	// identity fields and d.Id() before read

	var opRes map[string]interface{}
	err = TagsLocationOperationWaitTimeWithResponse(
		config, res, &opRes, "Creating LocationTagBinding", userAgent,
		d.Timeout(schema.TimeoutCreate))

	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error waiting to create LocationTagBinding: %s", err)
	}

	if _, ok := opRes["tagBindings"]; ok {
		opRes, err = flattenNestedTagsLocationTagBinding(d, meta, opRes)
		if err != nil {
			return fmt.Errorf("Error getting nested object from operation response: %s", err)
		}
		if opRes == nil {
			// Object isn't there any more - remove it from the state.
			d.SetId("")
			return fmt.Errorf("Error decoding response from operation, could not find nested object")
		}
	}
	if err := d.Set("name", flattenNestedTagsLocationTagBindingName(opRes["name"], d, config)); err != nil {
		return err
	}

	id, err := tpgresource.ReplaceVars(d, config, "{{location}}/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating LocationTagBinding %q: %#v", d.Id(), res)

	return resourceTagsLocationTagBindingRead(d, meta)
}

func resourceTagsLocationTagBindingRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{TagsLocationBasePath}}tagBindings/?parent={{parent}}&pageSize=300")
	if err != nil {
		return err
	}

	billingProject := ""

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
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("TagsLocationTagBinding %q", d.Id()))
	}
	log.Printf("[DEBUG] Skipping res with name for import = %#v,)", res)

	p, ok := res["tagBindings"]
	if !ok || p == nil {
		return nil
	}
	pView := p.([]interface{})

	//if there are more than 300 bindings - handling pagination over here
	if pageToken, ok := res["nextPageToken"].(string); ok {
		for pageToken != "" {
			url, err = transport_tpg.AddQueryParams(url, map[string]string{"pageToken": fmt.Sprintf("%s", res["nextPageToken"])})
			if err != nil {
				return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("TagsLocationTagBinding %q", d.Id()))
			}
			resp, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: userAgent,
			})
			if err != nil {
				return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("TagsLocationTagBinding %q", d.Id()))
			}
			if resp == nil {
				d.SetId("")
				return nil
			}
			v, ok := resp["tagBindings"]
			if !ok || v == nil {
				return nil
			}
			pView = append(pView, v.([]interface{})...)
			if token, ok := res["nextPageToken"]; ok {
				pageToken = token.(string)
			} else {
				pageToken = ""
			}
		}
	}

	newMap := make(map[string]interface{}, 1)
	newMap["tagBindings"] = pView

	res, err = flattenNestedTagsLocationTagBinding(d, meta, newMap)
	if err != nil {
		return err
	}

	if err := d.Set("name", flattenNestedTagsLocationTagBindingName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading LocationTagBinding: %s", err)
	}
	if err := d.Set("parent", flattenNestedTagsLocationTagBindingParent(res["parent"], d, config)); err != nil {
		return fmt.Errorf("Error reading LocationTagBinding: %s", err)
	}
	if err := d.Set("tag_value", flattenNestedTagsLocationTagBindingTagValue(res["tagValue"], d, config)); err != nil {
		return fmt.Errorf("Error reading LocationTagBinding: %s", err)
	}

	return nil
}

func resourceTagsLocationTagBindingDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	lockName, err := tpgresource.ReplaceVars(d, config, "tagBindings/{{parent}}")
	if err != nil {
		return err
	}
	transport_tpg.MutexStore.Lock(lockName)
	defer transport_tpg.MutexStore.Unlock(lockName)

	url, err := tpgresource.ReplaceVars(d, config, "{{TagsLocationBasePath}}{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting LocationTagBinding %q", d.Id())

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
		return transport_tpg.HandleNotFoundError(err, d, "LocationTagBinding")
	}

	err = TagsLocationOperationWaitTime(
		config, res, "Deleting LocationTagBinding", userAgent,
		d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting LocationTagBinding %q: %#v", d.Id(), res)
	return nil
}

func resourceTagsLocationTagBindingImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{"(?P<location>[^/]+)/tagBindings/(?P<parent>[^/]+)/tagValues/(?P<tag_value>[^/]+)"}, d, config); err != nil {
		return nil, err
	}

	parent := d.Get("parent").(string)
	parentProper := strings.ReplaceAll(parent, "%2F", "/")
	d.Set("parent", parentProper)
	d.Set("name", fmt.Sprintf("tagBindings/%s/tagValues/%s", parent, d.Get("tag_value").(string)))
	id, err := tpgresource.ReplaceVars(d, config, "{{location}}/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenNestedTagsLocationTagBindingName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNestedTagsLocationTagBindingParent(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNestedTagsLocationTagBindingTagValue(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandNestedTagsLocationTagBindingParent(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNestedTagsLocationTagBindingTagValue(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func flattenNestedTagsLocationTagBinding(d *schema.ResourceData, meta interface{}, res map[string]interface{}) (map[string]interface{}, error) {
	var v interface{}
	var ok bool

	v, ok = res["tagBindings"]
	if !ok || v == nil {
		return nil, nil
	}

	switch v.(type) {
	case []interface{}:
		log.Printf("[DEBUG] Hey it's in break = %#v,)", v)
		break
	case map[string]interface{}:
		// Construct list out of single nested resource
		v = []interface{}{v}
	default:
		return nil, fmt.Errorf("expected list or map for value tagBindings. Actual value: %v", v)
	}

	_, item, err := resourceTagsLocationTagBindingFindNestedObjectInList(d, meta, v.([]interface{}))
	if err != nil {
		return nil, err
	}
	return item, nil
}

func resourceTagsLocationTagBindingFindNestedObjectInList(d *schema.ResourceData, meta interface{}, items []interface{}) (index int, item map[string]interface{}, err error) {
	expectedName := d.Get("name")
	expectedFlattenedName := flattenNestedTagsLocationTagBindingName(expectedName, d, meta.(*transport_tpg.Config))

	// Search list for this resource.
	for idx, itemRaw := range items {
		if itemRaw == nil {
			continue
		}

		item := itemRaw.(map[string]interface{})
		itemName := flattenNestedTagsLocationTagBindingName(item["name"], d, meta.(*transport_tpg.Config))
		// IsEmptyValue check so that if one is nil and the other is "", that's considered a match
		if !(tpgresource.IsEmptyValue(reflect.ValueOf(itemName)) && tpgresource.IsEmptyValue(reflect.ValueOf(expectedFlattenedName))) && !reflect.DeepEqual(itemName, expectedFlattenedName) {
			log.Printf("[DEBUG] Skipping item with name= %#v, looking for %#v)", itemName, expectedFlattenedName)
			continue
		}
		return idx, item, nil
	}
	return -1, nil, nil
}
