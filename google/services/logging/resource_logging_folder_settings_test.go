// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package logging_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccLoggingFolderSettings_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
		"original_key":  acctest.BootstrapKMSKeyInLocation(t, "us-central1").CryptoKey.Name,
		"updated_key":   acctest.BootstrapKMSKeyInLocation(t, "us-east1").CryptoKey.Name,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingFolderSettings_onlyRequired(context),
			},
			{
				ResourceName:            "google_logging_folder_settings.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"folder"},
			},
			{
				Config: testAccLoggingFolderSettings_full(context),
			},
			{
				ResourceName:            "google_logging_folder_settings.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"folder"},
			},
		},
	})
}

func testAccLoggingFolderSettings_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_folder" "my_folder" {
  display_name = "tf-test-folder-%{random_suffix}"
  parent       = "organizations/%{org_id}"
  deletion_protection = false
}

data "google_logging_folder_settings" "settings" {
  folder = google_folder.my_folder.folder_id
}

resource "google_kms_crypto_key_iam_member" "iam_folder" {
  crypto_key_id = "%{original_key}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:${data.google_logging_folder_settings.settings.kms_service_account_id}"
}

data "google_logging_organization_settings" "settings" {
  organization = "%{org_id}"
}

resource "google_kms_crypto_key_iam_member" "iam_org" {
  crypto_key_id = "%{original_key}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:${data.google_logging_organization_settings.settings.kms_service_account_id}"
}

resource "google_logging_folder_settings" "example" {
  disable_default_sink = true
  folder               = google_folder.my_folder.folder_id
  kms_key_name         = "%{original_key}"
  storage_location     = "us-central1"
  depends_on   = [
    google_kms_crypto_key_iam_member.iam_folder,
	google_kms_crypto_key_iam_member.iam_org
  ]
}
`, context)
}

func testAccLoggingFolderSettings_onlyRequired(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_folder" "my_folder" {
  display_name = "tf-test-folder-%{random_suffix}"
  parent       = "organizations/%{org_id}"
  deletion_protection = false
}

data "google_logging_folder_settings" "settings" {
  folder = google_folder.my_folder.folder_id
}

resource "google_kms_crypto_key_iam_member" "iam_folder" {
  crypto_key_id = "%{original_key}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:${data.google_logging_folder_settings.settings.kms_service_account_id}"
}

data "google_logging_organization_settings" "settings" {
  organization = "%{org_id}"
}

resource "google_kms_crypto_key_iam_member" "iam_org" {
  crypto_key_id = "%{original_key}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:${data.google_logging_organization_settings.settings.kms_service_account_id}"
}

resource "google_logging_folder_settings" "example" {
  folder       = google_folder.my_folder.folder_id
  depends_on   = [
    google_kms_crypto_key_iam_member.iam_folder,
	google_kms_crypto_key_iam_member.iam_org
  ]
}
`, context)
}
