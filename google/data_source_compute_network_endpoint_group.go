package google

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceGoogleComputeNetworkEndpointGroup() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := datasourceSchemaFromResourceSchema(resourceComputeNetworkEndpointGroup().Schema)

	// Set 'Optional' schema elements
	addOptionalFieldsToSchema(dsSchema, "name")
	addOptionalFieldsToSchema(dsSchema, "zone")
	addOptionalFieldsToSchema(dsSchema, "self_link")

	return &schema.Resource{
		Read:   dataSourceComputeNetworkEndpointGroupRead,
		Schema: dsSchema,
	}
}

func dataSourceComputeNetworkEndpointGroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	if name, ok := d.GetOk("name"); ok {
		zone, err := getZone(d, config)
		if err != nil {
			return err
		}
		d.SetId(fmt.Sprintf("%s/%s", zone, name.(string)))
	} else if selfLink, ok := d.GetOk("self_link"); ok {
		parsed, err := ParseNetworkEndpointGroupFieldValue(selfLink.(string), d, config)
		if err != nil {
			return err
		}
		d.Set("name", parsed.Name)
		d.Set("zone", parsed.Zone)
		d.Set("project", parsed.Project)
		d.SetId(fmt.Sprintf("%s/%s", parsed.Zone, parsed.Name))
	} else {
		return errors.New("Must provide either `self_link` or `zone/name`")
	}

	return resourceComputeNetworkEndpointGroupRead(d, meta)
}
