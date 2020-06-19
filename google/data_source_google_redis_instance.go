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
	id, err := replaceVars(d, meta.(*Config), "projects/{{project}}/locations/{{region}}/instances/{{name}}")
	if err != nil {
		return err
	}
	d.SetId(id)

	return resourceRedisInstanceRead(d, meta)
}
