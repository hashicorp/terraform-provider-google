// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package logging_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccLoggingProjectSettings_datasource(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":    "tf-test-" + acctest.RandString(t, 10),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
	}
	resourceName := "data.google_logging_project_settings.settings"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingProjectSettings_datasource(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "kms_service_account_id"),
					resource.TestCheckResourceAttrSet(resourceName, "logging_service_account_id"),
				),
			},
		},
	})
}

func testAccLoggingProjectSettings_datasource(context map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_project" "default" {
		project_id      = "%{project_name}"
		name            = "%{project_name}"
		org_id          = "%{org_id}"
		billing_account = "%{billing_account}"
	}
	
	resource "google_project_service" "logging_service" {
		project = google_project.default.project_id
		service = "logging.googleapis.com"
	}

	data "google_logging_project_settings" "settings" {
		project = google_project_service.logging_service.project
	}
`, context)
}
