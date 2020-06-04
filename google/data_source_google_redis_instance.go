package google

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func dataSourceGoogleRedisInstance() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := datasourceSchemaFromResourceSchema(resourceRedisInstance().Schema)

	// Set 'Required' schema elements
	addRequiredFieldsToSchema(dsSchema, "name")

	// Set 'Optional' schema elements
	addOptionalFieldsToSchema(dsSchema, "project", "region")

	return &schema.Resource{
		Read:   dataSourceGoogleRedisInstanceRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleRedisInstanceRead(d *schema.ResourceData, meta interface{}) error {
	instanceName := d.Get("name").(string)

	d.SetId(instanceName)

	return resourceRedisInstanceRead(d, meta)
}
