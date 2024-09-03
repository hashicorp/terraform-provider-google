// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package networkconnectivity_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccNetworkConnectivitySpoke_LinkedVPCNetworkHandWritten(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"region":        envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkConnectivitySpokeDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConnectivitySpoke_LinkedVPCNetworkHandWritten(context),
			},
			{
				ResourceName:            "google_network_connectivity_spoke.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccNetworkConnectivitySpoke_LinkedVPCNetworkHandWrittenUpdate0(context),
			},
			{
				ResourceName:            "google_network_connectivity_spoke.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}
func TestAccNetworkConnectivitySpoke_LinkedVPCNetworkHandWrittenLongForm(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"region":        envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkConnectivitySpokeDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConnectivitySpoke_LinkedVPCNetworkHandWrittenLongForm(context),
			},
			{
				ResourceName:            "google_network_connectivity_spoke.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccNetworkConnectivitySpoke_LinkedVPCNetworkHandWrittenUpdate0LongForm(context),
			},
			{
				ResourceName:            "google_network_connectivity_spoke.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func TestAccNetworkConnectivitySpoke_RouterApplianceHandWritten(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"region":        envvar.GetTestRegionFromEnv(),
		"zone":          envvar.GetTestZoneFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkConnectivitySpokeDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConnectivitySpoke_RouterApplianceHandWritten(context),
			},
			{
				ResourceName:            "google_network_connectivity_spoke.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccNetworkConnectivitySpoke_RouterApplianceHandWrittenUpdate0(context),
			},
			{
				ResourceName:            "google_network_connectivity_spoke.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}
func TestAccNetworkConnectivitySpoke_RouterApplianceHandWrittenLongForm(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"region":        envvar.GetTestRegionFromEnv(),
		"zone":          envvar.GetTestZoneFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkConnectivitySpokeDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConnectivitySpoke_RouterApplianceHandWrittenLongForm(context),
			},
			{
				ResourceName:            "google_network_connectivity_spoke.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccNetworkConnectivitySpoke_RouterApplianceHandWrittenUpdate0LongForm(context),
			},
			{
				ResourceName:            "google_network_connectivity_spoke.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}
func testAccNetworkConnectivitySpoke_LinkedVPCNetworkHandWritten(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_compute_network" "network" {
  name                    = "tf-test-network%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_network_connectivity_hub" "basic_hub" {
  name        = "tf-test-hub%{random_suffix}"
  description = "A sample hub"
  labels = {
    label-two = "value-one"
  }
}

resource "google_network_connectivity_spoke" "primary" {
  name = "tf-test-name%{random_suffix}"
  location = "global"
  description = "A sample spoke with a linked routher appliance instance"
  labels = {
    label-one = "value-one"
  }
  hub = google_network_connectivity_hub.basic_hub.id
  linked_vpc_network {
    exclude_export_ranges = [
      "198.51.100.0/24",
      "10.10.0.0/16"
    ]
    include_export_ranges = [
      "198.51.100.0/23", 
      "10.0.0.0/8"
    ]
    uri = google_compute_network.network.self_link
  }
}
`, context)
}

func testAccNetworkConnectivitySpoke_LinkedVPCNetworkHandWrittenUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_compute_network" "network" {
  name                    = "tf-test-network%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_network_connectivity_hub" "basic_hub" {
  name        = "tf-test-hub%{random_suffix}"
  description = "A sample hub"
  labels = {
    label-two = "value-one"
  }
}

resource "google_network_connectivity_spoke" "primary" {
  name = "tf-test-name%{random_suffix}"
  location = "global"
  description = "A sample spoke with a linked routher appliance instance"
  labels = {
    label-one = "value-one"
  }
  hub = google_network_connectivity_hub.basic_hub.id
  linked_vpc_network {
    exclude_export_ranges = [
      "198.51.100.0/24",
      "10.10.0.0/16"
    ]
    include_export_ranges = [
      "198.51.100.0/23", 
      "10.0.0.0/8"
    ]
    uri = google_compute_network.network.self_link
  }
}
`, context)
}

func testAccNetworkConnectivitySpoke_RouterApplianceHandWritten(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_compute_network" "network" {
  name                    = "tf-test-network%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnetwork" {
  name          = "tf-test-subnet%{random_suffix}"
  ip_cidr_range = "10.0.0.0/28"
  region        = "%{region}"
  network       = google_compute_network.network.self_link
}

resource "google_compute_instance" "instance" {
  name         = "tf-test-instance%{random_suffix}"
  machine_type = "e2-medium"
  can_ip_forward = true
  zone         = "%{zone}"

  boot_disk {
    initialize_params {
      image = "projects/debian-cloud/global/images/debian-10-buster-v20210817"
    }
  }

  network_interface {
    subnetwork = google_compute_subnetwork.subnetwork.name
    network_ip = "10.0.0.2"
    access_config {
        network_tier = "PREMIUM"
    }
  }
}

resource "google_network_connectivity_hub" "basic_hub" {
  name        = "tf-test-hub%{random_suffix}"
  description = "A sample hub"
  labels = {
    label-two = "value-one"
  }
}

resource "google_network_connectivity_spoke" "primary" {
  name = "tf-test-name%{random_suffix}"
  location = "%{region}"
  description = "A sample spoke with a linked routher appliance instance"
  labels = {
    label-one = "value-one"
  }
  hub =  google_network_connectivity_hub.basic_hub.id
  linked_router_appliance_instances {
    instances {
        virtual_machine = google_compute_instance.instance.self_link
        ip_address = "10.0.0.2"
    }
    site_to_site_data_transfer = true
  }
}
`, context)
}

func testAccNetworkConnectivitySpoke_RouterApplianceHandWrittenUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_compute_network" "network" {
  name                    = "tf-test-network%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnetwork" {
  name          = "tf-test-subnet%{random_suffix}"
  ip_cidr_range = "10.0.0.0/28"
  region        = "%{region}"
  network       = google_compute_network.network.self_link
}

resource "google_compute_instance" "instance" {
  name         = "tf-test-instance%{random_suffix}"
  machine_type = "e2-medium"
  can_ip_forward = true
  zone         = "%{zone}"

  boot_disk {
    initialize_params {
      image = "projects/debian-cloud/global/images/debian-10-buster-v20210817"
    }
  }

  network_interface {
    subnetwork = google_compute_subnetwork.subnetwork.name
    network_ip = "10.0.0.2"
    access_config {
        network_tier = "PREMIUM"
    }
  }
}

resource "google_network_connectivity_hub" "basic_hub" {
  name        = "tf-test-hub%{random_suffix}"
  description = "A sample hub"
  labels = {
    label-two = "value-one"
  }
}

resource "google_network_connectivity_spoke" "primary" {
  name = "tf-test-name%{random_suffix}"
  location = "%{region}"
  description = "An UPDATED sample spoke with a linked routher appliance instance"
  labels = {
    label-two = "value-two"
  }
  hub = google_network_connectivity_hub.basic_hub.id
  linked_router_appliance_instances {
    instances {
        virtual_machine = google_compute_instance.instance.self_link
        ip_address = "10.0.0.2"
    }
    site_to_site_data_transfer = true
  }
}
`, context)
}
func testAccNetworkConnectivitySpoke_LinkedVPCNetworkHandWrittenLongForm(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_compute_network" "network" {
  name                    = "tf-test-network%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_network_connectivity_hub" "basic_hub" {
  name        = "tf-test-hub%{random_suffix}"
  description = "A sample hub"
  labels = {
    label-two = "value-one"
  }
}

resource "google_network_connectivity_spoke" "primary" {
  name = "tf-test-name%{random_suffix}"
  location = "global"
  description = "A sample spoke with a linked routher appliance instance"
  labels = {
    label-one = "value-one"
  }
  hub = google_network_connectivity_hub.basic_hub.id
  linked_vpc_network {
    exclude_export_ranges = [
      "198.51.100.0/24",
      "10.10.0.0/16"
    ]
    include_export_ranges = [
      "198.51.100.0/23", 
      "10.0.0.0/8"
    ]
    uri = google_compute_network.network.self_link
  }
}
`, context)
}

func testAccNetworkConnectivitySpoke_LinkedVPCNetworkHandWrittenUpdate0LongForm(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_compute_network" "network" {
  name                    = "tf-test-network%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_network_connectivity_hub" "basic_hub" {
  name        = "tf-test-hub%{random_suffix}"
  description = "A sample hub"
  labels = {
    label-two = "value-one"
  }
}

resource "google_network_connectivity_spoke" "primary" {
  name = "tf-test-name%{random_suffix}"
  location = "global"
  description = "A sample spoke with a linked routher appliance instance"
  labels = {
    label-one = "value-one"
  }
  hub = google_network_connectivity_hub.basic_hub.id
  linked_vpc_network {
    exclude_export_ranges = [
      "198.51.100.0/24",
      "10.10.0.0/16"
    ]
    include_export_ranges = [
      "198.51.100.0/23", 
      "10.0.0.0/8"
    ]
    uri = google_compute_network.network.self_link
  }
}
`, context)
}

func testAccNetworkConnectivitySpoke_RouterApplianceHandWrittenLongForm(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_compute_network" "network" {
  name                    = "tf-test-network%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnetwork" {
  name          = "tf-test-subnet%{random_suffix}"
  ip_cidr_range = "10.0.0.0/28"
  region        = "%{region}"
  network       = google_compute_network.network.self_link
}

resource "google_compute_instance" "instance" {
  name         = "tf-test-instance%{random_suffix}"
  machine_type = "e2-medium"
  can_ip_forward = true
  zone         = "%{zone}"

  boot_disk {
    initialize_params {
      image = "projects/debian-cloud/global/images/debian-10-buster-v20210817"
    }
  }

  network_interface {
    subnetwork = google_compute_subnetwork.subnetwork.name
    network_ip = "10.0.0.2"
    access_config {
        network_tier = "PREMIUM"
    }
  }
}

resource "google_network_connectivity_hub" "basic_hub" {
  name        = "tf-test-hub%{random_suffix}"
  description = "A sample hub"
  labels = {
    label-two = "value-one"
  }
}

resource "google_network_connectivity_spoke" "primary" {
  name = "tf-test-name%{random_suffix}"
  location = "%{region}"
  description = "A sample spoke with a linked routher appliance instance"
  labels = {
    label-one = "value-one"
  }
  hub =  google_network_connectivity_hub.basic_hub.id
  linked_router_appliance_instances {
    instances {
        virtual_machine = google_compute_instance.instance.self_link
        ip_address = "10.0.0.2"
    }
    site_to_site_data_transfer = true
  }
}
`, context)
}

func testAccNetworkConnectivitySpoke_RouterApplianceHandWrittenUpdate0LongForm(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_compute_network" "network" {
  name                    = "tf-test-network%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnetwork" {
  name          = "tf-test-subnet%{random_suffix}"
  ip_cidr_range = "10.0.0.0/28"
  region        = "%{region}"
  network       = google_compute_network.network.self_link
}

resource "google_compute_instance" "instance" {
  name         = "tf-test-instance%{random_suffix}"
  machine_type = "e2-medium"
  can_ip_forward = true
  zone         = "%{zone}"

  boot_disk {
    initialize_params {
      image = "projects/debian-cloud/global/images/debian-10-buster-v20210817"
    }
  }

  network_interface {
    subnetwork = google_compute_subnetwork.subnetwork.name
    network_ip = "10.0.0.2"
    access_config {
        network_tier = "PREMIUM"
    }
  }
}

resource "google_network_connectivity_hub" "basic_hub" {
  name        = "tf-test-hub%{random_suffix}"
  description = "A sample hub"
  labels = {
    label-two = "value-one"
  }
}

resource "google_network_connectivity_spoke" "primary" {
  name = "tf-test-name%{random_suffix}"
  location = "%{region}"
  description = "An UPDATED sample spoke with a linked routher appliance instance"
  labels = {
    label-two = "value-two"
  }
  hub = google_network_connectivity_hub.basic_hub.id
  linked_router_appliance_instances {
    instances {
        virtual_machine = google_compute_instance.instance.self_link
        ip_address = "10.0.0.2"
    }
    site_to_site_data_transfer = true
  }
}
`, context)
}
