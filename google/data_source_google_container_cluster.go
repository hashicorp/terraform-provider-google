package google

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGoogleContainerCluster() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := datasourceSchemaFromResourceSchema(resourceContainerCluster().Schema)

	// Set 'Required' schema elements
	addRequiredFieldsToSchema(dsSchema, "name", "zone")

	// Set 'Optional' schema elements
	addOptionalFieldsToSchema(dsSchema, "project")

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
