package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// TODO (chrisst)
// Make sure to get the ZONE from the compute instance if it's not there on project level
// remove debian-9 from hard coded compute instance
//

// Smoke Tests
// * test renaming a disk
// * test renaming a compute instance

// Acceptance Tests
// TEST - make sure count(N) results in N+1 disks on the instance

// Questions:
// How do I properly test

func TestAttachedDisk_basic(t *testing.T) {
	t.Parallel()

	diskName := acctest.RandomWithPrefix("tf-test")
	instanceName := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAttachedDiskDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAttachedDiskResource(diskName, instanceName),
			},
			resource.TestStep{
				ResourceName:      "google_compute_attached_disk.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

}

func testAccCheckAttachedDiskDestroy(s *terraform.State) error {
	fmt.Println("ZOMG testing destory")
	// config := testAccProvider.Meta().(*Config)

	// TODO (chrisst) - figure out how to test that things are deleted?
	// maybe check that instance + disk are deleted?

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_compute_attached_disk" {
			continue
		}
	}

	return nil
}

func testAttachedDiskResource(diskName, instanceName string) string {
	return fmt.Sprintf(`
resource "google_compute_disk" "test1" {
	name = "%s"
	zone = "us-central1-a"
}

resource "google_compute_instance" "test" {
	name         = "%s"
	machine_type = "f1-micro"
	zone         = "us-central1-a"
	
	lifecycle {
			ignore_changes = [
				"attached_disk"
			]
		}
	
	boot_disk {
		initialize_params {
		image = "debian-cloud/debian-9"
		}
	}
	
	network_interface {
		network = "default"
	}
}

resource "google_compute_attached_disk" "test" {
	attached_disk = "${google_compute_disk.test1.self_link}"
	attached_instance = "${google_compute_instance.test.self_link}"
}`, diskName, instanceName)
}
