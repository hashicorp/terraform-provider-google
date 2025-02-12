// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package gemini_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccGeminiReleaseChannelSetting_geminiReleaseChannelSettingBasicExample_update(t *testing.T) {
	t.Parallel()
	context := map[string]interface{}{
		"setting_id": fmt.Sprintf("tf-test-ls-%s", acctest.RandString(t, 10)),
	}
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiReleaseChannelSetting_geminiReleaseChannelSettingBasicExample_basic(context),
			},
			{
				ResourceName:            "google_gemini_release_channel_setting.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "release_channel_setting_id", "terraform_labels"},
			},
			{
				Config: testAccGeminiReleaseChannelSetting_geminiReleaseChannelSettingBasicExample_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_gemini_release_channel_setting.example", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_gemini_release_channel_setting.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "release_channel_setting_id", "terraform_labels"},
			},
		},
	})
}
func testAccGeminiReleaseChannelSetting_geminiReleaseChannelSettingBasicExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gemini_release_channel_setting" "example" {
    release_channel_setting_id = "%{setting_id}"
    location = "global"
    release_channel = "EXPERIMENTAL"
}
`, context)
}
func testAccGeminiReleaseChannelSetting_geminiReleaseChannelSettingBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gemini_release_channel_setting" "example" {
    release_channel_setting_id = "%{setting_id}"
    location = "global"
    labels = {"my_key" = "my_value"}
    release_channel = "STABLE"
}
`, context)
}
