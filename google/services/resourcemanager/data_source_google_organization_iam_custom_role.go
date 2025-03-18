// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func DataSourceGoogleOrganizationIamCustomRole() *schema.Resource {
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceGoogleOrganizationIamCustomRole().Schema)

	tpgresource.AddRequiredFieldsToSchema(dsSchema, "org_id")
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "role_id")

	return &schema.Resource{
		Read:   dataSourceOrganizationIamCustomRoleRead,
		Schema: dsSchema,
	}
}

func dataSourceOrganizationIamCustomRoleRead(d *schema.ResourceData, meta interface{}) error {
	orgId := d.Get("org_id").(string)
	roleId := d.Get("role_id").(string)
	d.SetId(fmt.Sprintf("organizations/%s/roles/%s", orgId, roleId))

	id := d.Id()

	if err := resourceGoogleOrganizationIamCustomRoleRead(d, meta); err != nil {
		return err
	}

	if d.Id() == "" {
		return fmt.Errorf("Role %s not found!", id)
	}

	return nil
}
