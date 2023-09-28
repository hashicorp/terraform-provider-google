// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func DataSourceGoogleFolderOrganizationPolicy() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceGoogleFolderOrganizationPolicy().Schema)

	tpgresource.AddRequiredFieldsToSchema(dsSchema, "folder")
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "constraint")

	return &schema.Resource{
		Read:   datasourceGoogleFolderOrganizationPolicyRead,
		Schema: dsSchema,
	}
}

func datasourceGoogleFolderOrganizationPolicyRead(d *schema.ResourceData, meta interface{}) error {

	id := fmt.Sprintf("%s/%s", d.Get("folder"), d.Get("constraint"))
	d.SetId(id)

	err := resourceGoogleFolderOrganizationPolicyRead(d, meta)
	if err != nil {
		return err
	}

	if d.Id() == "" {
		return fmt.Errorf("%s not found", id)
	}
	return nil
}
