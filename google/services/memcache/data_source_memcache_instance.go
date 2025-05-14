// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package memcache

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceMemcacheInstance() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceMemcacheInstance().Schema)

	// Set 'Required' schema elements
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "name")
	// Set 'Optional' schema elements
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project", "region")

	return &schema.Resource{
		Read:   dataSourceMemcacheInstanceRead,
		Schema: dsSchema,
	}
}

func dataSourceMemcacheInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	region, err := tpgresource.GetRegion(d, config)
	if err != nil {
		return err
	}

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{region}}/instances/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	// Setting location field, as this is set as a required field in instance resource to build the url
	d.Set("region", region)

	err = resourceMemcacheInstanceRead(d, meta)
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
