package google

import (
	"fmt"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func monitoringDashboardDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	computedFields := []string{"etag", "name"}

	oldMap, err := structure.ExpandJsonFromString(old)
	if err != nil {
		return false
	}

	newMap, err := structure.ExpandJsonFromString(new)
	if err != nil {
		return false
	}

	for _, f := range computedFields {
		delete(oldMap, f)
		delete(newMap, f)
	}

	return reflect.DeepEqual(oldMap, newMap)
}

func resourceMonitoringDashboard() *schema.Resource {
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

		Schema: map[string]*schema.Schema{
			"dashboard_json": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.ValidateJsonString,
				DiffSuppressFunc: monitoringDashboardDiffSuppress,
				StateFunc: func(v interface{}) string {
					json, _ := structure.NormalizeJsonString(v)
					return json
				},
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

func resourceMonitoringDashboardCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	obj, err := structure.ExpandJsonFromString(d.Get("dashboard_json").(string))
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{MonitoringBasePath}}v1/projects/{{project}}/dashboards")
	if err != nil {
		return err
	}
	res, err := sendRequestWithTimeout(config, "POST", project, url, obj, d.Timeout(schema.TimeoutCreate), isMonitoringConcurrentEditError)
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
	config := meta.(*Config)

	url := config.MonitoringBasePath + "v1/" + d.Id()

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	res, err := sendRequest(config, "GET", project, url, nil, isMonitoringConcurrentEditError)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("MonitoringDashboard %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading Dashboard: %s", err)
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
	config := meta.(*Config)

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

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	url := config.MonitoringBasePath + "v1/" + d.Id()
	_, err = sendRequestWithTimeout(config, "PATCH", project, url, nObj, d.Timeout(schema.TimeoutUpdate), isMonitoringConcurrentEditError)
	if err != nil {
		return fmt.Errorf("Error updating Dashboard %q: %s", d.Id(), err)
	}

	return resourceMonitoringDashboardRead(d, config)
}

func resourceMonitoringDashboardDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	url := config.MonitoringBasePath + "v1/" + d.Id()

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	_, err = sendRequestWithTimeout(config, "DELETE", project, url, nil, d.Timeout(schema.TimeoutDelete), isMonitoringConcurrentEditError)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("MonitoringDashboard %q", d.Id()))
	}

	return nil
}

func resourceMonitoringDashboardImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	// current import_formats can't import fields with forward slashes in their value
	parts, err := getImportIdQualifiers([]string{"projects/(?P<project>[^/]+)/dashboards/(?P<id>[^/]+)", "(?P<id>[^/]+)"}, d, config, d.Id())
	if err != nil {
		return nil, err
	}

	d.Set("project", parts["project"])
	d.SetId(fmt.Sprintf("projects/%s/dashboards/%s", parts["project"], parts["id"]))

	return []*schema.ResourceData{d}, nil
}
