package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceGoogleComputeBackendService() *schema.Resource {
	dsSchema := datasourceSchemaFromResourceSchema(resourceComputeBackendService().Schema)

	// Set 'Required' schema elements
	addRequiredFieldsToSchema(dsSchema, "name")

	// Set 'Optional' schema elements
	addOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceComputeBackendServiceRead,
		Schema: dsSchema,
	}
}

func dataSourceComputeBackendServiceRead(d *schema.ResourceData, meta interface{}) error {
	serviceName := d.Get("name").(string)

	d.SetId(serviceName)

	return resourceComputeBackendServiceRead(d, meta)
}
