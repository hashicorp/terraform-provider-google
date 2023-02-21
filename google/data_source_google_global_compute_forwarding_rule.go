package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleComputeGlobalForwardingRule() *schema.Resource {
	dsSchema := datasourceSchemaFromResourceSchema(ResourceComputeGlobalForwardingRule().Schema)

	// Set 'Required' schema elements
	addRequiredFieldsToSchema(dsSchema, "name")

	// Set 'Optional' schema elements
	addOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceGoogleComputeGlobalForwardingRuleRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleComputeGlobalForwardingRuleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	name := d.Get("name").(string)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("projects/%s/global/forwardingRules/%s", project, name))

	return resourceComputeGlobalForwardingRuleRead(d, meta)
}
