package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleComputeForwardingRule() *schema.Resource {
	dsSchema := datasourceSchemaFromResourceSchema(ResourceComputeForwardingRule().Schema)

	// Set 'Required' schema elements
	addRequiredFieldsToSchema(dsSchema, "name")

	// Set 'Optional' schema elements
	addOptionalFieldsToSchema(dsSchema, "project")
	addOptionalFieldsToSchema(dsSchema, "region")

	return &schema.Resource{
		Read:   dataSourceGoogleComputeForwardingRuleRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleComputeForwardingRuleRead(d *schema.ResourceData, meta interface{}) error {
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

	d.SetId(fmt.Sprintf("projects/%s/regions/%s/forwardingRules/%s", project, region, name))

	return resourceComputeForwardingRuleRead(d, meta)
}
