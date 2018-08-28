package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// TODO (chrisst)
// remove debian-9 from hard coded compute instance
//

// Acceptance Tests
// TEST - make sure count(N) results in N+1 disks on the instance

func TestAccAttachedDisk_basic(t *testing.T) {
	t.Parallel()

	diskName := acctest.RandomWithPrefix("tf-test-disk")
	instanceName := acctest.RandomWithPrefix("tf-test-inst")
	importID := fmt.Sprintf("%s/us-central1-a/%s:%s", getTestProjectFromEnv(), instanceName, diskName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAttachedDiskResource(diskName, instanceName) + testAttachedDiskResourceAttachment(),
			},
			resource.TestStep{
				ResourceName:      "google_compute_attached_disk.test",
				ImportStateId:     importID,
				ImportState:       true,
				ImportStateVerify: true,
			},
			resource.TestStep{
				Config: testAttachedDiskResource(diskName, instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccAttachedDiskIsNowDetached(instanceName, diskName),
				),
			},
		},
	})
}

func TestAccAttachedDisk_count(t *testing.T) {
	t.Parallel()

	diskPrefix := acctest.RandomWithPrefix("tf-test")
	instanceName := acctest.RandomWithPrefix("tf-test")
	count := 2

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAttachedDiskResourceCount(diskPrefix, instanceName, count),
				Check: resource.ComposeTestCheckFunc(
					testAccAttachedDiskContainsManyDisks(instanceName, count),
				),
			},
		},
	})

}

func testAccAttachedDiskIsNowDetached(instanceName, diskName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		instance, err := config.clientCompute.Instances.Get(getTestProjectFromEnv(), "us-central1-a", instanceName).Do()
		if err != nil {
			return err
		}

		ad := findDiskByName(instance.Disks, diskName)
		if ad != nil {
			return fmt.Errorf("compute disk is still attached to compute instance")
		}

		return nil
	}
}

func testAccAttachedDiskContainsManyDisks(instanceName string, count int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		instance, err := config.clientCompute.Instances.Get(getTestProjectFromEnv(), "us-central1-a", instanceName).Do()
		if err != nil {
			return err
		}

		// There will always be 1 extra disk because of the compute instance's boot disk
		if (count + 1) != len(instance.Disks) {
			return fmt.Errorf("expected %d disks to be attached, found %d", count+1, len(instance.Disks))
		}

		return nil
	}
}

func testAttachedDiskResourceAttachment() string {
	return fmt.Sprintf(`
resource "google_compute_attached_disk" "test" {
	attached_disk = "${google_compute_disk.test1.self_link}"
	attached_instance = "${google_compute_instance.test.self_link}"
}
	`)
}

func testAttachedDiskResource(diskName, instanceName string) string {
	return fmt.Sprintf(`
resource "google_compute_disk" "test1" {
	name = "%s"
	zone = "us-central1-a"
	size = 10
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
`, diskName, instanceName)
}

func testAttachedDiskResourceCount(diskPrefix, instanceName string, count int) string {
	return fmt.Sprintf(`
resource "google_compute_disk" "many" {
	name = "%s-${count.index}"
	zone = "us-central1-a"
	size = 10
	count = %d
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
	count = "${google_compute_disk.many.count}"
	attached_disk = "${google_compute_disk.many.*.self_link[count.index]}"
	attached_instance = "${google_compute_instance.test.self_link}"
}
`, diskPrefix, count, instanceName)
}
