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

package assuredworkloads

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	assuredworkloads "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/assuredworkloads"

	"github.com/hashicorp/terraform-provider-google/google/tpgdclresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceAssuredWorkloadsWorkload() *schema.Resource {
	return &schema.Resource{
		Create: resourceAssuredWorkloadsWorkloadCreate,
		Read:   resourceAssuredWorkloadsWorkloadRead,
		Update: resourceAssuredWorkloadsWorkloadUpdate,
		Delete: resourceAssuredWorkloadsWorkloadDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAssuredWorkloadsWorkloadImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"billing_account": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "Required. Input only. The billing account used for the resources which are direct children of workload. This billing account is initially associated with the resources created as part of Workload creation. After the initial creation of these resources, the customer can change the assigned billing account. The resource name has the form `billingAccounts/{billing_account_id}`. For example, 'billingAccounts/012345-567890-ABCDEF`.",
			},

			"compliance_regime": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Required. Immutable. Compliance Regime associated with this workload. Possible values: COMPLIANCE_REGIME_UNSPECIFIED, IL4, CJIS, FEDRAMP_HIGH, FEDRAMP_MODERATE, US_REGIONAL_ACCESS, HIPAA, EU_REGIONS_AND_SUPPORT, CA_REGIONS_AND_SUPPORT, ITAR, AU_REGIONS_AND_US_SUPPORT, ASSURED_WORKLOADS_FOR_PARTNERS",
			},

			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. The user-assigned display name of the Workload. When present it must be between 4 to 30 characters. Allowed characters are: lowercase and uppercase letters, numbers, hyphen, and spaces. Example: My Workload",
			},

			"location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The location for the resource",
			},

			"organization": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "The organization for the resource",
			},

			"kms_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Input only. Settings used to create a CMEK crypto key. When set a project with a KMS CMEK key is provisioned. This field is mandatory for a subset of Compliance Regimes.",
				MaxItems:    1,
				Elem:        AssuredWorkloadsWorkloadKmsSettingsSchema(),
			},

			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Optional. Labels applied to the workload.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"provisioned_resources_parent": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Input only. The parent resource for the resources managed by this Assured Workload. May be either an organization or a folder. Must be the same or a child of the Workload parent. If not specified all resources are created under the Workload parent. Formats: folders/{folder_id}, organizations/{organization_id}",
			},

			"resource_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Input only. Resource properties that are used to customize workload resources. These properties (such as custom project id) will be used to create workload resources if possible. This field is optional.",
				Elem:        AssuredWorkloadsWorkloadResourceSettingsSchema(),
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. Immutable. The Workload creation timestamp.",
			},

			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The resource name of the workload.",
			},

			"resources": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Output only. The resources associated with this workload. These resources will be created when creating the workload. If any of the projects already exist, the workload creation will fail. Always read only.",
				Elem:        AssuredWorkloadsWorkloadResourcesSchema(),
			},
		},
	}
}

func AssuredWorkloadsWorkloadKmsSettingsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"next_rotation_time": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Required. Input only. Immutable. The time at which the Key Management Service will automatically create a new version of the crypto key and mark it as the primary.",
			},

			"rotation_period": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Required. Input only. Immutable. will be advanced by this period when the Key Management Service automatically rotates a key. Must be at least 24 hours and at most 876,000 hours.",
			},
		},
	}
}

func AssuredWorkloadsWorkloadResourceSettingsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Resource identifier. For a project this represents project_number. If the project is already taken, the workload creation will fail.",
			},

			"resource_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Indicates the type of resource. This field should be specified to correspond the id to the right project type (CONSUMER_PROJECT or ENCRYPTION_KEYS_PROJECT) Possible values: RESOURCE_TYPE_UNSPECIFIED, CONSUMER_PROJECT, ENCRYPTION_KEYS_PROJECT, KEYRING, CONSUMER_FOLDER",
			},
		},
	}
}

func AssuredWorkloadsWorkloadResourcesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Resource identifier. For a project this represents project_number.",
			},

			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the type of resource. Possible values: RESOURCE_TYPE_UNSPECIFIED, CONSUMER_PROJECT, ENCRYPTION_KEYS_PROJECT, KEYRING, CONSUMER_FOLDER",
			},
		},
	}
}

func resourceAssuredWorkloadsWorkloadCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	obj := &assuredworkloads.Workload{
		BillingAccount:             dcl.String(d.Get("billing_account").(string)),
		ComplianceRegime:           assuredworkloads.WorkloadComplianceRegimeEnumRef(d.Get("compliance_regime").(string)),
		DisplayName:                dcl.String(d.Get("display_name").(string)),
		Location:                   dcl.String(d.Get("location").(string)),
		Organization:               dcl.String(d.Get("organization").(string)),
		KmsSettings:                expandAssuredWorkloadsWorkloadKmsSettings(d.Get("kms_settings")),
		Labels:                     tpgresource.CheckStringMap(d.Get("labels")),
		ProvisionedResourcesParent: dcl.String(d.Get("provisioned_resources_parent").(string)),
		ResourceSettings:           expandAssuredWorkloadsWorkloadResourceSettingsArray(d.Get("resource_settings")),
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
	billingProject := ""
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLAssuredWorkloadsClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutCreate))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyWorkload(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating Workload: %s", err)
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

	log.Printf("[DEBUG] Finished creating Workload %q: %#v", d.Id(), res)

	return resourceAssuredWorkloadsWorkloadRead(d, meta)
}

func resourceAssuredWorkloadsWorkloadRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	obj := &assuredworkloads.Workload{
		BillingAccount:             dcl.String(d.Get("billing_account").(string)),
		ComplianceRegime:           assuredworkloads.WorkloadComplianceRegimeEnumRef(d.Get("compliance_regime").(string)),
		DisplayName:                dcl.String(d.Get("display_name").(string)),
		Location:                   dcl.String(d.Get("location").(string)),
		Organization:               dcl.String(d.Get("organization").(string)),
		KmsSettings:                expandAssuredWorkloadsWorkloadKmsSettings(d.Get("kms_settings")),
		Labels:                     tpgresource.CheckStringMap(d.Get("labels")),
		ProvisionedResourcesParent: dcl.String(d.Get("provisioned_resources_parent").(string)),
		ResourceSettings:           expandAssuredWorkloadsWorkloadResourceSettingsArray(d.Get("resource_settings")),
		Name:                       dcl.StringOrNil(d.Get("name").(string)),
	}

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := ""
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLAssuredWorkloadsClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutRead))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.GetWorkload(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("AssuredWorkloadsWorkload %q", d.Id())
		return tpgdclresource.HandleNotFoundDCLError(err, d, resourceName)
	}

	if err = d.Set("billing_account", res.BillingAccount); err != nil {
		return fmt.Errorf("error setting billing_account in state: %s", err)
	}
	if err = d.Set("compliance_regime", res.ComplianceRegime); err != nil {
		return fmt.Errorf("error setting compliance_regime in state: %s", err)
	}
	if err = d.Set("display_name", res.DisplayName); err != nil {
		return fmt.Errorf("error setting display_name in state: %s", err)
	}
	if err = d.Set("location", res.Location); err != nil {
		return fmt.Errorf("error setting location in state: %s", err)
	}
	if err = d.Set("organization", res.Organization); err != nil {
		return fmt.Errorf("error setting organization in state: %s", err)
	}
	if err = d.Set("kms_settings", flattenAssuredWorkloadsWorkloadKmsSettings(res.KmsSettings)); err != nil {
		return fmt.Errorf("error setting kms_settings in state: %s", err)
	}
	if err = d.Set("labels", res.Labels); err != nil {
		return fmt.Errorf("error setting labels in state: %s", err)
	}
	if err = d.Set("provisioned_resources_parent", res.ProvisionedResourcesParent); err != nil {
		return fmt.Errorf("error setting provisioned_resources_parent in state: %s", err)
	}
	if err = d.Set("resource_settings", flattenAssuredWorkloadsWorkloadResourceSettingsArray(res.ResourceSettings)); err != nil {
		return fmt.Errorf("error setting resource_settings in state: %s", err)
	}
	if err = d.Set("create_time", res.CreateTime); err != nil {
		return fmt.Errorf("error setting create_time in state: %s", err)
	}
	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}
	if err = d.Set("resources", flattenAssuredWorkloadsWorkloadResourcesArray(res.Resources)); err != nil {
		return fmt.Errorf("error setting resources in state: %s", err)
	}

	return nil
}
func resourceAssuredWorkloadsWorkloadUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	obj := &assuredworkloads.Workload{
		BillingAccount:             dcl.String(d.Get("billing_account").(string)),
		ComplianceRegime:           assuredworkloads.WorkloadComplianceRegimeEnumRef(d.Get("compliance_regime").(string)),
		DisplayName:                dcl.String(d.Get("display_name").(string)),
		Location:                   dcl.String(d.Get("location").(string)),
		Organization:               dcl.String(d.Get("organization").(string)),
		KmsSettings:                expandAssuredWorkloadsWorkloadKmsSettings(d.Get("kms_settings")),
		Labels:                     tpgresource.CheckStringMap(d.Get("labels")),
		ProvisionedResourcesParent: dcl.String(d.Get("provisioned_resources_parent").(string)),
		ResourceSettings:           expandAssuredWorkloadsWorkloadResourceSettingsArray(d.Get("resource_settings")),
		Name:                       dcl.StringOrNil(d.Get("name").(string)),
	}
	// Construct state hint from old values
	old := &assuredworkloads.Workload{
		BillingAccount:             dcl.String(tpgdclresource.OldValue(d.GetChange("billing_account")).(string)),
		ComplianceRegime:           assuredworkloads.WorkloadComplianceRegimeEnumRef(tpgdclresource.OldValue(d.GetChange("compliance_regime")).(string)),
		DisplayName:                dcl.String(tpgdclresource.OldValue(d.GetChange("display_name")).(string)),
		Location:                   dcl.String(tpgdclresource.OldValue(d.GetChange("location")).(string)),
		Organization:               dcl.String(tpgdclresource.OldValue(d.GetChange("organization")).(string)),
		KmsSettings:                expandAssuredWorkloadsWorkloadKmsSettings(tpgdclresource.OldValue(d.GetChange("kms_settings"))),
		Labels:                     tpgresource.CheckStringMap(tpgdclresource.OldValue(d.GetChange("labels"))),
		ProvisionedResourcesParent: dcl.String(tpgdclresource.OldValue(d.GetChange("provisioned_resources_parent")).(string)),
		ResourceSettings:           expandAssuredWorkloadsWorkloadResourceSettingsArray(tpgdclresource.OldValue(d.GetChange("resource_settings"))),
		Name:                       dcl.StringOrNil(tpgdclresource.OldValue(d.GetChange("name")).(string)),
	}
	directive := tpgdclresource.UpdateDirective
	directive = append(directive, dcl.WithStateHint(old))
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLAssuredWorkloadsClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutUpdate))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyWorkload(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error updating Workload: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Workload %q: %#v", d.Id(), res)

	return resourceAssuredWorkloadsWorkloadRead(d, meta)
}

func resourceAssuredWorkloadsWorkloadDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	obj := &assuredworkloads.Workload{
		BillingAccount:             dcl.String(d.Get("billing_account").(string)),
		ComplianceRegime:           assuredworkloads.WorkloadComplianceRegimeEnumRef(d.Get("compliance_regime").(string)),
		DisplayName:                dcl.String(d.Get("display_name").(string)),
		Location:                   dcl.String(d.Get("location").(string)),
		Organization:               dcl.String(d.Get("organization").(string)),
		KmsSettings:                expandAssuredWorkloadsWorkloadKmsSettings(d.Get("kms_settings")),
		Labels:                     tpgresource.CheckStringMap(d.Get("labels")),
		ProvisionedResourcesParent: dcl.String(d.Get("provisioned_resources_parent").(string)),
		ResourceSettings:           expandAssuredWorkloadsWorkloadResourceSettingsArray(d.Get("resource_settings")),
		Name:                       dcl.StringOrNil(d.Get("name").(string)),
	}

	log.Printf("[DEBUG] Deleting Workload %q", d.Id())
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := ""
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLAssuredWorkloadsClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutDelete))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	if err := client.DeleteWorkload(context.Background(), obj); err != nil {
		return fmt.Errorf("Error deleting Workload: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting Workload %q", d.Id())
	return nil
}

func resourceAssuredWorkloadsWorkloadImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)

	if err := tpgresource.ParseImportId([]string{
		"organizations/(?P<organization>[^/]+)/locations/(?P<location>[^/]+)/workloads/(?P<name>[^/]+)",
		"(?P<organization>[^/]+)/(?P<location>[^/]+)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVarsForId(d, config, "organizations/{{organization}}/locations/{{location}}/workloads/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func expandAssuredWorkloadsWorkloadKmsSettings(o interface{}) *assuredworkloads.WorkloadKmsSettings {
	if o == nil {
		return assuredworkloads.EmptyWorkloadKmsSettings
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return assuredworkloads.EmptyWorkloadKmsSettings
	}
	obj := objArr[0].(map[string]interface{})
	return &assuredworkloads.WorkloadKmsSettings{
		NextRotationTime: dcl.String(obj["next_rotation_time"].(string)),
		RotationPeriod:   dcl.String(obj["rotation_period"].(string)),
	}
}

func flattenAssuredWorkloadsWorkloadKmsSettings(obj *assuredworkloads.WorkloadKmsSettings) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"next_rotation_time": obj.NextRotationTime,
		"rotation_period":    obj.RotationPeriod,
	}

	return []interface{}{transformed}

}
func expandAssuredWorkloadsWorkloadResourceSettingsArray(o interface{}) []assuredworkloads.WorkloadResourceSettings {
	if o == nil {
		return make([]assuredworkloads.WorkloadResourceSettings, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make([]assuredworkloads.WorkloadResourceSettings, 0)
	}

	items := make([]assuredworkloads.WorkloadResourceSettings, 0, len(objs))
	for _, item := range objs {
		i := expandAssuredWorkloadsWorkloadResourceSettings(item)
		items = append(items, *i)
	}

	return items
}

func expandAssuredWorkloadsWorkloadResourceSettings(o interface{}) *assuredworkloads.WorkloadResourceSettings {
	if o == nil {
		return assuredworkloads.EmptyWorkloadResourceSettings
	}

	obj := o.(map[string]interface{})
	return &assuredworkloads.WorkloadResourceSettings{
		ResourceId:   dcl.String(obj["resource_id"].(string)),
		ResourceType: assuredworkloads.WorkloadResourceSettingsResourceTypeEnumRef(obj["resource_type"].(string)),
	}
}

func flattenAssuredWorkloadsWorkloadResourceSettingsArray(objs []assuredworkloads.WorkloadResourceSettings) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenAssuredWorkloadsWorkloadResourceSettings(&item)
		items = append(items, i)
	}

	return items
}

func flattenAssuredWorkloadsWorkloadResourceSettings(obj *assuredworkloads.WorkloadResourceSettings) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"resource_id":   obj.ResourceId,
		"resource_type": obj.ResourceType,
	}

	return transformed

}

func flattenAssuredWorkloadsWorkloadResourcesArray(objs []assuredworkloads.WorkloadResources) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenAssuredWorkloadsWorkloadResources(&item)
		items = append(items, i)
	}

	return items
}

func flattenAssuredWorkloadsWorkloadResources(obj *assuredworkloads.WorkloadResources) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"resource_id":   obj.ResourceId,
		"resource_type": obj.ResourceType,
	}

	return transformed

}
