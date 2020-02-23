package google

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPubsubSubscriptionIamBinding(t *testing.T) {
	t.Parallel()

	topic := "tf-test-topic-iam-" + acctest.RandString(10)
	subscription := "tf-test-sub-iam-" + acctest.RandString(10)
	account := "tf-test-iam-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test IAM Binding creation
				Config: testAccPubsubSubscriptionIamBinding_basic(subscription, topic, account),
				Check: testAccCheckPubsubSubscriptionIam(subscription, "roles/pubsub.subscriber", []string{
					fmt.Sprintf("serviceAccount:%s-1@%s.iam.gserviceaccount.com", account, getTestProjectFromEnv()),
				}),
			},
			{
				// Test IAM Binding update
				Config: testAccPubsubSubscriptionIamBinding_update(subscription, topic, account),
				Check: testAccCheckPubsubSubscriptionIam(subscription, "roles/pubsub.subscriber", []string{
					fmt.Sprintf("serviceAccount:%s-1@%s.iam.gserviceaccount.com", account, getTestProjectFromEnv()),
					fmt.Sprintf("serviceAccount:%s-2@%s.iam.gserviceaccount.com", account, getTestProjectFromEnv()),
				}),
			},
			{
				ResourceName:      "google_pubsub_subscription_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("%s roles/pubsub.subscriber", getComputedSubscriptionName(getTestProjectFromEnv(), subscription)),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPubsubSubscriptionIamMember(t *testing.T) {
	t.Parallel()

	topic := "tf-test-topic-iam-" + acctest.RandString(10)
	subscription := "tf-test-sub-iam-" + acctest.RandString(10)
	account := "tf-test-iam-" + acctest.RandString(10)
	accountEmail := fmt.Sprintf("%s@%s.iam.gserviceaccount.com", account, getTestProjectFromEnv())

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccPubsubSubscriptionIamMember_basic(subscription, topic, account),
				Check: testAccCheckPubsubSubscriptionIam(subscription, "roles/pubsub.subscriber", []string{
					fmt.Sprintf("serviceAccount:%s", accountEmail),
				}),
			},
			{
				ResourceName:      "google_pubsub_subscription_iam_member.foo",
				ImportStateId:     fmt.Sprintf("%s roles/pubsub.subscriber serviceAccount:%s", getComputedSubscriptionName(getTestProjectFromEnv(), subscription), accountEmail),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPubsubSubscriptionIamPolicy(t *testing.T) {
	t.Parallel()

	topic := "tf-test-topic-iam-" + acctest.RandString(10)
	subscription := "tf-test-sub-iam-" + acctest.RandString(10)
	account := "tf-test-iam-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubSubscriptionIamPolicy_basic(subscription, topic, account, "roles/pubsub.subscriber"),
				Check: testAccCheckPubsubSubscriptionIam(subscription, "roles/pubsub.subscriber", []string{
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, getTestProjectFromEnv()),
				}),
			},
			{
				Config: testAccPubsubSubscriptionIamPolicy_basic(subscription, topic, account, "roles/pubsub.viewer"),
				Check: testAccCheckPubsubSubscriptionIam(subscription, "roles/pubsub.viewer", []string{
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, getTestProjectFromEnv()),
				}),
			},
			{
				ResourceName:      "google_pubsub_subscription_iam_policy.foo",
				ImportStateId:     getComputedSubscriptionName(getTestProjectFromEnv(), subscription),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPubsubSubscriptionIam(subscription, role string, members []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		p, err := config.clientPubsub.Projects.Subscriptions.GetIamPolicy(getComputedSubscriptionName(getTestProjectFromEnv(), subscription)).Do()
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

func testAccPubsubSubscriptionIamBinding_basic(subscription, topic, account string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "topic" {
  name = "%s"
}

resource "google_pubsub_subscription" "subscription" {
  name  = "%s"
  topic = google_pubsub_topic.topic.id
}

resource "google_service_account" "test-account-1" {
  account_id   = "%s-1"
  display_name = "Pubsub Subscription Iam Testing Account"
}

resource "google_pubsub_subscription_iam_binding" "foo" {
  subscription = google_pubsub_subscription.subscription.id
  role         = "roles/pubsub.subscriber"
  members = [
    "serviceAccount:${google_service_account.test-account-1.email}",
  ]
}
`, topic, subscription, account)
}

func testAccPubsubSubscriptionIamBinding_update(subscription, topic, account string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "topic" {
  name = "%s"
}

resource "google_pubsub_subscription" "subscription" {
  name  = "%s"
  topic = google_pubsub_topic.topic.id
}

resource "google_service_account" "test-account-1" {
  account_id   = "%s-1"
  display_name = "Pubsub Subscription Iam Testing Account"
}

resource "google_service_account" "test-account-2" {
  account_id   = "%s-2"
  display_name = "Pubsub Subscription Iam Testing Account"
}

resource "google_pubsub_subscription_iam_binding" "foo" {
  subscription = google_pubsub_subscription.subscription.id
  role         = "roles/pubsub.subscriber"
  members = [
    "serviceAccount:${google_service_account.test-account-1.email}",
    "serviceAccount:${google_service_account.test-account-2.email}",
  ]
}
`, topic, subscription, account, account)
}

func testAccPubsubSubscriptionIamMember_basic(subscription, topic, account string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "topic" {
  name = "%s"
}

resource "google_pubsub_subscription" "subscription" {
  name  = "%s"
  topic = google_pubsub_topic.topic.id
}

resource "google_service_account" "test-account" {
  account_id   = "%s"
  display_name = "Pubsub Subscription Iam Testing Account"
}

resource "google_pubsub_subscription_iam_member" "foo" {
  subscription = google_pubsub_subscription.subscription.id
  role         = "roles/pubsub.subscriber"
  member       = "serviceAccount:${google_service_account.test-account.email}"
}
`, topic, subscription, account)
}

func testAccPubsubSubscriptionIamPolicy_basic(subscription, topic, account, role string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "topic" {
  name = "%s"
}

resource "google_pubsub_subscription" "subscription" {
  name  = "%s"
  topic = google_pubsub_topic.topic.id
}

resource "google_service_account" "test-account" {
  account_id   = "%s"
  display_name = "Pubsub Subscription Iam Testing Account"
}

data "google_iam_policy" "foo" {
  binding {
    role    = "%s"
    members = ["serviceAccount:${google_service_account.test-account.email}"]
  }
}

resource "google_pubsub_subscription_iam_policy" "foo" {
  subscription = google_pubsub_subscription.subscription.id
  policy_data  = data.google_iam_policy.foo.policy_data
}
`, topic, subscription, account, role)
}
