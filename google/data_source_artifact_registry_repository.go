package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceArtifactRegistryRepository() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := datasourceSchemaFromResourceSchema(ResourceArtifactRegistryRepository().Schema)

	// Set 'Required' schema elements
	addRequiredFieldsToSchema(dsSchema, "repository_id", "location")

	// Set 'Optional' schema elements
	addOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceArtifactRegistryRepositoryRead,
		Schema: dsSchema,
	}
}

func dataSourceArtifactRegistryRepositoryRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	location, err := getLocation(d, config)
	if err != nil {
		return err
	}

	repository_id := d.Get("repository_id").(string)
	d.SetId(fmt.Sprintf("projects/%s/locations/%s/repositories/%s", project, location, repository_id))

	err = resourceArtifactRegistryRepositoryRead(d, meta)
	if err != nil {
		return err
	}

	return nil
}
