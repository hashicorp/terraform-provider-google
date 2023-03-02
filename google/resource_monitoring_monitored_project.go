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

func ResourceMonitoringMonitoredProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceMonitoringMonitoredProjectCreate,
		Read:   resourceMonitoringMonitoredProjectRead,
		Delete: resourceMonitoringMonitoredProjectDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMonitoringMonitoredProjectImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"metrics_scope": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Required. The resource name of the existing Metrics Scope that will monitor this project. Example: locations/global/metricsScopes/{SCOPING_PROJECT_ID_OR_NUMBER}",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Immutable. The resource name of the `MonitoredProject`. On input, the resource name includes the scoping project ID and monitored project ID. On output, it contains the equivalent project numbers. Example: `locations/global/metricsScopes/{SCOPING_PROJECT_ID_OR_NUMBER}/projects/{MONITORED_PROJECT_ID_OR_NUMBER}`",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The time when this `MonitoredProject` was created.",
			},
		},
	}
}

func resourceMonitoringMonitoredProjectCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	obj := &monitoring.MonitoredProject{
		MetricsScope: dcl.String(d.Get("metrics_scope").(string)),
		Name:         dcl.String(d.Get("name").(string)),
	}

	id, err := obj.ID()
	if err != nil {
		return fmt.Errorf("error constructing id: %s", err)
	}
	d.SetId(id)
	directive := CreateDirective
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := ""
	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := NewDCLMonitoringClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutCreate))
	client.Config.BasePath += "v1"
	if bp, err := replaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyMonitoredProject(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating MonitoredProject: %s", err)
	}

	log.Printf("[DEBUG] Finished creating MonitoredProject %q: %#v", d.Id(), res)

	return resourceMonitoringMonitoredProjectRead(d, meta)
}

func resourceMonitoringMonitoredProjectRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	obj := &monitoring.MonitoredProject{
		MetricsScope: dcl.String(d.Get("metrics_scope").(string)),
		Name:         dcl.String(d.Get("name").(string)),
	}

	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := ""
	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := NewDCLMonitoringClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutRead))
	client.Config.BasePath += "v1"
	if bp, err := replaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.GetMonitoredProject(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("MonitoringMonitoredProject %q", d.Id())
		return handleNotFoundDCLError(err, d, resourceName)
	}

	if err = d.Set("metrics_scope", res.MetricsScope); err != nil {
		return fmt.Errorf("error setting metrics_scope in state: %s", err)
	}
	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}
	if err = d.Set("create_time", res.CreateTime); err != nil {
		return fmt.Errorf("error setting create_time in state: %s", err)
	}

	return nil
}

func resourceMonitoringMonitoredProjectDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	obj := &monitoring.MonitoredProject{
		MetricsScope: dcl.String(d.Get("metrics_scope").(string)),
		Name:         dcl.String(d.Get("name").(string)),
	}

	log.Printf("[DEBUG] Deleting MonitoredProject %q", d.Id())
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := ""
	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := NewDCLMonitoringClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutDelete))
	client.Config.BasePath += "v1"
	if bp, err := replaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	if err := client.DeleteMonitoredProject(context.Background(), obj); err != nil {
		return fmt.Errorf("Error deleting MonitoredProject: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting MonitoredProject %q", d.Id())
	return nil
}

func resourceMonitoringMonitoredProjectImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	if err := parseImportId([]string{
		"locations/global/metricsScopes/(?P<metrics_scope>[^/]+)/projects/(?P<name>[^/]+)",
		"(?P<metrics_scope>[^/]+)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVarsForId(d, config, "locations/global/metricsScopes/{{metrics_scope}}/projects/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
