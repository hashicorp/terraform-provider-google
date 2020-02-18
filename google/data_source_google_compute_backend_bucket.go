package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceGoogleComputeBackendBucket() *schema.Resource {
	dsSchema := datasourceSchemaFromResourceSchema(resourceComputeBackendBucket().Schema)

	// Set 'Required' schema elements
	addRequiredFieldsToSchema(dsSchema, "name")

	// Set 'Optional' schema elements
	addOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceComputeBackendBucketRead,
		Schema: dsSchema,
	}
}

func dataSourceComputeBackendBucketRead(d *schema.ResourceData, meta interface{}) error {
	serviceName := d.Get("name").(string)

	d.SetId(serviceName)

	return resourceComputeBackendBucketRead(d, meta)
}
