package google

import (
	"context"
	"fmt"
	"strings"
	"testing"

	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	eventarc "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/eventarc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccEventarcTrigger_channel(t *testing.T) {
	t.Parallel()

	region := getTestRegionFromEnv()
	key1 := BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", region, "tf-bootstrap-key1")
	key2 := BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", region, "tf-bootstrap-key2")

	context := map[string]interface{}{
		"region":          region,
		"project_name":    getTestProjectFromEnv(),
		"service_account": getTestServiceAccountFromEnv(t),
		"key_ring":        GetResourceNameFromSelfLink(key1.KeyRing.Name),
		"key1":            GetResourceNameFromSelfLink(key1.CryptoKey.Name),
		"key2":            GetResourceNameFromSelfLink(key2.CryptoKey.Name),
		"random_suffix":   randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEventarcChannelTriggerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcTrigger_createTriggerWithChannelName(context),
			},
			{
				ResourceName:      "google_eventarc_trigger.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEventarcTrigger_createTriggerWithChannelName(context map[string]interface{}) string {
	return Nprintf(`
data "google_project" "test_project" {
	project_id  = "%{project_name}"
}

data "google_kms_key_ring" "test_key_ring" {
	name     = "%{key_ring}"
	location = "us-central1"
}

data "google_kms_crypto_key" "key1" {
	name     = "%{key1}"
	key_ring = data.google_kms_key_ring.test_key_ring.id
}


resource "google_kms_crypto_key_iam_binding" "key1_binding" {
	crypto_key_id = data.google_kms_crypto_key.key1.id
	role      = "roles/cloudkms.cryptoKeyEncrypterDecrypter"

	members = [
	"serviceAccount:service-${data.google_project.test_project.number}@gcp-sa-eventarc.iam.gserviceaccount.com",
	]
}

resource "google_eventarc_channel" "test_channel" {
	location = "%{region}"
	name     = "tf-test-channel%{random_suffix}"
	crypto_key_name =  data.google_kms_crypto_key.key1.id
	third_party_provider = "projects/${data.google_project.test_project.project_id}/locations/%{region}/providers/datadog"
	depends_on = [google_kms_crypto_key_iam_binding.key1_binding]
}

resource "google_cloud_run_service" "default" {
	name     = "tf-test-eventarc-service%{random_suffix}"
	location = "%{region}"

	metadata {
		namespace = "%{project_name}"
	}

	template {
		spec {
			containers {
				image = "gcr.io/cloudrun/hello"
				ports {
					container_port = 8080
				}
			}
			container_concurrency = 50
			timeout_seconds = 100
		}
	}

	traffic {
		percent         = 100
		latest_revision = true
	}
}

resource "google_eventarc_trigger" "primary" {
	name = "tf-test-trigger%{random_suffix}"
	location = "%{region}"
	matching_criteria {
		attribute = "type"
		value = "datadog.v1.alert"
	}
	destination {
		cloud_run_service {
			service = google_cloud_run_service.default.name
			region = "%{region}"
		}
	}
	service_account= "%{service_account}"

    channel = "projects/${data.google_project.test_project.project_id}/locations/%{region}/channels/${google_eventarc_channel.test_channel.name}"

    depends_on =[google_cloud_run_service.default,google_eventarc_channel.test_channel]
}
`, context)
}

func testAccCheckEventarcChannelTriggerDestroyProducer(t *testing.T) func(s *terraform.State) error {
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
				Channel:        dcl.StringOrNil(rs.Primary.Attributes["channel"]),
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
