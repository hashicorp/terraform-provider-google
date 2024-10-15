// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package netapp_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccNetappBackupVault_NetappBackupVaultExample_update(t *testing.T) {
	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetappBackupVault_NetappBackupVaultExample_basic(context),
			},
			{
				ResourceName:            "google_netapp_backup_vault.test_backup_vault",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "labels", "terraform_labels"},
			},
			{
				Config: testAccNetappBackupVault_NetappBackupVaultExample_update(context),
			},
			{
				ResourceName:            "google_netapp_backup_vault.test_backup_vault",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccNetappBackupVault_NetappBackupVaultExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_netapp_backup_vault" "test_backup_vault" {
  name = "tf-test-test-backup-vault%{random_suffix}"
  location = "us-east4"
}
`, context)
}

func testAccNetappBackupVault_NetappBackupVaultExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_netapp_backup_vault" "test_backup_vault" {
  name = "tf-test-test-backup-vault%{random_suffix}"
  location = "us-east4"
  description = "Terraform created vault"
  labels = { 
    "creator": "testuser",
	"foo": "bar",
  }
}
`, context)
}

func testAccCheckNetappBackupVaultDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_netapp_backup_vault" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{NetappBasePath}}projects/{{project}}/locations/{{location}}/backupVaults/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("NetappBackupVault still exists at %s", url)
			}
		}

		return nil
	}
}
