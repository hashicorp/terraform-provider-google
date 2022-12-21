package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleBeyondcorpAppGateway() *schema.Resource {

	dsSchema := datasourceSchemaFromResourceSchema(resourceBeyondcorpAppGateway().Schema)

	addRequiredFieldsToSchema(dsSchema, "name")

	addOptionalFieldsToSchema(dsSchema, "project")
	addOptionalFieldsToSchema(dsSchema, "region")

	return &schema.Resource{
		Read:   dataSourceGoogleBeyondcorpAppGatewayRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleBeyondcorpAppGatewayRead(d *schema.ResourceData, meta interface{}) error {
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

	d.SetId(fmt.Sprintf("projects/%s/locations/%s/appGateways/%s", project, region, name))

	return resourceBeyondcorpAppGatewayRead(d, meta)
}
