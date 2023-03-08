package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccApigeeSharedflowDeployment_apigeeSharedflowDeploymentTestExample(t *testing.T) {
	SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"org_id":          GetTestOrgFromEnv(t),
		"billing_account": GetTestBillingAccountFromEnv(t),
		"random_suffix":   RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccCheckApigeeSharedflowDeploymentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApigeeSharedflowDeployment_apigeeSharedflowDeploymentTestExample(context),
			},
			{
				ResourceName:            "google_apigee_sharedflow_deployment.sharedflow_deployment_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccApigeeSharedflowDeployment_apigeeSharedflowDeploymentTestExample(context map[string]interface{}) string {
	return Nprintf(`
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
  org_id   = google_apigee_organization.apigee_org.id
  name         = "tf-test%{random_suffix}"
  description  = "Apigee Environment"
  display_name = "environment-1"
}

resource "google_apigee_sharedflow" "test_apigee_sharedflow" {
  name            = "tf-test-apigee-sharedflow"
  org_id          = google_project.project.project_id
  config_bundle   = "./test-fixtures/apigee/apigee_sharedflow_bundle.zip"
  depends_on      = [google_apigee_organization.apigee_org]
}

resource "google_apigee_sharedflow_deployment" "sharedflow_deployment_test" {
  environment = google_apigee_environment.apigee_environment.name
  org_id = google_apigee_sharedflow.test_apigee_sharedflow.org_id
  revision = google_apigee_sharedflow.test_apigee_sharedflow.revision[length(google_apigee_sharedflow.test_apigee_sharedflow.revision)-1]
  sharedflow_id = google_apigee_sharedflow.test_apigee_sharedflow.name
}
`, context)
}

func testAccCheckApigeeSharedflowDeploymentDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_apigee_sharedflow_deployment" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := GoogleProviderConfig(t)

			url, err := replaceVarsForTest(config, rs, "{{ApigeeBasePath}}organizations/{{org_id}}/environments/{{environment}}/sharedflows/{{sharedflow_id}}/revisions/{{revision}}/deployments")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}
			_, err = SendRequest(config, "GET", billingProject, url, config.UserAgent, nil)
			if err == nil {
				return fmt.Errorf("ApigeeSharedFlow still exists at %s", url)
			}
		}

		return nil
	}
}
