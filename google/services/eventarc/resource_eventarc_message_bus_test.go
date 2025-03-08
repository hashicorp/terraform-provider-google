// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package eventarc_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// We make sure not to run tests in parallel, since only one MessageBus per project is supported.
func TestAccEventarcMessageBus(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		"basic":           testAccEventarcMessageBus_basic,
		"cryptoKey":       testAccEventarcMessageBus_cryptoKey,
		"update":          testAccEventarcMessageBus_update,
		"googleApiSource": testAccEventarcMessageBus_googleApiSource,
		"pipeline":        testAccEventarcMessageBus_pipeline,
	}

	for name, tc := range testCases {
		// shadow the tc variable into scope so that when
		// the loop continues, if t.Run hasn't executed tc(t)
		// yet, we don't have a race condition
		// see https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc(t)
		})
	}
}

func testAccEventarcMessageBus_basic(t *testing.T) {
	context := map[string]interface{}{
		"region":        envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}
	acctest.BootstrapIamMembers(t, []acctest.IamMember{
		{
			Member: "serviceAccount:service-{project_number}@gcp-sa-eventarc.iam.gserviceaccount.com",
			Role:   "roles/cloudkms.cryptoKeyEncrypterDecrypter",
		},
	})

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckEventarcMessageBusDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcMessageBus_basicCfg(context),
			},
			{
				ResourceName:            "google_eventarc_message_bus.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels", "annotations"},
			},
		},
	})
}

func testAccEventarcMessageBus_basicCfg(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_eventarc_message_bus" "primary" {
  location       = "%{region}"
  message_bus_id = "tf-test-messagebus%{random_suffix}"
  display_name   = "basic bus"
  labels = {
    test_label = "test-eventarc-label"
  }
  annotations = {
    test_annotation = "test-eventarc-annotation"
  }
}
`, context)
}

func testAccEventarcMessageBus_cryptoKey(t *testing.T) {
	region := envvar.GetTestRegionFromEnv()
	context := map[string]interface{}{
		"project_number": envvar.GetTestProjectNumberFromEnv(),
		"region":         region,
		"key":            acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", region, "tf-bootstrap-eventarc-messagebus-key").CryptoKey.Name,
		"random_suffix":  acctest.RandString(t, 10),
	}
	acctest.BootstrapIamMembers(t, []acctest.IamMember{
		{
			Member: "serviceAccount:service-{project_number}@gcp-sa-eventarc.iam.gserviceaccount.com",
			Role:   "roles/cloudkms.cryptoKeyEncrypterDecrypter",
		},
	})

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckEventarcMessageBusDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcMessageBus_cryptoKeyCfg(context),
			},
			{
				ResourceName:            "google_eventarc_message_bus.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels", "annotations"},
			},
		},
	})
}

func testAccEventarcMessageBus_cryptoKeyCfg(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_eventarc_message_bus" "primary" {
  location        = "%{region}"
  message_bus_id  = "tf-test-messagebus%{random_suffix}"
  crypto_key_name = "%{key}"
  logging_config {
    log_severity = "ALERT"
  }
}
`, context)
}

func testAccEventarcMessageBus_update(t *testing.T) {
	region := envvar.GetTestRegionFromEnv()
	context := map[string]interface{}{
		"project_number": envvar.GetTestProjectNumberFromEnv(),
		"region":         region,
		"key1":           acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", region, "tf-bootstrap-eventarc-messagebus-key1").CryptoKey.Name,
		"key2":           acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", region, "tf-bootstrap-eventarc-messagebus-key2").CryptoKey.Name,
		"random_suffix":  acctest.RandString(t, 10),
	}
	acctest.BootstrapIamMembers(t, []acctest.IamMember{
		{
			Member: "serviceAccount:service-{project_number}@gcp-sa-eventarc.iam.gserviceaccount.com",
			Role:   "roles/cloudkms.cryptoKeyEncrypterDecrypter",
		},
	})

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckEventarcMessageBusDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcMessageBus_setCfg(context),
			},
			{
				ResourceName:            "google_eventarc_message_bus.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels", "annotations"},
			},
			{
				Config: testAccEventarcMessageBus_updateCfg(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_eventarc_message_bus.primary", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_eventarc_message_bus.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels", "annotations"},
			},
			{
				Config: testAccEventarcMessageBus_deleteCfg(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_eventarc_message_bus.primary", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_eventarc_message_bus.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels", "annotations"},
			},
		},
	})
}

func testAccEventarcMessageBus_setCfg(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_eventarc_message_bus" "primary" {
  location        = "%{region}"
  message_bus_id  = "tf-test-messagebus%{random_suffix}"
  crypto_key_name = "%{key1}"
  display_name    = "message bus"
  logging_config {
    log_severity = "ALERT"
  }
}
`, context)
}

func testAccEventarcMessageBus_updateCfg(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_eventarc_message_bus" "primary" {
  location        = "%{region}"
  message_bus_id  = "tf-test-messagebus%{random_suffix}"
  crypto_key_name = "%{key2}"
  display_name    = "updated message bus"
  logging_config {
    log_severity = "DEBUG"
  }
}
`, context)
}

func testAccEventarcMessageBus_deleteCfg(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_eventarc_message_bus" "primary" {
  location        = "%{region}"
  message_bus_id  = "tf-test-messagebus%{random_suffix}"
  crypto_key_name = ""
  display_name    = "updated message bus"
  logging_config {
    log_severity = "DEBUG"
  }
}
`, context)
}

// Although this test is defined in resource_eventarc_message_bus_test, it is primarily
// concerned with testing the GoogleApiSource resource, which depends on a singleton MessageBus.
func testAccEventarcMessageBus_googleApiSource(t *testing.T) {
	region := envvar.GetTestRegionFromEnv()
	context := map[string]interface{}{
		"project_number": envvar.GetTestProjectNumberFromEnv(),
		"key1":           acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", region, "tf-bootstrap-eventarc-googleapisource-key1").CryptoKey.Name,
		"region":         region,
		"random_suffix":  acctest.RandString(t, 10),
	}
	acctest.BootstrapIamMembers(t, []acctest.IamMember{
		{
			Member: "serviceAccount:service-{project_number}@gcp-sa-eventarc.iam.gserviceaccount.com",
			Role:   "roles/cloudkms.cryptoKeyEncrypterDecrypter",
		},
	})

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckEventarcGoogleApiSourceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcMessageBus_googleApiSourceCfg(context),
			},
			{
				ResourceName:            "google_eventarc_google_api_source.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels", "annotations"},
			},
		},
	})
}

func testAccEventarcMessageBus_googleApiSourceCfg(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_eventarc_google_api_source" "primary" {
  location             = "%{region}"
  google_api_source_id = "tf-test-googleapisource%{random_suffix}"
  display_name         = "basic google api source"
  destination          = google_eventarc_message_bus.message_bus.id
  crypto_key_name      = "%{key1}"
  labels = {
    test_label = "test-eventarc-label"
  }
  annotations = {
    test_annotation = "test-eventarc-annotation"
  }
  logging_config {
    log_severity = "DEBUG"
  }
}
resource "google_eventarc_message_bus" "message_bus" {
  location       = "%{region}"
  message_bus_id = "tf-test-messagebus%{random_suffix}"
}
`, context)
}

// Although this test is defined in resource_eventarc_message_bus_test, it is primarily
// concerned with testing the Pipeline resource, which depends on a singleton MessageBus.
func testAccEventarcMessageBus_pipeline(t *testing.T) {
	context := map[string]interface{}{
		"project_id":              envvar.GetTestProjectFromEnv(),
		"region":                  envvar.GetTestRegionFromEnv(),
		"random_suffix":           acctest.RandString(t, 10),
		"network_attachment_name": acctest.BootstrapNetworkAttachment(t, "tf-test-eventarc-messagebus-na", acctest.BootstrapSubnet(t, "tf-test-eventarc-messagebus-subnet", acctest.BootstrapSharedTestNetwork(t, "tf-test-eventarc-messagebus-network"))),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckEventarcPipelineDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcMessageBus_pipelineCfg(context),
			},
			{
				ResourceName:            "google_eventarc_pipeline.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels", "annotations"},
			},
		},
	})
}

func testAccEventarcMessageBus_pipelineCfg(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_eventarc_pipeline" "primary" {
  location    = "%{region}"
  pipeline_id = "tf-test-some-pipeline%{random_suffix}"
  destinations {
    message_bus = google_eventarc_message_bus.primary.id
    network_config {
      network_attachment = "projects/%{project_id}/regions/%{region}/networkAttachments/%{network_attachment_name}"
    }
  }
}

resource "google_eventarc_message_bus" "primary" {
  location       = "%{region}"
  message_bus_id = "tf-test-messagebus%{random_suffix}"
}
`, context)
}

func testAccCheckEventarcMessageBusDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_eventarc_message_bus" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{EventarcBasePath}}projects/{{project}}/locations/{{location}}/messageBuses/{{name}}")
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
				return fmt.Errorf("EventarcMessageBus still exists at %s", url)
			}
		}

		return nil
	}
}
