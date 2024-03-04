// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package firebaseappcheck_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccFirebaseAppCheckDebugToken_firebaseAppCheckDebugTokenUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_id":   envvar.GetTestProjectFromEnv(),
		"display_name": "Debug Token 1",
		"token":        "5E728315-E121-467F-BCA1-1FE71130BB98",
	}

	contextUpdated := map[string]interface{}{
		"project_id":   envvar.GetTestProjectFromEnv(),
		"display_name": "Debug Token 2",
		"token":        "5E728315-E121-467F-BCA1-1FE71130BB98",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderBetaFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
			"time":   {},
		},
		CheckDestroy: testAccCheckFirebaseAppCheckDebugTokenDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFirebaseAppCheckDebugToken_firebaseAppCheckDebugTokenTemplate(context),
			},
			{
				ResourceName:            "google_firebase_app_check_debug_token.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"token", "app_id"},
			},
			{
				Config: testAccFirebaseAppCheckDebugToken_firebaseAppCheckDebugTokenTemplate(contextUpdated),
			},
			{
				ResourceName:            "google_firebase_app_check_debug_token.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"token", "app_id"},
			},
		},
	})
}

func testAccFirebaseAppCheckDebugToken_firebaseAppCheckDebugTokenTemplate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_firebase_web_app" "default" {
  provider = google-beta

  project = "%{project_id}"
  display_name = "Web App for debug token"
}

# It takes a while for App Check to recognize the new app
# If your app already exists, you don't have to wait 30 seconds.
resource "time_sleep" "wait_30s" {
  depends_on = [google_firebase_web_app.default]
  create_duration = "30s"
}

resource "google_firebase_app_check_debug_token" "default" {
  provider = google-beta

  project      = "%{project_id}"
  app_id       = google_firebase_web_app.default.app_id
  display_name = "%{display_name}"
  token        = "%{token}"

  depends_on = [time_sleep.wait_30s]
}
`, context)
}
