// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storagecontrol

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleStorageControlOrganizationIntelligenceConfig() *schema.Resource {

	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceStorageControlOrganizationIntelligenceConfig().Schema)
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "name")

	return &schema.Resource{
		Read:   dataSourceGoogleStorageControlOrganizationIntelligenceConfigRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleStorageControlOrganizationIntelligenceConfigRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	id, err := tpgresource.ReplaceVars(d, config, "organizations/{{name}}/locations/global/intelligenceConfig")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)
	err = resourceStorageControlOrganizationIntelligenceConfigRead(d, meta)
	if err != nil {
		return err
	}

	if d.Id() == "" {
		return fmt.Errorf("%s not found", id)
	}

	return nil
}
