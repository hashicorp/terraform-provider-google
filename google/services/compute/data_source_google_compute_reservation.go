// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleComputeReservation() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceComputeReservation().Schema)

	// Set 'Required' schema elements
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "name")
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "zone")

	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceGoogleComputeReservationRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleComputeReservationRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	err := resourceComputeReservationRead(d, meta)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	zone, err := tpgresource.GetZone(d, config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)

	d.SetId(fmt.Sprintf("projects/%s/zones/%s/reservations/%s", project, zone, name))
	return nil
}
