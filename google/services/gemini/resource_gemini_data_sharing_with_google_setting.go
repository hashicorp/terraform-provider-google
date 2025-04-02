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
//     Configuration: https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products/gemini/DataSharingWithGoogleSetting.yaml
//     Template:      https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/resource.go.tmpl
//
//     DO NOT EDIT this file directly. Any changes made to this file will be
//     overwritten during the next generation cycle.
//
// ----------------------------------------------------------------------------

package gemini

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceGeminiDataSharingWithGoogleSetting() *schema.Resource {
	return &schema.Resource{
		Create: resourceGeminiDataSharingWithGoogleSettingCreate,
		Read:   resourceGeminiDataSharingWithGoogleSettingRead,
		Update: resourceGeminiDataSharingWithGoogleSettingUpdate,
		Delete: resourceGeminiDataSharingWithGoogleSettingDelete,

		Importer: &schema.ResourceImporter{
			State: resourceGeminiDataSharingWithGoogleSettingImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			tpgresource.SetLabelsDiff,
			tpgresource.DefaultProviderProject,
		),

		Schema: map[string]*schema.Schema{
			"data_sharing_with_google_setting_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Id of the Data Sharing With Google Setting.`,
			},
			"enable_preview_data_sharing": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether preview data sharing should be enabled.`,
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Description: `Labels as key value pairs.

**Note**: This field is non-authoritative, and will only manage the labels present in your configuration.
Please refer to the field 'effective_labels' for all of the labels present on the resource.`,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Resource ID segment making up resource 'name'. It identifies the resource within its parent collection as described in https://google.aip.dev/122.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Create time stamp.`,
			},
			"effective_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: `All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `Identifier. Name of the resource.
Format:projects/{project}/locations/{location}/dataSharingWithGoogleSettings/{dataSharingWithGoogleSetting}`,
			},
			"terraform_labels": {
				Type:     schema.TypeMap,
				Computed: true,
				Description: `The combination of labels configured directly on the resource
 and default labels configured on the provider.`,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Update time stamp.`,
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

func resourceGeminiDataSharingWithGoogleSettingCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	enablePreviewDataSharingProp, err := expandGeminiDataSharingWithGoogleSettingEnablePreviewDataSharing(d.Get("enable_preview_data_sharing"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("enable_preview_data_sharing"); !tpgresource.IsEmptyValue(reflect.ValueOf(enablePreviewDataSharingProp)) && (ok || !reflect.DeepEqual(v, enablePreviewDataSharingProp)) {
		obj["enablePreviewDataSharing"] = enablePreviewDataSharingProp
	}
	labelsProp, err := expandGeminiDataSharingWithGoogleSettingEffectiveLabels(d.Get("effective_labels"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("effective_labels"); !tpgresource.IsEmptyValue(reflect.ValueOf(labelsProp)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}

	lockName, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/dataSharingWithGoogleSettings/{{data_sharing_with_google_setting_id}}")
	if err != nil {
		return err
	}
	transport_tpg.MutexStore.Lock(lockName)
	defer transport_tpg.MutexStore.Unlock(lockName)

	url, err := tpgresource.ReplaceVars(d, config, "{{GeminiBasePath}}projects/{{project}}/locations/{{location}}/dataSharingWithGoogleSettings?dataSharingWithGoogleSettingId={{data_sharing_with_google_setting_id}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new DataSharingWithGoogleSetting: %#v", obj)
	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for DataSharingWithGoogleSetting: %s", err)
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
		return fmt.Errorf("Error creating DataSharingWithGoogleSetting: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/dataSharingWithGoogleSettings/{{data_sharing_with_google_setting_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating DataSharingWithGoogleSetting %q: %#v", d.Id(), res)

	return resourceGeminiDataSharingWithGoogleSettingRead(d, meta)
}

func resourceGeminiDataSharingWithGoogleSettingRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{GeminiBasePath}}projects/{{project}}/locations/{{location}}/dataSharingWithGoogleSettings/{{data_sharing_with_google_setting_id}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for DataSharingWithGoogleSetting: %s", err)
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
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("GeminiDataSharingWithGoogleSetting %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading DataSharingWithGoogleSetting: %s", err)
	}

	if err := d.Set("name", flattenGeminiDataSharingWithGoogleSettingName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading DataSharingWithGoogleSetting: %s", err)
	}
	if err := d.Set("create_time", flattenGeminiDataSharingWithGoogleSettingCreateTime(res["createTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading DataSharingWithGoogleSetting: %s", err)
	}
	if err := d.Set("update_time", flattenGeminiDataSharingWithGoogleSettingUpdateTime(res["updateTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading DataSharingWithGoogleSetting: %s", err)
	}
	if err := d.Set("labels", flattenGeminiDataSharingWithGoogleSettingLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading DataSharingWithGoogleSetting: %s", err)
	}
	if err := d.Set("enable_preview_data_sharing", flattenGeminiDataSharingWithGoogleSettingEnablePreviewDataSharing(res["enablePreviewDataSharing"], d, config)); err != nil {
		return fmt.Errorf("Error reading DataSharingWithGoogleSetting: %s", err)
	}
	if err := d.Set("terraform_labels", flattenGeminiDataSharingWithGoogleSettingTerraformLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading DataSharingWithGoogleSetting: %s", err)
	}
	if err := d.Set("effective_labels", flattenGeminiDataSharingWithGoogleSettingEffectiveLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading DataSharingWithGoogleSetting: %s", err)
	}

	return nil
}

func resourceGeminiDataSharingWithGoogleSettingUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for DataSharingWithGoogleSetting: %s", err)
	}
	billingProject = project

	obj := make(map[string]interface{})
	enablePreviewDataSharingProp, err := expandGeminiDataSharingWithGoogleSettingEnablePreviewDataSharing(d.Get("enable_preview_data_sharing"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("enable_preview_data_sharing"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, enablePreviewDataSharingProp)) {
		obj["enablePreviewDataSharing"] = enablePreviewDataSharingProp
	}
	labelsProp, err := expandGeminiDataSharingWithGoogleSettingEffectiveLabels(d.Get("effective_labels"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("effective_labels"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}

	lockName, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/dataSharingWithGoogleSettings/{{data_sharing_with_google_setting_id}}")
	if err != nil {
		return err
	}
	transport_tpg.MutexStore.Lock(lockName)
	defer transport_tpg.MutexStore.Unlock(lockName)

	url, err := tpgresource.ReplaceVars(d, config, "{{GeminiBasePath}}projects/{{project}}/locations/{{location}}/dataSharingWithGoogleSettings/{{data_sharing_with_google_setting_id}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating DataSharingWithGoogleSetting %q: %#v", d.Id(), obj)
	headers := make(http.Header)
	updateMask := []string{}

	if d.HasChange("enable_preview_data_sharing") {
		updateMask = append(updateMask, "enablePreviewDataSharing")
	}

	if d.HasChange("effective_labels") {
		updateMask = append(updateMask, "labels")
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
			return fmt.Errorf("Error updating DataSharingWithGoogleSetting %q: %s", d.Id(), err)
		} else {
			log.Printf("[DEBUG] Finished updating DataSharingWithGoogleSetting %q: %#v", d.Id(), res)
		}

	}

	return resourceGeminiDataSharingWithGoogleSettingRead(d, meta)
}

func resourceGeminiDataSharingWithGoogleSettingDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for DataSharingWithGoogleSetting: %s", err)
	}
	billingProject = project

	lockName, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/dataSharingWithGoogleSettings/{{data_sharing_with_google_setting_id}}")
	if err != nil {
		return err
	}
	transport_tpg.MutexStore.Lock(lockName)
	defer transport_tpg.MutexStore.Unlock(lockName)

	url, err := tpgresource.ReplaceVars(d, config, "{{GeminiBasePath}}projects/{{project}}/locations/{{location}}/dataSharingWithGoogleSettings/{{data_sharing_with_google_setting_id}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)

	log.Printf("[DEBUG] Deleting DataSharingWithGoogleSetting %q", d.Id())
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
		return transport_tpg.HandleNotFoundError(err, d, "DataSharingWithGoogleSetting")
	}

	log.Printf("[DEBUG] Finished deleting DataSharingWithGoogleSetting %q: %#v", d.Id(), res)
	return nil
}

func resourceGeminiDataSharingWithGoogleSettingImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"^projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/dataSharingWithGoogleSettings/(?P<data_sharing_with_google_setting_id>[^/]+)$",
		"^(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<data_sharing_with_google_setting_id>[^/]+)$",
		"^(?P<location>[^/]+)/(?P<data_sharing_with_google_setting_id>[^/]+)$",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/dataSharingWithGoogleSettings/{{data_sharing_with_google_setting_id}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenGeminiDataSharingWithGoogleSettingName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenGeminiDataSharingWithGoogleSettingCreateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenGeminiDataSharingWithGoogleSettingUpdateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenGeminiDataSharingWithGoogleSettingLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}

	transformed := make(map[string]interface{})
	if l, ok := d.GetOkExists("labels"); ok {
		for k := range l.(map[string]interface{}) {
			transformed[k] = v.(map[string]interface{})[k]
		}
	}

	return transformed
}

func flattenGeminiDataSharingWithGoogleSettingEnablePreviewDataSharing(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenGeminiDataSharingWithGoogleSettingTerraformLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}

	transformed := make(map[string]interface{})
	if l, ok := d.GetOkExists("terraform_labels"); ok {
		for k := range l.(map[string]interface{}) {
			transformed[k] = v.(map[string]interface{})[k]
		}
	}

	return transformed
}

func flattenGeminiDataSharingWithGoogleSettingEffectiveLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandGeminiDataSharingWithGoogleSettingEnablePreviewDataSharing(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandGeminiDataSharingWithGoogleSettingEffectiveLabels(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}
