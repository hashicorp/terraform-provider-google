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

func TestAccGeminiLoggingSetting_geminiLoggingSettingBasicExample_update(t *testing.T) {
	t.Parallel()
	context := map[string]interface{}{
		"setting_id": fmt.Sprintf("tf-test-ls-%s", acctest.RandString(t, 10)),
	}
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiLoggingSetting_geminiLoggingSettingBasicExample_basic(context),
			},
			{
				ResourceName:            "google_gemini_logging_setting.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "logging_setting_id", "terraform_labels"},
			},
			{
				Config: testAccGeminiLoggingSetting_geminiLoggingSettingBasicExample_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_gemini_logging_setting.example", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_gemini_logging_setting.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "logging_setting_id", "terraform_labels"},
			},
		},
	})
}
func testAccGeminiLoggingSetting_geminiLoggingSettingBasicExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gemini_logging_setting" "example" {
    logging_setting_id = "%{setting_id}"
    location = "global"
    log_prompts_and_responses = true
    log_metadata = true
}
`, context)
}
func testAccGeminiLoggingSetting_geminiLoggingSettingBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gemini_logging_setting" "example" {
    logging_setting_id = "%{setting_id}"
    location = "global"
	labels = {"my_key" = "my_value"}
	log_prompts_and_responses = false
	log_metadata = false
}
`, context)
}
