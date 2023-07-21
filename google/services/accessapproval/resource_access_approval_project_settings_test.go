// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package accessapproval_test

import (
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Since access approval settings are hierarchical, and only one can exist per folder/project/org,
// and all refer to the same organization, they need to be run serially.
// See AccessApprovalOrganizationSettings for the test runner.
func testAccAccessApprovalProjectSettings(t *testing.T) {
	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"location":      envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAccessApprovalProjectSettingsDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAccessApprovalProjectSettings_basic(context),
			},
			{
				ResourceName:            "google_project_access_approval_settings.project_access_approval",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id"},
			},
			{
				Config: testAccAccessApprovalProjectSettings_update(context),
			},
			{
				ResourceName:            "google_project_access_approval_settings.project_access_approval",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id"},
			},
			{
				Config: testAccAccessApprovalProjectSettings_activeKeyVersion(context),
			},
			{
				ResourceName:            "google_project_access_approval_settings.project_access_approval",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id"},
			},
		},
	})
}

func testAccAccessApprovalProjectSettings_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project_access_approval_settings" "project_access_approval" {
  project_id          = "%{project}"

  enrolled_services {
    cloud_product = "all"
    enrollment_level = "BLOCK_ALL"
  }
}
`, context)
}

func testAccAccessApprovalProjectSettings_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project_access_approval_settings" "project_access_approval" {
  project_id          = "%{project}"
  notification_emails = ["testuser@example.com", "example.user@example.com"]

  enrolled_services {
    cloud_product = "all"
    enrollment_level = "BLOCK_ALL"
  }
}
`, context)
}

func testAccAccessApprovalProjectSettings_activeKeyVersion(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_kms_key_ring" "key_ring" {
  name     = "tf-test-%{random_suffix}"
  project  = "%{project}"
  location = "%{location}"
}

resource "google_kms_crypto_key" "crypto_key" {
  name            = "tf-test-%{random_suffix}"
  key_ring        = google_kms_key_ring.key_ring.id
  purpose         = "ASYMMETRIC_SIGN"

  version_template {
    algorithm = "EC_SIGN_P384_SHA384"
  }
}

data "google_access_approval_project_service_account" "aa_account" {
  project_id = "%{project}"
}

resource "google_kms_crypto_key_iam_member" "iam" {
  crypto_key_id = google_kms_crypto_key.crypto_key.id
  role          = "roles/cloudkms.signerVerifier"
  member        = "serviceAccount:${data.google_access_approval_project_service_account.aa_account.account_email}"
}

data "google_kms_crypto_key_version" "crypto_key_version" {
  crypto_key = google_kms_crypto_key.crypto_key.id
}

resource "google_project_access_approval_settings" "project_access_approval" {
  project_id          = "%{project}"

  enrolled_services {
    cloud_product = "all"
  }

  active_key_version = data.google_kms_crypto_key_version.crypto_key_version.name

  depends_on = [google_kms_crypto_key_iam_member.iam]
}
`, context)
}

func testAccCheckAccessApprovalProjectSettingsDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_project_access_approval_settings" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			log.Printf("[DEBUG] Ignoring destroy during test")
		}

		return nil
	}
}
