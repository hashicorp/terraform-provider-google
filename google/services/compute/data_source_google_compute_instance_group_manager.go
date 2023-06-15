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

func DataSourceGoogleComputeInstanceGroupManager() *schema.Resource {

	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceComputeInstanceGroupManager().Schema)
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "name", "self_link", "project", "zone")

	return &schema.Resource{
		Read:   dataSourceComputeInstanceGroupManagerRead,
		Schema: dsSchema,
	}
}

func dataSourceComputeInstanceGroupManagerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	if selfLink, ok := d.GetOk("self_link"); ok {
		parsed, err := tpgresource.ParseInstanceGroupFieldValue(selfLink.(string), d, config)
		if err != nil {
			return fmt.Errorf("InstanceGroup name, zone or project could not be parsed from %s", selfLink)
		}
		if err := d.Set("name", parsed.Name); err != nil {
			return fmt.Errorf("Error setting name: %s", err)
		}
		if err := d.Set("zone", parsed.Zone); err != nil {
			return fmt.Errorf("Error setting zone: %s", err)
		}
		if err := d.Set("project", parsed.Project); err != nil {
			return fmt.Errorf("Error setting project: %s", err)
		}
		d.SetId(fmt.Sprintf("projects/%s/zones/%s/instanceGroupManagers/%s", parsed.Project, parsed.Zone, parsed.Name))
	} else if name, ok := d.GetOk("name"); ok {
		zone, err := tpgresource.GetZone(d, config)
		if err != nil {
			return err
		}
		project, err := tpgresource.GetProject(d, config)
		if err != nil {
			return err
		}
		d.SetId(fmt.Sprintf("projects/%s/zones/%s/instanceGroupManagers/%s", project, zone, name.(string)))
	} else {
		return errors.New("Must provide either `self_link` or `zone/name`")
	}

	err := resourceComputeInstanceGroupManagerRead(d, meta)

	if err != nil {
		return err
	}
	if d.Id() == "" {
		return errors.New("Instance Manager Group not found")
	}
	return nil
}
