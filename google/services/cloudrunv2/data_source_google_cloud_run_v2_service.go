// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudrunv2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleCloudRunV2Service() *schema.Resource {
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceCloudRunV2Service().Schema)
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "name")
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "location")

	// Set 'Optional' schema elements
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceGoogleCloudRunV2ServiceRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleCloudRunV2ServiceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/services/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}

	d.SetId(id)

	err = resourceCloudRunV2ServiceRead(d, meta)
	if err != nil {
		return err
	}

	if d.Id() == "" {
		return fmt.Errorf("%s not found", id)
	}
	return nil
}
