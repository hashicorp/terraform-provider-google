// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dialogflowcx

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceDialogflowCXEnvironment() *schema.Resource {
	return &schema.Resource{
		Create: resourceDialogflowCXEnvironmentCreate,
		Read:   resourceDialogflowCXEnvironmentRead,
		Update: resourceDialogflowCXEnvironmentUpdate,
		Delete: resourceDialogflowCXEnvironmentDelete,

		Importer: &schema.ResourceImporter{
			State: resourceDialogflowCXEnvironmentImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Update: schema.DefaultTimeout(40 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(0, 64),
				Description:  `The human-readable name of the environment (unique in an agent). Limit of 64 characters.`,
			},
			"version_configs": {
				Type:        schema.TypeList,
				Required:    true,
				Description: `A list of configurations for flow versions. You should include version configs for all flows that are reachable from [Start Flow][Agent.start_flow] in the agent. Otherwise, an error will be returned.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Format: projects/{{project}}/locations/{{location}}/agents/{{agent}}/flows/{{flow}}/versions/{{version}}.`,
						},
					},
				},
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 500),
				Description:  `The human-readable description of the environment. The maximum length is 500 characters. If exceeded, the request is rejected.`,
			},
			"parent": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: `The Agent to create an Environment for.
Format: projects/<Project ID>/locations/<Location ID>/agents/<Agent ID>.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the environment.`,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Update time of this environment. A timestamp in RFC3339 UTC "Zulu" format, with nanosecond resolution and up to nine fractional digits. Examples: "2014-10-02T15:01:23Z" and "2014-10-02T15:01:23.045123456Z".`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceDialogflowCXEnvironmentCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	displayNameProp, err := expandDialogflowCXEnvironmentDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("display_name"); !tpgresource.IsEmptyValue(reflect.ValueOf(displayNameProp)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}
	descriptionProp, err := expandDialogflowCXEnvironmentDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	versionConfigsProp, err := expandDialogflowCXEnvironmentVersionConfigs(d.Get("version_configs"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("version_configs"); !tpgresource.IsEmptyValue(reflect.ValueOf(versionConfigsProp)) && (ok || !reflect.DeepEqual(v, versionConfigsProp)) {
		obj["versionConfigs"] = versionConfigsProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{DialogflowCXBasePath}}{{parent}}/environments")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new Environment: %#v", obj)
	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	// extract location from the parent
	location := ""

	if parts := regexp.MustCompile(`locations\/([^\/]*)\/`).FindStringSubmatch(d.Get("parent").(string)); parts != nil {
		location = parts[1]
	} else {
		return fmt.Errorf(
			"Saw %s when the parent is expected to contains location %s",
			d.Get("parent"),
			"projects/{{project}}/locations/{{location}}/...",
		)
	}

	url = strings.Replace(url, "-dialogflow", fmt.Sprintf("%s-dialogflow", location), 1)
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
		return fmt.Errorf("Error creating Environment: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "{{parent}}/environments/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	// Use the resource in the operation response to populate
	// identity fields and d.Id() before read
	var opRes map[string]interface{}
	err = DialogflowCXOperationWaitTimeWithResponse(
		config, res, &opRes, "Creating Environment", userAgent, location,
		d.Timeout(schema.TimeoutCreate))
	if err != nil {
		// The resource didn't actually create
		d.SetId("")

		return fmt.Errorf("Error waiting to create Environment: %s", err)
	}

	if err := d.Set("name", flattenDialogflowCXEnvironmentName(opRes["name"], d, config)); err != nil {
		return err
	}

	// This may have caused the ID to update - update it if so.
	id, err = tpgresource.ReplaceVars(d, config, "{{parent}}/environments/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating Environment %q: %#v", d.Id(), res)

	return resourceDialogflowCXEnvironmentRead(d, meta)
}

func resourceDialogflowCXEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{DialogflowCXBasePath}}{{parent}}/environments/{{name}}")
	if err != nil {
		return err
	}

	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	// extract location from the parent
	location := ""

	if parts := regexp.MustCompile(`locations\/([^\/]*)\/`).FindStringSubmatch(d.Get("parent").(string)); parts != nil {
		location = parts[1]
	} else {
		return fmt.Errorf(
			"Saw %s when the parent is expected to contains location %s",
			d.Get("parent"),
			"projects/{{project}}/locations/{{location}}/...",
		)
	}

	url = strings.Replace(url, "-dialogflow", fmt.Sprintf("%s-dialogflow", location), 1)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("DialogflowCXEnvironment %q", d.Id()))
	}

	if err := d.Set("name", flattenDialogflowCXEnvironmentName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading Environment: %s", err)
	}
	if err := d.Set("display_name", flattenDialogflowCXEnvironmentDisplayName(res["displayName"], d, config)); err != nil {
		return fmt.Errorf("Error reading Environment: %s", err)
	}
	if err := d.Set("description", flattenDialogflowCXEnvironmentDescription(res["description"], d, config)); err != nil {
		return fmt.Errorf("Error reading Environment: %s", err)
	}
	if err := d.Set("version_configs", flattenDialogflowCXEnvironmentVersionConfigs(res["versionConfigs"], d, config)); err != nil {
		return fmt.Errorf("Error reading Environment: %s", err)
	}
	if err := d.Set("update_time", flattenDialogflowCXEnvironmentUpdateTime(res["updateTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading Environment: %s", err)
	}

	return nil
}

func resourceDialogflowCXEnvironmentUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	obj := make(map[string]interface{})
	displayNameProp, err := expandDialogflowCXEnvironmentDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("display_name"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}
	descriptionProp, err := expandDialogflowCXEnvironmentDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	versionConfigsProp, err := expandDialogflowCXEnvironmentVersionConfigs(d.Get("version_configs"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("version_configs"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, versionConfigsProp)) {
		obj["versionConfigs"] = versionConfigsProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{DialogflowCXBasePath}}{{parent}}/environments/{{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating Environment %q: %#v", d.Id(), obj)
	updateMask := []string{}

	if d.HasChange("display_name") {
		updateMask = append(updateMask, "displayName")
	}

	if d.HasChange("description") {
		updateMask = append(updateMask, "description")
	}

	if d.HasChange("version_configs") {
		updateMask = append(updateMask, "versionConfigs")
	}
	// updateMask is a URL parameter but not present in the schema, so ReplaceVars
	// won't set it
	url, err = transport_tpg.AddQueryParams(url, map[string]string{"updateMask": strings.Join(updateMask, ",")})
	if err != nil {
		return err
	}

	// extract location from the parent
	location := ""

	if parts := regexp.MustCompile(`locations\/([^\/]*)\/`).FindStringSubmatch(d.Get("parent").(string)); parts != nil {
		location = parts[1]
	} else {
		return fmt.Errorf(
			"Saw %s when the parent is expected to contains location %s",
			d.Get("parent"),
			"projects/{{project}}/locations/{{location}}/...",
		)
	}

	url = strings.Replace(url, "-dialogflow", fmt.Sprintf("%s-dialogflow", location), 1)

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
		return fmt.Errorf("Error updating Environment %q: %s", d.Id(), err)
	} else {
		log.Printf("[DEBUG] Finished updating Environment %q: %#v", d.Id(), res)
	}

	err = DialogflowCXOperationWaitTime(
		config, res, "Updating Environment", userAgent, location,
		d.Timeout(schema.TimeoutUpdate))

	if err != nil {
		return err
	}

	return resourceDialogflowCXEnvironmentRead(d, meta)
}

func resourceDialogflowCXEnvironmentDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	url, err := tpgresource.ReplaceVars(d, config, "{{DialogflowCXBasePath}}{{parent}}/environments/{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}

	// extract location from the parent
	location := ""

	if parts := regexp.MustCompile(`locations\/([^\/]*)\/`).FindStringSubmatch(d.Get("parent").(string)); parts != nil {
		location = parts[1]
	} else {
		return fmt.Errorf(
			"Saw %s when the parent is expected to contains location %s",
			d.Get("parent"),
			"projects/{{project}}/locations/{{location}}/...",
		)
	}

	url = strings.Replace(url, "-dialogflow", fmt.Sprintf("%s-dialogflow", location), 1)
	log.Printf("[DEBUG] Deleting Environment %q", d.Id())

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
		return transport_tpg.HandleNotFoundError(err, d, "Environment")
	}

	err = DialogflowCXOperationWaitTime(
		config, res, "Deleting Environment", userAgent, location,
		d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting Environment %q: %#v", d.Id(), res)
	return nil
}

func resourceDialogflowCXEnvironmentImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)

	// current import_formats can't import fields with forward slashes in their value and parent contains slashes
	if err := tpgresource.ParseImportId([]string{
		"(?P<parent>.+)/environments/(?P<name>[^/]+)",
		"(?P<parent>.+)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "{{parent}}/environments/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenDialogflowCXEnvironmentName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	return tpgresource.NameFromSelfLinkStateFunc(v)
}

func flattenDialogflowCXEnvironmentDisplayName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDialogflowCXEnvironmentDescription(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDialogflowCXEnvironmentVersionConfigs(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"version": flattenDialogflowCXEnvironmentVersionConfigsVersion(original["version"], d, config),
		})
	}
	return transformed
}
func flattenDialogflowCXEnvironmentVersionConfigsVersion(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDialogflowCXEnvironmentUpdateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandDialogflowCXEnvironmentDisplayName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDialogflowCXEnvironmentDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDialogflowCXEnvironmentVersionConfigs(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedVersion, err := expandDialogflowCXEnvironmentVersionConfigsVersion(original["version"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedVersion); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["version"] = transformedVersion
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandDialogflowCXEnvironmentVersionConfigsVersion(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
