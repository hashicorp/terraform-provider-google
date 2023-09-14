// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package storageinsights

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"
)

func ResourceStorageInsightsReportConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceStorageInsightsReportConfigCreate,
		Read:   resourceStorageInsightsReportConfigRead,
		Update: resourceStorageInsightsReportConfigUpdate,
		Delete: resourceStorageInsightsReportConfigDelete,

		Importer: &schema.ResourceImporter{
			State: resourceStorageInsightsReportConfigImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			tpgresource.DefaultProviderProject,
		),

		Schema: map[string]*schema.Schema{
			"csv_options": {
				Type:        schema.TypeList,
				Required:    true,
				Description: `Options for configuring the format of the inventory report CSV file.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delimiter": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The delimiter used to separate the fields in the inventory report CSV file.`,
						},
						"header_required": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `The boolean that indicates whether or not headers are included in the inventory report CSV file.`,
						},
						"record_separator": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The character used to separate the records in the inventory report CSV file.`,
						},
					},
				},
			},
			"location": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: `The location of the ReportConfig. The source and destination buckets specified in the ReportConfig
must be in the same location.`,
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The editable display name of the inventory report configuration. Has a limit of 256 characters. Can be empty.`,
			},
			"frequency_options": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Options for configuring how inventory reports are generated.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"end_date": {
							Type:        schema.TypeList,
							Required:    true,
							Description: `The date to stop generating inventory reports. For example, {"day": 15, "month": 9, "year": 2022}.`,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"day": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: `The day of the month to stop generating inventory reports.`,
									},
									"month": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: `The month to stop generating inventory reports.`,
									},
									"year": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: `The year to stop generating inventory reports`,
									},
								},
							},
						},
						"frequency": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: verify.ValidateEnum([]string{"DAILY", "WEEKLY"}),
							Description:  `The frequency in which inventory reports are generated. Values are DAILY or WEEKLY. Possible values: ["DAILY", "WEEKLY"]`,
						},
						"start_date": {
							Type:        schema.TypeList,
							Required:    true,
							Description: `The date to start generating inventory reports. For example, {"day": 15, "month": 8, "year": 2022}.`,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"day": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: `The day of the month to start generating inventory reports.`,
									},
									"month": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: `The month to start generating inventory reports.`,
									},
									"year": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: `The year to start generating inventory reports`,
									},
								},
							},
						},
					},
				},
			},
			"object_metadata_report_options": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Options for including metadata in an inventory report.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metadata_fields": {
							Type:        schema.TypeList,
							Required:    true,
							Description: `The metadata fields included in an inventory report.`,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"storage_destination_options": {
							Type:        schema.TypeList,
							Required:    true,
							Description: `Options for where the inventory reports are stored.`,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The destination bucket that stores the generated inventory reports.`,
									},
									"destination_path": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The path within the destination bucket to store generated inventory reports.`,
									},
								},
							},
						},
						"storage_filters": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `A nested object resource`,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: `The filter to use when specifying which bucket to generate inventory reports for.`,
									},
								},
							},
						},
					},
				},
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The UUID of the inventory report configuration.`,
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

func resourceStorageInsightsReportConfigCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	frequencyOptionsProp, err := expandStorageInsightsReportConfigFrequencyOptions(d.Get("frequency_options"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("frequency_options"); !tpgresource.IsEmptyValue(reflect.ValueOf(frequencyOptionsProp)) && (ok || !reflect.DeepEqual(v, frequencyOptionsProp)) {
		obj["frequencyOptions"] = frequencyOptionsProp
	}
	csvOptionsProp, err := expandStorageInsightsReportConfigCsvOptions(d.Get("csv_options"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("csv_options"); !tpgresource.IsEmptyValue(reflect.ValueOf(csvOptionsProp)) && (ok || !reflect.DeepEqual(v, csvOptionsProp)) {
		obj["csvOptions"] = csvOptionsProp
	}
	objectMetadataReportOptionsProp, err := expandStorageInsightsReportConfigObjectMetadataReportOptions(d.Get("object_metadata_report_options"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("object_metadata_report_options"); !tpgresource.IsEmptyValue(reflect.ValueOf(objectMetadataReportOptionsProp)) && (ok || !reflect.DeepEqual(v, objectMetadataReportOptionsProp)) {
		obj["objectMetadataReportOptions"] = objectMetadataReportOptionsProp
	}
	displayNameProp, err := expandStorageInsightsReportConfigDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("display_name"); !tpgresource.IsEmptyValue(reflect.ValueOf(displayNameProp)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{StorageInsightsBasePath}}projects/{{project}}/locations/{{location}}/reportConfigs")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new ReportConfig: %#v", obj)
	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for ReportConfig: %s", err)
	}
	billingProject = project

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
		Timeout:   d.Timeout(schema.TimeoutCreate),
	})
	if err != nil {
		return fmt.Errorf("Error creating ReportConfig: %s", err)
	}
	if err := d.Set("name", flattenStorageInsightsReportConfigName(res["name"], d, config)); err != nil {
		return fmt.Errorf(`Error setting computed identity field "name": %s`, err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/reportConfigs/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating ReportConfig %q: %#v", d.Id(), res)

	return resourceStorageInsightsReportConfigRead(d, meta)
}

func resourceStorageInsightsReportConfigRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{StorageInsightsBasePath}}projects/{{project}}/locations/{{location}}/reportConfigs/{{name}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for ReportConfig: %s", err)
	}
	billingProject = project

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
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("StorageInsightsReportConfig %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading ReportConfig: %s", err)
	}

	if err := d.Set("name", flattenStorageInsightsReportConfigName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading ReportConfig: %s", err)
	}
	if err := d.Set("frequency_options", flattenStorageInsightsReportConfigFrequencyOptions(res["frequencyOptions"], d, config)); err != nil {
		return fmt.Errorf("Error reading ReportConfig: %s", err)
	}
	if err := d.Set("csv_options", flattenStorageInsightsReportConfigCsvOptions(res["csvOptions"], d, config)); err != nil {
		return fmt.Errorf("Error reading ReportConfig: %s", err)
	}
	if err := d.Set("object_metadata_report_options", flattenStorageInsightsReportConfigObjectMetadataReportOptions(res["objectMetadataReportOptions"], d, config)); err != nil {
		return fmt.Errorf("Error reading ReportConfig: %s", err)
	}
	if err := d.Set("display_name", flattenStorageInsightsReportConfigDisplayName(res["displayName"], d, config)); err != nil {
		return fmt.Errorf("Error reading ReportConfig: %s", err)
	}

	return nil
}

func resourceStorageInsightsReportConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for ReportConfig: %s", err)
	}
	billingProject = project

	obj := make(map[string]interface{})
	frequencyOptionsProp, err := expandStorageInsightsReportConfigFrequencyOptions(d.Get("frequency_options"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("frequency_options"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, frequencyOptionsProp)) {
		obj["frequencyOptions"] = frequencyOptionsProp
	}
	csvOptionsProp, err := expandStorageInsightsReportConfigCsvOptions(d.Get("csv_options"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("csv_options"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, csvOptionsProp)) {
		obj["csvOptions"] = csvOptionsProp
	}
	objectMetadataReportOptionsProp, err := expandStorageInsightsReportConfigObjectMetadataReportOptions(d.Get("object_metadata_report_options"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("object_metadata_report_options"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, objectMetadataReportOptionsProp)) {
		obj["objectMetadataReportOptions"] = objectMetadataReportOptionsProp
	}
	displayNameProp, err := expandStorageInsightsReportConfigDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("display_name"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{StorageInsightsBasePath}}projects/{{project}}/locations/{{location}}/reportConfigs/{{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating ReportConfig %q: %#v", d.Id(), obj)
	updateMask := []string{}

	if d.HasChange("frequency_options") {
		updateMask = append(updateMask, "frequencyOptions")
	}

	if d.HasChange("csv_options") {
		updateMask = append(updateMask, "csvOptions")
	}

	if d.HasChange("object_metadata_report_options") {
		updateMask = append(updateMask, "objectMetadataReportOptions.metadataFields",
			"objectMetadataReportOptions.storageDestinationOptions.bucket",
			"objectMetadataReportOptions.storageDestinationOptions.destinationPath")
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
		return fmt.Errorf("Error updating ReportConfig %q: %s", d.Id(), err)
	} else {
		log.Printf("[DEBUG] Finished updating ReportConfig %q: %#v", d.Id(), res)
	}

	return resourceStorageInsightsReportConfigRead(d, meta)
}

func resourceStorageInsightsReportConfigDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for ReportConfig: %s", err)
	}
	billingProject = project

	url, err := tpgresource.ReplaceVars(d, config, "{{StorageInsightsBasePath}}projects/{{project}}/locations/{{location}}/reportConfigs/{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting ReportConfig %q", d.Id())

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
		return transport_tpg.HandleNotFoundError(err, d, "ReportConfig")
	}

	log.Printf("[DEBUG] Finished deleting ReportConfig %q: %#v", d.Id(), res)
	return nil
}

func resourceStorageInsightsReportConfigImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/reportConfigs/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<name>[^/]+)",
		"(?P<location>[^/]+)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/reportConfigs/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenStorageInsightsReportConfigName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	return tpgresource.NameFromSelfLinkStateFunc(v)
}

func flattenStorageInsightsReportConfigFrequencyOptions(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["frequency"] =
		flattenStorageInsightsReportConfigFrequencyOptionsFrequency(original["frequency"], d, config)
	transformed["start_date"] =
		flattenStorageInsightsReportConfigFrequencyOptionsStartDate(original["startDate"], d, config)
	transformed["end_date"] =
		flattenStorageInsightsReportConfigFrequencyOptionsEndDate(original["endDate"], d, config)
	return []interface{}{transformed}
}
func flattenStorageInsightsReportConfigFrequencyOptionsFrequency(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenStorageInsightsReportConfigFrequencyOptionsStartDate(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["day"] =
		flattenStorageInsightsReportConfigFrequencyOptionsStartDateDay(original["day"], d, config)
	transformed["month"] =
		flattenStorageInsightsReportConfigFrequencyOptionsStartDateMonth(original["month"], d, config)
	transformed["year"] =
		flattenStorageInsightsReportConfigFrequencyOptionsStartDateYear(original["year"], d, config)
	return []interface{}{transformed}
}
func flattenStorageInsightsReportConfigFrequencyOptionsStartDateDay(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenStorageInsightsReportConfigFrequencyOptionsStartDateMonth(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenStorageInsightsReportConfigFrequencyOptionsStartDateYear(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenStorageInsightsReportConfigFrequencyOptionsEndDate(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["day"] =
		flattenStorageInsightsReportConfigFrequencyOptionsEndDateDay(original["day"], d, config)
	transformed["month"] =
		flattenStorageInsightsReportConfigFrequencyOptionsEndDateMonth(original["month"], d, config)
	transformed["year"] =
		flattenStorageInsightsReportConfigFrequencyOptionsEndDateYear(original["year"], d, config)
	return []interface{}{transformed}
}
func flattenStorageInsightsReportConfigFrequencyOptionsEndDateDay(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenStorageInsightsReportConfigFrequencyOptionsEndDateMonth(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenStorageInsightsReportConfigFrequencyOptionsEndDateYear(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenStorageInsightsReportConfigCsvOptions(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["record_separator"] =
		flattenStorageInsightsReportConfigCsvOptionsRecordSeparator(original["recordSeparator"], d, config)
	transformed["delimiter"] =
		flattenStorageInsightsReportConfigCsvOptionsDelimiter(original["delimiter"], d, config)
	transformed["header_required"] =
		flattenStorageInsightsReportConfigCsvOptionsHeaderRequired(original["headerRequired"], d, config)
	return []interface{}{transformed}
}
func flattenStorageInsightsReportConfigCsvOptionsRecordSeparator(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenStorageInsightsReportConfigCsvOptionsDelimiter(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenStorageInsightsReportConfigCsvOptionsHeaderRequired(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenStorageInsightsReportConfigObjectMetadataReportOptions(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["metadata_fields"] =
		flattenStorageInsightsReportConfigObjectMetadataReportOptionsMetadataFields(original["metadataFields"], d, config)
	transformed["storage_filters"] =
		flattenStorageInsightsReportConfigObjectMetadataReportOptionsStorageFilters(original["storageFilters"], d, config)
	transformed["storage_destination_options"] =
		flattenStorageInsightsReportConfigObjectMetadataReportOptionsStorageDestinationOptions(original["storageDestinationOptions"], d, config)
	return []interface{}{transformed}
}
func flattenStorageInsightsReportConfigObjectMetadataReportOptionsMetadataFields(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenStorageInsightsReportConfigObjectMetadataReportOptionsStorageFilters(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["bucket"] =
		flattenStorageInsightsReportConfigObjectMetadataReportOptionsStorageFiltersBucket(original["bucket"], d, config)
	return []interface{}{transformed}
}
func flattenStorageInsightsReportConfigObjectMetadataReportOptionsStorageFiltersBucket(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenStorageInsightsReportConfigObjectMetadataReportOptionsStorageDestinationOptions(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["bucket"] =
		flattenStorageInsightsReportConfigObjectMetadataReportOptionsStorageDestinationOptionsBucket(original["bucket"], d, config)
	transformed["destination_path"] =
		flattenStorageInsightsReportConfigObjectMetadataReportOptionsStorageDestinationOptionsDestinationPath(original["destinationPath"], d, config)
	return []interface{}{transformed}
}
func flattenStorageInsightsReportConfigObjectMetadataReportOptionsStorageDestinationOptionsBucket(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenStorageInsightsReportConfigObjectMetadataReportOptionsStorageDestinationOptionsDestinationPath(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenStorageInsightsReportConfigDisplayName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandStorageInsightsReportConfigFrequencyOptions(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedFrequency, err := expandStorageInsightsReportConfigFrequencyOptionsFrequency(original["frequency"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedFrequency); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["frequency"] = transformedFrequency
	}

	transformedStartDate, err := expandStorageInsightsReportConfigFrequencyOptionsStartDate(original["start_date"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedStartDate); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["startDate"] = transformedStartDate
	}

	transformedEndDate, err := expandStorageInsightsReportConfigFrequencyOptionsEndDate(original["end_date"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedEndDate); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["endDate"] = transformedEndDate
	}

	return transformed, nil
}

func expandStorageInsightsReportConfigFrequencyOptionsFrequency(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandStorageInsightsReportConfigFrequencyOptionsStartDate(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedDay, err := expandStorageInsightsReportConfigFrequencyOptionsStartDateDay(original["day"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDay); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["day"] = transformedDay
	}

	transformedMonth, err := expandStorageInsightsReportConfigFrequencyOptionsStartDateMonth(original["month"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMonth); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["month"] = transformedMonth
	}

	transformedYear, err := expandStorageInsightsReportConfigFrequencyOptionsStartDateYear(original["year"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedYear); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["year"] = transformedYear
	}

	return transformed, nil
}

func expandStorageInsightsReportConfigFrequencyOptionsStartDateDay(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandStorageInsightsReportConfigFrequencyOptionsStartDateMonth(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandStorageInsightsReportConfigFrequencyOptionsStartDateYear(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandStorageInsightsReportConfigFrequencyOptionsEndDate(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedDay, err := expandStorageInsightsReportConfigFrequencyOptionsEndDateDay(original["day"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDay); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["day"] = transformedDay
	}

	transformedMonth, err := expandStorageInsightsReportConfigFrequencyOptionsEndDateMonth(original["month"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMonth); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["month"] = transformedMonth
	}

	transformedYear, err := expandStorageInsightsReportConfigFrequencyOptionsEndDateYear(original["year"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedYear); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["year"] = transformedYear
	}

	return transformed, nil
}

func expandStorageInsightsReportConfigFrequencyOptionsEndDateDay(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandStorageInsightsReportConfigFrequencyOptionsEndDateMonth(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandStorageInsightsReportConfigFrequencyOptionsEndDateYear(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandStorageInsightsReportConfigCsvOptions(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedRecordSeparator, err := expandStorageInsightsReportConfigCsvOptionsRecordSeparator(original["record_separator"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedRecordSeparator); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["recordSeparator"] = transformedRecordSeparator
	}

	transformedDelimiter, err := expandStorageInsightsReportConfigCsvOptionsDelimiter(original["delimiter"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDelimiter); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["delimiter"] = transformedDelimiter
	}

	transformedHeaderRequired, err := expandStorageInsightsReportConfigCsvOptionsHeaderRequired(original["header_required"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedHeaderRequired); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["headerRequired"] = transformedHeaderRequired
	}

	return transformed, nil
}

func expandStorageInsightsReportConfigCsvOptionsRecordSeparator(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandStorageInsightsReportConfigCsvOptionsDelimiter(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandStorageInsightsReportConfigCsvOptionsHeaderRequired(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandStorageInsightsReportConfigObjectMetadataReportOptions(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedMetadataFields, err := expandStorageInsightsReportConfigObjectMetadataReportOptionsMetadataFields(original["metadata_fields"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMetadataFields); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["metadataFields"] = transformedMetadataFields
	}

	transformedStorageFilters, err := expandStorageInsightsReportConfigObjectMetadataReportOptionsStorageFilters(original["storage_filters"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedStorageFilters); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["storageFilters"] = transformedStorageFilters
	}

	transformedStorageDestinationOptions, err := expandStorageInsightsReportConfigObjectMetadataReportOptionsStorageDestinationOptions(original["storage_destination_options"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedStorageDestinationOptions); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["storageDestinationOptions"] = transformedStorageDestinationOptions
	}

	return transformed, nil
}

func expandStorageInsightsReportConfigObjectMetadataReportOptionsMetadataFields(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandStorageInsightsReportConfigObjectMetadataReportOptionsStorageFilters(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedBucket, err := expandStorageInsightsReportConfigObjectMetadataReportOptionsStorageFiltersBucket(original["bucket"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedBucket); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["bucket"] = transformedBucket
	}

	return transformed, nil
}

func expandStorageInsightsReportConfigObjectMetadataReportOptionsStorageFiltersBucket(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandStorageInsightsReportConfigObjectMetadataReportOptionsStorageDestinationOptions(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedBucket, err := expandStorageInsightsReportConfigObjectMetadataReportOptionsStorageDestinationOptionsBucket(original["bucket"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedBucket); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["bucket"] = transformedBucket
	}

	transformedDestinationPath, err := expandStorageInsightsReportConfigObjectMetadataReportOptionsStorageDestinationOptionsDestinationPath(original["destination_path"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDestinationPath); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["destinationPath"] = transformedDestinationPath
	}

	return transformed, nil
}

func expandStorageInsightsReportConfigObjectMetadataReportOptionsStorageDestinationOptionsBucket(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandStorageInsightsReportConfigObjectMetadataReportOptionsStorageDestinationOptionsDestinationPath(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandStorageInsightsReportConfigDisplayName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
