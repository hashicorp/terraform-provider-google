// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceComputeReservation(t *testing.T) {
	t.Parallel()

	reservationName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	rsName := "foobar"
	dsName := "my_reservation"
	rsFullName := fmt.Sprintf("google_compute_reservation.%s", rsName)
	dsFullName := fmt.Sprintf("data.google_compute_reservation.%s", dsName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeReservationDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComputeReservationConfig(reservationName, rsName, dsName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsFullName, "status", "READY"),
					acctest.CheckDataSourceStateMatchesResourceState(dsFullName, rsFullName),
				),
			},
		},
	})
}

func testAccDataSourceComputeReservationConfig(reservationName, rsName, dsName string) string {
	return fmt.Sprintf(`
resource "google_compute_reservation" "%s" {
  name = "%s"
  zone = "us-west1-a"

  specific_reservation {
    count = 1
    instance_properties {
      min_cpu_platform = "Intel Cascade Lake"
      machine_type     = "n2-standard-2"
    }
  }
}

data "google_compute_reservation" "%s" {
  name = google_compute_reservation.%s.name
  zone = "us-west1-a"
}
`, rsName, reservationName, dsName, rsName)
}
