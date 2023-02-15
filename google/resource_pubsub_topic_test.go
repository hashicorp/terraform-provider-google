package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPubsubTopic_update(t *testing.T) {
	t.Parallel()

	topic := fmt.Sprintf("tf-test-topic-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPubsubTopicDestroyProducer(t),
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
	})
}

func TestAccPubsubTopic_cmek(t *testing.T) {
	t.Parallel()

	kms := BootstrapKMSKey(t)
	topicName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPubsubTopicDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubTopic_cmek(topicName, kms.CryptoKey.Name),
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

func testAccPubsubTopic_cmek(topicName, kmsKey string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "topic" {
  name         = "%s"
  kms_key_name = "%s"
}
`, topicName, kmsKey)
}
