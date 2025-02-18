// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package backupdr_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"strings"
	"testing"
)

func TestAccDataSourceGoogleBackupDRManagementServer_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "backupdr-managementserver-basic"),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBackupDRManagementServerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleBackupDRManagementServer_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_backup_dr_management_server.foo", "google_backup_dr_management_server.foo"),
				),
			},
		},
	})
}

func testAccCheckBackupDRManagementServerDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_backup_dr_management_server" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, `{{BackupDRBasePath}}projects/{{project}}/locations/{{location}}/managementServers`)
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
				return fmt.Errorf("BackupDRManagementServer still exists at %s", url)
			}
		}

		return nil
	}
}

func testAccDataSourceGoogleBackupDRManagementServer_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_network" "default" {
  name = "%{network_name}"
}
resource "google_backup_dr_management_server" "foo" {
 location = "us-central1"
  name     = "tf-test-management-server%{random_suffix}"
  type     = "BACKUP_RESTORE" 
    networks {
    network      = data.google_compute_network.default.id
    peering_mode = "PRIVATE_SERVICE_ACCESS"
  }
}

data "google_backup_dr_management_server" "foo" {
  location =  "us-central1"
  depends_on = [ google_backup_dr_management_server.foo ]
}
`, context)
}
