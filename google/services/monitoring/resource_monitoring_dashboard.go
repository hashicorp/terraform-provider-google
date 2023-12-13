// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package monitoring

import (
	"fmt"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// This recursive function takes an old map and a new map and is intended to remove the computed keys
// from the old json string (stored in state) so that it doesn't show a diff if it's not defined in the
// new map's json string (defined in config)
func removeComputedKeys(old map[string]interface{}, new map[string]interface{}) map[string]interface{} {
	for k, v := range old {
		if _, ok := old[k]; ok && new[k] == nil {
			delete(old, k)
			continue
		}

		if reflect.ValueOf(v).Kind() == reflect.Map {
			old[k] = removeComputedKeys(v.(map[string]interface{}), new[k].(map[string]interface{}))
			continue
		}

		if reflect.ValueOf(v).Kind() == reflect.Slice {
			for i, j := range v.([]interface{}) {
				if reflect.ValueOf(j).Kind() == reflect.Map && len(new[k].([]interface{})) > i {
					old[k].([]interface{})[i] = removeComputedKeys(j.(map[string]interface{}), new[k].([]interface{})[i].(map[string]interface{}))
				}
			}
			continue
		}
	}

	return old
}

func monitoringDashboardDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	oldMap, err := structure.ExpandJsonFromString(old)
	if err != nil {
		return false
	}
	newMap, err := structure.ExpandJsonFromString(new)
	if err != nil {
		return false
	}

	oldMap = removeComputedKeys(oldMap, newMap)
	return reflect.DeepEqual(oldMap, newMap)
}

func ResourceMonitoringDashboard() *schema.Resource {
	return &schema.Resource{
		Create: resourceMonitoringDashboardCreate,
		Read:   resourceMonitoringDashboardRead,
		Update: resourceMonitoringDashboardUpdate,
		Delete: resourceMonitoringDashboardDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMonitoringDashboardImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Update: schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			tpgresource.DefaultProviderProject,
		),

		Schema: map[string]*schema.Schema{
			"dashboard_json": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: monitoringDashboardDiffSuppress,
				StateFunc: func(v interface{}) string {
					json, _ := structure.NormalizeJsonString(v)
					return json
				},
				Description: `The JSON representation of a dashboard, following the format at https://cloud.google.com/monitoring/api/ref_v3/rest/v1/projects.dashboards.`,
			},
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The ID of the project in which the resource belongs. If it is not provided, the provider project is used.`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceMonitoringDashboardCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj, err := structure.ExpandJsonFromString(d.Get("dashboard_json").(string))
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{MonitoringBasePath}}v1/projects/{{project}}/dashboards")
	if err != nil {
		return err
	}
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:               config,
		Method:               "POST",
		Project:              project,
		RawURL:               url,
		UserAgent:            userAgent,
		Body:                 obj,
		Timeout:              d.Timeout(schema.TimeoutCreate),
		ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.IsMonitoringConcurrentEditError},
	})
	if err != nil {
		return fmt.Errorf("Error creating Dashboard: %s", err)
	}

	name, ok := res["name"]
	if !ok {
		return fmt.Errorf("Create response didn't contain critical fields. Create may not have succeeded.")
	}
	d.SetId(name.(string))

	return resourceMonitoringDashboardRead(d, config)
}

func resourceMonitoringDashboardRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url := config.MonitoringBasePath + "v1/" + d.Id()

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:               config,
		Method:               "GET",
		Project:              project,
		RawURL:               url,
		UserAgent:            userAgent,
		ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.IsMonitoringConcurrentEditError},
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("MonitoringDashboard %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting Dashboard: %s", err)
	}

	str, err := structure.FlattenJsonToString(res)
	if err != nil {
		return fmt.Errorf("Error reading Dashboard: %s", err)
	}
	if err = d.Set("dashboard_json", str); err != nil {
		return fmt.Errorf("Error reading Dashboard: %s", err)
	}

	return nil
}

func resourceMonitoringDashboardUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	o, n := d.GetChange("dashboard_json")
	oObj, err := structure.ExpandJsonFromString(o.(string))
	if err != nil {
		return err
	}
	nObj, err := structure.ExpandJsonFromString(n.(string))
	if err != nil {
		return err
	}

	nObj["etag"] = oObj["etag"]

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	url := config.MonitoringBasePath + "v1/" + d.Id()
	_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:               config,
		Method:               "PATCH",
		Project:              project,
		RawURL:               url,
		UserAgent:            userAgent,
		Body:                 nObj,
		Timeout:              d.Timeout(schema.TimeoutUpdate),
		ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.IsMonitoringConcurrentEditError},
	})
	if err != nil {
		return fmt.Errorf("Error updating Dashboard %q: %s", d.Id(), err)
	}

	return resourceMonitoringDashboardRead(d, config)
}

func resourceMonitoringDashboardDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url := config.MonitoringBasePath + "v1/" + d.Id()

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:               config,
		Method:               "DELETE",
		Project:              project,
		RawURL:               url,
		UserAgent:            userAgent,
		Timeout:              d.Timeout(schema.TimeoutDelete),
		ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.IsMonitoringConcurrentEditError},
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("MonitoringDashboard %q", d.Id()))
	}

	return nil
}

func resourceMonitoringDashboardImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)

	// current import_formats can't import fields with forward slashes in their value
	parts, err := tpgresource.GetImportIdQualifiers([]string{"projects/(?P<project>[^/]+)/dashboards/(?P<id>[^/]+)", "(?P<id>[^/]+)"}, d, config, d.Id())
	if err != nil {
		return nil, err
	}

	if err := d.Set("project", parts["project"]); err != nil {
		return nil, fmt.Errorf("Error setting project: %s", err)
	}
	d.SetId(fmt.Sprintf("projects/%s/dashboards/%s", parts["project"], parts["id"]))

	return []*schema.ResourceData{d}, nil
}
