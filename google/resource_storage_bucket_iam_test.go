package google

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccStorageBucketIamBinding(t *testing.T) {
	t.Parallel()

	bucket := acctest.RandomWithPrefix("tf-test")
	account := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test IAM Binding creation
				Config: testAccStorageBucketIamBinding_basic(bucket, account),
				Check: testAccCheckGoogleStorageBucketIam(bucket, "roles/storage.objectViewer", []string{
					fmt.Sprintf("serviceAccount:%s-1@%s.iam.gserviceaccount.com", account, getTestProjectFromEnv()),
				}),
			},
			{
				// Test IAM Binding update
				Config: testAccStorageBucketIamBinding_update(bucket, account),
				Check: testAccCheckGoogleStorageBucketIam(bucket, "roles/storage.objectViewer", []string{
					fmt.Sprintf("serviceAccount:%s-1@%s.iam.gserviceaccount.com", account, getTestProjectFromEnv()),
					fmt.Sprintf("serviceAccount:%s-2@%s.iam.gserviceaccount.com", account, getTestProjectFromEnv()),
				}),
			},
		},
	})
}

func TestAccStorageBucketIamPolicy(t *testing.T) {
	t.Parallel()

	bucket := acctest.RandomWithPrefix("tf-test")
	account := acctest.RandomWithPrefix("tf-test")
	serviceAcct := getTestServiceAccountFromEnv(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test IAM Policy creation
				Config: testAccStorageBucketIamPolicy_basic(bucket, account, serviceAcct),
				Check: testAccCheckGoogleStorageBucketIam(bucket, "roles/storage.objectViewer", []string{
					fmt.Sprintf("serviceAccount:%s-1@%s.iam.gserviceaccount.com", account, getTestProjectFromEnv()),
				}),
			},
			{
				// Test IAM Policy update
				Config: testAccStorageBucketIamPolicy_update(bucket, account, serviceAcct),
				Check: testAccCheckGoogleStorageBucketIam(bucket, "roles/storage.objectViewer", []string{
					fmt.Sprintf("serviceAccount:%s-1@%s.iam.gserviceaccount.com", account, getTestProjectFromEnv()),
					fmt.Sprintf("serviceAccount:%s-2@%s.iam.gserviceaccount.com", account, getTestProjectFromEnv()),
				}),
			},
		},
	})
}

func TestAccStorageBucketIamMember(t *testing.T) {
	t.Parallel()

	bucket := acctest.RandomWithPrefix("tf-test")
	account := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccStorageBucketIamMember_basic(bucket, account),
				Check: testAccCheckGoogleStorageBucketIam(bucket, "roles/storage.admin", []string{
					fmt.Sprintf("serviceAccount:%s-1@%s.iam.gserviceaccount.com", account, getTestProjectFromEnv()),
				}),
			},
		},
	})
}

func testAccCheckGoogleStorageBucketIam(bucket, role string, members []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		p, err := config.clientStorage.Buckets.GetIamPolicy(bucket).Do()
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

func testAccStorageBucketIamPolicy_update(bucket, account, serviceAcct string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
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


data "google_iam_policy" "foo-policy" {
	binding {
		role = "roles/storage.objectViewer"

		members = [
			"serviceAccount:${google_service_account.test-account-1.email}",
			"serviceAccount:${google_service_account.test-account-2.email}",
		]
	}

	binding {
		role = "roles/storage.admin"
		members = [
			"serviceAccount:%s",
		]
	}
}

resource "google_storage_bucket_iam_policy" "bucket-binding" {
	bucket      = "${google_storage_bucket.bucket.name}"
	policy_data = "${data.google_iam_policy.foo-policy.policy_data}"
}

`, bucket, account, account, serviceAcct)
}

func testAccStorageBucketIamPolicy_basic(bucket, account, serviceAcct string) string {
	return fmt.Sprintf(`

resource "google_storage_bucket" "bucket" {
	name = "%s"
}

resource "google_service_account" "test-account-1" {
	account_id   = "%s-1"
	display_name = "Iam Testing Account"
}


data "google_iam_policy" "foo-policy" {
	binding {
		role = "roles/storage.objectViewer"
		members = [
			"serviceAccount:${google_service_account.test-account-1.email}",
		]
	}

	binding {
		role = "roles/storage.admin"
		members = [
			"serviceAccount:%s",
		]
	}
}

resource "google_storage_bucket_iam_policy" "bucket-binding" {
	bucket      = "${google_storage_bucket.bucket.name}"
	policy_data = "${data.google_iam_policy.foo-policy.policy_data}"
}


`, bucket, account, serviceAcct)
}

func testAccStorageBucketIamBinding_basic(bucket, account string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
	name = "%s"
}

resource "google_service_account" "test-account-1" {
	account_id   = "%s-1"
	display_name = "Iam Testing Account"
}

resource "google_storage_bucket_iam_binding" "foo" {
	bucket = "${google_storage_bucket.bucket.name}"
	role = "roles/storage.objectViewer"
	members = [
		"serviceAccount:${google_service_account.test-account-1.email}",
	]
}
`, bucket, account)
}

func testAccStorageBucketIamBinding_update(bucket, account string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
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

resource "google_storage_bucket_iam_binding" "foo" {
	bucket = "${google_storage_bucket.bucket.name}"
	role = "roles/storage.objectViewer"
	members = [
		"serviceAccount:${google_service_account.test-account-1.email}",
		"serviceAccount:${google_service_account.test-account-2.email}",
	]
}
`, bucket, account, account)
}

func testAccStorageBucketIamMember_basic(bucket, account string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
	name = "%s"
}

resource "google_service_account" "test-account-1" {
	account_id   = "%s-1"
	display_name = "Iam Testing Account"
}

resource "google_storage_bucket_iam_member" "foo" {
	bucket = "${google_storage_bucket.bucket.name}"
	role = "roles/storage.admin"
	member = "serviceAccount:${google_service_account.test-account-1.email}"
}
`, bucket, account)
}
