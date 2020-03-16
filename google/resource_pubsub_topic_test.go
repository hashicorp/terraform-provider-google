package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPubsubTopic_update(t *testing.T) {
	t.Parallel()

	topic := fmt.Sprintf("tf-test-topic-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPubsubTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubTopic_update(topic, "foo", "bar"),
			},
			{
				ResourceName:      "google_pubsub_topic.foo",
				ImportStateId:     topic,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPubsubTopic_updateWithRegion(topic, "wibble", "wobble", "us-central1"),
			},
			{
				ResourceName:      "google_pubsub_topic.foo",
				ImportStateId:     topic,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	}, testAccCheckPubsubTopicDestroyProducer)
}

func TestAccPubsubTopic_cmek(t *testing.T) {
	t.Parallel()

	kms := BootstrapKMSKey(t)
	pid := getTestProjectFromEnv()
	topicName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPubsubTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubTopic_cmek(pid, topicName, kms.CryptoKey.Name),
			},
			{
				ResourceName:      "google_pubsub_topic.topic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPubsubTopic_update(topic, key, value string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
  name = "%s"
  labels = {
    %s = "%s"
  }
}
`, topic, key, value)
}

func testAccPubsubTopic_updateWithRegion(topic, key, value, region string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
  name = "%s"
  labels = {
    %s = "%s"
  }

  message_storage_policy {
    allowed_persistence_regions = [
      "%s",
    ]
  }
}
`, topic, key, value, region)
}

func testAccPubsubTopic_cmek(pid, topicName, kmsKey string) string {
	return fmt.Sprintf(`
data "google_project" "project" {
  project_id = "%s"
}

resource "google_project_iam_member" "kms-project-binding" {
  project = data.google_project.project.project_id
  role    = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member  = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-pubsub.iam.gserviceaccount.com"
}

resource "google_pubsub_topic" "topic" {
  name         = "%s"
  project      = google_project_iam_member.kms-project-binding.project
  kms_key_name = "%s"
}
`, pid, topicName, kmsKey)
}

// Temporary until all destroy functions can be reworked to take a provider as an argument
func testAccCheckPubsubTopicDestroyProducer(provider *schema.Provider) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_pubsub_topic" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := provider.Meta().(*Config)

			url, err := replaceVarsForTest(config, rs, "{{PubsubBasePath}}projects/{{project}}/topics/{{name}}")
			if err != nil {
				return err
			}

			_, err = sendRequest(config, "GET", "", url, nil, pubsubTopicProjectNotReady)
			if err == nil {
				return fmt.Errorf("PubsubTopic still exists at %s", url)
			}
		}

		return nil
	}
}
