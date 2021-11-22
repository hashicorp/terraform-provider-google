// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: DCL     ***
//
// ----------------------------------------------------------------------------
//
//     This file is managed by Magic Modules (https://github.com/GoogleCloudPlatform/magic-modules)
//     and is based on the DCL (https://github.com/GoogleCloudPlatform/declarative-resource-client-library).
//     Changes will need to be made to the DCL or Magic Modules instead of here.
//
//     We are not currently able to accept contributions to this file. If changes
//     are required, please file an issue at https://github.com/hashicorp/terraform-provider-google/issues/new/choose
//
// ----------------------------------------------------------------------------

package google

import (
	"context"
	"fmt"
	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	eventarc "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/eventarc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strings"
	"testing"
)

func TestAccEventarcTrigger_BasicHandWritten(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  getTestProjectFromEnv(),
		"region":        getTestRegionFromEnv(),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEventarcTriggerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcTrigger_BasicHandWritten(context),
			},
			{
				ResourceName:      "google_eventarc_trigger.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEventarcTrigger_BasicHandWrittenUpdate0(context),
			},
			{
				ResourceName:      "google_eventarc_trigger.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEventarcTrigger_BasicHandWrittenUpdate1(context),
			},
			{
				ResourceName:      "google_eventarc_trigger.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEventarcTrigger_BasicHandWritten(context map[string]interface{}) string {
	return Nprintf(`
resource "google_eventarc_trigger" "primary" {
	name = "tf-test-name%{random_suffix}"
	location = "europe-west1"
	matching_criteria {
		attribute = "type"
		value = "google.cloud.pubsub.topic.v1.messagePublished"
	}
	destination {
		cloud_run_service {
			service = google_cloud_run_service.default.name
			region = "europe-west1"
		}
	}
	labels = {
		foo = "bar"
	}
}

resource "google_pubsub_topic" "foo" {
	name = "tf-test-topic%{random_suffix}"
}

resource "google_cloud_run_service" "default" {
	name     = "tf-test-eventarc-service%{random_suffix}"
	location = "europe-west1"

	metadata {
		namespace = "%{project_name}"
	}

	template {
		spec {
			containers {
				image = "gcr.io/cloudrun/hello"
				args  = ["arrgs"]
			}
		container_concurrency = 50
		}
	}

	traffic {
		percent         = 100
		latest_revision = true
	}
}

`, context)
}

func testAccEventarcTrigger_BasicHandWrittenUpdate0(context map[string]interface{}) string {
	return Nprintf(`
resource "google_eventarc_trigger" "primary" {
	name = "tf-test-name%{random_suffix}"
	location = "europe-west1"
	matching_criteria {
		attribute = "type"
		value = "google.cloud.pubsub.topic.v1.messagePublished"
	}
	destination {
		cloud_run_service {
			service = google_cloud_run_service.default.name
			region = "europe-west1"
		}
	}
	transport {
		pubsub {
			topic = google_pubsub_topic.foo.id
		}
	}
}

resource "google_pubsub_topic" "foo" {
	name = "tf-test-topic%{random_suffix}"
}

resource "google_cloud_run_service" "default" {
	name     = "tf-test-eventarc-service%{random_suffix}"
	location = "europe-west1"

	metadata {
		namespace = "%{project_name}"
	}

	template {
		spec {
			containers {
				image = "gcr.io/cloudrun/hello"
				args  = ["arrgs"]
			}
		container_concurrency = 50
		}
	}

	traffic {
		percent         = 100
		latest_revision = true
	}
}

resource "google_cloud_run_service" "default2" {
	name     = "tf-test-eventarc-service%{random_suffix}2"
	location = "europe-north1"

	metadata {
		namespace = "%{project_name}"
	}

	template {
		spec {
			containers {
				image = "gcr.io/cloudrun/hello"
				args  = ["arrgs"]
			}
		container_concurrency = 50
		}
	}

	traffic {
		percent         = 100
		latest_revision = true
	}
}

`, context)
}

func testAccEventarcTrigger_BasicHandWrittenUpdate1(context map[string]interface{}) string {
	return Nprintf(`
resource "google_eventarc_trigger" "primary" {
	name = "tf-test-name%{random_suffix}"
	location = "europe-west1"
	matching_criteria {
		attribute = "type"
		value = "google.cloud.pubsub.topic.v1.messagePublished"
	}
	destination {
		cloud_run_service {
			service = google_cloud_run_service.default2.name
			region = "europe-north1"
		}
	}
	transport {
		pubsub {
			topic = google_pubsub_topic.foo.id
		}
	}
	labels = {
		foo = "bar"
	}
	service_account = google_service_account.eventarc-sa.email
}

resource "google_service_account" "eventarc-sa" {
	account_id   = "tf-test-sa%{random_suffix}"
	display_name = "Test Service Account"
}

resource "google_pubsub_topic" "foo" {
	name = "tf-test-topic%{random_suffix}"
}

resource "google_cloud_run_service" "default" {
	name     = "tf-test-eventarc-service%{random_suffix}"
	location = "europe-west1"

	metadata {
		namespace = "%{project_name}"
	}

	template {
		spec {
			containers {
				image = "gcr.io/cloudrun/hello"
				args  = ["arrgs"]
			}
		container_concurrency = 50
		}
	}

	traffic {
		percent         = 100
		latest_revision = true
	}
}

resource "google_cloud_run_service" "default2" {
	name     = "tf-test-eventarc-service%{random_suffix}2"
	location = "europe-north1"

	metadata {
		namespace = "%{project_name}"
	}

	template {
		spec {
			containers {
				image = "gcr.io/cloudrun/hello"
				args  = ["arrgs"]
			}
		container_concurrency = 50
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
			if rs.Type != "rs.google_eventarc_trigger" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := googleProviderConfig(t)

			billingProject := ""
			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			obj := &eventarc.Trigger{
				Location:       dcl.String(rs.Primary.Attributes["location"]),
				Name:           dcl.String(rs.Primary.Attributes["name"]),
				Project:        dcl.StringOrNil(rs.Primary.Attributes["project"]),
				ServiceAccount: dcl.String(rs.Primary.Attributes["service_account"]),
				CreateTime:     dcl.StringOrNil(rs.Primary.Attributes["create_time"]),
				Etag:           dcl.StringOrNil(rs.Primary.Attributes["etag"]),
				Uid:            dcl.StringOrNil(rs.Primary.Attributes["uid"]),
				UpdateTime:     dcl.StringOrNil(rs.Primary.Attributes["update_time"]),
			}

			client := NewDCLEventarcClient(config, config.userAgent, billingProject, 0)
			_, err := client.GetTrigger(context.Background(), obj)
			if err == nil {
				return fmt.Errorf("google_eventarc_trigger still exists %v", obj)
			}
		}
		return nil
	}
}
