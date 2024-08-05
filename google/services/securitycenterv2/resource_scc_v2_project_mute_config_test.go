// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package securitycenterv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccSecurityCenterV2ProjectMuteConfig_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterV2ProjectMuteConfig_basic(context),
			},
			{
				ResourceName:      "google_scc_v2_project_mute_config.default",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"project", "location", "mute_config_id",
				},
			},
			{
				Config: testAccSecurityCenterV2ProjectMuteConfig_update(context),
			},
			{
				ResourceName:      "google_scc_v2_project_mute_config.default",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"project", "location", "mute_config_id",
				},
			},
		},
	})
}

func testAccSecurityCenterV2ProjectMuteConfig_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_scc_v2_project_mute_config" "default" {
  mute_config_id = "tf-test-config-%{random_suffix}"
  project   = "%{project}"
  location       = "global"
  description    = "A test project mute config"
  filter         = "severity = \"HIGH\""
  type           = "STATIC"
}
`, context)
}

func testAccSecurityCenterV2ProjectMuteConfig_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_scc_v2_project_mute_config" "default" {
  mute_config_id = "tf-test-config-%{random_suffix}"
  project   = "%{project}"
  location       = "global"
  description    = "An updated test project mute config"
  filter         = "severity = \"HIGH\""
  type           = "STATIC"
}
`, context)
}
