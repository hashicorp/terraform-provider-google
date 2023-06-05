// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: DCL     ***
//
// ----------------------------------------------------------------------------
//
//     This file is managed by Magic Modules (https://github.com/GoogleCloudPlatform/magic-modules)
//     and is based on the DCL (https://github.com/GoogleCloudPlatform/declarative-resource-client-library).
//     Changes will need to be made to the DCL or Magic Modules instead of here.
//
//     We are not currently able to accept contributions to this file. If changes
//     are required, please file an issue at https://github.com/hashicorp/terraform-provider-google/issues/new/choose
//
// ----------------------------------------------------------------------------

package apikeys

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	apikeys "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/apikeys"

	"github.com/hashicorp/terraform-provider-google/google/tpgdclresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceApikeysKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceApikeysKeyCreate,
		Read:   resourceApikeysKeyRead,
		Update: resourceApikeysKeyUpdate,
		Delete: resourceApikeysKeyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceApikeysKeyImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The resource name of the key. The name must be unique within the project, must conform with RFC-1034, is restricted to lower-cased letters, and has a maximum length of 63 characters. In another word, the name must match the regular expression: `[a-z]([a-z0-9-]{0,61}[a-z0-9])?`.",
			},

			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Human-readable display name of this API key. Modifiable by user.",
			},

			"project": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "The project for the resource",
			},

			"restrictions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Key restrictions.",
				MaxItems:    1,
				Elem:        ApikeysKeyRestrictionsSchema(),
			},

			"key_string": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Output only. An encrypted and signed value held by this key. This field can be accessed only through the `GetKeyString` method.",
			},

			"uid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. Unique id in UUID4 format.",
			},
		},
	}
}

func ApikeysKeyRestrictionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"android_key_restrictions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Android apps that are allowed to use the key.",
				MaxItems:    1,
				Elem:        ApikeysKeyRestrictionsAndroidKeyRestrictionsSchema(),
			},

			"api_targets": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A restriction for a specific service and optionally one or more specific methods. Requests are allowed if they match any of these restrictions. If no restrictions are specified, all targets are allowed.",
				Elem:        ApikeysKeyRestrictionsApiTargetsSchema(),
			},

			"browser_key_restrictions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The HTTP referrers (websites) that are allowed to use the key.",
				MaxItems:    1,
				Elem:        ApikeysKeyRestrictionsBrowserKeyRestrictionsSchema(),
			},

			"ios_key_restrictions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The iOS apps that are allowed to use the key.",
				MaxItems:    1,
				Elem:        ApikeysKeyRestrictionsIosKeyRestrictionsSchema(),
			},

			"server_key_restrictions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The IP addresses of callers that are allowed to use the key.",
				MaxItems:    1,
				Elem:        ApikeysKeyRestrictionsServerKeyRestrictionsSchema(),
			},
		},
	}
}

func ApikeysKeyRestrictionsAndroidKeyRestrictionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"allowed_applications": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "A list of Android applications that are allowed to make API calls with this key.",
				Elem:        ApikeysKeyRestrictionsAndroidKeyRestrictionsAllowedApplicationsSchema(),
			},
		},
	}
}

func ApikeysKeyRestrictionsAndroidKeyRestrictionsAllowedApplicationsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"package_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The package name of the application.",
			},

			"sha1_fingerprint": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The SHA1 fingerprint of the application. For example, both sha1 formats are acceptable : DA:39:A3:EE:5E:6B:4B:0D:32:55:BF:EF:95:60:18:90:AF:D8:07:09 or DA39A3EE5E6B4B0D3255BFEF95601890AFD80709. Output format is the latter.",
			},
		},
	}
}

func ApikeysKeyRestrictionsApiTargetsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"service": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The service for this restriction. It should be the canonical service name, for example: `translate.googleapis.com`. You can use `gcloud services list` to get a list of services that are enabled in the project.",
			},

			"methods": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. List of one or more methods that can be called. If empty, all methods for the service are allowed. A wildcard (*) can be used as the last symbol. Valid examples: `google.cloud.translate.v2.TranslateService.GetSupportedLanguage` `TranslateText` `Get*` `translate.googleapis.com.Get*`",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func ApikeysKeyRestrictionsBrowserKeyRestrictionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"allowed_referrers": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "A list of regular expressions for the referrer URLs that are allowed to make API calls with this key.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func ApikeysKeyRestrictionsIosKeyRestrictionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"allowed_bundle_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "A list of bundle IDs that are allowed when making API calls with this key.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func ApikeysKeyRestrictionsServerKeyRestrictionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"allowed_ips": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "A list of the caller IP addresses that are allowed to make API calls with this key.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceApikeysKeyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &apikeys.Key{
		Name:         dcl.String(d.Get("name").(string)),
		DisplayName:  dcl.String(d.Get("display_name").(string)),
		Project:      dcl.String(project),
		Restrictions: expandApikeysKeyRestrictions(d.Get("restrictions")),
	}

	id, err := obj.ID()
	if err != nil {
		return fmt.Errorf("error constructing id: %s", err)
	}
	d.SetId(id)
	directive := tpgdclresource.CreateDirective
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLApikeysClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutCreate))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyKey(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating Key: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Key %q: %#v", d.Id(), res)

	return resourceApikeysKeyRead(d, meta)
}

func resourceApikeysKeyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &apikeys.Key{
		Name:         dcl.String(d.Get("name").(string)),
		DisplayName:  dcl.String(d.Get("display_name").(string)),
		Project:      dcl.String(project),
		Restrictions: expandApikeysKeyRestrictions(d.Get("restrictions")),
	}

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLApikeysClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutRead))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.GetKey(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("ApikeysKey %q", d.Id())
		return tpgdclresource.HandleNotFoundDCLError(err, d, resourceName)
	}

	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}
	if err = d.Set("display_name", res.DisplayName); err != nil {
		return fmt.Errorf("error setting display_name in state: %s", err)
	}
	if err = d.Set("project", res.Project); err != nil {
		return fmt.Errorf("error setting project in state: %s", err)
	}
	if err = d.Set("restrictions", flattenApikeysKeyRestrictions(res.Restrictions)); err != nil {
		return fmt.Errorf("error setting restrictions in state: %s", err)
	}
	if err = d.Set("key_string", res.KeyString); err != nil {
		return fmt.Errorf("error setting key_string in state: %s", err)
	}
	if err = d.Set("uid", res.Uid); err != nil {
		return fmt.Errorf("error setting uid in state: %s", err)
	}

	return nil
}
func resourceApikeysKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &apikeys.Key{
		Name:         dcl.String(d.Get("name").(string)),
		DisplayName:  dcl.String(d.Get("display_name").(string)),
		Project:      dcl.String(project),
		Restrictions: expandApikeysKeyRestrictions(d.Get("restrictions")),
	}
	directive := tpgdclresource.UpdateDirective
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLApikeysClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutUpdate))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyKey(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error updating Key: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Key %q: %#v", d.Id(), res)

	return resourceApikeysKeyRead(d, meta)
}

func resourceApikeysKeyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &apikeys.Key{
		Name:         dcl.String(d.Get("name").(string)),
		DisplayName:  dcl.String(d.Get("display_name").(string)),
		Project:      dcl.String(project),
		Restrictions: expandApikeysKeyRestrictions(d.Get("restrictions")),
	}

	log.Printf("[DEBUG] Deleting Key %q", d.Id())
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLApikeysClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutDelete))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	if err := client.DeleteKey(context.Background(), obj); err != nil {
		return fmt.Errorf("Error deleting Key: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting Key %q", d.Id())
	return nil
}

func resourceApikeysKeyImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)

	if err := tpgresource.ParseImportId([]string{
		"projects/(?P<project>[^/]+)/locations/global/keys/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<name>[^/]+)",
		"(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVarsForId(d, config, "projects/{{project}}/locations/global/keys/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func expandApikeysKeyRestrictions(o interface{}) *apikeys.KeyRestrictions {
	if o == nil {
		return apikeys.EmptyKeyRestrictions
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return apikeys.EmptyKeyRestrictions
	}
	obj := objArr[0].(map[string]interface{})
	return &apikeys.KeyRestrictions{
		AndroidKeyRestrictions: expandApikeysKeyRestrictionsAndroidKeyRestrictions(obj["android_key_restrictions"]),
		ApiTargets:             expandApikeysKeyRestrictionsApiTargetsArray(obj["api_targets"]),
		BrowserKeyRestrictions: expandApikeysKeyRestrictionsBrowserKeyRestrictions(obj["browser_key_restrictions"]),
		IosKeyRestrictions:     expandApikeysKeyRestrictionsIosKeyRestrictions(obj["ios_key_restrictions"]),
		ServerKeyRestrictions:  expandApikeysKeyRestrictionsServerKeyRestrictions(obj["server_key_restrictions"]),
	}
}

func flattenApikeysKeyRestrictions(obj *apikeys.KeyRestrictions) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"android_key_restrictions": flattenApikeysKeyRestrictionsAndroidKeyRestrictions(obj.AndroidKeyRestrictions),
		"api_targets":              flattenApikeysKeyRestrictionsApiTargetsArray(obj.ApiTargets),
		"browser_key_restrictions": flattenApikeysKeyRestrictionsBrowserKeyRestrictions(obj.BrowserKeyRestrictions),
		"ios_key_restrictions":     flattenApikeysKeyRestrictionsIosKeyRestrictions(obj.IosKeyRestrictions),
		"server_key_restrictions":  flattenApikeysKeyRestrictionsServerKeyRestrictions(obj.ServerKeyRestrictions),
	}

	return []interface{}{transformed}

}

func expandApikeysKeyRestrictionsAndroidKeyRestrictions(o interface{}) *apikeys.KeyRestrictionsAndroidKeyRestrictions {
	if o == nil {
		return apikeys.EmptyKeyRestrictionsAndroidKeyRestrictions
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return apikeys.EmptyKeyRestrictionsAndroidKeyRestrictions
	}
	obj := objArr[0].(map[string]interface{})
	return &apikeys.KeyRestrictionsAndroidKeyRestrictions{
		AllowedApplications: expandApikeysKeyRestrictionsAndroidKeyRestrictionsAllowedApplicationsArray(obj["allowed_applications"]),
	}
}

func flattenApikeysKeyRestrictionsAndroidKeyRestrictions(obj *apikeys.KeyRestrictionsAndroidKeyRestrictions) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"allowed_applications": flattenApikeysKeyRestrictionsAndroidKeyRestrictionsAllowedApplicationsArray(obj.AllowedApplications),
	}

	return []interface{}{transformed}

}
func expandApikeysKeyRestrictionsAndroidKeyRestrictionsAllowedApplicationsArray(o interface{}) []apikeys.KeyRestrictionsAndroidKeyRestrictionsAllowedApplications {
	if o == nil {
		return make([]apikeys.KeyRestrictionsAndroidKeyRestrictionsAllowedApplications, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make([]apikeys.KeyRestrictionsAndroidKeyRestrictionsAllowedApplications, 0)
	}

	items := make([]apikeys.KeyRestrictionsAndroidKeyRestrictionsAllowedApplications, 0, len(objs))
	for _, item := range objs {
		i := expandApikeysKeyRestrictionsAndroidKeyRestrictionsAllowedApplications(item)
		items = append(items, *i)
	}

	return items
}

func expandApikeysKeyRestrictionsAndroidKeyRestrictionsAllowedApplications(o interface{}) *apikeys.KeyRestrictionsAndroidKeyRestrictionsAllowedApplications {
	if o == nil {
		return apikeys.EmptyKeyRestrictionsAndroidKeyRestrictionsAllowedApplications
	}

	obj := o.(map[string]interface{})
	return &apikeys.KeyRestrictionsAndroidKeyRestrictionsAllowedApplications{
		PackageName:     dcl.String(obj["package_name"].(string)),
		Sha1Fingerprint: dcl.String(obj["sha1_fingerprint"].(string)),
	}
}

func flattenApikeysKeyRestrictionsAndroidKeyRestrictionsAllowedApplicationsArray(objs []apikeys.KeyRestrictionsAndroidKeyRestrictionsAllowedApplications) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenApikeysKeyRestrictionsAndroidKeyRestrictionsAllowedApplications(&item)
		items = append(items, i)
	}

	return items
}

func flattenApikeysKeyRestrictionsAndroidKeyRestrictionsAllowedApplications(obj *apikeys.KeyRestrictionsAndroidKeyRestrictionsAllowedApplications) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"package_name":     obj.PackageName,
		"sha1_fingerprint": obj.Sha1Fingerprint,
	}

	return transformed

}
func expandApikeysKeyRestrictionsApiTargetsArray(o interface{}) []apikeys.KeyRestrictionsApiTargets {
	if o == nil {
		return make([]apikeys.KeyRestrictionsApiTargets, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make([]apikeys.KeyRestrictionsApiTargets, 0)
	}

	items := make([]apikeys.KeyRestrictionsApiTargets, 0, len(objs))
	for _, item := range objs {
		i := expandApikeysKeyRestrictionsApiTargets(item)
		items = append(items, *i)
	}

	return items
}

func expandApikeysKeyRestrictionsApiTargets(o interface{}) *apikeys.KeyRestrictionsApiTargets {
	if o == nil {
		return apikeys.EmptyKeyRestrictionsApiTargets
	}

	obj := o.(map[string]interface{})
	return &apikeys.KeyRestrictionsApiTargets{
		Service: dcl.String(obj["service"].(string)),
		Methods: tpgdclresource.ExpandStringArray(obj["methods"]),
	}
}

func flattenApikeysKeyRestrictionsApiTargetsArray(objs []apikeys.KeyRestrictionsApiTargets) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenApikeysKeyRestrictionsApiTargets(&item)
		items = append(items, i)
	}

	return items
}

func flattenApikeysKeyRestrictionsApiTargets(obj *apikeys.KeyRestrictionsApiTargets) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"service": obj.Service,
		"methods": obj.Methods,
	}

	return transformed

}

func expandApikeysKeyRestrictionsBrowserKeyRestrictions(o interface{}) *apikeys.KeyRestrictionsBrowserKeyRestrictions {
	if o == nil {
		return apikeys.EmptyKeyRestrictionsBrowserKeyRestrictions
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return apikeys.EmptyKeyRestrictionsBrowserKeyRestrictions
	}
	obj := objArr[0].(map[string]interface{})
	return &apikeys.KeyRestrictionsBrowserKeyRestrictions{
		AllowedReferrers: tpgdclresource.ExpandStringArray(obj["allowed_referrers"]),
	}
}

func flattenApikeysKeyRestrictionsBrowserKeyRestrictions(obj *apikeys.KeyRestrictionsBrowserKeyRestrictions) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"allowed_referrers": obj.AllowedReferrers,
	}

	return []interface{}{transformed}

}

func expandApikeysKeyRestrictionsIosKeyRestrictions(o interface{}) *apikeys.KeyRestrictionsIosKeyRestrictions {
	if o == nil {
		return apikeys.EmptyKeyRestrictionsIosKeyRestrictions
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return apikeys.EmptyKeyRestrictionsIosKeyRestrictions
	}
	obj := objArr[0].(map[string]interface{})
	return &apikeys.KeyRestrictionsIosKeyRestrictions{
		AllowedBundleIds: tpgdclresource.ExpandStringArray(obj["allowed_bundle_ids"]),
	}
}

func flattenApikeysKeyRestrictionsIosKeyRestrictions(obj *apikeys.KeyRestrictionsIosKeyRestrictions) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"allowed_bundle_ids": obj.AllowedBundleIds,
	}

	return []interface{}{transformed}

}

func expandApikeysKeyRestrictionsServerKeyRestrictions(o interface{}) *apikeys.KeyRestrictionsServerKeyRestrictions {
	if o == nil {
		return apikeys.EmptyKeyRestrictionsServerKeyRestrictions
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return apikeys.EmptyKeyRestrictionsServerKeyRestrictions
	}
	obj := objArr[0].(map[string]interface{})
	return &apikeys.KeyRestrictionsServerKeyRestrictions{
		AllowedIps: tpgdclresource.ExpandStringArray(obj["allowed_ips"]),
	}
}

func flattenApikeysKeyRestrictionsServerKeyRestrictions(obj *apikeys.KeyRestrictionsServerKeyRestrictions) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"allowed_ips": obj.AllowedIps,
	}

	return []interface{}{transformed}

}
