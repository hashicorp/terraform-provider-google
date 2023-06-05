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

package dataplex

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	dataplex "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/dataplex"

	"github.com/hashicorp/terraform-provider-google/google/tpgdclresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceDataplexZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataplexZoneCreate,
		Read:   resourceDataplexZoneRead,
		Update: resourceDataplexZoneUpdate,
		Delete: resourceDataplexZoneDelete,

		Importer: &schema.ResourceImporter{
			State: resourceDataplexZoneImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"discovery_spec": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Required. Specification of the discovery feature applied to data in this zone.",
				MaxItems:    1,
				Elem:        DataplexZoneDiscoverySpecSchema(),
			},

			"lake": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The lake for the resource",
			},

			"location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The location for the resource",
			},

			"name": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "The name of the zone.",
			},

			"resource_spec": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "Required. Immutable. Specification of the resources that are referenced by the assets within this zone.",
				MaxItems:    1,
				Elem:        DataplexZoneResourceSpecSchema(),
			},

			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Required. Immutable. The type of the zone. Possible values: TYPE_UNSPECIFIED, RAW, CURATED",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional. Description of the zone.",
			},

			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional. User friendly display name.",
			},

			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Optional. User defined labels for the zone.",
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

			"asset_status": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Output only. Aggregated status of the underlying assets of the zone.",
				Elem:        DataplexZoneAssetStatusSchema(),
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The time when the zone was created.",
			},

			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. Current state of the zone. Possible values: STATE_UNSPECIFIED, ACTIVE, CREATING, DELETING, ACTION_REQUIRED",
			},

			"uid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. System generated globally unique ID for the zone. This ID will be different if the zone is deleted and re-created with the same name.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The time when the zone was last updated.",
			},
		},
	}
}

func DataplexZoneDiscoverySpecSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Required. Whether discovery is enabled.",
			},

			"csv_options": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: "Optional. Configuration for CSV data.",
				MaxItems:    1,
				Elem:        DataplexZoneDiscoverySpecCsvOptionsSchema(),
			},

			"exclude_patterns": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. The list of patterns to apply for selecting data to exclude during discovery. For Cloud Storage bucket assets, these are interpreted as glob patterns used to match object names. For BigQuery dataset assets, these are interpreted as patterns to match table names.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"include_patterns": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. The list of patterns to apply for selecting data to include during discovery if only a subset of the data should considered. For Cloud Storage bucket assets, these are interpreted as glob patterns used to match object names. For BigQuery dataset assets, these are interpreted as patterns to match table names.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"json_options": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: "Optional. Configuration for Json data.",
				MaxItems:    1,
				Elem:        DataplexZoneDiscoverySpecJsonOptionsSchema(),
			},

			"schedule": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Optional. Cron schedule (https://en.wikipedia.org/wiki/Cron) for running discovery periodically. Successive discovery runs must be scheduled at least 60 minutes apart. The default value is to run discovery every 60 minutes. To explicitly set a timezone to the cron tab, apply a prefix in the cron tab: \"CRON_TZ=${IANA_TIME_ZONE}\" or TZ=${IANA_TIME_ZONE}\". The ${IANA_TIME_ZONE} may only be a valid string from IANA time zone database. For example, \"CRON_TZ=America/New_York 1 * * * *\", or \"TZ=America/New_York 1 * * * *\".",
			},
		},
	}
}

func DataplexZoneDiscoverySpecCsvOptionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"delimiter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional. The delimiter being used to separate values. This defaults to ','.",
			},

			"disable_type_inference": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Optional. Whether to disable the inference of data type for CSV data. If true, all columns will be registered as strings.",
			},

			"encoding": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional. The character encoding of the data. The default is UTF-8.",
			},

			"header_rows": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Optional. The number of rows to interpret as header rows that should be skipped when reading data rows.",
			},
		},
	}
}

func DataplexZoneDiscoverySpecJsonOptionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"disable_type_inference": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Optional. Whether to disable the inference of data type for Json data. If true, all columns will be registered as their primitive types (strings, number or boolean).",
			},

			"encoding": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional. The character encoding of the data. The default is UTF-8.",
			},
		},
	}
}

func DataplexZoneResourceSpecSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"location_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Required. Immutable. The location type of the resources that are allowed to be attached to the assets within this zone. Possible values: LOCATION_TYPE_UNSPECIFIED, SINGLE_REGION, MULTI_REGION",
			},
		},
	}
}

func DataplexZoneAssetStatusSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"active_assets": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of active assets.",
			},

			"security_policy_applying_assets": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of assets that are in process of updating the security policy on attached resources.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update time of the status.",
			},
		},
	}
}

func resourceDataplexZoneCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &dataplex.Zone{
		DiscoverySpec: expandDataplexZoneDiscoverySpec(d.Get("discovery_spec")),
		Lake:          dcl.String(d.Get("lake").(string)),
		Location:      dcl.String(d.Get("location").(string)),
		Name:          dcl.String(d.Get("name").(string)),
		ResourceSpec:  expandDataplexZoneResourceSpec(d.Get("resource_spec")),
		Type:          dataplex.ZoneTypeEnumRef(d.Get("type").(string)),
		Description:   dcl.String(d.Get("description").(string)),
		DisplayName:   dcl.String(d.Get("display_name").(string)),
		Labels:        tpgresource.CheckStringMap(d.Get("labels")),
		Project:       dcl.String(project),
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
	client := transport_tpg.NewDCLDataplexClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutCreate))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyZone(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating Zone: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Zone %q: %#v", d.Id(), res)

	return resourceDataplexZoneRead(d, meta)
}

func resourceDataplexZoneRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &dataplex.Zone{
		DiscoverySpec: expandDataplexZoneDiscoverySpec(d.Get("discovery_spec")),
		Lake:          dcl.String(d.Get("lake").(string)),
		Location:      dcl.String(d.Get("location").(string)),
		Name:          dcl.String(d.Get("name").(string)),
		ResourceSpec:  expandDataplexZoneResourceSpec(d.Get("resource_spec")),
		Type:          dataplex.ZoneTypeEnumRef(d.Get("type").(string)),
		Description:   dcl.String(d.Get("description").(string)),
		DisplayName:   dcl.String(d.Get("display_name").(string)),
		Labels:        tpgresource.CheckStringMap(d.Get("labels")),
		Project:       dcl.String(project),
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
	client := transport_tpg.NewDCLDataplexClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutRead))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.GetZone(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("DataplexZone %q", d.Id())
		return tpgdclresource.HandleNotFoundDCLError(err, d, resourceName)
	}

	if err = d.Set("discovery_spec", flattenDataplexZoneDiscoverySpec(res.DiscoverySpec)); err != nil {
		return fmt.Errorf("error setting discovery_spec in state: %s", err)
	}
	if err = d.Set("lake", res.Lake); err != nil {
		return fmt.Errorf("error setting lake in state: %s", err)
	}
	if err = d.Set("location", res.Location); err != nil {
		return fmt.Errorf("error setting location in state: %s", err)
	}
	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}
	if err = d.Set("resource_spec", flattenDataplexZoneResourceSpec(res.ResourceSpec)); err != nil {
		return fmt.Errorf("error setting resource_spec in state: %s", err)
	}
	if err = d.Set("type", res.Type); err != nil {
		return fmt.Errorf("error setting type in state: %s", err)
	}
	if err = d.Set("description", res.Description); err != nil {
		return fmt.Errorf("error setting description in state: %s", err)
	}
	if err = d.Set("display_name", res.DisplayName); err != nil {
		return fmt.Errorf("error setting display_name in state: %s", err)
	}
	if err = d.Set("labels", res.Labels); err != nil {
		return fmt.Errorf("error setting labels in state: %s", err)
	}
	if err = d.Set("project", res.Project); err != nil {
		return fmt.Errorf("error setting project in state: %s", err)
	}
	if err = d.Set("asset_status", flattenDataplexZoneAssetStatus(res.AssetStatus)); err != nil {
		return fmt.Errorf("error setting asset_status in state: %s", err)
	}
	if err = d.Set("create_time", res.CreateTime); err != nil {
		return fmt.Errorf("error setting create_time in state: %s", err)
	}
	if err = d.Set("state", res.State); err != nil {
		return fmt.Errorf("error setting state in state: %s", err)
	}
	if err = d.Set("uid", res.Uid); err != nil {
		return fmt.Errorf("error setting uid in state: %s", err)
	}
	if err = d.Set("update_time", res.UpdateTime); err != nil {
		return fmt.Errorf("error setting update_time in state: %s", err)
	}

	return nil
}
func resourceDataplexZoneUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &dataplex.Zone{
		DiscoverySpec: expandDataplexZoneDiscoverySpec(d.Get("discovery_spec")),
		Lake:          dcl.String(d.Get("lake").(string)),
		Location:      dcl.String(d.Get("location").(string)),
		Name:          dcl.String(d.Get("name").(string)),
		ResourceSpec:  expandDataplexZoneResourceSpec(d.Get("resource_spec")),
		Type:          dataplex.ZoneTypeEnumRef(d.Get("type").(string)),
		Description:   dcl.String(d.Get("description").(string)),
		DisplayName:   dcl.String(d.Get("display_name").(string)),
		Labels:        tpgresource.CheckStringMap(d.Get("labels")),
		Project:       dcl.String(project),
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
	client := transport_tpg.NewDCLDataplexClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutUpdate))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyZone(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error updating Zone: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Zone %q: %#v", d.Id(), res)

	return resourceDataplexZoneRead(d, meta)
}

func resourceDataplexZoneDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &dataplex.Zone{
		DiscoverySpec: expandDataplexZoneDiscoverySpec(d.Get("discovery_spec")),
		Lake:          dcl.String(d.Get("lake").(string)),
		Location:      dcl.String(d.Get("location").(string)),
		Name:          dcl.String(d.Get("name").(string)),
		ResourceSpec:  expandDataplexZoneResourceSpec(d.Get("resource_spec")),
		Type:          dataplex.ZoneTypeEnumRef(d.Get("type").(string)),
		Description:   dcl.String(d.Get("description").(string)),
		DisplayName:   dcl.String(d.Get("display_name").(string)),
		Labels:        tpgresource.CheckStringMap(d.Get("labels")),
		Project:       dcl.String(project),
	}

	log.Printf("[DEBUG] Deleting Zone %q", d.Id())
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLDataplexClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutDelete))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	if err := client.DeleteZone(context.Background(), obj); err != nil {
		return fmt.Errorf("Error deleting Zone: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting Zone %q", d.Id())
	return nil
}

func resourceDataplexZoneImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)

	if err := tpgresource.ParseImportId([]string{
		"projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/lakes/(?P<lake>[^/]+)/zones/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<lake>[^/]+)/(?P<name>[^/]+)",
		"(?P<location>[^/]+)/(?P<lake>[^/]+)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVarsForId(d, config, "projects/{{project}}/locations/{{location}}/lakes/{{lake}}/zones/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func expandDataplexZoneDiscoverySpec(o interface{}) *dataplex.ZoneDiscoverySpec {
	if o == nil {
		return dataplex.EmptyZoneDiscoverySpec
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataplex.EmptyZoneDiscoverySpec
	}
	obj := objArr[0].(map[string]interface{})
	return &dataplex.ZoneDiscoverySpec{
		Enabled:         dcl.Bool(obj["enabled"].(bool)),
		CsvOptions:      expandDataplexZoneDiscoverySpecCsvOptions(obj["csv_options"]),
		ExcludePatterns: tpgdclresource.ExpandStringArray(obj["exclude_patterns"]),
		IncludePatterns: tpgdclresource.ExpandStringArray(obj["include_patterns"]),
		JsonOptions:     expandDataplexZoneDiscoverySpecJsonOptions(obj["json_options"]),
		Schedule:        dcl.StringOrNil(obj["schedule"].(string)),
	}
}

func flattenDataplexZoneDiscoverySpec(obj *dataplex.ZoneDiscoverySpec) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"enabled":          obj.Enabled,
		"csv_options":      flattenDataplexZoneDiscoverySpecCsvOptions(obj.CsvOptions),
		"exclude_patterns": obj.ExcludePatterns,
		"include_patterns": obj.IncludePatterns,
		"json_options":     flattenDataplexZoneDiscoverySpecJsonOptions(obj.JsonOptions),
		"schedule":         obj.Schedule,
	}

	return []interface{}{transformed}

}

func expandDataplexZoneDiscoverySpecCsvOptions(o interface{}) *dataplex.ZoneDiscoverySpecCsvOptions {
	if o == nil {
		return nil
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return nil
	}
	obj := objArr[0].(map[string]interface{})
	return &dataplex.ZoneDiscoverySpecCsvOptions{
		Delimiter:            dcl.String(obj["delimiter"].(string)),
		DisableTypeInference: dcl.Bool(obj["disable_type_inference"].(bool)),
		Encoding:             dcl.String(obj["encoding"].(string)),
		HeaderRows:           dcl.Int64(int64(obj["header_rows"].(int))),
	}
}

func flattenDataplexZoneDiscoverySpecCsvOptions(obj *dataplex.ZoneDiscoverySpecCsvOptions) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"delimiter":              obj.Delimiter,
		"disable_type_inference": obj.DisableTypeInference,
		"encoding":               obj.Encoding,
		"header_rows":            obj.HeaderRows,
	}

	return []interface{}{transformed}

}

func expandDataplexZoneDiscoverySpecJsonOptions(o interface{}) *dataplex.ZoneDiscoverySpecJsonOptions {
	if o == nil {
		return nil
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return nil
	}
	obj := objArr[0].(map[string]interface{})
	return &dataplex.ZoneDiscoverySpecJsonOptions{
		DisableTypeInference: dcl.Bool(obj["disable_type_inference"].(bool)),
		Encoding:             dcl.String(obj["encoding"].(string)),
	}
}

func flattenDataplexZoneDiscoverySpecJsonOptions(obj *dataplex.ZoneDiscoverySpecJsonOptions) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"disable_type_inference": obj.DisableTypeInference,
		"encoding":               obj.Encoding,
	}

	return []interface{}{transformed}

}

func expandDataplexZoneResourceSpec(o interface{}) *dataplex.ZoneResourceSpec {
	if o == nil {
		return dataplex.EmptyZoneResourceSpec
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataplex.EmptyZoneResourceSpec
	}
	obj := objArr[0].(map[string]interface{})
	return &dataplex.ZoneResourceSpec{
		LocationType: dataplex.ZoneResourceSpecLocationTypeEnumRef(obj["location_type"].(string)),
	}
}

func flattenDataplexZoneResourceSpec(obj *dataplex.ZoneResourceSpec) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"location_type": obj.LocationType,
	}

	return []interface{}{transformed}

}

func flattenDataplexZoneAssetStatus(obj *dataplex.ZoneAssetStatus) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"active_assets":                   obj.ActiveAssets,
		"security_policy_applying_assets": obj.SecurityPolicyApplyingAssets,
		"update_time":                     obj.UpdateTime,
	}

	return []interface{}{transformed}

}
