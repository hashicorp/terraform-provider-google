package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccPubsubTopic_update(t *testing.T) {
	t.Parallel()

	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
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
				Config: testAccPubsubTopic_update(topic, "wibble", "wobble"),
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
