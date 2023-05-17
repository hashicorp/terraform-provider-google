package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccApigeeKeystoresAliasesPkcs12_ApigeeKeystoresAliasesPkcs12Example(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":          acctest.GetTestOrgFromEnv(t),
		"billing_account": acctest.GetTestBillingAccountFromEnv(t),
		"random_suffix":   RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckApigeeKeystoresAliasesPkcs12DestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApigeeKeystoresAliasesPkcs12_ApigeeKeystoresAliasesPkcs12Example(context),
			},
			{
				ResourceName:            "google_apigee_keystores_aliases_pkcs12.apigee_environment_keystore_aliases_pkcs",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"file", "filehash", "password", "org_id", "environment"},
			},
		},
	})
}

func testAccApigeeKeystoresAliasesPkcs12_ApigeeKeystoresAliasesPkcs12Example(context map[string]interface{}) string {
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

resource "google_apigee_environment" "apigee_environment_keystore" {
  org_id       = google_apigee_organization.apigee_org.id
  name         = "tf-test%{random_suffix}"
  description  = "Apigee Environment"
  display_name = "environment-1"
}

resource "google_apigee_env_keystore" "apigee_environment_keystore_alias" {
  name       = "tf-test-keystore%{random_suffix}"
  env_id     = google_apigee_environment.apigee_environment_keystore.id
}

resource "google_apigee_keystores_aliases_pkcs12" "apigee_environment_keystore_aliases_pkcs" {
  environment 			= google_apigee_environment.apigee_environment_keystore.name
  org_id				= google_apigee_organization.apigee_org.name
  keystore				= google_apigee_env_keystore.apigee_environment_keystore_alias.name
  alias                 = "tf-test%{random_suffix}"
  file                  = "./test-fixtures/apigee/keyStore.p12"
  filehash				= filemd5("./test-fixtures/apigee/keyStore.p12")
  password              = sensitive("abcd")
}
`, context)
}

func testAccCheckApigeeKeystoresAliasesPkcs12DestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_apigee_env_keystore_aliases_pkcs" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{ApigeeBasePath}}{{keystore_id}}/aliases/{{alias}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("ApigeeEnvKeystoreAliasesPkcs still exists at %s", url)
			}
		}

		return nil
	}
}
