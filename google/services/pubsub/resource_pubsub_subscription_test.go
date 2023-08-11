// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package pubsub_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/pubsub"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccPubsubSubscription_emptyTTL(t *testing.T) {
	t.Parallel()

	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))
	subscription := fmt.Sprintf("tf-test-sub-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubSubscriptionDestroyProducer(t),
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

func TestAccPubsubSubscription_basic(t *testing.T) {
	t.Parallel()

	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))
	subscription := fmt.Sprintf("tf-test-sub-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubSubscriptionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubSubscription_basic(topic, subscription, "bar", 20, false),
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

	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))
	subscriptionShort := fmt.Sprintf("tf-test-sub-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubSubscriptionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubSubscription_basic(topic, subscriptionShort, "bar", 20, false),
			},
			{
				ResourceName:      "google_pubsub_subscription.foo",
				ImportStateId:     subscriptionShort,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPubsubSubscription_basic(topic, subscriptionShort, "baz", 30, true),
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

	topicFoo := fmt.Sprintf("tf-test-topic-foo-%s", acctest.RandString(t, 10))
	subscription := fmt.Sprintf("tf-test-sub-foo-%s", acctest.RandString(t, 10))
	saAccount := fmt.Sprintf("tf-test-pubsub-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubSubscriptionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubSubscription_push(topicFoo, saAccount, subscription),
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

func TestAccPubsubSubscription_pushNoWrapper(t *testing.T) {
	t.Parallel()

	topicFoo := fmt.Sprintf("tf-test-topic-foo-%s", acctest.RandString(t, 10))
	subscription := fmt.Sprintf("tf-test-sub-foo-%s", acctest.RandString(t, 10))
	saAccount := fmt.Sprintf("tf-test-pubsub-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubSubscriptionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubSubscription_pushNoWrapper(topicFoo, saAccount, subscription),
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

// Context: hashicorp/terraform-provider-google#4993
// This test makes a call to GET an subscription before it is actually created.
// The PubSub API negative-caches responses so this tests we are
// correctly polling for existence post-creation.
func TestAccPubsubSubscription_pollOnCreate(t *testing.T) {
	t.Parallel()

	topic := fmt.Sprintf("tf-test-topic-foo-%s", acctest.RandString(t, 10))
	subscription := fmt.Sprintf("tf-test-topic-foo-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubSubscriptionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// Create only the topic
				Config: testAccPubsubSubscription_topicOnly(topic),
				// Read from non-existent subscription created in next step
				// so API negative-caches result
				Check: testAccCheckPubsubSubscriptionCache404(t, subscription),
			},
			{
				// Create the subscription - if the polling fails,
				// the test step will fail because the read post-create
				// will have removed the resource from state.
				Config: testAccPubsubSubscription_pollOnCreate(topic, subscription),
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
  name  = "%s"
  topic = google_pubsub_topic.foo.id

  message_retention_duration = "1200s"
  retain_acked_messages      = true
  ack_deadline_seconds       = 20
  expiration_policy {
    ttl = ""
  }
  enable_message_ordering    = false
}
`, topic, subscription)
}

func testAccPubsubSubscription_push(topicFoo, saAccount, subscription string) string {
	return fmt.Sprintf(`
data "google_project" "project" { }

resource "google_service_account" "pub_sub_service_account" {
  account_id = "%s"
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
  topic                = google_pubsub_topic.foo.name
  ack_deadline_seconds = 10
  push_config {
    push_endpoint = "https://${data.google_project.project.project_id}.appspot.com"
    oidc_token {
      service_account_email = google_service_account.pub_sub_service_account.email
    }
  }
}
`, saAccount, topicFoo, subscription)
}

func testAccPubsubSubscription_pushNoWrapper(topicFoo, saAccount, subscription string) string {
	return fmt.Sprintf(`
data "google_project" "project" { }

resource "google_service_account" "pub_sub_service_account" {
  account_id = "%s"
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
  topic                = google_pubsub_topic.foo.name
  ack_deadline_seconds = 10
  push_config {
    push_endpoint = "https://${data.google_project.project.project_id}.appspot.com"
    oidc_token {
      service_account_email = google_service_account.pub_sub_service_account.email
    }
    no_wrapper {
      write_metadata = true
    }
  }
}
`, saAccount, topicFoo, subscription)
}

func testAccPubsubSubscription_basic(topic, subscription, label string, deadline int, exactlyOnceDelivery bool) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
  name = "%s"
}

resource "google_pubsub_subscription" "foo" {
  name   = "%s"
  topic  = google_pubsub_topic.foo.id
  filter = "attributes.foo = \"bar\""
  labels = {
    foo = "%s"
  }
  retry_policy {
    minimum_backoff = "60.0s"
  }
  ack_deadline_seconds = %d
  enable_exactly_once_delivery = %t
}
`, topic, subscription, label, deadline, exactlyOnceDelivery)
}

func testAccPubsubSubscription_topicOnly(topic string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
  name = "%s"
}
`, topic)
}

func testAccPubsubSubscription_pollOnCreate(topic, subscription string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
  name = "%s"
}

resource "google_pubsub_subscription" "foo" {
  name  = "%s"
  topic = google_pubsub_topic.foo.id
}
`, topic, subscription)
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
		computedTopicName := pubsub.GetComputedTopicName(testCase.project, testCase.topic)
		if computedTopicName != testCase.expected {
			t.Fatalf("bad computed topic name: %s' => expected %s", computedTopicName, testCase.expected)
		}
	}
}

func testAccCheckPubsubSubscriptionCache404(t *testing.T, subName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)
		url := fmt.Sprintf("%sprojects/%s/subscriptions/%s", config.PubsubBasePath, envvar.GetTestProjectFromEnv(), subName)
		resp, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			RawURL:    url,
			UserAgent: config.UserAgent,
		})
		if err == nil {
			return fmt.Errorf("Expected Pubsub Subscription %q not to exist, was found", resp["name"])
		}
		if !transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
			return fmt.Errorf("Got non-404 error while trying to read Pubsub Subscription %q: %v", subName, err)
		}
		return nil
	}
}
