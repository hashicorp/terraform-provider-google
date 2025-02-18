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
//     Configuration: https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products/securitycenter/EventThreatDetectionCustomModule.yaml
//     Template:      https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/resource.go.tmpl
//
//     DO NOT EDIT this file directly. Any changes made to this file will be
//     overwritten during the next generation cycle.
//
// ----------------------------------------------------------------------------

package securitycenter

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"
)

func ResourceSecurityCenterEventThreatDetectionCustomModule() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityCenterEventThreatDetectionCustomModuleCreate,
		Read:   resourceSecurityCenterEventThreatDetectionCustomModuleRead,
		Update: resourceSecurityCenterEventThreatDetectionCustomModuleUpdate,
		Delete: resourceSecurityCenterEventThreatDetectionCustomModuleDelete,

		Importer: &schema.ResourceImporter{
			State: resourceSecurityCenterEventThreatDetectionCustomModuleImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"config": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
				StateFunc:    func(v interface{}) string { s, _ := structure.NormalizeJsonString(v); return s },
				Description: `Config for the module. For the resident module, its config value is defined at this level.
For the inherited module, its config value is inherited from the ancestor module.`,
			},
			"enablement_state": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: verify.ValidateEnum([]string{"ENABLED", "DISABLED"}),
				Description:  `The state of enablement for the module at the given level of the hierarchy. Possible values: ["ENABLED", "DISABLED"]`,
			},
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Numerical ID of the parent organization.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Immutable. Type for the module. e.g. CONFIGURABLE_BAD_IP.`,
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The human readable name to be displayed for the module.`,
			},
			"last_editor": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The editor that last updated the custom module`,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The resource name of the Event Threat Detection custom module.
Its format is "organizations/{organization}/eventThreatDetectionSettings/customModules/{module}".`,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The time at which the custom module was last updated.

A timestamp in RFC3339 UTC "Zulu" format, with nanosecond resolution and
up to nine fractional digits. Examples: "2014-10-02T15:01:23Z" and "2014-10-02T15:01:23.045123456Z".`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceSecurityCenterEventThreatDetectionCustomModuleCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	configProp, err := expandSecurityCenterEventThreatDetectionCustomModuleConfig(d.Get("config"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("config"); !tpgresource.IsEmptyValue(reflect.ValueOf(configProp)) && (ok || !reflect.DeepEqual(v, configProp)) {
		obj["config"] = configProp
	}
	enablementStateProp, err := expandSecurityCenterEventThreatDetectionCustomModuleEnablementState(d.Get("enablement_state"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("enablement_state"); !tpgresource.IsEmptyValue(reflect.ValueOf(enablementStateProp)) && (ok || !reflect.DeepEqual(v, enablementStateProp)) {
		obj["enablementState"] = enablementStateProp
	}
	typeProp, err := expandSecurityCenterEventThreatDetectionCustomModuleType(d.Get("type"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("type"); !tpgresource.IsEmptyValue(reflect.ValueOf(typeProp)) && (ok || !reflect.DeepEqual(v, typeProp)) {
		obj["type"] = typeProp
	}
	displayNameProp, err := expandSecurityCenterEventThreatDetectionCustomModuleDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("display_name"); !tpgresource.IsEmptyValue(reflect.ValueOf(displayNameProp)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}

	lockName, err := tpgresource.ReplaceVars(d, config, "organizations/{{organization}}/eventThreatDetectionSettings/customModules")
	if err != nil {
		return err
	}
	transport_tpg.MutexStore.Lock(lockName)
	defer transport_tpg.MutexStore.Unlock(lockName)

	url, err := tpgresource.ReplaceVars(d, config, "{{SecurityCenterBasePath}}organizations/{{organization}}/eventThreatDetectionSettings/customModules")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new EventThreatDetectionCustomModule: %#v", obj)
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
		return fmt.Errorf("Error creating EventThreatDetectionCustomModule: %s", err)
	}
	if err := d.Set("name", flattenSecurityCenterEventThreatDetectionCustomModuleName(res["name"], d, config)); err != nil {
		return fmt.Errorf(`Error setting computed identity field "name": %s`, err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "organizations/{{organization}}/eventThreatDetectionSettings/customModules/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating EventThreatDetectionCustomModule %q: %#v", d.Id(), res)

	return resourceSecurityCenterEventThreatDetectionCustomModuleRead(d, meta)
}

func resourceSecurityCenterEventThreatDetectionCustomModuleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{SecurityCenterBasePath}}organizations/{{organization}}/eventThreatDetectionSettings/customModules/{{name}}")
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
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("SecurityCenterEventThreatDetectionCustomModule %q", d.Id()))
	}

	if err := d.Set("name", flattenSecurityCenterEventThreatDetectionCustomModuleName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading EventThreatDetectionCustomModule: %s", err)
	}
	if err := d.Set("config", flattenSecurityCenterEventThreatDetectionCustomModuleConfig(res["config"], d, config)); err != nil {
		return fmt.Errorf("Error reading EventThreatDetectionCustomModule: %s", err)
	}
	if err := d.Set("enablement_state", flattenSecurityCenterEventThreatDetectionCustomModuleEnablementState(res["enablementState"], d, config)); err != nil {
		return fmt.Errorf("Error reading EventThreatDetectionCustomModule: %s", err)
	}
	if err := d.Set("type", flattenSecurityCenterEventThreatDetectionCustomModuleType(res["type"], d, config)); err != nil {
		return fmt.Errorf("Error reading EventThreatDetectionCustomModule: %s", err)
	}
	if err := d.Set("display_name", flattenSecurityCenterEventThreatDetectionCustomModuleDisplayName(res["displayName"], d, config)); err != nil {
		return fmt.Errorf("Error reading EventThreatDetectionCustomModule: %s", err)
	}
	if err := d.Set("update_time", flattenSecurityCenterEventThreatDetectionCustomModuleUpdateTime(res["updateTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading EventThreatDetectionCustomModule: %s", err)
	}
	if err := d.Set("last_editor", flattenSecurityCenterEventThreatDetectionCustomModuleLastEditor(res["lastEditor"], d, config)); err != nil {
		return fmt.Errorf("Error reading EventThreatDetectionCustomModule: %s", err)
	}

	return nil
}

func resourceSecurityCenterEventThreatDetectionCustomModuleUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	obj := make(map[string]interface{})
	configProp, err := expandSecurityCenterEventThreatDetectionCustomModuleConfig(d.Get("config"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("config"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, configProp)) {
		obj["config"] = configProp
	}
	enablementStateProp, err := expandSecurityCenterEventThreatDetectionCustomModuleEnablementState(d.Get("enablement_state"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("enablement_state"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, enablementStateProp)) {
		obj["enablementState"] = enablementStateProp
	}
	displayNameProp, err := expandSecurityCenterEventThreatDetectionCustomModuleDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("display_name"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}

	lockName, err := tpgresource.ReplaceVars(d, config, "organizations/{{organization}}/eventThreatDetectionSettings/customModules")
	if err != nil {
		return err
	}
	transport_tpg.MutexStore.Lock(lockName)
	defer transport_tpg.MutexStore.Unlock(lockName)

	url, err := tpgresource.ReplaceVars(d, config, "{{SecurityCenterBasePath}}organizations/{{organization}}/eventThreatDetectionSettings/customModules/{{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating EventThreatDetectionCustomModule %q: %#v", d.Id(), obj)
	headers := make(http.Header)
	updateMask := []string{}

	if d.HasChange("config") {
		updateMask = append(updateMask, "config")
	}

	if d.HasChange("enablement_state") {
		updateMask = append(updateMask, "enablementState")
	}

	if d.HasChange("display_name") {
		updateMask = append(updateMask, "displayName")
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

	// if updateMask is empty we are not updating anything so skip the post
	if len(updateMask) > 0 {
		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "PATCH",
			Project:   billingProject,
			RawURL:    url,
			UserAgent: userAgent,
			Body:      obj,
			Timeout:   d.Timeout(schema.TimeoutUpdate),
			Headers:   headers,
		})

		if err != nil {
			return fmt.Errorf("Error updating EventThreatDetectionCustomModule %q: %s", d.Id(), err)
		} else {
			log.Printf("[DEBUG] Finished updating EventThreatDetectionCustomModule %q: %#v", d.Id(), res)
		}

	}

	return resourceSecurityCenterEventThreatDetectionCustomModuleRead(d, meta)
}

func resourceSecurityCenterEventThreatDetectionCustomModuleDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	lockName, err := tpgresource.ReplaceVars(d, config, "organizations/{{organization}}/eventThreatDetectionSettings/customModules")
	if err != nil {
		return err
	}
	transport_tpg.MutexStore.Lock(lockName)
	defer transport_tpg.MutexStore.Unlock(lockName)

	url, err := tpgresource.ReplaceVars(d, config, "{{SecurityCenterBasePath}}organizations/{{organization}}/eventThreatDetectionSettings/customModules/{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)

	log.Printf("[DEBUG] Deleting EventThreatDetectionCustomModule %q", d.Id())
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "DELETE",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutDelete),
		Headers:   headers,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, "EventThreatDetectionCustomModule")
	}

	log.Printf("[DEBUG] Finished deleting EventThreatDetectionCustomModule %q: %#v", d.Id(), res)
	return nil
}

func resourceSecurityCenterEventThreatDetectionCustomModuleImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"^organizations/(?P<organization>[^/]+)/eventThreatDetectionSettings/customModules/(?P<name>[^/]+)$",
		"^(?P<organization>[^/]+)/(?P<name>[^/]+)$",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "organizations/{{organization}}/eventThreatDetectionSettings/customModules/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenSecurityCenterEventThreatDetectionCustomModuleName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	return tpgresource.NameFromSelfLinkStateFunc(v)
}

func flattenSecurityCenterEventThreatDetectionCustomModuleConfig(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		// TODO: return error once https://github.com/GoogleCloudPlatform/magic-modules/issues/3257 is fixed.
		log.Printf("[ERROR] failed to marshal schema to JSON: %v", err)
	}
	return string(b)
}

func flattenSecurityCenterEventThreatDetectionCustomModuleEnablementState(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenSecurityCenterEventThreatDetectionCustomModuleType(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenSecurityCenterEventThreatDetectionCustomModuleDisplayName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenSecurityCenterEventThreatDetectionCustomModuleUpdateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenSecurityCenterEventThreatDetectionCustomModuleLastEditor(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandSecurityCenterEventThreatDetectionCustomModuleConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	b := []byte(v.(string))
	if len(b) == 0 {
		return nil, nil
	}
	m := make(map[string]interface{})
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func expandSecurityCenterEventThreatDetectionCustomModuleEnablementState(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandSecurityCenterEventThreatDetectionCustomModuleType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandSecurityCenterEventThreatDetectionCustomModuleDisplayName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
