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

func DataSourceGoogleComputeRegionInstanceGroupManager() *schema.Resource {

	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceComputeRegionInstanceGroupManager().Schema)
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "name", "self_link", "project", "region")

	return &schema.Resource{
		Read:   dataSourceComputeRegionInstanceGroupManagerRead,
		Schema: dsSchema,
	}
}

func dataSourceComputeRegionInstanceGroupManagerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	if selfLink, ok := d.Get("self_link").(string); ok && selfLink != "" {
		parsed, err := tpgresource.ParseRegionalInstanceGroupManagersFieldValue(selfLink, d, config)
		if err != nil {
			return fmt.Errorf("InstanceGroup name, region or project could not be parsed from %s: %v", selfLink, err)
		}
		if err := d.Set("name", parsed.Name); err != nil {
			return fmt.Errorf("Error setting name: %s", err)
		}
		if err := d.Set("region", parsed.Region); err != nil {
			return fmt.Errorf("Error setting region: %s", err)
		}
		if err := d.Set("project", parsed.Project); err != nil {
			return fmt.Errorf("Error setting project: %s", err)
		}
		d.SetId(fmt.Sprintf("projects/%s/regions/%s/instanceGroupManagers/%s", parsed.Project, parsed.Region, parsed.Name))
	} else if name, ok := d.Get("name").(string); ok && name != "" {
		region, err := tpgresource.GetRegion(d, config)
		if err != nil {
			return err
		}
		project, err := tpgresource.GetProject(d, config)
		if err != nil {
			return err
		}
		d.SetId(fmt.Sprintf("projects/%s/regions/%s/instanceGroupManagers/%s", project, region, name))
	} else {
		return errors.New("Must provide either `self_link` or `region/name`")
	}

	err := resourceComputeRegionInstanceGroupManagerRead(d, meta)

	if err != nil {
		return err
	}
	if d.Id() == "" {
		return errors.New("Regional Instance Manager Group not found")
	}
	return nil
}
