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

package google

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	monitoring "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/monitoring"
)

func resourceMonitoringMetricsScope() *schema.Resource {
	return &schema.Resource{
		Create: resourceMonitoringMetricsScopeCreate,
		Read:   resourceMonitoringMetricsScopeRead,
		Delete: resourceMonitoringMetricsScopeDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMonitoringMetricsScopeImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Immutable. The resource name of the Monitoring Metrics Scope. On input, the resource name can be specified with the scoping project ID or number. On output, the resource name is specified with the scoping project number. Example: `locations/global/metricsScopes/{SCOPING_PROJECT_ID_OR_NUMBER}`",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The time when this `Metrics Scope` was created.",
			},

			"monitored_projects": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Output only. The list of projects monitored by this `Metrics Scope`.",
				Elem:        MonitoringMetricsScopeMonitoredProjectsSchema(),
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The time when this `Metrics Scope` record was last updated.",
			},
		},
	}
}

func MonitoringMetricsScopeMonitoredProjectsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The time when this `MonitoredProject` was created.",
			},

			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Immutable. The resource name of the `MonitoredProject`. On input, the resource name includes the scoping project ID and monitored project ID. On output, it contains the equivalent project numbers. Example: `locations/global/metricsScopes/{SCOPING_PROJECT_ID_OR_NUMBER}/projects/{MONITORED_PROJECT_ID_OR_NUMBER}`",
			},
		},
	}
}

func resourceMonitoringMetricsScopeCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	obj := &monitoring.MetricsScope{
		Name: dcl.String(d.Get("name").(string)),
	}

	id, err := replaceVarsForId(d, config, "locations/global/metricsScopes/{{name}}")
	if err != nil {
		return fmt.Errorf("error constructing id: %s", err)
	}
	d.SetId(id)
	createDirective := CreateDirective
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}
	billingProject := ""
	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := NewDCLMonitoringClient(config, userAgent, billingProject)
	client.Config.BasePath += "v1"
	res, err := client.ApplyMetricsScope(context.Background(), obj, createDirective...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating MetricsScope: %s", err)
	}

	log.Printf("[DEBUG] Finished creating MetricsScope %q: %#v", d.Id(), res)

	return resourceMonitoringMetricsScopeRead(d, meta)
}

func resourceMonitoringMetricsScopeRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	obj := &monitoring.MetricsScope{
		Name: dcl.String(d.Get("name").(string)),
	}

	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}
	billingProject := ""
	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := NewDCLMonitoringClient(config, userAgent, billingProject)
	client.Config.BasePath += "v1"
	res, err := client.GetMetricsScope(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("MonitoringMetricsScope %q", d.Id())
		return handleNotFoundDCLError(err, d, resourceName)
	}

	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}
	if err = d.Set("create_time", res.CreateTime); err != nil {
		return fmt.Errorf("error setting create_time in state: %s", err)
	}
	if err = d.Set("monitored_projects", flattenMonitoringMetricsScopeMonitoredProjectsArray(res.MonitoredProjects)); err != nil {
		return fmt.Errorf("error setting monitored_projects in state: %s", err)
	}
	if err = d.Set("update_time", res.UpdateTime); err != nil {
		return fmt.Errorf("error setting update_time in state: %s", err)
	}

	return nil
}

func resourceMonitoringMetricsScopeDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf(`[WARNING] MonitoringMetricsScope resources cannot be deleted from Google Cloud.
The resource %s will be removed from Terraform state, but will still be present on Google Cloud.`, d.Id())
	return nil
}

func resourceMonitoringMetricsScopeImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"locations/global/metricsScopes/(?P<name>[^/]+)",
		"(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVarsForId(d, config, "locations/global/metricsScopes/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenMonitoringMetricsScopeMonitoredProjectsArray(objs []monitoring.MetricsScopeMonitoredProjects) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenMonitoringMetricsScopeMonitoredProjects(&item)
		items = append(items, i)
	}

	return items
}

func flattenMonitoringMetricsScopeMonitoredProjects(obj *monitoring.MetricsScopeMonitoredProjects) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"create_time": obj.CreateTime,
		"name":        obj.Name,
	}

	return transformed

}
