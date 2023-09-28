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

package bigqueryanalyticshub

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceBigqueryAnalyticsHubDataExchange() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigqueryAnalyticsHubDataExchangeCreate,
		Read:   resourceBigqueryAnalyticsHubDataExchangeRead,
		Update: resourceBigqueryAnalyticsHubDataExchangeUpdate,
		Delete: resourceBigqueryAnalyticsHubDataExchangeDelete,

		Importer: &schema.ResourceImporter{
			State: resourceBigqueryAnalyticsHubDataExchangeImport,
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
			"data_exchange_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the data exchange. Must contain only Unicode letters, numbers (0-9), underscores (_). Should not use characters that require URL-escaping, or characters outside of ASCII, spaces.`,
			},
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Human-readable display name of the data exchange. The display name must contain only Unicode letters, numbers (0-9), underscores (_), dashes (-), spaces ( ), and must not start or end with spaces.`,
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the location this data exchange.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Description of the data exchange.`,
			},
			"documentation": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Documentation describing the data exchange.`,
			},
			"icon": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Base64 encoded image representing the data exchange.`,
			},
			"primary_contact": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Email or URL of the primary point of contact of the data exchange.`,
			},
			"listing_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Number of listings contained in the data exchange.`,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The resource name of the data exchange, for example:
"projects/myproject/locations/US/dataExchanges/123"`,
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

func resourceBigqueryAnalyticsHubDataExchangeCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	displayNameProp, err := expandBigqueryAnalyticsHubDataExchangeDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("display_name"); !tpgresource.IsEmptyValue(reflect.ValueOf(displayNameProp)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}
	descriptionProp, err := expandBigqueryAnalyticsHubDataExchangeDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	primaryContactProp, err := expandBigqueryAnalyticsHubDataExchangePrimaryContact(d.Get("primary_contact"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("primary_contact"); !tpgresource.IsEmptyValue(reflect.ValueOf(primaryContactProp)) && (ok || !reflect.DeepEqual(v, primaryContactProp)) {
		obj["primaryContact"] = primaryContactProp
	}
	documentationProp, err := expandBigqueryAnalyticsHubDataExchangeDocumentation(d.Get("documentation"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("documentation"); !tpgresource.IsEmptyValue(reflect.ValueOf(documentationProp)) && (ok || !reflect.DeepEqual(v, documentationProp)) {
		obj["documentation"] = documentationProp
	}
	iconProp, err := expandBigqueryAnalyticsHubDataExchangeIcon(d.Get("icon"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("icon"); !tpgresource.IsEmptyValue(reflect.ValueOf(iconProp)) && (ok || !reflect.DeepEqual(v, iconProp)) {
		obj["icon"] = iconProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{BigqueryAnalyticsHubBasePath}}projects/{{project}}/locations/{{location}}/dataExchanges?data_exchange_id={{data_exchange_id}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new DataExchange: %#v", obj)
	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for DataExchange: %s", err)
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
		return fmt.Errorf("Error creating DataExchange: %s", err)
	}
	if err := d.Set("name", flattenBigqueryAnalyticsHubDataExchangeName(res["name"], d, config)); err != nil {
		return fmt.Errorf(`Error setting computed identity field "name": %s`, err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/dataExchanges/{{data_exchange_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating DataExchange %q: %#v", d.Id(), res)

	return resourceBigqueryAnalyticsHubDataExchangeRead(d, meta)
}

func resourceBigqueryAnalyticsHubDataExchangeRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{BigqueryAnalyticsHubBasePath}}projects/{{project}}/locations/{{location}}/dataExchanges/{{data_exchange_id}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for DataExchange: %s", err)
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
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("BigqueryAnalyticsHubDataExchange %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading DataExchange: %s", err)
	}

	if err := d.Set("name", flattenBigqueryAnalyticsHubDataExchangeName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading DataExchange: %s", err)
	}
	if err := d.Set("display_name", flattenBigqueryAnalyticsHubDataExchangeDisplayName(res["displayName"], d, config)); err != nil {
		return fmt.Errorf("Error reading DataExchange: %s", err)
	}
	if err := d.Set("description", flattenBigqueryAnalyticsHubDataExchangeDescription(res["description"], d, config)); err != nil {
		return fmt.Errorf("Error reading DataExchange: %s", err)
	}
	if err := d.Set("primary_contact", flattenBigqueryAnalyticsHubDataExchangePrimaryContact(res["primaryContact"], d, config)); err != nil {
		return fmt.Errorf("Error reading DataExchange: %s", err)
	}
	if err := d.Set("documentation", flattenBigqueryAnalyticsHubDataExchangeDocumentation(res["documentation"], d, config)); err != nil {
		return fmt.Errorf("Error reading DataExchange: %s", err)
	}
	if err := d.Set("listing_count", flattenBigqueryAnalyticsHubDataExchangeListingCount(res["listingCount"], d, config)); err != nil {
		return fmt.Errorf("Error reading DataExchange: %s", err)
	}
	if err := d.Set("icon", flattenBigqueryAnalyticsHubDataExchangeIcon(res["icon"], d, config)); err != nil {
		return fmt.Errorf("Error reading DataExchange: %s", err)
	}

	return nil
}

func resourceBigqueryAnalyticsHubDataExchangeUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for DataExchange: %s", err)
	}
	billingProject = project

	obj := make(map[string]interface{})
	displayNameProp, err := expandBigqueryAnalyticsHubDataExchangeDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("display_name"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}
	descriptionProp, err := expandBigqueryAnalyticsHubDataExchangeDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	primaryContactProp, err := expandBigqueryAnalyticsHubDataExchangePrimaryContact(d.Get("primary_contact"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("primary_contact"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, primaryContactProp)) {
		obj["primaryContact"] = primaryContactProp
	}
	documentationProp, err := expandBigqueryAnalyticsHubDataExchangeDocumentation(d.Get("documentation"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("documentation"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, documentationProp)) {
		obj["documentation"] = documentationProp
	}
	iconProp, err := expandBigqueryAnalyticsHubDataExchangeIcon(d.Get("icon"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("icon"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, iconProp)) {
		obj["icon"] = iconProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{BigqueryAnalyticsHubBasePath}}projects/{{project}}/locations/{{location}}/dataExchanges/{{data_exchange_id}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating DataExchange %q: %#v", d.Id(), obj)
	updateMask := []string{}

	if d.HasChange("display_name") {
		updateMask = append(updateMask, "displayName")
	}

	if d.HasChange("description") {
		updateMask = append(updateMask, "description")
	}

	if d.HasChange("primary_contact") {
		updateMask = append(updateMask, "primaryContact")
	}

	if d.HasChange("documentation") {
		updateMask = append(updateMask, "documentation")
	}

	if d.HasChange("icon") {
		updateMask = append(updateMask, "icon")
	}
	// updateMask is a URL parameter but not present in the schema, so ReplaceVars
	// won't set it
	url, err = transport_tpg.AddQueryParams(url, map[string]string{"updateMask": strings.Join(updateMask, ",")})
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
		return fmt.Errorf("Error updating DataExchange %q: %s", d.Id(), err)
	} else {
		log.Printf("[DEBUG] Finished updating DataExchange %q: %#v", d.Id(), res)
	}

	return resourceBigqueryAnalyticsHubDataExchangeRead(d, meta)
}

func resourceBigqueryAnalyticsHubDataExchangeDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for DataExchange: %s", err)
	}
	billingProject = project

	url, err := tpgresource.ReplaceVars(d, config, "{{BigqueryAnalyticsHubBasePath}}projects/{{project}}/locations/{{location}}/dataExchanges/{{data_exchange_id}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting DataExchange %q", d.Id())

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
		return transport_tpg.HandleNotFoundError(err, d, "DataExchange")
	}

	log.Printf("[DEBUG] Finished deleting DataExchange %q: %#v", d.Id(), res)
	return nil
}

func resourceBigqueryAnalyticsHubDataExchangeImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"^projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/dataExchanges/(?P<data_exchange_id>[^/]+)$",
		"^(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<data_exchange_id>[^/]+)$",
		"^(?P<location>[^/]+)/(?P<data_exchange_id>[^/]+)$",
		"^(?P<data_exchange_id>[^/]+)$",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/dataExchanges/{{data_exchange_id}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenBigqueryAnalyticsHubDataExchangeName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenBigqueryAnalyticsHubDataExchangeDisplayName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenBigqueryAnalyticsHubDataExchangeDescription(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenBigqueryAnalyticsHubDataExchangePrimaryContact(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenBigqueryAnalyticsHubDataExchangeDocumentation(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenBigqueryAnalyticsHubDataExchangeListingCount(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenBigqueryAnalyticsHubDataExchangeIcon(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandBigqueryAnalyticsHubDataExchangeDisplayName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandBigqueryAnalyticsHubDataExchangeDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandBigqueryAnalyticsHubDataExchangePrimaryContact(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandBigqueryAnalyticsHubDataExchangeDocumentation(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandBigqueryAnalyticsHubDataExchangeIcon(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
