// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleProjectIamCustomRole() *schema.Resource {
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceGoogleProjectIamCustomRole().Schema)

	dsSchema["project"].Computed = false
	dsSchema["project"].Optional = true
	dsSchema["role_id"].Computed = false
	dsSchema["role_id"].Required = true

	return &schema.Resource{
		Read:   dataSourceProjectIamCustomRoleRead,
		Schema: dsSchema,
	}
}

func dataSourceProjectIamCustomRoleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for service accounts: %s", err)
	}

	roleId := d.Get("role_id").(string)
	d.SetId(fmt.Sprintf("projects/%s/roles/%s", project, roleId))

	id := d.Id()

	if err := resourceGoogleProjectIamCustomRoleRead(d, meta); err != nil {
		return err
	}

	if d.Id() == "" {
		return fmt.Errorf("Role %s not found!", id)
	}

	return nil
}
