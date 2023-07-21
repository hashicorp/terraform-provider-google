// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccComputeNetworkEndpoint_networkEndpointsBasic(t *testing.T) {
	// Multiple fine-grained resources
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"default_port":  90,
		"modified_port": 100,
		"add1_port":     101,
		"add2_port":     102,
	}
	negId := fmt.Sprintf("projects/%s/zones/%s/networkEndpointGroups/tf-test-neg-%s",
		envvar.GetTestProjectFromEnv(), envvar.GetTestZoneFromEnv(), context["random_suffix"])

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Create one endpoint
				Config: testAccComputeNetworkEndpoint_networkEndpointsBasic(context),
			},
			{
				ResourceName:      "google_compute_network_endpoint.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Force-recreate old endpoint
				Config: testAccComputeNetworkEndpoint_networkEndpointsModified(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeNetworkEndpointWithPortsDestroyed(t, negId, "90"),
				),
			},
			{
				ResourceName:      "google_compute_network_endpoint.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Add two new endpoints
				Config: testAccComputeNetworkEndpoint_networkEndpointsAdditional(context),
			},
			{
				ResourceName:      "google_compute_network_endpoint.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_compute_network_endpoint.add1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_compute_network_endpoint.add2",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// delete all endpoints
				Config: testAccComputeNetworkEndpoint_noNetworkEndpoints(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeNetworkEndpointWithPortsDestroyed(t, negId, "100"),
				),
			},
		},
	})
}

func testAccComputeNetworkEndpoint_networkEndpointsBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network_endpoint" "default" {
  zone                   = "us-central1-a"
  network_endpoint_group = google_compute_network_endpoint_group.neg.id

  instance   = google_compute_instance.default.name
  ip_address = google_compute_instance.default.network_interface[0].network_ip
  port       = google_compute_network_endpoint_group.neg.default_port
}
`, context) + testAccComputeNetworkEndpoint_noNetworkEndpoints(context)
}

func testAccComputeNetworkEndpoint_networkEndpointsModified(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network_endpoint" "default" {
  zone                   = "us-central1-a"
  network_endpoint_group = google_compute_network_endpoint_group.neg.name

  instance   = google_compute_instance.default.name
  ip_address = google_compute_instance.default.network_interface[0].network_ip
  port       = "%{modified_port}"
}
`, context) + testAccComputeNetworkEndpoint_noNetworkEndpoints(context)
}

func testAccComputeNetworkEndpoint_networkEndpointsAdditional(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network_endpoint" "default" {
  zone                   = "us-central1-a"
  network_endpoint_group = google_compute_network_endpoint_group.neg.id

  instance   = google_compute_instance.default.name
  ip_address = google_compute_instance.default.network_interface[0].network_ip
  port       = "%{modified_port}"
}

resource "google_compute_network_endpoint" "add1" {
  zone                   = "us-central1-a"
  network_endpoint_group = google_compute_network_endpoint_group.neg.id

  instance   = google_compute_instance.default.name
  ip_address = google_compute_instance.default.network_interface[0].network_ip
  port       = "%{add1_port}"
}

resource "google_compute_network_endpoint" "add2" {
  zone                   = "us-central1-a"
  network_endpoint_group = google_compute_network_endpoint_group.neg.name

  instance   = google_compute_instance.default.name
  ip_address = google_compute_instance.default.network_interface[0].network_ip
  port       = "%{add2_port}"
}
`, context) + testAccComputeNetworkEndpoint_noNetworkEndpoints(context)
}

func testAccComputeNetworkEndpoint_noNetworkEndpoints(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network_endpoint_group" "neg" {
  name         = "tf-test-neg-%{random_suffix}"
  zone         = "us-central1-a"
  network      = google_compute_network.default.self_link
  subnetwork   = google_compute_subnetwork.default.self_link
  default_port = "%{default_port}"
}

resource "google_compute_network" "default" {
  name                    = "tf-test-neg-network-%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "default" {
  name          = "tf-test-neg-subnetwork-%{random_suffix}"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.default.self_link
}

resource "google_compute_instance" "default" {
  name         = "tf-test-neg-%{random_suffix}"
  machine_type = "e2-medium"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    subnetwork = google_compute_subnetwork.default.self_link
    access_config {
    }
  }
}

data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}
`, context)
}

// testAccCheckComputeNetworkEndpointDestroyed makes sure the endpoint with
// given Terraform resource name and previous information (obtained from Exists)
// was destroyed properly.
func testAccCheckComputeNetworkEndpointWithPortsDestroyed(t *testing.T, negId string, ports ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		foundPorts, err := testAccComputeNetworkEndpointListEndpointPorts(t, negId)
		if err != nil {
			return fmt.Errorf("unable to confirm endpoints with ports %+v was destroyed: %v", ports, err)
		}
		for _, p := range ports {
			if _, ok := foundPorts[p]; ok {
				return fmt.Errorf("network endpoint with port %s still exists", p)
			}
		}

		return nil
	}
}

func testAccComputeNetworkEndpointListEndpointPorts(t *testing.T, negId string) (map[string]struct{}, error) {
	config := acctest.GoogleProviderConfig(t)

	url := fmt.Sprintf("https://www.googleapis.com/compute/v1/%s/listNetworkEndpoints", negId)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		RawURL:    url,
		UserAgent: config.UserAgent,
	})
	if err != nil {
		return nil, err
	}

	v, ok := res["items"]
	if !ok || v == nil {
		return nil, nil
	}
	items := v.([]interface{})
	ports := make(map[string]struct{})
	for _, item := range items {
		endptWithHealth := item.(map[string]interface{})
		v, ok := endptWithHealth["networkEndpoint"]
		if !ok || v == nil {
			continue
		}
		endpt := v.(map[string]interface{})
		ports[fmt.Sprintf("%v", endpt["port"])] = struct{}{}
	}
	return ports, nil
}
