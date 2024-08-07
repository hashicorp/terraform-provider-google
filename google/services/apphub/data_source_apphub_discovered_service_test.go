// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package apphub_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceApphubDiscoveredService_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testDataSourceApphubDiscoveredService_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_apphub_discovered_service.catalog-service", "name"),
				),
			},
		},
	})
}

func testDataSourceApphubDiscoveredService_basic(context map[string]interface{}) string {
	return acctest.Nprintf(
		`
resource "google_project" "service_project" {
	project_id ="tf-test-ah-%{random_suffix}"
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
  location = "us-central1"
  # ServiceReference | Application Hub | Google Cloud
  # Using this reference means that this resource will not be provisioned until the forwarding rule is fully created
  service_uri = "//compute.googleapis.com/${google_compute_forwarding_rule.forwarding_rule.id}"
	depends_on = [time_sleep.wait_120s_for_resource_ingestion]
}

# VPC network
resource "google_compute_network" "ilb_network" {
  name                    = "ilb-network-%{random_suffix}"
  project                 = google_project.service_project.project_id
  auto_create_subnetworks = false
  depends_on = [time_sleep.wait_120s]
}

# backend subnet
resource "google_compute_subnetwork" "ilb_subnet" {
  name          			 = "ilb-subnet-%{random_suffix}"
  project       			 = google_project.service_project.project_id
  ip_cidr_range 			 = "10.0.1.0/24"
  region        			 = "us-central1"
  network       			 = google_compute_network.ilb_network.id
}

# forwarding rule
resource "google_compute_forwarding_rule" "forwarding_rule" {
  name                  = "forwarding-rule-%{random_suffix}"
  project               = google_project.service_project.project_id
  region                = "us-central1"
  ip_version            = "IPV4"
  load_balancing_scheme = "INTERNAL"
  all_ports             = true
  backend_service       = google_compute_region_backend_service.backend.id
  network               = google_compute_network.ilb_network.id
  subnetwork            = google_compute_subnetwork.ilb_subnet.id
}

resource "time_sleep" "wait_120s_for_resource_ingestion" {
  depends_on = [google_compute_forwarding_rule.forwarding_rule]
  create_duration = "120s"
}

# backend service
resource "google_compute_region_backend_service" "backend" {
  name                  = "backend-service-%{random_suffix}"
  project               = google_project.service_project.project_id
  region                = "us-central1"
  health_checks         = [google_compute_health_check.default.id]
}
    
# health check
resource "google_compute_health_check" "default" {
  name     					 		= "health-check-%{random_suffix}"
  project  					 		= google_project.service_project.project_id
  check_interval_sec 		= 1
  timeout_sec        		= 1

  tcp_health_check {
    port = "80"
  }
  depends_on = [time_sleep.wait_120s]
}
`, context)
}
