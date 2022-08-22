package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleIapClient() *schema.Resource {

	dsSchema := datasourceSchemaFromResourceSchema(resourceIapClient().Schema)
	addRequiredFieldsToSchema(dsSchema, "brand", "client_id")

	return &schema.Resource{
		Read:   dataSourceGoogleIapClientRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleIapClientRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	id, err := replaceVars(d, config, "{{brand}}/identityAwareProxyClients/{{client_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)
	return resourceIapClientRead(d, meta)
}
