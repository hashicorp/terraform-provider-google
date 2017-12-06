package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"reflect"
	"sort"
	"testing"
)

const DEFAULT_STORAGE_BUCKET_TEST_ROLE = "roles/storage.objectViewer"

func TestAccGoogleStorageBucketIamBinding(t *testing.T) {
	t.Parallel()

	bucket := acctest.RandomWithPrefix("tf-test")
	account := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test IAM Binding creation
				Config: testAccGoogleStorageBucketIamBinding_basic(bucket, account),
				Check: testAccCheckGoogleStorageBucketIam(bucket, DEFAULT_STORAGE_BUCKET_TEST_ROLE, []string{
					fmt.Sprintf("serviceAccount:%s-1@%s.iam.gserviceaccount.com", account, getTestProjectFromEnv()),
				}),
			},
			{
				// Test IAM Binding update
				Config: testAccGoogleStorageBucketIamBinding_update(bucket, account),
				Check: testAccCheckGoogleStorageBucketIam(bucket, DEFAULT_STORAGE_BUCKET_TEST_ROLE, []string{
					fmt.Sprintf("serviceAccount:%s-1@%s.iam.gserviceaccount.com", account, getTestProjectFromEnv()),
					fmt.Sprintf("serviceAccount:%s-2@%s.iam.gserviceaccount.com", account, getTestProjectFromEnv()),
				}),
			},
		},
	})
}

func TestAccGoogleStorageBucketIamMember(t *testing.T) {
	t.Parallel()

	bucket := acctest.RandomWithPrefix("tf-test")
	account := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccGoogleStorageBucketIamMember_basic(bucket, account),
				Check: testAccCheckGoogleStorageBucketIam(bucket, DEFAULT_STORAGE_BUCKET_TEST_ROLE, []string{
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

func testAccGoogleStorageBucketIamBinding_basic(bucket, account string) string {
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
  role = "%s"
  members = [
    "serviceAccount:${google_service_account.test-account-1.email}",
  ]
}
`, bucket, account, DEFAULT_STORAGE_BUCKET_TEST_ROLE)
}

func testAccGoogleStorageBucketIamBinding_update(bucket, account string) string {
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
  role = "%s"
  members = [
    "serviceAccount:${google_service_account.test-account-1.email}",
    "serviceAccount:${google_service_account.test-account-2.email}",
  ]
}
`, bucket, account, account, DEFAULT_STORAGE_BUCKET_TEST_ROLE)
}

func testAccGoogleStorageBucketIamMember_basic(bucket, account string) string {
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
  role = "%s"
  member = "serviceAccount:${google_service_account.test-account-1.email}"
}
`, bucket, account, DEFAULT_STORAGE_BUCKET_TEST_ROLE)
}
