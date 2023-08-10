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
// and all refer to the same organization, they need to be run serially
func TestAccAccessApprovalSettings(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		"folder":       testAccAccessApprovalFolderSettings,
		"project":      testAccAccessApprovalProjectSettings,
		"organization": testAccAccessApprovalOrganizationSettings,
	}

	for name, tc := range testCases {
		// shadow the tc variable into scope so that when
		// the loop continues, if t.Run hasn't executed tc(t)
		// yet, we don't have a race condition
		// see https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc(t)
		})
	}
}

func testAccAccessApprovalOrganizationSettings(t *testing.T) {
	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"location":      envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAccessApprovalOrganizationSettingsDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAccessApprovalOrganizationSettings_full(context),
			},
			{
				ResourceName:            "google_organization_access_approval_settings.organization_access_approval",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organization_id"},
			},
			{
				Config: testAccAccessApprovalOrganizationSettings_update(context),
			},
			{
				ResourceName:            "google_organization_access_approval_settings.organization_access_approval",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organization_id"},
			},
			{
				Config: testAccAccessApprovalOrganizationSettings_activeKeyVersion(context),
			},
			{
				ResourceName:            "google_organization_access_approval_settings.organization_access_approval",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organization_id"},
			},
		},
	})
}

func testAccAccessApprovalOrganizationSettings_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_organization_access_approval_settings" "organization_access_approval" {
  organization_id     = "%{org_id}"
  notification_emails = ["testuser@example.com"]

  enrolled_services {
    cloud_product = "App Engine"
  }

  enrolled_services {
    cloud_product = "dataflow.googleapis.com"
    enrollment_level = "BLOCK_ALL"
  }
}
`, context)
}

func testAccAccessApprovalOrganizationSettings_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_organization_access_approval_settings" "organization_access_approval" {
  organization_id     = "%{org_id}"
  notification_emails = ["testuser@example.com", "example.user@example.com"]

  enrolled_services {
    cloud_product = "all"
    enrollment_level = "BLOCK_ALL"
  }
}
`, context)
}

func testAccAccessApprovalOrganizationSettings_activeKeyVersion(context map[string]interface{}) string {
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

data "google_access_approval_organization_service_account" "aa_account" {
  organization_id = "%{org_id}"
}

resource "google_kms_crypto_key_iam_member" "iam" {
  crypto_key_id = google_kms_crypto_key.crypto_key.id
  role          = "roles/cloudkms.signerVerifier"
  member        = "serviceAccount:${data.google_access_approval_organization_service_account.aa_account.account_email}"
}

data "google_kms_crypto_key_version" "crypto_key_version" {
  crypto_key = google_kms_crypto_key.crypto_key.id
}

resource "google_organization_access_approval_settings" "organization_access_approval" {
  organization_id     = "%{org_id}"

  enrolled_services {
    cloud_product = "all"
  }

  active_key_version = data.google_kms_crypto_key_version.crypto_key_version.name

  depends_on = [google_kms_crypto_key_iam_member.iam]
}
`, context)
}

func testAccCheckAccessApprovalOrganizationSettingsDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_organization_access_approval_settings" {
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
