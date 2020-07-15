package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceSqlDatabaseInstance() *schema.Resource {

	dsSchema := datasourceSchemaFromResourceSchema(resourceSqlDatabaseInstance().Schema)
	addRequiredFieldsToSchema(dsSchema, "name")

	return &schema.Resource{
		Read:   dataSourceSqlDatabaseInstanceRead,
		Schema: dsSchema,
	}
}

func dataSourceSqlDatabaseInstanceRead(d *schema.ResourceData, meta interface{}) error {

	return resourceSqlDatabaseInstanceRead(d, meta)

}
