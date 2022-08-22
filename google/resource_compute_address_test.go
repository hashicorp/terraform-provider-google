package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeAddress_networkTier(t *testing.T) {
	t.Parallel()

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeAddress_networkTier(randString(t, 10)),
			},
			{
				ResourceName:      "google_compute_address.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeAddress_internal(t *testing.T) {
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeAddress_internal(randString(t, 10)),
			},
			{
				ResourceName:      "google_compute_address.internal",
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				ResourceName:      "google_compute_address.internal_with_subnet",
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				ResourceName:      "google_compute_address.internal_with_subnet_and_address",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeAddress_internal(i string) string {
	return fmt.Sprintf(`
resource "google_compute_address" "internal" {
  name         = "address-test-internal-%s"
  address_type = "INTERNAL"
  region       = "us-east1"
}

resource "google_compute_network" "default" {
  name = "tf-test-network-test-%s"
}

resource "google_compute_subnetwork" "foo" {
  name          = "subnetwork-test-%s"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-east1"
  network       = google_compute_network.default.self_link
}

resource "google_compute_address" "internal_with_subnet" {
  name         = "address-test-internal-with-subnet-%s"
  subnetwork   = google_compute_subnetwork.foo.self_link
  address_type = "INTERNAL"
  region       = "us-east1"
}

// We can't test the address alone, because we don't know what IP range the
// default subnetwork uses.
resource "google_compute_address" "internal_with_subnet_and_address" {
  name         = "address-test-internal-with-subnet-and-address-%s"
  subnetwork   = google_compute_subnetwork.foo.self_link
  address_type = "INTERNAL"
  address      = "10.0.42.42"
  region       = "us-east1"
}
`,
		i, // google_compute_address.internal name
		i, // google_compute_network.default name
		i, // google_compute_subnetwork.foo name
		i, // google_compute_address.internal_with_subnet_name
		i, // google_compute_address.internal_with_subnet_and_address name
	)
}

func testAccComputeAddress_networkTier(i string) string {
	return fmt.Sprintf(`
resource "google_compute_address" "foobar" {
  name         = "address-test-%s"
  network_tier = "STANDARD"
}
`, i)
}
