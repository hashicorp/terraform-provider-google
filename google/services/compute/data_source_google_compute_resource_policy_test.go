// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceComputeResourcePolicy(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)

	rsName := "foo_" + randomSuffix
	rsFullName := fmt.Sprintf("google_compute_resource_policy.%s", rsName)
	dsName := "my_policy_" + randomSuffix
	dsFullName := fmt.Sprintf("data.google_compute_resource_policy.%s", dsName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataSourceComputeResourcePolicyDestroy(t, rsFullName),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComputeResourcePolicyConfig(rsName, dsName, randomSuffix),
				Check:  acctest.CheckDataSourceStateMatchesResourceState(rsFullName, dsFullName),
			},
		},
	})
}

func testAccCheckDataSourceComputeResourcePolicyDestroy(t *testing.T, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_resource_policy" {
				continue
			}

			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			policyAttrs := rs.Primary.Attributes

			_, err := config.NewComputeClient(config.UserAgent).ResourcePolicies.Get(
				config.Project, policyAttrs["region"], policyAttrs["name"]).Do()
			if err == nil {
				return fmt.Errorf("Resource Policy still exists")
			}
		}

		return nil
	}
}

func testAccDataSourceComputeResourcePolicyConfig(rsName, dsName, randomSuffix string) string {
	return fmt.Sprintf(`
resource "google_compute_resource_policy" "%s" {
  name   = "policy-%s"
  region = "us-central1"
  snapshot_schedule_policy {
    schedule {
      daily_schedule {
        days_in_cycle = 1
        start_time    = "04:00"
      }
    }
  }
}

data "google_compute_resource_policy" "%s" {
  name     = google_compute_resource_policy.%s.name
  region   = google_compute_resource_policy.%s.region
}
`, rsName, randomSuffix, dsName, rsName, rsName)
}
