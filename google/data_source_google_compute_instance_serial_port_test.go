package google

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceComputeInstanceSerialPort_basic(t *testing.T) {
	instanceName := fmt.Sprintf("tf-test-serial-data-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceSerialPort(instanceName),
				Check: resource.ComposeTestCheckFunc(
					// Contents of serial port output include lots of initialization logging
					resource.TestMatchResourceAttr("data.google_compute_instance_serial_port.serial", "contents",
						regexp.MustCompile("Initializing cgroup subsys")),
				),
			},
		},
	})
}

func testAccComputeInstanceSerialPort(instanceName string) string {
	return fmt.Sprintf(`
resource "google_compute_instance" "default" {
	name				 = "%s"
	machine_type = "n1-standard-1"
	zone				 = "us-central1-a"

	boot_disk {
		initialize_params {
			image = "debian-8-jessie-v20160803"
		}
	}

	// Local SSD disk
	scratch_disk {
		interface = "SCSI"
	}

	network_interface {
		network = "default"

		access_config {
			// Ephemeral IP
		}
	}

	metadata = {
		foo = "bar"
		serial-port-logging-enable = "TRUE"
		windows-keys = jsonencode(
			{
				email		 = "example.user@example.com"
				expireOn = "2020-04-14T01:37:19Z"
				exponent = "AQAB"
				modulus	 = "wgsquN4IBNPqIUnu+h/5Za1kujb2YRhX1vCQVQAkBwnWigcCqOBVfRa5JoZfx6KIvEXjWqa77jPvlsxM4WPqnDIM2qiK36up3SKkYwFjff6F2ni/ry8vrwXCX3sGZ1hbIHlK0O012HpA3ISeEswVZmX2X67naOvJXfY5v0hGPWqCADao+xVxrmxsZD4IWnKl1UaZzI5lhAzr8fw6utHwx1EZ/MSgsEki6tujcZfN+GUDRnmJGQSnPTXmsf7Q4DKreTZk49cuyB3prV91S0x3DYjCUpSXrkVy1Ha5XicGD/q+ystuFsJnrrhbNXJbpSjM6sjo/aduAkZJl4FmOt0R7Q=="
				userName = "example-user"
			}
		)
	}

	service_account {
		scopes = ["userinfo-email", "compute-ro", "storage-ro"]
	}
}

data "google_compute_instance_serial_port" "serial" {
	instance = google_compute_instance.default.name
	zone		 = google_compute_instance.default.zone
	port		 = 1
}
`, instanceName)
}
