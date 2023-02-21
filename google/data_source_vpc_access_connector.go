package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceVPCAccessConnector() *schema.Resource {

	dsSchema := datasourceSchemaFromResourceSchema(ResourceVPCAccessConnector().Schema)
	addRequiredFieldsToSchema(dsSchema, "name")
	addOptionalFieldsToSchema(dsSchema, "project", "region")

	return &schema.Resource{
		Read:   dataSourceVPCAccessConnectorRead,
		Schema: dsSchema,
	}
}

func dataSourceVPCAccessConnectorRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	id, err := replaceVars(d, config, "projects/{{project}}/locations/{{region}}/connectors/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}

	d.SetId(id)

	return resourceVPCAccessConnectorRead(d, meta)
}
