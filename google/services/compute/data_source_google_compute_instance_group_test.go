// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func TestAccDataSourceGoogleComputeInstanceGroup_basic(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceGoogleComputeInstanceGroupConfig(acctest.RandString(t, 10), acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceGoogleComputeInstanceGroup("data.google_compute_instance_group.test"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleComputeInstanceGroup_withNamedPort(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceGoogleComputeInstanceGroupConfigWithNamedPort(acctest.RandString(t, 10), acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceGoogleComputeInstanceGroup("data.google_compute_instance_group.test"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleComputeInstanceGroup_fromIGM(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceGoogleComputeInstanceGroup_fromIGM(fmt.Sprintf("tf-test-igm-%d", acctest.RandInt(t)), fmt.Sprintf("tf-test-igm-%d", acctest.RandInt(t))),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_compute_instance_group.test", "instances.#", "10"),
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
			return fmt.Errorf("cant' find data source called %s in state", dsFullName)
		}

		rs, ok := s.RootModule().Resources[rsFullName]
		if !ok {
			return fmt.Errorf("can't find resource called %s in state", rsFullName)
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
			"size",
		}

		for _, attrToTest := range attrsToTest {
			if dsAttrs[attrToTest] != rsAttrs[attrToTest] {
				return fmt.Errorf("%s is %s; want %s", attrToTest, dsAttrs[attrToTest], rsAttrs[attrToTest])
			}
		}

		if !tpgresource.CompareSelfLinkOrResourceName("", dsAttrs["self_link"], rsAttrs["self_link"], nil) && dsAttrs["self_link"] != rsAttrs["self_link"] {
			return fmt.Errorf("self link does not match: %s vs %s", dsAttrs["self_link"], rsAttrs["self_link"])
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

		for k, dsAttr := range dsInstancesValues {
			rsAttr := rsInstancesValues[k]
			if !tpgresource.CompareSelfLinkOrResourceName("", dsAttr, rsAttr, nil) && dsAttr != rsAttr {
				return fmt.Errorf("instance expected value %s did not match real value %s. expected list of instances %v, received %v", rsAttr, dsAttr, rsInstancesValues, dsInstancesValues)
			}
		}

		return nil
	}
}

func testAccCheckDataSourceGoogleComputeInstanceGroupConfig(instanceName, igName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "test" {
  name         = "tf-test-%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
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
  zone = google_compute_instance.test.zone

  instances = [
    google_compute_instance.test.self_link,
  ]
}

data "google_compute_instance_group" "test" {
  name = google_compute_instance_group.test.name
  zone = google_compute_instance_group.test.zone
}
`, instanceName, igName)
}

func testAccCheckDataSourceGoogleComputeInstanceGroupConfigWithNamedPort(instanceName, igName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "test" {
  name         = "tf-test-%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
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
  zone = google_compute_instance.test.zone

  named_port {
    name = "http"
    port = "8080"
  }

  named_port {
    name = "https"
    port = "8443"
  }

  instances = [
    google_compute_instance.test.self_link,
  ]
}

data "google_compute_instance_group" "test" {
  name = google_compute_instance_group.test.name
  zone = google_compute_instance_group.test.zone
}
`, instanceName, igName)
}

func testAccCheckDataSourceGoogleComputeInstanceGroup_fromIGM(igmName, secondIgmName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "igm-basic" {
  name         = "%s"
  machine_type = "e2-medium"

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }
}

resource "google_compute_instance_group_manager" "igm" {
  name              = "%s"
  version {
    instance_template = google_compute_instance_template.igm-basic.self_link
    name              = "primary"
  }
  base_instance_name = "igm"
  zone               = "us-central1-a"
  target_size        = 10

  wait_for_instances = true
}

data "google_compute_instance_group" "test" {
  self_link = google_compute_instance_group_manager.igm.instance_group
}
`, igmName, secondIgmName)
}
