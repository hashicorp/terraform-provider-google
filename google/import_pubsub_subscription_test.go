package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccPubsubSubscription_import(t *testing.T) {
	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(10))
	subscription := fmt.Sprintf("tf-test-sub-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPubsubTopicDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccPubsubSubscription_basic(topic, subscription),
			},
			resource.TestStep{
				ResourceName:      "google_pubsub_subscription.foobar_sub",
				ImportStateId:     subscription,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
