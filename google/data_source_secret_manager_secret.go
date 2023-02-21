package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSecretManagerSecret() *schema.Resource {

	dsSchema := datasourceSchemaFromResourceSchema(ResourceSecretManagerSecret().Schema)
	addRequiredFieldsToSchema(dsSchema, "secret_id")
	addOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceSecretManagerSecretRead,
		Schema: dsSchema,
	}
}

func dataSourceSecretManagerSecretRead(d *schema.ResourceData, meta interface{}) error {
	id, err := replaceVars(d, meta.(*Config), "projects/{{project}}/secrets/{{secret_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)
	return resourceSecretManagerSecretRead(d, meta)
}
