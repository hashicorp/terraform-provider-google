// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package vmwareengine

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceVmwareengineCluster() *schema.Resource {

	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceVmwareengineCluster().Schema)
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "parent", "name")
	return &schema.Resource{
		Read:   dataSourceVmwareengineClusterRead,
		Schema: dsSchema,
	}
}

func dataSourceVmwareengineClusterRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "{{parent}}/clusters/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)
	err = resourceVmwareengineClusterRead(d, meta)
	if err != nil {
		return err
	}

	if d.Id() == "" {
		return fmt.Errorf("%s not found", id)
	}
	return nil
}
