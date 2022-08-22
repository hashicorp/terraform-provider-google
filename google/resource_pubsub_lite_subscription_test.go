package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPubsubLiteSubscription_pubsubLiteSubscription_deliveryRequirement_update(t *testing.T) {
	t.Parallel()

	topic := fmt.Sprintf("tf-test-topic-foo-%s", randString(t, 10))
	subscription := fmt.Sprintf("tf-test-topic-foo-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPubsubLiteSubscriptionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubLiteSubscription_pubsubLiteSubscription_deliveryRequirement_update(topic, subscription, "DELIVER_AFTER_STORED"),
			},
			{
				ResourceName:            "google_pubsub_lite_subscription.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"topic", "region", "zone", "name"},
			},
			{
				Config: testAccPubsubLiteSubscription_pubsubLiteSubscription_deliveryRequirement_update(topic, subscription, "DELIVER_IMMEDIATELY"),
			},
			{
				ResourceName:            "google_pubsub_lite_subscription.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"topic", "region", "zone", "name"},
			},
		},
	})
}

func testAccPubsubLiteSubscription_pubsubLiteSubscription_deliveryRequirement_update(topic, subscription, delivery string) string {
	return fmt.Sprintf(`
  resource "google_pubsub_lite_topic" "example" {
	name = "%s"
	partition_config {
	  count = 1
	  capacity {
		publish_mib_per_sec = 4
		subscribe_mib_per_sec = 8
	  }
	}
  
	retention_config {
	  per_partition_bytes = 32212254720
	}
  }
  
  resource "google_pubsub_lite_subscription" "example" {
	name  = "%s"
	topic = google_pubsub_lite_topic.example.name
	delivery_config {
	  delivery_requirement = "%s"
	}
  }
`, topic, subscription, delivery)
}
