package google

import (
	"fmt"
	neturl "net/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type monitoringServiceTypeStateSetter func(map[string]interface{}, *schema.ResourceData, interface{}) error

// dataSourceMonitoringServiceType creates a Datasource resource for a type of service. It takes
// - schema for identifying the service, specific to the type (AppEngine moduleId)
// - list query filter to filter a specific service (type, ID) from the list of services for a parent
// - typeFlattenF for reading the service-specific schema (typeSchema)
func dataSourceMonitoringServiceType(
	typeSchema map[string]*schema.Schema,
	listFilter string,
	typeStateSetter monitoringServiceTypeStateSetter) *schema.Resource {

	// Convert monitoring schema to ds schema
	dsSchema := datasourceSchemaFromResourceSchema(resourceMonitoringService().Schema)
	addOptionalFieldsToSchema(dsSchema, "project")

	// Add schema specific to the service type
	dsSchema = mergeSchemas(typeSchema, dsSchema)

	return &schema.Resource{
		Read:   dataSourceMonitoringServiceTypeReadFromList(listFilter, typeStateSetter),
		Schema: dsSchema,
	}
}

// dataSourceMonitoringServiceRead returns a ReadFunc that calls service.list with proper filters
// to identify both the type of service and underlying service resource.
// It takes the list query filter (i.e. ?filter=$listFilter) and a ReadFunc to handle reading any type-specific schema.
func dataSourceMonitoringServiceTypeReadFromList(listFilter string, typeStateSetter monitoringServiceTypeStateSetter) schema.ReadFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)

		project, err := getProject(d, config)
		if err != nil {
			return err
		}

		filters, err := replaceVars(d, config, listFilter)
		if err != nil {
			return err
		}

		listUrlTmpl := "{{MonitoringBasePath}}v3/projects/{{project}}/services?filter=" + neturl.QueryEscape(filters)
		url, err := replaceVars(d, config, listUrlTmpl)
		if err != nil {
			return err
		}

		resp, err := sendRequest(config, "GET", project, url, nil, isMonitoringRetryableError)
		if err != nil {
			return fmt.Errorf("unable to list Monitoring Service for data source: %v", err)
		}

		v, ok := resp["services"]
		if !ok || v == nil {
			return fmt.Errorf("no Monitoring Services found for data source")
		}
		ls, ok := v.([]interface{})
		if !ok {
			return fmt.Errorf("no Monitoring Services found for data source")
		}
		if len(ls) == 0 {
			return fmt.Errorf("no Monitoring Services found for data source")
		}
		if len(ls) > 1 {
			return fmt.Errorf("more than one Monitoring Services with given identifier found")
		}
		res := ls[0].(map[string]interface{})

		if err := d.Set("project", project); err != nil {
			return fmt.Errorf("Error reading Service: %s", err)
		}
		if err := d.Set("display_name", flattenMonitoringServiceDisplayName(res["displayName"], d, config)); err != nil {
			return fmt.Errorf("Error reading Service: %s", err)
		}
		if err := d.Set("telemetry", flattenMonitoringServiceTelemetry(res["telemetry"], d, config)); err != nil {
			return fmt.Errorf("Error reading Service: %s", err)
		}
		if err := d.Set("service_id", flattenMonitoringServiceServiceId(res["name"], d, config)); err != nil {
			return fmt.Errorf("Error reading Service: %s", err)
		}
		if err := typeStateSetter(res, d, config); err != nil {
			return fmt.Errorf("Error reading Service: %s", err)
		}

		name := flattenMonitoringServiceName(res["name"], d, config).(string)
		d.Set("name", name)
		d.SetId(name)

		return nil
	}
}
