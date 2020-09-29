package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSqlDatabaseInstance() *schema.Resource {

	dsSchema := datasourceSchemaFromResourceSchema(resourceSqlDatabaseInstance().Schema)
	addRequiredFieldsToSchema(dsSchema, "name")
	addOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceSqlDatabaseInstanceRead,
		Schema: dsSchema,
	}
}

func dataSourceSqlDatabaseInstanceRead(d *schema.ResourceData, meta interface{}) error {

	return resourceSqlDatabaseInstanceRead(d, meta)

}
