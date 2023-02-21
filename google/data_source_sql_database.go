package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSqlDatabase() *schema.Resource {

	dsSchema := datasourceSchemaFromResourceSchema(ResourceSQLDatabase().Schema)
	addRequiredFieldsToSchema(dsSchema, "name")
	addRequiredFieldsToSchema(dsSchema, "instance")
	addOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceSqlDatabaseRead,
		Schema: dsSchema,
	}
}

func dataSourceSqlDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Database: %s", err)
	}
	d.SetId(fmt.Sprintf("projects/%s/instances/%s/databases/%s", project, d.Get("instance").(string), d.Get("name").(string)))
	err = resourceSQLDatabaseRead(d, meta)
	if err != nil {
		return err
	}
	if err := d.Set("deletion_policy", nil); err != nil {
		return fmt.Errorf("Error setting deletion_policy: %s", err)
	}
	return nil
}
