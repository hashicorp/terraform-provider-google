package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccFilestoreInstance_basic(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFilestoreInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccFilestoreInstance_basic(name),
			},
			resource.TestStep{
				ResourceName:      "google_filestore_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccFilestoreInstance_update(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFilestoreInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccFilestoreInstance_update(name),
			},
			resource.TestStep{
				ResourceName:      "google_filestore_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
			resource.TestStep{
				Config: testAccFilestoreInstance_update2(name),
			},
			resource.TestStep{
				ResourceName:      "google_filestore_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckFilestoreInstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_filestore_instance" {
			continue
		}

		redisIdParts := strings.Split(rs.Primary.ID, "/")
		if len(redisIdParts) != 3 {
			return fmt.Errorf("Unexpected resource ID %s, expected {project}/{region}/{name}", rs.Primary.ID)
		}

		project, region, inst := redisIdParts[0], redisIdParts[1], redisIdParts[2]

		name := fmt.Sprintf("projects/%s/locations/%s/instances/%s", project, region, inst)
		_, err := config.clientFilestore.Projects.Locations.Get(name).Do()
		if err == nil {
			return fmt.Errorf("Filestore instance still exists")
		}
	}

	return nil
}

func testAccFilestoreInstance_basic(name string) string {
	return fmt.Sprintf(`
resource "google_filestore_instance" "instance" {
  name = "tf-instance-%s"
  zone = "us-central1-b"
  file_shares {
    capacity_gb = 2560
    name = "share"
  }
  networks {
    network = "default"
    modes = ["MODE_IPV4"]
  }
	labels {
		foo = "bar"
	}
  tier = "PREMIUM"
}
`, name)
}

func testAccFilestoreInstance_update(name string) string {
	return fmt.Sprintf(`
resource "google_filestore_instance" "instance" {
  name = "tf-instance-%s"
  zone = "us-central1-b"
  file_shares {
    capacity_gb = 2660
    name = "share"
  }
  networks {
    network = "default"
    modes = ["MODE_IPV4"]
  }
	labels {
		baz = "qux"
	}
  tier = "PREMIUM"
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
    name = "share"
  }
  networks {
    network = "default"
    modes = ["MODE_IPV4"]
  }
  tier = "PREMIUM"
	description = "A modified instance created during testing."
}`, name)
}
