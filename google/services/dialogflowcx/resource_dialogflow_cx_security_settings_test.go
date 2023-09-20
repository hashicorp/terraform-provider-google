// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dialogflowcx_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDialogflowCXSecuritySettings_dialogflowcxSecuritySettings_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":         envvar.GetTestProjectFromEnv(),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDialogflowCXSecuritySettings_dialogflowcxSecuritySettings_basic(context),
			},
			{
				ResourceName:            "google_dialogflow_cx_security_settings.basic_security_settings",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
			{
				Config: testAccDialogflowCXSecuritySettings_dialogflowcxSecuritySettings_update(context),
			},
			{
				ResourceName:            "google_dialogflow_cx_security_settings.basic_security_settings",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func testAccDialogflowCXSecuritySettings_dialogflowcxSecuritySettings_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dialogflow_cx_security_settings" "basic_security_settings" {
  display_name          = "tf-test-dialogflowcx-security-settings%{random_suffix}"
  location              = "global"
  purge_data_types      = []
  retention_window_days = 7
}
`, context)
}

func testAccDialogflowCXSecuritySettings_dialogflowcxSecuritySettings_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_inspect_template" "inspect" {
  parent       = "projects/%{project}/locations/global"
  display_name = "tf-test-dialogflowcx-inspect-template%{random_suffix}"
  inspect_config {
    info_types {
      name = "EMAIL_ADDRESS"
    }
  }
}

resource "google_data_loss_prevention_deidentify_template" "deidentify" {
  parent       = "projects/%{project}/locations/global"
  display_name = "tf-test-dialogflowcx-deidentify-template%{random_suffix}"
  deidentify_config {
    info_type_transformations {
      transformations {
        primitive_transformation {
          replace_config {
            new_value {
              string_value = "[REDACTED]"
            }
          }
        }
      }
    }
  }
}

resource "google_storage_bucket" "bucket" {
  name                        = "tf-test-dialogflowcx-bucket%{random_suffix}"
  location                    = "US"
  uniform_bucket_level_access = true
}

resource "google_dialogflow_cx_security_settings" "basic_security_settings" {
  display_name        = "tf-test-dialogflowcx-security-settings%{random_suffix}"
  location            = "global"
  redaction_strategy  = "REDACT_WITH_SERVICE"
  redaction_scope     = "REDACT_DISK_STORAGE"
  inspect_template    = google_data_loss_prevention_inspect_template.inspect.id
  deidentify_template = google_data_loss_prevention_deidentify_template.deidentify.id
  purge_data_types    = ["DIALOGFLOW_HISTORY"]
  audio_export_settings {
    gcs_bucket             = google_storage_bucket.bucket.id
    audio_export_pattern   = "export"
    enable_audio_redaction = true
    audio_format           = "OGG"
  }
  insights_export_settings {
    enable_insights_export = true
  }
  retention_strategy = "REMOVE_AFTER_CONVERSATION"
}
`, context)
}
