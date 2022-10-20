package google

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleComputeRegionNetworkEndpointGroup() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := datasourceSchemaFromResourceSchema(resourceComputeRegionNetworkEndpointGroup().Schema)

	addOptionalFieldsToSchema(dsSchema, "name")
	addOptionalFieldsToSchema(dsSchema, "region")
	addOptionalFieldsToSchema(dsSchema, "project")
	addOptionalFieldsToSchema(dsSchema, "self_link")

	return &schema.Resource{
		Read:   dataSourceComputeRegionNetworkEndpointGroupRead,
		Schema: dsSchema,
	}
}

func dataSourceComputeRegionNetworkEndpointGroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	if name, ok := d.GetOk("name"); ok {
		project, err := getProject(d, config)
		if err != nil {
			return err
		}
		region, err := getRegion(d, config)
		if err != nil {
			return err
		}

		d.SetId(fmt.Sprintf("projects/%s/regions/%s/networkEndpointGroups/%s", project, region, name.(string)))
	} else if selfLink, ok := d.GetOk("self_link"); ok {
		parsed, err := ParseNetworkEndpointGroupRegionalFieldValue(selfLink.(string), d, config)
		if err != nil {
			return err
		}
		if err := d.Set("name", parsed.Name); err != nil {
			return fmt.Errorf("Error setting name: %s", err)
		}
		if err := d.Set("project", parsed.Project); err != nil {
			return fmt.Errorf("Error setting project: %s", err)
		}
		if err := d.Set("region", parsed.Region); err != nil {
			return fmt.Errorf("Error setting region: %s", err)
		}

		d.SetId(fmt.Sprintf("projects/%s/regions/%s/networkEndpointGroups/%s", parsed.Project, parsed.Region, parsed.Name))
	} else {
		return errors.New("Must provide either `self_link` or `region/name`")
	}

	return resourceComputeRegionNetworkEndpointGroupRead(d, meta)
}
