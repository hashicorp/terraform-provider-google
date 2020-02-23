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
	config := meta.(*Config)

	clusterName := d.Get("name").(string)

	location, err := getLocation(d, config)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	d.SetId(containerClusterFullName(project, location, clusterName))

	return resourceContainerClusterRead(d, meta)
}
