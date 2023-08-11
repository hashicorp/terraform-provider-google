// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package vertexai_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccVertexAIIndexEndpoint_updated(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"network_name":  acctest.BootstrapSharedTestNetwork(t, "vertex-ai-index-endpoint-update"),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVertexAIIndexEndpointDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIIndexEndpoint_basic(context),
			},
			{
				ResourceName:            "google_vertex_ai_index_endpoint.index_endpoint",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "region"},
			},
			{
				Config: testAccVertexAIIndexEndpoint_updated(context),
			},
			{
				ResourceName:            "google_vertex_ai_index_endpoint.index_endpoint",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "region"},
			},
		},
	})
}

func testAccVertexAIIndexEndpoint_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_index_endpoint" "index_endpoint" {
  display_name = "sample-endpoint"
  description  = "A sample vertex endpoint"
  region       = "us-central1"
  labels       = {
    label-one = "value-one"
  }
  network      = "projects/${data.google_project.project.number}/global/networks/${data.google_compute_network.vertex_network.name}"
  depends_on   = [
    google_service_networking_connection.vertex_vpc_connection
  ]
}
resource "google_service_networking_connection" "vertex_vpc_connection" {
  network                 = data.google_compute_network.vertex_network.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.vertex_range.name]
}
resource "google_compute_global_address" "vertex_range" {
  name          = "tf-test-address-name%{random_suffix}"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 24
  network       = data.google_compute_network.vertex_network.id
}
data "google_compute_network" "vertex_network" {
  name       = "%{network_name}"
}
data "google_project" "project" {}
`, context)
}

func testAccVertexAIIndexEndpoint_updated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_index_endpoint" "index_endpoint" {
  display_name = "sample-endpoint-updated"
  description  = "A sample vertex endpoint (updated)"
  region       = "us-central1"
  labels       = {
    label-one = "value-one"
    label-two = "value-two"
  }
  network      = "projects/${data.google_project.project.number}/global/networks/${data.google_compute_network.vertex_network.name}"
  depends_on   = [
    google_service_networking_connection.vertex_vpc_connection
  ]
}
resource "google_service_networking_connection" "vertex_vpc_connection" {
  network                 = data.google_compute_network.vertex_network.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.vertex_range.name]
}
resource "google_compute_global_address" "vertex_range" {
  name          = "tf-test-address-name%{random_suffix}"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 24
  network       = data.google_compute_network.vertex_network.id
}
data "google_compute_network" "vertex_network" {
  name       = "%{network_name}"
}
data "google_project" "project" {}
`, context)
}
