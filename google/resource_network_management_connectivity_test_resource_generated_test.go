// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccNetworkManagementConnectivityTest_networkManagementConnectivityTestInstancesExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkManagementConnectivityTestDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkManagementConnectivityTest_networkManagementConnectivityTestInstancesExample(context),
			},
			{
				ResourceName:      "google_network_management_connectivity_test.instance-test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNetworkManagementConnectivityTest_networkManagementConnectivityTestInstancesExample(context map[string]interface{}) string {
	return Nprintf(`
resource "google_network_management_connectivity_test" "instance-test" {
  name = "tf-test-conn-test-instances%{random_suffix}"
  source {
    instance = google_compute_instance.source.id
  }

  destination {
    instance = google_compute_instance.destination.id
  }

  protocol = "TCP"
}

resource "google_compute_instance" "source" {
  name = "tf-test-source-vm%{random_suffix}"
  machine_type = "e2-medium"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.debian_9.id
    }
  }

  network_interface {
    network = google_compute_network.vpc.id
    access_config {
    }
  }
}

resource "google_compute_instance" "destination" {
  name = "tf-test-dest-vm%{random_suffix}"
  machine_type = "e2-medium"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.debian_9.id
    }
  }

  network_interface {
    network = google_compute_network.vpc.id
    access_config {
    }
  }
}

resource "google_compute_network" "vpc" {
  name = "tf-test-conn-test-net%{random_suffix}"
}

data "google_compute_image" "debian_9" {
  family  = "debian-11"
  project = "debian-cloud"
}
`, context)
}

func TestAccNetworkManagementConnectivityTest_networkManagementConnectivityTestAddressesExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkManagementConnectivityTestDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkManagementConnectivityTest_networkManagementConnectivityTestAddressesExample(context),
			},
			{
				ResourceName:      "google_network_management_connectivity_test.address-test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNetworkManagementConnectivityTest_networkManagementConnectivityTestAddressesExample(context map[string]interface{}) string {
	return Nprintf(`
resource "google_network_management_connectivity_test" "address-test" {
  name = "tf-test-conn-test-addr%{random_suffix}"
  source {
      ip_address = google_compute_address.source-addr.address
      project_id = google_compute_address.source-addr.project
      network = google_compute_network.vpc.id
      network_type = "GCP_NETWORK"
  }

  destination {
      ip_address = google_compute_address.dest-addr.address
      project_id = google_compute_address.dest-addr.project
      network = google_compute_network.vpc.id
  }

  protocol = "UDP"
}

resource "google_compute_network" "vpc" {
  name = "tf-test-connectivity-vpc%{random_suffix}"
}

resource "google_compute_subnetwork" "subnet" {
  name          = "tf-test-connectivity-vpc%{random_suffix}-subnet"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.vpc.id
}

resource "google_compute_address" "source-addr" {
  name         = "tf-test-src-addr%{random_suffix}"
  subnetwork   = google_compute_subnetwork.subnet.id
  address_type = "INTERNAL"
  address      = "10.0.42.42"
  region       = "us-central1"
}

resource "google_compute_address" "dest-addr" {
  name         = "tf-test-dest-addr%{random_suffix}"
  subnetwork   = google_compute_subnetwork.subnet.id
  address_type = "INTERNAL"
  address      = "10.0.43.43"
  region       = "us-central1"
}
`, context)
}

func testAccCheckNetworkManagementConnectivityTestDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_network_management_connectivity_test" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{NetworkManagementBasePath}}projects/{{project}}/locations/global/connectivityTests/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("NetworkManagementConnectivityTest still exists at %s", url)
			}
		}

		return nil
	}
}
