// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package apigee_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccApigeeTargetServer_apigeeTargetServerTest_update(t *testing.T) {
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
			"time": {},
		},
		CheckDestroy: testAccCheckApigeeTargetServerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApigeeTargetServer_apigeeTargetServerTest_basic(context),
			},
			{
				ResourceName:            "google_apigee_target_server.apigee_target_server",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"env_id"},
			},
			{
				Config: testAccApigeeTargetServer_apigeeTargetServerTest_createWithoutProtocol(context),
			},
			{
				ResourceName:            "google_apigee_target_server.apigee_target_server",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"env_id"},
			},
		},
	})
}

func testAccApigeeTargetServer_apigeeTargetServerTest_basic(context map[string]interface{}) string {
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
  project    = google_project.project.project_id
  service    = "servicenetworking.googleapis.com"
  depends_on = [google_project_service.apigee]
}

resource "google_project_service" "compute" {
  project    = google_project.project.project_id
  service    = "compute.googleapis.com"
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

resource "google_apigee_target_server" "apigee_target_server" {
  name        = "tf-test-target-server%{random_suffix}"
  description = "Apigee Target Server"
  protocol    = "HTTP"
  host        = "abc.foo.com"
  port        = 8080
  env_id      = google_apigee_environment.apigee_environment.id
}
`, context)
}

func testAccApigeeTargetServer_apigeeTargetServerTest_createWithoutProtocol(context map[string]interface{}) string {
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
  project    = google_project.project.project_id
  service    = "servicenetworking.googleapis.com"
  depends_on = [google_project_service.apigee]
}

resource "google_project_service" "compute" {
  project    = google_project.project.project_id
  service    = "compute.googleapis.com"
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

resource "google_apigee_target_server" "apigee_target_server" {
  name        = "tf-test-target-server%{random_suffix}"
  description = "Apigee Target Server"
  host        = "abc.foo.com"
  port        = 8080
  env_id      = google_apigee_environment.apigee_environment.id
}
`, context)
}

func TestAccApigeeTargetServer_apigeeTargetServerTest_clientAuthEnabled(t *testing.T) {
	t.Parallel()
	// Skipping VCR tests; google_apigee_keystores_aliases_key_cert_file resource uses multipart boundary which by default is random. Currently this is incompatible with VCR.
	acctest.SkipIfVcr(t)
	context := map[string]interface{}{
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckApigeeTargetServerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApigeeTargetServer_apigeeTargetServerTest_full(context),
			},
			{
				ResourceName:            "google_apigee_target_server.apigee_target_server",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"env_id"},
			},
			{
				Config: testAccApigeeTargetServer_apigeeTargetServerTest_update(context),
			},
			{
				ResourceName:            "google_apigee_target_server.apigee_target_server",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"env_id"},
			},
		},
	})
}

func testAccApigeeTargetServer_apigeeTargetServerTest_full(context map[string]interface{}) string {
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
  project    = google_project.project.project_id
  service    = "servicenetworking.googleapis.com"
  depends_on = [google_project_service.apigee]
}

resource "google_project_service" "compute" {
  project    = google_project.project.project_id
  service    = "compute.googleapis.com"
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

resource "google_apigee_env_keystore" "apigee_environment_keystore" {
	name   = "tf-test-keystore%{random_suffix}"
  env_id = google_apigee_environment.apigee_environment.id
}
   
resource "google_apigee_keystores_aliases_key_cert_file" "apigee_test_alias" {
	alias       = "tf-test-alias%{random_suffix}"
	org_id      = google_apigee_organization.apigee_org.name
	environment = google_apigee_environment.apigee_environment.name
	cert        = file("./test-fixtures/apigee_keystore_alias_test_cert.pem")
	key         = sensitive(file("./test-fixtures/apigee_keystore_alias_test_key.pem"))
	password    = sensitive("password")
	keystore    = google_apigee_env_keystore.apigee_environment_keystore.name
}  
   
resource "google_apigee_target_server" "apigee_target_server"{
	env_id                    = google_apigee_environment.apigee_environment.id
	is_enabled                = true
	name                      = "tf-test-target-server%{random_suffix}"
	host                      = "host.test.com"
	port                      = 443
	protocol                  = "HTTP"
	s_sl_info{
	 enabled                  = true
	 ciphers                  = ["TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA38411"]
	 client_auth_enabled      = true
	 ignore_validation_errors = true
	 key_alias                = google_apigee_keystores_aliases_key_cert_file.apigee_test_alias.alias
	 key_store                = google_apigee_env_keystore.apigee_environment_keystore.name
	 protocols                = ["TLSv1.1"]
	 trust_store              = google_apigee_env_keystore.apigee_environment_keystore.name
     enforce                  = false
   common_name{
    value                   = "testCn"
    wildcard_match          = true
   }
	}
	depends_on = [ 
    google_apigee_env_keystore.apigee_environment_keystore,
    google_apigee_keystores_aliases_key_cert_file.apigee_test_alias
  ]
	}
`, context)
}

func testAccApigeeTargetServer_apigeeTargetServerTest_update(context map[string]interface{}) string {
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
  project    = google_project.project.project_id
  service    = "servicenetworking.googleapis.com"
  depends_on = [google_project_service.apigee]
}

resource "google_project_service" "compute" {
  project    = google_project.project.project_id
  service    = "compute.googleapis.com"
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

resource "google_apigee_env_keystore" "apigee_environment_keystore2" {
	name   = "tf-test-keystore2%{random_suffix}"
	env_id = google_apigee_environment.apigee_environment.id
}
   
resource "google_apigee_keystores_aliases_key_cert_file" "apigee_test_alias2" {
	alias       = "tf-test-alias%{random_suffix}"
	org_id      = google_apigee_organization.apigee_org.name
	environment = google_apigee_environment.apigee_environment.name
	cert        = file("./test-fixtures/apigee_keystore_alias_test_cert2.pem")
	key         = sensitive(file("./test-fixtures/apigee_keystore_alias_test_key.pem"))
	password    = sensitive("password")
	keystore    = google_apigee_env_keystore.apigee_environment_keystore2.name
}  
   
resource "google_apigee_target_server" "apigee_target_server"{
	env_id                     = google_apigee_environment.apigee_environment.id
	is_enabled                 = true
	name                       = "tf-test-target-server%{random_suffix}"
	host                       = "host.test.com"
	port                       = 8443
	protocol                   = "GRPC"
  s_sl_info{
	  enabled                  = true
	  ciphers                  = ["TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384"]
	  client_auth_enabled      = true
	  ignore_validation_errors = true
	  key_alias                = google_apigee_keystores_aliases_key_cert_file.apigee_test_alias2.alias 
	  key_store                = google_apigee_env_keystore.apigee_environment_keystore2.name
	  protocols                = ["TLSv1.2", "TLSv1.1"]
	  trust_store              = google_apigee_env_keystore.apigee_environment_keystore2.name
      enforce                  = true
	}
	depends_on                 = [ 
    google_apigee_env_keystore.apigee_environment_keystore2,
    google_apigee_keystores_aliases_key_cert_file.apigee_test_alias2
  ]
}	
`, context)
}
