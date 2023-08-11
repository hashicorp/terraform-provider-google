// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccComputeReservation_update(t *testing.T) {
	t.Parallel()

	reservationName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeReservationDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeReservation_basic(reservationName, "2"),
			},
			{
				ResourceName:      "google_compute_reservation.reservation",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeReservation_basic(reservationName, "1"),
			},
			{
				ResourceName:      "google_compute_reservation.reservation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeReservation_basic(reservationName, count string) string {
	return fmt.Sprintf(`
resource "google_compute_reservation" "reservation" {
  name = "%s"
  zone = "us-central1-a"

  specific_reservation {
    count = %s
    instance_properties {
      min_cpu_platform = "Intel Cascade Lake"
      machine_type     = "n2-standard-2"
    }
  }
}
`, reservationName, count)
}
