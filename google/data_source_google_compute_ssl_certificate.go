package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceGoogleComputeSslCertificate() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := datasourceSchemaFromResourceSchema(resourceComputeSslCertificate().Schema)

	// Set 'Required' schema elements
	addRequiredFieldsToSchema(dsSchema, "name")

	// Set 'Optional' schema elements
	addOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceComputeSslCertificateRead,
		Schema: dsSchema,
	}
}

func dataSourceComputeSslCertificateRead(d *schema.ResourceData, meta interface{}) error {
	certificateName := d.Get("name").(string)

	d.SetId(certificateName)

	return resourceComputeSslCertificateRead(d, meta)
}
