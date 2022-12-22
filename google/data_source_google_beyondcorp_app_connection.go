package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleBeyondcorpAppConnection() *schema.Resource {

	dsSchema := datasourceSchemaFromResourceSchema(resourceBeyondcorpAppConnection().Schema)

	addRequiredFieldsToSchema(dsSchema, "name")

	addOptionalFieldsToSchema(dsSchema, "project")
	addOptionalFieldsToSchema(dsSchema, "region")

	return &schema.Resource{
		Read:   dataSourceGoogleBeyondcorpAppConnectionRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleBeyondcorpAppConnectionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	name := d.Get("name").(string)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("projects/%s/locations/%s/appConnections/%s", project, region, name))

	return resourceBeyondcorpAppConnectionRead(d, meta)
}
