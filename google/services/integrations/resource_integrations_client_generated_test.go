// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package integrations_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccIntegrationsClient_integrationsClientBasicExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIntegrationsClientDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationsClient_integrationsClientBasicExample(context),
			},
			{
				ResourceName:            "google_integrations_client.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cloud_kms_config", "create_sample_integrations", "create_sample_workflows", "location", "provision_gmek", "run_as_service_account"},
			},
		},
	})
}

func testAccIntegrationsClient_integrationsClientBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_integrations_client" "example" {
  location = "us-central1"
}
`, context)
}

func TestAccIntegrationsClient_integrationsClientFullExample(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIntegrationsClientDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationsClient_integrationsClientFullExample(context),
			},
			{
				ResourceName:            "google_integrations_client.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cloud_kms_config", "create_sample_integrations", "create_sample_workflows", "location", "provision_gmek", "run_as_service_account"},
			},
		},
	})
}

func testAccIntegrationsClient_integrationsClientFullExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "test_project" {
}

resource "google_kms_key_ring" "keyring" {
  name     = "tf-test-my-keyring%{random_suffix}"
  location = "us-east1"
}

resource "google_kms_crypto_key" "cryptokey" {
  name = "crypto-key-example"
  key_ring = google_kms_key_ring.keyring.id
  rotation_period = "7776000s"
}

resource "google_kms_crypto_key_version" "test_key" {
  crypto_key = google_kms_crypto_key.cryptokey.id
}

resource "google_service_account" "service_account" {
  account_id   = "tf-test-my-service-acc%{random_suffix}"
  display_name = "Service Account"
}

resource "google_integrations_client" "example" {
  location = "us-east1"
  create_sample_integrations = true
  run_as_service_account = google_service_account.service_account.email
  cloud_kms_config {
    kms_location = "us-east1"
    kms_ring = google_kms_key_ring.keyring.id
    key = google_kms_crypto_key.cryptokey.id
    key_version = google_kms_crypto_key_version.test_key.id
    kms_project_id = data.google_project.test_project.project_id
  }
}
`, context)
}

func TestAccIntegrationsClient_integrationsClientDeprecatedFieldsExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIntegrationsClientDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationsClient_integrationsClientDeprecatedFieldsExample(context),
			},
			{
				ResourceName:            "google_integrations_client.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cloud_kms_config", "create_sample_integrations", "create_sample_workflows", "location", "provision_gmek", "run_as_service_account"},
			},
		},
	})
}

func testAccIntegrationsClient_integrationsClientDeprecatedFieldsExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_integrations_client" "example" {
  location = "asia-south1"
  provision_gmek = true
  create_sample_workflows = true
}
`, context)
}

func testAccCheckIntegrationsClientDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_integrations_client" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{IntegrationsBasePath}}projects/{{project}}/locations/{{location}}/clients")
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
				return fmt.Errorf("IntegrationsClient still exists at %s", url)
			}
		}

		return nil
	}
}
