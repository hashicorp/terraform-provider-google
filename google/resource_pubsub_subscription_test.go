package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccPubsubSubscription_basic(t *testing.T) {
	t.Parallel()

	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(10))
	subscription := fmt.Sprintf("tf-test-sub-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPubsubSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubSubscription_basic(topic, subscription),
			},
			resource.TestStep{
				ResourceName:      "google_pubsub_subscription.foo",
				ImportStateId:     subscription,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TODO: Add acceptance test for push delivery.
//
// Testing push endpoints is tricky for the following reason:
// - You need a publicly accessible HTTPS server to handle POST requests in order to receive push messages.
// - The server must present a valid SSL certificate signed by a certificate authority
// - The server must be routable by DNS.
// - You also need to validate that you own the domain (or have equivalent access to the endpoint).
// - Finally, you must register the endpoint domain with the GCP project.
//
// An easy way to test this would be to create an App Engine Hello World app. With AppEngine, SSL certificate, DNS and domain registry is handled for us.
// App Engine is not yet supported by Terraform but once it is, it will provide an easy path to testing push configs.
// Another option would be to use Cloud Functions once Terraform support is added.
func testAccCheckPubsubSubscriptionDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_pubsub_subscription" {
			continue
		}

		config := testAccProvider.Meta().(*Config)
		sub, _ := config.clientPubsub.Projects.Subscriptions.Get(rs.Primary.ID).Do()
		if sub != nil {
			return fmt.Errorf("Subscription still present")
		}
	}

	return nil
}

func testAccPubsubSubscription_basic(topic, subscription string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
	name = "%s"
}

resource "google_pubsub_subscription" "foo" {
	name                 = "%s"
	topic                = "${google_pubsub_topic.foo.name}"
	ack_deadline_seconds = 20
}`, topic, subscription)
}

func TestGetComputedTopicName(t *testing.T) {
	type testData struct {
		project  string
		topic    string
		expected string
	}

	var testCases = []testData{
		testData{
			project:  "my-project",
			topic:    "my-topic",
			expected: "projects/my-project/topics/my-topic",
		},
		testData{
			project:  "my-project",
			topic:    "projects/another-project/topics/my-topic",
			expected: "projects/another-project/topics/my-topic",
		},
	}

	for _, testCase := range testCases {
		computedTopicName := getComputedTopicName(testCase.project, testCase.topic)
		if computedTopicName != testCase.expected {
			t.Fatalf("bad computed topic name: %s' => expected %s", computedTopicName, testCase.expected)
		}
	}
}
