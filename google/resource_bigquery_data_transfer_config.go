// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
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

package google

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBigqueryDataTransferConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigqueryDataTransferConfigCreate,
		Read:   resourceBigqueryDataTransferConfigRead,
		Update: resourceBigqueryDataTransferConfigUpdate,
		Delete: resourceBigqueryDataTransferConfigDelete,

		Importer: &schema.ResourceImporter{
			State: resourceBigqueryDataTransferConfigImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(240 * time.Second),
			Update: schema.DefaultTimeout(240 * time.Second),
			Delete: schema.DefaultTimeout(240 * time.Second),
		},

		Schema: map[string]*schema.Schema{
			"data_source_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"destination_dataset_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"params": {
				Type:     schema.TypeMap,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"data_refresh_window_days": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"location": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "US",
			},
			"schedule": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceBigqueryDataTransferConfigCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	obj := make(map[string]interface{})
	displayNameProp, err := expandBigqueryDataTransferConfigDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("display_name"); !isEmptyValue(reflect.ValueOf(displayNameProp)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}
	destinationDatasetIdProp, err := expandBigqueryDataTransferConfigDestinationDatasetId(d.Get("destination_dataset_id"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("destination_dataset_id"); !isEmptyValue(reflect.ValueOf(destinationDatasetIdProp)) && (ok || !reflect.DeepEqual(v, destinationDatasetIdProp)) {
		obj["destinationDatasetId"] = destinationDatasetIdProp
	}
	dataSourceIdProp, err := expandBigqueryDataTransferConfigDataSourceId(d.Get("data_source_id"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("data_source_id"); !isEmptyValue(reflect.ValueOf(dataSourceIdProp)) && (ok || !reflect.DeepEqual(v, dataSourceIdProp)) {
		obj["dataSourceId"] = dataSourceIdProp
	}
	scheduleProp, err := expandBigqueryDataTransferConfigSchedule(d.Get("schedule"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("schedule"); !isEmptyValue(reflect.ValueOf(scheduleProp)) && (ok || !reflect.DeepEqual(v, scheduleProp)) {
		obj["schedule"] = scheduleProp
	}
	dataRefreshWindowDaysProp, err := expandBigqueryDataTransferConfigDataRefreshWindowDays(d.Get("data_refresh_window_days"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("data_refresh_window_days"); !isEmptyValue(reflect.ValueOf(dataRefreshWindowDaysProp)) && (ok || !reflect.DeepEqual(v, dataRefreshWindowDaysProp)) {
		obj["dataRefreshWindowDays"] = dataRefreshWindowDaysProp
	}
	disabledProp, err := expandBigqueryDataTransferConfigDisabled(d.Get("disabled"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("disabled"); !isEmptyValue(reflect.ValueOf(disabledProp)) && (ok || !reflect.DeepEqual(v, disabledProp)) {
		obj["disabled"] = disabledProp
	}
	paramsProp, err := expandBigqueryDataTransferConfigParams(d.Get("params"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("params"); !isEmptyValue(reflect.ValueOf(paramsProp)) && (ok || !reflect.DeepEqual(v, paramsProp)) {
		obj["params"] = paramsProp
	}

	url, err := replaceVars(d, config, "https://bigquerydatatransfer.googleapis.com/v1/projects/{{project}}/locations/{{location}}/transferConfigs")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new Config: %#v", obj)
	res, err := sendRequestWithTimeout(config, "POST", url, obj, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating Config: %s", err)
	}

	// Store the ID now
	id, err := replaceVars(d, config, "{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating Config %q: %#v", d.Id(), res)

	return resourceBigqueryDataTransferConfigRead(d, meta)
}

func resourceBigqueryDataTransferConfigRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	url, err := replaceVars(d, config, "https://bigquerydatatransfer.googleapis.com/v1/projects/{{project}}/locations/{{location}}/transferConfigs/{{name}}")
	if err != nil {
		return err
	}

	res, err := sendRequest(config, "GET", url, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("BigqueryDataTransferConfig %q", d.Id()))
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading Config: %s", err)
	}

	if err := d.Set("display_name", flattenBigqueryDataTransferConfigDisplayName(res["displayName"], d)); err != nil {
		return fmt.Errorf("Error reading Config: %s", err)
	}
	if err := d.Set("destination_dataset_id", flattenBigqueryDataTransferConfigDestinationDatasetId(res["destinationDatasetId"], d)); err != nil {
		return fmt.Errorf("Error reading Config: %s", err)
	}
	if err := d.Set("data_source_id", flattenBigqueryDataTransferConfigDataSourceId(res["dataSourceId"], d)); err != nil {
		return fmt.Errorf("Error reading Config: %s", err)
	}
	if err := d.Set("schedule", flattenBigqueryDataTransferConfigSchedule(res["schedule"], d)); err != nil {
		return fmt.Errorf("Error reading Config: %s", err)
	}
	if err := d.Set("data_refresh_window_days", flattenBigqueryDataTransferConfigDataRefreshWindowDays(res["dataRefreshWindowDays"], d)); err != nil {
		return fmt.Errorf("Error reading Config: %s", err)
	}
	if err := d.Set("disabled", flattenBigqueryDataTransferConfigDisabled(res["disabled"], d)); err != nil {
		return fmt.Errorf("Error reading Config: %s", err)
	}
	if err := d.Set("params", flattenBigqueryDataTransferConfigParams(res["params"], d)); err != nil {
		return fmt.Errorf("Error reading Config: %s", err)
	}

	return nil
}

func resourceBigqueryDataTransferConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	obj := make(map[string]interface{})
	displayNameProp, err := expandBigqueryDataTransferConfigDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("display_name"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}
	destinationDatasetIdProp, err := expandBigqueryDataTransferConfigDestinationDatasetId(d.Get("destination_dataset_id"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("destination_dataset_id"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, destinationDatasetIdProp)) {
		obj["destinationDatasetId"] = destinationDatasetIdProp
	}
	dataSourceIdProp, err := expandBigqueryDataTransferConfigDataSourceId(d.Get("data_source_id"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("data_source_id"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, dataSourceIdProp)) {
		obj["dataSourceId"] = dataSourceIdProp
	}
	scheduleProp, err := expandBigqueryDataTransferConfigSchedule(d.Get("schedule"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("schedule"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, scheduleProp)) {
		obj["schedule"] = scheduleProp
	}
	dataRefreshWindowDaysProp, err := expandBigqueryDataTransferConfigDataRefreshWindowDays(d.Get("data_refresh_window_days"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("data_refresh_window_days"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, dataRefreshWindowDaysProp)) {
		obj["dataRefreshWindowDays"] = dataRefreshWindowDaysProp
	}
	disabledProp, err := expandBigqueryDataTransferConfigDisabled(d.Get("disabled"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("disabled"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, disabledProp)) {
		obj["disabled"] = disabledProp
	}
	paramsProp, err := expandBigqueryDataTransferConfigParams(d.Get("params"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("params"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, paramsProp)) {
		obj["params"] = paramsProp
	}

	url, err := replaceVars(d, config, "https://bigquerydatatransfer.googleapis.com/v1/projects/{{project}}/locations/{{location}}/transferConfigs/{{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating Config %q: %#v", d.Id(), obj)
	_, err = sendRequestWithTimeout(config, "PUT", url, obj, d.Timeout(schema.TimeoutUpdate))

	if err != nil {
		return fmt.Errorf("Error updating Config %q: %s", d.Id(), err)
	}

	return resourceBigqueryDataTransferConfigRead(d, meta)
}

func resourceBigqueryDataTransferConfigDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	url, err := replaceVars(d, config, "https://bigquerydatatransfer.googleapis.com/v1/projects/{{project}}/locations/{{location}}/transferConfigs/{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting Config %q", d.Id())
	res, err := sendRequestWithTimeout(config, "DELETE", url, obj, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return handleNotFoundError(err, d, "Config")
	}

	log.Printf("[DEBUG] Finished deleting Config %q: %#v", d.Id(), res)
	return nil
}

func resourceBigqueryDataTransferConfigImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{"projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/transferConfigs/(?P<name>[^/]+)", "(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<name>[^/]+)", "(?P<location>[^/]+)/(?P<name>[^/]+)"}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenBigqueryDataTransferConfigDisplayName(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenBigqueryDataTransferConfigDestinationDatasetId(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenBigqueryDataTransferConfigDataSourceId(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenBigqueryDataTransferConfigSchedule(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenBigqueryDataTransferConfigDataRefreshWindowDays(v interface{}, d *schema.ResourceData) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := strconv.ParseInt(strVal, 10, 64); err == nil {
			return intVal
		} // let terraform core handle it if we can't convert the string to an int.
	}
	return v
}

func flattenBigqueryDataTransferConfigDisabled(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenBigqueryDataTransferConfigParams(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandBigqueryDataTransferConfigDisplayName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandBigqueryDataTransferConfigDestinationDatasetId(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandBigqueryDataTransferConfigDataSourceId(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandBigqueryDataTransferConfigSchedule(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandBigqueryDataTransferConfigDataRefreshWindowDays(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandBigqueryDataTransferConfigDisabled(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandBigqueryDataTransferConfigParams(v interface{}, d TerraformResourceData, config *Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}
