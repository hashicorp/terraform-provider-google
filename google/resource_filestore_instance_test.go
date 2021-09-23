package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFilestoreInstance_update(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", randInt(t))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFilestoreInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFilestoreInstance_update(name),
			},
			{
				ResourceName:      "google_filestore_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFilestoreInstance_update2(name),
			},
			{
				ResourceName:      "google_filestore_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccFilestoreInstance_update(name string) string {
	return fmt.Sprintf(`
resource "google_filestore_instance" "instance" {
  name = "tf-instance-%s"
  zone = "us-central1-b"
  file_shares {
    capacity_gb = 2660
    name        = "share"
  }
  networks {
    network = "default"
    modes   = ["MODE_IPV4"]
  }
  labels = {
    baz = "qux"
  }
  tier        = "PREMIUM"
  description = "An instance created during testing."
}
`, name)
}

func testAccFilestoreInstance_update2(name string) string {
	return fmt.Sprintf(`
resource "google_filestore_instance" "instance" {
  name = "tf-instance-%s"
  zone = "us-central1-b"
  file_shares {
    capacity_gb = 2760
    name        = "share"
  }
  networks {
    network = "default"
    modes   = ["MODE_IPV4"]
  }
  tier        = "PREMIUM"
  description = "A modified instance created during testing."
}
`, name)
}

func TestAccFilestoreInstance_reservedIpRange_update(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", randInt(t))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFilestoreInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFilestoreInstance_reservedIpRange_update(name),
			},
			{
				ResourceName:      "google_filestore_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFilestoreInstance_reservedIpRange_update2(name),
			},
			{
				ResourceName:      "google_filestore_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccFilestoreInstance_reservedIpRange_update(name string) string {
	return fmt.Sprintf(`
resource "google_filestore_instance" "instance" {
  name = "tf-instance-%s"
  zone = "us-central1-b"
  tier    = "BASIC_HDD"

  file_shares {
    capacity_gb = 1024
    name        = "share1"
  }

  networks {
    network           = "default"
    modes             = ["MODE_IPV4"]
    reserved_ip_range = "172.19.30.0/29"
  }
}
`, name)
}

func testAccFilestoreInstance_reservedIpRange_update2(name string) string {
	return fmt.Sprintf(`
resource "google_filestore_instance" "instance" {
  name = "tf-instance-%s"
  zone = "us-central1-b"
  tier    = "BASIC_HDD"

  file_shares {
    capacity_gb = 1024
    name        = "share1"
  }

  networks {
    network           = "default"
    modes             = ["MODE_IPV4"]
    reserved_ip_range = "172.19.31.0/29"
  }
}
`, name)
}
