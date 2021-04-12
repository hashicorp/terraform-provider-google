package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccEventarcTrigger_basic(t *testing.T) {
	// DCL currently fails due to transport modification
	skipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
		"project":       getTestProjectFromEnv(),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: funcAccTestEventarcTriggerCheckDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcTrigger_basic(context),
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "google_eventarc_trigger.trigger",
			},
		},
	})
}

func TestAccEventarcTrigger_transport(t *testing.T) {
	// DCL currently fails due to transport modification
	skipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
		"project":       getTestProjectFromEnv(),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: funcAccTestEventarcTriggerCheckDestroy(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcTrigger_transport(context),
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "google_eventarc_trigger.trigger",
			},
			{
				Config: testAccEventarcTrigger_transportUpdate(context),
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "google_eventarc_trigger.trigger",
			},
			{
				Config: testAccEventarcTrigger_transport(context),
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "google_eventarc_trigger.trigger",
			},
		},
	})
}

func testAccEventarcTrigger_basic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_eventarc_trigger" "trigger" {
	name = "trigger%{random_suffix}"
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
	name = "topic%{random_suffix}"
}

resource "google_cloud_run_service" "default" {
	name     = "service-eventarc%{random_suffix}"
	location = "europe-west1"

	metadata {
		namespace = "%{project}"
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

func testAccEventarcTrigger_transport(context map[string]interface{}) string {
	return Nprintf(`
resource "google_eventarc_trigger" "trigger" {
	name = "trigger%{random_suffix}"
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
	name = "topic%{random_suffix}"
}

resource "google_cloud_run_service" "default" {
	name     = "service-eventarc%{random_suffix}"
	location = "europe-west1"

	metadata {
		namespace = "%{project}"
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
	name     = "service-eventarc2%{random_suffix}"
	location = "europe-north1"

	metadata {
		namespace = "%{project}"
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

func testAccEventarcTrigger_transportUpdate(context map[string]interface{}) string {
	return Nprintf(`
resource "google_eventarc_trigger" "trigger" {
	name = "trigger%{random_suffix}"
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
	account_id   = "sa%{random_suffix}"
	display_name = "Test Service Account"
}

resource "google_pubsub_topic" "foo" {
	name = "topic%{random_suffix}"
}

resource "google_cloud_run_service" "default" {
	name     = "service-eventarc%{random_suffix}"
	location = "europe-west1"

	metadata {
		namespace = "%{project}"
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
	name     = "service-eventarc2%{random_suffix}"
	location = "europe-north1"

	metadata {
		namespace = "%{project}"
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

func funcAccTestEventarcTriggerCheckDestroy(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_eventarc_trigger" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := googleProviderConfig(t)

			url, err := replaceVarsForTest(config, rs, "{{EventarcBasePath}}projects/{{project}}/locations/{{location}}/triggers/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = sendRequest(config, "GET", billingProject, url, config.userAgent, nil)
			if err == nil {
				return fmt.Errorf("EventarcTrigger still exists at %s", url)
			}
		}

		return nil
	}
}
