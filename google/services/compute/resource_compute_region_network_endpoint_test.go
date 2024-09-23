// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccComputeRegionNetworkEndpoint_regionNetworkEndpointBasic(t *testing.T) {
	// Multiple fine-grained resources
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"modified_port": 100,
		"add1_port":     101,
		"add2_port":     102,
	}
	negId := fmt.Sprintf("projects/%s/regions/%s/networkEndpointGroups/tf-test-neg-%s",
		envvar.GetTestProjectFromEnv(), envvar.GetTestRegionFromEnv(), context["random_suffix"])

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Create one endpoint
				Config: testAccComputeRegionNetworkEndpoint_regionNetworkEndpointBasic(context),
			},
			{
				ResourceName:      "google_compute_region_network_endpoint.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Force-recreate old endpoint
				Config: testAccComputeRegionNetworkEndpoint_regionNetworkEndpointsModified(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeRegionNetworkEndpointWithPortsDestroyed(t, negId, "90"),
				),
			},
			{
				ResourceName:      "google_compute_region_network_endpoint.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Add two new endpoints
				Config: testAccComputeRegionNetworkEndpoint_regionNetworkEndpointsAdditional(context),
			},
			{
				ResourceName:      "google_compute_region_network_endpoint.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_compute_region_network_endpoint.add1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_compute_region_network_endpoint.add2",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Remove add1 and add2 endpoints
				Config: testAccComputeRegionNetworkEndpoint_regionNetworkEndpointsModified(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeRegionNetworkEndpointWithPortsDestroyed(t, negId, "90"),
				),
			},
			{
				ResourceName:      "google_compute_region_network_endpoint.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Delete all endpoints
				Config: testAccComputeRegionNetworkEndpoint_noRegionNetworkEndpoints(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeRegionNetworkEndpointWithPortsDestroyed(t, negId, "100"),
				),
			},
		},
	})
}

func testAccComputeRegionNetworkEndpoint_regionNetworkEndpointBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_region_network_endpoint" "default" {
  region                        = "us-central1"
  region_network_endpoint_group = google_compute_region_network_endpoint_group.neg.id

  ip_address = "8.8.8.8"
  port       = 443
}
`, context) + testAccComputeRegionNetworkEndpoint_noRegionNetworkEndpoints(context)
}

func testAccComputeRegionNetworkEndpoint_regionNetworkEndpointsModified(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_region_network_endpoint" "default" {
  region                        = "us-central1"
  region_network_endpoint_group = google_compute_region_network_endpoint_group.neg.name

  ip_address = "8.8.8.8"
  port       = "%{modified_port}"
}
`, context) + testAccComputeRegionNetworkEndpoint_noRegionNetworkEndpoints(context)
}

func testAccComputeRegionNetworkEndpoint_regionNetworkEndpointsAdditional(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_region_network_endpoint" "default" {
  region                        = "us-central1"
  region_network_endpoint_group = google_compute_region_network_endpoint_group.neg.id

  ip_address = "8.8.8.8"
  port       = "%{modified_port}"
}

resource "google_compute_region_network_endpoint" "add1" {
  region                        = "us-central1"
  region_network_endpoint_group = google_compute_region_network_endpoint_group.neg.id

  ip_address = "8.8.8.8"
  port       = "%{add1_port}"
}

resource "google_compute_region_network_endpoint" "add2" {
  region                        = "us-central1"
  region_network_endpoint_group = google_compute_region_network_endpoint_group.neg.name

  ip_address = "8.8.8.8"
  port       = "%{add2_port}"
}
`, context) + testAccComputeRegionNetworkEndpoint_noRegionNetworkEndpoints(context)
}

func TestAccComputeRegionNetworkEndpoint_regionNetworkEndpointPortmapExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	negId := fmt.Sprintf("projects/%s/regions/%s/networkEndpointGroups/tf-test-portmap-neg%s",
		envvar.GetTestProjectFromEnv(), envvar.GetTestRegionFromEnv(), context["random_suffix"])

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionNetworkEndpoint_regionNetworkEndpointPortmapExample(context),
			},
			{
				ResourceName:            "google_compute_region_network_endpoint.region_network_endpoint_portmap",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance", "region", "region_network_endpoint_group"},
			},
			{
				// Delete all endpoints
				Config: testAccComputeRegionNetworkEndpoint_regionNetworkEndpointPortmapNoEndpointExample(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeRegionNetworkEndpointWithPortsDestroyed(t, negId, "80"),
				),
			},
		},
	})
}

func testAccComputeRegionNetworkEndpoint_regionNetworkEndpointPortmapNoEndpointExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "default" {
  name                    = "network%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "default" {
  name          = "subnetwork%{random_suffix}"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.default.id
}

resource "google_compute_region_network_endpoint_group" default {
  name                  = "tf-test-portmap-neg%{random_suffix}"
  region                = "us-central1"
  network               = google_compute_network.default.id
  subnetwork            = google_compute_subnetwork.default.id

  network_endpoint_type = "GCE_VM_IP_PORTMAP"
}

data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "default" {
  name         = "instance%{random_suffix}"
  machine_type = "e2-medium"
  zone = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    subnetwork = google_compute_subnetwork.default.id
    access_config {
    }
  }
}
`, context)
}

func testAccComputeRegionNetworkEndpoint_regionNetworkEndpointPortmapExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_region_network_endpoint" "region_network_endpoint_portmap" {
  region_network_endpoint_group = google_compute_region_network_endpoint_group.default.name
  region = "us-central1"
  instance   = google_compute_instance.default.self_link
  port       = 80
  ip_address = google_compute_instance.default.network_interface[0].network_ip
  client_destination_port = 8080
}
`, context) + testAccComputeRegionNetworkEndpoint_regionNetworkEndpointPortmapNoEndpointExample(context)
}

func testAccComputeRegionNetworkEndpoint_noRegionNetworkEndpoints(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_region_network_endpoint_group" "neg" {
  name                  = "tf-test-neg-%{random_suffix}"
  region                = "us-central1"
  network               = google_compute_network.default.self_link
  network_endpoint_type = "INTERNET_IP_PORT"
}

resource "google_compute_network" "default" {
  name                    = "tf-test-neg-network-%{random_suffix}"
  auto_create_subnetworks = false
}
`, context)
}

// testAccCheckComputeRegionNetworkEndpointDestroyed makes sure the endpoint with
// given Terraform resource name and previous information (obtained from Exists)
// was destroyed properly.
func testAccCheckComputeRegionNetworkEndpointWithPortsDestroyed(t *testing.T, negId string, ports ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		foundPorts, err := testAccComputeRegionNetworkEndpointListEndpointPorts(t, negId)
		if err != nil {
			return fmt.Errorf("unable to confirm endpoints with ports %+v was destroyed: %v", ports, err)
		}
		for _, p := range ports {
			if _, ok := foundPorts[p]; ok {
				return fmt.Errorf("region network endpoint with port %s still exists", p)
			}
		}

		return nil
	}
}

func testAccComputeRegionNetworkEndpointListEndpointPorts(t *testing.T, negId string) (map[string]struct{}, error) {
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
