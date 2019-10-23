package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccPubsubSubscription_emptyTTL(t *testing.T) {
	t.Parallel()

	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(10))
	subscription := fmt.Sprintf("projects/%s/subscriptions/tf-test-sub-%s", getTestProjectFromEnv(), acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPubsubSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubSubscription_emptyTTL(topic, subscription),
			},
			{
				ResourceName:      "google_pubsub_subscription.foo",
				ImportStateId:     subscription,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPubsubSubscription_fullName(t *testing.T) {
	t.Parallel()

	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(10))
	subscription := fmt.Sprintf("projects/%s/subscriptions/tf-test-sub-%s", getTestProjectFromEnv(), acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPubsubSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubSubscription_fullName(topic, subscription, "bar", 20),
			},
			{
				ResourceName:      "google_pubsub_subscription.foo",
				ImportStateId:     subscription,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPubsubSubscription_update(t *testing.T) {
	t.Parallel()

	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(10))
	subscriptionShort := fmt.Sprintf("tf-test-sub-%s", acctest.RandString(10))
	subscriptionLong := fmt.Sprintf("projects/%s/subscriptions/%s", getTestProjectFromEnv(), subscriptionShort)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPubsubSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubSubscription_fullName(topic, subscriptionLong, "bar", 20),
			},
			{
				ResourceName:      "google_pubsub_subscription.foo",
				ImportStateId:     subscriptionLong,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPubsubSubscription_fullName(topic, subscriptionLong, "baz", 30),
				Check: resource.TestCheckResourceAttr(
					"google_pubsub_subscription.foo", "path", subscriptionLong,
				),
			},
			{
				ResourceName:      "google_pubsub_subscription.foo",
				ImportStateId:     subscriptionLong,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPubsubSubscription_fullName(topic, subscriptionShort, "baz", 30),
				Check: resource.TestCheckResourceAttr(
					"google_pubsub_subscription.foo", "path", subscriptionLong,
				),
			},
			{
				ResourceName:      "google_pubsub_subscription.foo",
				ImportStateId:     subscriptionShort,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPubsubSubscription_push(t *testing.T) {
	t.Parallel()

	topicFoo := fmt.Sprintf("tf-test-topic-foo-%s", acctest.RandString(10))
	subscription := fmt.Sprintf("projects/%s/subscriptions/tf-test-topic-foo-%s", getTestProjectFromEnv(), acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPubsubSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubSubscription_push(topicFoo, subscription),
			},
			{
				ResourceName:      "google_pubsub_subscription.foo",
				ImportStateId:     subscription,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPubsubSubscription_emptyTTL(topic, subscription string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
	name = "%s"
}

resource "google_pubsub_subscription" "foo" {
	name                 = "%s"
	topic                = "${google_pubsub_topic.foo.id}"

	message_retention_duration = "1200s"
	retain_acked_messages = true
	ack_deadline_seconds = 20
	expiration_policy {}
}
`, topic, subscription)
}

func testAccPubsubSubscription_push(topicFoo string, subscription string) string {
	return fmt.Sprintf(`
data "google_project" "project" {}

resource "google_service_account" "pub_sub_service_account" {
  account_id = "my-super-service"
}

data "google_iam_policy" "admin" {
  binding {
    role = "roles/projects.topics.publish"

    members = [
      "serviceAccount:${google_service_account.pub_sub_service_account.email}",
    ]
  }
}

resource "google_pubsub_topic" "foo" {
	name = "%s"
}

resource "google_pubsub_subscription" "foo" {
  name                 = "%s"
  topic                = "${google_pubsub_topic.foo.name}"
  ack_deadline_seconds = 10
  push_config {
    push_endpoint = "https://${data.google_project.project.project_id}.appspot.com"
    oidc_token {
      service_account_email = "${google_service_account.pub_sub_service_account.email}"
    }
  }
}
`, topicFoo, subscription)
}

func testAccPubsubSubscription_fullName(topic, subscription, label string, deadline int) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
	name = "%s"
}

resource "google_pubsub_subscription" "foo" {
	name                 = "%s"
	topic                = "${google_pubsub_topic.foo.id}"
	labels = {
		foo = "%s"
	}
	ack_deadline_seconds = %d
}
`, topic, subscription, label, deadline)
}

func TestGetComputedTopicName(t *testing.T) {
	type testData struct {
		project  string
		topic    string
		expected string
	}

	var testCases = []testData{
		{
			project:  "my-project",
			topic:    "my-topic",
			expected: "projects/my-project/topics/my-topic",
		},
		{
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
