// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package logging_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccLoggingOrganizationSettings_update(t *testing.T) {
	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgTargetFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
		"original_key":  acctest.BootstrapKMSKeyInLocation(t, "us-central1").CryptoKey.Name,
		"updated_key":   acctest.BootstrapKMSKeyInLocation(t, "us-east1").CryptoKey.Name,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingOrganizationSettings_onlyRequired(context),
			},
			{
				ResourceName:            "google_logging_organization_settings.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organization"},
			},
			{
				Config: testAccLoggingOrganizationSettings_full(context),
			},
			{
				ResourceName:            "google_logging_organization_settings.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organization"},
			},
		},
	})
}

func testAccLoggingOrganizationSettings_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_logging_organization_settings" "example" {
  disable_default_sink = false
  kms_key_name         = "%{original_key}"
  organization         = "%{org_id}"
  storage_location     = "us-central1"
  depends_on           = [ google_kms_crypto_key_iam_member.iam ]
}

data "google_logging_organization_settings" "settings" {
  organization = "%{org_id}"
}

resource "google_kms_crypto_key_iam_member" "iam" {
  crypto_key_id = "%{original_key}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:${data.google_logging_organization_settings.settings.kms_service_account_id}"
}
`, context)
}

func testAccLoggingOrganizationSettings_onlyRequired(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_logging_organization_settings" "example" {
  organization = "%{org_id}"
}
`, context)
}
