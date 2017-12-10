package google

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGoogleContainerCluster() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := datasourceSchemaFromResourceSchema(resourceContainerCluster().Schema)

	// Fixup the schema flags of the Required and Optional input attributes
	fixDatasourceSchemaFlags(dsSchema, true, "name", "zone")
	fixDatasourceSchemaFlags(dsSchema, false, "project")

	dsSchema["project"].Optional = true

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
