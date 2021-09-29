package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleStorageBucket() *schema.Resource {

	dsSchema := datasourceSchemaFromResourceSchema(resourceStorageBucket().Schema)

	addRequiredFieldsToSchema(dsSchema, "name")

	return &schema.Resource{
		Read:   dataSourceGoogleStorageBucketRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleStorageBucketRead(d *schema.ResourceData, meta interface{}) error {

	bucket := d.Get("name").(string)
	d.SetId(bucket)

	return resourceStorageBucketRead(d, meta)
}
