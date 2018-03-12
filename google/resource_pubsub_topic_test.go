package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccPubsubTopic_basic(t *testing.T) {
	t.Parallel()

	topicName := acctest.RandomWithPrefix("tf-test-topic")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPubsubTopicDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccPubsubTopic_basic(topicName),
			},
			// Check importing with just the topic name
			resource.TestStep{
				ResourceName:            "google_pubsub_topic.foo",
				ImportStateId:           topicName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			// Check importing with the full resource id
			resource.TestStep{
				ResourceName:            "google_pubsub_topic.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

func testAccCheckPubsubTopicDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_pubsub_topic" {
			continue
		}

		config := testAccProvider.Meta().(*Config)
		topic, _ := config.clientPubsub.Projects.Topics.Get(rs.Primary.ID).Do()
		if topic != nil {
			return fmt.Errorf("Topic still present")
		}
	}

	return nil
}

func testAccPubsubTopic_basic(name string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
	name = "%s"
}`, name)
}
