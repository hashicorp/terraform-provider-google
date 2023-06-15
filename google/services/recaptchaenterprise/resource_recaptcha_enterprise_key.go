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

package recaptchaenterprise

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	recaptchaenterprise "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/recaptchaenterprise"

	"github.com/hashicorp/terraform-provider-google/google/tpgdclresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceRecaptchaEnterpriseKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceRecaptchaEnterpriseKeyCreate,
		Read:   resourceRecaptchaEnterpriseKeyRead,
		Update: resourceRecaptchaEnterpriseKeyUpdate,
		Delete: resourceRecaptchaEnterpriseKeyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceRecaptchaEnterpriseKeyImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Human-readable display name of this key. Modifiable by user.",
			},

			"android_settings": {
				Type:          schema.TypeList,
				Optional:      true,
				Description:   "Settings for keys that can be used by Android apps.",
				MaxItems:      1,
				Elem:          RecaptchaEnterpriseKeyAndroidSettingsSchema(),
				ConflictsWith: []string{"web_settings", "ios_settings"},
			},

			"ios_settings": {
				Type:          schema.TypeList,
				Optional:      true,
				Description:   "Settings for keys that can be used by iOS apps.",
				MaxItems:      1,
				Elem:          RecaptchaEnterpriseKeyIosSettingsSchema(),
				ConflictsWith: []string{"web_settings", "android_settings"},
			},

			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "See [Creating and managing labels](https://cloud.google.com/recaptcha-enterprise/docs/labels).",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"project": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "The project for the resource",
			},

			"testing_options": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Options for user acceptance testing.",
				MaxItems:    1,
				Elem:        RecaptchaEnterpriseKeyTestingOptionsSchema(),
			},

			"web_settings": {
				Type:          schema.TypeList,
				Optional:      true,
				Description:   "Settings for keys that can be used by websites.",
				MaxItems:      1,
				Elem:          RecaptchaEnterpriseKeyWebSettingsSchema(),
				ConflictsWith: []string{"android_settings", "ios_settings"},
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp corresponding to the creation of this Key.",
			},

			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource name for the Key in the format \"projects/{project}/keys/{key}\".",
			},
		},
	}
}

func RecaptchaEnterpriseKeyAndroidSettingsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"allow_all_package_names": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If set to true, it means allowed_package_names will not be enforced.",
			},

			"allowed_package_names": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Android package names of apps allowed to use the key. Example: 'com.companyname.appname'",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func RecaptchaEnterpriseKeyIosSettingsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"allow_all_bundle_ids": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If set to true, it means allowed_bundle_ids will not be enforced.",
			},

			"allowed_bundle_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "iOS bundle ids of apps allowed to use the key. Example: 'com.companyname.productname.appname'",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func RecaptchaEnterpriseKeyTestingOptionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"testing_challenge": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "For challenge-based keys only (CHECKBOX, INVISIBLE), all challenge requests for this site will return nocaptcha if NOCAPTCHA, or an unsolvable challenge if UNSOLVABLE_CHALLENGE. Possible values: TESTING_CHALLENGE_UNSPECIFIED, NOCAPTCHA, UNSOLVABLE_CHALLENGE",
			},

			"testing_score": {
				Type:        schema.TypeFloat,
				Optional:    true,
				ForceNew:    true,
				Description: "All assessments for this Key will return this score. Must be between 0 (likely not legitimate) and 1 (likely legitimate) inclusive.",
			},
		},
	}
}

func RecaptchaEnterpriseKeyWebSettingsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"integration_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Required. Describes how this key is integrated with the website. Possible values: SCORE, CHECKBOX, INVISIBLE",
			},

			"allow_all_domains": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If set to true, it means allowed_domains will not be enforced.",
			},

			"allow_amp_traffic": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If set to true, the key can be used on AMP (Accelerated Mobile Pages) websites. This is supported only for the SCORE integration type.",
			},

			"allowed_domains": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Domains or subdomains of websites allowed to use the key. All subdomains of an allowed domain are automatically allowed. A valid domain requires a host and must not include any path, port, query or fragment. Examples: 'example.com' or 'subdomain.example.com'",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"challenge_security_preference": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Settings for the frequency and difficulty at which this key triggers captcha challenges. This should only be specified for IntegrationTypes CHECKBOX and INVISIBLE. Possible values: CHALLENGE_SECURITY_PREFERENCE_UNSPECIFIED, USABILITY, BALANCE, SECURITY",
			},
		},
	}
}

func resourceRecaptchaEnterpriseKeyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &recaptchaenterprise.Key{
		DisplayName:     dcl.String(d.Get("display_name").(string)),
		AndroidSettings: expandRecaptchaEnterpriseKeyAndroidSettings(d.Get("android_settings")),
		IosSettings:     expandRecaptchaEnterpriseKeyIosSettings(d.Get("ios_settings")),
		Labels:          tpgresource.CheckStringMap(d.Get("labels")),
		Project:         dcl.String(project),
		TestingOptions:  expandRecaptchaEnterpriseKeyTestingOptions(d.Get("testing_options")),
		WebSettings:     expandRecaptchaEnterpriseKeyWebSettings(d.Get("web_settings")),
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
	client := transport_tpg.NewDCLRecaptchaEnterpriseClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutCreate))
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

	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}
	// ID has a server-generated value, set again after creation.

	id, err = res.ID()
	if err != nil {
		return fmt.Errorf("error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating Key %q: %#v", d.Id(), res)

	return resourceRecaptchaEnterpriseKeyRead(d, meta)
}

func resourceRecaptchaEnterpriseKeyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &recaptchaenterprise.Key{
		DisplayName:     dcl.String(d.Get("display_name").(string)),
		AndroidSettings: expandRecaptchaEnterpriseKeyAndroidSettings(d.Get("android_settings")),
		IosSettings:     expandRecaptchaEnterpriseKeyIosSettings(d.Get("ios_settings")),
		Labels:          tpgresource.CheckStringMap(d.Get("labels")),
		Project:         dcl.String(project),
		TestingOptions:  expandRecaptchaEnterpriseKeyTestingOptions(d.Get("testing_options")),
		WebSettings:     expandRecaptchaEnterpriseKeyWebSettings(d.Get("web_settings")),
		Name:            dcl.StringOrNil(d.Get("name").(string)),
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
	client := transport_tpg.NewDCLRecaptchaEnterpriseClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutRead))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.GetKey(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("RecaptchaEnterpriseKey %q", d.Id())
		return tpgdclresource.HandleNotFoundDCLError(err, d, resourceName)
	}

	if err = d.Set("display_name", res.DisplayName); err != nil {
		return fmt.Errorf("error setting display_name in state: %s", err)
	}
	if err = d.Set("android_settings", flattenRecaptchaEnterpriseKeyAndroidSettings(res.AndroidSettings)); err != nil {
		return fmt.Errorf("error setting android_settings in state: %s", err)
	}
	if err = d.Set("ios_settings", flattenRecaptchaEnterpriseKeyIosSettings(res.IosSettings)); err != nil {
		return fmt.Errorf("error setting ios_settings in state: %s", err)
	}
	if err = d.Set("labels", res.Labels); err != nil {
		return fmt.Errorf("error setting labels in state: %s", err)
	}
	if err = d.Set("project", res.Project); err != nil {
		return fmt.Errorf("error setting project in state: %s", err)
	}
	if err = d.Set("testing_options", flattenRecaptchaEnterpriseKeyTestingOptions(res.TestingOptions)); err != nil {
		return fmt.Errorf("error setting testing_options in state: %s", err)
	}
	if err = d.Set("web_settings", flattenRecaptchaEnterpriseKeyWebSettings(res.WebSettings)); err != nil {
		return fmt.Errorf("error setting web_settings in state: %s", err)
	}
	if err = d.Set("create_time", res.CreateTime); err != nil {
		return fmt.Errorf("error setting create_time in state: %s", err)
	}
	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}

	return nil
}
func resourceRecaptchaEnterpriseKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &recaptchaenterprise.Key{
		DisplayName:     dcl.String(d.Get("display_name").(string)),
		AndroidSettings: expandRecaptchaEnterpriseKeyAndroidSettings(d.Get("android_settings")),
		IosSettings:     expandRecaptchaEnterpriseKeyIosSettings(d.Get("ios_settings")),
		Labels:          tpgresource.CheckStringMap(d.Get("labels")),
		Project:         dcl.String(project),
		TestingOptions:  expandRecaptchaEnterpriseKeyTestingOptions(d.Get("testing_options")),
		WebSettings:     expandRecaptchaEnterpriseKeyWebSettings(d.Get("web_settings")),
		Name:            dcl.StringOrNil(d.Get("name").(string)),
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
	client := transport_tpg.NewDCLRecaptchaEnterpriseClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutUpdate))
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

	return resourceRecaptchaEnterpriseKeyRead(d, meta)
}

func resourceRecaptchaEnterpriseKeyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &recaptchaenterprise.Key{
		DisplayName:     dcl.String(d.Get("display_name").(string)),
		AndroidSettings: expandRecaptchaEnterpriseKeyAndroidSettings(d.Get("android_settings")),
		IosSettings:     expandRecaptchaEnterpriseKeyIosSettings(d.Get("ios_settings")),
		Labels:          tpgresource.CheckStringMap(d.Get("labels")),
		Project:         dcl.String(project),
		TestingOptions:  expandRecaptchaEnterpriseKeyTestingOptions(d.Get("testing_options")),
		WebSettings:     expandRecaptchaEnterpriseKeyWebSettings(d.Get("web_settings")),
		Name:            dcl.StringOrNil(d.Get("name").(string)),
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
	client := transport_tpg.NewDCLRecaptchaEnterpriseClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutDelete))
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

func resourceRecaptchaEnterpriseKeyImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)

	if err := tpgresource.ParseImportId([]string{
		"projects/(?P<project>[^/]+)/keys/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<name>[^/]+)",
		"(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVarsForId(d, config, "projects/{{project}}/keys/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func expandRecaptchaEnterpriseKeyAndroidSettings(o interface{}) *recaptchaenterprise.KeyAndroidSettings {
	if o == nil {
		return recaptchaenterprise.EmptyKeyAndroidSettings
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return recaptchaenterprise.EmptyKeyAndroidSettings
	}
	obj := objArr[0].(map[string]interface{})
	return &recaptchaenterprise.KeyAndroidSettings{
		AllowAllPackageNames: dcl.Bool(obj["allow_all_package_names"].(bool)),
		AllowedPackageNames:  tpgdclresource.ExpandStringArray(obj["allowed_package_names"]),
	}
}

func flattenRecaptchaEnterpriseKeyAndroidSettings(obj *recaptchaenterprise.KeyAndroidSettings) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"allow_all_package_names": obj.AllowAllPackageNames,
		"allowed_package_names":   obj.AllowedPackageNames,
	}

	return []interface{}{transformed}

}

func expandRecaptchaEnterpriseKeyIosSettings(o interface{}) *recaptchaenterprise.KeyIosSettings {
	if o == nil {
		return recaptchaenterprise.EmptyKeyIosSettings
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return recaptchaenterprise.EmptyKeyIosSettings
	}
	obj := objArr[0].(map[string]interface{})
	return &recaptchaenterprise.KeyIosSettings{
		AllowAllBundleIds: dcl.Bool(obj["allow_all_bundle_ids"].(bool)),
		AllowedBundleIds:  tpgdclresource.ExpandStringArray(obj["allowed_bundle_ids"]),
	}
}

func flattenRecaptchaEnterpriseKeyIosSettings(obj *recaptchaenterprise.KeyIosSettings) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"allow_all_bundle_ids": obj.AllowAllBundleIds,
		"allowed_bundle_ids":   obj.AllowedBundleIds,
	}

	return []interface{}{transformed}

}

func expandRecaptchaEnterpriseKeyTestingOptions(o interface{}) *recaptchaenterprise.KeyTestingOptions {
	if o == nil {
		return recaptchaenterprise.EmptyKeyTestingOptions
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return recaptchaenterprise.EmptyKeyTestingOptions
	}
	obj := objArr[0].(map[string]interface{})
	return &recaptchaenterprise.KeyTestingOptions{
		TestingChallenge: recaptchaenterprise.KeyTestingOptionsTestingChallengeEnumRef(obj["testing_challenge"].(string)),
		TestingScore:     dcl.Float64(obj["testing_score"].(float64)),
	}
}

func flattenRecaptchaEnterpriseKeyTestingOptions(obj *recaptchaenterprise.KeyTestingOptions) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"testing_challenge": obj.TestingChallenge,
		"testing_score":     obj.TestingScore,
	}

	return []interface{}{transformed}

}

func expandRecaptchaEnterpriseKeyWebSettings(o interface{}) *recaptchaenterprise.KeyWebSettings {
	if o == nil {
		return recaptchaenterprise.EmptyKeyWebSettings
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return recaptchaenterprise.EmptyKeyWebSettings
	}
	obj := objArr[0].(map[string]interface{})
	return &recaptchaenterprise.KeyWebSettings{
		IntegrationType:             recaptchaenterprise.KeyWebSettingsIntegrationTypeEnumRef(obj["integration_type"].(string)),
		AllowAllDomains:             dcl.Bool(obj["allow_all_domains"].(bool)),
		AllowAmpTraffic:             dcl.Bool(obj["allow_amp_traffic"].(bool)),
		AllowedDomains:              tpgdclresource.ExpandStringArray(obj["allowed_domains"]),
		ChallengeSecurityPreference: recaptchaenterprise.KeyWebSettingsChallengeSecurityPreferenceEnumRef(obj["challenge_security_preference"].(string)),
	}
}

func flattenRecaptchaEnterpriseKeyWebSettings(obj *recaptchaenterprise.KeyWebSettings) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"integration_type":              obj.IntegrationType,
		"allow_all_domains":             obj.AllowAllDomains,
		"allow_amp_traffic":             obj.AllowAmpTraffic,
		"allowed_domains":               obj.AllowedDomains,
		"challenge_security_preference": obj.ChallengeSecurityPreference,
	}

	return []interface{}{transformed}

}
