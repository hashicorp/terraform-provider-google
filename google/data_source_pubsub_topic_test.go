package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGooglePubsubTopic_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPubsubTopicDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGooglePubsubTopic_basic(context),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceState("data.google_pubsub_topic.foo", "google_pubsub_topic.foo"),
				),
			},
		},
	})
}

func TestAccDataSourceGooglePubsubTopic_optionalProject(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPubsubTopicDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGooglePubsubTopic_optionalProject(context),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceState("data.google_pubsub_topic.foo", "google_pubsub_topic.foo"),
				),
			},
		},
	})
}

func testAccDataSourceGooglePubsubTopic_basic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_pubsub_topic" "foo" {
  name     = "tf-test-pubsub-%{random_suffix}"
}

data "google_pubsub_topic" "foo" {
  name     = google_pubsub_topic.foo.name
  project  = google_pubsub_topic.foo.project
}
`, context)
}

func testAccDataSourceGooglePubsubTopic_optionalProject(context map[string]interface{}) string {
	return Nprintf(`
resource "google_pubsub_topic" "foo" {
  name     = "tf-test-pubsub-%{random_suffix}"
}

data "google_pubsub_topic" "foo" {
  name     = google_pubsub_topic.foo.name
}
`, context)
}
