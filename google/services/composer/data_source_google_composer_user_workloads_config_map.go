// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package composer

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleComposerUserWorkloadsConfigMap() *schema.Resource {
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceComposerUserWorkloadsConfigMap().Schema)

	// Set 'Required' schema elements
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "environment", "name")

	// Set 'Optional' schema elements
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project", "region")

	return &schema.Resource{
		Read:   dataSourceGoogleComposerUserWorkloadsConfigMapRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleComposerUserWorkloadsConfigMapRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{region}}/environments/{{environment}}/userWorkloadsConfigMaps/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = resourceComposerUserWorkloadsConfigMapRead(d, meta)
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
