package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleComposerEnvironment() *schema.Resource {
	dsSchema := datasourceSchemaFromResourceSchema(ResourceComposerEnvironment().Schema)

	// Set 'Required' schema elements
	addRequiredFieldsToSchema(dsSchema, "name")

	// Set 'Optional' schema elements
	addOptionalFieldsToSchema(dsSchema, "project", "region")

	return &schema.Resource{
		Read:   dataSourceGoogleComposerEnvironmentRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleComposerEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	region, err := getRegion(d, config)
	if err != nil {
		return err
	}
	envName := d.Get("name").(string)

	d.SetId(fmt.Sprintf("projects/%s/locations/%s/environments/%s", project, region, envName))

	return resourceComposerEnvironmentRead(d, meta)
}
