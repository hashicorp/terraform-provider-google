// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package netapp_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccNetappBackupPolicy_NetappBackupPolicyFullExample_update(t *testing.T) {
	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetappBackupPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetappBackupPolicy_NetappBackupPolicyFullExample_basic(context),
			},
			{
				ResourceName:            "google_netapp_backup_policy.test_backup_policy_full",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "labels", "terraform_labels"},
			},
			{
				Config: testAccNetappBackupPolicy_NetappBackupPolicyFullExample_updates(context),
			},
			{
				ResourceName:            "google_netapp_backup_policy.test_backup_policy_full",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "labels", "terraform_labels"},
			}, {
				Config: testAccNetappBackupPolicy_NetappBackupPolicyFullExample_disable(context),
			},
			{
				ResourceName:            "google_netapp_backup_policy.test_backup_policy_full",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "labels", "terraform_labels"},
			},
		},
	})
}

// Setup minimal policy
func testAccNetappBackupPolicy_NetappBackupPolicyFullExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_netapp_backup_policy" "test_backup_policy_full" {
  name          = "tf-test-test-backup-policy-full%{random_suffix}"
  location = "us-east4"
  daily_backup_limit   = 2
  weekly_backup_limit  = 0
  monthly_backup_limit = 0
}
`, context)
}

// Update all fields
func testAccNetappBackupPolicy_NetappBackupPolicyFullExample_updates(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_netapp_backup_policy" "test_backup_policy_full" {
  name          = "tf-test-test-backup-policy-full%{random_suffix}"
  location = "us-east4"
  daily_backup_limit   = 6
  weekly_backup_limit  = 4
  monthly_backup_limit = 3
  description = "TF test backup schedule"
  enabled = true
  labels = {
    "foo" = "bar"
  }
}
`, context)
}

// test disabling the policy
func testAccNetappBackupPolicy_NetappBackupPolicyFullExample_disable(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_netapp_backup_policy" "test_backup_policy_full" {
  name          = "tf-test-test-backup-policy-full%{random_suffix}"
  location = "us-east4"
  daily_backup_limit   = 2
  weekly_backup_limit  = 1
  monthly_backup_limit = 1
  enabled = false
}
`, context)
}
