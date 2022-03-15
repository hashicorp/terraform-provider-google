// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: DCL     ***
//
// ----------------------------------------------------------------------------
//
//     This file is managed by Magic Modules (https://github.com/GoogleCloudPlatform/magic-modules)
//     and is based on the DCL (https://github.com/GoogleCloudPlatform/declarative-resource-client-library).
//     Changes will need to be made to the DCL or Magic Modules instead of here.
//
//     We are not currently able to accept contributions to this file. If changes
//     are required, please file an issue at https://github.com/hashicorp/terraform-provider-google/issues/new/choose
//
// ----------------------------------------------------------------------------

package google

import (
	"context"
	"fmt"
	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	networkconnectivity "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/networkconnectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strings"
	"testing"
)

func TestAccNetworkConnectivitySpoke_RouterApplianceHandWritten(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  getTestProjectFromEnv(),
		"region":        getTestRegionFromEnv(),
		"zone":          getTestZoneFromEnv(),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkConnectivitySpokeDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConnectivitySpoke_RouterApplianceHandWritten(context),
			},
			{
				ResourceName:      "google_network_connectivity_spoke.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNetworkConnectivitySpoke_RouterApplianceHandWrittenUpdate0(context),
			},
			{
				ResourceName:      "google_network_connectivity_spoke.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNetworkConnectivitySpoke_RouterApplianceHandWritten(context map[string]interface{}) string {
	return Nprintf(`

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
	return Nprintf(`

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

func testAccCheckNetworkConnectivitySpokeDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "rs.google_network_connectivity_spoke" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := googleProviderConfig(t)

			billingProject := ""
			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			obj := &networkconnectivity.Spoke{
				Hub:         dcl.String(rs.Primary.Attributes["hub"]),
				Location:    dcl.String(rs.Primary.Attributes["location"]),
				Name:        dcl.String(rs.Primary.Attributes["name"]),
				Description: dcl.String(rs.Primary.Attributes["description"]),
				Project:     dcl.StringOrNil(rs.Primary.Attributes["project"]),
				CreateTime:  dcl.StringOrNil(rs.Primary.Attributes["create_time"]),
				State:       networkconnectivity.SpokeStateEnumRef(rs.Primary.Attributes["state"]),
				UniqueId:    dcl.StringOrNil(rs.Primary.Attributes["unique_id"]),
				UpdateTime:  dcl.StringOrNil(rs.Primary.Attributes["update_time"]),
			}

			client := NewDCLNetworkConnectivityClient(config, config.userAgent, billingProject, 0)
			_, err := client.GetSpoke(context.Background(), obj)
			if err == nil {
				return fmt.Errorf("google_network_connectivity_spoke still exists %v", obj)
			}
		}
		return nil
	}
}
