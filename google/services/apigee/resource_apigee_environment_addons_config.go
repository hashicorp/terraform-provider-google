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
//     Configuration: https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products/apigee/EnvironmentAddonsConfig.yaml
//     Template:      https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/resource.go.tmpl
//
//     DO NOT EDIT this file directly. Any changes made to this file will be
//     overwritten during the next generation cycle.
//
// ----------------------------------------------------------------------------

package apigee

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceApigeeEnvironmentAddonsConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceApigeeEnvironmentAddonsConfigCreate,
		Read:   resourceApigeeEnvironmentAddonsConfigRead,
		Update: resourceApigeeEnvironmentAddonsConfigUpdate,
		Delete: resourceApigeeEnvironmentAddonsConfigDelete,

		Importer: &schema.ResourceImporter{
			State: resourceApigeeEnvironmentAddonsConfigImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(0 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"env_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: `The Apigee environment group associated with the Apigee environment,
in the format 'organizations/{{org_name}}/environments/{{env_name}}'.`,
			},
			"analytics_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Flag to enable/disable Analytics.`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceApigeeEnvironmentAddonsConfigCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	analyticsEnabledProp, err := expandApigeeEnvironmentAddonsConfigAnalyticsEnabled(d.Get("analytics_enabled"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("analytics_enabled"); !tpgresource.IsEmptyValue(reflect.ValueOf(analyticsEnabledProp)) && (ok || !reflect.DeepEqual(v, analyticsEnabledProp)) {
		obj["analyticsEnabled"] = analyticsEnabledProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ApigeeBasePath}}{{env_id}}/addonsConfig:setAddonEnablement")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new EnvironmentAddonsConfig: %#v", obj)
	billingProject := ""

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
		return fmt.Errorf("Error creating EnvironmentAddonsConfig: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "{{env_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = ApigeeOperationWaitTime(
		config, res, "Creating EnvironmentAddonsConfig", userAgent,
		d.Timeout(schema.TimeoutCreate))

	if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error waiting to create EnvironmentAddonsConfig: %s", err)
	}

	log.Printf("[DEBUG] Finished creating EnvironmentAddonsConfig %q: %#v", d.Id(), res)

	return resourceApigeeEnvironmentAddonsConfigRead(d, meta)
}

func resourceApigeeEnvironmentAddonsConfigRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ApigeeBasePath}}{{env_id}}/addonsConfig")
	if err != nil {
		return err
	}

	billingProject := ""

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
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("ApigeeEnvironmentAddonsConfig %q", d.Id()))
	}

	res, err = resourceApigeeEnvironmentAddonsConfigDecoder(d, meta, res)
	if err != nil {
		return err
	}

	if res == nil {
		// Decoding the object has resulted in it being gone. It may be marked deleted
		log.Printf("[DEBUG] Removing ApigeeEnvironmentAddonsConfig because it no longer exists.")
		d.SetId("")
		return nil
	}

	if err := d.Set("analytics_enabled", flattenApigeeEnvironmentAddonsConfigAnalyticsEnabled(res["analyticsEnabled"], d, config)); err != nil {
		return fmt.Errorf("Error reading EnvironmentAddonsConfig: %s", err)
	}

	return nil
}

func resourceApigeeEnvironmentAddonsConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	obj := make(map[string]interface{})
	analyticsEnabledProp, err := expandApigeeEnvironmentAddonsConfigAnalyticsEnabled(d.Get("analytics_enabled"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("analytics_enabled"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, analyticsEnabledProp)) {
		obj["analyticsEnabled"] = analyticsEnabledProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ApigeeBasePath}}{{env_id}}/addonsConfig:setAddonEnablement")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating EnvironmentAddonsConfig %q: %#v", d.Id(), obj)
	headers := make(http.Header)

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
		Timeout:   d.Timeout(schema.TimeoutUpdate),
		Headers:   headers,
	})

	if err != nil {
		return fmt.Errorf("Error updating EnvironmentAddonsConfig %q: %s", d.Id(), err)
	} else {
		log.Printf("[DEBUG] Finished updating EnvironmentAddonsConfig %q: %#v", d.Id(), res)
	}

	err = ApigeeOperationWaitTime(
		config, res, "Updating EnvironmentAddonsConfig", userAgent,
		d.Timeout(schema.TimeoutUpdate))

	if err != nil {
		return err
	}

	return resourceApigeeEnvironmentAddonsConfigRead(d, meta)
}

func resourceApigeeEnvironmentAddonsConfigDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARNING] Apigee EnvironmentAddonsConfig resources"+
		" cannot be deleted from Google Cloud. The resource %s will be removed from Terraform"+
		" state, but will still be present on Google Cloud.", d.Id())
	d.SetId("")

	return nil
}

func resourceApigeeEnvironmentAddonsConfigImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)

	// current import_formats cannot import fields with forward slashes in their value
	if err := tpgresource.ParseImportId([]string{"(?P<env_id>.+)"}, d, config); err != nil {
		return nil, err
	}

	id := d.Get("env_id").(string)
	nameParts := strings.Split(id, "/")
	if len(nameParts) != 4 {
		return nil, fmt.Errorf("env is expected to have shape organizations/{{org_id}}/environments/{{env}}, got %s instead", id)
	}
	d.SetId(id)
	return []*schema.ResourceData{d}, nil
}

func flattenApigeeEnvironmentAddonsConfigAnalyticsEnabled(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandApigeeEnvironmentAddonsConfigAnalyticsEnabled(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func resourceApigeeEnvironmentAddonsConfigDecoder(d *schema.ResourceData, meta interface{}, res map[string]interface{}) (map[string]interface{}, error) {
	res["analyticsEnabled"] = res["analyticsConfig"].(map[string]interface{})["enabled"]
	return res, nil
}
