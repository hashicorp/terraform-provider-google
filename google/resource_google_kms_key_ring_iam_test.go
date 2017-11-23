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

func TestAccGoogleKmsKeyRingIamBinding(t *testing.T) {
	t.Parallel()

	projectId := getTestProjectFromEnv()
	account := acctest.RandomWithPrefix("tf-test")
	roleId := "roles/cloudkms.cryptoKeyDecrypter"
	keyRingName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test Iam Binding creation
				Config: testAccGoogleKmsKeyRingIamBinding_basic(projectId, account, keyRingName, roleId),
				Check: testAccCheckGoogleKmsKeyRingIamBindingExists("foo", roleId, []string{
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, projectId),
				}),
			},
			{
				// Test Iam Binding update
				Config: testAccGoogleKmsKeyRingIamBinding_update(projectId, account, keyRingName, roleId),
				Check: testAccCheckGoogleKmsKeyRingIamBindingExists("foo", roleId, []string{
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, projectId),
					fmt.Sprintf("serviceAccount:%s-2@%s.iam.gserviceaccount.com", account, projectId),
				}),
			},
		},
	})
}

func TestAccGoogleKmsKeyRingIamMember(t *testing.T) {
	t.Parallel()

	projectId := getTestProjectFromEnv()
	account := acctest.RandomWithPrefix("tf-test")
	roleId := "roles/cloudkms.cryptoKeyEncrypter"
	keyRingName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccGoogleKmsKeyRingIamMember_basic(projectId, account, keyRingName, roleId),
				Check: testAccCheckGoogleKmsKeyRingIamMemberExists("foo", roleId,
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, projectId),
				),
			},
		},
	})
}

func testAccCheckGoogleKmsKeyRingIamBindingExists(bindingResourceName, roleId string, members []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		bindingRs, ok := s.RootModule().Resources[fmt.Sprintf("google_kms_key_ring_iam_binding.%s", bindingResourceName)]
		if !ok {
			return fmt.Errorf("Not found: %s", bindingResourceName)
		}

		config := testAccProvider.Meta().(*Config)
		keyRingId, err := parseKmsKeyRingId(bindingRs.Primary.Attributes["key_ring_id"], config)

		if err != nil {
			return err
		}

		p, err := config.clientKms.Projects.Locations.KeyRings.GetIamPolicy(keyRingId.keyRingId()).Do()
		if err != nil {
			return err
		}

		for _, binding := range p.Bindings {
			if binding.Role == roleId {
				sort.Strings(members)
				sort.Strings(binding.Members)

				if reflect.DeepEqual(members, binding.Members) {
					return nil
				}

				return fmt.Errorf("Binding found but expected members is %v, got %v", members, binding.Members)
			}
		}

		return fmt.Errorf("No binding for role %q", roleId)
	}
}

func testAccCheckGoogleKmsKeyRingIamMemberExists(n, role, member string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["google_kms_key_ring_iam_member."+n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		config := testAccProvider.Meta().(*Config)
		keyRingId, err := parseKmsKeyRingId(rs.Primary.Attributes["key_ring_id"], config)

		if err != nil {
			return err
		}

		p, err := config.clientKms.Projects.Locations.KeyRings.GetIamPolicy(keyRingId.keyRingId()).Do()
		if err != nil {
			return err
		}

		for _, binding := range p.Bindings {
			if binding.Role == role {
				for _, m := range binding.Members {
					if m == member {
						return nil
					}
				}

				return fmt.Errorf("Missing member %q, got %v", member, binding.Members)
			}
		}

		return fmt.Errorf("No binding for role %q", role)
	}
}

// We are using a custom role since iam_binding is authoritative on the member list and
// we want to avoid removing members from an existing role to prevent unwanted side effects.
func testAccGoogleKmsKeyRingIamBinding_basic(projectId, account, keyRingName, roleId string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  project      = "%s"
  account_id   = "%s"
  display_name = "Iam Testing Account"
}

resource "google_kms_key_ring" "key_ring" {
  project  = "%s"
  location = "us-central1"
  name     = "%s"
}

resource "google_kms_key_ring_iam_binding" "foo" {
  key_ring_id = "${google_kms_key_ring.key_ring.id}"
  role        = "%s"
  members     = ["serviceAccount:${google_service_account.test_account.email}"]
}
`, projectId, account, projectId, keyRingName, roleId)
}

func testAccGoogleKmsKeyRingIamBinding_update(projectId, account, keyRingName, roleId string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  project      = "%s"
  account_id   = "%s"
  display_name = "Iam Testing Account"
}

resource "google_service_account" "test_account_2" {
  project      = "%s"
  account_id   = "%s-2"
  display_name = "Iam Testing Account"
}

resource "google_kms_key_ring" "key_ring" {
  project  = "%s"
  location = "us-central1"
  name     = "%s"
}

resource "google_kms_key_ring_iam_binding" "foo" {
  key_ring_id  = "${google_kms_key_ring.key_ring.id}"
  role         = "%s"
  members      = [
    "serviceAccount:${google_service_account.test_account.email}",
    "serviceAccount:${google_service_account.test_account_2.email}"
  ]
}
`, projectId, account, projectId, account, projectId, keyRingName, roleId)
}

func testAccGoogleKmsKeyRingIamMember_basic(projectId, account, keyRingName, roleId string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  project      = "%s"
  account_id   = "%s"
  display_name = "Iam Testing Account"
}

resource "google_kms_key_ring" "key_ring" {
  project  = "%s"
  location = "us-central1"
  name     = "%s"
}

resource "google_kms_key_ring_iam_member" "foo" {
  key_ring_id = "${google_kms_key_ring.key_ring.id}"
  role        = "%s"
  member      = "serviceAccount:${google_service_account.test_account.email}"
}
`, projectId, account, projectId, keyRingName, roleId)
}
