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

func TestAccKmsCryptoKeyIamBinding(t *testing.T) {
	t.Parallel()

	orgId := getTestOrgFromEnv(t)
	projectId := acctest.RandomWithPrefix("tf-test")
	billingAccount := getTestBillingAccountFromEnv(t)
	account := acctest.RandomWithPrefix("tf-test")
	roleId := "roles/cloudkms.cryptoKeyDecrypter"
	keyRingName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	keyRingId := &kmsKeyRingId{
		Project:  projectId,
		Location: DEFAULT_KMS_TEST_LOCATION,
		Name:     keyRingName,
	}
	cryptoKeyName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test Iam Binding creation
				Config: testAccKmsCryptoKeyIamBinding_basic(projectId, orgId, billingAccount, account, keyRingName, cryptoKeyName, roleId),
				Check: testAccCheckGoogleKmsCryptoKeyIamBindingExists("foo", roleId, []string{
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, projectId),
				}),
			},
			{
				ResourceName:      "google_kms_crypto_key_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("%s/%s %s", keyRingId.terraformId(), cryptoKeyName, roleId),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test Iam Binding update
				Config: testAccKmsCryptoKeyIamBinding_update(projectId, orgId, billingAccount, account, keyRingName, cryptoKeyName, roleId),
				Check: testAccCheckGoogleKmsCryptoKeyIamBindingExists("foo", roleId, []string{
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, projectId),
					fmt.Sprintf("serviceAccount:%s-2@%s.iam.gserviceaccount.com", account, projectId),
				}),
			},
			{
				ResourceName:      "google_kms_crypto_key_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("%s/%s %s", keyRingId.terraformId(), cryptoKeyName, roleId),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccKmsCryptoKeyIamMember(t *testing.T) {
	t.Parallel()

	orgId := getTestOrgFromEnv(t)
	projectId := acctest.RandomWithPrefix("tf-test")
	billingAccount := getTestBillingAccountFromEnv(t)
	account := acctest.RandomWithPrefix("tf-test")
	roleId := "roles/cloudkms.cryptoKeyEncrypter"
	keyRingName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	keyRingId := &kmsKeyRingId{
		Project:  projectId,
		Location: DEFAULT_KMS_TEST_LOCATION,
		Name:     keyRingName,
	}
	cryptoKeyName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccKmsCryptoKeyIamMember_basic(projectId, orgId, billingAccount, account, keyRingName, cryptoKeyName, roleId),
				Check: testAccCheckGoogleKmsCryptoKeyIamMemberExists("foo", roleId,
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, projectId),
				),
			},
			{
				ResourceName:      "google_kms_crypto_key_iam_member.foo",
				ImportStateId:     fmt.Sprintf("%s/%s %s serviceAccount:%s@%s.iam.gserviceaccount.com", keyRingId.terraformId(), cryptoKeyName, roleId, account, projectId),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGoogleKmsCryptoKeyIamBindingExists(bindingResourceName, roleId string, members []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		bindingRs, ok := s.RootModule().Resources[fmt.Sprintf("google_kms_crypto_key_iam_binding.%s", bindingResourceName)]
		if !ok {
			return fmt.Errorf("Not found: %s", bindingResourceName)
		}

		config := testAccProvider.Meta().(*Config)
		cryptoKeyId, err := parseKmsCryptoKeyId(bindingRs.Primary.Attributes["crypto_key_id"], config)

		if err != nil {
			return err
		}

		p, err := config.clientKms.Projects.Locations.KeyRings.CryptoKeys.GetIamPolicy(cryptoKeyId.cryptoKeyId()).Do()
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

func testAccCheckGoogleKmsCryptoKeyIamMemberExists(n, role, member string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["google_kms_crypto_key_iam_member."+n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		config := testAccProvider.Meta().(*Config)
		cryptoKeyId, err := parseKmsCryptoKeyId(rs.Primary.Attributes["crypto_key_id"], config)

		if err != nil {
			return err
		}

		p, err := config.clientKms.Projects.Locations.KeyRings.GetIamPolicy(cryptoKeyId.cryptoKeyId()).Do()
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
func testAccKmsCryptoKeyIamBinding_basic(projectId, orgId, billingAccount, account, keyRingName, cryptoKeyName, roleId string) string {
	return fmt.Sprintf(`
resource "google_project" "test_project" {
  name            = "Test project"
  project_id      = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_service" "kms" {
  project = google_project.test_project.project_id
  service = "cloudkms.googleapis.com"
}

resource "google_project_service" "iam" {
  project = google_project_service.kms.project
  service = "iam.googleapis.com"
}

resource "google_service_account" "test_account" {
  project      = google_project_service.iam.project
  account_id   = "%s"
  display_name = "Kms Crypto Key Iam Testing Account"
}

resource "google_kms_key_ring" "key_ring" {
  project  = google_project_service.iam.project
  location = "us-central1"
  name     = "%s"
}

resource "google_kms_crypto_key" "crypto_key" {
  key_ring = google_kms_key_ring.key_ring.id
  name     = "%s"
}

resource "google_kms_crypto_key_iam_binding" "foo" {
  crypto_key_id = google_kms_crypto_key.crypto_key.id
  role          = "%s"
  members       = ["serviceAccount:${google_service_account.test_account.email}"]
}
`, projectId, orgId, billingAccount, account, keyRingName, cryptoKeyName, roleId)
}

func testAccKmsCryptoKeyIamBinding_update(projectId, orgId, billingAccount, account, keyRingName, cryptoKeyName, roleId string) string {
	return fmt.Sprintf(`
resource "google_project" "test_project" {
  name            = "Test project"
  project_id      = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_service" "kms" {
  project = google_project.test_project.project_id
  service = "cloudkms.googleapis.com"
}

resource "google_project_service" "iam" {
  project = google_project_service.kms.project
  service = "iam.googleapis.com"
}

resource "google_service_account" "test_account" {
  project      = google_project_service.iam.project
  account_id   = "%s"
  display_name = "Kms Crypto Key Iam Testing Account"
}

resource "google_service_account" "test_account_2" {
  project      = google_project_service.iam.project
  account_id   = "%s-2"
  display_name = "Kms Crypto Key Iam Testing Account"
}

resource "google_kms_key_ring" "key_ring" {
  project  = google_project_service.iam.project
  location = "us-central1"
  name     = "%s"
}

resource "google_kms_crypto_key" "crypto_key" {
  key_ring = google_kms_key_ring.key_ring.id
  name     = "%s"
}

resource "google_kms_crypto_key_iam_binding" "foo" {
  crypto_key_id = google_kms_crypto_key.crypto_key.id
  role          = "%s"
  members = [
    "serviceAccount:${google_service_account.test_account.email}",
    "serviceAccount:${google_service_account.test_account_2.email}",
  ]
}
`, projectId, orgId, billingAccount, account, account, keyRingName, cryptoKeyName, roleId)
}

func testAccKmsCryptoKeyIamMember_basic(projectId, orgId, billingAccount, account, keyRingName, cryptoKeyName, roleId string) string {
	return fmt.Sprintf(`
resource "google_project" "test_project" {
  name            = "Test project"
  project_id      = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_service" "kms" {
  project = google_project.test_project.project_id
  service = "cloudkms.googleapis.com"
}

resource "google_project_service" "iam" {
  project = google_project_service.kms.project
  service = "iam.googleapis.com"
}

resource "google_service_account" "test_account" {
  project      = google_project_service.iam.project
  account_id   = "%s"
  display_name = "Kms Crypto Key Iam Testing Account"
}

resource "google_kms_key_ring" "key_ring" {
  project  = google_project_service.iam.project
  location = "us-central1"
  name     = "%s"
}

resource "google_kms_crypto_key" "crypto_key" {
  key_ring = google_kms_key_ring.key_ring.id
  name     = "%s"
}

resource "google_kms_crypto_key_iam_member" "foo" {
  crypto_key_id = google_kms_crypto_key.crypto_key.id
  role          = "%s"
  member        = "serviceAccount:${google_service_account.test_account.email}"
}
`, projectId, orgId, billingAccount, account, keyRingName, cryptoKeyName, roleId)
}
