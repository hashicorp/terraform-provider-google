// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package apphub_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccApphubService_serviceUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
			"time":   {},
		},
		CheckDestroy: testAccCheckApphubServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApphubService_apphubServiceFullExample(context),
			},
			{
				ResourceName:            "google_apphub_service.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "application_id", "service_id"},
			},
			{
				Config: testAccApphubService_apphubServiceUpdate(context),
			},
			{
				ResourceName:            "google_apphub_service.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "application_id", "service_id"},
			},
		},
	})
}

func testAccApphubService_apphubServiceUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_apphub_application" "application" {
  location = "us-central1"
  application_id = "tf-test-example-application-1%{random_suffix}"
  scope {
    type = "REGIONAL"
  }
}

resource "google_project" "service_project" {
  project_id ="tf-test-project-1%{random_suffix}"
  name = "Service Project"
  org_id = "%{org_id}"
  billing_account = "%{billing_account}"
}

# Enable Compute API
resource "google_project_service" "compute_service_project" {
  project = google_project.service_project.project_id
  service = "compute.googleapis.com"
}

resource "time_sleep" "wait_120s" {
  depends_on = [google_project_service.compute_service_project]

  create_duration = "120s"
}

resource "google_apphub_service_project_attachment" "service_project_attachment" {
  service_project_attachment_id = google_project.service_project.project_id
  depends_on = [time_sleep.wait_120s]
}

# discovered service block
data "google_apphub_discovered_service" "catalog-service" {
  provider = google
  location = "us-central1"
  service_uri = "//compute.googleapis.com/${google_compute_forwarding_rule.forwarding_rule.id}"
  depends_on = [google_apphub_service_project_attachment.service_project_attachment, time_sleep.wait_120s_for_resource_ingestion]
}

resource "time_sleep" "wait_120s_for_resource_ingestion" {
  depends_on = [google_compute_forwarding_rule.forwarding_rule]
  create_duration = "120s"
}

resource "google_apphub_service" "example" {
  location = "us-central1"
  application_id = google_apphub_application.application.application_id
  service_id = google_compute_forwarding_rule.forwarding_rule.name
  discovered_service = data.google_apphub_discovered_service.catalog-service.name
}


#creates service


# VPC network
resource "google_compute_network" "ilb_network" {
  name                    = "tf-test-l7-ilb-network%{random_suffix}"
  project                 = google_project.service_project.project_id
  auto_create_subnetworks = false
  depends_on = [time_sleep.wait_120s]
}


# backend subnet
resource "google_compute_subnetwork" "ilb_subnet" {
  name          = "tf-test-l7-ilb-subnet%{random_suffix}"
  project       = google_project.service_project.project_id
  ip_cidr_range = "10.0.1.0/24"
  region        = "us-central1"
  network       = google_compute_network.ilb_network.id
}

# forwarding rule
resource "google_compute_forwarding_rule" "forwarding_rule" {
  name                  ="tf-test-l7-ilb-forwarding-rule%{random_suffix}"
  project               = google_project.service_project.project_id
  region                = "us-central1"
  ip_version            = "IPV4"
  load_balancing_scheme = "INTERNAL"
  all_ports             = true
  backend_service       = google_compute_region_backend_service.backend.id
  network               = google_compute_network.ilb_network.id
  subnetwork            = google_compute_subnetwork.ilb_subnet.id
}



# backend service
resource "google_compute_region_backend_service" "backend" {
  name                  = "tf-test-l7-ilb-backend-subnet%{random_suffix}"
  project               = google_project.service_project.project_id
  region                = "us-central1"
  health_checks         = [google_compute_health_check.default.id]
}

# health check
resource "google_compute_health_check" "default" {
  name     = "tf-test-l7-ilb-hc%{random_suffix}"
  project  = google_project.service_project.project_id
  check_interval_sec = 1
  timeout_sec        = 1
  tcp_health_check {
    port = "80"
  }
  depends_on = [time_sleep.wait_120s]
}
`, context)
}
