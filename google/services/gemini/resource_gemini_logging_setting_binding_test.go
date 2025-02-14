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

func TestAccGeminiLoggingSettingBinding_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"logging_setting_id": fmt.Sprintf("tf-test-ls-%s", acctest.RandString(t, 10)),
		"setting_binding_id": fmt.Sprintf("tf-test-lsb-%s", acctest.RandString(t, 10)),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiLoggingSettingBinding_basic(context),
			},
			{
				ResourceName:            "google_gemini_logging_setting_binding.basic_binding",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "logging_setting_id", "terraform_labels"},
			},
			{
				Config: testAccGeminiLoggingSettingBinding_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_gemini_logging_setting_binding.basic_binding", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_gemini_logging_setting_binding.basic_binding",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "logging_setting_id", "terraform_labels"},
			},
		},
	})
}

func testAccGeminiLoggingSettingBinding_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_gemini_logging_setting" "basic" {
    logging_setting_id = "%{logging_setting_id}"
	location = "global"
	labels = {"my_key" = "my_value"}
    log_prompts_and_responses = true
}

resource "google_gemini_logging_setting_binding" "basic_binding" {
    logging_setting_id = google_gemini_logging_setting.basic.logging_setting_id
    setting_binding_id = "%{setting_binding_id}"
	location = "global"
	target = "projects/${data.google_project.project.number}"
}
`, context)
}

func testAccGeminiLoggingSettingBinding_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_gemini_logging_setting" "basic" {
    logging_setting_id = "%{logging_setting_id}"
	location = "global"
	labels = {"my_key" = "my_value"}
    log_prompts_and_responses = true
}

resource "google_gemini_logging_setting_binding" "basic_binding" {
    logging_setting_id = google_gemini_logging_setting.basic.logging_setting_id
    setting_binding_id = "%{setting_binding_id}"
	location = "global"
	target = "projects/${data.google_project.project.number}"
	labels = {"my_key" = "my_value"}
	product = "GEMINI_CODE_ASSIST"
}
`, context)
}
