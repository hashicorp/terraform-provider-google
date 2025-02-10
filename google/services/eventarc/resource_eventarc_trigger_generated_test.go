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

package eventarc_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccEventarcTrigger_eventarcTriggerWithCloudRunDestinationExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckEventarcTriggerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcTrigger_eventarcTriggerWithCloudRunDestinationExample(context),
			},
			{
				ResourceName:            "google_eventarc_trigger.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "terraform_labels"},
			},
		},
	})
}

func testAccEventarcTrigger_eventarcTriggerWithCloudRunDestinationExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_eventarc_trigger" "primary" {
  name     = "tf-test-some-trigger%{random_suffix}"
  location = "us-central1"
  matching_criteria {
    attribute = "type"
    value     = "google.cloud.pubsub.topic.v1.messagePublished"
  }
  destination {
    cloud_run_service {
      service = google_cloud_run_service.default.name
      region  = "us-central1"
    }
  }
  labels = {
    foo = "bar"
  }
  transport {
    pubsub {
      topic = google_pubsub_topic.foo.id
    }
  }
}

resource "google_pubsub_topic" "foo" {
  name = "tf-test-some-topic%{random_suffix}"
}

resource "google_cloud_run_service" "default" {
  name     = "tf-test-some-service%{random_suffix}"
  location = "us-central1"

  template {
    spec {
      containers {
        image = "gcr.io/cloudrun/hello"
        ports {
          container_port = 8080
        }
      }
      container_concurrency = 50
      timeout_seconds       = 100
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}
`, context)
}

func TestAccEventarcTrigger_eventarcTriggerWithHttpDestinationExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_id":              envvar.GetTestProjectFromEnv(),
		"service_account":         envvar.GetTestServiceAccountFromEnv(t),
		"network_attachment_name": acctest.BootstrapNetworkAttachment(t, "tf-test-eventarc-trigger-na", acctest.BootstrapSubnet(t, "tf-test-eventarc-trigger-subnet", acctest.BootstrapSharedTestNetwork(t, "tf-test-eventarc-trigger-network"))),
		"random_suffix":           acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckEventarcTriggerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcTrigger_eventarcTriggerWithHttpDestinationExample(context),
			},
			{
				ResourceName:            "google_eventarc_trigger.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "terraform_labels"},
			},
		},
	})
}

func testAccEventarcTrigger_eventarcTriggerWithHttpDestinationExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_eventarc_trigger" "primary" {
  name     = "tf-test-some-trigger%{random_suffix}"
  location = "us-central1"
  matching_criteria {
    attribute = "type"
    value     = "google.cloud.pubsub.topic.v1.messagePublished"
  }
  destination {
    http_endpoint {
      uri = "http://10.77.0.0:80/route"
    }
    network_config {
      network_attachment = "projects/%{project_id}/regions/us-central1/networkAttachments/%{network_attachment_name}"
    }
  }
  service_account = "%{service_account}"
}
`, context)
}

func TestAccEventarcTrigger_eventarcTriggerWithChannelCmekExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_id":      envvar.GetTestProjectFromEnv(),
		"project_number":  envvar.GetTestProjectNumberFromEnv(),
		"service_account": envvar.GetTestServiceAccountFromEnv(t),
		"key_name":        acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us-central1", "tf-bootstrap-eventarc-trigger-key").CryptoKey.Name,
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckEventarcTriggerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcTrigger_eventarcTriggerWithChannelCmekExample(context),
			},
			{
				ResourceName:            "google_eventarc_trigger.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "terraform_labels"},
			},
		},
	})
}

func testAccEventarcTrigger_eventarcTriggerWithChannelCmekExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_kms_crypto_key_iam_member" "key_member" {
  crypto_key_id = "%{key_name}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-%{project_number}@gcp-sa-eventarc.iam.gserviceaccount.com"
}

resource "google_eventarc_channel" "test_channel" {
  location             = "us-central1"
  name                 = "tf-test-some-channel%{random_suffix}"
  crypto_key_name      = "%{key_name}"
  third_party_provider = "projects/%{project_id}/locations/us-central1/providers/datadog"
  depends_on           = [google_kms_crypto_key_iam_member.key_member]
}

resource "google_cloud_run_service" "default" {
  name     = "tf-test-some-service%{random_suffix}"
  location = "us-central1"

  template {
    spec {
      containers {
        image = "gcr.io/cloudrun/hello"
        ports {
          container_port = 8080
        }
      }
      container_concurrency = 50
      timeout_seconds       = 100
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

resource "google_eventarc_trigger" "primary" {
  name     = "tf-test-some-trigger%{random_suffix}"
  location = "us-central1"
  matching_criteria {
    attribute = "type"
    value     = "datadog.v1.alert"
  }
  destination {
    cloud_run_service {
      service = google_cloud_run_service.default.name
      region  = "us-central1"
    }
  }
  service_account = "%{service_account}"
  channel         = google_eventarc_channel.test_channel.id
}
`, context)
}

func TestAccEventarcTrigger_eventarcTriggerWithWorkflowDestinationExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"service_account": envvar.GetTestServiceAccountFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckEventarcTriggerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcTrigger_eventarcTriggerWithWorkflowDestinationExample(context),
			},
			{
				ResourceName:            "google_eventarc_trigger.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "terraform_labels"},
			},
		},
	})
}

func testAccEventarcTrigger_eventarcTriggerWithWorkflowDestinationExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_eventarc_trigger" "primary" {
  name     = "tf-test-some-trigger%{random_suffix}"
  location = "us-central1"
  matching_criteria {
    attribute = "type"
    value     = "google.cloud.pubsub.topic.v1.messagePublished"
  }
  destination {
    workflow = google_workflows_workflow.workflow.id
  }
  service_account = "%{service_account}"
}

resource "google_workflows_workflow" "workflow" {
  name                = "tf-test-some-workflow%{random_suffix}"
  deletion_protection = false
  region              = "us-central1"
  source_contents     = <<-EOF
  # This is a sample workflow, feel free to replace it with your source code
  #
  # This workflow does the following:
  # - reads current time and date information from an external API and stores
  #   the response in CurrentDateTime variable
  # - retrieves a list of Wikipedia articles related to the day of the week
  #   from CurrentDateTime
  # - returns the list of articles as an output of the workflow
  # FYI, In terraform you need to escape the $$ or it will cause errors.

  - getCurrentTime:
      call: http.get
      args:
          url: $${sys.get_env("url")}
      result: CurrentDateTime
  - readWikipedia:
      call: http.get
      args:
          url: https://en.wikipedia.org/w/api.php
          query:
              action: opensearch
              search: $${CurrentDateTime.body.dayOfTheWeek}
      result: WikiResult
  - returnOutput:
      return: $${WikiResult.body[1]}
EOF
}
`, context)
}

func TestAccEventarcTrigger_eventarcTriggerWithPathPatternFilterExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckEventarcTriggerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcTrigger_eventarcTriggerWithPathPatternFilterExample(context),
			},
			{
				ResourceName:            "google_eventarc_trigger.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "terraform_labels"},
			},
		},
	})
}

func testAccEventarcTrigger_eventarcTriggerWithPathPatternFilterExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_eventarc_trigger" "primary" {
  name     = "tf-test-some-trigger%{random_suffix}"
  location = "us-central1"
  matching_criteria {
    attribute = "type"
    value     = "google.cloud.eventarc.trigger.v1.created"
  }
  matching_criteria {
    attribute = "trigger"
    operator  = "match-path-pattern"
    value     = "trigger-with-wildcard-*"
  }
  destination {
    cloud_run_service {
      service = google_cloud_run_service.default.name
      region  = "us-central1"
    }
  }
  labels = {
    foo = "bar"
  }
  event_data_content_type = "application/protobuf"
  service_account         = google_service_account.trigger_service_account.email
  depends_on              = [google_project_iam_member.event_receiver]
}

resource "google_service_account" "trigger_service_account" {
  account_id = "tf-test-trigger-sa%{random_suffix}"
}

resource "google_project_iam_member" "event_receiver" {
  project = google_service_account.trigger_service_account.project
  role    = "roles/eventarc.eventReceiver"
  member  = "serviceAccount:${google_service_account.trigger_service_account.email}"
}

resource "google_cloud_run_service" "default" {
  name     = "tf-test-some-service%{random_suffix}"
  location = "us-central1"

  template {
    spec {
      containers {
        image = "gcr.io/cloudrun/hello"
        ports {
          container_port = 8080
        }
      }
      container_concurrency = 50
      timeout_seconds       = 100
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}
`, context)
}

func testAccCheckEventarcTriggerDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_eventarc_trigger" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{EventarcBasePath}}projects/{{project}}/locations/{{location}}/triggers/{{name}}")
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
				return fmt.Errorf("EventarcTrigger still exists at %s", url)
			}
		}

		return nil
	}
}
