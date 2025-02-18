// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package backupdr_test

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"testing"
)

func TestAccDataSourceGoogleBackupDRBackupPlanAssociation_basic(t *testing.T) {
	t.Parallel()
	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleBackupDRBackupPlanAssociation_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores("data.google_backup_dr_backup_plan_association.bpa-test", "google_backup_dr_backup_plan_association.bpa", map[string]struct{}{
						"resource": {},
					},
					),
				),
			},
		},
	})
}

func testAccDataSourceGoogleBackupDRBackupPlanAssociation_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_service_account" "default" {
  account_id   = "tf-test-my-custom-%{random_suffix}"
  display_name = "Custom SA for VM Instance"
}

resource "google_compute_instance" "default" {
  name         = "tf-test-compute-instance-%{random_suffix}"
  machine_type = "n2-standard-2"
  zone         = "us-central1-a"
  tags = ["foo", "bar"]
  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
      labels = {
        my_label = "value"
      }
    }
  }
  // Local SSD disk
  scratch_disk {
    interface = "NVME"
  }
  network_interface {
    network = "default"
    access_config {
      // Ephemeral public IP
    }
  }
  service_account {
    # Google recommends custom service accounts that have cloud-platform scope and permissions granted via IAM Roles.
    email  = google_service_account.default.email
    scopes = ["cloud-platform"]
  }
}
resource "google_backup_dr_backup_vault" "my-backup-vault" {
    location ="us-central1"
    backup_vault_id    = "tf-test-bv-%{random_suffix}"
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

resource "google_backup_dr_backup_plan" "foo" {
  location       = "us-central1"
  backup_plan_id = "tf-test-bp-test-%{random_suffix}"
  resource_type  = "compute.googleapis.com/Instance"
  backup_vault   = google_backup_dr_backup_vault.my-backup-vault.name

  backup_rules {
    rule_id                = "rule-1"
    backup_retention_days  = 2

    standard_schedule {
      recurrence_type     = "HOURLY"
      hourly_frequency    = 6
      time_zone           = "UTC"

      backup_window {
        start_hour_of_day = 12
        end_hour_of_day   = 18
      }
    }
  }
}

resource "google_backup_dr_backup_plan_association" "bpa" { 
  location = "us-central1" 
  backup_plan_association_id = "tf-test-bpa-test-%{random_suffix}"
  resource =   google_compute_instance.default.id
  resource_type= "compute.googleapis.com/Instance"
  backup_plan = google_backup_dr_backup_plan.foo.name
}

data "google_backup_dr_backup_plan_association" "bpa-test" {
  location =  "us-central1"
  backup_plan_association_id="tf-test-bpa-test-%{random_suffix}"
  depends_on= [ google_backup_dr_backup_plan_association.bpa ]
  }
`, context)
}
