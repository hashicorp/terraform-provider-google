// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package servicenetworking_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccServiceNetworkingPeeredDNSDomain_basic(t *testing.T) {
	t.Parallel()
	org := envvar.GetTestOrgFromEnv(t)
	billingId := envvar.GetTestBillingAccountFromEnv(t)

	project := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	service := "servicenetworking.googleapis.com"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceNetworkingPeeredDNSDomain_basic(project, org, billingId, name, service),
			},
		},
	})
}

func testAccServiceNetworkingPeeredDNSDomain_basic(project, org, billing, name, service string) string {
	return fmt.Sprintf(`
resource "google_project" "host" {
  project_id      = "%s"
  name            = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_service" "host-compute" {
  project = google_project.host.project_id
  service = "compute.googleapis.com"
}

resource "google_project_service" "host" {
  project = google_project.host.project_id
  service = "%s"
  depends_on = [google_project_service.host-compute]
}

resource "google_compute_network" "test" {
	name                    = "test-network"
	project                 = google_project.host.project_id
	routing_mode            = "GLOBAL"
  depends_on              = [google_project_service.host-compute]
}

resource "google_compute_global_address" "host-private-access" {
  name          = "%s-ip"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 24
  address       = "192.168.255.0"
  network       = "test-network"
  project       = google_project.host.project_id

	depends_on = [
		google_project_service.host-compute,
		google_project_service.host,
		google_compute_network.test,
	]
}

resource "google_service_networking_connection" "host-private-access" {
  network                 = google_compute_network.test.id
  service                 = "%s"
  reserved_peering_ranges = [google_compute_global_address.host-private-access.name]

	depends_on = [
		google_project_service.host,
		google_compute_network.test,
		google_compute_global_address.host-private-access,
	]
}

resource "google_service_networking_peered_dns_domain" "test" {
  name       = "%s"
	project    = google_project.host.number
	network    = "test-network"
	dns_suffix = "example.com."
  service    = "%s"
	depends_on = [
		google_compute_network.test,
		google_service_networking_connection.host-private-access,
  ]
}
`, project, project, org, billing, service, project, service, name, service)
}
