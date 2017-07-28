package google

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccGoogleComputeInstanceGroupInstances_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleComputeInstanceGroupInstancesConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleComputeInstanceGroupInstances("data.google_compute_instance_group_instances.all"),
				),
			},
		},
	})
}

func testAccCheckGoogleComputeInstanceGroupInstances(dataSourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[dataSourceName]
		if !ok {
			return fmt.Errorf("Can't find instance group instances data source: %s", dataSourceName)
		}

		if ds.Primary.ID == "" {
			return fmt.Errorf("%s data source ID not set", dataSourceName)
		}

		dsAttrs := ds.Primary.Attributes

		countInstances, ok := dsAttrs["instances.#"]
		if !ok {
			return errors.New("can't find 'instances' attribute")
		}

		nOfInstances, err := strconv.Atoi(countInstances)
		if err != nil {
			return errors.New("failed to read number of instances")
		}

		countNames, ok := dsAttrs["names.#"]
		if !ok {
			return errors.New("can't find 'names' attribute")
		}

		nOfNames, err := strconv.Atoi(countNames)
		if err != nil {
			return errors.New("failed to read number of instances")
		}

		for i := 0; i < nOfInstances; i++ {
			idx := "instances." + strconv.Itoa(i)
			v, ok := dsAttrs[idx]
			if !ok {
				return fmt.Errorf("instance list is corrupt (%q not found), this is definitely a bug", idx)
			}
			if len(v) < 1 {
				return fmt.Errorf("Empty instance value (%q), this is definitely a bug", idx)
			}
		}

		for i := 0; i < nOfNames; i++ {
			idx := "names." + strconv.Itoa(i)
			v, ok := dsAttrs[idx]
			if !ok {
				return fmt.Errorf("name list is corrupt (%q not found), this is definitely a bug", idx)
			}
			if len(v) < 1 {
				return fmt.Errorf("Empty name value (%q), this is definitely a bug", idx)
			}
		}

		return nil
	}
}

var testAccCheckGoogleComputeInstanceGroupInstancesConfig = fmt.Sprintf(`
resource "google_compute_instance" "test" {
  name         = "tf-test-%s"
  machine_type = "n1-standard-1"
  zone         = "us-central1-a"

  disk {
    image = "debian-cloud/debian-8"
  }

  network_interface {
    network = "default"

    access_config {
      // Ephemeral IP
    }
  }
}

resource "google_compute_instance_group" "test" {
  name = "tf-test-%s"
  zone = "${google_compute_instance.test.zone}"

  instances = [
    "${google_compute_instance.test.self_link}",
  ]
}

data "google_compute_instance_group_instances" "all" {
	name = "${google_compute_instance_group.test.name}"
	zone = "${google_compute_instance_group.test.zone}"
}
`, acctest.RandString(10), acctest.RandString(10))
