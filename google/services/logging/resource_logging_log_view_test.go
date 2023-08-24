// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package logging_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccLoggingLogView_loggingLogViewBasicExampleUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingLogViewDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingLogView_loggingLogViewBasicExample(context),
			},
			{
				ResourceName:            "google_logging_log_view.logging_log_view",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "location", "bucket"},
			},
			{
				Config: testAccLoggingLogView_loggingLogViewBasicExampleUpdate(context),
			},
			{
				ResourceName:            "google_logging_log_view.logging_log_view",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "location", "bucket"},
			},
		},
	})
}

func testAccLoggingLogView_loggingLogViewBasicExampleUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_logging_project_bucket_config" "logging_log_view" {
    project        = "%{project}"
    location       = "global"
    retention_days = 30
    bucket_id      = "_Default"
}

resource "google_logging_log_view" "logging_log_view" {
  name        = "tf-test-view%{random_suffix}"
  bucket      = google_logging_project_bucket_config.logging_log_view.id
  description = "An updated logging view configured with Terraform"
  filter      = "SOURCE(\"projects/myproject\") AND resource.type = \"gce_instance\""
}
`, context)
}
