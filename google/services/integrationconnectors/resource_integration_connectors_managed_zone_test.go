// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package integrationconnectors_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccIntegrationConnectorsManagedZone_integrationConnectorsManagedZoneExample_update(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIntegrationConnectorsManagedZoneDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationConnectorsManagedZone_integrationConnectorsManagedZoneExample_full(context),
			},
			{
				ResourceName:            "google_integration_connectors_managed_zone.testmanagedzone",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "labels", "terraform_labels"},
			},
			{
				Config: testAccIntegrationConnectorsManagedZone_integrationConnectorsManagedZoneExample_update(context),
			},
			{
				ResourceName:            "google_integration_connectors_managed_zone.testmanagedzone",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccIntegrationConnectorsManagedZone_integrationConnectorsManagedZoneExample_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "target_project" {
  project_id      = "tf-test%{random_suffix}"
  name            = "tf-test%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}

resource "google_project_iam_member" "dns_peer_binding" {
  project = google_project.target_project.project_id
  role    = "roles/dns.peer"
  member  = "serviceAccount:service-${data.google_project.test_project.number}@gcp-sa-connectors.iam.gserviceaccount.com"
}

resource "google_project_service" "dns" {
  project = google_project.target_project.project_id
  service = "dns.googleapis.com"
}

resource "google_project_service" "compute" {
  project = google_project.target_project.project_id
  service = "compute.googleapis.com"
}

resource "google_compute_network" "network" {
  project = google_project.target_project.project_id
  name                    = "test"
  auto_create_subnetworks = false
  depends_on = [google_project_service.compute]
}

resource "google_dns_managed_zone" "zone" {
  name        = "tf-test-dns%{random_suffix}"
  dns_name    = "private%{random_suffix}.example.com."
  visibility  = "private"

  private_visibility_config {
    networks {
      network_url = google_compute_network.network.id
    }
  }
  depends_on = [google_project_service.dns]
}

data "google_project" "test_project" {
}

resource "google_integration_connectors_managed_zone" "testmanagedzone" {
  name     = "test%{random_suffix}"
  description = "tf created description"
  labels = {
    intent = "example"
  }
  target_project = google_project.target_project.project_id
  target_vpc="test"
  dns=google_dns_managed_zone.zone.dns_name
  depends_on = [google_project_iam_member.dns_peer_binding,google_dns_managed_zone.zone]
}
`, context)
}

func testAccIntegrationConnectorsManagedZone_integrationConnectorsManagedZoneExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "target_project" {
  project_id      = "tf-test%{random_suffix}"
  name            = "tf-test%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}

resource "google_project_iam_member" "dns_peer_binding" {
  project = google_project.target_project.project_id
  role    = "roles/dns.peer"
  member  = "serviceAccount:service-${data.google_project.test_project.number}@gcp-sa-connectors.iam.gserviceaccount.com"
}

resource "google_project_service" "dns" {
  project = google_project.target_project.project_id
  service = "dns.googleapis.com"
}

resource "google_project_service" "compute" {
  project = google_project.target_project.project_id
  service = "compute.googleapis.com"
}

resource "google_compute_network" "network" {
  project = google_project.target_project.project_id
  name                    = "test"
  auto_create_subnetworks = false
  depends_on = [google_project_service.compute]
}

resource "google_dns_managed_zone" "zone" {
  name        = "tf-test-dns%{random_suffix}"
  dns_name    = "private%{random_suffix}.example.com."
  visibility  = "private"

  private_visibility_config {
    networks {
      network_url = google_compute_network.network.id
    }
  }
  depends_on = [google_project_service.dns]
}

data "google_project" "test_project" {
}

resource "google_integration_connectors_managed_zone" "testmanagedzone" {
  name     = "test%{random_suffix}"
  description = "tf updated description"
  labels = {
    intent = "example"
  }
  target_project = google_project.target_project.project_id
  target_vpc="test"
  dns=google_dns_managed_zone.zone.dns_name
  depends_on = [google_project_iam_member.dns_peer_binding,google_dns_managed_zone.zone]
}
`, context)
}
