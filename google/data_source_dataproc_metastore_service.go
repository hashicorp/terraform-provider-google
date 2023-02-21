package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDataprocMetastoreService() *schema.Resource {

	dsSchema := datasourceSchemaFromResourceSchema(ResourceDataprocMetastoreService().Schema)
	addRequiredFieldsToSchema(dsSchema, "service_id")
	addRequiredFieldsToSchema(dsSchema, "location")
	addOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceDataprocMetastoreServiceRead,
		Schema: dsSchema,
	}
}

func dataSourceDataprocMetastoreServiceRead(d *schema.ResourceData, meta interface{}) error {
	id, err := replaceVars(d, meta.(*Config), "projects/{{project}}/locations/{{location}}/services/{{service_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)
	return resourceDataprocMetastoreServiceRead(d, meta)
}
