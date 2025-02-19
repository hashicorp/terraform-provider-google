// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package alloydb

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceAlloydbDatabaseCluster() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceAlloydbCluster().Schema)

	// Set 'Required' schema elements
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "cluster_id")

	// Set 'Optional' schema elements
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "location", "project")

	return &schema.Resource{
		Read:   dataSourceAlloydbDatabaseClusterRead,
		Schema: dsSchema,
	}
}

func dataSourceAlloydbDatabaseClusterRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	// Get location for setting location field in resource
	location, err := tpgresource.GetLocation(d, config)
	if err != nil {
		return err
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/clusters/{{cluster_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	// Setting location field as this is set as a required field in cluster resource
	d.Set("location", location)

	err = resourceAlloydbClusterRead(d, meta)
	if err != nil {
		return err
	}

	if err := tpgresource.SetDataSourceLabels(d); err != nil {
		return err
	}

	if d.Id() == "" {
		return fmt.Errorf("%s not found", id)
	}
	return nil
}
