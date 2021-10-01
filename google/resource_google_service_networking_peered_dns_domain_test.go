package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccServiceNetworkingPeeredDNSDomain_basic(t *testing.T) {
	t.Parallel()
	org := getTestOrgFromEnv(t)
	billingId := getTestBillingAccountFromEnv(t)

	project := fmt.Sprintf("tf-test-%d", randInt(t))
	name := fmt.Sprintf("test-name-%d", randInt(t))
	service := "servicenetworking.googleapis.com"

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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
}

resource "google_compute_network" "test" {
	name                    = "test-network"
	project                 = google_project.host.project_id
	routing_mode            = "GLOBAL"
  depends_on              = [google_project_service.host-compute]
}

resource "google_compute_global_address" "host-private-access" {
  name          = "private-ip-alloc-host"
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
`, project, project, org, billing, service, service, name, service)
}
