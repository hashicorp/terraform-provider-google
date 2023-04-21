package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleProjectService() *schema.Resource {

	dsSchema := datasourceSchemaFromResourceSchema(ResourceGoogleProjectService().Schema)
	addRequiredFieldsToSchema(dsSchema, "service")
	addOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceGoogleProjectServiceRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleProjectServiceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	id, err := ReplaceVars(d, config, "{{project}}/{{service}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)
	return resourceGoogleProjectServiceRead(d, meta)
}
