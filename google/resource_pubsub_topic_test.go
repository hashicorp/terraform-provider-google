package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccPubsubTopic_fullName(t *testing.T) {
	t.Parallel()

	topicName := fmt.Sprintf("projects/%s/topics/tf-test-topic-%s", getTestProjectFromEnv(), acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPubsubTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubTopic_fullName(topicName),
			},
			// Check importing with just the topic name
			{
				ResourceName:            "google_pubsub_topic.foo",
				ImportStateId:           topicName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			// Check importing with the full resource id
			{
				ResourceName:            "google_pubsub_topic.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

func testAccPubsubTopic_fullName(name string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
	name = "%s"
}`, name)
}
