package google

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func DataSourceGoogleComputeHealthCheck() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := datasourceSchemaFromResourceSchema(ResourceComputeHealthCheck().Schema)

	// Set 'Required' schema elements
	addRequiredFieldsToSchema(dsSchema, "name")

	// Set 'Optional' schema elements
	addOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceGoogleComputeHealthCheckRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleComputeHealthCheckRead(d *schema.ResourceData, meta interface{}) error {
	id, err := replaceVars(d, meta.(*Config), "projects/{{project}}/global/healthChecks/{{name}}")
	if err != nil {
		return err
	}
	d.SetId(id)

	return resourceComputeHealthCheckRead(d, meta)
}
