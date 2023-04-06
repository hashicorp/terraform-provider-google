package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccApigeeKeystoresAliasesKeyCertFile_apigeeKeystoresAliasesKeyCertFileTestExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":          GetTestOrgFromEnv(t),
		"billing_account": GetTestBillingAccountFromEnv(t),
		"random_suffix":   RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckApigeeKeystoresAliasesKeyCertFileDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApigeeKeystoresAliasesKeyCertFile_apigeeKeystoresAliasesKeyCertFileTestExample(context),
			},
			{
				ResourceName:            "google_apigee_keystores_aliases_key_cert_file.test-alias",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cert", "key", "environment", "keystore", "org_id", "password"},
			},
			{
				Config: testAccApigeeKeystoresAliasesKeyCertFile_apigeeKeystoresAliasesKeyCertFileTestExampleUpdate(context),
			},
			{
				ResourceName:            "google_apigee_keystores_aliases_key_cert_file.test-alias",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cert", "key", "environment", "keystore", "org_id", "password"},
			},
		},
	})
}

func testAccApigeeKeystoresAliasesKeyCertFile_apigeeKeystoresAliasesKeyCertFileTestExample(context map[string]interface{}) string {
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

resource "google_apigee_environment" "apigee_environment_keystore_alias" {
  org_id       = google_apigee_organization.apigee_org.id
  name         = "tf-test%{random_suffix}"
  description  = "Apigee Environment"
  display_name = "environment-1"
}

resource "google_apigee_env_keystore" "apigee_environment_keystore" {
  name       = "tf-test-keystore%{random_suffix}"
  env_id     = google_apigee_environment.apigee_environment_keystore_alias.id
}

resource "google_apigee_keystores_aliases_key_cert_file" "test-alias" {
	alias = "tf-test%{random_suffix}"
	org_id =  google_apigee_organization.apigee_org.name
	environment = google_apigee_environment.apigee_environment_keystore_alias.name
	cert = file("./test-fixtures/apigee/apigee_keystore_alias_test_cert.pem")
	key = sensitive(file("./test-fixtures/apigee/apigee_keystore_alias_test_key.pem"))
	password = sensitive("password")
	keystore = google_apigee_env_keystore.apigee_environment_keystore.name
}
`, context)
}

func testAccCheckApigeeKeystoresAliasesKeyCertFileDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_apigee_keystores_aliases_key_cert_file" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := GoogleProviderConfig(t)

			url, err := replaceVarsForTest(config, rs, "{{ApigeeBasePath}}organizations/{{org_id}}/environments/{{environment}}/keystores/{{keystore}}/aliases/{{alias}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = SendRequest(config, "GET", billingProject, url, config.UserAgent, nil)
			if err == nil {
				return fmt.Errorf("ApigeeKeystoresAliasesKeyCertFile still exists at %s", url)
			}
		}

		return nil
	}
}

func testAccApigeeKeystoresAliasesKeyCertFile_apigeeKeystoresAliasesKeyCertFileTestExampleUpdate(context map[string]interface{}) string {
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

resource "google_apigee_environment" "apigee_environment_keystore_alias" {
  org_id       = google_apigee_organization.apigee_org.id
  name         = "tf-test%{random_suffix}"
  description  = "Apigee Environment"
  display_name = "environment-1"
}

resource "google_apigee_env_keystore" "apigee_environment_keystore" {
  name       = "tf-test-keystore%{random_suffix}"
  env_id     = google_apigee_environment.apigee_environment_keystore_alias.id
}

resource "google_apigee_keystores_aliases_key_cert_file" "test-alias" {
	alias = "tf-test%{random_suffix}"
	org_id =  google_apigee_organization.apigee_org.name
	environment = google_apigee_environment.apigee_environment_keystore_alias.name
	cert = file("./test-fixtures/apigee/apigee_keystore_alias_test_cert2.pem")
	key = sensitive(file("./test-fixtures/apigee/apigee_keystore_alias_test_key.pem"))
	password = sensitive("password")
	keystore = google_apigee_env_keystore.apigee_environment_keystore.name
}  


`, context)
}
