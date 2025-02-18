// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package backupdr_test

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"testing"
)

func TestAccDataSourceGoogleBackupDRBackupPlan_basic(t *testing.T) {
	t.Parallel()
	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBackupDRBackupPlanDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleBackupDRBackupPlan_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_backup_dr_backup_plan.fetch-bp", "google_backup_dr_backup_plan.test"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleBackupDRBackupPlan_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_backup_dr_backup_vault" "my-backup-vault-1" {
    location ="us-central1"
    backup_vault_id    = "bv-%{random_suffix}"
    description = "This is a second backup vault built by Terraform."
    backup_minimum_enforced_retention_duration = "100000s"
    labels = {
      foo = "bar1"
      bar = "baz1"
    }
    annotations = {
      annotations1 = "bar1"
      annotations2 = "baz1"
    }
    force_update = "true"
    force_delete = "true"
    allow_missing = "true" 
}


resource "google_backup_dr_backup_plan" "test" { 
  location = "us-central1" 
  backup_plan_id = "bp-test-%{random_suffix}"
  resource_type= "compute.googleapis.com/Instance"
  backup_vault = google_backup_dr_backup_vault.my-backup-vault-1.name
  depends_on=[ google_backup_dr_backup_vault.my-backup-vault-1 ]
  lifecycle {
    ignore_changes = [backup_vault]
  }
  backup_rules {
	rule_id = "rule-1"
	backup_retention_days = 5
	standard_schedule {
	  recurrence_type = "HOURLY"
	   hourly_frequency = 6
	    time_zone = "UTC"
	     backup_window{
		start_hour_of_day = 0
		end_hour_of_day = 24
      }
    }
	}
}

data "google_backup_dr_backup_plan" "fetch-bp" {
  location =  "us-central1"
  backup_plan_id="bp-test-%{random_suffix}"
  depends_on= [ google_backup_dr_backup_plan.test ]
  }
`, context)
}
