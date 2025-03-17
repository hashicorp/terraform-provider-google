// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storagecontrol_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceGoogleStorageControlProjectIntelligenceConfig_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"project":       envvar.GetTestProjectFromEnv(),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleStorageControlProjectIntelligenceConfig_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_storage_control_project_intelligence_config.project_storage_intelligence", "google_storage_control_project_intelligence_config.project_storage_intelligence"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleStorageControlProjectIntelligenceConfig_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_control_project_intelligence_config" "project_storage_intelligence" {
  name = "%{project}"
  edition_config = "STANDARD"
}

data "google_storage_control_project_intelligence_config" "project_storage_intelligence" {
  name = google_storage_control_project_intelligence_config.project_storage_intelligence.name
}
`, context)
}
