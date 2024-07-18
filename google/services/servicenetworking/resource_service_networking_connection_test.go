// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package servicenetworking_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccServiceNetworkingConnection_create(t *testing.T) {
	t.Parallel()

	network := fmt.Sprintf("tf-test-service-networking-connection-create-%s", acctest.RandString(t, 10))
	addr := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	service := "servicenetworking.googleapis.com"
	org_id := envvar.GetTestOrgFromEnv(t)
	billing_account := envvar.GetTestBillingAccountFromEnv(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testServiceNetworkingConnectionDestroy(t, service, network),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceNetworkingConnection(network, addr, "servicenetworking.googleapis.com", org_id, billing_account),
			},
			{
				ResourceName:      "google_service_networking_connection.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccServiceNetworkingConnection_abandon(t *testing.T) {
	t.Parallel()

	network := fmt.Sprintf("tf-test-service-networking-connection-abandon-%s", acctest.RandString(t, 10))
	addr := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	service := "servicenetworking.googleapis.com"
	org_id := envvar.GetTestOrgFromEnv(t)
	billing_account := envvar.GetTestBillingAccountFromEnv(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testServiceNetworkingConnectionDestroyAbandon(t, service, network),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceNetworkingConnectionToBeAbandoned(network, addr, "servicenetworking.googleapis.com", org_id, billing_account),
			},
			{
				ResourceName:      "google_service_networking_connection.foobar",
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func TestAccServiceNetworkingConnection_update(t *testing.T) {
	t.Parallel()

	network := fmt.Sprintf("tf-test-service-networking-connection-update-%s", acctest.RandString(t, 10))
	addr1 := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	addr2 := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	service := "servicenetworking.googleapis.com"
	org_id := envvar.GetTestOrgFromEnv(t)
	billing_account := envvar.GetTestBillingAccountFromEnv(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testServiceNetworkingConnectionDestroy(t, service, network),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceNetworkingConnection(network, addr1, "servicenetworking.googleapis.com", org_id, billing_account),
			},
			{
				ResourceName:      "google_service_networking_connection.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccServiceNetworkingConnection(network, addr2, "servicenetworking.googleapis.com", org_id, billing_account),
			},
			{
				ResourceName:      "google_service_networking_connection.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

}

func testServiceNetworkingConnectionDestroy(t *testing.T, parent, network string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)
		parentService := "services/" + parent
		networkName := fmt.Sprintf("projects/%s/global/networks/%s", envvar.GetTestProjectFromEnv(), network)
		listCall := config.NewServiceNetworkingClient(config.UserAgent).Services.Connections.List(parentService).Network(networkName)
		if config.UserProjectOverride {
			listCall.Header().Add("X-Goog-User-Project", envvar.GetTestProjectFromEnv())
		}
		response, err := listCall.Do()
		if err != nil {
			return err
		}

		for _, c := range response.Connections {
			if c.Network == networkName {
				return fmt.Errorf("Found %s which should have been destroyed.", networkName)
			}
		}

		return nil
	}
}

func testServiceNetworkingConnectionDestroyAbandon(t *testing.T, parent, network string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)
		parentService := "services/" + parent
		networkName := fmt.Sprintf("projects/%s/global/networks/%s", envvar.GetTestProjectFromEnv(), network)
		listCall := config.NewServiceNetworkingClient(config.UserAgent).Services.Connections.List(parentService).Network(networkName)
		if config.UserProjectOverride {
			listCall.Header().Add("X-Goog-User-Project", envvar.GetTestProjectFromEnv())
		}
		response, err := listCall.Do()
		if err != nil {
			return err
		}

		for _, c := range response.Connections {
			if c.Network == networkName {
				return fmt.Errorf("Found %s which should have been destroyed.", networkName)
			}
		}

		return nil
	}
}

func testAccServiceNetworkingConnection(networkName, addressRangeName, serviceName, org_id, billing_account string) string {
	return fmt.Sprintf(`
resource "google_project" "project" {
  project_id      = "%s"
  name            = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_service" "servicenetworking" {
  project = google_project.project.project_id
  service = "servicenetworking.googleapis.com"
}

resource "google_compute_network" "servicenet" {
  name = "%s"
  depends_on = [google_project_service.servicenetworking]
}

resource "google_compute_global_address" "foobar" {
  name          = "%s"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.servicenet.self_link
  depends_on = [google_project_service.servicenetworking]
}

resource "google_service_networking_connection" "foobar" {
  network                 = google_compute_network.servicenet.self_link
  service                 = "%s"
  reserved_peering_ranges = [google_compute_global_address.foobar.name]
  depends_on = [google_project_service.servicenetworking]
}
`, addressRangeName, addressRangeName, org_id, billing_account, networkName, addressRangeName, serviceName)
}

func testAccServiceNetworkingConnectionToBeAbandoned(networkName, addressRangeName, serviceName, org_id, billing_account string) string {
	return fmt.Sprintf(`
resource "google_project" "project" {
  project_id      = "%s"
  name            = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_service" "servicenetworking" {
  project = google_project.project.project_id
  service = "servicenetworking.googleapis.com"
}

resource "google_compute_network" "servicenet" {
  name = "%s"
  depends_on = [google_project_service.servicenetworking]
}

resource "google_compute_global_address" "foobar" {
  name          = "%s"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.servicenet.self_link
  depends_on = [google_project_service.servicenetworking]
}

resource "google_service_networking_connection" "foobar" {
  network                 = google_compute_network.servicenet.self_link
  service                 = "%s"
  reserved_peering_ranges = [google_compute_global_address.foobar.name]
  depends_on = [google_project_service.servicenetworking]
  deletion_policy = "ABANDON"
}
`, addressRangeName, addressRangeName, org_id, billing_account, networkName, addressRangeName, serviceName)
}
