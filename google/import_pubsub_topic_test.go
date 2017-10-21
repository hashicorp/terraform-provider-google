package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccPubsubTopic_import(t *testing.T) {
	t.Parallel()

	topicName := fmt.Sprintf("tf-test-topic-%d", acctest.RandInt())
	conf := fmt.Sprintf(`
		resource "google_pubsub_topic" "tf-test" {
			name = "%s"
		}`, topicName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPubsubTopicDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: conf,
			},
			resource.TestStep{
				ResourceName:            "google_pubsub_topic.tf-test",
				ImportStateId:           topicName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}
