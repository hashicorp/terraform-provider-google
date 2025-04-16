// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package apigee_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccApigeeEnvReferences_apigeeEnvironmentReferenceTest_Update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckApigeeEnvReferencesDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApigeeEnvReferences_apigeeEnvironmentReferenceTest_full(context),
			},
			{
				ResourceName:            "google_apigee_env_references.apigee_environment_reference",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"env_id"},
			},
			{
				Config: testAccApigeeEnvReferences_apigeeEnvironmentReferenceTest_update(context),
			},
			{
				ResourceName:            "google_apigee_env_references.apigee_environment_reference",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"env_id"},
			},
		},
	})
}

func testAccApigeeEnvReferences_apigeeEnvironmentReferenceTest_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  project_id      = "tf-test%{random_suffix}"
  name            = "tf-test%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
  deletion_policy = "DELETE"
}

resource "google_project_service" "apigee" {
  project = google_project.project.project_id
  service = "apigee.googleapis.com"
}

resource "google_project_service" "servicenetworking" {
  project = google_project.project.project_id
  service = "servicenetworking.googleapis.com"
  depends_on = [google_project_service.apigee]
}

resource "google_project_service" "compute" {
  project = google_project.project.project_id
  service = "compute.googleapis.com"
  depends_on = [google_project_service.servicenetworking]
}

resource "time_sleep" "wait_120_seconds" {
  create_duration = "120s"
  depends_on = [google_project_service.compute]
}

resource "google_compute_network" "apigee_network" {
  name       = "apigee-network"
  project    = google_project.project.project_id
  depends_on = [time_sleep.wait_120_seconds]
}

resource "google_compute_global_address" "apigee_range" {
  name          = "apigee-range"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.apigee_network.id
  project       = google_project.project.project_id
}

resource "google_service_networking_connection" "apigee_vpc_connection" {
  network                 = google_compute_network.apigee_network.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.apigee_range.name]
  depends_on              = [google_project_service.servicenetworking]
}

resource "google_apigee_organization" "apigee_org" {
  analytics_region   = "us-central1"
  project_id         = google_project.project.project_id
  authorized_network = google_compute_network.apigee_network.id
  depends_on         = [
    google_service_networking_connection.apigee_vpc_connection,
    google_project_service.apigee,
  ]
}

resource "google_apigee_environment" "apigee_environment" {
  org_id       = google_apigee_organization.apigee_org.id
  name         = "tf-test%{random_suffix}"
  description  = "Apigee Environment"
  display_name = "environment-1"
}

resource "google_apigee_env_keystore" "apigee_environment_keystore_1" {
  name       = "tf-test-keystore1%{random_suffix}"
  env_id     = google_apigee_environment.apigee_environment.id
}

resource "google_apigee_env_references" "apigee_environment_reference" {
  env_id         = google_apigee_environment.apigee_environment.id
  name           = "tf-test-reference%{random_suffix}"
  resource_type  = "KeyStore"
  refers         = google_apigee_env_keystore.apigee_environment_keystore_1.name
  depends_on = [google_apigee_env_keystore.apigee_environment_keystore_1]
}

resource "google_apigee_env_keystore" "apigee_environment_keystore_2" {
  name       = "tf-test-keystore2%{random_suffix}"
  env_id     = google_apigee_environment.apigee_environment.id
  depends_on = [google_apigee_env_references.apigee_environment_reference]
}
`, context)
}

func testAccApigeeEnvReferences_apigeeEnvironmentReferenceTest_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  project_id      = "tf-test%{random_suffix}"
  name            = "tf-test%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
  deletion_policy = "DELETE"
}

resource "google_project_service" "apigee" {
  project = google_project.project.project_id
  service = "apigee.googleapis.com"
}

resource "google_project_service" "servicenetworking" {
  project = google_project.project.project_id
  service = "servicenetworking.googleapis.com"
  depends_on = [google_project_service.apigee]
}

resource "google_project_service" "compute" {
  project = google_project.project.project_id
  service = "compute.googleapis.com"
  depends_on = [google_project_service.servicenetworking]
}

resource "time_sleep" "wait_120_seconds" {
  create_duration = "120s"
  depends_on = [google_project_service.compute]
}

resource "google_compute_network" "apigee_network" {
  name       = "apigee-network"
  project    = google_project.project.project_id
  depends_on = [time_sleep.wait_120_seconds]
}

resource "google_compute_global_address" "apigee_range" {
  name          = "apigee-range"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.apigee_network.id
  project       = google_project.project.project_id
}

resource "google_service_networking_connection" "apigee_vpc_connection" {
  network                 = google_compute_network.apigee_network.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.apigee_range.name]
  depends_on              = [google_project_service.servicenetworking]
}

resource "google_apigee_organization" "apigee_org" {
  analytics_region   = "us-central1"
  project_id         = google_project.project.project_id
  authorized_network = google_compute_network.apigee_network.id
  depends_on         = [
    google_service_networking_connection.apigee_vpc_connection,
    google_project_service.apigee,
  ]
}

resource "google_apigee_environment" "apigee_environment" {
  org_id       = google_apigee_organization.apigee_org.id
  name         = "tf-test%{random_suffix}"
  description  = "Apigee Environment"
  display_name = "environment-1"
}

resource "google_apigee_env_keystore" "apigee_environment_keystore_2" {
  name       = "tf-test-keystore2%{random_suffix}"
  env_id     = google_apigee_environment.apigee_environment.id
}

resource "google_apigee_env_references" "apigee_environment_reference" {
  env_id         = google_apigee_environment.apigee_environment.id
  name           = "tf-test-reference%{random_suffix}"
  resource_type  = "KeyStore"
  refers         = google_apigee_env_keystore.apigee_environment_keystore_2.name
  depends_on = [google_apigee_env_keystore.apigee_environment_keystore_2]
}

resource "google_apigee_env_keystore" "apigee_environment_keystore_1" {
  name       = "tf-test-keystore1%{random_suffix}"
  env_id     = google_apigee_environment.apigee_environment.id
  depends_on = [google_apigee_env_references.apigee_environment_reference]
}
`, context)
}
