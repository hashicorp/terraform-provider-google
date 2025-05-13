// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package lustre

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceLustreInstance() *schema.Resource {

	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceLustreInstance().Schema)

	dsScema_zone := map[string]*schema.Schema{
		"zone": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: `Zone of Lustre instance`,
		},
	}

	// Set 'Required' schema elements from resource
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "instance_id")

	// Set 'Optional' schema elements from resource
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project")

	// Merge schema elements
	dsSchema_m := tpgresource.MergeSchemas(dsScema_zone, dsSchema)

	return &schema.Resource{
		Read:   dataSourceLustreInstanceRead,
		Schema: dsSchema_m,
	}
}

func dataSourceLustreInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	//  Get required fields for ID
	instance_id := d.Get("instance_id").(string)

	zone, err := tpgresource.GetZone(d, config)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	// Set the ID
	id := fmt.Sprintf("projects/%s/locations/%s/instances/%s", project, zone, instance_id)
	d.SetId(id)

	// Setting location field for url_param_only field
	d.Set("location", zone)

	err = resourceLustreInstanceRead(d, meta)
	if err != nil {
		return err
	}

	if err := tpgresource.SetDataSourceLabels(d); err != nil {
		return err
	}

	if d.Id() == "" {
		return fmt.Errorf("%s not found", d.Id())
	}

	return nil
}
