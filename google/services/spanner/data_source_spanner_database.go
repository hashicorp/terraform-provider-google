// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package spanner

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceSpannerDatabase() *schema.Resource {

	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceSpannerDatabase().Schema)

	tpgresource.AddRequiredFieldsToSchema(dsSchema, "name")
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "instance")

	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceSpannerDatabaseRead,
		Schema: dsSchema,
	}
}

func dataSourceSpannerDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	id, err := tpgresource.ReplaceVars(d, config, "{{instance}}/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}

	d.SetId(id)

	err = resourceSpannerDatabaseRead(d, meta)
	if err != nil {
		return err
	}

	if d.Id() == "" {
		return fmt.Errorf("%s not found", id)
	}

	return nil
}
