package google

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPubsubTopicIamBinding(t *testing.T) {
	t.Parallel()

	topic := "tf-test-topic-iam-" + randString(t, 10)
	account := "tf-test-topic-iam-" + randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test IAM Binding creation
				Config: testAccPubsubTopicIamBinding_basic(topic, account),
				Check: testAccCheckPubsubTopicIam(t, topic, "roles/pubsub.publisher", []string{
					fmt.Sprintf("serviceAccount:%s-1@%s.iam.gserviceaccount.com", account, getTestProjectFromEnv()),
				}),
			},
			{
				ResourceName:      "google_pubsub_topic_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("%s roles/pubsub.publisher", getComputedTopicName(getTestProjectFromEnv(), topic)),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test IAM Binding update
				Config: testAccPubsubTopicIamBinding_update(topic, account),
				Check: testAccCheckPubsubTopicIam(t, topic, "roles/pubsub.publisher", []string{
					fmt.Sprintf("serviceAccount:%s-1@%s.iam.gserviceaccount.com", account, getTestProjectFromEnv()),
					fmt.Sprintf("serviceAccount:%s-2@%s.iam.gserviceaccount.com", account, getTestProjectFromEnv()),
				}),
			},
			{
				ResourceName:      "google_pubsub_topic_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("%s roles/pubsub.publisher", getComputedTopicName(getTestProjectFromEnv(), topic)),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPubsubTopicIamBinding_topicName(t *testing.T) {
	t.Parallel()

	topic := "tf-test-topic-iam-" + randString(t, 10)
	account := "tf-test-topic-iam-" + randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test IAM Binding creation
				Config: testAccPubsubTopicIamBinding_topicName(topic, account),
				Check: testAccCheckPubsubTopicIam(t, topic, "roles/pubsub.publisher", []string{
					fmt.Sprintf("serviceAccount:%s-1@%s.iam.gserviceaccount.com", account, getTestProjectFromEnv()),
				}),
			},
			// No import step- imports want the resource to be defined using the full id as the topic
		},
	})
}

func TestAccPubsubTopicIamMember(t *testing.T) {
	t.Parallel()

	topic := "tf-test-topic-iam-" + randString(t, 10)
	account := "tf-test-topic-iam-" + randString(t, 10)
	accountEmail := fmt.Sprintf("%s@%s.iam.gserviceaccount.com", account, getTestProjectFromEnv())

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccPubsubTopicIamMember_basic(topic, account),
				Check: testAccCheckPubsubTopicIam(t, topic, "roles/pubsub.publisher", []string{
					fmt.Sprintf("serviceAccount:%s", accountEmail),
				}),
			},
			{
				ResourceName:      "google_pubsub_topic_iam_member.foo",
				ImportStateId:     fmt.Sprintf("%s roles/pubsub.publisher serviceAccount:%s", getComputedTopicName(getTestProjectFromEnv(), topic), accountEmail),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPubsubTopicIamPolicy(t *testing.T) {
	t.Parallel()

	topic := "tf-test-topic-iam-" + randString(t, 10)
	account := "tf-test-topic-iam-" + randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubTopicIamPolicy_basic(topic, account, "roles/pubsub.publisher"),
				Check: testAccCheckPubsubTopicIam(t, topic, "roles/pubsub.publisher", []string{
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, getTestProjectFromEnv()),
				}),
			},
			{
				Config: testAccPubsubTopicIamPolicy_basic(topic, account, "roles/pubsub.subscriber"),
				Check: testAccCheckPubsubTopicIam(t, topic, "roles/pubsub.subscriber", []string{
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, getTestProjectFromEnv()),
				}),
			},
			{
				ResourceName:      "google_pubsub_topic_iam_policy.foo",
				ImportStateId:     getComputedTopicName(getTestProjectFromEnv(), topic),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPubsubTopicIam(t *testing.T, topic, role string, members []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)
		p, err := config.clientPubsub.Projects.Topics.GetIamPolicy(getComputedTopicName(getTestProjectFromEnv(), topic)).Do()
		if err != nil {
			return err
		}

		for _, binding := range p.Bindings {
			if binding.Role == role {
				sort.Strings(members)
				sort.Strings(binding.Members)

				if reflect.DeepEqual(members, binding.Members) {
					return nil
				}

				return fmt.Errorf("Binding found but expected members is %v, got %v", members, binding.Members)
			}
		}

		return fmt.Errorf("No binding for role %q", role)
	}
}

func testAccPubsubTopicIamBinding_topicName(topic, account string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "topic" {
  name = "%s"
}

resource "google_service_account" "test-account-1" {
  account_id   = "%s-1"
  display_name = "Iam Testing Account"
}

resource "google_pubsub_topic_iam_binding" "foo" {
  project = "%s"
  topic   = google_pubsub_topic.topic.name
  role    = "roles/pubsub.publisher"
  members = [
    "serviceAccount:${google_service_account.test-account-1.email}",
  ]
}
`, topic, account, getTestProjectFromEnv())
}

func testAccPubsubTopicIamBinding_basic(topic, account string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "topic" {
  name = "%s"
}

resource "google_service_account" "test-account-1" {
  account_id   = "%s-1"
  display_name = "Iam Testing Account"
}

resource "google_pubsub_topic_iam_binding" "foo" {
  # use the id instead of the name because it's more compatible with import
  topic = google_pubsub_topic.topic.id
  role  = "roles/pubsub.publisher"
  members = [
    "serviceAccount:${google_service_account.test-account-1.email}",
  ]
}
`, topic, account)
}

func testAccPubsubTopicIamBinding_update(topic, account string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "topic" {
  name = "%s"
}

resource "google_service_account" "test-account-1" {
  account_id   = "%s-1"
  display_name = "Iam Testing Account"
}

resource "google_service_account" "test-account-2" {
  account_id   = "%s-2"
  display_name = "Iam Testing Account"
}

resource "google_pubsub_topic_iam_binding" "foo" {
  # use the id instead of the name because it's more compatible with import
  topic = google_pubsub_topic.topic.id
  role  = "roles/pubsub.publisher"
  members = [
    "serviceAccount:${google_service_account.test-account-1.email}",
    "serviceAccount:${google_service_account.test-account-2.email}",
  ]
}
`, topic, account, account)
}

func testAccPubsubTopicIamMember_basic(topic, account string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "topic" {
  name = "%s"
}

resource "google_service_account" "test-account" {
  account_id   = "%s"
  display_name = "Iam Testing Account"
}

resource "google_pubsub_topic_iam_member" "foo" {
  topic  = google_pubsub_topic.topic.id
  role   = "roles/pubsub.publisher"
  member = "serviceAccount:${google_service_account.test-account.email}"
}
`, topic, account)
}

func testAccPubsubTopicIamPolicy_basic(topic, account, role string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "topic" {
  name = "%s"
}

resource "google_service_account" "test-account" {
  account_id   = "%s"
  display_name = "Iam Testing Account"
}

data "google_iam_policy" "foo" {
  binding {
    role    = "%s"
    members = ["serviceAccount:${google_service_account.test-account.email}"]
  }
}

resource "google_pubsub_topic_iam_policy" "foo" {
  topic       = google_pubsub_topic.topic.id
  policy_data = data.google_iam_policy.foo.policy_data
}
`, topic, account, role)
}
