// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package netapp_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccNetappbackupPolicy_netappBackupPolicyFullExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetappbackupPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetappbackupPolicy_netappBackupPolicyFullExample_basic(context),
			},
			{
				ResourceName:            "google_netapp_backup_policy.test_backup_policy_full",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "labels", "terraform_labels"},
			},
			{
				Config: testAccNetappbackupPolicy_netappBackupPolicyFullExample_updates(context),
			},
			{
				ResourceName:            "google_netapp_backup_policy.test_backup_policy_full",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "labels", "terraform_labels"},
			}, {
				Config: testAccNetappbackupPolicy_netappBackupPolicyFullExample_disable(context),
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
func testAccNetappbackupPolicy_netappBackupPolicyFullExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_netapp_backup_policy" "test_backup_policy_full" {
  name          = "tf-test-test-backup-policy-full%{random_suffix}"
  location = "us-central1"
  daily_backup_limit   = 2
  weekly_backup_limit  = 0
  monthly_backup_limit = 0
}
`, context)
}

// Update all fields
func testAccNetappbackupPolicy_netappBackupPolicyFullExample_updates(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_netapp_backup_policy" "test_backup_policy_full" {
  name          = "tf-test-test-backup-policy-full%{random_suffix}"
  location = "us-central1"
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
func testAccNetappbackupPolicy_netappBackupPolicyFullExample_disable(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_netapp_backup_policy" "test_backup_policy_full" {
  name          = "tf-test-test-backup-policy-full%{random_suffix}"
  location = "us-central1"
  daily_backup_limit   = 2
  weekly_backup_limit  = 1
  monthly_backup_limit = 1
  enabled = false
}
`, context)
}
