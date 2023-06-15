// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package apigee

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceApigeeFlowhook() *schema.Resource {
	return &schema.Resource{
		Create: resourceApigeeFlowhookCreate,
		Read:   resourceApigeeFlowhookRead,
		Delete: resourceApigeeFlowhookDelete,

		Importer: &schema.ResourceImporter{
			State: resourceApigeeFlowhookImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Description of the flow hook.`,
			},
			"environment": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The resource ID of the environment.`,
			},
			"flow_hook_point": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Where in the API call flow the flow hook is invoked. Must be one of PreProxyFlowHook, PostProxyFlowHook, PreTargetFlowHook, or PostTargetFlowHook.`,
			},
			"org_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The Apigee Organization associated with the environment`,
			},
			"sharedflow": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Id of the Sharedflow attaching to a flowhook point.`,
			},
			"continue_on_error": {
				Type:        schema.TypeBool,
				ForceNew:    true,
				Optional:    true,
				Default:     true,
				Description: `Flag that specifies whether execution should continue if the flow hook throws an exception. Set to true to continue execution. Set to false to stop execution if the flow hook throws an exception. Defaults to true.`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceApigeeFlowhookCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	descriptionProp, err := expandApigeeFlowhookDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	sharedflowProp, err := expandApigeeFlowhookSharedflow(d.Get("sharedflow"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("sharedflow"); !tpgresource.IsEmptyValue(reflect.ValueOf(sharedflowProp)) && (ok || !reflect.DeepEqual(v, sharedflowProp)) {
		obj["sharedFlow"] = sharedflowProp
	}
	continue_on_errorProp, err := expandApigeeFlowhookContinueOnError(d.Get("continue_on_error"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("continue_on_error"); !tpgresource.IsEmptyValue(reflect.ValueOf(continue_on_errorProp)) && (ok || !reflect.DeepEqual(v, continue_on_errorProp)) {
		obj["continueOnError"] = continue_on_errorProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ApigeeBasePath}}organizations/{{org_id}}/environments/{{environment}}/flowhooks/{{flow_hook_point}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new Flowhook: %#v", obj)
	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "PUT",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutCreate),
	})
	if err != nil {
		return fmt.Errorf("Error creating Flowhook: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "organizations/{{org_id}}/environments/{{environment}}/flowhooks/{{flow_hook_point}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating Flowhook %q: %#v", d.Id(), res)

	return resourceApigeeFlowhookRead(d, meta)
}

func resourceApigeeFlowhookRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ApigeeBasePath}}organizations/{{org_id}}/environments/{{environment}}/flowhooks/{{flow_hook_point}}")
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
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("ApigeeFlowhook %q", d.Id()))
	}
	if res["sharedFlow"] == nil || res["sharedFlow"].(string) == "" {
		//if response does not contain shared_flow field, then nothing is attached to this flowhook, we treat this "binding" resource non-existent
		d.SetId("")
		return nil
	}
	if err := d.Set("description", flattenApigeeFlowhookDescription(res["description"], d, config)); err != nil {
		return fmt.Errorf("Error reading Flowhook: %s", err)
	}
	if err := d.Set("sharedflow", flattenApigeeFlowhookSharedflow(res["sharedFlow"], d, config)); err != nil {
		return fmt.Errorf("Error reading Flowhook: %s", err)
	}
	if err := d.Set("continue_on_error", flattenApigeeFlowhookContinueOnError(res["continueOnError"], d, config)); err != nil {
		return fmt.Errorf("Error reading Flowhook: %s", err)
	}

	return nil
}

func resourceApigeeFlowhookDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	url, err := tpgresource.ReplaceVars(d, config, "{{ApigeeBasePath}}organizations/{{org_id}}/environments/{{environment}}/flowhooks/{{flow_hook_point}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting Flowhook %q", d.Id())

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
		return transport_tpg.HandleNotFoundError(err, d, "Flowhook")
	}

	log.Printf("[DEBUG] Finished deleting Flowhook %q: %#v", d.Id(), res)
	return nil
}

func resourceApigeeFlowhookImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"organizations/(?P<org_id>[^/]+)/environments/(?P<environment>[^/]+)/flowhooks/(?P<flow_hook_point>[^/]+)",
		"(?P<org_id>[^/]+)/(?P<environment>[^/]+)/(?P<flow_hook_point>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "organizations/{{org_id}}/environments/{{environment}}/flowhooks/{{flow_hook_point}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenApigeeFlowhookDescription(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeFlowhookSharedflow(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApigeeFlowhookContinueOnError(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandApigeeFlowhookDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeFlowhookSharedflow(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApigeeFlowhookContinueOnError(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
