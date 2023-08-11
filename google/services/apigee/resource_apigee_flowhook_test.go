// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package apigee_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccApigeeFlowhook_apigeeFlowhookTestExample(t *testing.T) {
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
		CheckDestroy:             testAccCheckApigeeFlowhookDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApigeeFlowhook_apigeeFlowhookTestExample(context),
			},
			{
				ResourceName:            "google_apigee_flowhook.flowhook_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccApigeeFlowhook_apigeeFlowhookTestExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  project_id      = "tf-test%{random_suffix}"
  name            = "tf-test%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
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

resource "google_compute_network" "apigee_network" {
  name       = "apigee-network"
  project    = google_project.project.project_id
  depends_on = [google_project_service.compute]
}

resource "google_compute_global_address" "apigee_range" {
  name          = "tf-test-apigee-range%{random_suffix}"
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
  org_id   = google_apigee_organization.apigee_org.id
  name         = "tf-test%{random_suffix}"
  description  = "Apigee Environment"
  display_name = "environment-1"
}

resource "google_apigee_sharedflow" "test_apigee_sharedflow" {
  name            = "tf-test-apigee-sharedflow"
  org_id          = google_project.project.project_id
  config_bundle   = "./test-fixtures/apigee_sharedflow_bundle.zip"
  depends_on      = [google_apigee_organization.apigee_org]
}

resource "google_apigee_sharedflow_deployment" "sharedflow_deployment_test" {
  environment = google_apigee_environment.apigee_environment.name
  org_id = google_apigee_sharedflow.test_apigee_sharedflow.org_id
  revision = google_apigee_sharedflow.test_apigee_sharedflow.revision[length(google_apigee_sharedflow.test_apigee_sharedflow.revision)-1]
  sharedflow_id = google_apigee_sharedflow.test_apigee_sharedflow.name
}

resource "google_apigee_flowhook" "flowhook_test" {
	environment = google_apigee_sharedflow_deployment.sharedflow_deployment_test.environment
	org_id = google_apigee_sharedflow.test_apigee_sharedflow.org_id
	flow_hook_point = "PreProxyFlowHook"
	sharedflow = google_apigee_sharedflow.test_apigee_sharedflow.name
	description = "test flowhook"
	continue_on_error = true
  }
`, context)
}

func testAccCheckApigeeFlowhookDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_apigee_flowhook" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{ApigeeBasePath}}organizations/{{org_id}}/environments/{{environment}}/flowhooks/{{flow_hook_point}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}
			res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			// Flowhooks always exist, we treat the binding as a removable resource, thus we check if the sharedFlow field to detect sharedflow attachment
			if err == nil && res != nil && res["sharedFlow"] != nil {
				return fmt.Errorf("Flowhook still has an attachment at %s", url)
			}
		}

		return nil
	}
}
