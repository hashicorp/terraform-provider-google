package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeReservation_update(t *testing.T) {
	t.Parallel()

	reservationName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeReservationDestroyProducer(t),
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
