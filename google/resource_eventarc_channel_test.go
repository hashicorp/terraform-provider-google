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

func TestAccEventarcChannel_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"region":        getTestRegionFromEnv(),
		"project_name":  getTestProjectFromEnv(),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEventarcChannelDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcChannel_basic(context),
			},
			{
				ResourceName:      "google_eventarc_channel.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccEventarcChannel_cryptoKeyUpdate(t *testing.T) {
	t.Parallel()

	region := getTestRegionFromEnv()
	key1 := BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", region, "tf-bootstrap-key1")
	key2 := BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", region, "tf-bootstrap-key2")

	context := map[string]interface{}{
		"region":        region,
		"project_name":  getTestProjectFromEnv(),
		"key_ring":      GetResourceNameFromSelfLink(key1.KeyRing.Name),
		"key1":          GetResourceNameFromSelfLink(key1.CryptoKey.Name),
		"key2":          GetResourceNameFromSelfLink(key2.CryptoKey.Name),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEventarcChannelDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcChannel_setCryptoKey(context),
			},
			{
				ResourceName:      "google_eventarc_channel.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEventarcChannel_cryptoKeyUpdate(context),
			},
			{
				ResourceName:      "google_eventarc_channel.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEventarcChannel_basic(context map[string]interface{}) string {
	return Nprintf(`
data "google_project" "test_project" {
	project_id  = "%{project_name}"
}

resource "google_eventarc_channel" "primary" {
	location = "%{region}"
	name     = "tf-test-name%{random_suffix}"
	third_party_provider = "projects/${data.google_project.test_project.project_id}/locations/%{region}/providers/datadog"
}
`, context)
}

func testAccEventarcChannel_setCryptoKey(context map[string]interface{}) string {
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

resource "google_eventarc_channel" "primary" {
	location = "%{region}"
	name     = "tf-test-name%{random_suffix}"
	crypto_key_name =  data.google_kms_crypto_key.key1.id
	third_party_provider = "projects/${data.google_project.test_project.project_id}/locations/%{region}/providers/datadog"
	depends_on = [google_kms_crypto_key_iam_binding.key1_binding]
}
`, context)
}

func testAccEventarcChannel_cryptoKeyUpdate(context map[string]interface{}) string {
	return Nprintf(`
data "google_project" "test_project" {
	project_id  = "%{project_name}"
}

data "google_kms_key_ring" "test_key_ring" {
	name     = "%{key_ring}"
	location = "us-central1"
}
	
data "google_kms_crypto_key" "key2" {
	name     = "%{key2}"
	key_ring = data.google_kms_key_ring.test_key_ring.id
}

resource "google_kms_crypto_key_iam_binding" "key2_binding" {
	crypto_key_id = data.google_kms_crypto_key.key2.id
	role      = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
	
	members = [
	"serviceAccount:service-${data.google_project.test_project.number}@gcp-sa-eventarc.iam.gserviceaccount.com",
	]
}

resource "google_eventarc_channel" "primary" {
	location = "%{region}"
	name     = "tf-test-name%{random_suffix}"
	crypto_key_name= data.google_kms_crypto_key.key2.id
	third_party_provider = "projects/${data.google_project.test_project.project_id}/locations/%{region}/providers/datadog"
	depends_on = [google_kms_crypto_key_iam_binding.key2_binding]
}
`, context)
}

func testAccCheckEventarcChannelDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "rs.google_eventarc_channel" {
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

			obj := &eventarc.Channel{
				Location:           dcl.String(rs.Primary.Attributes["location"]),
				Name:               dcl.String(rs.Primary.Attributes["name"]),
				CryptoKeyName:      dcl.String(rs.Primary.Attributes["crypto_key_name"]),
				Project:            dcl.StringOrNil(rs.Primary.Attributes["project"]),
				ThirdPartyProvider: dcl.String(rs.Primary.Attributes["third_party_provider"]),
				ActivationToken:    dcl.StringOrNil(rs.Primary.Attributes["activation_token"]),
				CreateTime:         dcl.StringOrNil(rs.Primary.Attributes["create_time"]),
				PubsubTopic:        dcl.StringOrNil(rs.Primary.Attributes["pubsub_topic"]),
				State:              eventarc.ChannelStateEnumRef(rs.Primary.Attributes["state"]),
				Uid:                dcl.StringOrNil(rs.Primary.Attributes["uid"]),
				UpdateTime:         dcl.StringOrNil(rs.Primary.Attributes["update_time"]),
			}

			client := NewDCLEventarcClient(config, config.userAgent, billingProject, 0)
			_, err := client.GetChannel(context.Background(), obj)
			if err == nil {
				return fmt.Errorf("google_eventarc_channel still exists %v", obj)
			}
		}
		return nil
	}
}
