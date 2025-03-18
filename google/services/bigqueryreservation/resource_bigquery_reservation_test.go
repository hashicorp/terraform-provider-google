// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigqueryreservation_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccBigqueryReservation_withDisasterRecovery_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryReservation_withDisasterRecovery_basic(context),
			},
			{
				ResourceName:      "google_bigquery_reservation.reservation",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigqueryReservation_withDisasterRecovery_update(context),
			},
			{
				ResourceName:      "google_bigquery_reservation.reservation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccBigqueryReservation_withDisasterRecovery_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_reservation" "reservation" {
  name           = "tf-test-reservation-%{random_suffix}"
  location       = "us-west2"
  secondary_location = "us-west1"

  // Set to 0 for testing purposes
  // In reality this would be larger than zero
  slot_capacity     = 0
  edition = "ENTERPRISE_PLUS"
  ignore_idle_slots = true
  concurrency       = 0
  autoscale {
    max_slots = 100
  }
}
`, context)
}

func testAccBigqueryReservation_withDisasterRecovery_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_reservation" "reservation" {
  name           = "tf-test-reservation-%{random_suffix}"
  location       = "us-west2"

  // secondary_location is removed. Direct value update (e.g. "us-west1" to "us-east1") is not supported.

  // Set to 0 for testing purposes
  // In reality this would be larger than zero
  slot_capacity     = 0
  edition = "ENTERPRISE_PLUS"
  ignore_idle_slots = true
  concurrency       = 0
  autoscale {
    max_slots = 100
  }
}
`, context)
}
