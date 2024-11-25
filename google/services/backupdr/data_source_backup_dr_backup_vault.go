// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package backupdr

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleCloudBackupDRBackupVault() *schema.Resource {

	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceBackupDRBackupVault().Schema)

	tpgresource.AddRequiredFieldsToSchema(dsSchema, "backup_vault_id", "location")

	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceGoogleCloudBackupDRBackupVaultRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleCloudBackupDRBackupVaultRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	location, err := tpgresource.GetLocation(d, config)
	if err != nil {
		return err
	}
	backup_vault_id := d.Get("backup_vault_id").(string)
	id := fmt.Sprintf("projects/%s/locations/%s/backupVaults/%s", project, location, backup_vault_id)
	d.SetId(id)
	err = resourceBackupDRBackupVaultRead(d, meta)
	if err != nil {
		return err
	}
	if d.Id() == "" {
		return fmt.Errorf("%s not found", id)
	}

	return nil
}
