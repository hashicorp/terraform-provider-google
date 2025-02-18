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

func TestAccGeminiDataSharingWithGoogleSetting_geminiDataSharingWithGoogleSettingBasicExample_update(t *testing.T) {
	t.Parallel()
	context := map[string]interface{}{
		"setting_id": fmt.Sprintf("tf-test-ls-%s", acctest.RandString(t, 10)),
	}
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDataSharingWithGoogleSetting_geminiDataSharingWithGoogleSettingBasicExample_basic(context),
			},
			{
				ResourceName:            "google_gemini_data_sharing_with_google_setting.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "data_sharing_with_google_setting_id", "terraform_labels"},
			},
			{
				Config: testAccGeminiDataSharingWithGoogleSetting_geminiDataSharingWithGoogleSettingBasicExample_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_gemini_data_sharing_with_google_setting.example", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_gemini_data_sharing_with_google_setting.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "data_sharing_with_google_setting_id", "terraform_labels"},
			},
		},
	})
}
func testAccGeminiDataSharingWithGoogleSetting_geminiDataSharingWithGoogleSettingBasicExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gemini_data_sharing_with_google_setting" "example" {
    data_sharing_with_google_setting_id = "%{setting_id}"
    location = "global"
    enable_preview_data_sharing = true
}
`, context)
}
func testAccGeminiDataSharingWithGoogleSetting_geminiDataSharingWithGoogleSettingBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gemini_data_sharing_with_google_setting" "example" {
    data_sharing_with_google_setting_id = "%{setting_id}"
    location = "global"
    labels = {"my_key" = "my_value"}
    enable_preview_data_sharing = false
}
`, context)
}
