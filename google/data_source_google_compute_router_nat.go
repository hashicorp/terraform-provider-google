package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleComputeRouterNat() *schema.Resource {

	dsSchema := datasourceSchemaFromResourceSchema(ResourceComputeRouterNat().Schema)

	addRequiredFieldsToSchema(dsSchema, "name", "router")
	addOptionalFieldsToSchema(dsSchema, "project", "region")

	return &schema.Resource{
		Read:   dataSourceGoogleComputeRouterNatRead,
		Schema: dsSchema,
	}

}

func dataSourceGoogleComputeRouterNatRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	id, err := replaceVars(d, config, "{{project}}/{{region}}/{{router}}/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return resourceComputeRouterNatRead(d, meta)
}
