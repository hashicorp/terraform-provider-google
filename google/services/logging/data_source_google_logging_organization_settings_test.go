// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package logging_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccLoggingOrganizationSettings_datasource(t *testing.T) {
	// google_logging_organization_settings is a singleton, and multiple tests mutate it.
	orgSettingsMu.Lock()
	t.Cleanup(orgSettingsMu.Unlock)

	context := map[string]interface{}{
		"org_id":       envvar.GetTestOrgFromEnv(t),
		"original_key": acctest.BootstrapKMSKeyInLocation(t, "us-central1").CryptoKey.Name,
	}
	resourceName := "data.google_logging_organization_settings.settings"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingOrganizationSettings_full(context),
			},
			{
				Config: testAccLoggingOrganizationSettings_datasource(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "kms_service_account_id"),
					resource.TestCheckResourceAttrSet(resourceName, "logging_service_account_id"),
				),
			},
		},
	})
}

func testAccLoggingOrganizationSettings_datasource(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_logging_organization_settings" "settings" {
	organization = "%{org_id}"
}
`, context)
}
