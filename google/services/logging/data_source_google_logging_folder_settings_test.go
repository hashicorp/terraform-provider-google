// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package logging_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccLoggingFolderSettings_datasource(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"folder_name": "tf-test-" + acctest.RandString(t, 10),
		"org_id":      envvar.GetTestOrgFromEnv(t),
	}
	resourceName := "data.google_logging_folder_settings.settings"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingFolderSettings_datasource(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "kms_service_account_id"),
					resource.TestCheckResourceAttrSet(resourceName, "logging_service_account_id"),
				),
			},
		},
	})
}

func testAccLoggingFolderSettings_datasource(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_folder" "default" {
	display_name = "%{folder_name}"
	parent       = "organizations/%{org_id}"
}

data "google_logging_folder_settings" "settings" {
	folder = google_folder.default.folder_id
}
`, context)
}
