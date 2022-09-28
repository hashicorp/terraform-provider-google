package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleCloudFunctions2Function() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := datasourceSchemaFromResourceSchema(resourceCloudfunctions2function().Schema)

	// Set 'Required' schema elements
	addRequiredFieldsToSchema(dsSchema, "name", "location")

	// Set 'Optional' schema elements
	addOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceGoogleCloudFunctions2FunctionRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleCloudFunctions2FunctionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("projects/%s/locations/%s/functions/%s", project, d.Get("location").(string), d.Get("name").(string)))

	err = resourceCloudfunctions2functionRead(d, meta)
	if err != nil {
		return err
	}

	return nil
}
