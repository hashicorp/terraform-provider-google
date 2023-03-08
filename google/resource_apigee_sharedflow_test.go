package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccApigeeSharedFlow_apigeeSharedflowTestExample(t *testing.T) {
	SkipIfVcr(t)
	t.Parallel()

	fmt.Printf("from t: org_id %s", GetTestOrgFromEnv(t))

	context := map[string]interface{}{
		"org_id":          GetTestOrgFromEnv(t),
		"billing_account": GetTestBillingAccountFromEnv(t),
		"random_suffix":   RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccCheckApigeeSharedFlowDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApigeeSharedFlow_apigeeSharedflowTestExample(context),
			},
			{
				ResourceName:            "google_apigee_sharedflow.test_apigee_sharedflow",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"config_bundle", "detect_md5hash", "md5hash"},
			},
			{
				Config: testAccApigeeSharedFlow_apigeeSharedflowTestExampleUpdate(context),
			},
			{
				ResourceName:            "google_apigee_sharedflow.test_apigee_sharedflow",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"config_bundle", "detect_md5hash", "md5hash"},
			},
		},
	})
}

func testAccApigeeSharedFlow_apigeeSharedflowTestExample(context map[string]interface{}) string {
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

resource "google_apigee_sharedflow" "test_apigee_sharedflow" {
  name            = "tf-test-apigee-sharedflow"
  org_id          = google_project.project.project_id
  config_bundle   = "./test-fixtures/apigee/apigee_sharedflow_bundle.zip"
  depends_on      = [google_apigee_organization.apigee_org]
}
`, context)
}

func testAccCheckApigeeSharedFlowDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_apigee_shared_flow" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := GoogleProviderConfig(t)

			url, err := replaceVarsForTest(config, rs, "{{ApigeeBasePath}}organizations/{{org_id}}/sharedflows/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}
			fmt.Printf("testAccCheckApigeeSharedFlowDestroyProducer, url %s", url)
			_, err = SendRequest(config, "GET", billingProject, url, config.UserAgent, nil)
			if err == nil {
				return fmt.Errorf("ApigeeSharedFlow still exists at %s", url)
			}
		}

		return nil
	}
}

func testAccApigeeSharedFlow_apigeeSharedflowTestExampleUpdate(context map[string]interface{}) string {
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

resource "google_apigee_sharedflow" "test_apigee_sharedflow" {
  name            = "tf-test-apigee-sharedflow"
  org_id          = google_project.project.project_id
  config_bundle   = "./test-fixtures/apigee/apigee_sharedflow_bundle2.zip"
  depends_on      = [google_apigee_organization.apigee_org]
}
`, context)
}
