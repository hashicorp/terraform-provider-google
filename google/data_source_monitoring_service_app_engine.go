// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceMonitoringServiceAppEngine() *schema.Resource {
	aeSchema := map[string]*schema.Schema{
		"module_id": {
			Type:     schema.TypeString,
			Required: true,
			Description: `The ID of the App Engine module underlying this service. 
Corresponds to the 'moduleId' resource label for a 'gae_app'
monitored resource(see https://cloud.google.com/monitoring/api/resources#tag_gae_app)`,
		},
	}
	filter := `app_engine.module_id="{{module_id}}"`
	return dataSourceMonitoringServiceType(aeSchema, filter, dataSourceMonitoringServiceAppEngineRead)
}

func dataSourceMonitoringServiceAppEngineRead(res map[string]interface{}, d *schema.ResourceData, meta interface{}) error {
	var appEngine map[string]interface{}
	if v, ok := res["app_engine"]; ok {
		appEngine = v.(map[string]interface{})
	}
	if len(appEngine) == 0 {
		return nil
	}

	if err := d.Set("module_id", appEngine["module_id"]); err != nil {
		return err
	}
	return nil
}
