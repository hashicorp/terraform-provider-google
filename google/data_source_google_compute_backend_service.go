package google

import (
	"fmt"

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
	config := meta.(*Config)

	serviceName := d.Get("name").(string)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("projects/%s/global/backendServices/%s", project, serviceName))

	return resourceComputeBackendServiceRead(d, meta)
}
