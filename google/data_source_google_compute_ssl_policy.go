package google

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGoogleComputeSslPolicy() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := datasourceSchemaFromResourceSchema(resourceComputeSslPolicy().Schema)

	// Set 'Required' schema elements
	addRequiredFieldsToSchema(dsSchema, "name")

	// Set 'Optional' schema elements
	addOptionalFieldsToSchema(dsSchema, "custom_features", "description", "min_tls_version", "profile", "project")

	return &schema.Resource{
		Read:   datasourceComputeSslPolicyRead,
		Schema: dsSchema,
	}
}

func datasourceComputeSslPolicyRead(d *schema.ResourceData, meta interface{}) error {
	policyName := d.Get("name").(string)

	d.SetId(policyName)

	return resourceComputeSslPolicyRead(d, meta)
}
