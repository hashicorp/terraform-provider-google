package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceGoogleProjectOrganizationPolicy() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := datasourceSchemaFromResourceSchema(resourceGoogleProjectOrganizationPolicy().Schema)

	addRequiredFieldsToSchema(dsSchema, "project")
	addRequiredFieldsToSchema(dsSchema, "constraint")

	return &schema.Resource{
		Read:   datasourceGoogleProjectOrganizationPolicyRead,
		Schema: dsSchema,
	}
}

func datasourceGoogleProjectOrganizationPolicyRead(d *schema.ResourceData, meta interface{}) error {

	d.SetId(fmt.Sprintf("%s:%s", d.Get("project"), d.Get("constraint")))

	return resourceGoogleProjectOrganizationPolicyRead(d, meta)
}
