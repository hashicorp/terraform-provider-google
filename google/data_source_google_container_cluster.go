package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceGoogleContainerCluster() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := datasourceSchemaFromResourceSchema(resourceContainerCluster().Schema)

	// Set 'Required' schema elements
	addRequiredFieldsToSchema(dsSchema, "name")

	// Set 'Optional' schema elements
	addOptionalFieldsToSchema(dsSchema, "project", "zone", "region", "location")

	return &schema.Resource{
		Read:   datasourceContainerClusterRead,
		Schema: dsSchema,
	}
}

func datasourceContainerClusterRead(d *schema.ResourceData, meta interface{}) error {
	clusterName := d.Get("name").(string)

	d.SetId(clusterName)

	return resourceContainerClusterRead(d, meta)
}
