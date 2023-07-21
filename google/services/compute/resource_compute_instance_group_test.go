// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"google.golang.org/api/compute/v1"
)

func TestAccComputeInstanceGroup_basic(t *testing.T) {
	t.Parallel()

	var instanceGroup compute.InstanceGroup
	var resourceName = "google_compute_instance_group.basic"
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var zone = "us-central1-c"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComputeInstanceGroup_destroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceGroup_basic(zone, instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccComputeInstanceGroup_exists(
						t, "google_compute_instance_group.basic", &instanceGroup),
					testAccComputeInstanceGroup_exists(
						t, "google_compute_instance_group.empty", &instanceGroup),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("%s/%s/%s", envvar.GetTestProjectFromEnv(), zone, instanceName),
			},
		},
	})
}

func TestAccComputeInstanceGroup_rename(t *testing.T) {
	t.Parallel()

	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var instanceGroupName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var backendName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var healthName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComputeInstanceGroup_destroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceGroup_rename(instanceName, instanceGroupName, backendName, healthName),
			},
			{
				ResourceName:      "google_compute_instance_group.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeInstanceGroup_rename(instanceName, instanceGroupName+"2", backendName, healthName),
			},
			{
				ResourceName:      "google_compute_instance_group.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceGroup_update(t *testing.T) {
	t.Parallel()

	var instanceGroup compute.InstanceGroup
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComputeInstanceGroup_destroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceGroup_update(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccComputeInstanceGroup_exists(
						t, "google_compute_instance_group.update", &instanceGroup),
					testAccComputeInstanceGroup_named_ports(
						t,
						"google_compute_instance_group.update",
						map[string]int64{"http": 8080, "https": 8443},
						&instanceGroup),
				),
			},
			{
				Config: testAccComputeInstanceGroup_update2(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccComputeInstanceGroup_exists(
						t, "google_compute_instance_group.update", &instanceGroup),
					testAccComputeInstanceGroup_updated(
						t, "google_compute_instance_group.update", 1, &instanceGroup),
					testAccComputeInstanceGroup_named_ports(
						t,
						"google_compute_instance_group.update",
						map[string]int64{"http": 8081, "test": 8444},
						&instanceGroup),
				),
			},
		},
	})
}

func TestAccComputeInstanceGroup_outOfOrderInstances(t *testing.T) {
	t.Parallel()

	var instanceGroup compute.InstanceGroup
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComputeInstanceGroup_destroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceGroup_outOfOrderInstances(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccComputeInstanceGroup_exists(
						t, "google_compute_instance_group.group", &instanceGroup),
				),
			},
		},
	})
}

func TestAccComputeInstanceGroup_network(t *testing.T) {
	t.Parallel()

	var instanceGroup compute.InstanceGroup
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComputeInstanceGroup_destroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceGroup_network(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccComputeInstanceGroup_exists(
						t, "google_compute_instance_group.with_instance", &instanceGroup),
					testAccComputeInstanceGroup_hasCorrectNetwork(
						t, "google_compute_instance_group.with_instance", "google_compute_network.ig_network", &instanceGroup),
					testAccComputeInstanceGroup_exists(
						t, "google_compute_instance_group.without_instance", &instanceGroup),
					testAccComputeInstanceGroup_hasCorrectNetwork(
						t, "google_compute_instance_group.without_instance", "google_compute_network.ig_network", &instanceGroup),
				),
			},
		},
	})
}

func testAccComputeInstanceGroup_destroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_instance_group" {
				continue
			}
			_, err := config.NewComputeClient(config.UserAgent).InstanceGroups.Get(
				config.Project, rs.Primary.Attributes["zone"], rs.Primary.Attributes["name"]).Do()
			if err == nil {
				return fmt.Errorf("InstanceGroup still exists")
			}
		}

		return nil
	}
}

func testAccComputeInstanceGroup_exists(t *testing.T, n string, instanceGroup *compute.InstanceGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acctest.GoogleProviderConfig(t)

		found, err := config.NewComputeClient(config.UserAgent).InstanceGroups.Get(
			config.Project, rs.Primary.Attributes["zone"], rs.Primary.Attributes["name"]).Do()
		if err != nil {
			return err
		}

		*instanceGroup = *found

		return nil
	}
}

func testAccComputeInstanceGroup_updated(t *testing.T, n string, size int64, instanceGroup *compute.InstanceGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acctest.GoogleProviderConfig(t)

		instanceGroup, err := config.NewComputeClient(config.UserAgent).InstanceGroups.Get(
			config.Project, rs.Primary.Attributes["zone"], rs.Primary.Attributes["name"]).Do()
		if err != nil {
			return err
		}

		// Cannot check the target pool as the instance creation is asynchronous.  However, can
		// check the target_size.
		if instanceGroup.Size != size {
			return fmt.Errorf("instance count incorrect. saw real value %v instead of expected value %v", instanceGroup.Size, size)
		}

		return nil
	}
}

func testAccComputeInstanceGroup_named_ports(t *testing.T, n string, np map[string]int64, instanceGroup *compute.InstanceGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acctest.GoogleProviderConfig(t)

		instanceGroup, err := config.NewComputeClient(config.UserAgent).InstanceGroups.Get(
			config.Project, rs.Primary.Attributes["zone"], rs.Primary.Attributes["name"]).Do()
		if err != nil {
			return err
		}

		var found bool
		for _, namedPort := range instanceGroup.NamedPorts {
			found = false
			for name, port := range np {
				if namedPort.Name == name && namedPort.Port == port {
					found = true
				}
			}
			if !found {
				return fmt.Errorf("named port incorrect")
			}
		}

		return nil
	}
}

func testAccComputeInstanceGroup_hasCorrectNetwork(t *testing.T, nInstanceGroup string, nNetwork string, instanceGroup *compute.InstanceGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		rsInstanceGroup, ok := s.RootModule().Resources[nInstanceGroup]
		if !ok {
			return fmt.Errorf("Not found: %s", nInstanceGroup)
		}
		if rsInstanceGroup.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		instanceGroup, err := config.NewComputeClient(config.UserAgent).InstanceGroups.Get(
			config.Project, rsInstanceGroup.Primary.Attributes["zone"], rsInstanceGroup.Primary.Attributes["name"]).Do()
		if err != nil {
			return err
		}

		rsNetwork, ok := s.RootModule().Resources[nNetwork]
		if !ok {
			return fmt.Errorf("Not found: %s", nNetwork)
		}
		if rsNetwork.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		network, err := config.NewComputeClient(config.UserAgent).Networks.Get(
			config.Project, rsNetwork.Primary.Attributes["name"]).Do()
		if err != nil {
			return err
		}

		if instanceGroup.Network != network.SelfLink {
			return fmt.Errorf("network incorrect: actual=%s vs expected=%s", instanceGroup.Network, network.SelfLink)
		}

		return nil
	}
}

func testAccComputeInstanceGroup_basic(zone, instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "ig_instance" {
  name           = "%s"
  machine_type   = "e2-medium"
  can_ip_forward = false
  zone           = "us-central1-c"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }
}

resource "google_compute_instance_group" "basic" {
  description = "Terraform test instance group"
  name        = "%s"
  zone        = "%s"
  instances   = [google_compute_instance.ig_instance.id]
  named_port {
    name = "http"
    port = "8080"
  }
  named_port {
    name = "https"
    port = "8443"
  }
}

resource "google_compute_instance_group" "empty" {
  description = "Terraform test instance group empty"
  name        = "%s-empty"
  zone        = "%s"
  named_port {
    name = "http"
    port = "8080"
  }
  named_port {
    name = "https"
    port = "8443"
  }
}
`, instance, instance, zone, instance, zone)
}

func testAccComputeInstanceGroup_rename(instance, instanceGroup, backend, health string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "ig_instance" {
  name           = "%s"
  machine_type   = "e2-medium"
  can_ip_forward = false
  zone           = "us-central1-c"
  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }
}

resource "google_compute_instance_group" "basic" {
  name      = "%s"
  zone      = "us-central1-c"
  instances = [google_compute_instance.ig_instance.self_link]
  named_port {
    name = "http"
    port = "8080"
  }

  named_port {
    name = "https"
    port = "8443"
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "google_compute_backend_service" "default_backend" {
  name      = "%s"
  port_name = "https"
  protocol  = "HTTPS"

  backend {
    group = google_compute_instance_group.basic.self_link
  }

  health_checks = [
    google_compute_https_health_check.healthcheck.self_link,
  ]
}

resource "google_compute_https_health_check" "healthcheck" {
  name         = "%s"
  request_path = "/health_check"
}
`, instance, instanceGroup, backend, health)
}

func testAccComputeInstanceGroup_update(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "ig_instance" {
  name           = "%s-${count.index}"
  machine_type   = "e2-medium"
  can_ip_forward = false
  zone           = "us-central1-c"
  count          = 2

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }
}

resource "google_compute_instance_group" "update" {
  description = "Terraform test instance group"
  name        = "%s"
  zone        = "us-central1-c"
  instances   = google_compute_instance.ig_instance.*.self_link
  named_port {
    name = "http"
    port = "8080"
  }
  named_port {
    name = "https"
    port = "8443"
  }
}
`, instance, instance)
}

// Change IGM's instance template and target size
func testAccComputeInstanceGroup_update2(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "ig_instance" {
  name           = "%s-${count.index}"
  machine_type   = "e2-medium"
  can_ip_forward = false
  zone           = "us-central1-c"
  count          = 1

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }
}

resource "google_compute_instance_group" "update" {
  description = "Terraform test instance group"
  name        = "%s"
  zone        = "us-central1-c"
  instances   = google_compute_instance.ig_instance.*.self_link

  named_port {
    name = "http"
    port = "8081"
  }
  named_port {
    name = "test"
    port = "8444"
  }
}
`, instance, instance)
}

func testAccComputeInstanceGroup_outOfOrderInstances(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "ig_instance" {
  name           = "%s-1"
  machine_type   = "e2-medium"
  can_ip_forward = false
  zone           = "us-central1-c"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }
}

resource "google_compute_instance" "ig_instance_2" {
  name           = "%s-2"
  machine_type   = "e2-medium"
  can_ip_forward = false
  zone           = "us-central1-c"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }
}

resource "google_compute_instance_group" "group" {
  description = "Terraform test instance group"
  name        = "%s"
  zone        = "us-central1-c"
  instances   = [google_compute_instance.ig_instance_2.self_link, google_compute_instance.ig_instance.self_link]
  named_port {
    name = "http"
    port = "8080"
  }
  named_port {
    name = "https"
    port = "8443"
  }
}
`, instance, instance, instance)
}

func testAccComputeInstanceGroup_network(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_network" "ig_network" {
  name                    = "%[1]s"
  auto_create_subnetworks = true
}

resource "google_compute_instance" "ig_instance" {
  name           = "%[1]s"
  machine_type   = "e2-medium"
  can_ip_forward = false
  zone           = "us-central1-c"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = google_compute_network.ig_network.name
  }
}

resource "google_compute_instance_group" "with_instance" {
  description = "Terraform test instance group"
  name        = "%[1]s-with-instance"
  zone        = "us-central1-c"
  instances   = [google_compute_instance.ig_instance.self_link]
}

resource "google_compute_instance_group" "without_instance" {
  description = "Terraform test instance group"
  name        = "%[1]s-without-instance"
  zone        = "us-central1-c"
  network     = google_compute_network.ig_network.self_link
}
`, instance)
}
