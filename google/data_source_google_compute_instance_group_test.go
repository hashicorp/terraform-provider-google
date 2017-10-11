package google

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceGoogleComputeInstanceGroup_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceGoogleComputeInstanceGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceGoogleComputeInstanceGroup("data.google_compute_instance_group.test"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleComputeInstanceGroup_withNamedPort(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceGoogleComputeInstanceGroupConfigWithNamedPort,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceGoogleComputeInstanceGroup("data.google_compute_instance_group.test"),
				),
			},
		},
	})
}

func testAccCheckDataSourceGoogleComputeInstanceGroup(dataSourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		dsFullName := "data.google_compute_instance_group.test"
		rsFullName := "google_compute_instance_group.test"
		ds, ok := s.RootModule().Resources[dsFullName]
		if !ok {
			return fmt.Errorf("cant' find resource called %s in state", dsFullName)
		}

		rs, ok := s.RootModule().Resources[rsFullName]
		if !ok {
			return fmt.Errorf("can't find data source called %s in state", rsFullName)
		}

		dsAttrs := ds.Primary.Attributes
		rsAttrs := rs.Primary.Attributes

		attrsToTest := []string{
			"id",
			"name",
			"zone",
			"project",
			"description",
			"network",
			"self_link",
			"size",
		}

		for _, attrToTest := range attrsToTest {
			if dsAttrs[attrToTest] != rsAttrs[attrToTest] {
				return fmt.Errorf("%s is %s; want %s", attrToTest, dsAttrs[attrToTest], rsAttrs[attrToTest])
			}
		}

		dsNamedPortsCount, ok := dsAttrs["named_port.#"]
		if !ok {
			return errors.New("can't find 'named_port' attribute in data source")
		}

		dsNoOfNamedPorts, err := strconv.Atoi(dsNamedPortsCount)
		if err != nil {
			return errors.New("failed to read number of named ports in data source")
		}

		rsNamedPortsCount, ok := rsAttrs["named_port.#"]
		if !ok {
			return errors.New("can't find 'named_port' attribute in resource")
		}

		rsNoOfNamedPorts, err := strconv.Atoi(rsNamedPortsCount)
		if err != nil {
			return errors.New("failed to read number of named ports in resource")
		}

		if dsNoOfNamedPorts != rsNoOfNamedPorts {
			return fmt.Errorf(
				"expected %d number of named port, received %d, this is most likely a bug",
				rsNoOfNamedPorts,
				dsNoOfNamedPorts,
			)
		}

		namedPortItemKeys := []string{"name", "value"}
		for i := 0; i < dsNoOfNamedPorts; i++ {
			for _, key := range namedPortItemKeys {
				idx := fmt.Sprintf("named_port.%d.%s", i, key)
				if dsAttrs[idx] != rsAttrs[idx] {
					return fmt.Errorf("%s is %s; want %s", idx, dsAttrs[idx], rsAttrs[idx])
				}
			}
		}

		dsInstancesCount, ok := dsAttrs["instances.#"]
		if !ok {
			return errors.New("can't find 'instances' attribute in data source")
		}

		dsNoOfInstances, err := strconv.Atoi(dsInstancesCount)
		if err != nil {
			return errors.New("failed to read number of named ports in data source")
		}

		rsInstancesCount, ok := rsAttrs["instances.#"]
		if !ok {
			return errors.New("can't find 'instances' attribute in resource")
		}

		rsNoOfInstances, err := strconv.Atoi(rsInstancesCount)
		if err != nil {
			return errors.New("failed to read number of instances in resource")
		}

		if dsNoOfInstances != rsNoOfInstances {
			return fmt.Errorf(
				"expected %d number of instances, received %d, this is most likely a bug",
				rsNoOfInstances,
				dsNoOfInstances,
			)
		}

		// We don't know the exact keys of the elements, so go through the whole list looking for matching ones
		dsInstancesValues := []string{}
		for k, v := range dsAttrs {
			if strings.HasPrefix(k, "instances") && !strings.HasSuffix(k, "#") {
				dsInstancesValues = append(dsInstancesValues, v)
			}
		}

		rsInstancesValues := []string{}
		for k, v := range rsAttrs {
			if strings.HasPrefix(k, "instances") && !strings.HasSuffix(k, "#") {
				rsInstancesValues = append(rsInstancesValues, v)
			}
		}

		sort.Strings(dsInstancesValues)
		sort.Strings(rsInstancesValues)

		if !reflect.DeepEqual(dsInstancesValues, rsInstancesValues) {
			return fmt.Errorf("expected %v list of instances, received %v", rsInstancesValues, dsInstancesValues)
		}

		return nil
	}
}

var testAccCheckDataSourceGoogleComputeInstanceGroupConfig = fmt.Sprintf(`
resource "google_compute_instance" "test" {
  name         = "tf-test-%s"
  machine_type = "n1-standard-1"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-8"
    }
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

data "google_compute_instance_group" "test" {
  name = "${google_compute_instance_group.test.name}"
  zone = "${google_compute_instance_group.test.zone}"
}
`, acctest.RandString(10), acctest.RandString(10))

var testAccCheckDataSourceGoogleComputeInstanceGroupConfigWithNamedPort = fmt.Sprintf(`
resource "google_compute_instance" "test" {
  name         = "tf-test-%s"
  machine_type = "n1-standard-1"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-8"
    }
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

  named_port {
    name = "http"
    port = "8080"
  }

  named_port {
    name = "https"
    port = "8443"
  }

  instances = [
    "${google_compute_instance.test.self_link}",
  ]
}

data "google_compute_instance_group" "test" {
  name = "${google_compute_instance_group.test.name}"
  zone = "${google_compute_instance_group.test.zone}"
}
`, acctest.RandString(10), acctest.RandString(10))
