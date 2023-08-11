// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package pubsublite_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccPubsubLiteTopic_pubsubLiteTopic_count_update(t *testing.T) {
	t.Parallel()

	topic := fmt.Sprintf("tf-test-topic-foo-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubLiteTopicDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubLiteTopic_pubsubLiteTopic_count_update(topic, "1"),
			},
			{
				ResourceName:            "google_pubsub_lite_topic.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region", "zone", "name"},
			},
			{
				Config: testAccPubsubLiteTopic_pubsubLiteTopic_count_update(topic, "2"),
			},
			{
				ResourceName:            "google_pubsub_lite_topic.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region", "zone", "name"},
			},
		},
	})
}

func testAccPubsubLiteTopic_pubsubLiteTopic_count_update(topic, count string) string {
	return fmt.Sprintf(`
resource "google_pubsub_lite_topic" "example" {
  name = "%s"

  partition_config {
    count = %s
    capacity {
      publish_mib_per_sec = 4
      subscribe_mib_per_sec = 8
    }
  }

  retention_config {
    per_partition_bytes = 32212254720
  }
}
`, topic, count)
}
