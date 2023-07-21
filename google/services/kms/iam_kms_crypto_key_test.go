// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package kms_test

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/kms"
)

func TestAccKmsCryptoKeyIamBinding(t *testing.T) {
	t.Parallel()

	orgId := envvar.GetTestOrgFromEnv(t)
	projectId := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	billingAccount := envvar.GetTestBillingAccountFromEnv(t)
	account := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	roleId := "roles/cloudkms.cryptoKeyDecrypter"
	keyRingName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	keyRingId := &kms.KmsKeyRingId{
		Project:  projectId,
		Location: DEFAULT_KMS_TEST_LOCATION,
		Name:     keyRingName,
	}
	cryptoKeyName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test Iam Binding creation
				Config: testAccKmsCryptoKeyIamBinding_basic(projectId, orgId, billingAccount, account, keyRingName, cryptoKeyName, roleId),
				Check: testAccCheckGoogleKmsCryptoKeyIamBindingExists(t, "foo", roleId, []string{
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, projectId),
				}),
			},
			{
				ResourceName:      "google_kms_crypto_key_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("%s/%s %s", keyRingId.TerraformId(), cryptoKeyName, roleId),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test Iam Binding update
				Config: testAccKmsCryptoKeyIamBinding_update(projectId, orgId, billingAccount, account, keyRingName, cryptoKeyName, roleId),
				Check: testAccCheckGoogleKmsCryptoKeyIamBindingExists(t, "foo", roleId, []string{
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, projectId),
					fmt.Sprintf("serviceAccount:%s-2@%s.iam.gserviceaccount.com", account, projectId),
				}),
			},
			{
				ResourceName:      "google_kms_crypto_key_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("%s/%s %s", keyRingId.TerraformId(), cryptoKeyName, roleId),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccKmsCryptoKeyIamMember(t *testing.T) {
	t.Parallel()

	orgId := envvar.GetTestOrgFromEnv(t)
	projectId := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	billingAccount := envvar.GetTestBillingAccountFromEnv(t)
	account := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	roleId := "roles/cloudkms.cryptoKeyEncrypter"
	keyRingName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	keyRingId := &kms.KmsKeyRingId{
		Project:  projectId,
		Location: DEFAULT_KMS_TEST_LOCATION,
		Name:     keyRingName,
	}
	cryptoKeyName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccKmsCryptoKeyIamMember_basic(projectId, orgId, billingAccount, account, keyRingName, cryptoKeyName, roleId),
				Check: testAccCheckGoogleKmsCryptoKeyIamMemberExists(t, "foo", roleId,
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, projectId),
				),
			},
			{
				ResourceName:      "google_kms_crypto_key_iam_member.foo",
				ImportStateId:     fmt.Sprintf("%s/%s %s serviceAccount:%s@%s.iam.gserviceaccount.com", keyRingId.TerraformId(), cryptoKeyName, roleId, account, projectId),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccKmsCryptoKeyIamPolicy(t *testing.T) {
	t.Parallel()

	orgId := envvar.GetTestOrgFromEnv(t)
	projectId := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	billingAccount := envvar.GetTestBillingAccountFromEnv(t)
	account := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	roleId := "roles/cloudkms.cryptoKeyEncrypter"
	keyRingName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	keyRingId := &kms.KmsKeyRingId{
		Project:  projectId,
		Location: DEFAULT_KMS_TEST_LOCATION,
		Name:     keyRingName,
	}
	cryptoKeyName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccKmsCryptoKeyIamPolicy_basic(projectId, orgId, billingAccount, account, keyRingName, cryptoKeyName, roleId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleCryptoKmsKeyIam(t, "foo", roleId, []string{
						fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, projectId),
					}),
					resource.TestCheckResourceAttrSet("data.google_kms_crypto_key_iam_policy.foo", "policy_data"),
				),
			},
			{
				ResourceName:      "google_kms_crypto_key_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("%s/%s", keyRingId.TerraformId(), cryptoKeyName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGoogleKmsCryptoKeyIamBindingExists(t *testing.T, bindingResourceName, roleId string, members []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		bindingRs, ok := s.RootModule().Resources[fmt.Sprintf("google_kms_crypto_key_iam_binding.%s", bindingResourceName)]
		if !ok {
			return fmt.Errorf("Not found: %s", bindingResourceName)
		}

		config := acctest.GoogleProviderConfig(t)
		cryptoKeyId, err := kms.ParseKmsCryptoKeyId(bindingRs.Primary.Attributes["crypto_key_id"], config)

		if err != nil {
			return err
		}

		p, err := config.NewKmsClient(config.UserAgent).Projects.Locations.KeyRings.CryptoKeys.GetIamPolicy(cryptoKeyId.CryptoKeyId()).Do()
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

func testAccCheckGoogleKmsCryptoKeyIamMemberExists(t *testing.T, n, role, member string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["google_kms_crypto_key_iam_member."+n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		config := acctest.GoogleProviderConfig(t)
		cryptoKeyId, err := kms.ParseKmsCryptoKeyId(rs.Primary.Attributes["crypto_key_id"], config)

		if err != nil {
			return err
		}

		p, err := config.NewKmsClient(config.UserAgent).Projects.Locations.KeyRings.GetIamPolicy(cryptoKeyId.CryptoKeyId()).Do()
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

func testAccCheckGoogleCryptoKmsKeyIam(t *testing.T, n, role string, members []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["google_kms_crypto_key_iam_policy."+n]
		if !ok {
			return fmt.Errorf("IAM policy resource not found")
		}

		config := acctest.GoogleProviderConfig(t)
		cryptoKeyId, err := kms.ParseKmsCryptoKeyId(rs.Primary.Attributes["crypto_key_id"], config)

		if err != nil {
			return err
		}

		p, err := config.NewKmsClient(config.UserAgent).Projects.Locations.KeyRings.GetIamPolicy(cryptoKeyId.CryptoKeyId()).Do()
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
			} else {
				return fmt.Errorf("Binding found but not expected for role: %v", binding.Role)
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

func testAccKmsCryptoKeyIamPolicy_basic(projectId, orgId, billingAccount, account, keyRingName, cryptoKeyName, roleId string) string {
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

data "google_iam_policy" "foo" {
  binding {
    role = "%s"
    members = ["serviceAccount:${google_service_account.test_account.email}"]
  }
}

resource "google_kms_crypto_key_iam_policy" "foo" {
  crypto_key_id = google_kms_crypto_key.crypto_key.id
  policy_data = data.google_iam_policy.foo.policy_data
}

data "google_kms_crypto_key_iam_policy" "foo" {
  crypto_key_id = google_kms_crypto_key.crypto_key.id
}
`, projectId, orgId, billingAccount, account, keyRingName, cryptoKeyName, roleId)
}
