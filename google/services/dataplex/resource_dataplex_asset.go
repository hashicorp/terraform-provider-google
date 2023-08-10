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

func ResourceDataplexAsset() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataplexAssetCreate,
		Read:   resourceDataplexAssetRead,
		Update: resourceDataplexAssetUpdate,
		Delete: resourceDataplexAssetDelete,

		Importer: &schema.ResourceImporter{
			State: resourceDataplexAssetImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"dataplex_zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The zone for the resource",
			},

			"discovery_spec": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Required. Specification of the discovery feature applied to data referenced by this asset. When this spec is left unset, the asset will use the spec set on the parent zone.",
				MaxItems:    1,
				Elem:        DataplexAssetDiscoverySpecSchema(),
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
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the asset.",
			},

			"resource_spec": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Required. Immutable. Specification of the resource that is referenced by this asset.",
				MaxItems:    1,
				Elem:        DataplexAssetResourceSpecSchema(),
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional. Description of the asset.",
			},

			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional. User friendly display name.",
			},

			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Optional. User defined labels for the asset.",
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

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The time when the asset was created.",
			},

			"discovery_status": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Output only. Status of the discovery feature applied to data referenced by this asset.",
				Elem:        DataplexAssetDiscoveryStatusSchema(),
			},

			"resource_status": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Output only. Status of the resource referenced by this asset.",
				Elem:        DataplexAssetResourceStatusSchema(),
			},

			"security_status": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Output only. Status of the security policy applied to resource referenced by this asset.",
				Elem:        DataplexAssetSecurityStatusSchema(),
			},

			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. Current state of the asset. Possible values: STATE_UNSPECIFIED, ACTIVE, CREATING, DELETING, ACTION_REQUIRED",
			},

			"uid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. System generated globally unique ID for the asset. This ID will be different if the asset is deleted and re-created with the same name.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The time when the asset was last updated.",
			},
		},
	}
}

func DataplexAssetDiscoverySpecSchema() *schema.Resource {
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
				Elem:        DataplexAssetDiscoverySpecCsvOptionsSchema(),
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
				Elem:        DataplexAssetDiscoverySpecJsonOptionsSchema(),
			},

			"schedule": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional. Cron schedule (https://en.wikipedia.org/wiki/Cron) for running discovery periodically. Successive discovery runs must be scheduled at least 60 minutes apart. The default value is to run discovery every 60 minutes. To explicitly set a timezone to the cron tab, apply a prefix in the cron tab: \"CRON_TZ=${IANA_TIME_ZONE}\" or TZ=${IANA_TIME_ZONE}\". The ${IANA_TIME_ZONE} may only be a valid string from IANA time zone database. For example, \"CRON_TZ=America/New_York 1 * * * *\", or \"TZ=America/New_York 1 * * * *\".",
			},
		},
	}
}

func DataplexAssetDiscoverySpecCsvOptionsSchema() *schema.Resource {
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

func DataplexAssetDiscoverySpecJsonOptionsSchema() *schema.Resource {
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

func DataplexAssetResourceSpecSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Required. Immutable. Type of resource. Possible values: STORAGE_BUCKET, BIGQUERY_DATASET",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Immutable. Relative name of the cloud resource that contains the data that is being managed within a lake. For example: `projects/{project_number}/buckets/{bucket_id}` `projects/{project_number}/datasets/{dataset_id}`",
			},

			"read_access_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Optional. Determines how read permissions are handled for each asset and their associated tables. Only available to storage buckets assets. Possible values: DIRECT, MANAGED",
			},
		},
	}
}

func DataplexAssetDiscoveryStatusSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"last_run_duration": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The duration of the last discovery run.",
			},

			"last_run_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The start time of the last discovery run.",
			},

			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Additional information about the current state.",
			},

			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the discovery feature. Possible values: STATE_UNSPECIFIED, SCHEDULED, IN_PROGRESS, PAUSED, DISABLED",
			},

			"stats": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Data Stats of the asset reported by discovery.",
				Elem:        DataplexAssetDiscoveryStatusStatsSchema(),
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update time of the status.",
			},
		},
	}
}

func DataplexAssetDiscoveryStatusStatsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"data_items": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The count of data items within the referenced resource.",
			},

			"data_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of stored data bytes within the referenced resource.",
			},

			"filesets": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The count of fileset entities within the referenced resource.",
			},

			"tables": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The count of table entities within the referenced resource.",
			},
		},
	}
}

func DataplexAssetResourceStatusSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Additional information about the current state.",
			},

			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current state of the managed resource. Possible values: STATE_UNSPECIFIED, READY, ERROR",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update time of the status.",
			},
		},
	}
}

func DataplexAssetSecurityStatusSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Additional information about the current state.",
			},

			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current state of the security policy applied to the attached resource. Possible values: STATE_UNSPECIFIED, READY, APPLYING, ERROR",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update time of the status.",
			},
		},
	}
}

func resourceDataplexAssetCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &dataplex.Asset{
		DataplexZone:  dcl.String(d.Get("dataplex_zone").(string)),
		DiscoverySpec: expandDataplexAssetDiscoverySpec(d.Get("discovery_spec")),
		Lake:          dcl.String(d.Get("lake").(string)),
		Location:      dcl.String(d.Get("location").(string)),
		Name:          dcl.String(d.Get("name").(string)),
		ResourceSpec:  expandDataplexAssetResourceSpec(d.Get("resource_spec")),
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
	res, err := client.ApplyAsset(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating Asset: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Asset %q: %#v", d.Id(), res)

	return resourceDataplexAssetRead(d, meta)
}

func resourceDataplexAssetRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &dataplex.Asset{
		DataplexZone:  dcl.String(d.Get("dataplex_zone").(string)),
		DiscoverySpec: expandDataplexAssetDiscoverySpec(d.Get("discovery_spec")),
		Lake:          dcl.String(d.Get("lake").(string)),
		Location:      dcl.String(d.Get("location").(string)),
		Name:          dcl.String(d.Get("name").(string)),
		ResourceSpec:  expandDataplexAssetResourceSpec(d.Get("resource_spec")),
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
	res, err := client.GetAsset(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("DataplexAsset %q", d.Id())
		return tpgdclresource.HandleNotFoundDCLError(err, d, resourceName)
	}

	if err = d.Set("dataplex_zone", res.DataplexZone); err != nil {
		return fmt.Errorf("error setting dataplex_zone in state: %s", err)
	}
	if err = d.Set("discovery_spec", flattenDataplexAssetDiscoverySpec(res.DiscoverySpec)); err != nil {
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
	if err = d.Set("resource_spec", flattenDataplexAssetResourceSpec(res.ResourceSpec)); err != nil {
		return fmt.Errorf("error setting resource_spec in state: %s", err)
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
	if err = d.Set("create_time", res.CreateTime); err != nil {
		return fmt.Errorf("error setting create_time in state: %s", err)
	}
	if err = d.Set("discovery_status", flattenDataplexAssetDiscoveryStatus(res.DiscoveryStatus)); err != nil {
		return fmt.Errorf("error setting discovery_status in state: %s", err)
	}
	if err = d.Set("resource_status", flattenDataplexAssetResourceStatus(res.ResourceStatus)); err != nil {
		return fmt.Errorf("error setting resource_status in state: %s", err)
	}
	if err = d.Set("security_status", flattenDataplexAssetSecurityStatus(res.SecurityStatus)); err != nil {
		return fmt.Errorf("error setting security_status in state: %s", err)
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
func resourceDataplexAssetUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &dataplex.Asset{
		DataplexZone:  dcl.String(d.Get("dataplex_zone").(string)),
		DiscoverySpec: expandDataplexAssetDiscoverySpec(d.Get("discovery_spec")),
		Lake:          dcl.String(d.Get("lake").(string)),
		Location:      dcl.String(d.Get("location").(string)),
		Name:          dcl.String(d.Get("name").(string)),
		ResourceSpec:  expandDataplexAssetResourceSpec(d.Get("resource_spec")),
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
	res, err := client.ApplyAsset(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error updating Asset: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Asset %q: %#v", d.Id(), res)

	return resourceDataplexAssetRead(d, meta)
}

func resourceDataplexAssetDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &dataplex.Asset{
		DataplexZone:  dcl.String(d.Get("dataplex_zone").(string)),
		DiscoverySpec: expandDataplexAssetDiscoverySpec(d.Get("discovery_spec")),
		Lake:          dcl.String(d.Get("lake").(string)),
		Location:      dcl.String(d.Get("location").(string)),
		Name:          dcl.String(d.Get("name").(string)),
		ResourceSpec:  expandDataplexAssetResourceSpec(d.Get("resource_spec")),
		Description:   dcl.String(d.Get("description").(string)),
		DisplayName:   dcl.String(d.Get("display_name").(string)),
		Labels:        tpgresource.CheckStringMap(d.Get("labels")),
		Project:       dcl.String(project),
	}

	log.Printf("[DEBUG] Deleting Asset %q", d.Id())
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
	if err := client.DeleteAsset(context.Background(), obj); err != nil {
		return fmt.Errorf("Error deleting Asset: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting Asset %q", d.Id())
	return nil
}

func resourceDataplexAssetImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)

	if err := tpgresource.ParseImportId([]string{
		"projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/lakes/(?P<lake>[^/]+)/zones/(?P<dataplex_zone>[^/]+)/assets/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<lake>[^/]+)/(?P<dataplex_zone>[^/]+)/(?P<name>[^/]+)",
		"(?P<location>[^/]+)/(?P<lake>[^/]+)/(?P<dataplex_zone>[^/]+)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVarsForId(d, config, "projects/{{project}}/locations/{{location}}/lakes/{{lake}}/zones/{{dataplex_zone}}/assets/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func expandDataplexAssetDiscoverySpec(o interface{}) *dataplex.AssetDiscoverySpec {
	if o == nil {
		return dataplex.EmptyAssetDiscoverySpec
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataplex.EmptyAssetDiscoverySpec
	}
	obj := objArr[0].(map[string]interface{})
	return &dataplex.AssetDiscoverySpec{
		Enabled:         dcl.Bool(obj["enabled"].(bool)),
		CsvOptions:      expandDataplexAssetDiscoverySpecCsvOptions(obj["csv_options"]),
		ExcludePatterns: tpgdclresource.ExpandStringArray(obj["exclude_patterns"]),
		IncludePatterns: tpgdclresource.ExpandStringArray(obj["include_patterns"]),
		JsonOptions:     expandDataplexAssetDiscoverySpecJsonOptions(obj["json_options"]),
		Schedule:        dcl.String(obj["schedule"].(string)),
	}
}

func flattenDataplexAssetDiscoverySpec(obj *dataplex.AssetDiscoverySpec) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"enabled":          obj.Enabled,
		"csv_options":      flattenDataplexAssetDiscoverySpecCsvOptions(obj.CsvOptions),
		"exclude_patterns": obj.ExcludePatterns,
		"include_patterns": obj.IncludePatterns,
		"json_options":     flattenDataplexAssetDiscoverySpecJsonOptions(obj.JsonOptions),
		"schedule":         obj.Schedule,
	}

	return []interface{}{transformed}

}

func expandDataplexAssetDiscoverySpecCsvOptions(o interface{}) *dataplex.AssetDiscoverySpecCsvOptions {
	if o == nil {
		return nil
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return nil
	}
	obj := objArr[0].(map[string]interface{})
	return &dataplex.AssetDiscoverySpecCsvOptions{
		Delimiter:            dcl.String(obj["delimiter"].(string)),
		DisableTypeInference: dcl.Bool(obj["disable_type_inference"].(bool)),
		Encoding:             dcl.String(obj["encoding"].(string)),
		HeaderRows:           dcl.Int64(int64(obj["header_rows"].(int))),
	}
}

func flattenDataplexAssetDiscoverySpecCsvOptions(obj *dataplex.AssetDiscoverySpecCsvOptions) interface{} {
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

func expandDataplexAssetDiscoverySpecJsonOptions(o interface{}) *dataplex.AssetDiscoverySpecJsonOptions {
	if o == nil {
		return nil
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return nil
	}
	obj := objArr[0].(map[string]interface{})
	return &dataplex.AssetDiscoverySpecJsonOptions{
		DisableTypeInference: dcl.Bool(obj["disable_type_inference"].(bool)),
		Encoding:             dcl.String(obj["encoding"].(string)),
	}
}

func flattenDataplexAssetDiscoverySpecJsonOptions(obj *dataplex.AssetDiscoverySpecJsonOptions) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"disable_type_inference": obj.DisableTypeInference,
		"encoding":               obj.Encoding,
	}

	return []interface{}{transformed}

}

func expandDataplexAssetResourceSpec(o interface{}) *dataplex.AssetResourceSpec {
	if o == nil {
		return dataplex.EmptyAssetResourceSpec
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataplex.EmptyAssetResourceSpec
	}
	obj := objArr[0].(map[string]interface{})
	return &dataplex.AssetResourceSpec{
		Type:           dataplex.AssetResourceSpecTypeEnumRef(obj["type"].(string)),
		Name:           dcl.String(obj["name"].(string)),
		ReadAccessMode: dataplex.AssetResourceSpecReadAccessModeEnumRef(obj["read_access_mode"].(string)),
	}
}

func flattenDataplexAssetResourceSpec(obj *dataplex.AssetResourceSpec) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"type":             obj.Type,
		"name":             obj.Name,
		"read_access_mode": obj.ReadAccessMode,
	}

	return []interface{}{transformed}

}

func flattenDataplexAssetDiscoveryStatus(obj *dataplex.AssetDiscoveryStatus) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"last_run_duration": obj.LastRunDuration,
		"last_run_time":     obj.LastRunTime,
		"message":           obj.Message,
		"state":             obj.State,
		"stats":             flattenDataplexAssetDiscoveryStatusStats(obj.Stats),
		"update_time":       obj.UpdateTime,
	}

	return []interface{}{transformed}

}

func flattenDataplexAssetDiscoveryStatusStats(obj *dataplex.AssetDiscoveryStatusStats) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"data_items": obj.DataItems,
		"data_size":  obj.DataSize,
		"filesets":   obj.Filesets,
		"tables":     obj.Tables,
	}

	return []interface{}{transformed}

}

func flattenDataplexAssetResourceStatus(obj *dataplex.AssetResourceStatus) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"message":     obj.Message,
		"state":       obj.State,
		"update_time": obj.UpdateTime,
	}

	return []interface{}{transformed}

}

func flattenDataplexAssetSecurityStatus(obj *dataplex.AssetSecurityStatus) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"message":     obj.Message,
		"state":       obj.State,
		"update_time": obj.UpdateTime,
	}

	return []interface{}{transformed}

}
