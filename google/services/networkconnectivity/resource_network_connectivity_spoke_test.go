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
			{
				Config: testAccNetworkConnectivitySpoke_RouterApplianceHandWrittenUpdate1(context),
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
			{
				Config: testAccNetworkConnectivitySpoke_RouterApplianceHandWrittenUpdate1LongForm(context),
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

func TestAccNetworkConnectivitySpoke_VPNTunnelHandWrittenHandWritten(t *testing.T) {
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
				Config: testAccNetworkConnectivitySpoke_VPNTunnelHandWrittenHandWritten(context),
			},
			{
				ResourceName:            "google_network_connectivity_spoke.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccNetworkConnectivitySpoke_VPNTunnelHandWrittenHandWrittenUpdate0(context),
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

func TestAccNetworkConnectivitySpoke_InterconnectAttachmentHandWrittenHandWritten(t *testing.T) {
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
				Config: testAccNetworkConnectivitySpoke_InterconnectAttachmentHandWrittenHandWritten(context),
			},
			{
				ResourceName:            "google_network_connectivity_spoke.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccNetworkConnectivitySpoke_InterconnectAttachmentHandWrittenHandWrittenUpdate0(context),
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

resource "google_compute_instance" "router-instance1" {
  name         = "tf-test-router-instance1%{random_suffix}"
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
  description = "A sample spoke with a single linked routher appliance instance"
  labels = {
    label-one = "value-one"
  }
  hub =  google_network_connectivity_hub.basic_hub.id
  linked_router_appliance_instances {
    instances {
        virtual_machine = google_compute_instance.router-instance1.self_link
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

resource "google_compute_instance" "router-instance1" {
  name         = "tf-test-router-instance1%{random_suffix}"
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
  description = "An UPDATED sample spoke with a single linked routher appliance instance"
  labels = {
    label-two = "value-two"
  }
  hub = google_network_connectivity_hub.basic_hub.id
  linked_router_appliance_instances {
    instances {
        virtual_machine = google_compute_instance.router-instance1.self_link
        ip_address = "10.0.0.2"
    }
    site_to_site_data_transfer = true
  }
}
`, context)
}

func testAccNetworkConnectivitySpoke_RouterApplianceHandWrittenUpdate1(context map[string]interface{}) string {
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

resource "google_compute_instance" "router-instance1" {
  name         = "tf-test-router-instance1%{random_suffix}"
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

resource "google_compute_instance" "router-instance2" {
  name         = "tf-test-router-instance2%{random_suffix}"
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
    network_ip = "10.0.0.3"
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
  description = "An UPDATED sample spoke with two linked routher appliance instances"
  labels = {
    label-two = "value-two"
  }
  hub = google_network_connectivity_hub.basic_hub.id
  linked_router_appliance_instances {
    instances {
        virtual_machine = google_compute_instance.router-instance1.self_link
        ip_address = "10.0.0.2"
    }
    instances {
        virtual_machine = google_compute_instance.router-instance2.self_link
        ip_address = "10.0.0.3"
    }
    include_import_ranges = ["ALL_IPV4_RANGES"]
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

resource "google_compute_instance" "router-instance1" {
  name         = "tf-test-router-instance1%{random_suffix}"
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
  description = "A sample spoke with a single linked routher appliance instance"
  labels = {
    label-one = "value-one"
  }
  hub =  google_network_connectivity_hub.basic_hub.id
  linked_router_appliance_instances {
    instances {
        virtual_machine = google_compute_instance.router-instance1.self_link
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

resource "google_compute_instance" "router-instance1" {
  name         = "tf-test-router-instance1%{random_suffix}"
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
  description = "An UPDATED sample spoke with a single linked routher appliance instance"
  labels = {
    label-two = "value-two"
  }
  hub = google_network_connectivity_hub.basic_hub.id
  linked_router_appliance_instances {
    instances {
        virtual_machine = google_compute_instance.router-instance1.self_link
        ip_address = "10.0.0.2"
    }
    site_to_site_data_transfer = true
  }
}
`, context)
}

func testAccNetworkConnectivitySpoke_RouterApplianceHandWrittenUpdate1LongForm(context map[string]interface{}) string {
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

resource "google_compute_instance" "router-instance1" {
  name         = "tf-test-router-instance1%{random_suffix}"
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

resource "google_compute_instance" "router-instance2" {
  name         = "tf-test-router-instance2%{random_suffix}"
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
    network_ip = "10.0.0.3"
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
  description = "An UPDATED sample spoke with two linked routher appliance instances"
  labels = {
    label-two = "value-two"
  }
  hub = google_network_connectivity_hub.basic_hub.id
  linked_router_appliance_instances {
    instances {
        virtual_machine = google_compute_instance.router-instance1.self_link
        ip_address = "10.0.0.2"
    }
    instances {
        virtual_machine = google_compute_instance.router-instance2.self_link
        ip_address = "10.0.0.3"
    }
    include_import_ranges = ["ALL_IPV4_RANGES"]
    site_to_site_data_transfer = true
  }
}
`, context)
}

func testAccNetworkConnectivitySpoke_VPNTunnelHandWrittenHandWritten(context map[string]interface{}) string {
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

resource "google_compute_ha_vpn_gateway" "gateway" {
  name    = "tf-test-gw%{random_suffix}"
  network = google_compute_network.network.id
}

resource "google_compute_external_vpn_gateway" "external_vpn_gw" {
  name            = "tf-test-external-gw%{random_suffix}"
  redundancy_type = "SINGLE_IP_INTERNALLY_REDUNDANT"
  description     = "An externally managed VPN gateway"
  interface {
    id         = 0
    ip_address = "8.8.8.8"
  }
}

resource "google_compute_router" "router" {
  name    = "tf-test-router%{random_suffix}"
  region  = "%{region}"
  network = google_compute_network.network.name
  bgp {
    asn = 64514
  }
}

resource "google_compute_vpn_tunnel" "tunnel" {
  name                            = "tf-test-tunnel%{random_suffix}"
  region                          = "%{region}"
  vpn_gateway                     = google_compute_ha_vpn_gateway.gateway.id
  peer_external_gateway           = google_compute_external_vpn_gateway.external_vpn_gw.id
  peer_external_gateway_interface = 0
  shared_secret                   = "a secret message"
  router                          = google_compute_router.router.id
  vpn_gateway_interface           = 0
}

resource "google_compute_router_interface" "router_interface" {
  name       = "tf-test-ri%{random_suffix}"
  router     = google_compute_router.router.name
  region     = "%{region}"
  ip_range   = "169.254.0.1/30"
  vpn_tunnel = google_compute_vpn_tunnel.tunnel.name
}

resource "google_compute_router_peer" "router_peer" {
  name                      = "tf-test-peer%{random_suffix}"
  router                    = google_compute_router.router.name
  region                    = "%{region}"
  peer_ip_address           = "169.254.0.2"
  peer_asn                  = 64515
  advertised_route_priority = 100
  interface                 = google_compute_router_interface.router_interface.name
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
  description = "A sample spoke with a linked VPN Tunnel, no include_import_ranges yet"
  labels = {
    label-one = "value-one"
  }
  hub = google_network_connectivity_hub.basic_hub.id
  linked_vpn_tunnels {
    uris                       = [google_compute_vpn_tunnel.tunnel.self_link]
    site_to_site_data_transfer = true
  }
}
`, context)
}

func testAccNetworkConnectivitySpoke_VPNTunnelHandWrittenHandWrittenUpdate0(context map[string]interface{}) string {
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

resource "google_compute_ha_vpn_gateway" "gateway" {
  name    = "tf-test-gw%{random_suffix}"
  network = google_compute_network.network.id
}

resource "google_compute_external_vpn_gateway" "external_vpn_gw" {
  name            = "tf-test-external-gw%{random_suffix}"
  redundancy_type = "SINGLE_IP_INTERNALLY_REDUNDANT"
  description     = "An externally managed VPN gateway"
  interface {
    id         = 0
    ip_address = "8.8.8.8"
  }
}

resource "google_compute_router" "router" {
  name    = "tf-test-router%{random_suffix}"
  region  = "%{region}"
  network = google_compute_network.network.name
  bgp {
    asn = 64514
  }
}

resource "google_compute_vpn_tunnel" "tunnel" {
  name                            = "tf-test-tunnel%{random_suffix}"
  region                          = "%{region}"
  vpn_gateway                     = google_compute_ha_vpn_gateway.gateway.id
  peer_external_gateway           = google_compute_external_vpn_gateway.external_vpn_gw.id
  peer_external_gateway_interface = 0
  shared_secret                   = "a secret message"
  router                          = google_compute_router.router.id
  vpn_gateway_interface           = 0
}

resource "google_compute_router_interface" "router_interface" {
  name       = "tf-test-ri%{random_suffix}"
  router     = google_compute_router.router.name
  region     = "%{region}"
  ip_range   = "169.254.0.1/30"
  vpn_tunnel = google_compute_vpn_tunnel.tunnel.name
}

resource "google_compute_router_peer" "router_peer" {
  name                      = "tf-test-peer%{random_suffix}"
  router                    = google_compute_router.router.name
  region                    = "%{region}"
  peer_ip_address           = "169.254.0.2"
  peer_asn                  = 64515
  advertised_route_priority = 100
  interface                 = google_compute_router_interface.router_interface.name
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
  description = "An UPDATED sample spoke with a linked VPN Tunnel, now includes ALL_IPV4_RANGES"
  labels = {
    label-one = "value-one"
  }
  hub = google_network_connectivity_hub.basic_hub.id
  linked_vpn_tunnels {
    uris                       = [google_compute_vpn_tunnel.tunnel.self_link]
    site_to_site_data_transfer = true
    include_import_ranges = ["ALL_IPV4_RANGES"]
  }
}
`, context)
}

func testAccNetworkConnectivitySpoke_InterconnectAttachmentHandWrittenHandWritten(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_network_connectivity_hub" "basic_hub" {
  name        = "tf-test-hub%{random_suffix}"
  description = "A sample hub"
  labels = {
    label-two = "value-one"
  }
}

resource "google_compute_network" "network" {
  name                    = "tf-test-network%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_router" "router" {
  name    = "tf-test-router%{random_suffix}"
  region  = "%{region}"
  network = google_compute_network.network.name
  bgp {
    asn = 16550
  }
}

resource "google_compute_interconnect_attachment" "interconnect_attachment" {
  name                     = "tf-test-ia%{random_suffix}"
  edge_availability_domain = "AVAILABILITY_DOMAIN_1"
  type                     = "PARTNER"
  router                   = google_compute_router.router.id
  mtu                      = 1500
  region                   = "%{region}"
}

resource "google_network_connectivity_spoke" "primary" {
  name        = "tf-test-spoke-ia%{random_suffix}"
  location    = "%{region}"
  description = "A sample spoke with a linked interconnect_attachment, no include_import_ranges yet"
  labels = {
    label-one = "value-one"
  }
  hub = google_network_connectivity_hub.basic_hub.id
  linked_interconnect_attachments {
    uris                       = [google_compute_interconnect_attachment.interconnect_attachment.self_link]
    site_to_site_data_transfer = true
    # include_import_ranges not set initially
  }
}
`, context)
}

func testAccNetworkConnectivitySpoke_InterconnectAttachmentHandWrittenHandWrittenUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_network_connectivity_hub" "basic_hub" {
  name        = "tf-test-hub%{random_suffix}"
  description = "A sample hub"
  labels = {
    label-two = "value-one"
  }
}

resource "google_compute_network" "network" {
  name                    = "tf-test-network%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_router" "router" {
  name    = "tf-test-router%{random_suffix}"
  region  = "%{region}"
  network = google_compute_network.network.name
  bgp {
    asn = 16550
  }
}

resource "google_compute_interconnect_attachment" "interconnect_attachment" {
  name                     = "tf-test-ia%{random_suffix}"
  edge_availability_domain = "AVAILABILITY_DOMAIN_1"
  type                     = "PARTNER"
  router                   = google_compute_router.router.id
  mtu                      = 1500
  region                   = "%{region}"
}

resource "google_network_connectivity_spoke" "primary" {
  name        = "tf-test-spoke-ia%{random_suffix}"
  location    = "%{region}"
  description = "An updated sample spoke with interconnect_attachment, now includes ALL_IPV4_RANGES"
  labels = {
    label-one = "value-one"
  }
  hub = google_network_connectivity_hub.basic_hub.id
  linked_interconnect_attachments {
    uris                       = [google_compute_interconnect_attachment.interconnect_attachment.self_link]
    site_to_site_data_transfer = true
    include_import_ranges      = ["ALL_IPV4_RANGES"]
  }
}
`, context)
}
