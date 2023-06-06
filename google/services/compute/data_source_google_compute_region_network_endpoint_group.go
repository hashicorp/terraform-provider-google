// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleComputeRegionNetworkEndpointGroup() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceComputeRegionNetworkEndpointGroup().Schema)

	tpgresource.AddOptionalFieldsToSchema(dsSchema, "name")
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "region")
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project")
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "self_link")

	return &schema.Resource{
		Read:   dataSourceComputeRegionNetworkEndpointGroupRead,
		Schema: dsSchema,
	}
}

func dataSourceComputeRegionNetworkEndpointGroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	if name, ok := d.GetOk("name"); ok {
		project, err := tpgresource.GetProject(d, config)
		if err != nil {
			return err
		}
		region, err := tpgresource.GetRegion(d, config)
		if err != nil {
			return err
		}

		d.SetId(fmt.Sprintf("projects/%s/regions/%s/networkEndpointGroups/%s", project, region, name.(string)))
	} else if selfLink, ok := d.GetOk("self_link"); ok {
		parsed, err := tpgresource.ParseNetworkEndpointGroupRegionalFieldValue(selfLink.(string), d, config)
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
