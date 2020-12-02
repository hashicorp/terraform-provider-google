package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleComputeResourcePolicy() *schema.Resource {
	dsSchema := datasourceSchemaFromResourceSchema(resourceComputeResourcePolicy().Schema)

	addRequiredFieldsToSchema(dsSchema, "name")
	addOptionalFieldsToSchema(dsSchema, "region")
	addOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceGoogleComputeResourcePolicyRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleComputeResourcePolicyRead(d *schema.ResourceData, meta interface{}) error {
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

	d.SetId(fmt.Sprintf("projects/%s/regions/%s/resourcePolicies/%s", project, region, name))

	return resourceComputeResourcePolicyRead(d, meta)
}
