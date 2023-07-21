// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccComputeNetworkEndpoints_networkEndpointsBasic(t *testing.T) {
	t.Parallel()

	// detachNetworkEndpoints call ordering is not guaranteed, causing VCR to rerecord
	acctest.SkipIfVcr(t)

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"default_port":  90,
		"modified_port": 100,
		"add1_port":     101,
		"add2_port":     102,
		"add3_port":     103,
		"add4_port":     104,
		"add5_port":     105,
	}
	negId := fmt.Sprintf("projects/%s/zones/%s/networkEndpointGroups/tf-test-neg-%s",
		envvar.GetTestProjectFromEnv(), envvar.GetTestZoneFromEnv(), context["random_suffix"])

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Create one endpoint
				Config: testAccComputeNetworkEndpoints_networkEndpointsBase(context),
			},
			{
				ResourceName:      "google_compute_network_endpoints.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Force-recreate old endpoint
				Config: testAccComputeNetworkEndpoints_networkEndpointsModified(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeNetworkEndpointsWithPortsDestroyed(t, negId, "90"),
				),
			},
			{
				ResourceName:      "google_compute_network_endpoints.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Add four new endpoints
				Config: testAccComputeNetworkEndpoints_networkEndpointsAdditional(context),
			},
			{
				ResourceName:      "google_compute_network_endpoints.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Add enough endpoints to trigger pagination
				Config: testAccComputeNetworkEndpoints_networkEndpointsPaginated(context, 100, 1300),
			},
			{
				ResourceName:      "google_compute_network_endpoints.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Remove enough endpoints to trigger pagination
				Config: testAccComputeNetworkEndpoints_networkEndpointsPaginated(context, 700, 1900),
			},
			{
				ResourceName:      "google_compute_network_endpoints.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// delete all endpoints
				Config: testAccComputeNetworkEndpoints_noNetworkEndpoints(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeNetworkEndpointsWithAllEndpointsDestroyed(t, negId),
				),
			},
		},
	})
}

func testAccComputeNetworkEndpoints_networkEndpointsBase(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network_endpoints" "default" {
  zone                   = "us-central1-a"
  network_endpoint_group = google_compute_network_endpoint_group.neg.id

  network_endpoints {
    instance   = google_compute_instance.default.name
    ip_address = google_compute_instance.default.network_interface[0].network_ip
    port       = google_compute_network_endpoint_group.neg.default_port
  }
}
`, context) + testAccComputeNetworkEndpoints_noNetworkEndpoints(context)
}

func testAccComputeNetworkEndpoints_networkEndpointsModified(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network_endpoints" "default" {
  zone                   = "us-central1-a"
  network_endpoint_group = google_compute_network_endpoint_group.neg.name

  network_endpoints {
    instance   = google_compute_instance.default.name
    ip_address = google_compute_instance.default.network_interface[0].network_ip
    port       = "%{modified_port}"
  }
}
`, context) + testAccComputeNetworkEndpoints_noNetworkEndpoints(context)
}

func testAccComputeNetworkEndpoints_networkEndpointsAdditional(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network_endpoints" "default" {
  zone                   = "us-central1-a"
  network_endpoint_group = google_compute_network_endpoint_group.neg.id
 
  network_endpoints {
    instance   = google_compute_instance.default.name
    ip_address = google_compute_instance.default.network_interface[0].network_ip
    port       = "%{modified_port}"
  }
  network_endpoints {
    instance   = google_compute_instance.default.name
    ip_address = google_compute_instance.default.network_interface[0].network_ip
    port       = "%{add1_port}"
  }
  network_endpoints {
    instance   = google_compute_instance.default.name
    ip_address = google_compute_instance.default.network_interface[0].network_ip
    port       = "%{add2_port}"
  }
  network_endpoints {
    instance   = google_compute_instance.default.name
    ip_address = google_compute_instance.default.network_interface[0].network_ip
    port       = "%{add3_port}"
  }
  network_endpoints {
    instance   = google_compute_instance.default.name
    ip_address = google_compute_instance.default.network_interface[0].network_ip
    port       = "%{add4_port}"
  }
  network_endpoints {
    instance   = google_compute_instance.default.name
    ip_address = google_compute_instance.default.network_interface[0].network_ip
    port       = "%{add5_port}"
  }
}
`, context) + testAccComputeNetworkEndpoints_noNetworkEndpoints(context)
}

func testAccComputeNetworkEndpoints_networkEndpointsPaginated(context map[string]interface{}, lower, upper int) string {
	context["for_each"] = networkEndpointsGenerateRanges(lower, upper)
	return acctest.Nprintf(`
resource "google_compute_network_endpoints" "default" {
  zone                   = "us-central1-a"
  network_endpoint_group = google_compute_network_endpoint_group.neg.name

  dynamic "network_endpoints" {
    for_each = %{for_each}
    content {
      instance   = google_compute_instance.default.name
      ip_address = google_compute_instance.default.network_interface[0].network_ip
      port       = network_endpoints.value
    }
  }
}
`, context) + testAccComputeNetworkEndpoints_noNetworkEndpoints(context)
}

// Terraform `range` can only generate a list of 1024 elements, so we need to
// concat them to get a longer list
func networkEndpointsGenerateRanges(lower, upper int) string {
	var ranges []string
	l := lower
	for l < upper {
		u := l + 1024
		if u > upper {
			u = upper
		}
		ranges = append(ranges, fmt.Sprintf("range(%d, %d)", l, u))
		l += 1024
	}
	return fmt.Sprintf("concat(%s)", strings.Join(ranges, ", "))
}

func testAccComputeNetworkEndpoints_noNetworkEndpoints(context map[string]interface{}) string {
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
func testAccCheckComputeNetworkEndpointsWithPortsDestroyed(t *testing.T, negId string, ports ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		foundPorts, err := testAccComputeNetworkEndpointsListEndpointPorts(t, negId)
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

func testAccCheckComputeNetworkEndpointsWithAllEndpointsDestroyed(t *testing.T, negId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		endpoints, err := testAccComputeNetworkEndpointsListEndpoints(t, negId)
		if err != nil {
			return fmt.Errorf("unable to confirm all endpoints were destroyed: %v", err)
		}
		if len(endpoints) > 0 {
			return fmt.Errorf("Not all network endpoints were deleted: %v", endpoints)
		}
		return nil
	}
}

func testAccComputeNetworkEndpointsListEndpoints(t *testing.T, negId string) ([]interface{}, error) {
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
	return v.([]interface{}), nil
}

func testAccComputeNetworkEndpointsListEndpointPorts(t *testing.T, negId string) (map[string]struct{}, error) {
	items, err := testAccComputeNetworkEndpointsListEndpoints(t, negId)
	if err != nil {
		return nil, err
	}
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
