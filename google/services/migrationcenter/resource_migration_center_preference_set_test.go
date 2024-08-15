// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package migrationcenter_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccMigrationCenterPreferenceSet_preferenceSetUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMigrationCenterPreferenceSetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMigrationCenterPreferenceSet_preferenceSetStart(context),
			},
			{
				ResourceName:            "google_migration_center_preference_set.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "preference_set_id"},
			},
		},
	})
}

func testAccMigrationCenterPreferenceSet_preferenceSetStart(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_migration_center_preference_set" "default" {
  location          = "us-central1"
  preference_set_id = "tf-test-preference-set-test%{random_suffix}"
  description       = "Terraform integration test description"
  display_name      = "Terraform integration test display"
  virtual_machine_preferences {
    vmware_engine_preferences {
      cpu_overcommit_ratio = 1.5
      memory_overcommit_ratio = 2.0
    }
    sizing_optimization_strategy = "SIZING_OPTIMIZATION_STRATEGY_SAME_AS_SOURCE"
  }
}
`, context)
}

func testAccMigrationCenterPreferenceSet_preferenceSetUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_migration_center_preference_set" "default" {
  location          = "us-central1"
  preference_set_id = "tf-test-preference-set-test%{random_suffix}"
  description       = "Terraform integration test updated description"
  display_name      = "Terraform integration test updated display"
  virtual_machine_preferences {
    vmware_engine_preferences {
      cpu_overcommit_ratio = 1.4
    }
    sizing_optimization_strategy = "SIZING_OPTIMIZATION_STRATEGY_MODERATE"
    commitment_plan = "COMMITMENT_PLAN_ONE_YEAR"
    preferred_regions = ["us-central1"]
  }
}
`, context)
}
