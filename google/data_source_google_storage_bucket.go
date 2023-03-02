package google

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleStorageBucket() *schema.Resource {

	dsSchema := datasourceSchemaFromResourceSchema(ResourceStorageBucket().Schema)

	addRequiredFieldsToSchema(dsSchema, "name")

	return &schema.Resource{
		Read:   dataSourceGoogleStorageBucketRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleStorageBucketRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	// Get the bucket and acl
	bucket := d.Get("name").(string)

	res, err := config.NewStorageClient(userAgent).Buckets.Get(bucket).Do()
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Read bucket %v at location %v\n\n", res.Name, res.SelfLink)

	return setStorageBucket(d, config, res, bucket, userAgent)
}
