// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package backupdr_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceGoogleCloudBackupDRBackup_basic(t *testing.T) {
	t.Parallel()

	project := envvar.GetTestProjectFromEnv()
	location := "us-central1"
	backupVaultId := "bv-test"
	dataSourceId := "ds-test"

	name := fmt.Sprintf("projects/%s/locations/%s/backupVaults/%s/dataSources/%s/backups", project, location, backupVaultId, dataSourceId)

	context := map[string]interface{}{
		"backup_vault_id": backupVaultId,
		"data_source_id":  dataSourceId,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleCloudBackupDRBackup_basic(context),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("data.google_backup_dr_backup.foo", "name", name)),
			},
		},
	})
}

func testAccDataSourceGoogleCloudBackupDRBackup_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

data "google_backup_dr_backup" "foo" {
  project = data.google_project.project.project_id
  location      = "us-central1"
  backup_vault_id = "%{backup_vault_id}"
  data_source_id = "%{data_source_id}"
}

`, context)
}
